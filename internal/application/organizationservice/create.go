package organizationservice

import (
	"context"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type CreateOrganizationRequest struct {
	Name string `validate:"required,min=2,max=255"`
}

type CreateOrganizationService interface {
	CreateOrganization(ctx context.Context, req *CreateOrganizationRequest) (*domain.Organization, error)
}

type createOrganizationService struct {
	organizationRepository domain.OrganizationRepository
	auther                 authzservice.AuthZ
	t                      txmanager.TxManager
}

var _ CreateOrganizationService = (*createOrganizationService)(nil)

func NewCreateOrganizationService(
	auther authzservice.AuthZ,
	organizationRepository domain.OrganizationRepository,
	txManager txmanager.TxManager,
) CreateOrganizationService {
	return &createOrganizationService{
		organizationRepository: organizationRepository,
		auther:                 auther,
		t:                      txManager,
	}
}

func (s *createOrganizationService) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest) (*domain.Organization, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("CreateOrganizationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	organization := &domain.Organization{
		ID:   uuid.New(),
		Name: req.Name,
	}

	err = s.t.WithinTx(ctx, func(ctx context.Context) error {
		return s.organizationRepository.Create(ctx, organization)
	})
	if err != nil {
		return nil, err
	}

	return organization, nil
}

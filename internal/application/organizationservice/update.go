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

type UpdateOrganizationRequest struct {
	Name string `validate:"required,min=2,max=255"`
}

type UpdateOrganizationService interface {
	UpdateOrganization(ctx context.Context, id uuid.UUID, req *UpdateOrganizationRequest) (*domain.Organization, error)
}

type updateOrganizationService struct {
	organizationRepository domain.OrganizationRepository
	auther                 authzservice.AuthZ
	t                      txmanager.TxManager
}

var _ UpdateOrganizationService = (*updateOrganizationService)(nil)

func NewUpdateOrganizationService(
	auther authzservice.AuthZ,
	organizationRepository domain.OrganizationRepository,
	txManager txmanager.TxManager,
) UpdateOrganizationService {
	return &updateOrganizationService{
		organizationRepository: organizationRepository,
		auther:                 auther,
		t:                      txManager,
	}
}

func (s *updateOrganizationService) UpdateOrganization(ctx context.Context, id uuid.UUID, req *UpdateOrganizationRequest) (*domain.Organization, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateOrganizationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	organization, err := s.organizationRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	organization.Name = req.Name

	err = s.t.WithinTx(ctx, func(ctx context.Context) error {
		return s.organizationRepository.Update(ctx, organization)
	})
	if err != nil {
		return nil, err
	}

	return organization, nil
}

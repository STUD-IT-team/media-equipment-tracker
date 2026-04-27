package organizationservice

import (
	"context"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type DeleteOrganizationService interface {
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
}

type deleteOrganizationService struct {
	organizationRepository domain.OrganizationRepository
	auther                 authzservice.AuthZ
	t                      txmanager.TxManager
}

var _ DeleteOrganizationService = (*deleteOrganizationService)(nil)

func NewDeleteOrganizationService(
	auther authzservice.AuthZ,
	organizationRepository domain.OrganizationRepository,
	txManager txmanager.TxManager,
) DeleteOrganizationService {
	return &deleteOrganizationService{
		organizationRepository: organizationRepository,
		auther:                 auther,
		t:                      txManager,
	}
}

func (s *deleteOrganizationService) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	err = s.t.WithinTx(ctx, func(ctx context.Context) error {
		organization, err := s.organizationRepository.Get(ctx, id, domain.WithUsers(), domain.WithEquipmentInvocations(), domain.WithStudioInvocations())
		if err != nil {
			return err
		}

		if len(organization.Users) > 0 {
			return ErrOrganizationHasUsers
		}
		if len(organization.EquipmentInvocations) > 0 || len(organization.StudioInvocations) > 0 {
			return ErrOrganizationHasInvocations
		}

		return s.organizationRepository.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	return nil
}

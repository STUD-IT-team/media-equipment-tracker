package departmentservice

import (
	"context"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type DeleteDepartmentService interface {
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
}

type deleteDepartmentService struct {
	departmentRepository domain.DepartmentRepository
	auther               authzservice.AuthZ
	t                    txmanager.TxManager
}

var _ DeleteDepartmentService = (*deleteDepartmentService)(nil)

func NewDeleteDepartmentService(
	auther authzservice.AuthZ,
	departmentRepository domain.DepartmentRepository,
	txManager txmanager.TxManager,
) DeleteDepartmentService {
	return &deleteDepartmentService{
		departmentRepository: departmentRepository,
		auther:               auther,
		t:                    txManager,
	}
}

func (s *deleteDepartmentService) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	err = s.t.WithinTx(ctx, func(ctx context.Context) error {
		department, err := s.departmentRepository.Get(ctx, id, domain.DepartmentWithUsers(), domain.DepartmentWithEquipmentInvocations(), domain.DepartmentWithStudioInvocations())
		if err != nil {
			return err
		}

		if len(department.Users) > 0 {
			return ErrDepartmentHasUsers
		}
		if len(department.EquipmentInvocations) > 0 || len(department.StudioInvocations) > 0 {
			return ErrDepartmentHasInvocations
		}

		return s.departmentRepository.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	return nil
}

package departmentservice

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

type UpdateDepartmentRequest struct {
	Name string `validate:"required,min=2,max=255"`
}

type UpdateDepartmentService interface {
	UpdateDepartment(ctx context.Context, id uuid.UUID, req *UpdateDepartmentRequest) (*domain.Department, error)
}

type updateDepartmentService struct {
	departmentRepository domain.DepartmentRepository
	auther               authzservice.AuthZ
	t                    txmanager.TxManager
}

var _ UpdateDepartmentService = (*updateDepartmentService)(nil)

func NewUpdateDepartmentService(
	auther authzservice.AuthZ,
	departmentRepository domain.DepartmentRepository,
	txManager txmanager.TxManager,
) UpdateDepartmentService {
	return &updateDepartmentService{
		departmentRepository: departmentRepository,
		auther:               auther,
		t:                    txManager,
	}
}

func (s *updateDepartmentService) UpdateDepartment(ctx context.Context, id uuid.UUID, req *UpdateDepartmentRequest) (*domain.Department, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateDepartmentRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	department, err := s.departmentRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	department.Name = req.Name

	err = s.t.WithinTx(ctx, func(ctx context.Context) error {
		return s.departmentRepository.Update(ctx, department)
	})
	if err != nil {
		return nil, err
	}

	return department, nil
}

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

type CreateDepartmentRequest struct {
	Name string `validate:"required,min=2,max=255"`
}

type CreateDepartmentService interface {
	CreateDepartment(ctx context.Context, req *CreateDepartmentRequest) (*domain.Department, error)
}

type createDepartmentService struct {
	departmentRepository domain.DepartmentRepository
	auther               authzservice.AuthZ
	t                    txmanager.TxManager
}

var _ CreateDepartmentService = (*createDepartmentService)(nil)

func NewCreateDepartmentService(
	auther authzservice.AuthZ,
	departmentRepository domain.DepartmentRepository,
	txManager txmanager.TxManager,
) CreateDepartmentService {
	return &createDepartmentService{
		departmentRepository: departmentRepository,
		auther:               auther,
		t:                    txManager,
	}
}

func (s *createDepartmentService) CreateDepartment(ctx context.Context, req *CreateDepartmentRequest) (*domain.Department, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("CreateDepartmentRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	department := &domain.Department{
		ID:   uuid.New(),
		Name: req.Name,
	}

	err = s.departmentRepository.Create(ctx, department)
	if err != nil {
		return nil, err
	}

	return department, nil
}

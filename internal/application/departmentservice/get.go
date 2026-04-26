package departmentservice

import (
	"context"

	"media-equipment-tracker/internal/domain"

	"github.com/google/uuid"
)

type GetDepartmentService interface {
	GetDepartment(ctx context.Context, id uuid.UUID) (*domain.Department, error)
}

type getDepartmentService struct {
	departmentRepository domain.DepartmentRepository
}

var _ GetDepartmentService = (*getDepartmentService)(nil)

func NewGetDepartmentService(
	departmentRepository domain.DepartmentRepository,
) GetDepartmentService {
	return &getDepartmentService{
		departmentRepository: departmentRepository,
	}
}

func (s *getDepartmentService) GetDepartment(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	return s.departmentRepository.Get(ctx, id, domain.DepartmentWithUsers(), domain.DepartmentWithEquipment())
}

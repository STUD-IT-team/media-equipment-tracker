package departmentservice

import (
	"context"
	"strings"

	"media-equipment-tracker/internal/domain"
)

type ListDepartmentService interface {
	ListDepartments(ctx context.Context, search string) ([]*domain.Department, error)
}

type listDepartmentService struct {
	departmentRepository domain.DepartmentRepository
}

var _ ListDepartmentService = (*listDepartmentService)(nil)

func NewListDepartmentService(
	departmentRepository domain.DepartmentRepository,
) ListDepartmentService {
	return &listDepartmentService{
		departmentRepository: departmentRepository,
	}
}

func (s *listDepartmentService) ListDepartments(ctx context.Context, search string) ([]*domain.Department, error) {
	departments, err := s.departmentRepository.List(ctx, domain.DepartmentWithUsers(), domain.DepartmentWithEquipment())
	if err != nil {
		return nil, err
	}

	// Simple search by name (case insensitive)
	if search != "" {
		filtered := make([]*domain.Department, 0)
		for _, dep := range departments {
			if strings.Contains(strings.ToLower(dep.Name), strings.ToLower(search)) {
				filtered = append(filtered, dep)
			}
		}
		return filtered, nil
	}

	return departments, nil
}

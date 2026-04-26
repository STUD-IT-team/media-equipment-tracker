package departmentservice

import (
	"context"

	"media-equipment-tracker/internal/domain"

	"github.com/google/uuid"
)

type GetDepartmentEquipmentService interface {
	GetDepartmentEquipment(ctx context.Context, id uuid.UUID, availableOnly bool) ([]*domain.Equipment, error)
}

type getDepartmentEquipmentService struct {
	departmentRepository domain.DepartmentRepository
}

var _ GetDepartmentEquipmentService = (*getDepartmentEquipmentService)(nil)

func NewGetDepartmentEquipmentService(
	departmentRepository domain.DepartmentRepository,
) GetDepartmentEquipmentService {
	return &getDepartmentEquipmentService{
		departmentRepository: departmentRepository,
	}
}

func (s *getDepartmentEquipmentService) GetDepartmentEquipment(ctx context.Context, id uuid.UUID, availableOnly bool) ([]*domain.Equipment, error) {
	department, err := s.departmentRepository.Get(ctx, id, domain.DepartmentWithEquipment())
	if err != nil {
		return nil, err
	}

	equipment := department.Equipment
	if availableOnly {
		filtered := make([]*domain.Equipment, 0)
		for _, eq := range equipment {
			if eq.Status == domain.EquipmentStatusAvailable {
				filtered = append(filtered, eq)
			}
		}
		return filtered, nil
	}

	return equipment, nil
}

package equipmentservice

import (
	"context"
	"fmt"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/txmanager"
	"slices"

	"github.com/google/uuid"
)

type DepartmentAssociationService interface {
	UpdateAssociations(ctx context.Context, equipment *domain.Equipment, newDepIds []uuid.UUID) error
}

type departmentAssociationService struct {
	departmentRepository domain.DepartmentRepository
	txManager            txmanager.TxManager
}

func NewDepartmentAssociationService(
	departmentRepository domain.DepartmentRepository,
	txManager txmanager.TxManager,
) DepartmentAssociationService {
	return &departmentAssociationService{
		departmentRepository: departmentRepository,
		txManager:            txManager,
	}
}

func (s *departmentAssociationService) UpdateAssociations(ctx context.Context, equipment *domain.Equipment, newDepIds []uuid.UUID) error {
	if equipment == nil || newDepIds == nil {
		return fmt.Errorf("invalid arguments")
	}

	toDelete := make(map[uuid.UUID]struct{}, len(equipment.Departments))
	for _, eq := range equipment.Departments {
		toDelete[eq.ID] = struct{}{}
	}

	toAdd := make(map[uuid.UUID]struct{})
	for _, newId := range newDepIds {
		if _, ok := toDelete[newId]; !ok {
			toAdd[newId] = struct{}{}
		} else {
			delete(toDelete, newId)
		}
	}

	err := s.txManager.WithinTx(ctx, func(ctx context.Context) error {
		for depId := range toDelete {
			department, err := s.departmentRepository.Get(ctx, depId, domain.DepartmentWithEquipment())
			if err != nil {
				return err
			}
			department.Equipment = slices.DeleteFunc(department.Equipment, func(e *domain.Equipment) bool {
				return e.ID == equipment.ID
			})
			if err := s.departmentRepository.Update(ctx, department); err != nil {
				return err
			}
			equipment.Departments = slices.DeleteFunc(equipment.Departments, func(d *domain.Department) bool {
				return d.ID == depId
			})
		}

		for depId := range toAdd {
			department, err := s.departmentRepository.Get(ctx, depId, domain.DepartmentWithEquipment())
			if err != nil {
				return err
			}
			department.Equipment = append(department.Equipment, equipment)
			if err := s.departmentRepository.Update(ctx, department); err != nil {
				return err
			}
			equipment.Departments = append(equipment.Departments, department)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

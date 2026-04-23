package equipmentservice

import (
	"context"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type EquipmentService interface {
	SearchEquipmentService
	CreateEquipmentService
	UpdateEquipmentService
	AvailabilityEquipmentService
	Get(ctx context.Context, id uuid.UUID) (*domain.Equipment, error)
	GetByInventoryNumber(ctx context.Context, inventoryNumber string) (*domain.Equipment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type equipmentService struct {
	SearchEquipmentService
	CreateEquipmentService
	UpdateEquipmentService
	AvailabilityEquipmentService
	equipmentRepository domain.EquipmentRepository
	auther              authzservice.AuthZ
	t                   txmanager.TxManager
}

var _ EquipmentService = (*equipmentService)(nil)

func NewEquipmentService(
	auther authzservice.AuthZ,
	equipmentRepository domain.EquipmentRepository,
	searchRepository SearchEquipmentRepository,
	searchInvocationRepository invocationservice.SearchInvocationRepository,
	departmentRepository domain.DepartmentRepository,
	txManager txmanager.TxManager,
) EquipmentService {
	return &equipmentService{
		equipmentRepository:          equipmentRepository,
		SearchEquipmentService:       NewSearchEquipmentService(searchRepository),
		CreateEquipmentService:       NewCreateEquipmentService(auther, equipmentRepository, NewDepartmentAssociationService(departmentRepository, txManager), txManager),
		UpdateEquipmentService:       NewUpdateEquipmentService(auther, equipmentRepository, NewDepartmentAssociationService(departmentRepository, txManager), txManager),
		AvailabilityEquipmentService: NewAvailabilityEquipmentService(equipmentRepository, searchInvocationRepository, txManager),
		auther:                       auther,
		t:                            txManager,
	}
}

func (s *equipmentService) Get(ctx context.Context, id uuid.UUID) (*domain.Equipment, error) {
	return s.equipmentRepository.Get(ctx, id, domain.EquipmentWithDepartments(), domain.EquipmentWithCurrentInvocation(), domain.EquipmentWithEquipmentInInvocations())
}

func (s *equipmentService) GetByInventoryNumber(ctx context.Context, inventoryNumber string) (*domain.Equipment, error) {
	return s.equipmentRepository.GetByInventoryNumber(ctx, inventoryNumber, domain.EquipmentWithDepartments(), domain.EquipmentWithCurrentInvocation(), domain.EquipmentWithEquipmentInInvocations())
}

func (s *equipmentService) Delete(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	err = s.t.WithinTx(ctx, func(ctx context.Context) error {
		equipment, err := s.equipmentRepository.Get(ctx, id, domain.EquipmentWithEquipmentInInvocations())
		if err != nil {
			return err
		}
		if len(equipment.Invocations) > 0 {
			return ErrEquipmentHasInvocations
		}

		return s.equipmentRepository.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	return nil
}

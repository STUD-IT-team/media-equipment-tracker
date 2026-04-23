package equipmentservice

import (
	"context"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"
	"slices"

	"github.com/google/uuid"
)

type UpdateEquipmentService interface {
	UpdateEquipment(ctx context.Context, req UpdateEquipmentRequest) (*domain.Equipment, error)
}

type UpdateEquipmentRequest struct {
	ID                 uuid.UUID               `validate:"required,uuid4"`
	InventoryNumber    *string                 `validate:"omitempty,min=4,max=14"`
	Name               *string                 `validate:"omitempty,min=4,max=50"`
	ShortName          *string                 `validate:"omitempty,min=4,max=20"`
	Category           *string                 `validate:"omitempty,min=4,max=20"`
	AvailableToTrainee *bool                   `validate:"omitempty"`
	Status             *domain.EquipmentStatus `validate:"omitempty"`
	Departments        []uuid.UUID             `validate:"omitempty,dive,uuid4"`
}

type updateEquipmentService struct {
	auther              authzservice.AuthZ
	equipmentRepository domain.EquipmentRepository
	txManager           txmanager.TxManager
	depAssocService     DepartmentAssociationService
}

var _ UpdateEquipmentService = (*updateEquipmentService)(nil)

func NewUpdateEquipmentService(
	auther authzservice.AuthZ,
	equipmentRepository domain.EquipmentRepository,
	depAssocService DepartmentAssociationService,
	txManager txmanager.TxManager,
) UpdateEquipmentService {
	return &updateEquipmentService{
		auther:              auther,
		equipmentRepository: equipmentRepository,
		depAssocService:     depAssocService,
		txManager:           txManager,
	}
}

func (s *updateEquipmentService) UpdateEquipment(ctx context.Context, req UpdateEquipmentRequest) (*domain.Equipment, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateEquipmentRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	var equipment *domain.Equipment
	err = s.txManager.WithinTx(ctx, func(ctx context.Context) error {

		equipment, err = s.equipmentRepository.Get(ctx, req.ID, domain.EquipmentWithDepartments())
		if err != nil {
			return err
		}
		if req.InventoryNumber != nil {
			equipment.InventoryNumber = *req.InventoryNumber
		}
		if req.Name != nil {
			equipment.Name = *req.Name
		}
		if req.ShortName != nil {
			equipment.ShortName = *req.ShortName
		}
		if req.Category != nil {
			equipment.Category = *req.Category
		}
		if req.AvailableToTrainee != nil {
			equipment.AvailableToTrainee = *req.AvailableToTrainee
		}
		if req.Status != nil {
			equipment.Status = *req.Status
		}

		err = s.equipmentRepository.Update(ctx, equipment)
		if err != nil {
			return err
		}

		if req.Departments != nil {
			err = s.depAssocService.UpdateAssociations(ctx, equipment, req.Departments)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return equipment, nil
}

package equipmentservice

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

type CreateEquipmentRequest struct {
	InventoryNumber    string                 `validate:"required,min=4,max=14"`
	Name               string                 `validate:"required,min=4,max=50"`
	ShortName          string                 `validate:"required,min=4,max=20"`
	Category           string                 `validate:"required,min=4,max=20"`
	AvailableToTrainee bool                   `validate:"required"`
	Status             domain.EquipmentStatus `validate:"required"`
	Departments        []uuid.UUID
}

type CreateEquipmentService interface {
	CreateEquipment(ctx context.Context, req *CreateEquipmentRequest) (*domain.Equipment, error)
}

type createEquipmentService struct {
	auther              authzservice.AuthZ
	equipmentRepository domain.EquipmentRepository
	depAssocService     DepartmentAssociationService
	txManager           txmanager.TxManager
}

var _ CreateEquipmentService = (*createEquipmentService)(nil)

func NewCreateEquipmentService(
	auther authzservice.AuthZ,
	equipmentRepository domain.EquipmentRepository,
	depAssocService DepartmentAssociationService,
	txManager txmanager.TxManager,
) CreateEquipmentService {
	return &createEquipmentService{
		auther:              auther,
		equipmentRepository: equipmentRepository,
		depAssocService:     depAssocService,
		txManager:           txManager,
	}
}

func (s *createEquipmentService) CreateEquipment(ctx context.Context, req *CreateEquipmentRequest) (*domain.Equipment, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("CreateEquipmentRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	equipment := &domain.Equipment{
		ID:                 uuid.New(),
		InventoryNumber:    req.InventoryNumber,
		Name:               req.Name,
		ShortName:          req.ShortName,
		Category:           req.Category,
		AvailableToTrainee: req.AvailableToTrainee,
		Status:             req.Status,
		Departments:        make([]*domain.Department, 0, len(req.Departments)),
	}

	depIDs := req.Departments
	if depIDs == nil {
		depIDs = []uuid.UUID{}
	}

	err = s.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.equipmentRepository.Create(ctx, equipment); err != nil {
			return err
		}

		err = s.depAssocService.UpdateAssociations(ctx, equipment, depIDs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return equipment, nil
}

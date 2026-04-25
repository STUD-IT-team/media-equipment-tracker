package equipmentservice

import (
	"context"
	"time"

	"media-equipment-tracker/internal/application/invocationservice/invocationsearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type AvailabilityEquipmentService interface {
	Availability(ctx context.Context, req EquipmentAvailabilityRequest) (*EquipmentAvailabilityResponse, error)
}

type availabilityEquipmentService struct {
	equipmentRepository domain.EquipmentRepository
	searchRepository    invocationsearch.SearchInvocationRepository
	txManager           txmanager.TxManager
}

type EquipmentAvailabilityRequest struct {
	ID        uuid.UUID `validate:"required"`
	StartTime time.Time `validate:"required"`
	EndTime   time.Time `validate:"required,gtfield=StartTime"`
}

type EquipmentAvailabilityResponse struct {
	Available              bool
	Equipment              *domain.Equipment
	ConflictingInvocations []*domain.EquipmentInvocation
}

func NewAvailabilityEquipmentService(
	equipmentRepository domain.EquipmentRepository,
	searchRepository invocationsearch.SearchInvocationRepository,
	txManager txmanager.TxManager,
) AvailabilityEquipmentService {
	return &availabilityEquipmentService{
		equipmentRepository: equipmentRepository,
		searchRepository:    searchRepository,
		txManager:           txManager,
	}
}

func (s *availabilityEquipmentService) Availability(ctx context.Context, req EquipmentAvailabilityRequest) (*EquipmentAvailabilityResponse, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("EquipmentAvailabilityRequest", err.Error())
	}

	var equipment *domain.Equipment
	var invocations []*domain.EquipmentInvocation
	var available bool

	if err := s.txManager.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		equipment, err = s.equipmentRepository.Get(ctx, req.ID, domain.EquipmentWithEquipmentInInvocations())
		if err != nil {
			return err
		}

		if equipment.Status == domain.EquipmentStatusUnavailable || equipment.Status == domain.EquipmentStatusUnderMaintenance {
			available = false
			invocations = []*domain.EquipmentInvocation{}
			return nil
		}

		invocations, err = s.searchRepository.Search(ctx, &invocationsearch.SearchInvocationRequest{
			EquipmentIDs: []uuid.UUID{equipment.ID},
			StartTime:    &req.StartTime,
			EndTime:      &req.EndTime,
		})
		if err != nil {
			return err
		}
		available = len(invocations) == 0
		return nil
	}); err != nil {
		return nil, err
	}

	return &EquipmentAvailabilityResponse{
		Available:              available,
		Equipment:              equipment,
		ConflictingInvocations: invocations,
	}, nil
}

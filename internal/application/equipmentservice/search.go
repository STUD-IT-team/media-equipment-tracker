package equipmentservice

import (
	"context"
	"time"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
)

type SearchEquipmentRepository interface {
	Search(ctx context.Context, search *SearchEquipmentRequest, with ...domain.EquipmentOption) ([]*domain.Equipment, error)
}

type SearchEquipmentService interface {
	Search(ctx context.Context, search *SearchEquipmentRequest, with ...domain.EquipmentOption) ([]*domain.Equipment, error)
}

type searchEquipmentService struct {
	searchRepository SearchEquipmentRepository
}

var _ SearchEquipmentService = (*searchEquipmentService)(nil)

func NewSearchEquipmentService(searchRepository SearchEquipmentRepository) SearchEquipmentService {
	return &searchEquipmentService{searchRepository: searchRepository}
}

type SearchEquipmentRequest struct {
	Categories         []string
	Statuses           []domain.EquipmentStatus
	AvailableToTrainee *bool
	DepartmentIDs      []uuid.UUID
	AvailableAt        *time.Time
	SearchString       *string
}

func (r *SearchEquipmentRequest) nilize() {
	if r.Categories != nil && len(r.Categories) == 0 {
		r.Categories = nil
	}
	if r.Statuses != nil && len(r.Statuses) == 0 {
		r.Statuses = nil
	}
	if r.DepartmentIDs != nil && len(r.DepartmentIDs) == 0 {
		r.DepartmentIDs = nil
	}
	if r.SearchString != nil && r.SearchString == nil {
		r.SearchString = nil
	}
}

func (s *searchEquipmentService) Search(ctx context.Context, search *SearchEquipmentRequest, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	search.nilize()

	if search.AvailableAt != nil && search.AvailableAt.Before(time.Now()) {
		return nil, errs.NewValidationError("AvailableAt", "AvailableAt must be in the future")
	}

	if search.Statuses == nil {
		search.Statuses = []domain.EquipmentStatus{domain.EquipmentStatusAvailable, domain.EquipmentStatusIssued}
	}

	equipment, err := s.searchRepository.Search(ctx, search, with...)
	if err != nil {
		return nil, err
	}

	return equipment, nil
}

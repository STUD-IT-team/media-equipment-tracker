package invocationsearch

import (
	"context"
	"time"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"

	"github.com/google/uuid"
)

type SearchInvocationRepository interface {
	Search(ctx context.Context, search *SearchInvocationRequest, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error)
}

type SearchInvocationService interface {
	Search(ctx context.Context, search *SearchInvocationRequest, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error)
}

type searchInvocationService struct {
	searchRepository SearchInvocationRepository
}

var _ SearchInvocationService = (*searchInvocationService)(nil)

func NewSearchInvocationService(searchRepository SearchInvocationRepository) SearchInvocationService {
	return &searchInvocationService{searchRepository: searchRepository}
}

type SearchInvocationRequest struct {
	SearchString *string
	Statuses     []domain.EquipmentInvocationStatus
	EquipmentIDs []uuid.UUID

	StartTime *time.Time `validate:"omitempty"`
	EndTime   *time.Time `validate:"omitempty,gtfield=StartTime"`

	AdminID *uuid.UUID
	UserID  *uuid.UUID

	DepartmentID   *uuid.UUID
	OrganizationID *uuid.UUID
}

func (r *SearchInvocationRequest) nilize() {
	if r.SearchString != nil && *r.SearchString == "" {
		r.SearchString = nil
	}
	if r.Statuses != nil && len(r.Statuses) == 0 {
		r.Statuses = nil
	}
	if r.EquipmentIDs != nil && len(r.EquipmentIDs) == 0 {
		r.EquipmentIDs = nil
	}
	if r.AdminID != nil && *r.AdminID == uuid.Nil {
		r.AdminID = nil
	}
	if r.UserID != nil && *r.UserID == uuid.Nil {
		r.UserID = nil
	}
	if r.DepartmentID != nil && *r.DepartmentID == uuid.Nil {
		r.DepartmentID = nil
	}
	if r.OrganizationID != nil && *r.OrganizationID == uuid.Nil {
		r.OrganizationID = nil
	}
}

func (s *searchInvocationService) Search(ctx context.Context, search *SearchInvocationRequest, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
	if err := validate.ValidateStruct(search); err != nil {
		return nil, errs.NewValidationError("SearchInvocationRequest", err.Error())
	}
	search.nilize()

	if search.Statuses == nil {
		search.Statuses = []domain.EquipmentInvocationStatus{domain.InvocationCreated, domain.InvocationUnderReview, domain.InvocationChangesRequired, domain.InvocationApproved, domain.InvocationEquipmentIssued, domain.InvocationEquipmentReturned, domain.InvocationCompleted, domain.InvocationCancelled}
	}

	invocation, err := s.searchRepository.Search(ctx, search, with...)
	if err != nil {
		return nil, err
	}

	return invocation, nil
}

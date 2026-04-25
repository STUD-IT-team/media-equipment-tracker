package invocationsearch

import (
	"context"
	"time"

	"media-equipment-tracker/internal/domain"

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

	StartTime *time.Time
	EndTime   *time.Time

	AdminID *uuid.UUID
	UserID  *uuid.UUID
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
	if r.StartTime != nil && r.StartTime.Before(time.Now()) {
		r.StartTime = nil
	}
	if r.EndTime != nil && r.EndTime.Before(time.Now()) {
		r.EndTime = nil
	}
	if r.AdminID != nil && *r.AdminID == uuid.Nil {
		r.AdminID = nil
	}
	if r.UserID != nil && *r.UserID == uuid.Nil {
		r.UserID = nil
	}
}

func (s *searchInvocationService) Search(ctx context.Context, search *SearchInvocationRequest, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
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

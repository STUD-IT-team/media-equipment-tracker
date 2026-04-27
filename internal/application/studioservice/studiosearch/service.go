package studiosearch

import (
	"context"
	"time"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"

	"github.com/google/uuid"
)

type SearchStudioInvocationRequest struct {
	SearchString   *string
	Statuses       []domain.StudioInvocationStatus
	StartTime      *time.Time `validate:"omitempty"`
	EndTime        *time.Time `validate:"omitempty,gtfield=StartTime"`
	AdminID        *uuid.UUID
	UserID         *uuid.UUID
	DepartmentID   *uuid.UUID
	OrganizationID *uuid.UUID
}

type SearchStudioInvocationRepository interface {
	Search(ctx context.Context, search *SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error)
}

type SearchStudioService interface {
	Search(ctx context.Context, search *SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error)
}

type searchStudioService struct {
	searchRepository SearchStudioInvocationRepository
}

var _ SearchStudioService = (*searchStudioService)(nil)

func NewSearchStudioService(searchRepository SearchStudioInvocationRepository) SearchStudioService {
	return &searchStudioService{searchRepository: searchRepository}
}

func (r *SearchStudioInvocationRequest) nilize() {
	if r.SearchString != nil && *r.SearchString == "" {
		r.SearchString = nil
	}
	if r.Statuses != nil && len(r.Statuses) == 0 {
		r.Statuses = nil
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

func (s *searchStudioService) Search(ctx context.Context, search *SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	if err := validate.ValidateStruct(search); err != nil {
		return nil, errs.NewValidationError("SearchStudioInvocationRequest", err.Error())
	}
	search.nilize()

	if search.Statuses == nil {
		search.Statuses = []domain.StudioInvocationStatus{domain.StudioCreated, domain.StudioUnderReview, domain.StudioChangesRequired, domain.StudioApproved}
	}

	invocation, err := s.searchRepository.Search(ctx, search, with...)
	if err != nil {
		return nil, err
	}

	return invocation, nil
}

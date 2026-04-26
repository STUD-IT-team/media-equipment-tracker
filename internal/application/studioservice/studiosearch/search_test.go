//go:build unit

package studiosearch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type SearchStudioServiceSuite struct {
	suite.Suite

	svc         studiosearch.SearchStudioService
	searchRepo  *SearchStudioInvocationRepoMock
}

func (s *SearchStudioServiceSuite) SetupTest() {
	s.searchRepo = new(SearchStudioInvocationRepoMock)
	s.svc = studiosearch.NewSearchStudioService(s.searchRepo)
}

func (s *SearchStudioServiceSuite) TestSearch_Success_WithDefaultStatuses() {
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := &studiosearch.SearchStudioInvocationRequest{
		StartTime: &start,
		EndTime:   &end,
	}

	expectedResults := []*domain.StudioInvocation{{ID: uuid.New()}}

	ctx := context.Background()
	s.searchRepo.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return r.StartTime != nil && r.EndTime != nil && len(r.Statuses) == 4
	}), mock.Anything).Return(expectedResults, nil)

	result, err := s.svc.Search(ctx, req)

	s.NoError(err)
	s.Equal(expectedResults, result)
}

func (s *SearchStudioServiceSuite) TestSearch_Success_WithCustomStatuses() {
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := &studiosearch.SearchStudioInvocationRequest{
		StartTime: &start,
		EndTime:   &end,
		Statuses:  []domain.StudioInvocationStatus{domain.StudioApproved, domain.StudioCompleted},
	}

	expectedResults := []*domain.StudioInvocation{{ID: uuid.New()}}

	ctx := context.Background()
	s.searchRepo.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return r.StartTime != nil && r.EndTime != nil && len(r.Statuses) == 2 &&
			r.Statuses[0] == domain.StudioApproved && r.Statuses[1] == domain.StudioCompleted
	}), mock.Anything).Return(expectedResults, nil)

	result, err := s.svc.Search(ctx, req)

	s.NoError(err)
	s.Equal(expectedResults, result)
}

func (s *SearchStudioServiceSuite) TestSearch_Success_WithFilters() {
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	userID := uuid.New()
	depID := uuid.New()
	adminID := uuid.New()
	searchStr := "test event"
	req := &studiosearch.SearchStudioInvocationRequest{
		SearchString:   &searchStr,
		StartTime:      &start,
		EndTime:        &end,
		UserID:         &userID,
		DepartmentID:   &depID,
		AdminID:        &adminID,
		Statuses:       []domain.StudioInvocationStatus{domain.StudioUnderReview},
	}

	expectedResults := []*domain.StudioInvocation{{ID: uuid.New()}}

	ctx := context.Background()
	s.searchRepo.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return r.SearchString != nil && *r.SearchString == searchStr &&
			r.UserID != nil && *r.UserID == userID &&
			r.DepartmentID != nil && *r.DepartmentID == depID &&
			r.AdminID != nil && *r.AdminID == adminID &&
			len(r.Statuses) == 1 && r.Statuses[0] == domain.StudioUnderReview
	}), mock.Anything).Return(expectedResults, nil)

	result, err := s.svc.Search(ctx, req)

	s.NoError(err)
	s.Equal(expectedResults, result)
}

func (s *SearchStudioServiceSuite) TestSearch_InvalidRequest_EndBeforeStart() {
	req := &studiosearch.SearchStudioInvocationRequest{
		StartTime: &time.Time{}, // zero time
		EndTime:   &time.Time{}, // same zero time - invalid
	}

	ctx := context.Background()

	_, err := s.svc.Search(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *SearchStudioServiceSuite) TestSearch_RepositoryError() {
	req := &studiosearch.SearchStudioInvocationRequest{
		Statuses: []domain.StudioInvocationStatus{domain.StudioCreated},
	}

	ctx := context.Background()
	s.searchRepo.On("Search", ctx, mock.Anything, mock.Anything).Return(([]*domain.StudioInvocation)(nil), errs.NewRepositoryError("search", errors.New("database error")))

	_, err := s.svc.Search(ctx, req)

	s.Error(err)
	s.True(errs.IsRepositoryError(err))
}

func (s *SearchStudioServiceSuite) TestSearch_EmptyStatusesDefaults() {
	req := &studiosearch.SearchStudioInvocationRequest{
		Statuses: []domain.StudioInvocationStatus{}, // empty slice
	}

	expectedResults := []*domain.StudioInvocation{{ID: uuid.New()}}

	ctx := context.Background()
	s.searchRepo.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return len(r.Statuses) == 4 // should be set to defaults
	}), mock.Anything).Return(expectedResults, nil)

	result, err := s.svc.Search(ctx, req)

	s.NoError(err)
	s.Equal(expectedResults, result)
}

func TestSearchStudioServiceSuite(t *testing.T) {
	suite.Run(t, new(SearchStudioServiceSuite))
}

type SearchStudioInvocationRepoMock struct{ mock.Mock }

func (m *SearchStudioInvocationRepoMock) Search(ctx context.Context, search *studiosearch.SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	args := m.Called(ctx, search, with)
	return args.Get(0).([]*domain.StudioInvocation), args.Error(1)
}
//go:build unit

package studioservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	studiosearch "media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type AvailabilityStudioSuite struct {
	suite.Suite

	svc       studioservice.AvailabilityStudioService
	searchSvc *SearchStudioServiceMock
}

func (s *AvailabilityStudioSuite) SetupTest() {
	s.searchSvc = new(SearchStudioServiceMock)

	s.svc = studioservice.NewStudioInvocationAvailabilityService(s.searchSvc)
}

func (s *AvailabilityStudioSuite) TestAvailability_Available_NoConflicts() {
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := studioservice.AvailabilityStudioInvocationRequest{
		StartTime: start,
		EndTime:   end,
	}

	ctx := context.Background()
	s.searchSvc.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return r.StartTime != nil && r.EndTime != nil
	}), mock.Anything).Return([]*domain.StudioInvocation{}, nil)

	result, err := s.svc.Availability(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.True(result.Available)
	s.Empty(result.ConflictingInvocations)
}

func (s *AvailabilityStudioSuite) TestAvailability_NotAvailable_HasConflicts() {
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := studioservice.AvailabilityStudioInvocationRequest{
		StartTime: start,
		EndTime:   end,
	}

	conflictingInv := &domain.StudioInvocation{ID: uuid.New()}

	ctx := context.Background()
	s.searchSvc.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return r.StartTime != nil && r.EndTime != nil
	}), mock.Anything).Return([]*domain.StudioInvocation{conflictingInv}, nil)

	result, err := s.svc.Availability(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.False(result.Available)
	s.Len(result.ConflictingInvocations, 1)
	s.Equal(conflictingInv, result.ConflictingInvocations[0])
}

func (s *AvailabilityStudioSuite) TestAvailability_ExcludeOwnInvocation() {
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	ownID := uuid.New()
	req := studioservice.AvailabilityStudioInvocationRequest{
		StartTime:     start,
		EndTime:       end,
		ForInvocation: &ownID,
	}

	conflictingInv := &domain.StudioInvocation{ID: uuid.New()}
	ownInv := &domain.StudioInvocation{ID: ownID}

	ctx := context.Background()
	s.searchSvc.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return r.StartTime != nil && r.EndTime != nil
	}), mock.Anything).Return([]*domain.StudioInvocation{conflictingInv, ownInv}, nil)

	result, err := s.svc.Availability(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.False(result.Available)
	s.Len(result.ConflictingInvocations, 1)
	s.Equal(conflictingInv, result.ConflictingInvocations[0])
}

func (s *AvailabilityStudioSuite) TestAvailability_InvalidRequest_EndBeforeStart() {
	req := studioservice.AvailabilityStudioInvocationRequest{
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now(), // before start
	}

	ctx := context.Background()

	_, err := s.svc.Availability(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *AvailabilityStudioSuite) TestAvailability_SearchError() {
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := studioservice.AvailabilityStudioInvocationRequest{
		StartTime: start,
		EndTime:   end,
	}

	ctx := context.Background()
	s.searchSvc.On("Search", ctx, mock.Anything, mock.Anything).Return(([]*domain.StudioInvocation)(nil), errs.NewRepositoryError("search", errors.New("database error")))

	_, err := s.svc.Availability(ctx, req)

	s.Error(err)
	s.True(errs.IsRepositoryError(err))
}

func TestAvailabilityStudioSuite(t *testing.T) {
	suite.Run(t, new(AvailabilityStudioSuite))
}
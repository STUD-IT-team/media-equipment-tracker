//go:build unit

package invocationservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"
)

type InvocationSuite struct {
	suite.Suite

	svc        invocationservice.SearchInvocationService
	searchRepo *SearchInvocationRepoMock
}

func (s *InvocationSuite) SetupTest() {
	s.searchRepo = new(SearchInvocationRepoMock)
	s.svc = invocationservice.NewSearchInvocationService(s.searchRepo)
}

func (s *InvocationSuite) TestSearch_Success() {
	req := &invocationservice.SearchInvocationRequest{
		SearchString: stringPtr("test"),
		Statuses:     []domain.EquipmentInvocationStatus{domain.InvocationCreated},
		EquipmentIDs: []uuid.UUID{uuid.New()},
		StartTime:    &time.Time{},
		EndTime:      &time.Time{},
	}
	invocations := []*domain.EquipmentInvocation{{ID: uuid.New()}}

	s.searchRepo.On("Search", mock.Anything, req, mock.Anything).Return(invocations, nil)

	result, err := s.svc.Search(context.Background(), req)

	s.NoError(err)
	s.Equal(invocations, result)
}

func (s *InvocationSuite) TestSearch_NilizeAndSetStatuses() {
	req := &invocationservice.SearchInvocationRequest{
		Statuses: nil, // should set default statuses
	}
	invocations := []*domain.EquipmentInvocation{}

	s.searchRepo.On("Search", mock.Anything, mock.MatchedBy(func(r *invocationservice.SearchInvocationRequest) bool {
		return r.Statuses != nil && len(r.Statuses) > 0
	}), mock.Anything).Return(invocations, nil)

	result, err := s.svc.Search(context.Background(), req)

	s.NoError(err)
	s.Equal(invocations, result)
}

// Helper
func stringPtr(s string) *string {
	return &s
}

// Mock
type SearchInvocationRepoMock struct{ mock.Mock }

func (m *SearchInvocationRepoMock) Search(ctx context.Context, search *invocationservice.SearchInvocationRequest, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
	args := m.Called(ctx, search, with)
	return args.Get(0).([]*domain.EquipmentInvocation), args.Error(1)
}

func TestInvocationSuite(t *testing.T) {
	suite.Run(t, new(InvocationSuite))
}

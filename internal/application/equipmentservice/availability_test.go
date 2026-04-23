//go:build unit

package equipmentservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"
)

type AvailabilityEquipmentSuite struct {
	suite.Suite

	svc                  equipmentservice.AvailabilityEquipmentService
	equipmentRepo        *EquipmentRepoMock
	searchInvocationRepo *SearchInvocationRepoMock
	txManager            *TxManagerMock
}

func (s *AvailabilityEquipmentSuite) SetupTest() {
	s.equipmentRepo = new(EquipmentRepoMock)
	s.searchInvocationRepo = new(SearchInvocationRepoMock)
	s.txManager = new(TxManagerMock)

	s.svc = equipmentservice.NewAvailabilityEquipmentService(s.equipmentRepo, s.searchInvocationRepo, s.txManager)
}

func (s *AvailabilityEquipmentSuite) TestAvailability_Available() {
	eqID := uuid.New()
	eq := &domain.Equipment{
		ID:     eqID,
		Status: domain.EquipmentStatusAvailable,
	}
	req := equipmentservice.EquipmentAvailabilityRequest{
		ID:        eqID,
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	ctx := context.Background()
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.equipmentRepo.On("Get", mock.Anything, eqID, mock.Anything).Return(eq, nil)
		s.searchInvocationRepo.On("Search", mock.Anything, mock.AnythingOfType("*invocationservice.SearchInvocationRequest"), mock.Anything).Return([]*domain.EquipmentInvocation{}, nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.Availability(ctx, req)

	s.NoError(err)
	s.True(result.Available)
	s.Empty(result.ConflictingInvocations)
}

func (s *AvailabilityEquipmentSuite) TestAvailability_UnavailableStatus() {
	eqID := uuid.New()
	eq := &domain.Equipment{
		ID:     eqID,
		Status: domain.EquipmentStatusUnavailable,
	}
	req := equipmentservice.EquipmentAvailabilityRequest{
		ID:        eqID,
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	ctx := context.Background()
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.equipmentRepo.On("Get", mock.Anything, eqID, mock.Anything).Return(eq, nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.Availability(ctx, req)

	s.NoError(err)
	s.False(result.Available)
	s.Empty(result.ConflictingInvocations)
}

func (s *AvailabilityEquipmentSuite) TestAvailability_HasConflicts() {
	eqID := uuid.New()
	eq := &domain.Equipment{
		ID:     eqID,
		Status: domain.EquipmentStatusAvailable,
	}
	invocations := []*domain.EquipmentInvocation{{ID: uuid.New()}}
	req := equipmentservice.EquipmentAvailabilityRequest{
		ID:        eqID,
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	ctx := context.Background()
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.equipmentRepo.On("Get", mock.Anything, eqID, mock.Anything).Return(eq, nil)
		s.searchInvocationRepo.On("Search", mock.Anything, mock.AnythingOfType("*invocationservice.SearchInvocationRequest"), mock.Anything).Return(invocations, nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.Availability(ctx, req)

	s.NoError(err)
	s.False(result.Available)
	s.Equal(invocations, result.ConflictingInvocations)
}

func TestAvailabilityEquipmentSuite(t *testing.T) {
	suite.Run(t, new(AvailabilityEquipmentSuite))
}

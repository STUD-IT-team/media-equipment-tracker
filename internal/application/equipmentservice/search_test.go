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
	"media-equipment-tracker/internal/domain/errs"
)

type SearchEquipmentSuite struct {
	suite.Suite

	svc        equipmentservice.SearchEquipmentService
	searchRepo *SearchEquipmentRepoMock
}

func (s *SearchEquipmentSuite) SetupTest() {
	s.searchRepo = new(SearchEquipmentRepoMock)
	s.svc = equipmentservice.NewSearchEquipmentService(s.searchRepo)
}

func (s *SearchEquipmentSuite) TestSearch_Success() {
	req := &equipmentservice.SearchEquipmentRequest{
		Categories: []string{"Test"},
		Statuses:   []domain.EquipmentStatus{domain.EquipmentStatusAvailable},
	}
	equipments := []*domain.Equipment{{ID: uuid.New()}}

	s.searchRepo.On("Search", mock.Anything, req, mock.Anything).Return(equipments, nil)

	result, err := s.svc.Search(context.Background(), req)

	s.NoError(err)
	s.Equal(equipments, result)
}

func (s *SearchEquipmentSuite) TestSearch_SetDefaultStatuses() {
	req := &equipmentservice.SearchEquipmentRequest{
		Statuses: nil, // should set defaults
	}
	equipments := []*domain.Equipment{}

	s.searchRepo.On("Search", mock.Anything, mock.MatchedBy(func(r *equipmentservice.SearchEquipmentRequest) bool {
		return r.Statuses != nil && len(r.Statuses) > 0
	}), mock.Anything).Return(equipments, nil)

	result, err := s.svc.Search(context.Background(), req)

	s.NoError(err)
	s.Equal(equipments, result)
}

func (s *SearchEquipmentSuite) TestSearch_InvalidAvailableAt() {
	pastTime := time.Now().Add(-time.Hour)
	req := &equipmentservice.SearchEquipmentRequest{
		AvailableAt: &pastTime,
	}

	_, err := s.svc.Search(context.Background(), req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestSearchEquipmentSuite(t *testing.T) {
	suite.Run(t, new(SearchEquipmentSuite))
}

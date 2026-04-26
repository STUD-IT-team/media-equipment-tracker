//go:build unit

package departmentservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/departmentservice"
	"media-equipment-tracker/internal/domain"
)

type GetDepartmentEquipmentSuite struct {
	suite.Suite

	svc            departmentservice.GetDepartmentEquipmentService
	departmentRepo *DepartmentRepoMock
}

func (s *GetDepartmentEquipmentSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)

	s.svc = departmentservice.NewGetDepartmentEquipmentService(s.departmentRepo)
}

func (s *GetDepartmentEquipmentSuite) TestGetDepartmentEquipment_Success() {
	depID := uuid.New()
	equipment := []*domain.Equipment{
		{ID: uuid.New(), Name: "Equipment1", Status: domain.EquipmentStatusAvailable},
		{ID: uuid.New(), Name: "Equipment2", Status: domain.EquipmentStatusIssued},
	}
	department := &domain.Department{ID: depID, Name: "Test Department", Equipment: equipment}

	ctx := context.Background()
	s.departmentRepo.On("Get", ctx, depID, mock.Anything).Return(department, nil)

	result, err := s.svc.GetDepartmentEquipment(ctx, depID, false)

	s.NoError(err)
	s.Len(result, 2)
}

func (s *GetDepartmentEquipmentSuite) TestGetDepartmentEquipment_AvailableOnly() {
	depID := uuid.New()
	equipment := []*domain.Equipment{
		{ID: uuid.New(), Name: "Equipment1", Status: domain.EquipmentStatusAvailable},
		{ID: uuid.New(), Name: "Equipment2", Status: domain.EquipmentStatusIssued},
	}
	department := &domain.Department{ID: depID, Name: "Test Department", Equipment: equipment}

	ctx := context.Background()
	s.departmentRepo.On("Get", ctx, depID, mock.Anything).Return(department, nil)

	result, err := s.svc.GetDepartmentEquipment(ctx, depID, true)

	s.NoError(err)
	s.Len(result, 1)
	s.Equal(domain.EquipmentStatusAvailable, result[0].Status)
}

func TestGetDepartmentEquipmentSuite(t *testing.T) {
	suite.Run(t, new(GetDepartmentEquipmentSuite))
}

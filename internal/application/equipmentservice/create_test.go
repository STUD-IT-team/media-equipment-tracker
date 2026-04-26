//go:build unit

package equipmentservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type CreateEquipmentSuite struct {
	suite.Suite

	svc             equipmentservice.CreateEquipmentService
	equipmentRepo   *EquipmentRepoMock
	auther          *AuthZMock
	txManager       *TxManagerMock
	depAssocService *DepartmentAssociationServiceMock
}

func (s *CreateEquipmentSuite) SetupTest() {
	s.equipmentRepo = new(EquipmentRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)
	s.depAssocService = new(DepartmentAssociationServiceMock)

	s.svc = equipmentservice.NewCreateEquipmentService(s.auther, s.equipmentRepo, s.depAssocService, s.txManager)
}

func (s *CreateEquipmentSuite) TestCreateEquipment_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	req := &equipmentservice.CreateEquipmentRequest{
		InventoryNumber:    "INV001",
		Name:               "Test Equipment",
		ShortName:          "Test",
		Category:           "Category",
		AvailableToTrainee: true,
		Status:             domain.EquipmentStatusAvailable,
		Departments:        []uuid.UUID{uuid.New()},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.equipmentRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Equipment")).Return(nil)
		s.depAssocService.On("UpdateAssociations", mock.Anything, mock.AnythingOfType("*domain.Equipment"), req.Departments).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.CreateEquipment(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.InventoryNumber, result.InventoryNumber)
}

func (s *CreateEquipmentSuite) TestCreateEquipment_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	req := &equipmentservice.CreateEquipmentRequest{
		InventoryNumber:    "INV001",
		Name:               "Test",
		ShortName:          "Test",
		Category:           "Category",
		Status:             domain.EquipmentStatusAvailable,
		AvailableToTrainee: true,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.CreateEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CreateEquipmentSuite) TestCreateEquipment_NoDepartments() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	req := &equipmentservice.CreateEquipmentRequest{
		InventoryNumber:    "INV001",
		Name:               "Test Equipment",
		ShortName:          "Test",
		Category:           "Category",
		AvailableToTrainee: true,
		Status:             domain.EquipmentStatusAvailable,
		Departments:        []uuid.UUID{}, // empty
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.equipmentRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Equipment")).Return(nil)
		s.depAssocService.On("UpdateAssociations", mock.Anything, mock.AnythingOfType("*domain.Equipment"), []uuid.UUID{}).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.CreateEquipment(ctx, req)

	s.NoError(err)
	s.NotNil(result)
}

// Mock for DepartmentAssociationService
type DepartmentAssociationServiceMock struct{ mock.Mock }

func (m *DepartmentAssociationServiceMock) UpdateAssociations(ctx context.Context, equipment *domain.Equipment, newDepIDs []uuid.UUID) error {
	return m.Called(ctx, equipment, newDepIDs).Error(0)
}

func TestCreateEquipmentSuite(t *testing.T) {
	suite.Run(t, new(CreateEquipmentSuite))
}

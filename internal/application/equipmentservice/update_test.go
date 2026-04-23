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

type UpdateEquipmentSuite struct {
	suite.Suite

	svc             equipmentservice.UpdateEquipmentService
	equipmentRepo   *EquipmentRepoMock
	auther          *AuthZMock
	txManager       *TxManagerMock
	depAssocService *DepartmentAssociationServiceMock
}

func (s *UpdateEquipmentSuite) SetupTest() {
	s.equipmentRepo = new(EquipmentRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)
	s.depAssocService = new(DepartmentAssociationServiceMock)

	s.svc = equipmentservice.NewUpdateEquipmentService(s.auther, s.equipmentRepo, s.depAssocService, s.txManager)
}

func (s *UpdateEquipmentSuite) TestUpdateEquipment_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	eqID := uuid.New()
	eq := &domain.Equipment{
		ID:              eqID,
		InventoryNumber: "OLD",
		Name:            "Old Name",
		Departments:     []*domain.Department{},
	}
	newName := "New Name"
	req := &equipmentservice.UpdateEquipmentRequest{
		ID:   eqID,
		Name: &newName,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.equipmentRepo.On("Get", mock.Anything, eqID, mock.Anything).Return(eq, nil)
		s.equipmentRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Equipment")).Return(nil).Run(func(args mock.Arguments) {
			updatedEq := args.Get(1).(*domain.Equipment)
			s.Equal(newName, updatedEq.Name)
		})
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.UpdateEquipment(ctx, req)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *UpdateEquipmentSuite) TestUpdateEquipment_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	req := &equipmentservice.UpdateEquipmentRequest{
		ID:   uuid.New(),
		Name: stringPtr("New Name"),
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.UpdateEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *UpdateEquipmentSuite) TestUpdateEquipment_WithDepartments() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	eqID := uuid.New()
	eq := &domain.Equipment{
		ID:              eqID,
		InventoryNumber: "OLD",
		Name:            "Old Name",
		Departments:     []*domain.Department{},
	}
	newDeps := []uuid.UUID{uuid.New()}
	req := &equipmentservice.UpdateEquipmentRequest{
		ID:          eqID,
		Name:        stringPtr("New Name"),
		Departments: newDeps,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.equipmentRepo.On("Get", mock.Anything, eqID, mock.Anything).Return(eq, nil)
		s.equipmentRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Equipment")).Return(nil)
		s.depAssocService.On("UpdateAssociations", mock.Anything, eq, newDeps).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.UpdateEquipment(ctx, req)

	s.NoError(err)
	s.Equal(*req.Name, result.Name)
	// s.Equal(req.)
	resDeps := make([]uuid.UUID, 0, len(result.Departments))
	for _, dep := range result.Departments {
		resDeps = append(resDeps, dep.ID)
	}
	s.Equal(newDeps, resDeps)
}

func stringPtr(s string) *string {
	return &s
}

func TestUpdateEquipmentSuite(t *testing.T) {
	suite.Run(t, new(UpdateEquipmentSuite))
}

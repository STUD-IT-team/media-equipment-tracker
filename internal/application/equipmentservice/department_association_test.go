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
)

type DepartmentAssociationSuite struct {
	suite.Suite

	svc            equipmentservice.DepartmentAssociationService
	departmentRepo *DepartmentRepoMock
	txManager      *TxManagerMock
}

func (s *DepartmentAssociationSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)
	s.txManager = new(TxManagerMock)

	s.svc = equipmentservice.NewDepartmentAssociationService(s.departmentRepo, s.txManager)
}

func (s *DepartmentAssociationSuite) TestUpdateAssociations_Add() {
	eq := &domain.Equipment{
		ID:          uuid.New(),
		Departments: []*domain.Department{},
	}
	depID := uuid.New()
	dep := &domain.Department{
		ID:        depID,
		Equipment: []*domain.Equipment{},
	}

	ctx := context.Background()
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.departmentRepo.On("Get", mock.Anything, depID, mock.Anything).Return(dep, nil)
		s.departmentRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Department")).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	err := s.svc.UpdateAssociations(ctx, eq, []uuid.UUID{depID})

	s.NoError(err)
}

func (s *DepartmentAssociationSuite) TestUpdateAssociations_Remove() {
	depID := uuid.New()
	eqID := uuid.New()
	dep := &domain.Department{
		ID:        depID,
		Equipment: []*domain.Equipment{{ID: eqID}},
	}
	eq := &domain.Equipment{
		ID:          eqID,
		Departments: []*domain.Department{dep},
	}

	ctx := context.Background()
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.departmentRepo.On("Get", mock.Anything, depID, mock.Anything).Return(dep, nil)
		s.departmentRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Department")).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	err := s.svc.UpdateAssociations(ctx, eq, []uuid.UUID{}) // remove

	s.NoError(err)
}

func TestDepartmentAssociationSuite(t *testing.T) {
	suite.Run(t, new(DepartmentAssociationSuite))
}

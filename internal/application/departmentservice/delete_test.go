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
	"media-equipment-tracker/internal/domain/errs"
)

type DeleteDepartmentSuite struct {
	suite.Suite

	svc         departmentservice.DeleteDepartmentService
	departmentRepo *DepartmentRepoMock
	auther       *AuthZMock
	txManager    *TxManagerMock
}

func (s *DeleteDepartmentSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = departmentservice.NewDeleteDepartmentService(s.auther, s.departmentRepo, s.txManager)
}

func (s *DeleteDepartmentSuite) TestDeleteDepartment_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	depID := uuid.New()
	department := &domain.Department{ID: depID, Name: "Test Department", Users: []*domain.UserDepartment{}, EquipmentInvocations: []*domain.EquipmentInvocation{}, StudioInvocations: []*domain.StudioInvocation{}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.departmentRepo.On("Get", mock.Anything, depID, mock.Anything).Return(department, nil)
		s.departmentRepo.On("Delete", mock.Anything, depID).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	err := s.svc.DeleteDepartment(ctx, depID)

	s.NoError(err)
}

func (s *DeleteDepartmentSuite) TestDeleteDepartment_HasUsers() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	depID := uuid.New()
	department := &domain.Department{ID: depID, Name: "Test Department", Users: []*domain.UserDepartment{{}}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(departmentservice.ErrDepartmentHasUsers).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.departmentRepo.On("Get", mock.Anything, depID, mock.Anything).Return(department, nil)
		fn(ctx)
	})

	err := s.svc.DeleteDepartment(ctx, depID)

	s.Error(err)
	s.Equal(departmentservice.ErrDepartmentHasUsers, err)
}

func (s *DeleteDepartmentSuite) TestDeleteDepartment_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	depID := uuid.New()

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	err := s.svc.DeleteDepartment(ctx, depID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func TestDeleteDepartmentSuite(t *testing.T) {
	suite.Run(t, new(DeleteDepartmentSuite))
}
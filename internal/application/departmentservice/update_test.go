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

type UpdateDepartmentSuite struct {
	suite.Suite

	svc         departmentservice.UpdateDepartmentService
	departmentRepo *DepartmentRepoMock
	auther       *AuthZMock
	txManager    *TxManagerMock
}

func (s *UpdateDepartmentSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = departmentservice.NewUpdateDepartmentService(s.auther, s.departmentRepo, s.txManager)
}

func (s *UpdateDepartmentSuite) TestUpdateDepartment_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	depID := uuid.New()
	req := &departmentservice.UpdateDepartmentRequest{
		Name: "Updated Department",
	}
	existingDep := &domain.Department{ID: depID, Name: "Old Name"}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.departmentRepo.On("Get", ctx, depID, mock.Anything).Return(existingDep, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.departmentRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Department")).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.UpdateDepartment(ctx, depID, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *UpdateDepartmentSuite) TestUpdateDepartment_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	depID := uuid.New()
	req := &departmentservice.UpdateDepartmentRequest{
		Name: "Updated Department",
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.UpdateDepartment(ctx, depID, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func TestUpdateDepartmentSuite(t *testing.T) {
	suite.Run(t, new(UpdateDepartmentSuite))
}
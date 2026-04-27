//go:build unit

package departmentservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/departmentservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type CreateDepartmentSuite struct {
	suite.Suite

	svc            departmentservice.CreateDepartmentService
	departmentRepo *DepartmentRepoMock
	auther         *AuthZMock
}

func (s *CreateDepartmentSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)
	s.auther = new(AuthZMock)

	s.svc = departmentservice.NewCreateDepartmentService(s.auther, s.departmentRepo)
}

func (s *CreateDepartmentSuite) TestCreateDepartment_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	req := &departmentservice.CreateDepartmentRequest{
		Name: "Test Department",
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.departmentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Department")).Return(nil)

	result, err := s.svc.CreateDepartment(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *CreateDepartmentSuite) TestCreateDepartment_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	req := &departmentservice.CreateDepartmentRequest{
		Name: "Test Department",
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.CreateDepartment(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CreateDepartmentSuite) TestCreateDepartment_ValidationError() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	req := &departmentservice.CreateDepartmentRequest{
		Name: "", // invalid
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.CreateDepartment(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestCreateDepartmentSuite(t *testing.T) {
	suite.Run(t, new(CreateDepartmentSuite))
}

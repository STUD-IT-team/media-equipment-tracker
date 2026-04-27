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

type GetDepartmentUsersSuite struct {
	suite.Suite

	svc            departmentservice.GetDepartmentUsersService
	departmentRepo *DepartmentRepoMock
}

func (s *GetDepartmentUsersSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)

	s.svc = departmentservice.NewGetDepartmentUsersService(s.departmentRepo)
}

func (s *GetDepartmentUsersSuite) TestGetDepartmentUsers_Success() {
	depID := uuid.New()
	users := []*domain.UserDepartment{
		{User: &domain.User{ID: uuid.New(), FullName: "User1"}, Role: domain.RoleTrainee},
		{User: &domain.User{ID: uuid.New(), FullName: "User2"}, Role: domain.RoleActivist},
	}
	department := &domain.Department{ID: depID, Name: "Test Department", Users: users}

	ctx := context.Background()
	s.departmentRepo.On("Get", ctx, depID, mock.Anything).Return(department, nil)

	result, err := s.svc.GetDepartmentUsers(ctx, depID, nil)

	s.NoError(err)
	s.Len(result, 2)
}

func (s *GetDepartmentUsersSuite) TestGetDepartmentUsers_WithRole() {
	depID := uuid.New()
	role := domain.RoleActivist
	users := []*domain.UserDepartment{
		{User: &domain.User{ID: uuid.New(), FullName: "User1"}, Role: domain.RoleTrainee},
		{User: &domain.User{ID: uuid.New(), FullName: "User2"}, Role: domain.RoleActivist},
	}
	department := &domain.Department{ID: depID, Name: "Test Department", Users: users}

	ctx := context.Background()
	s.departmentRepo.On("Get", ctx, depID, mock.Anything).Return(department, nil)

	result, err := s.svc.GetDepartmentUsers(ctx, depID, &role)

	s.NoError(err)
	s.Len(result, 1)
	s.Equal(domain.RoleActivist, result[0].Role)
}

func TestGetDepartmentUsersSuite(t *testing.T) {
	suite.Run(t, new(GetDepartmentUsersSuite))
}

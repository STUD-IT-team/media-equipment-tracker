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

type ListDepartmentSuite struct {
	suite.Suite

	svc         departmentservice.ListDepartmentService
	departmentRepo *DepartmentRepoMock
}

func (s *ListDepartmentSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)

	s.svc = departmentservice.NewListDepartmentService(s.departmentRepo)
}

func (s *ListDepartmentSuite) TestListDepartments_Success() {
	departments := []*domain.Department{
		{ID: uuid.New(), Name: "Dep1"},
		{ID: uuid.New(), Name: "Dep2"},
	}

	ctx := context.Background()
	s.departmentRepo.On("List", ctx, mock.Anything).Return(departments, nil)

	result, err := s.svc.ListDepartments(ctx, "")

	s.NoError(err)
	s.Len(result, 2)
}

func (s *ListDepartmentSuite) TestListDepartments_WithSearch() {
	departments := []*domain.Department{
		{ID: uuid.New(), Name: "Test Department"},
		{ID: uuid.New(), Name: "Other Department"},
	}

	ctx := context.Background()
	s.departmentRepo.On("List", ctx, mock.Anything).Return(departments, nil)

	result, err := s.svc.ListDepartments(ctx, "test")

	s.NoError(err)
	s.Len(result, 1)
	s.Equal("Test Department", result[0].Name)
}

func TestListDepartmentSuite(t *testing.T) {
	suite.Run(t, new(ListDepartmentSuite))
}
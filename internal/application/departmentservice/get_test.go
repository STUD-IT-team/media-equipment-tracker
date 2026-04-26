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

type GetDepartmentSuite struct {
	suite.Suite

	svc         departmentservice.GetDepartmentService
	departmentRepo *DepartmentRepoMock
}

func (s *GetDepartmentSuite) SetupTest() {
	s.departmentRepo = new(DepartmentRepoMock)

	s.svc = departmentservice.NewGetDepartmentService(s.departmentRepo)
}

func (s *GetDepartmentSuite) TestGetDepartment_Success() {
	depID := uuid.New()
	department := &domain.Department{ID: depID, Name: "Test Department"}

	ctx := context.Background()
	s.departmentRepo.On("Get", ctx, depID, mock.Anything).Return(department, nil)

	result, err := s.svc.GetDepartment(ctx, depID)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(department.Name, result.Name)
}

func TestGetDepartmentSuite(t *testing.T) {
	suite.Run(t, new(GetDepartmentSuite))
}
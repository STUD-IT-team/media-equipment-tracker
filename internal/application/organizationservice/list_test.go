//go:build unit

package organizationservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/organizationservice"
	"media-equipment-tracker/internal/domain"
)

type ListOrganizationSuite struct {
	suite.Suite

	svc         organizationservice.ListOrganizationService
	organizationRepo *OrganizationRepoMock
}

func (s *ListOrganizationSuite) SetupTest() {
	s.organizationRepo = new(OrganizationRepoMock)

	s.svc = organizationservice.NewListOrganizationService(s.organizationRepo)
}

func (s *ListOrganizationSuite) TestListOrganizations_Success() {
	organizations := []*domain.Organization{
		{ID: uuid.New(), Name: "Org1"},
		{ID: uuid.New(), Name: "Org2"},
	}

	ctx := context.Background()
	s.organizationRepo.On("List", ctx, mock.Anything).Return(organizations, nil)

	result, err := s.svc.ListOrganizations(ctx, "")

	s.NoError(err)
	s.Len(result, 2)
}

func (s *ListOrganizationSuite) TestListOrganizations_WithSearch() {
	organizations := []*domain.Organization{
		{ID: uuid.New(), Name: "Test Organization"},
		{ID: uuid.New(), Name: "Other Organization"},
	}

	ctx := context.Background()
	s.organizationRepo.On("List", ctx, mock.Anything).Return(organizations, nil)

	result, err := s.svc.ListOrganizations(ctx, "test")

	s.NoError(err)
	s.Len(result, 1)
	s.Equal("Test Organization", result[0].Name)
}

func TestListOrganizationSuite(t *testing.T) {
	suite.Run(t, new(ListOrganizationSuite))
}
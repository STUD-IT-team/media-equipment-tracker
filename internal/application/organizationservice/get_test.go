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

type GetOrganizationSuite struct {
	suite.Suite

	svc              organizationservice.GetOrganizationService
	organizationRepo *OrganizationRepoMock
}

func (s *GetOrganizationSuite) SetupTest() {
	s.organizationRepo = new(OrganizationRepoMock)

	s.svc = organizationservice.NewGetOrganizationService(s.organizationRepo)
}

func (s *GetOrganizationSuite) TestGetOrganization_Success() {
	orgID := uuid.New()
	organization := &domain.Organization{ID: orgID, Name: "Test Organization"}

	ctx := context.Background()
	s.organizationRepo.On("Get", ctx, orgID, mock.Anything).Return(organization, nil)

	result, err := s.svc.GetOrganization(ctx, orgID)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(organization.Name, result.Name)
}

func TestGetOrganizationSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationSuite))
}

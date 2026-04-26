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

type GetOrganizationUsersSuite struct {
	suite.Suite

	svc              organizationservice.GetOrganizationUsersService
	organizationRepo *OrganizationRepoMock
}

func (s *GetOrganizationUsersSuite) SetupTest() {
	s.organizationRepo = new(OrganizationRepoMock)

	s.svc = organizationservice.NewGetOrganizationUsersService(s.organizationRepo)
}

func (s *GetOrganizationUsersSuite) TestGetOrganizationUsers_Success() {
	orgID := uuid.New()
	users := []*domain.User{
		{ID: uuid.New(), FullName: "User1"},
		{ID: uuid.New(), FullName: "User2"},
	}
	organization := &domain.Organization{ID: orgID, Name: "Test Organization", Users: users}

	ctx := context.Background()
	s.organizationRepo.On("Get", ctx, orgID, mock.Anything).Return(organization, nil)

	result, err := s.svc.GetOrganizationUsers(ctx, orgID)

	s.NoError(err)
	s.Len(result, 2)
}

func TestGetOrganizationUsersSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationUsersSuite))
}

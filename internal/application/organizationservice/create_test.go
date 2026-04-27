//go:build unit

package organizationservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/organizationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type CreateOrganizationSuite struct {
	suite.Suite

	svc              organizationservice.CreateOrganizationService
	organizationRepo *OrganizationRepoMock
	auther           *AuthZMock
}

func (s *CreateOrganizationSuite) SetupTest() {
	s.organizationRepo = new(OrganizationRepoMock)
	s.auther = new(AuthZMock)

	s.svc = organizationservice.NewCreateOrganizationService(s.auther, s.organizationRepo)
}

func (s *CreateOrganizationSuite) TestCreateOrganization_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	req := &organizationservice.CreateOrganizationRequest{
		Name: "Test Organization",
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.organizationRepo.On("Create", ctx, mock.AnythingOfType("*domain.Organization")).Return(nil)

	result, err := s.svc.CreateOrganization(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *CreateOrganizationSuite) TestCreateOrganization_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	req := &organizationservice.CreateOrganizationRequest{
		Name: "Test Organization",
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.CreateOrganization(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CreateOrganizationSuite) TestCreateOrganization_ValidationError() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	req := &organizationservice.CreateOrganizationRequest{
		Name: "", // invalid
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.CreateOrganization(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestCreateOrganizationSuite(t *testing.T) {
	suite.Run(t, new(CreateOrganizationSuite))
}

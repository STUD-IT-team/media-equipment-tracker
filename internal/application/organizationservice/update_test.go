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
	"media-equipment-tracker/internal/domain/errs"
)

type UpdateOrganizationSuite struct {
	suite.Suite

	svc         organizationservice.UpdateOrganizationService
	organizationRepo *OrganizationRepoMock
	auther       *AuthZMock
	txManager    *TxManagerMock
}

func (s *UpdateOrganizationSuite) SetupTest() {
	s.organizationRepo = new(OrganizationRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = organizationservice.NewUpdateOrganizationService(s.auther, s.organizationRepo, s.txManager)
}

func (s *UpdateOrganizationSuite) TestUpdateOrganization_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	orgID := uuid.New()
	req := &organizationservice.UpdateOrganizationRequest{
		Name: "Updated Organization",
	}
	existingOrg := &domain.Organization{ID: orgID, Name: "Old Name"}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.organizationRepo.On("Get", ctx, orgID, mock.Anything).Return(existingOrg, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.organizationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Organization")).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	result, err := s.svc.UpdateOrganization(ctx, orgID, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *UpdateOrganizationSuite) TestUpdateOrganization_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	orgID := uuid.New()
	req := &organizationservice.UpdateOrganizationRequest{
		Name: "Updated Organization",
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.UpdateOrganization(ctx, orgID, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func TestUpdateOrganizationSuite(t *testing.T) {
	suite.Run(t, new(UpdateOrganizationSuite))
}
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

type DeleteOrganizationSuite struct {
	suite.Suite

	svc         organizationservice.DeleteOrganizationService
	organizationRepo *OrganizationRepoMock
	auther       *AuthZMock
	txManager    *TxManagerMock
}

func (s *DeleteOrganizationSuite) SetupTest() {
	s.organizationRepo = new(OrganizationRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = organizationservice.NewDeleteOrganizationService(s.auther, s.organizationRepo, s.txManager)
}

func (s *DeleteOrganizationSuite) TestDeleteOrganization_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	orgID := uuid.New()
	organization := &domain.Organization{ID: orgID, Name: "Test Organization", Users: []*domain.User{}, EquipmentInvocations: []*domain.EquipmentInvocation{}, StudioInvocations: []*domain.StudioInvocation{}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.organizationRepo.On("Get", mock.Anything, orgID, mock.Anything).Return(organization, nil)
		s.organizationRepo.On("Delete", mock.Anything, orgID).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	err := s.svc.DeleteOrganization(ctx, orgID)

	s.NoError(err)
}

func (s *DeleteOrganizationSuite) TestDeleteOrganization_HasUsers() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}
	orgID := uuid.New()
	organization := &domain.Organization{ID: orgID, Name: "Test Organization", Users: []*domain.User{{}}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(organizationservice.ErrOrganizationHasUsers).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		s.organizationRepo.On("Get", mock.Anything, orgID, mock.Anything).Return(organization, nil)
		fn(ctx)
	})

	err := s.svc.DeleteOrganization(ctx, orgID)

	s.Error(err)
	s.Equal(organizationservice.ErrOrganizationHasUsers, err)
}

func (s *DeleteOrganizationSuite) TestDeleteOrganization_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}
	orgID := uuid.New()

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	err := s.svc.DeleteOrganization(ctx, orgID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func TestDeleteOrganizationSuite(t *testing.T) {
	suite.Run(t, new(DeleteOrganizationSuite))
}
//go:build unit

package invocationservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type InvocationSuite struct {
	suite.Suite

	svc              invocationservice.InvocationService
	invocationRepo   *InvocationRepoMock
	searchRepo       *SearchInvocationRepoMock
	departmentRepo   *DepartmentRepoMock
	organizationRepo *OrganizationRepoMock
	availabilitySvc  *AvailabilityEquipmentServiceMock
	accessService    *AccessServiceMock
	txManager        *TxManagerMock
	auther           *AuthZMock
}

func (s *InvocationSuite) SetupTest() {
	s.invocationRepo = new(InvocationRepoMock)
	s.searchRepo = new(SearchInvocationRepoMock)
	s.departmentRepo = new(DepartmentRepoMock)
	s.organizationRepo = new(OrganizationRepoMock)
	s.availabilitySvc = new(AvailabilityEquipmentServiceMock)
	s.accessService = new(AccessServiceMock)
	s.txManager = new(TxManagerMock)
	s.auther = new(AuthZMock)

	s.svc = invocationservice.NewInvocationService(
		s.invocationRepo,
		s.searchRepo,
		s.departmentRepo,
		s.organizationRepo,
		s.availabilitySvc,
		nil,
		s.accessService,
		s.txManager,
		s.auther,
	)
}

func (s *InvocationSuite) TestGet_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, UserID: payload.UserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	result, err := s.svc.Get(ctx, invID)

	s.NoError(err)
	s.Equal(inv, result)
}

func (s *InvocationSuite) TestGet_AdminAccess() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	otherUserID := uuid.New()
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, UserID: otherUserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	result, err := s.svc.Get(ctx, invID)

	s.NoError(err)
	s.Equal(inv, result)
}

func (s *InvocationSuite) TestGet_NoAccess() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	otherUserID := uuid.New()
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, UserID: otherUserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	_, err := s.svc.Get(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *InvocationSuite) TestDelete_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, UserID: payload.UserID, Status: domain.InvocationCreated}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Delete", mock.Anything, invID).Return(nil)

	err := s.svc.Delete(ctx, invID)

	s.NoError(err)
}

func (s *InvocationSuite) TestDelete_NoAccess() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	otherUserID := uuid.New()
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, UserID: otherUserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Delete", mock.Anything, invID).Return(nil)

	err := s.svc.Delete(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *InvocationSuite) TestDelete_InvalidStatus() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, Status: domain.InvocationEquipmentIssued}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.Delete(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestInvocationSuite(t *testing.T) {
	suite.Run(t, new(InvocationSuite))
}

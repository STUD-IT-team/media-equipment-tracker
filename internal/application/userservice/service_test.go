//go:build unit

package userservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/userservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type UserServiceSuite struct {
	suite.Suite

	svc     userservice.UserService
	userRep *UserRepoMock
	authz   *AuthZMock
	txm     *TxManagerMock
}

func (s *UserServiceSuite) SetupTest() {
	s.userRep = new(UserRepoMock)
	s.authz = new(AuthZMock)
	s.txm = new(TxManagerMock)

	svc, err := userservice.NewUserService(s.userRep, s.authz, s.txm)
	s.NoError(err)
	s.svc = svc
}

func (s *UserServiceSuite) TestGet_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	userID := uuid.New()
	user := &domain.User{ID: userID}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, userID, mock.Anything).Return(user, nil)

	result, err := s.svc.Get(ctx, userID)

	s.NoError(err)
	s.Equal(user, result)
}

func (s *UserServiceSuite) TestGet_NoAdminRole() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	userID := uuid.New()

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.Get(ctx, userID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *UserServiceSuite) TestGetAll_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	users := []*domain.User{{ID: uuid.New()}}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("List", mock.Anything, mock.Anything).Return(users, nil)

	result, err := s.svc.GetAll(ctx)

	s.NoError(err)
	s.Equal(users, result)
}

func (s *UserServiceSuite) TestGetAll_NoAdminRole() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.GetAll(ctx)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *UserServiceSuite) TestDelete_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	userID := uuid.New()
	user := &domain.User{ID: userID, EquipmentInvocations: []*domain.EquipmentInvocation{}, StudioInvocations: []*domain.StudioInvocation{}}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, userID, mock.Anything).Return(user, nil)
	s.userRep.On("Delete", mock.Anything, userID).Return(nil)

	err := s.svc.Delete(ctx, userID)

	s.NoError(err)
}

func (s *UserServiceSuite) TestDelete_NoAdminRole() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	userID := uuid.New()

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	err := s.svc.Delete(ctx, userID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *UserServiceSuite) TestDelete_Self() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	userID := payload.UserID
	user := &domain.User{ID: userID, EquipmentInvocations: []*domain.EquipmentInvocation{}, StudioInvocations: []*domain.StudioInvocation{}}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, userID, mock.Anything).Return(user, nil)

	err := s.svc.Delete(ctx, userID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UserServiceSuite) TestDelete_WithInvocations() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	userID := uuid.New()
	user := &domain.User{ID: userID, EquipmentInvocations: []*domain.EquipmentInvocation{{ID: uuid.New()}}, StudioInvocations: []*domain.StudioInvocation{}}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, userID, mock.Anything).Return(user, nil)

	err := s.svc.Delete(ctx, userID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
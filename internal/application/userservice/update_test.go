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

type UpdateUserServiceSuite struct {
	suite.Suite

	svc     userservice.UpdateUserService
	userRep *UserRepoMock
	txm     *TxManagerMock
	auther  *AuthZMock
}

func (s *UpdateUserServiceSuite) SetupTest() {
	s.userRep = new(UserRepoMock)
	s.txm = new(TxManagerMock)
	s.auther = new(AuthZMock)

	s.svc = userservice.NewUpdateUserService(s.userRep, s.txm, s.auther)
}

func (s *UpdateUserServiceSuite) TestUpdate_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	userID := uuid.New()
	user := &domain.User{ID: userID, FullName: "old", Email: "old@example.com"}
	newName := "new name"
	newEmail := "new@example.com"

	req := &userservice.UpdateUserRequest{
		ID:       userID,
		FullName: &newName,
		Email:    &newEmail,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, userID, mock.Anything).Return(user, nil)
	s.userRep.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		updatedUser := args.Get(1).(*domain.User)
		s.Equal(newName, updatedUser.FullName)
		s.Equal(newEmail, updatedUser.Email)
	})
	s.userRep.On("Reload", mock.Anything, user, mock.Anything).Return(nil)

	result, err := s.svc.Update(ctx, req)

	s.NoError(err)
	s.Equal(user, result)
	s.Equal(newName, result.FullName)
	s.Equal(newEmail, result.Email)
}

func (s *UpdateUserServiceSuite) TestUpdate_InvalidRequest() {
	req := &userservice.UpdateUserRequest{
		ID:    uuid.New(),
		Email: stringPtr("invalid-email"),
	}

	ctx := context.Background()

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateUserServiceSuite) TestUpdate_NoAdminRole() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	req := &userservice.UpdateUserRequest{ID: uuid.New()}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *UpdateUserServiceSuite) TestUpdateNice_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	userID := uuid.New()
	user := &domain.User{ID: userID, Nice: 50}

	req := &userservice.UpdateNiceRequest{
		ID:   userID,
		Nice: 75,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, userID, mock.Anything).Return(user, nil)
	s.userRep.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		updatedUser := args.Get(1).(*domain.User)
		s.Equal(75, updatedUser.Nice)
	})

	result, err := s.svc.UpdateNice(ctx, req)

	s.NoError(err)
	s.Equal(user, result)
	s.Equal(75, result.Nice)
}

func (s *UpdateUserServiceSuite) TestUpdateNice_InvalidNice() {
	req := &userservice.UpdateNiceRequest{
		ID:   uuid.New(),
		Nice: 150, // over 100
	}

	ctx := context.Background()

	_, err := s.svc.UpdateNice(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateUserServiceSuite) TestUpdateNice_NoAdminRole() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	req := &userservice.UpdateNiceRequest{ID: uuid.New(), Nice: 50}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.UpdateNice(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func stringPtr(s string) *string { return &s }

func TestUpdateUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserServiceSuite))
}
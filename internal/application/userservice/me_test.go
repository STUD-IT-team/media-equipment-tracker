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

type MeUserServiceSuite struct {
	suite.Suite

	svc     userservice.MeUserService
	userRep *UserRepoMock
	authz   *AuthZMock
	txm     *TxManagerMock
}

func (s *MeUserServiceSuite) SetupTest() {
	s.userRep = new(UserRepoMock)
	s.authz = new(AuthZMock)
	s.txm = new(TxManagerMock)

	s.svc = userservice.NewMeUserService(s.userRep, s.authz, s.txm)
}

func (s *MeUserServiceSuite) TestMe_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	user := &domain.User{ID: payload.UserID}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, payload.UserID, mock.Anything).Return(user, nil)

	result, err := s.svc.Me(ctx)

	s.NoError(err)
	s.Equal(user, result)
}

func (s *MeUserServiceSuite) TestUpdateSelf_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	user := &domain.User{ID: payload.UserID, FullName: "old", Email: "old@example.com"}
	newName := "new name"
	newEmail := "new@example.com"

	req := &userservice.UpdateSelfRequest{
		FullName: &newName,
		Email:    &newEmail,
	}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, payload.UserID, mock.Anything).Return(user, nil)
	s.userRep.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		updatedUser := args.Get(1).(*domain.User)
		s.Equal(newName, updatedUser.FullName)
		s.Equal(newEmail, updatedUser.Email)
	})

	result, err := s.svc.UpdateSelf(ctx, req)

	s.NoError(err)
	s.Equal(user, result)
	s.Equal(newName, result.FullName)
	s.Equal(newEmail, result.Email)
}

func (s *MeUserServiceSuite) TestUpdateSelf_InvalidRequest() {
	req := &userservice.UpdateSelfRequest{
		Email: stringPtr("invalid-email"),
	}

	ctx := context.Background()

	_, err := s.svc.UpdateSelf(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *MeUserServiceSuite) TestInvocations_StudioType_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	user := &domain.User{
		ID:               payload.UserID,
		StudioInvocations: []*domain.StudioInvocation{{ID: uuid.New(), Status: domain.StudioApproved}},
	}

	req := &userservice.MyInvocationsRequest{
		Type: userservice.Studio,
	}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, payload.UserID, mock.Anything).Return(user, nil)

	result, err := s.svc.Invocations(ctx, req)

	s.NoError(err)
	s.Len(result.StudioInvocations, 1)
	s.Len(result.EquipmentInvocations, 0)
	s.Len(result.AdminStudioInvocations, 0)
	s.Len(result.AdminEquipmentInvocations, 0)
}

func (s *MeUserServiceSuite) TestInvocations_EquipmentType_WithStatuses_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	user := &domain.User{
		ID:                   payload.UserID,
		EquipmentInvocations: []*domain.EquipmentInvocation{{ID: uuid.New(), Status: domain.InvocationApproved}},
		AdminEquipmentInvocations: []*domain.EquipmentInvocation{{ID: uuid.New(), Status: domain.InvocationUnderReview}},
	}

	req := &userservice.MyInvocationsRequest{
		Type:              userservice.Equipment,
		EquipmentStatuses: []domain.EquipmentInvocationStatus{domain.InvocationApproved},
	}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, payload.UserID, mock.Anything).Return(user, nil)

	result, err := s.svc.Invocations(ctx, req)

	s.NoError(err)
	s.Len(result.EquipmentInvocations, 1)
	s.Len(result.AdminEquipmentInvocations, 0)
	s.Len(result.StudioInvocations, 0)
	s.Len(result.AdminStudioInvocations, 0)
}

func (s *MeUserServiceSuite) TestInvocations_AllType_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	user := &domain.User{
		ID:                        payload.UserID,
		EquipmentInvocations:      []*domain.EquipmentInvocation{{ID: uuid.New(), Status: domain.InvocationCreated}},
		StudioInvocations:         []*domain.StudioInvocation{{ID: uuid.New(), Status: domain.StudioCreated}},
		AdminEquipmentInvocations: []*domain.EquipmentInvocation{{ID: uuid.New(), Status: domain.InvocationCreated}},
		AdminStudioInvocations:    []*domain.StudioInvocation{{ID: uuid.New(), Status: domain.StudioCreated}},
	}

	req := &userservice.MyInvocationsRequest{
		Type: userservice.All,
	}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.userRep.On("Get", mock.Anything, payload.UserID, mock.Anything).Return(user, nil)

	result, err := s.svc.Invocations(ctx, req)

	s.NoError(err)
	s.Len(result.EquipmentInvocations, 1)
	s.Len(result.StudioInvocations, 1)
	s.Len(result.AdminEquipmentInvocations, 1)
	s.Len(result.AdminStudioInvocations, 1)
}

func (s *MeUserServiceSuite) TestInvocations_InvalidType() {
	req := &userservice.MyInvocationsRequest{
		Type: "invalid",
	}

	ctx := context.Background()
	s.authz.On("TokenPayloadFromContext", ctx).Return(domain.TokenPayload{}, nil)

	_, err := s.svc.Invocations(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestMeUserServiceSuite(t *testing.T) {
	suite.Run(t, new(MeUserServiceSuite))
}
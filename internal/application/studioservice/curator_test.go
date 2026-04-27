//go:build unit

package studioservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type CuratorStudioSuite struct {
	suite.Suite

	svc        studioservice.CuratorStudioService
	studioRepo *StudioInvocationRepoMock
	auther     *AuthZMock
	txManager  *TxManagerMock
}

func (s *CuratorStudioSuite) SetupTest() {
	s.studioRepo = new(StudioInvocationRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = studioservice.NewCuratorStudioService(
		s.studioRepo,
		s.txManager,
		s.auther,
	)
}

func (s *CuratorStudioSuite) TestBecome_Success() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:     invID,
		Status: domain.StudioCreated,
		// AdminID is nil
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.Become(ctx, invID)

	s.NoError(err)
	s.Equal(payload.UserID, *existingInv.AdminID)
	s.Equal(domain.StudioUnderReview, existingInv.Status)
}

func (s *CuratorStudioSuite) TestBecome_NotAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	invID := uuid.New()

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	err := s.svc.Become(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CuratorStudioSuite) TestBecome_AlreadyCompleted() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:     invID,
		Status: domain.StudioCompleted,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)

	err := s.svc.Become(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorStudioSuite) TestBecome_AlreadyHasAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	adminID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioCreated,
		AdminID: &adminID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)

	err := s.svc.Become(ctx, invID)

	s.Error(err)
	s.True(errs.IsEntityAlreadyExistsError(err))
}

func (s *CuratorStudioSuite) TestApprove_Success() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioUnderReview,
		AdminID: &payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.Approve(ctx, invID)

	s.NoError(err)
	s.Equal(domain.StudioApproved, existingInv.Status)
}

func (s *CuratorStudioSuite) TestApprove_NotCurator() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	otherAdminID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioUnderReview,
		AdminID: &otherAdminID, // different admin
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)

	err := s.svc.Approve(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CuratorStudioSuite) TestApprove_WrongStatus() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioCreated, // wrong status
		AdminID: &payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)

	err := s.svc.Approve(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorStudioSuite) TestRequestChanges_Success() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioUnderReview,
		AdminID: &payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.RequestChanges(ctx, invID)

	s.NoError(err)
	s.Equal(domain.StudioChangesRequired, existingInv.Status)
}

func (s *CuratorStudioSuite) TestRequestChanges_FromApprovedStatus() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioApproved,
		AdminID: &payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.RequestChanges(ctx, invID)

	s.NoError(err)
	s.Equal(domain.StudioChangesRequired, existingInv.Status)
}

func (s *CuratorStudioSuite) TestRequestChanges_WrongStatus() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioCreated, // wrong status
		AdminID: &payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)

	err := s.svc.RequestChanges(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorStudioSuite) TestComplete_Success() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioApproved,
		AdminID: &payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.Complete(ctx, invID)

	s.NoError(err)
	s.Equal(domain.StudioCompleted, existingInv.Status)
}

func (s *CuratorStudioSuite) TestComplete_WrongStatus() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		Status:  domain.StudioUnderReview, // wrong status
		AdminID: &payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)

	err := s.svc.Complete(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorStudioSuite) TestCancel_Success_ByOwner() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:     invID,
		UserID: payload.UserID, // owner
		Status: domain.StudioCreated,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.Cancel(ctx, invID)

	s.NoError(err)
	s.Equal(domain.StudioCancelled, existingInv.Status)
}

func (s *CuratorStudioSuite) TestCancel_Success_ByAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:     invID,
		UserID: uuid.New(), // different user
		Status: domain.StudioUnderReview,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.Cancel(ctx, invID)

	s.NoError(err)
	s.Equal(domain.StudioCancelled, existingInv.Status)
}

func (s *CuratorStudioSuite) TestCancel_Success_ByCurator() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:      invID,
		UserID:  uuid.New(),      // different user
		AdminID: &payload.UserID, // current user is curator
		Status:  domain.StudioUnderReview,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	err := s.svc.Cancel(ctx, invID)

	s.NoError(err)
	s.Equal(domain.StudioCancelled, existingInv.Status)
}

func (s *CuratorStudioSuite) TestCancel_AlreadyCompleted() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:     invID,
		UserID: payload.UserID,
		Status: domain.StudioCompleted,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything).Return(existingInv, nil)

	err := s.svc.Cancel(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorStudioSuite) TestCancel_NotAuthorized() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	invID := uuid.New()
	existingInv := &domain.StudioInvocation{
		ID:     invID,
		UserID: uuid.New(), // different user
		Status: domain.StudioCreated,
		// No admin assigned
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything).Return(existingInv, nil)

	err := s.svc.Cancel(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func TestCuratorStudioSuite(t *testing.T) {
	suite.Run(t, new(CuratorStudioSuite))
}

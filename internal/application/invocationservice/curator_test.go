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

type CuratorInvocationSuite struct {
	suite.Suite

	svc                invocationservice.CuratorInvocationService
	invocationRepo     *InvocationRepoMock
	updateEquipmentSvc *UpdateEquipmentServiceMock
	txManager          *TxManagerMock
	auther             *AuthZMock
}

func (s *CuratorInvocationSuite) SetupTest() {
	s.invocationRepo = new(InvocationRepoMock)
	s.updateEquipmentSvc = new(UpdateEquipmentServiceMock)
	s.txManager = new(TxManagerMock)
	s.auther = new(AuthZMock)

	s.svc = invocationservice.NewCuratorInvocationService(
		s.invocationRepo,
		s.updateEquipmentSvc,
		s.txManager,
		s.auther,
	)
}

func (s *CuratorInvocationSuite) TestBecome_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, Status: domain.InvocationCreated}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil).Run(func(args mock.Arguments) {
		updatedInv := args.Get(1).(*domain.EquipmentInvocation)
		s.Equal(&payload.UserID, updatedInv.AdminID)
		s.Equal(domain.InvocationUnderReview, updatedInv.Status)
	})

	err := s.svc.Become(ctx, invID)

	s.NoError(err)
}

func (s *CuratorInvocationSuite) TestBecome_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	invID := uuid.New()

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	err := s.svc.Become(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CuratorInvocationSuite) TestBecome_AlreadyHasAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	otherUserID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &otherUserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.Become(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorInvocationSuite) TestApprove_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &payload.UserID, Status: domain.InvocationUnderReview}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil).Run(func(args mock.Arguments) {
		updatedInv := args.Get(1).(*domain.EquipmentInvocation)
		s.Equal(domain.InvocationApproved, updatedInv.Status)
	})

	err := s.svc.Approve(ctx, invID)

	s.NoError(err)
}

func (s *CuratorInvocationSuite) TestApprove_NotAdminForThis() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	otherUserID := uuid.New()
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &otherUserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.Approve(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CuratorInvocationSuite) TestApprove_NotUnderReview() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &payload.UserID, Status: domain.InvocationChangesRequired}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.Approve(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorInvocationSuite) TestRequestChanges_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &payload.UserID, Status: domain.InvocationApproved}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil).Run(func(args mock.Arguments) {
		updatedInv := args.Get(1).(*domain.EquipmentInvocation)
		s.Equal(domain.InvocationChangesRequired, updatedInv.Status)
	})

	err := s.svc.RequestChanges(ctx, invID)

	s.NoError(err)
}

func (s *CuratorInvocationSuite) TestRequestChanges_NoAdminForThis() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	otherUserID := uuid.New()
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &otherUserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.RequestChanges(ctx, invID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CuratorInvocationSuite) TestRequestChanges_NotUnderReview() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &payload.UserID, Status: domain.InvocationEquipmentIssued}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.RequestChanges(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorInvocationSuite) TestIssue_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	eqID := uuid.New()
	inv := &domain.EquipmentInvocation{
		ID:      invID,
		AdminID: &payload.UserID,
		Status:  domain.InvocationEquipmentIssued,
		Equipment: []*domain.EquipmentInInvocation{
			{EquipmentID: eqID, Status: domain.EquipmentNotIssued, Equipment: &domain.Equipment{ID: eqID, Status: domain.EquipmentStatusAvailable}},
		},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.updateEquipmentSvc.On("UpdateEquipment", mock.Anything, mock.AnythingOfType("*equipmentservice.UpdateEquipmentRequest")).Return(&domain.Equipment{ID: eqID, Status: domain.EquipmentStatusIssued}, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil)

	err := s.svc.Issue(ctx, invID, eqID)

	s.NoError(err)
	s.Equal(domain.InvocationEquipmentIssued, inv.Status)
	s.Equal(domain.EquipmentIssued, inv.Equipment[0].Status)
}

func (s *CuratorInvocationSuite) TestIssue_NoAdminForThis() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	otherUserID := uuid.New()
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &otherUserID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.Issue(ctx, invID, uuid.New())

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *CuratorInvocationSuite) TestIssue_EquipmentNotAvailable() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	eqID := uuid.New()
	inv := &domain.EquipmentInvocation{
		ID:      invID,
		AdminID: &payload.UserID,
		Status:  domain.InvocationApproved,
		Equipment: []*domain.EquipmentInInvocation{
			{EquipmentID: eqID, Status: domain.EquipmentNotIssued, Equipment: &domain.Equipment{ID: eqID, Status: domain.EquipmentStatusUnavailable}},
		},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	err := s.svc.Issue(ctx, invID, eqID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CuratorInvocationSuite) TestReturn_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	eqID := uuid.New()
	inv := &domain.EquipmentInvocation{
		ID:      invID,
		AdminID: &payload.UserID,
		Status:  domain.InvocationEquipmentIssued,
		Equipment: []*domain.EquipmentInInvocation{
			{EquipmentID: eqID, Status: domain.EquipmentIssued, Equipment: &domain.Equipment{ID: eqID, Status: domain.EquipmentStatusIssued}},
		},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.updateEquipmentSvc.On("UpdateEquipment", mock.Anything, mock.AnythingOfType("*equipmentservice.UpdateEquipmentRequest")).Return(&domain.Equipment{ID: eqID, Status: domain.EquipmentStatusAvailable}, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil)

	err := s.svc.Return(ctx, invID, eqID)

	s.NoError(err)
	s.Equal(domain.EquipmentReturned, inv.Equipment[0].Status)
	s.Equal(domain.InvocationEquipmentReturned, inv.Status)
}

func (s *CuratorInvocationSuite) TestComplete_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &payload.UserID, Status: domain.InvocationEquipmentReturned}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil).Run(func(args mock.Arguments) {
		updatedInv := args.Get(1).(*domain.EquipmentInvocation)
		s.Equal(domain.InvocationCompleted, updatedInv.Status)
	})

	err := s.svc.Complete(ctx, invID)

	s.NoError(err)
}

func (s *CuratorInvocationSuite) TestCancel_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &payload.UserID, Status: domain.InvocationUnderReview}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil).Run(func(args mock.Arguments) {
		updatedInv := args.Get(1).(*domain.EquipmentInvocation)
		s.Equal(domain.InvocationCancelled, updatedInv.Status)
	})
	err := s.svc.Cancel(ctx, invID)

	s.NoError(err)
}

func (s *CuratorInvocationSuite) TestCancel_EquipmentIssued() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	inv := &domain.EquipmentInvocation{ID: invID, AdminID: &payload.UserID, Status: domain.InvocationEquipmentIssued}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	err := s.svc.Cancel(ctx, invID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestCuratorInvocationSuite(t *testing.T) {
	suite.Run(t, new(CuratorInvocationSuite))
}

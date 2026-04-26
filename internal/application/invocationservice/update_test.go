//go:build unit

package invocationservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type UpdateInvocationSuite struct {
	suite.Suite

	svc             invocationservice.UpdateInvocationService
	invocationRepo  *InvocationRepoMock
	availabilitySvc *AvailabilityEquipmentServiceMock
	accessService   *AccessServiceMock
	txManager       *TxManagerMock
	auther          *AuthZMock
}

func (s *UpdateInvocationSuite) SetupTest() {
	s.invocationRepo = new(InvocationRepoMock)
	s.availabilitySvc = new(AvailabilityEquipmentServiceMock)
	s.accessService = new(AccessServiceMock)
	s.txManager = new(TxManagerMock)
	s.auther = new(AuthZMock)

	s.svc = invocationservice.NewUpdateInvocationService(
		s.invocationRepo,
		s.availabilitySvc,
		s.accessService,
		s.txManager,
		s.auther,
	)
}

func (s *UpdateInvocationSuite) TestUpdate_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	newName := "New Event Name"
	req := &invocationservice.UpdateInvocationRequest{
		ID:        invID,
		EventName: &newName,
	}

	inv := &domain.EquipmentInvocation{
		ID:        invID,
		UserID:    payload.UserID,
		EventName: "Old Name",
		Status:    domain.InvocationCreated,
		Equipment: []*domain.EquipmentInInvocation{},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil)
	s.invocationRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation"), mock.Anything).Return(nil)
	result, err := s.svc.Update(ctx, req)

	s.NoError(err)
	s.Equal(newName, result.EventName)
}

func (s *UpdateInvocationSuite) TestUpdate_InvalidRequest() {
	req := &invocationservice.UpdateInvocationRequest{
		ID:        uuid.New(),
		EventName: stringPtr(""), // too short
	}

	ctx := context.Background()

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateInvocationSuite) TestUpdate_NoPermission() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	otherUserID := uuid.New()
	invID := uuid.New()
	req := &invocationservice.UpdateInvocationRequest{
		ID:        invID,
		EventName: stringPtr("New Name"),
	}

	inv := &domain.EquipmentInvocation{
		ID:     invID,
		UserID: otherUserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *UpdateInvocationSuite) TestUpdate_AlreadyIssued() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	req := &invocationservice.UpdateInvocationRequest{
		ID:        invID,
		EventName: stringPtr("New Name"),
	}

	inv := &domain.EquipmentInvocation{
		ID:     invID,
		UserID: payload.UserID,
		Status: domain.InvocationEquipmentIssued,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateInvocationSuite) TestUpdate_WithEquipment() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	eqID := uuid.New()
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := &invocationservice.UpdateInvocationRequest{
		ID:           invID,
		StartTime:    &start,
		EndTime:      &end,
		EquipmentIDs: []uuid.UUID{eqID},
	}

	inv := &domain.EquipmentInvocation{
		ID:        invID,
		UserID:    payload.UserID,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Minute),
		Status:    domain.InvocationCreated,
		Equipment: []*domain.EquipmentInInvocation{},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("equipmentservice.EquipmentAvailabilityRequest")).Return(&equipmentservice.EquipmentAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToEquipment", mock.Anything, mock.AnythingOfType("*accessservice.HaveEquipmentAccessRequest")).Return(true, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil)
	s.invocationRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation"), mock.Anything).Return(nil)

	result, err := s.svc.Update(ctx, req)

	s.NoError(err)
	s.Equal(start, result.StartTime)
	s.Len(result.Equipment, 1)
	s.Equal(eqID, result.Equipment[0].EquipmentID)
}

func (s *UpdateInvocationSuite) TestUpdate_EquipmentNotAvailable() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	eqID := uuid.New()
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := &invocationservice.UpdateInvocationRequest{
		ID:           invID,
		StartTime:    &start,
		EndTime:      &end,
		EquipmentIDs: []uuid.UUID{eqID},
	}

	inv := &domain.EquipmentInvocation{
		ID:        invID,
		UserID:    payload.UserID,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Minute),
		Status:    domain.InvocationCreated,
		Equipment: []*domain.EquipmentInInvocation{},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("equipmentservice.EquipmentAvailabilityRequest")).Return(&equipmentservice.EquipmentAvailabilityResponse{Available: false}, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateInvocationSuite) TestUpdate_NoAccess() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	invID := uuid.New()
	eqID := uuid.New()
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	req := &invocationservice.UpdateInvocationRequest{
		ID:           invID,
		StartTime:    &start,
		EndTime:      &end,
		EquipmentIDs: []uuid.UUID{eqID},
	}

	inv := &domain.EquipmentInvocation{
		ID:     invID,
		UserID: payload.UserID,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("equipmentservice.EquipmentAvailabilityRequest")).Return(&equipmentservice.EquipmentAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToEquipment", mock.Anything, mock.AnythingOfType("*accessservice.HaveEquipmentAccessRequest")).Return(false, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsEquipmentAccessError(err))
}

func (s *UpdateInvocationSuite) TestUpdateStatus_Success() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}, UserID: uuid.New()}
	invID := uuid.New()
	comment := "Approved"
	req := &invocationservice.UpdateInvocationStatusRequest{
		ID:             invID,
		Status:         domain.InvocationApproved,
		CuratorComment: &comment,
	}

	inv := &domain.EquipmentInvocation{ID: invID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	s.invocationRepo.On("Get", mock.Anything, invID, mock.Anything).Return(inv, nil)
	s.invocationRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil).Run(func(args mock.Arguments) {
		updatedInv := args.Get(1).(*domain.EquipmentInvocation)
		s.Equal(domain.InvocationApproved, updatedInv.Status)
		s.Equal(comment, updatedInv.CuratorComment)
	})

	result, err := s.svc.UpdateStatus(ctx, req)

	s.NoError(err)
	s.Equal(domain.InvocationApproved, result.Status)
}

func (s *UpdateInvocationSuite) TestUpdateStatus_NoAdmin() {
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}, UserID: uuid.New()}
	req := &invocationservice.UpdateInvocationStatusRequest{
		ID:     uuid.New(),
		Status: domain.InvocationApproved,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.UpdateStatus(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func stringPtr(s string) *string {
	return &s
}

func TestUpdateInvocationSuite(t *testing.T) {
	suite.Run(t, new(UpdateInvocationSuite))
}

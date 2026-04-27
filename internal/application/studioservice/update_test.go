//go:build unit

package studioservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type UpdateStudioSuite struct {
	suite.Suite

	svc              studioservice.UpdateStudioService
	studioRepo       *StudioInvocationRepoMock
	departmentRepo   *DepartmentRepoMock
	organizationRepo *OrganizationRepoMock
	availabilitySvc  *AvailabilityStudioServiceMock
	accessService    *AccessServiceMock
	auther           *AuthZMock
	txManager        *TxManagerMock
}

func (s *UpdateStudioSuite) SetupTest() {
	s.studioRepo = new(StudioInvocationRepoMock)
	s.departmentRepo = new(DepartmentRepoMock)
	s.organizationRepo = new(OrganizationRepoMock)
	s.availabilitySvc = new(AvailabilityStudioServiceMock)
	s.accessService = new(AccessServiceMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = studioservice.NewUpdateStudioService(
		s.studioRepo,
		s.departmentRepo,
		s.organizationRepo,
		s.availabilitySvc,
		s.accessService,
		s.txManager,
		s.auther,
	)
}

func (s *UpdateStudioSuite) TestUpdate_Success_FullUpdate() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	invID := uuid.New()
	orgID := uuid.New()
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:                  invID,
		EventName:           stringPtr("Updated Event"),
		ShootingDescription: stringPtr("Updated shooting"),
		StartTime:           &start,
		EndTime:             &end,
		NeedsChromakey:      boolPtr(false),
		NeedsCyclorama:      boolPtr(true),
		NeedsBlackFabric:    boolPtr(false),
	}

	existingInv := &domain.StudioInvocation{
		ID:                  invID,
		UserID:              payload.UserID,
		EventName:           "Original Event",
		ShootingDescription: "Original shooting",
		StartTime:           time.Now(),
		EndTime:             time.Now().Add(time.Hour),
		Status:              domain.StudioCreated,
		OrganizationID:      &orgID,
		Organization:        &domain.Organization{ID: orgID},
		NeedsChromakey:      true,
		NeedsCyclorama:      false,
		NeedsBlackFabric:    true,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.organizationRepo.On("Reload", ctx, mock.AnythingOfType("*domain.Organization"), mock.Anything).Return(nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToStudio", mock.Anything, mock.AnythingOfType("*accessservice.HaveStudioAccessRequest")).Return(true, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)
	s.studioRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation"), mock.Anything).Return(nil)

	result, err := s.svc.Update(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(*req.EventName, result.EventName)
	s.Equal(*req.ShootingDescription, result.ShootingDescription)
	s.Equal(*req.StartTime, result.StartTime)
	s.Equal(*req.EndTime, result.EndTime)
	s.Equal(*req.NeedsChromakey, result.NeedsChromakey)
	s.Equal(*req.NeedsCyclorama, result.NeedsCyclorama)
	s.Equal(*req.NeedsBlackFabric, result.NeedsBlackFabric)
}

func (s *UpdateStudioSuite) TestUpdate_Success_PartialUpdate() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	invID := uuid.New()
	depID := uuid.New()
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:        invID,
		EventName: stringPtr("Updated Event"),
		// Other fields nil - partial update
	}

	existingInv := &domain.StudioInvocation{
		ID:                  invID,
		UserID:              payload.UserID,
		EventName:           "Original Event",
		ShootingDescription: "Original shooting",
		StartTime:           time.Now(),
		EndTime:             time.Now().Add(time.Hour),
		Status:              domain.StudioCreated,
		DepartmentID:        &depID,
		Department:          &domain.Department{ID: depID},
		NeedsChromakey:      true,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.departmentRepo.On("Reload", ctx, mock.AnythingOfType("*domain.Department"), mock.Anything).Return(nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToStudio", mock.Anything, mock.AnythingOfType("*accessservice.HaveStudioAccessRequest")).Return(true, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)
	s.studioRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation"), mock.Anything).Return(nil)

	result, err := s.svc.Update(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(*req.EventName, result.EventName)
	// Other fields should remain unchanged
	s.Equal("Original shooting", result.ShootingDescription)
	s.Equal(true, result.NeedsChromakey)
}

func (s *UpdateStudioSuite) TestUpdate_StatusChangesRequiredToUnderReview() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	invID := uuid.New()
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:        invID,
		EventName: stringPtr("Updated Event"),
	}

	existingInv := &domain.StudioInvocation{
		ID:        invID,
		UserID:    payload.UserID,
		EventName: "Original Event",
		Status:    domain.StudioChangesRequired,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToStudio", mock.Anything, mock.AnythingOfType("*accessservice.HaveStudioAccessRequest")).Return(true, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)
	s.studioRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation"), mock.Anything).Return(nil)

	result, err := s.svc.Update(ctx, req)

	s.NoError(err)
	s.Equal(domain.StudioUnderReview, result.Status)
}

func (s *UpdateStudioSuite) TestUpdate_InvalidRequest() {
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:        uuid.New(),
		EventName: stringPtr(""), // too short
	}

	ctx := context.Background()

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateStudioSuite) TestUpdate_NotOwnerAndNotAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	invID := uuid.New()
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:        invID,
		EventName: stringPtr("Updated Event"),
	}

	existingInv := &domain.StudioInvocation{
		ID:     invID,
		UserID: uuid.New(), // different user
		Status: domain.StudioCreated,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingInv, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *UpdateStudioSuite) TestUpdate_AlreadyCompleted() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:        invID,
		EventName: stringPtr("Updated Event"),
	}

	existingInv := &domain.StudioInvocation{
		ID:     invID,
		UserID: payload.UserID,
		Status: domain.StudioCompleted,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingInv, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateStudioSuite) TestUpdate_TimeNotAvailable() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	invID := uuid.New()
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:        invID,
		StartTime: &start,
		EndTime:   &end,
	}

	existingInv := &domain.StudioInvocation{
		ID:        invID,
		UserID:    payload.UserID,
		Status:    domain.StudioCreated,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: false}, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateStudioSuite) TestUpdate_NoAccessToStudio() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	invID := uuid.New()
	orgID := uuid.New()
	req := &studioservice.UpdateStudioInvocationRequest{
		ID:        invID,
		EventName: stringPtr("Updated Event"),
	}

	existingInv := &domain.StudioInvocation{
		ID:             invID,
		UserID:         payload.UserID,
		Status:         domain.StudioCreated,
		OrganizationID: &orgID,
		Organization:   &domain.Organization{ID: orgID},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(existingInv, nil)
	s.organizationRepo.On("Reload", ctx, mock.AnythingOfType("*domain.Organization"), mock.Anything).Return(nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToStudio", mock.Anything, mock.AnythingOfType("*accessservice.HaveStudioAccessRequest")).Return(false, nil)

	_, err := s.svc.Update(ctx, req)

	s.Error(err)
	s.True(errs.IsStudioAccessError(err))
}

func (s *UpdateStudioSuite) TestUpdateStatus_Success_Approve() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	req := &studioservice.UpdateStudioInvocationStatusRequest{
		ID:     invID,
		Status: domain.StudioApproved,
	}

	existingInv := &domain.StudioInvocation{
		ID:     invID,
		Status: domain.StudioUnderReview,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	result, err := s.svc.UpdateStatus(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(domain.StudioApproved, result.Status)
}

func (s *UpdateStudioSuite) TestUpdateStatus_Success_RequestChanges() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	invID := uuid.New()
	comment := "Please update the description"
	req := &studioservice.UpdateStudioInvocationStatusRequest{
		ID:             invID,
		Status:         domain.StudioChangesRequired,
		CuratorComment: &comment,
	}

	existingInv := &domain.StudioInvocation{
		ID:     invID,
		Status: domain.StudioUnderReview,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, invID, mock.Anything).Return(existingInv, nil)
	s.studioRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)

	result, err := s.svc.UpdateStatus(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(domain.StudioChangesRequired, result.Status)
	s.Equal(comment, *result.CuratorComment)
}

func (s *UpdateStudioSuite) TestUpdateStatus_InvalidRequest() {
	req := &studioservice.UpdateStudioInvocationStatusRequest{
		ID:             uuid.New(),
		Status:         domain.StudioApproved,
		CuratorComment: stringPtr(""), // too short
	}

	ctx := context.Background()

	_, err := s.svc.UpdateStatus(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *UpdateStudioSuite) TestUpdateStatus_NotAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	req := &studioservice.UpdateStudioInvocationStatusRequest{
		ID:     uuid.New(),
		Status: domain.StudioApproved,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.UpdateStatus(ctx, req)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func TestUpdateStudioSuite(t *testing.T) {
	suite.Run(t, new(UpdateStudioSuite))
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

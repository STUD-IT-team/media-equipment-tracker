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

type CreateStudioSuite struct {
	suite.Suite

	svc              studioservice.CreateStudioService
	studioRepo       *StudioInvocationRepoMock
	departmentRepo   *DepartmentRepoMock
	organizationRepo *OrganizationRepoMock
	availabilitySvc  *AvailabilityStudioServiceMock
	accessService    *AccessServiceMock
	auther           *AuthZMock
	txManager        *TxManagerMock
}

func (s *CreateStudioSuite) SetupTest() {
	s.studioRepo = new(StudioInvocationRepoMock)
	s.departmentRepo = new(DepartmentRepoMock)
	s.organizationRepo = new(OrganizationRepoMock)
	s.availabilitySvc = new(AvailabilityStudioServiceMock)
	s.accessService = new(AccessServiceMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = studioservice.NewCreateStudioService(
		s.studioRepo,
		s.departmentRepo,
		s.organizationRepo,
		s.availabilitySvc,
		s.accessService,
		s.txManager,
		s.auther,
	)
}

func (s *CreateStudioSuite) TestCreate_Success_WithOrganization() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	orgID := uuid.New()
	req := &studioservice.CreateStudioInvocationRequest{
		EventName:           "Test Event",
		ShootingDescription: "Test shooting",
		StartTime:           start,
		EndTime:             end,
		OrganizationID:      &orgID,
		NeedsChromakey:      true,
		NeedsCyclorama:      true,
		NeedsBlackFabric:    false,
	}

	org := &domain.Organization{ID: orgID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	s.organizationRepo.On("Get", mock.Anything, orgID, mock.Anything).Return(org, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToStudio", mock.Anything, mock.AnythingOfType("*accessservice.HaveStudioAccessRequest")).Return(true, nil)
	s.studioRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)
	s.studioRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation"), mock.Anything).Return(nil)

	result, err := s.svc.Create(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.EventName, result.EventName)
	s.Equal(req.ShootingDescription, result.ShootingDescription)
	s.Equal(payload.UserID, result.UserID)
	s.Equal(domain.StudioCreated, result.Status)
	s.Equal(orgID, *result.OrganizationID)
	s.Equal(req.NeedsChromakey, result.NeedsChromakey)
	s.Equal(req.NeedsCyclorama, result.NeedsCyclorama)
	s.Equal(req.NeedsBlackFabric, result.NeedsBlackFabric)
}

func (s *CreateStudioSuite) TestCreate_Success_WithDepartment() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	depID := uuid.New()
	req := &studioservice.CreateStudioInvocationRequest{
		EventName:           "Test Event",
		ShootingDescription: "Test shooting",
		StartTime:           start,
		EndTime:             end,
		DepartmentID:        &depID,
		NeedsChromakey:      false,
		NeedsCyclorama:      false,
		NeedsBlackFabric:    true,
	}

	dep := &domain.Department{ID: depID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.departmentRepo.On("Get", mock.Anything, depID, mock.Anything).Return(dep, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToStudio", mock.Anything, mock.AnythingOfType("*accessservice.HaveStudioAccessRequest")).Return(true, nil)
	s.studioRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation")).Return(nil)
	s.studioRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.StudioInvocation"), mock.Anything).Return(nil)

	result, err := s.svc.Create(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(depID, *result.DepartmentID)
}

func (s *CreateStudioSuite) TestCreate_InvalidRequest() {
	req := &studioservice.CreateStudioInvocationRequest{
		EventName: "", // invalid: too short
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now(), // invalid: before start
	}

	ctx := context.Background()

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CreateStudioSuite) TestCreate_TimeNotAvailable() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	orgID := uuid.New()
	req := &studioservice.CreateStudioInvocationRequest{
		EventName:           "Test Event",
		ShootingDescription: "Test shooting",
		StartTime:           start,
		EndTime:             end,
		OrganizationID:      &orgID,
		NeedsChromakey:      true,
	}

	org := &domain.Organization{ID: orgID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.organizationRepo.On("Get", mock.Anything, orgID, mock.Anything).Return(org, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: false}, nil)

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CreateStudioSuite) TestCreate_NoAccessToStudio() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	depID := uuid.New()
	req := &studioservice.CreateStudioInvocationRequest{
		EventName:           "Test Event",
		ShootingDescription: "Test shooting",
		StartTime:           start,
		EndTime:             end,
		DepartmentID:        &depID,
		NeedsChromakey:      true,
		NeedsCyclorama:      true,
		NeedsBlackFabric:    true,
	}

	dep := &domain.Department{ID: depID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.departmentRepo.On("Get", mock.Anything, depID, mock.Anything).Return(dep, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("studioservice.AvailabilityStudioInvocationRequest")).Return(&studioservice.StudioInvocationAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToStudio", mock.Anything, mock.AnythingOfType("*accessservice.HaveStudioAccessRequest")).Return(false, nil)

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsStudioAccessError(err))
}

func (s *CreateStudioSuite) TestCreate_NoOrgOrDepSpecified() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	req := &studioservice.CreateStudioInvocationRequest{
		EventName:           "Test Event",
		ShootingDescription: "Test shooting",
		StartTime:           time.Now().Add(time.Hour),
		EndTime:             time.Now().Add(2 * time.Hour),
		// Neither OrganizationID nor DepartmentID set
		NeedsChromakey: true,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CreateStudioSuite) TestCreate_BothOrgAndDepSpecified() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	orgID := uuid.New()
	depID := uuid.New()
	req := &studioservice.CreateStudioInvocationRequest{
		EventName:           "Test Event",
		ShootingDescription: "Test shooting",
		StartTime:           time.Now().Add(time.Hour),
		EndTime:             time.Now().Add(2 * time.Hour),
		OrganizationID:      &orgID, // Both set
		DepartmentID:        &depID,
		NeedsChromakey:      true,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestCreateStudioSuite(t *testing.T) {
	suite.Run(t, new(CreateStudioSuite))
}

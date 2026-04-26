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

type CreateInvocationSuite struct {
	suite.Suite

	svc              invocationservice.CreateInvocationService
	invocationRepo   *InvocationRepoMock
	departmentRepo   *DepartmentRepoMock
	organizationRepo *OrganizationRepoMock
	accessService    *AccessServiceMock
	availabilitySvc  *AvailabilityEquipmentServiceMock
	auther           *AuthZMock
	txManager        *TxManagerMock
}

func (s *CreateInvocationSuite) SetupTest() {
	s.invocationRepo = new(InvocationRepoMock)
	s.departmentRepo = new(DepartmentRepoMock)
	s.organizationRepo = new(OrganizationRepoMock)
	s.accessService = new(AccessServiceMock)
	s.availabilitySvc = new(AvailabilityEquipmentServiceMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)

	s.svc = invocationservice.NewCreateInvocationService(
		s.invocationRepo,
		s.departmentRepo,
		s.organizationRepo,
		s.accessService,
		s.availabilitySvc,
		s.auther,
		s.txManager,
	)
}

func (s *CreateInvocationSuite) TestCreate_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	orgID := uuid.New()
	eqID := uuid.New()
	req := &invocationservice.CreateInvocationRequest{
		EventName:      "Test Event",
		StartTime:      start,
		EndTime:        end,
		OrganizationID: &orgID,
		EquipmentIDs:   []uuid.UUID{eqID},
	}

	org := &domain.Organization{ID: orgID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	s.organizationRepo.On("Get", mock.Anything, orgID, mock.Anything).Return(org, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("equipmentservice.EquipmentAvailabilityRequest")).Return(&equipmentservice.EquipmentAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToEquipment", mock.Anything, mock.AnythingOfType("*accessservice.HaveEquipmentAccessRequest")).Return(true, nil)
	s.invocationRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil)
	s.invocationRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation"), mock.Anything).Return(nil)
	result, err := s.svc.Create(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.EventName, result.EventName)
	s.Equal(payload.UserID, result.UserID)
	s.Equal(domain.InvocationCreated, result.Status)
}

func (s *CreateInvocationSuite) TestCreate_InvalidRequest() {
	req := &invocationservice.CreateInvocationRequest{
		EventName: "", // invalid
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now(),
	}

	ctx := context.Background()

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CreateInvocationSuite) TestCreate_EquipmentNotAvailable() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	orgID := uuid.New()
	eqID := uuid.New()
	req := &invocationservice.CreateInvocationRequest{
		EventName:      "Test Event",
		StartTime:      start,
		EndTime:        end,
		OrganizationID: &orgID,
		EquipmentIDs:   []uuid.UUID{eqID},
	}

	org := &domain.Organization{ID: orgID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	s.organizationRepo.On("Get", mock.Anything, orgID, mock.Anything).Return(org, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("equipmentservice.EquipmentAvailabilityRequest")).Return(&equipmentservice.EquipmentAvailabilityResponse{Available: false}, nil)

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *CreateInvocationSuite) TestCreate_NoAccessToEquipment() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	orgID := uuid.New()
	eqID := uuid.New()
	req := &invocationservice.CreateInvocationRequest{
		EventName:      "Test Event",
		StartTime:      start,
		EndTime:        end,
		OrganizationID: &orgID,
		EquipmentIDs:   []uuid.UUID{eqID},
	}

	org := &domain.Organization{ID: orgID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.organizationRepo.On("Get", mock.Anything, orgID, mock.Anything).Return(org, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("equipmentservice.EquipmentAvailabilityRequest")).Return(&equipmentservice.EquipmentAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToEquipment", mock.Anything, mock.AnythingOfType("*accessservice.HaveEquipmentAccessRequest")).Return(false, nil)

	_, err := s.svc.Create(ctx, req)

	s.Error(err)
	s.True(errs.IsEquipmentAccessError(err))
}

func (s *CreateInvocationSuite) TestCreate_WithDepartment() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	depID := uuid.New()
	eqID := uuid.New()
	req := &invocationservice.CreateInvocationRequest{
		EventName:    "Test Event",
		StartTime:    start,
		EndTime:      end,
		DepartmentID: &depID,
		EquipmentIDs: []uuid.UUID{eqID},
	}

	dep := &domain.Department{ID: depID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.departmentRepo.On("Get", mock.Anything, depID, mock.Anything).Return(dep, nil)
	s.availabilitySvc.On("Availability", mock.Anything, mock.AnythingOfType("equipmentservice.EquipmentAvailabilityRequest")).Return(&equipmentservice.EquipmentAvailabilityResponse{Available: true}, nil)
	s.accessService.On("HaveAccessToEquipment", mock.Anything, mock.AnythingOfType("*accessservice.HaveEquipmentAccessRequest")).Return(true, nil)
	s.invocationRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation")).Return(nil)
	s.invocationRepo.On("Reload", mock.Anything, mock.AnythingOfType("*domain.EquipmentInvocation"), mock.Anything).Return(nil)

	result, err := s.svc.Create(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(depID, *result.DepartmentID)
}

func TestCreateInvocationSuite(t *testing.T) {
	suite.Run(t, new(CreateInvocationSuite))
}

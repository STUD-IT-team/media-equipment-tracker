//go:build unit

package accessservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/accessservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type AccessServiceSuite struct {
	suite.Suite

	svc    accessservice.AccessService
	auther *AuthZMock
}

func (s *AccessServiceSuite) SetupTest() {
	s.auther = new(AuthZMock)
	s.svc = accessservice.NewAccessService(s.auther)
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Success_Trainee() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	userID := payload.UserID
	eqID := uuid.New()
	depID := uuid.New()

	dep := &domain.Department{
		ID: depID,
		Users: []*domain.UserDepartment{
			{UserID: userID, Role: domain.RoleTrainee, User: &domain.User{Nice: 100}},
		},
		Equipment: []*domain.Equipment{
			{ID: eqID, AvailableToTrainee: true},
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	access, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.NoError(err)
	s.True(access)
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Success_Activist() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	userID := payload.UserID
	eqID := uuid.New()
	depID := uuid.New()

	dep := &domain.Department{
		ID: depID,
		Users: []*domain.UserDepartment{
			{UserID: userID, Role: domain.RoleActivist, User: &domain.User{Nice: 70}},
		},
		Equipment: []*domain.Equipment{
			{ID: eqID, AvailableToTrainee: false},
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	access, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.NoError(err)
	s.True(access)
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Organization_Success() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	userID := payload.UserID
	eqID := uuid.New()
	orgID := uuid.New()

	org := &domain.Organization{
		ID: orgID,
		Users: []*domain.User{
			{ID: userID, Nice: 100},
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID:  eqID,
		Organization: org,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	access, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.NoError(err)
	s.True(access)
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Fail_UsersNotLoaded() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	eqID := uuid.New()

	dep := &domain.Department{
		Users: nil, // not loaded
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Fail_EquipmentNotLoaded() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	eqID := uuid.New()

	dep := &domain.Department{
		Users:     []*domain.UserDepartment{},
		Equipment: nil, // not loaded
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Fail_UserNotInDepartment() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	eqID := uuid.New()

	dep := &domain.Department{
		Users: []*domain.UserDepartment{
			{UserID: uuid.New()}, // different user
		},
		Equipment: []*domain.Equipment{
			{ID: eqID},
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsEquipmentAccessError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Fail_EquipmentNotOwned() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	userID := payload.UserID
	eqID := uuid.New()

	dep := &domain.Department{
		Users: []*domain.UserDepartment{
			{UserID: userID, Role: domain.RoleTrainee},
		},
		Equipment: []*domain.Equipment{
			{ID: uuid.New()}, // different equipment
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsEquipmentAccessError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Fail_TraineeNoAccess() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	userID := payload.UserID
	eqID := uuid.New()

	dep := &domain.Department{
		Users: []*domain.UserDepartment{
			{UserID: userID, Role: domain.RoleTrainee, User: &domain.User{Nice: 100}},
		},
		Equipment: []*domain.Equipment{
			{ID: eqID, AvailableToTrainee: false},
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsEquipmentAccessError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Department_Fail_ActivistNoAccess() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	userID := payload.UserID
	eqID := uuid.New()

	dep := &domain.Department{
		Users: []*domain.UserDepartment{
			{UserID: userID, Role: domain.RoleActivist, User: &domain.User{Nice: 50}}, // nice < 60
		},
		Equipment: []*domain.Equipment{
			{ID: eqID, AvailableToTrainee: false},
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		Department:  dep,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsEquipmentAccessError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Organization_Fail_UsersNotLoaded() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	eqID := uuid.New()

	org := &domain.Organization{
		Users: nil, // not loaded
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID:  eqID,
		Organization: org,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Organization_Fail_UserNotInOrganization() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	eqID := uuid.New()

	org := &domain.Organization{
		Users: []*domain.User{
			{ID: uuid.New()}, // different user
		},
	}

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID:  eqID,
		Organization: org,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsEquipmentAccessError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Fail_NoDepartmentOrOrganization() {
	payload := domain.TokenPayload{UserID: uuid.New()}
	eqID := uuid.New()

	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: eqID,
		// neither Department nor Organization
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *AccessServiceSuite) TestHaveAccessToEquipment_Fail_AutherError() {
	req := &accessservice.HaveEquipmentAccessRequest{
		EquipmentID: uuid.New(),
		Department:  &domain.Department{},
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(domain.TokenPayload{}, errors.New("auth error"))

	_, err := s.svc.HaveAccessToEquipment(ctx, req)

	s.Error(err)
}

type AuthZMock struct{ mock.Mock }

func (m *AuthZMock) Authorize(ctx context.Context, payload domain.TokenPayload) context.Context {
	return m.Called(ctx, payload).Get(0).(context.Context)
}
func (m *AuthZMock) TokenPayloadFromContext(ctx context.Context) (domain.TokenPayload, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.TokenPayload), args.Error(1)
}

func TestAccessServiceSuite(t *testing.T) {
	suite.Run(t, new(AccessServiceSuite))
}

//go:build unit

package studioservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/studioservice"
	studiosearch "media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type StudioServiceSuite struct {
	suite.Suite

	svc *studioservice.StudioService

	studioRepo       *StudioInvocationRepoMock
	departmentRepo   *DepartmentRepoMock
	organizationRepo *OrganizationRepoMock
	accessService    *AccessServiceMock
	searchRepo       *SearchStudioInvocationRepoMock
	auther           *AuthZMock
	txm              *TxManagerMock
}

func (s *StudioServiceSuite) SetupTest() {
	s.studioRepo = new(StudioInvocationRepoMock)
	s.departmentRepo = new(DepartmentRepoMock)
	s.organizationRepo = new(OrganizationRepoMock)
	s.accessService = new(AccessServiceMock)
	s.searchRepo = new(SearchStudioInvocationRepoMock)
	s.auther = new(AuthZMock)
	s.txm = new(TxManagerMock)

	// Create search service for availability
	searchSvc := studiosearch.NewSearchStudioService(s.searchRepo)

	svc := studioservice.NewStudioService(
		s.studioRepo,
		s.searchRepo,
		s.departmentRepo,
		s.organizationRepo,
		s.accessService,
		searchSvc, // for availability
		s.auther,
		s.txm,
	)
	s.svc = &svc
}

func (s *StudioServiceSuite) TestGet_Success_Admin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	studioID := uuid.New()
	expectedStudio := &domain.StudioInvocation{ID: studioID}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, studioID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedStudio, nil)

	result, err := (*s.svc).Get(ctx, studioID)

	s.NoError(err)
	s.Equal(expectedStudio, result)
}

func (s *StudioServiceSuite) TestGet_Fail_NotAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	studioID := uuid.New()

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := (*s.svc).Get(ctx, studioID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *StudioServiceSuite) TestGet_Fail_AutherError() {
	studioID := uuid.New()

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(domain.TokenPayload{}, errors.New("auth error"))

	_, err := (*s.svc).Get(ctx, studioID)

	s.Error(err)
}

func (s *StudioServiceSuite) TestDelete_Success_Owner() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	userID := payload.UserID
	studioID := uuid.New()
	studio := &domain.StudioInvocation{
		ID:     studioID,
		UserID: userID,
		Status: domain.StudioCreated,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, studioID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(studio, nil)
	s.studioRepo.On("Delete", ctx, studioID).Return(nil)

	err := (*s.svc).Delete(ctx, studioID)

	s.NoError(err)
}

func (s *StudioServiceSuite) TestDelete_Success_Admin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	studioID := uuid.New()
	studio := &domain.StudioInvocation{
		ID:     studioID,
		UserID: uuid.New(), // different user
		Status: domain.StudioCreated,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, studioID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(studio, nil)
	s.studioRepo.On("Delete", ctx, studioID).Return(nil)

	err := (*s.svc).Delete(ctx, studioID)

	s.NoError(err)
}

func (s *StudioServiceSuite) TestDelete_Fail_NonOwnerNonAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	studioID := uuid.New()
	studio := &domain.StudioInvocation{
		ID:     studioID,
		UserID: uuid.New(), // different user
		Status: domain.StudioCreated,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, studioID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(studio, nil)

	err := (*s.svc).Delete(ctx, studioID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *StudioServiceSuite) TestDelete_Fail_CompletedStatus() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	studioID := uuid.New()
	studio := &domain.StudioInvocation{
		ID:     studioID,
		UserID: uuid.New(),
		Status: domain.StudioCompleted,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, studioID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(studio, nil)

	err := (*s.svc).Delete(ctx, studioID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *StudioServiceSuite) TestDelete_Fail_CancelledStatus() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	studioID := uuid.New()
	studio := &domain.StudioInvocation{
		ID:     studioID,
		UserID: uuid.New(),
		Status: domain.StudioCancelled,
	}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.studioRepo.On("Get", ctx, studioID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(studio, nil)

	err := (*s.svc).Delete(ctx, studioID)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func (s *StudioServiceSuite) TestDelete_Fail_AutherError() {
	studioID := uuid.New()

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(domain.TokenPayload{}, errors.New("auth error"))

	err := (*s.svc).Delete(ctx, studioID)

	s.Error(err)
}

func (s *StudioServiceSuite) TestSearch_Success_Admin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{domain.AdminRole}}
	searchReq := &studiosearch.SearchStudioInvocationRequest{}
	expectedStudios := []*domain.StudioInvocation{{ID: uuid.New()}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.searchRepo.On("Search", ctx, mock.MatchedBy(func(req *studiosearch.SearchStudioInvocationRequest) bool {
		return req.Statuses != nil && len(req.Statuses) == 4
	}), mock.Anything).Return(expectedStudios, nil)

	result, err := (*s.svc).Search(ctx, searchReq)

	s.NoError(err)
	s.Equal(expectedStudios, result)
}

func (s *StudioServiceSuite) TestSearch_Fail_NotAdmin() {
	payload := domain.TokenPayload{UserID: uuid.New(), Roles: []domain.RoleAuth{}}
	searchReq := &studiosearch.SearchStudioInvocationRequest{}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	_, err := (*s.svc).Search(ctx, searchReq)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *StudioServiceSuite) TestSearch_Fail_AutherError() {
	searchReq := &studiosearch.SearchStudioInvocationRequest{}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(domain.TokenPayload{}, errors.New("auth error"))

	_, err := (*s.svc).Search(ctx, searchReq)

	s.Error(err)
}

func (s *StudioServiceSuite) TestSchedule_Success() {
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	req := &studioservice.StudioScheduleRequest{
		StartTime: startTime,
		EndTime:   endTime,
	}
	expectedStudios := []*domain.StudioInvocation{{ID: uuid.New()}}

	ctx := context.Background()
	s.searchRepo.On("Search", ctx, mock.MatchedBy(func(r *studiosearch.SearchStudioInvocationRequest) bool {
		return r.StartTime != nil && r.EndTime != nil && len(r.Statuses) == 4
	}), mock.Anything).Return(expectedStudios, nil)

	result, err := (*s.svc).Schedule(ctx, req)

	s.NoError(err)
	s.Equal(expectedStudios, result)
}

func (s *StudioServiceSuite) TestSchedule_Fail_ValidationError() {
	req := &studioservice.StudioScheduleRequest{
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now(), // EndTime before StartTime
	}

	ctx := context.Background()

	_, err := (*s.svc).Schedule(ctx, req)

	s.Error(err)
	s.True(errs.IsValidationError(err))
}

func TestStudioServiceSuite(t *testing.T) {
	suite.Run(t, new(StudioServiceSuite))
}
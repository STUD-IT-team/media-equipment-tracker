package studioservice

import (
	"context"
	"slices"
	"time"

	"media-equipment-tracker/internal/application/accessservice"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type StudioService interface {
	studiosearch.SearchStudioService
	CreateStudioService
	UpdateStudioService
	CuratorStudioService
	Get(ctx context.Context, id uuid.UUID) (*domain.StudioInvocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Schedule(ctx context.Context, req *StudioScheduleRequest) ([]*domain.StudioInvocation, error)
}

type studioService struct {
	CreateStudioService
	UpdateStudioService
	CuratorStudioService

	studiorepo domain.StudioInvocationRepository
	search     studiosearch.SearchStudioService
	auther     authzservice.AuthZ
	txm        txmanager.TxManager
}

var _ StudioService = (*studioService)(nil)

func NewStudioService(
	studiorepo domain.StudioInvocationRepository,
	searchrepo studiosearch.SearchStudioInvocationRepository,
	departmentRepo domain.DepartmentRepository,
	organizationRepo domain.OrganizationRepository,
	accessService accessservice.AccessService,
	auther authzservice.AuthZ,
	txm txmanager.TxManager,
) StudioService {
	return &studioService{
		studiorepo: studiorepo,
		search:     studiosearch.NewSearchStudioService(searchrepo),
		CreateStudioService: NewCreateStudioService(
			studiorepo,
			departmentRepo,
			organizationRepo,
			NewStudioInvocationAvailabilityService(searchrepo),
			accessService,
			txm,
			auther,
		),
		UpdateStudioService: NewUpdateStudioService(
			studiorepo,
			departmentRepo,
			organizationRepo,
			NewStudioInvocationAvailabilityService(searchrepo),
			accessService,
			txm,
			auther,
		),
		CuratorStudioService: NewCuratorStudioService(
			studiorepo,
			txm,
			auther,
		),
		auther: auther,
		txm:    txm,
	}
}

func (s *studioService) Get(ctx context.Context, id uuid.UUID) (*domain.StudioInvocation, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	return s.studiorepo.Get(ctx, id, domain.StudioInvocationWithDepartment(), domain.StudioInvocationWithOrganization(), domain.StudioInvocationWithUser(), domain.StudioInvocationWithAdmin())
}

func (s *studioService) Delete(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		inv, err := s.studiorepo.Get(ctx, id, domain.StudioInvocationWithDepartment(), domain.StudioInvocationWithOrganization(), domain.StudioInvocationWithUser(), domain.StudioInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.Status == domain.StudioCompleted || inv.Status == domain.StudioCancelled {
			return errs.NewValidationError("Status", "can't delete already completed or cancelled invocation")
		}

		if inv.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		return s.studiorepo.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *studioService) Search(ctx context.Context, search *studiosearch.SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	return s.search.Search(ctx, search, with...)
}

type StudioScheduleRequest struct {
	StartTime time.Time `validate:"required"`
	EndTime   time.Time `validate:"required,gtfield=StartTime"`
}

func (s *studioService) Schedule(ctx context.Context, req *StudioScheduleRequest) ([]*domain.StudioInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("StudioScheduleRequest", err.Error())
	}

	return s.search.Search(ctx, &studiosearch.SearchStudioInvocationRequest{
		StartTime: &req.StartTime,
		EndTime:   &req.EndTime,
		Statuses: []domain.StudioInvocationStatus{
			domain.StudioCreated,
			domain.StudioApproved,
			domain.StudioUnderReview,
			domain.StudioChangesRequired,
		},
	})
}

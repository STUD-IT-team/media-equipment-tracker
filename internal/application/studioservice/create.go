package studioservice

import (
	"context"
	"time"

	"media-equipment-tracker/internal/application/accessservice"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type CreateStudioInvocationRequest struct {
	EventName           string     `validate:"required,min=4,max=255"`
	ShootingDescription string     `validate:"required,min=4,max=2000"`
	StartTime           time.Time  `validate:"required"`
	EndTime             time.Time  `validate:"required,gtfield=StartTime"`
	OrganizationID      *uuid.UUID `validate:"omitempty,required_without=DepartmentID"`
	DepartmentID        *uuid.UUID `validate:"omitempty,required_without=OrganizationID"`

	NeedsChromakey   bool `validate:"required"`
	NeedsCyclorama   bool `validate:"required"`
	NeedsBlackFabric bool `validate:"required"`
}

type CreateStudioService interface {
	Create(ctx context.Context, req *CreateStudioInvocationRequest) (*domain.StudioInvocation, error)
}

type createStudioService struct {
	studioRepo       domain.StudioInvocationRepository
	departmentRepo   domain.DepartmentRepository
	organizationRepo domain.OrganizationRepository

	availabilityService AvailabilityStudioService
	accessService       accessservice.AccessService

	txm    txmanager.TxManager
	auther authzservice.AuthZ
}

func NewCreateStudioService(
	studioRepo domain.StudioInvocationRepository,
	departmentRepo domain.DepartmentRepository,
	organizationRepo domain.OrganizationRepository,
	availabilityService AvailabilityStudioService,
	accessService accessservice.AccessService,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) CreateStudioService {
	return &createStudioService{
		studioRepo:          studioRepo,
		departmentRepo:      departmentRepo,
		organizationRepo:    organizationRepo,
		availabilityService: availabilityService,
		accessService:       accessService,
		txm:                 txm,
		auther:              auther,
	}
}

var _ CreateStudioService = (*createStudioService)(nil)

func (s *createStudioService) Create(ctx context.Context, req *CreateStudioInvocationRequest) (*domain.StudioInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("CreateStudioInvocationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var org *domain.Organization
	var dep *domain.Department
	var inv *domain.StudioInvocation

	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		if req.OrganizationID != nil {
			org, err = s.organizationRepo.Get(ctx, *req.OrganizationID, domain.WithUsers())
			if err != nil {
				return err
			}
		} else if req.DepartmentID != nil {
			dep, err = s.departmentRepo.Get(ctx, *req.DepartmentID, domain.DepartmentWithUsers())
			if err != nil {
				return err
			}
		}

		resp, err := s.availabilityService.Availability(ctx, AvailabilityStudioInvocationRequest{
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
		})

		if err != nil {
			return err
		}

		if !resp.Available {
			return errs.NewValidationError("Time", "Time is not available")
		}

		access, err := s.accessService.HaveAccessToStudio(ctx, &accessservice.HaveStudioAccessRequest{
			Organization: org,
			Department:   dep,
		})

		if err != nil {
			return err
		}

		if !access {
			return errs.NewStudioAccessError("user do not have access to studio")
		}

		inv = &domain.StudioInvocation{
			ID:                  uuid.New(),
			EventName:           req.EventName,
			ShootingDescription: req.ShootingDescription,
			StartTime:           req.StartTime,
			EndTime:             req.EndTime,
			NeedsChromakey:      req.NeedsChromakey,
			NeedsCyclorama:      req.NeedsCyclorama,
			NeedsBlackFabric:    req.NeedsBlackFabric,
			Status:              domain.StudioCreated,
			OrganizationID:      req.OrganizationID,
			DepartmentID:        req.DepartmentID,
			UserID:              payload.UserID,
		}

		err = s.studioRepo.Create(ctx, inv)
		if err != nil {
			return err
		}

		err = s.studioRepo.Reload(ctx, inv, domain.StudioInvocationWithDepartment(), domain.StudioInvocationWithOrganization(), domain.StudioInvocationWithUser())
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return inv, nil
}

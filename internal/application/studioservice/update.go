package studioservice

import (
	"context"
	"slices"
	"time"

	"media-equipment-tracker/internal/application/accessservice"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type UpdateStudioInvocationRequest struct {
	ID                  uuid.UUID  `validate:"required"`
	EventName           *string    `validate:"omitempty,min=4,max=255"`
	ShootingDescription *string    `validate:"omitempty,min=4,max=2000"`
	StartTime           *time.Time `validate:"omitempty"`
	EndTime             *time.Time `validate:"omitempty,gtfield=StartTime"`
	NeedsChromakey      *bool      `validate:"omitempty"`
	NeedsCyclorama      *bool      `validate:"omitempty"`
	NeedsBlackFabric    *bool      `validate:"omitempty"`
}

type UpdateStudioService interface {
	Update(ctx context.Context, req *UpdateStudioInvocationRequest) (*domain.StudioInvocation, error)
	UpdateStatus(ctx context.Context, req *UpdateStudioInvocationStatusRequest) (*domain.StudioInvocation, error)
}

type updateStudioService struct {
	studioRepo       domain.StudioInvocationRepository
	departmentRepo   domain.DepartmentRepository
	organizationRepo domain.OrganizationRepository

	availabilityService AvailabilityStudioService
	accessService       accessservice.AccessService

	txm    txmanager.TxManager
	auther authzservice.AuthZ
}

var _ UpdateStudioService = (*updateStudioService)(nil)

func NewUpdateStudioService(
	studioRepo domain.StudioInvocationRepository,
	departmentRepo domain.DepartmentRepository,
	organizationRepo domain.OrganizationRepository,
	availabilityService AvailabilityStudioService,
	accessService accessservice.AccessService,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) UpdateStudioService {
	return &updateStudioService{
		studioRepo:          studioRepo,
		departmentRepo:      departmentRepo,
		organizationRepo:    organizationRepo,
		availabilityService: availabilityService,
		accessService:       accessService,
		txm:                 txm,
		auther:              auther,
	}
}

func (s *updateStudioService) Update(ctx context.Context, req *UpdateStudioInvocationRequest) (*domain.StudioInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateStudioInvocationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var org *domain.Organization
	var dep *domain.Department
	var inv *domain.StudioInvocation

	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		inv, err = s.studioRepo.Get(ctx, req.ID, domain.StudioInvocationWithDepartment(), domain.StudioInvocationWithOrganization(), domain.StudioInvocationWithUser())
		if err != nil {
			return err
		}

		if inv.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status == domain.StudioCompleted || inv.Status == domain.StudioCancelled {
			return errs.NewValidationError("Status", "can't update already completed or cancelled invocation")
		}

		if inv.OrganizationID != nil {
			org = inv.Organization
			err = s.organizationRepo.Reload(ctx, org, domain.WithUsers())
			if err != nil {
				return err
			}
		} else if inv.DepartmentID != nil {
			dep = inv.Department
			err = s.departmentRepo.Reload(ctx, dep, domain.DepartmentWithUsers())
			if err != nil {
				return err
			}
		}

		if req.EventName != nil {
			inv.EventName = *req.EventName
		}
		if req.ShootingDescription != nil {
			inv.ShootingDescription = *req.ShootingDescription
		}
		if req.StartTime != nil {
			inv.StartTime = *req.StartTime
		}
		if req.EndTime != nil {
			inv.EndTime = *req.EndTime
		}
		if req.NeedsChromakey != nil {
			inv.NeedsChromakey = *req.NeedsChromakey
		}
		if req.NeedsCyclorama != nil {
			inv.NeedsCyclorama = *req.NeedsCyclorama
		}
		if req.NeedsBlackFabric != nil {
			inv.NeedsBlackFabric = *req.NeedsBlackFabric
		}

		if inv.Status == domain.StudioChangesRequired {
			inv.Status = domain.StudioUnderReview
		}

		resp, err := s.availabilityService.Availability(ctx, AvailabilityStudioInvocationRequest{
			StartTime: inv.StartTime,
			EndTime:   inv.EndTime,
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

		err = s.studioRepo.Update(ctx, inv)
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

type UpdateStudioInvocationStatusRequest struct {
	ID             uuid.UUID                     `validate:"required"`
	Status         domain.StudioInvocationStatus `validate:"required"`
	CuratorComment *string                       `validate:"omitempty,min=4,max=2000"`
}

func (s *updateStudioService) UpdateStatus(ctx context.Context, req *UpdateStudioInvocationStatusRequest) (*domain.StudioInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateStudioInvocationStatusRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	var inv *domain.StudioInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		inv, err = s.studioRepo.Get(ctx, req.ID, domain.StudioInvocationWithAdmin())
		if err != nil {
			return err
		}

		inv.Status = req.Status
		if req.CuratorComment != nil {
			inv.CuratorComment = req.CuratorComment
		}

		err = s.studioRepo.Update(ctx, inv)
		if err != nil {
			return err
		}

		return nil
	})

	return inv, nil
}

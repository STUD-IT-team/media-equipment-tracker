package invocationservice

import (
	"context"
	"fmt"
	"slices"
	"time"

	"media-equipment-tracker/internal/application/accessservice"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type UpdateInvocationRequest struct {
	ID           uuid.UUID   `validate:"required"`
	EventName    *string     `validate:"omitempty,min=4,max=255"`
	StartTime    *time.Time  `validate:"omitempty"`
	EndTime      *time.Time  `validate:"omitempty,gtfield=StartTime"`
	EquipmentIDs []uuid.UUID `validate:"omitempty,min=1,dive"`
}

type UpdateInvocationStatusRequest struct {
	ID             uuid.UUID                        `validate:"required"`
	Status         domain.EquipmentInvocationStatus `validate:"required"`
	CuratorComment *string                          `validate:"omitempty,min=4,max=2000"`
}

type UpdateInvocationService interface {
	Update(ctx context.Context, req *UpdateInvocationRequest) (*domain.EquipmentInvocation, error)
	UpdateStatus(ctx context.Context, req *UpdateInvocationStatusRequest) (*domain.EquipmentInvocation, error)
}

type updateInvocationService struct {
	invocationRepo      domain.EquipmentInvocationRepository
	departmentRepo      domain.DepartmentRepository
	organizationRepo    domain.OrganizationRepository
	availabilityService equipmentservice.AvailabilityEquipmentService
	accessService       accessservice.AccessService
	txm                 txmanager.TxManager
	auther              authzservice.AuthZ
}

func NewUpdateInvocationService(
	invocationRepo domain.EquipmentInvocationRepository,
	departmentRepo domain.DepartmentRepository,
	organizationRepo domain.OrganizationRepository,
	availabilityService equipmentservice.AvailabilityEquipmentService,
	accessService accessservice.AccessService,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) UpdateInvocationService {
	return &updateInvocationService{
		invocationRepo:      invocationRepo,
		departmentRepo:      departmentRepo,
		organizationRepo:    organizationRepo,
		availabilityService: availabilityService,
		accessService:       accessService,
		txm:                 txm,
		auther:              auther,
	}
}

var _ UpdateInvocationService = (*updateInvocationService)(nil)

func (s *updateInvocationService) Update(ctx context.Context, req *UpdateInvocationRequest) (*domain.EquipmentInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateInvocationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, req.ID, domain.EquipmentInvocationWithEquipment(), domain.EquipmentInvocationWithDepartment(), domain.EquipmentInvocationWithOrganization())
		if err != nil {
			return err
		}

		if inv.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status == domain.InvocationEquipmentIssued ||
			inv.Status == domain.InvocationEquipmentReturned ||
			inv.Status == domain.InvocationCompleted ||
			inv.Status == domain.InvocationCancelled {
			return errs.NewValidationError("Status", "can't update already issued invocation")
		}

		if req.EventName != nil {
			inv.EventName = *req.EventName
		}
		if req.StartTime != nil {
			inv.StartTime = *req.StartTime
		}
		if req.EndTime != nil {
			inv.EndTime = *req.EndTime
		}

		if inv.EndTime.Before(inv.StartTime) {
			return errs.NewValidationError("EndTime", "EndTime must be after StartTime")
		}

		if req.EquipmentIDs != nil {

			inv.Equipment = make([]*domain.EquipmentInInvocation, 0, len(req.EquipmentIDs))
			for _, id := range req.EquipmentIDs {
				inv.Equipment = append(inv.Equipment, &domain.EquipmentInInvocation{
					InvocationID: inv.ID,
					EquipmentID:  id,
					Status:       domain.EquipmentNotIssued,
				})
			}

			if inv.DepartmentID != nil {
				err = s.departmentRepo.Reload(ctx, inv.Department, domain.DepartmentWithUsers(), domain.DepartmentWithEquipment())
				if err != nil {
					return err
				}
			} else if inv.OrganizationID != nil {
				err = s.organizationRepo.Reload(ctx, inv.Organization, domain.WithUsers())
				if err != nil {
					return err
				}
			}

			for _, eq := range inv.Equipment {
				// Проверка что обрудование доступно физически
				resp, err := s.availabilityService.Availability(ctx, equipmentservice.EquipmentAvailabilityRequest{
					ID:            eq.EquipmentID,
					StartTime:     inv.StartTime,
					EndTime:       inv.EndTime,
					ForInvocation: &inv.ID,
				})

				if err != nil {
					return err
				}
				if !resp.Available {
					return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s is not available", eq.EquipmentID))
				}

				// и душевно
				access, err := s.accessService.HaveAccessToEquipment(ctx, &accessservice.HaveEquipmentAccessRequest{
					EquipmentID:  eq.EquipmentID,
					Department:   inv.Department,
					Organization: inv.Organization,
				})

				if err != nil {
					return err
				}

				if !access {
					return errs.NewEquipmentAccessError(fmt.Sprintf("Equipment %s is no accessible for current user (undefined reason)", eq.EquipmentID))
				}
			}
		}

		// Обновляем статус, чтобы заявку снова проверили
		if inv.Status == domain.InvocationApproved || inv.Status == domain.InvocationChangesRequired {
			inv.Status = domain.InvocationUnderReview
		}

		err = s.invocationRepo.Update(ctx, inv)
		if err != nil {
			return err
		}

		err = s.invocationRepo.Reload(ctx, inv, domain.EquipmentInvocationWithDepartment(), domain.EquipmentInvocationWithOrganization(), domain.EquipmentInvocationWithUser(), domain.EquipmentInvocationWithAdmin(), domain.EquipmentInvocationWithEquipment())
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

func (s *updateInvocationService) UpdateStatus(ctx context.Context, req *UpdateInvocationStatusRequest) (*domain.EquipmentInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateInvocationStatusRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, req.ID, domain.EquipmentInvocationWithAdmin())
		if err != nil {
			return err
		}

		inv.Status = req.Status
		if req.CuratorComment != nil {
			inv.CuratorComment = *req.CuratorComment
		}

		err = s.invocationRepo.Update(ctx, inv)
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

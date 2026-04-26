package invocationservice

import (
	"context"
	"fmt"
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

type CreateInvocationService interface {
	Create(ctx context.Context, req *CreateInvocationRequest) (*domain.EquipmentInvocation, error)
}

type createInvocationService struct {
	auther              authzservice.AuthZ
	invocationRepo      domain.EquipmentInvocationRepository
	departmentRepo      domain.DepartmentRepository
	organizationRepo    domain.OrganizationRepository
	accessService       accessservice.AccessService
	availabilityService equipmentservice.AvailabilityEquipmentService
	txm                 txmanager.TxManager
}

var _ CreateInvocationService = (*createInvocationService)(nil)

func NewCreateInvocationService(
	invocationRepo domain.EquipmentInvocationRepository,
	departmentRepo domain.DepartmentRepository,
	organizationRepo domain.OrganizationRepository,
	accessService accessservice.AccessService,
	availabilityService equipmentservice.AvailabilityEquipmentService,
	auther authzservice.AuthZ,
	txm txmanager.TxManager,
) CreateInvocationService {
	return &createInvocationService{
		invocationRepo:      invocationRepo,
		departmentRepo:      departmentRepo,
		organizationRepo:    organizationRepo,
		accessService:       accessService,
		availabilityService: availabilityService,
		auther:              auther,
		txm:                 txm,
	}
}

type CreateInvocationRequest struct {
	EventName      string      `validate:"required,min=4,max=255"`
	StartTime      time.Time   `validate:"required"`
	EndTime        time.Time   `validate:"required,gtfield=StartTime"`
	OrganizationID *uuid.UUID  `validate:"omitempty,required_without=DepartmentID"`
	DepartmentID   *uuid.UUID  `validate:"omitempty,required_without=OrganizationID"`
	EquipmentIDs   []uuid.UUID `validate:"required,min=1,dive"`
}

func (s *createInvocationService) Create(ctx context.Context, req *CreateInvocationRequest) (*domain.EquipmentInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("CreateInvocationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	inv := &domain.EquipmentInvocation{
		ID:             uuid.New(),
		EventName:      req.EventName,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		OrganizationID: req.OrganizationID,
		DepartmentID:   req.DepartmentID,
		UserID:         payload.UserID,
		CuratorComment: "",
		Status:         domain.InvocationCreated,
	}

	eqs := make([]*domain.EquipmentInInvocation, 0, len(req.EquipmentIDs))
	for _, id := range req.EquipmentIDs {
		eqs = append(eqs, &domain.EquipmentInInvocation{
			EquipmentID:  id,
			InvocationID: inv.ID,
			Status:       domain.EquipmentNotIssued,
		})
	}

	inv.Equipment = eqs

	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var org *domain.Organization
		var dep *domain.Department
		var err error

		if req.OrganizationID != nil {
			org, err = s.organizationRepo.Get(ctx, *req.OrganizationID, domain.WithUsers())
			if err != nil {
				return err
			}
		} else {
			dep, err = s.departmentRepo.Get(ctx, *req.DepartmentID, domain.DepartmentWithUsers(), domain.DepartmentWithEquipment())
			if err != nil {
				return err
			}
		}

		for _, eq := range eqs {
			// Проверка что обрудование доступно физически
			resp, err := s.availabilityService.Availability(ctx, equipmentservice.EquipmentAvailabilityRequest{
				ID:        eq.EquipmentID,
				StartTime: inv.StartTime,
				EndTime:   inv.EndTime,
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
				Department:   dep,
				Organization: org,
			})

			if err != nil {
				return err
			}

			if !access {
				return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s is not available", eq.EquipmentID))
			}
		}

		err = s.invocationRepo.Create(ctx, inv)
		if err != nil {
			return err
		}

		err = s.invocationRepo.Reload(ctx, inv, domain.EquipmentInvocationWithDepartment(), domain.EquipmentInvocationWithOrganization(), domain.EquipmentInvocationWithUser(), domain.EquipmentInvocationWithEquipment())
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

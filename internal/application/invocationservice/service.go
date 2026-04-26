package invocationservice

import (
	"context"
	"fmt"
	"slices"

	"media-equipment-tracker/internal/application/accessservice"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/application/invocationservice/invocationsearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type InvocationService interface {
	invocationsearch.SearchInvocationService
	CreateInvocationService
	UpdateInvocationService
	CuratorInvocationService
	Get(ctx context.Context, id uuid.UUID) (*domain.EquipmentInvocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type invocationService struct {
	invocationsearch.SearchInvocationService
	CreateInvocationService
	UpdateInvocationService
	CuratorInvocationService
	invocationRepository domain.EquipmentInvocationRepository
	txm                  txmanager.TxManager
	auther               authzservice.AuthZ
}

func NewInvocationService(
	invocationRepository domain.EquipmentInvocationRepository,
	searchInvocationRepository invocationsearch.SearchInvocationRepository,
	departmentRepository domain.DepartmentRepository,
	organizationRepository domain.OrganizationRepository,
	availabilityService equipmentservice.AvailabilityEquipmentService,
	updateEquipmentService equipmentservice.UpdateEquipmentService,
	accessService accessservice.AccessService,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) InvocationService {
	return &invocationService{
		invocationRepository:     invocationRepository,
		SearchInvocationService:  invocationsearch.NewSearchInvocationService(searchInvocationRepository),
		CreateInvocationService:  NewCreateInvocationService(invocationRepository, departmentRepository, organizationRepository, accessService, availabilityService, auther, txm),
		UpdateInvocationService:  NewUpdateInvocationService(invocationRepository, departmentRepository, organizationRepository, availabilityService, accessService, txm, auther),
		CuratorInvocationService: NewCuratorInvocationService(invocationRepository, updateEquipmentService, txm, auther),
		txm:                      txm,
		auther:                   auther,
	}
}

var _ InvocationService = (*invocationService)(nil)

func (s *invocationService) Get(ctx context.Context, id uuid.UUID) (*domain.EquipmentInvocation, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	inv, err := s.invocationRepository.Get(ctx, id, domain.EquipmentInvocationWithUser(), domain.EquipmentInvocationWithDepartment(), domain.EquipmentInvocationWithOrganization(), domain.EquipmentInvocationWithEquipment())
	if err != nil {
		return nil, err
	}

	if inv.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	return inv, nil
}

func (s *invocationService) Delete(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}

	inv, err := s.invocationRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	if inv.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
		return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	if !slices.Contains([]domain.EquipmentInvocationStatus{
		domain.InvocationApproved,
		domain.InvocationChangesRequired,
		domain.InvocationCreated,
		domain.InvocationUnderReview,
	}, inv.Status) {
		return errs.NewValidationError("Status", fmt.Sprintf("can't delete invocation in status %s", inv.Status))
	}

	err = s.invocationRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

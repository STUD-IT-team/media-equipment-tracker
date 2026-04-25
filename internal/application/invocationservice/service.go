package invocationservice

import (
	"context"
	"slices"

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
	Get(ctx context.Context, id uuid.UUID) (*domain.EquipmentInvocation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type invocationService struct {
	invocationsearch.SearchInvocationService
	CreateInvocationService
	UpdateInvocationService
	invocationRepository domain.EquipmentInvocationRepository
	txm                  txmanager.TxManager
	auther               authzservice.AuthZ
}

func NewInvocationService(
	invocationRepository domain.EquipmentInvocationRepository,
	searchInvocationRepository invocationsearch.SearchInvocationRepository,
	availabilityService equipmentservice.AvailabilityEquipmentService,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) InvocationService {
	return &invocationService{
		invocationRepository:    invocationRepository,
		SearchInvocationService: invocationsearch.NewSearchInvocationService(searchInvocationRepository),
		CreateInvocationService: NewCreateInvocationService(invocationRepository, availabilityService, auther, txm),
		UpdateInvocationService: NewUpdateInvocationService(invocationRepository, availabilityService, txm, auther),
		txm:                     txm,
		auther:                  auther,
	}
}

var _ InvocationService = (*invocationService)(nil)

func (s *invocationService) Get(ctx context.Context, id uuid.UUID) (*domain.EquipmentInvocation, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	inv, err := s.invocationRepository.Get(ctx, id)
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

	return nil
}

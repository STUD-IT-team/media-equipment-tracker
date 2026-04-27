package messageequipmentservice

import (
	"context"
	"slices"

	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"github.com/google/uuid"
)

type SearchMessageEquipmentInvocationService interface {
	ListByInvocation(ctx context.Context, invocationID uuid.UUID, offset, limit int, with ...domain.MessageEquipmentInvocationOption) ([]*domain.MessageEquipmentInvocation, error)
}

type searchMessageEquipmentInvocationService struct {
	messageRepository domain.MessageEquipmentInvocationRepository
	invocationService invocationservice.InvocationService
	txm               txmanager.TxManager
	auther            authzservice.AuthZ
}

func NewSearchMessageEquipmentInvocationService(messageRepository domain.MessageEquipmentInvocationRepository, invocationService invocationservice.InvocationService, txm txmanager.TxManager, auther authzservice.AuthZ) SearchMessageEquipmentInvocationService {
	return &searchMessageEquipmentInvocationService{
		messageRepository: messageRepository,
		invocationService: invocationService,
		txm:               txm,
		auther:            auther,
	}
}

func (s *searchMessageEquipmentInvocationService) ListByInvocation(ctx context.Context, invocationID uuid.UUID, offset, limit int, with ...domain.MessageEquipmentInvocationOption) ([]*domain.MessageEquipmentInvocation, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the invocation
	invocation, err := s.invocationService.Get(ctx, invocationID)
	if err != nil {
		return nil, err
	}

	if invocation.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	var messages []*domain.MessageEquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		messages, err = s.messageRepository.GetInvocation(ctx, invocationID, with...)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Apply pagination
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = len(messages)
	}
	if offset >= len(messages) {
		return []*domain.MessageEquipmentInvocation{}, nil
	}
	end := offset + limit
	if end > len(messages) {
		end = len(messages)
	}
	return messages[offset:end], nil
}
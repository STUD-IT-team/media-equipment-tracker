package messagestudioservice

import (
	"context"
	"slices"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"github.com/google/uuid"
)

type SearchMessageStudioInvocationService interface {
	ListByInvocation(ctx context.Context, invocationID uuid.UUID, offset, limit int, with ...domain.MessageStudioInvocationOption) ([]*domain.MessageStudioInvocation, error)
}

type searchMessageStudioInvocationService struct {
	messageRepository domain.MessageStudioInvocationRepository
	studioInvocationRepository domain.StudioInvocationRepository
	txm               txmanager.TxManager
	auther            authzservice.AuthZ
}

func NewSearchMessageStudioInvocationService(messageRepository domain.MessageStudioInvocationRepository, studioInvocationRepository domain.StudioInvocationRepository, txm txmanager.TxManager, auther authzservice.AuthZ) SearchMessageStudioInvocationService {
	return &searchMessageStudioInvocationService{
		messageRepository:          messageRepository,
		studioInvocationRepository: studioInvocationRepository,
		txm:                        txm,
		auther:                     auther,
	}
}

func (s *searchMessageStudioInvocationService) ListByInvocation(ctx context.Context, invocationID uuid.UUID, offset, limit int, with ...domain.MessageStudioInvocationOption) ([]*domain.MessageStudioInvocation, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the invocation
	studioInvocation, err := s.studioInvocationRepository.Get(ctx, invocationID)
	if err != nil {
		return nil, err
	}

	if studioInvocation.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	var messages []*domain.MessageStudioInvocation
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
		return []*domain.MessageStudioInvocation{}, nil
	}
	end := offset + limit
	if end > len(messages) {
		end = len(messages)
	}
	return messages[offset:end], nil
}
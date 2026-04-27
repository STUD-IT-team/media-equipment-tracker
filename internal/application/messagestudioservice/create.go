package messagestudioservice

import (
	"context"
	"slices"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"github.com/google/uuid"
)

type CreateMessageStudioInvocationRequest struct {
	Content     string `validate:"required,max=2000"`
	RecipientID uuid.UUID `validate:"required"`
}

type CreateMessageStudioInvocationService interface {
	CreateMessageStudioInvocation(ctx context.Context, invocationID uuid.UUID, req *CreateMessageStudioInvocationRequest) (*domain.MessageStudioInvocation, error)
}

type createMessageStudioInvocationService struct {
	messageRepository domain.MessageStudioInvocationRepository
	studioInvocationRepository domain.StudioInvocationRepository
	txManager         txmanager.TxManager
	auther            authzservice.AuthZ
}

func NewCreateMessageStudioInvocationService(
	messageRepository domain.MessageStudioInvocationRepository,
	studioInvocationRepository domain.StudioInvocationRepository,
	txManager txmanager.TxManager,
	auther authzservice.AuthZ,
) CreateMessageStudioInvocationService {
	return &createMessageStudioInvocationService{
		messageRepository:          messageRepository,
		studioInvocationRepository: studioInvocationRepository,
		txManager:                  txManager,
		auther:                     auther,
	}
}

func (s *createMessageStudioInvocationService) CreateMessageStudioInvocation(ctx context.Context, invocationID uuid.UUID, req *CreateMessageStudioInvocationRequest) (*domain.MessageStudioInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("CreateMessageStudioInvocationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the invocation (only the creator or admin can send messages)
	studioInvocation, err := s.studioInvocationRepository.Get(ctx, invocationID)
	if err != nil {
		return nil, err
	}

	if studioInvocation.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	message := &domain.MessageStudioInvocation{
		ID:           uuid.New(),
		Content:      req.Content,
		InvocationID: invocationID,
		SenderID:     payload.UserID,
		RecipientID:  req.RecipientID,
	}

	err = s.txManager.WithinTx(ctx, func(ctx context.Context) error {
		return s.messageRepository.Create(ctx, message)
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}
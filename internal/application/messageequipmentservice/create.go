package messageequipmentservice

import (
	"context"
	"slices"

	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"github.com/google/uuid"
)

type CreateMessageEquipmentInvocationRequest struct {
	Content     string `validate:"required,max=2000"`
	RecipientID uuid.UUID `validate:"required"`
}

type CreateMessageEquipmentInvocationService interface {
	CreateMessageEquipmentInvocation(ctx context.Context, invocationID uuid.UUID, req *CreateMessageEquipmentInvocationRequest) (*domain.MessageEquipmentInvocation, error)
}

type createMessageEquipmentInvocationService struct {
	messageRepository domain.MessageEquipmentInvocationRepository
	invocationService invocationservice.InvocationService
	txManager         txmanager.TxManager
	auther            authzservice.AuthZ
}

func NewCreateMessageEquipmentInvocationService(
	messageRepository domain.MessageEquipmentInvocationRepository,
	invocationService invocationservice.InvocationService,
	txManager txmanager.TxManager,
	auther authzservice.AuthZ,
) CreateMessageEquipmentInvocationService {
	return &createMessageEquipmentInvocationService{
		messageRepository: messageRepository,
		invocationService: invocationService,
		txManager:         txManager,
		auther:            auther,
	}
}

func (s *createMessageEquipmentInvocationService) CreateMessageEquipmentInvocation(ctx context.Context, invocationID uuid.UUID, req *CreateMessageEquipmentInvocationRequest) (*domain.MessageEquipmentInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("CreateMessageEquipmentInvocationRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the invocation (only the creator or admin can send messages)
	invocation, err := s.invocationService.Get(ctx, invocationID)
	if err != nil {
		return nil, err
	}

	if invocation.UserID != payload.UserID && !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	message := &domain.MessageEquipmentInvocation{
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
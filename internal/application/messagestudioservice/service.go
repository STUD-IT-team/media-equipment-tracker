package messagestudioservice

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/txmanager"

	authzservice "media-equipment-tracker/internal/application/authz_service"
)

type MessageStudioInvocationService interface {
	SearchMessageStudioInvocationService
	CreateMessageStudioInvocationService
}

type messageStudioInvocationService struct {
	SearchMessageStudioInvocationService
	CreateMessageStudioInvocationService
	txm               txmanager.TxManager
	auther            authzservice.AuthZ
}

func NewMessageStudioInvocationService(
	messageRepository domain.MessageStudioInvocationRepository,
	studioInvocationRepository domain.StudioInvocationRepository,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) MessageStudioInvocationService {
	return &messageStudioInvocationService{
		SearchMessageStudioInvocationService: NewSearchMessageStudioInvocationService(messageRepository, studioInvocationRepository, txm, auther),
		CreateMessageStudioInvocationService: NewCreateMessageStudioInvocationService(messageRepository, studioInvocationRepository, txm, auther),
		txm:                                    txm,
		auther:                                 auther,
	}
}
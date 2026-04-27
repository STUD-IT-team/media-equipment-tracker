package messageequipmentservice

import (
	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/txmanager"

	authzservice "media-equipment-tracker/internal/application/authz_service"
)

type MessageEquipmentInvocationService interface {
	SearchMessageEquipmentInvocationService
	CreateMessageEquipmentInvocationService
}

type messageEquipmentInvocationService struct {
	SearchMessageEquipmentInvocationService
	CreateMessageEquipmentInvocationService
	txm               txmanager.TxManager
	auther            authzservice.AuthZ
}

func NewMessageEquipmentInvocationService(
	messageRepository domain.MessageEquipmentInvocationRepository,
	invocationService invocationservice.InvocationService,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) MessageEquipmentInvocationService {
	return &messageEquipmentInvocationService{
		SearchMessageEquipmentInvocationService: NewSearchMessageEquipmentInvocationService(messageRepository, invocationService, txm, auther),
		CreateMessageEquipmentInvocationService: NewCreateMessageEquipmentInvocationService(messageRepository, invocationService, txm, auther),
		txm:                                     txm,
		auther:                                  auther,
	}
}
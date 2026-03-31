package domain

import (
	"context"

	"github.com/google/uuid"
)

type MessageEquipmentInvocationRepository interface {
	Create(ctx context.Context, message *MessageEquipmentInvocation) error
	GetByID(ctx context.Context, id uuid.UUID) (*MessageEquipmentInvocation, error)
	GetAll(ctx context.Context) ([]*MessageEquipmentInvocation, error)
	Update(ctx context.Context, message *MessageEquipmentInvocation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

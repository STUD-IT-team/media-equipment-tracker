package domain

import (
	"context"

	"github.com/google/uuid"
)

type MessageStudioInvocationRepository interface {
	Create(ctx context.Context, message *MessageStudioInvocation) error
	GetByID(ctx context.Context, id uuid.UUID) (*MessageStudioInvocation, error)
	GetAll(ctx context.Context) ([]*MessageStudioInvocation, error)
	Update(ctx context.Context, message *MessageStudioInvocation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

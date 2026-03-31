package domain

import (
	"context"

	"github.com/google/uuid"
)

type StudioInvocationRepository interface {
	Create(ctx context.Context, studioInvocation *StudioInvocation) error
	GetByID(ctx context.Context, id uuid.UUID) (*StudioInvocation, error)
	GetAll(ctx context.Context) ([]*StudioInvocation, error)
	Update(ctx context.Context, studioInvocation *StudioInvocation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

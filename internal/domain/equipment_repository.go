package domain

import (
	"context"

	"github.com/google/uuid"
)

type EquipmentRepository interface {
	Create(ctx context.Context, equipment *Equipment) error
	GetByID(ctx context.Context, id uuid.UUID) (*Equipment, error)
	GetAll(ctx context.Context) ([]*Equipment, error)
	Update(ctx context.Context, equipment *Equipment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

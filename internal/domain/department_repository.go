package domain

import (
	"context"

	"github.com/google/uuid"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) error
	GetByID(ctx context.Context, id uuid.UUID) (*Department, error)
	GetAll(ctx context.Context) ([]*Department, error)
	Update(ctx context.Context, department *Department) error
	Delete(ctx context.Context, id uuid.UUID) error
}

package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization *Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	GetAll(ctx context.Context) ([]*Organization, error)
	Update(ctx context.Context, organization *Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
}

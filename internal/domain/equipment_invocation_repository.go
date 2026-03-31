package domain

import (
	"context"

	"github.com/google/uuid"
)

type EquipmentInvocationRepository interface {
	Create(ctx context.Context, equipmentInvocation *EquipmentInvocation) error
	GetByID(ctx context.Context, id uuid.UUID) (*EquipmentInvocation, error)
	GetAll(ctx context.Context) ([]*EquipmentInvocation, error)
	Update(ctx context.Context, equipmentInvocation *EquipmentInvocation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

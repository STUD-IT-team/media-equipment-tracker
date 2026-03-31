package domain

import (
	"context"

	"github.com/google/uuid"
)

type EquipmentInInvocationRepository interface {
	Create(ctx context.Context, equipmentInInvocation *EquipmentInInvocation) error
	GetByInvocationIDAndEquipmentID(ctx context.Context, invocationID, equipmentID uuid.UUID) (*EquipmentInInvocation, error)
	GetAllByInvocationID(ctx context.Context, invocationID uuid.UUID) ([]*EquipmentInInvocation, error)
	GetAll(ctx context.Context) ([]*EquipmentInInvocation, error)
	Update(ctx context.Context, equipmentInInvocation *EquipmentInInvocation) error
	Delete(ctx context.Context, invocationID, equipmentID uuid.UUID) error
}

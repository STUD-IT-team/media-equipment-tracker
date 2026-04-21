package domain

import (
	"context"

	"github.com/google/uuid"
)

type EquipmentInInvocation struct {
	InvocationID uuid.UUID                   `gorm:"type:uuid;primaryKey"`
	EquipmentID  uuid.UUID                   `gorm:"type:uuid;primaryKey"`
	Status       EquipmentInInvocationStatus `gorm:"type:equipment_in_invocation_status"`

	Invocation *EquipmentInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
	Equipment  *Equipment           `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE"`
}

func (EquipmentInInvocation) TableName() string {
	return "equipment_in_invocation"
}

type EquipmentInInvocationStatus string

const (
	EquipmentNotIssued EquipmentInInvocationStatus = "not_issued"
	EquipmentIssued    EquipmentInInvocationStatus = "issued"
	EquipmentReturned  EquipmentInInvocationStatus = "returned"
)

type EquipmentInInvocationOptions struct {
	relations      []string
	withInvocation bool
	withEquipment  bool
}

func (o *EquipmentInInvocationOptions) Relations() []string {
	return o.relations
}

type EquipmentInInvocationOption func(*EquipmentInInvocationOptions)

func WithInvocation() EquipmentInInvocationOption {
	return func(options *EquipmentInInvocationOptions) {
		options.withInvocation = true
		options.relations = append(options.relations, "Invocation")
	}
}

func WithEquipment() EquipmentInInvocationOption {
	return func(options *EquipmentInInvocationOptions) {
		options.withEquipment = true
		options.relations = append(options.relations, "Equipment")
	}
}

type EquipmentInInvocationRepository interface {
	Get(ctx context.Context, invocationID uuid.UUID, equipmentID uuid.UUID, with ...EquipmentInInvocationOption) (*EquipmentInInvocation, error)
	GetByInvocation(ctx context.Context, invocationID uuid.UUID, with ...EquipmentInInvocationOption) ([]*EquipmentInInvocation, error)
	GetByEquipment(ctx context.Context, equipmentID uuid.UUID, with ...EquipmentInInvocationOption) ([]*EquipmentInInvocation, error)
	Create(ctx context.Context, equipmentInInvocation *EquipmentInInvocation) error
	Update(ctx context.Context, equipmentInInvocation *EquipmentInInvocation) error
	Delete(ctx context.Context, invocationID uuid.UUID, equipmentID uuid.UUID) error
}

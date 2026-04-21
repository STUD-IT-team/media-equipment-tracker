package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MessageEquipmentInvocation struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Content  string    `gorm:"type:text;not null"`
	SentTime time.Time `gorm:"column:sent_time;type:timestamptz;default:CURRENT_TIMESTAMP"`

	InvocationID uuid.UUID `gorm:"column:invocation_id;type:uuid;not null"`

	SenderID uuid.UUID `gorm:"column:sender_id;type:uuid;not null"`
	// Не сохраняется (только ID)
	Sender *User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`

	RecipientID uuid.UUID `gorm:"column:recipient_id;type:uuid;not null"`
	// Не сохраняется (только ID)
	Recipient *User `gorm:"foreignKey:RecipientID;constraint:OnDelete:CASCADE"`
}

func (MessageEquipmentInvocation) TableName() string {
	return "message_equipment_invocation"
}

type MessageEquipmentInvocationOptions struct {
	relations     []string
	WithSender    bool
	WithRecipient bool
}

func (o *MessageEquipmentInvocationOptions) Relations() []string {
	return o.relations
}

type MessageEquipmentInvocationOption func(options *MessageEquipmentInvocationOptions)

func MessageEquipmentInvocationWithSender() MessageEquipmentInvocationOption {
	return func(options *MessageEquipmentInvocationOptions) {
		options.WithSender = true
		options.relations = append(options.relations, "Sender")
	}
}

func MessageEquipmentInvocationWithRecipient() MessageEquipmentInvocationOption {
	return func(options *MessageEquipmentInvocationOptions) {
		options.WithRecipient = true
		options.relations = append(options.relations, "Recipient")
	}
}

type MessageEquipmentInvocationRepository interface {
	Get(ctx context.Context, id uuid.UUID, with ...MessageEquipmentInvocationOption) (*MessageEquipmentInvocation, error)
	GetInvocation(ctx context.Context, invocationID uuid.UUID, with ...MessageEquipmentInvocationOption) ([]*MessageEquipmentInvocation, error)
	Create(ctx context.Context, messageEquipmentInvocation *MessageEquipmentInvocation) error
	Update(ctx context.Context, messageEquipmentInvocation *MessageEquipmentInvocation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

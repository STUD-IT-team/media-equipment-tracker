package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageStudioInvocation struct {
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

func (MessageStudioInvocation) TableName() string {
	return "message_studio_invocation"
}

type MessageStudioInvocationOptions struct {
	relations     []string
	WithSender    bool
	WithRecipient bool
}

func (o *MessageStudioInvocationOptions) Relations() []string {
	return o.relations
}

type MessageStudioInvocationOption func(options *MessageStudioInvocationOptions)

func MessageStudioInvocationWithSender() MessageStudioInvocationOption {
	return func(options *MessageStudioInvocationOptions) {
		options.WithSender = true
		options.relations = append(options.relations, "Sender")
	}
}

func MessageStudioInvocationWithRecipient() MessageStudioInvocationOption {
	return func(options *MessageStudioInvocationOptions) {
		options.WithRecipient = true
		options.relations = append(options.relations, "Recipient")
	}
}

type MessageStudioInvocationRepository interface {
	Get(id uuid.UUID, with ...MessageStudioInvocationOption) (*MessageStudioInvocation, error)
	GetInvocation(invocationID uuid.UUID, with ...MessageStudioInvocationOption) ([]*MessageStudioInvocation, error)
	Create(messageStudioInvocation *MessageStudioInvocation) error
	Update(messageStudioInvocation *MessageStudioInvocation) error
	Delete(id uuid.UUID) error
}

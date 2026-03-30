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

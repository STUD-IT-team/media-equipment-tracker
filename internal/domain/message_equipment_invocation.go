package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageEquipmentInvocation struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Content  string    `gorm:"type:text;not null"`
	SentTime time.Time `gorm:"column:sent_time;type:timestamptz;default:CURRENT_TIMESTAMP"`

	//InvocationID uuid.UUID `gorm:"column:invocation_id;type:uuid;not null"`
	//SenderID     uuid.UUID `gorm:"column:sender_id;type:uuid;not null"`
	//RecipientID  uuid.UUID `gorm:"column:recipient_id;type:uuid;not null"`

	Invocation *EquipmentInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
	Sender     *User                `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Recipient  *User                `gorm:"foreignKey:RecipientID;constraint:OnDelete:CASCADE"`
}

func (MessageEquipmentInvocation) TableName() string {
	return "message_equipment_invocation"
}

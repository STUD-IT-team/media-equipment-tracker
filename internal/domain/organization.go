package domain

import (
	"github.com/google/uuid"
)

type Organization struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"type:varchar(255);not null"`

	// Не сохраняется
	Users []User `gorm:"many2many:user_organization;foreignKey:ID;joinForeignKey:organization_id;References:ID;joinReferences:user_id"`

	// Не сохраняется
	EquipmentInvocations []EquipmentInvocation `gorm:"foreignKey:OrganizationID"`

	// Не сохраняется
	StudioInvocations []StudioInvocation `gorm:"foreignKey:OrganizationID"`
}

func (Organization) TableName() string {
	return "organization"
}

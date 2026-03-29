package domain

import "github.com/google/uuid"

type EquipmentInInvocation struct {
	InvocationID uuid.UUID                   `gorm:"type:uuid;primaryKey"`
	EquipmentID  uuid.UUID                   `gorm:"type:uuid;primaryKey"`
	Status       EquipmentInInvocationStatus `gorm:"type:equipment_in_invocation_status"`

	//Invocation *EquipmentInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
	Equipment *Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE"`
}

func (EquipmentInInvocation) TableName() string {
	return "equipment_in_invocation"
}

type InvocationsForEquipment struct {
	InvocationID uuid.UUID                   `gorm:"type:uuid;primaryKey"`
	EquipmentID  uuid.UUID                   `gorm:"type:uuid;primaryKey"`
	Status       EquipmentInInvocationStatus `gorm:"type:equipment_in_invocation_status"`

	Invocation *EquipmentInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
	//Equipment  *Equipment           `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE"`
}

func (InvocationsForEquipment) TableName() string {
	return "equipment_in_invocation"
}

type EquipmentInInvocationStatus string

const (
	EquipmentNotIssued EquipmentInInvocationStatus = "not_issued"
	EquipmentIssued    EquipmentInInvocationStatus = "issued"
	EquipmentReturned  EquipmentInInvocationStatus = "returned"
)

package domain

import "github.com/google/uuid"

type Equipment struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primaryKey"`
	InventoryNumber     string          `gorm:"column:inventory_number;type:varchar(100);unique;not null;index:idx_equipment_inventory"`
	Name                string          `gorm:"type:varchar(255);not null"`
	ShortName           string          `gorm:"column:short_name;type:varchar(100)"`
	Category            string          `gorm:"type:varchar(100);not null"`
	AvailableToTrainee  bool            `gorm:"column:available_to_trainee;type:boolean;default:false"`
	Status              EquipmentStatus `gorm:"type:equipment_status"`
	CurrentInvocationID *uuid.UUID      `gorm:"column:current_invocation_id;type:uuid"`

	Departments []*Department `gorm:"many2many:equipment_department;foreignKey:ID;joinForeignKey:equipment_id;References:ID;joinReferences:department_id"`
	//CurrentInvocation      *EquipmentInvocation     `gorm:"foreignKey:CurrentInvocationID;constraint:OnDelete:SET NULL"`
	EquipmentInInvocations []*EquipmentInInvocation `gorm:"foreignKey:EquipmentID"`
}

func (Equipment) TableName() string {
	return "equipment"
}

type EquipmentStatus string

const (
	EquipmentStatusAvailable        EquipmentStatus = "available"
	EquipmentStatusIssued           EquipmentStatus = "issued"
	EquipmentStatusUnderMaintenance EquipmentStatus = "under_maintenance"
	EquipmentStatusUnavailable      EquipmentStatus = "unavailable"
)

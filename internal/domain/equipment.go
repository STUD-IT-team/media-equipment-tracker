package domain

import (
	"github.com/google/uuid"
)

type Equipment struct {
	ID                 uuid.UUID       `gorm:"type:uuid;primaryKey"`
	InventoryNumber    string          `gorm:"column:inventory_number;type:varchar(100);unique;not null;index:idx_equipment_inventory"`
	Name               string          `gorm:"type:varchar(255);not null"`
	ShortName          string          `gorm:"column:short_name;type:varchar(100)"`
	Category           string          `gorm:"type:varchar(100);not null"`
	AvailableToTrainee bool            `gorm:"column:available_to_trainee;type:boolean;default:false"`
	Status             EquipmentStatus `gorm:"type:equipment_status"`

	CurrentInvocationID *uuid.UUID `gorm:"column:current_invocation_id;type:uuid"`
	// Не сохраняется (только ID)
	CurrentInvocation *EquipmentInvocation `gorm:"foreignKey:CurrentInvocationID;constraint:OnDelete:SET NULL"`

	// Не сохраняется
	Invocations []*EquipmentInInvocation `gorm:"foreignKey:EquipmentID"`

	// Не сохраняется
	Departments []*Department `gorm:"many2many:equipment_department;foreignKey:ID;joinForeignKey:equipment_id;References:ID;joinReferences:department_id"`
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

type EquipmentOptions struct {
	relations                  []string
	withDepartments            bool
	withCurrentInvocation      bool
	withEquipmentInInvocations bool
}

type EquipmentOption func(options *EquipmentOptions)

func EquipmentWithDepartments() EquipmentOption {
	return func(options *EquipmentOptions) {
		options.withDepartments = true
		options.relations = append(options.relations, "Departments")
	}
}

func EquipmentWithCurrentInvocation() EquipmentOption {
	return func(options *EquipmentOptions) {
		options.withCurrentInvocation = true
		options.relations = append(options.relations, "CurrentInvocation")
	}
}

func EquipmentWithEquipmentInInvocations() EquipmentOption {
	return func(options *EquipmentOptions) {
		options.withEquipmentInInvocations = true
		options.relations = append(options.relations, "Invocations")
	}
}

type EquipmentRepository interface {
	Get(id uuid.UUID, with ...EquipmentOption) (*Equipment, error)
	GetUnoccupied(with ...EquipmentOption) ([]*Equipment, error)
	List(with ...EquipmentOption) ([]*Equipment, error)
	Reload(equipment *Equipment, with ...EquipmentOption) error

	Create(equipment *Equipment) error
	Update(equipment *Equipment) error
	Delete(id uuid.UUID) error
}

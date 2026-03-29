package domain

import (
	"github.com/google/uuid"
)

type Department struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"type:varchar(255);not null"`

	//Users                []User                `gorm:"many2many:user_department;foreignKey:ID;joinForeignKey:department_id;References:ID;joinReferences:user_id"`
	//Equipment            []Equipment           `gorm:"many2many:equipment_department;foreignKey:ID;joinForeignKey:department_id;References:ID;joinReferences:equipment_id"`
	//EquipmentInvocations []EquipmentInvocation `gorm:"foreignKey:DepartmentID"`
	//StudioInvocations    []StudioInvocation    `gorm:"foreignKey:DepartmentID"`
}

func (Department) TableName() string {
	return "department"
}

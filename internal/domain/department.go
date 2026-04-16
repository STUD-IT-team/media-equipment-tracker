package domain

import (
	"github.com/google/uuid"
)

type Department struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"type:varchar(255);not null"`

	// Не сохраняется
	Users []*User `gorm:"many2many:user_department;foreignKey:ID;joinForeignKey:department_id;References:ID;joinReferences:user_id"`

	// Сохраняется (только связь, не создаёт и не обновляет Equipment)
	Equipment []*Equipment `gorm:"many2many:equipment_department;foreignKey:ID;joinForeignKey:department_id;References:ID;joinReferences:equipment_id"`

	// Не сохраняется
	EquipmentInvocations []*EquipmentInvocation `gorm:"foreignKey:DepartmentID"`

	// Не сохраняется
	StudioInvocations []*StudioInvocation `gorm:"foreignKey:DepartmentID"`
}

func (Department) TableName() string {
	return "department"
}

type DepartmentOptions struct {
	relations                []string
	withUsers                bool
	withEquipment            bool
	withStudioInvocations    bool
	withEquipmentInvocations bool
}

func (o *DepartmentOptions) Relations() []string {
	return o.relations
}

type DepartmentOption func(options *DepartmentOptions)

func DepartmentWithUsers() DepartmentOption {
	return func(options *DepartmentOptions) {
		options.withUsers = true
		options.relations = append(options.relations, "Users")
	}
}

func DepartmentWithEquipment() DepartmentOption {
	return func(options *DepartmentOptions) {
		options.withEquipment = true
		options.relations = append(options.relations, "Equipment")
	}
}

func DepartmentWithStudioInvocations() DepartmentOption {
	return func(options *DepartmentOptions) {
		options.withStudioInvocations = true
		options.relations = append(options.relations, "StudioInvocations")
	}
}

func DepartmentWithEquipmentInvocations() DepartmentOption {
	return func(options *DepartmentOptions) {
		options.withEquipmentInvocations = true
		options.relations = append(options.relations, "EquipmentInvocations")
	}
}

type DepartmentRepository interface {
	Get(id uuid.UUID, with ...DepartmentOption) (*Department, error)
	List(with ...DepartmentOption) ([]*Department, error)
	Reload(dep *Department, with ...DepartmentOption) error
	Create(department *Department) error
	Update(department *Department) error
	Delete(id uuid.UUID) error
}

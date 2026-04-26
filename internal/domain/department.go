package domain

import (
	"context"

	"github.com/google/uuid"
)

type Department struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"type:varchar(255);not null"`

	// Не сохраняется
	Users []*UserDepartment `gorm:"foreignKey:DepartmentID"`

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
		options.relations = append(options.relations, "Users.User")
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
	Get(ctx context.Context, id uuid.UUID, with ...DepartmentOption) (*Department, error)
	List(ctx context.Context, with ...DepartmentOption) ([]*Department, error)
	Reload(ctx context.Context, dep *Department, with ...DepartmentOption) error
	Create(ctx context.Context, department *Department) error
	Update(ctx context.Context, department *Department) error
	Delete(ctx context.Context, id uuid.UUID) error
}

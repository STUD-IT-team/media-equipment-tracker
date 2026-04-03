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

type OrganizationOptions struct {
	relations                []string
	withUsers                bool
	withEquipmentInvocations bool
	withStudioInvocations    bool
}

func (o *OrganizationOptions) Relations() []string {
	return o.relations
}

type OrganizationOption func(options *OrganizationOptions)

func WithUsers() OrganizationOption {
	return func(options *OrganizationOptions) {
		options.withUsers = true
		options.relations = append(options.relations, "Users")
	}
}

func WithEquipmentInvocations() OrganizationOption {
	return func(options *OrganizationOptions) {
		options.withEquipmentInvocations = true
		options.relations = append(options.relations, "EquipmentInvocations")
	}
}

func WithStudioInvocations() OrganizationOption {
	return func(options *OrganizationOptions) {
		options.withStudioInvocations = true
		options.relations = append(options.relations, "StudioInvocations")
	}
}

type OrganizationRepository interface {
	Get(id uuid.UUID, with ...OrganizationOption) (*Organization, error)
	List(with ...OrganizationOption) ([]*Organization, error)
	Reload(organization *Organization, with ...OrganizationOption) error
	Create(organization *Organization) error
	Update(organization *Organization) error
	Delete(id uuid.UUID) error
}

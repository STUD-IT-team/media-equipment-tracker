package domain

import (
	"time"

	"github.com/google/uuid"
)

type EquipmentInvocation struct {
	ID                  uuid.UUID                 `gorm:"type:uuid;primaryKey"`
	EventName           string                    `gorm:"column:event_name;type:varchar(255);not null"`
	StartTime           time.Time                 `gorm:"column:start_time;type:timestamptz;not null;index:idx_equipment_invocation_start_time"`
	EndTime             time.Time                 `gorm:"column:end_time;type:timestamptz;not null;index:idx_equipment_invocation_end_time"`
	EquipmentReturnTime time.Time                 `gorm:"column:equipment_return_time;type:timestamptz;not null"`
	SdCardReturnTime    time.Time                 `gorm:"column:sd_card_return_time;type:timestamptz;not null"`
	Status              EquipmentInvocationStatus `gorm:"type:equipment_invocation_status"`
	CuratorComment      string                    `gorm:"column:curator_comment;type:text"`

	OrganizationID *uuid.UUID `gorm:"column:organization_id;type:uuid"`
	// Не сохраняется (только ID)
	Organization *Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"`

	DepartmentID *uuid.UUID `gorm:"column:department_id;type:uuid"`
	// Не сохраняется (только ID)
	Department *Department `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL"`

	UserID uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	// Не сохраняется (только ID)
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`

	AdminID *uuid.UUID `gorm:"column:admin_id;type:uuid"`
	// Не сохраняется (только ID)
	Admin *User `gorm:"foreignKey:AdminID;constraint:OnDelete:SET NULL"`

	// Сохраняется (только сама связь, не создаёт и не обновляет Equipment)
	Equipment []*EquipmentInInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
}

func (EquipmentInvocation) TableName() string {
	return "equipment_invocation"
}

type EquipmentInvocationStatus string

const (
	InvocationCreated           EquipmentInvocationStatus = "created"
	InvocationUnderReview       EquipmentInvocationStatus = "under_review"
	InvocationChangesRequired   EquipmentInvocationStatus = "changes_required"
	InvocationApproved          EquipmentInvocationStatus = "approved"
	InvocationEquipmentIssued   EquipmentInvocationStatus = "equipment_issued"
	InvocationEquipmentReturned EquipmentInvocationStatus = "equipment_returned"
	InvocationCompleted         EquipmentInvocationStatus = "completed"
	InvocationCancelled         EquipmentInvocationStatus = "cancelled"
)

type EquipmentInInvocationOptions struct {
	relations        []string
	withOrganization bool
	withDepartment   bool
	withUser         bool
	withAdmin        bool
	withEquipment    bool
}

type EquipmentInInvocationOption func(options *EquipmentInInvocationOptions)

func EquipmentInvocationWithOrganization() EquipmentInInvocationOption {
	return func(options *EquipmentInInvocationOptions) {
		options.withOrganization = true
		options.relations = append(options.relations, "Organization")
	}
}

func EquipmentInvocationWithDepartment() EquipmentInInvocationOption {
	return func(options *EquipmentInInvocationOptions) {
		options.withDepartment = true
		options.relations = append(options.relations, "Department")
	}
}

func EquipmentInvocationWithUser() EquipmentInInvocationOption {
	return func(options *EquipmentInInvocationOptions) {
		options.withUser = true
		options.relations = append(options.relations, "User")
	}
}

func EquipmentInvocationWithAdmin() EquipmentInInvocationOption {
	return func(options *EquipmentInInvocationOptions) {
		options.withAdmin = true
		options.relations = append(options.relations, "Admin")
	}
}

func EquipmentInvocationWithEquipment() EquipmentInInvocationOption {
	return func(options *EquipmentInInvocationOptions) {
		options.withEquipment = true
		options.relations = append(options.relations, "Equipment")
	}
}

type EquipmentInvocationRepository interface {
	Get(id uuid.UUID, with ...EquipmentInInvocationOption) (*EquipmentInvocation, error)
	List(with ...EquipmentInInvocationOption) ([]*EquipmentInvocation, error)
	Reload(equipmentInvocation *EquipmentInvocation, with ...EquipmentInInvocationOption) error

	Create(equipmentInvocation *EquipmentInvocation) error
	Update(equipmentInvocation *EquipmentInvocation) error
	Delete(id uuid.UUID) error
}

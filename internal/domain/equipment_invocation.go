package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EquipmentInvocation struct {
	ID                  uuid.UUID                 `gorm:"type:uuid;primaryKey"`
	EventName           string                    `gorm:"column:event_name;type:varchar(255);not null"`
	StartTime           time.Time                 `gorm:"column:start_time;type:timestamptz;not null;index:idx_equipment_invocation_start_time"`
	EndTime             time.Time                 `gorm:"column:end_time;type:timestamptz;not null;index:idx_equipment_invocation_end_time"`
	EquipmentReturnTime *time.Time                `gorm:"column:equipment_return_time;type:timestamptz"`
	SdCardReturnTime    *time.Time                `gorm:"column:sd_card_return_time;type:timestamptz"`
	Status              EquipmentInvocationStatus `gorm:"type:equipment_invocation_status"`
	CuratorComment      string                    `gorm:"column:curator_comment;type:text"`

	OrganizationID *uuid.UUID `gorm:"column:organization_id;type:uuid"`
	// Сохраняется (только сама связь, не создаёт и не обновляет Organization)
	Organization *Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"`

	DepartmentID *uuid.UUID `gorm:"column:department_id;type:uuid"`
	// Сохраняется (только сама связь, не создаёт и не обновляет Department)
	Department *Department `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL"`

	UserID uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	// Cохраняется (только сама связь, не создаёт и не обновляет User)
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`

	AdminID *uuid.UUID `gorm:"column:admin_id;type:uuid"`
	// Cохраняется (только сама связь, не создаёт и не обновляет User)
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

type EquipmentInvocationOptions struct {
	relations        []string
	withOrganization bool
	withDepartment   bool
	withUser         bool
	withAdmin        bool
	withEquipment    bool
}

func (o *EquipmentInvocationOptions) Relations() []string {
	return o.relations
}

type EquipmentInvocationOption func(options *EquipmentInvocationOptions)

func EquipmentInvocationWithOrganization() EquipmentInvocationOption {
	return func(options *EquipmentInvocationOptions) {
		options.withOrganization = true
		options.relations = append(options.relations, "Organization")
	}
}

func EquipmentInvocationWithDepartment() EquipmentInvocationOption {
	return func(options *EquipmentInvocationOptions) {
		options.withDepartment = true
		options.relations = append(options.relations, "Department")
	}
}

func EquipmentInvocationWithUser() EquipmentInvocationOption {
	return func(options *EquipmentInvocationOptions) {
		options.withUser = true
		options.relations = append(options.relations, "User")
	}
}

func EquipmentInvocationWithAdmin() EquipmentInvocationOption {
	return func(options *EquipmentInvocationOptions) {
		options.withAdmin = true
		options.relations = append(options.relations, "Admin")
	}
}

func EquipmentInvocationWithEquipment() EquipmentInvocationOption {
	return func(options *EquipmentInvocationOptions) {
		options.withEquipment = true
		options.relations = append(options.relations, "Equipment")
	}
}

type EquipmentInvocationRepository interface {
	Get(ctx context.Context, id uuid.UUID, with ...EquipmentInvocationOption) (*EquipmentInvocation, error)
	List(ctx context.Context, with ...EquipmentInvocationOption) ([]*EquipmentInvocation, error)
	Reload(ctx context.Context, equipmentInvocation *EquipmentInvocation, with ...EquipmentInvocationOption) error

	Create(ctx context.Context, equipmentInvocation *EquipmentInvocation) error
	Update(ctx context.Context, equipmentInvocation *EquipmentInvocation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

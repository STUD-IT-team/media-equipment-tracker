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

	//OrganizationID *uuid.UUID `gorm:"column:organization_id;type:uuid"`
	//DepartmentID   *uuid.UUID `gorm:"column:department_id;type:uuid"`
	//UserID         uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	//AdminID        *uuid.UUID `gorm:"column:admin_id;type:uuid"`

	Organization           *Organization            `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"`
	Department             *Department              `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL"`
	User                   *User                    `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	Admin                  *User                    `gorm:"foreignKey:AdminID;constraint:OnDelete:SET NULL"`
	EquipmentInInvocations []*EquipmentInInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
	//Messages               []*MessageEquipmentInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
	Equipments []Equipment `gorm:"many2many:equipment_in_invocation;foreignKey:ID;joinForeignKey:invocation_id;References:ID;joinReferences:equipment_id"`
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

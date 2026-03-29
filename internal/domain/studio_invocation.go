package domain

import (
	"time"

	"github.com/google/uuid"
)

type StudioInvocation struct {
	ID                  uuid.UUID              `gorm:"type:uuid;primaryKey"`
	EventName           string                 `gorm:"column:event_name;type:varchar(255);not null"`
	ShootingDescription string                 `gorm:"column:shooting_description;type:text;not null"`
	StartTime           time.Time              `gorm:"column:start_time;type:timestamptz;not null;index:idx_studio_invocation_start_time"`
	EndTime             time.Time              `gorm:"column:end_time;type:timestamptz;not null;index:idx_studio_invocation_end_time"`
	NeedsChromakey      bool                   `gorm:"column:needs_chromakey;type:boolean;default:false"`
	NeedsCyclorama      bool                   `gorm:"column:needs_cyclorama;type:boolean;default:false"`
	NeedsBlackFabric    bool                   `gorm:"column:needs_black_fabric;type:boolean;default:false"`
	Status              StudioInvocationStatus `gorm:"type:studio_invocation_status;not null"`
	CuratorComment      string                 `gorm:"column:curator_comment;type:text"`

	//OrganizationID uuid.UUID `gorm:"column:organization_id;type:uuid;not null"`
	//DepartmentID   uuid.UUID `gorm:"column:department_id;type:uuid;not null"`
	//UserID         uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	//AdminID        uuid.UUID `gorm:"column:admin_id;type:uuid;not null"`

	Organization *Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"`
	Department   *Department   `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL"`
	User         *User         `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	Admin        *User         `gorm:"foreignKey:AdminID;constraint:OnDelete:SET NULL"`
	//Messages     []MessageStudioInvocation `gorm:"foreignKey:InvocationID;constraint:OnDelete:CASCADE"`
}

func (StudioInvocation) TableName() string {
	return "studio_invocation"
}

type StudioInvocationStatus string

const (
	StudioCreated         StudioInvocationStatus = "created"
	StudioUnderReview     StudioInvocationStatus = "under_review"
	StudioChangesRequired StudioInvocationStatus = "changes_required"
	StudioApproved        StudioInvocationStatus = "approved"
	StudioCompleted       StudioInvocationStatus = "completed"
	StudioCancelled       StudioInvocationStatus = "cancelled"
)

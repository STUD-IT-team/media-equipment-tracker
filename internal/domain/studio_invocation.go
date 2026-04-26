package domain

import (
	"context"
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
	CuratorComment      *string                `gorm:"column:curator_comment;type:text"`

	OrganizationID *uuid.UUID `gorm:"column:organization_id;type:uuid"`
	// Не сохраняется (только ID)
	Organization *Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"`

	DepartmentID *uuid.UUID `gorm:"column:department_id;type:uuid"`
	// Не сохраняется (только ID)
	Department *Department `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL"`

	AdminID *uuid.UUID `gorm:"column:admin_id;type:uuid"`
	// Не сохраняется (только ID)
	Admin *User `gorm:"foreignKey:AdminID;constraint:OnDelete:SET NULL"`

	UserID uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	// Не сохраняется (только ID)
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
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

type StudioInvocationOptions struct {
	relations        []string
	withOrganization bool
	withDepartment   bool
	withAdmin        bool
	withUser         bool
}

func (o *StudioInvocationOptions) Relations() []string {
	return o.relations
}

type StudioInvocationOption func(options *StudioInvocationOptions)

func StudioInvocationWithOrganization() StudioInvocationOption {
	return func(options *StudioInvocationOptions) {
		options.withOrganization = true
		options.relations = append(options.relations, "Organization")
	}
}

func StudioInvocationWithDepartment() StudioInvocationOption {
	return func(options *StudioInvocationOptions) {
		options.withDepartment = true
		options.relations = append(options.relations, "Department")
	}
}

func StudioInvocationWithAdmin() StudioInvocationOption {
	return func(options *StudioInvocationOptions) {
		options.withAdmin = true
		options.relations = append(options.relations, "Admin")
	}
}

func StudioInvocationWithUser() StudioInvocationOption {
	return func(options *StudioInvocationOptions) {
		options.withUser = true
		options.relations = append(options.relations, "User")
	}
}

type StudioInvocationRepository interface {
	Get(ctx context.Context, id uuid.UUID, with ...StudioInvocationOption) (*StudioInvocation, error)
	List(ctx context.Context, with ...StudioInvocationOption) ([]*StudioInvocation, error)
	Reload(ctx context.Context, studioInvocation *StudioInvocation, with ...StudioInvocationOption) error

	Create(ctx context.Context, studioInvocation *StudioInvocation) error
	Update(ctx context.Context, studioInvocation *StudioInvocation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

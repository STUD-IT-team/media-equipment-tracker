package dto

import (
	"time"

	"media-equipment-tracker/internal/domain"
)

type ShortStudioItem struct {
	ID                  string  `json:"id"`
	EventName           string  `json:"event_name"`
	ShootingDescription string  `json:"shooting_description"`
	StartTime           string  `json:"start_time"`
	EndTime             string  `json:"end_time"`
	Status              string  `json:"status"`
	CuratorComment      *string `json:"curator_comment,omitempty"`
	AdminID             *string `json:"admin_id,omitempty"`
	UserID              string  `json:"user_id"`
	DepartmentID        *string `json:"department_id,omitempty"`
	OrganizationID      *string `json:"organization_id,omitempty"`
}

func ShortFromStudio(studio *domain.StudioInvocation) ShortStudioItem {
	return ShortStudioItem{
		ID:                  studio.ID.String(),
		EventName:           studio.EventName,
		ShootingDescription: studio.ShootingDescription,
		StartTime:           studio.StartTime.Format(time.RFC3339),
		EndTime:             studio.EndTime.Format(time.RFC3339),
		Status:              string(studio.Status),
		CuratorComment:      studio.CuratorComment,
		AdminID: func() *string {
			if studio.AdminID == nil {
				return nil
			}
			s := studio.AdminID.String()
			return &s
		}(),
		UserID: studio.UserID.String(),
		DepartmentID: func() *string {
			if studio.DepartmentID == nil {
				return nil
			}
			s := studio.DepartmentID.String()
			return &s
		}(),
		OrganizationID: func() *string {
			if studio.OrganizationID == nil {
				return nil
			}
			s := studio.OrganizationID.String()
			return &s
		}(),
	}
}

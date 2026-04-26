package dto

import (
	"time"

	"media-equipment-tracker/internal/domain"
)

type ShortInvocationItem struct {
	ID             string  `json:"id"`
	EventName      string  `json:"event_name"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	Status         string  `json:"status"`
	CuratorComment *string `json:"curator,omitempty"`
	AdminID        *string `json:"admin_id,omitempty"`
	UserID         string  `json:"user_id"`
	DepartmentID   *string `json:"department_id,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	EquipmentCount int     `json:"equipment_count"`
}

func ShortFromInvocation(inv *domain.EquipmentInvocation) ShortInvocationItem {
	return ShortInvocationItem{
		ID:        inv.ID.String(),
		EventName: inv.EventName,
		StartTime: inv.StartTime.Format(time.RFC3339),
		EndTime:   inv.EndTime.Format(time.RFC3339),
		Status:    string(inv.Status),
		CuratorComment: func() *string {
			if inv.CuratorComment == "" {
				return nil
			}
			return &inv.CuratorComment
		}(),
		AdminID: func() *string {
			if inv.AdminID == nil {
				return nil
			}
			s := inv.AdminID.String()
			return &s
		}(),
		UserID: inv.UserID.String(),
		DepartmentID: func() *string {
			if inv.DepartmentID == nil {
				return nil
			}
			s := inv.DepartmentID.String()
			return &s
		}(),
		OrganizationID: func() *string {
			if inv.OrganizationID == nil {
				return nil
			}
			s := inv.OrganizationID.String()
			return &s
		}(),
		EquipmentCount: len(inv.Equipment),
	}
}

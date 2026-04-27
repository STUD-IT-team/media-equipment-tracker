package dto

import (
	"time"

	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeserializeGetStudioRequest(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

type GetStudioResponse struct {
	ID                  string                 `json:"id"`
	EventName           string                 `json:"event_name"`
	ShootingDescription string                 `json:"shooting_description"`
	StartTime           string                 `json:"start_time"`
	EndTime             string                 `json:"end_time"`
	Status              string                 `json:"status"`
	CuratorComment      *string                `json:"curator_comment,omitempty"`
	Admin               *getStudioAdmin        `json:"admin,omitempty"`
	User                *getStudioUser         `json:"user"`
	Department          *getStudioDepartment   `json:"department,omitempty"`
	Organization        *getStudioOrganization `json:"organization,omitempty"`
}

type getStudioUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
	Nice    int    `json:"nice"`
}

type getStudioDepartment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getStudioOrganization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getStudioAdmin struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
	Nice    int    `json:"nice"`
}

func SerializeGetStudioResponse(_ *gin.Context, inv *domain.StudioInvocation) any {
	return GetStudioResponse{
		ID:                  inv.ID.String(),
		EventName:           inv.EventName,
		ShootingDescription: inv.ShootingDescription,
		StartTime:           inv.StartTime.Format(time.RFC3339),
		EndTime:             inv.EndTime.Format(time.RFC3339),
		Status:              string(inv.Status),
		CuratorComment:      inv.CuratorComment,
		Admin: func() *getStudioAdmin {
			if inv.AdminID == nil {
				return nil
			}
			admin := getStudioAdmin{
				ID:      inv.AdminID.String(),
				Name:    inv.Admin.FullName,
				IsAdmin: inv.Admin.IsAdmin,
				Nice:    inv.Admin.Nice,
			}
			return &admin
		}(),
		User: func() *getStudioUser {
			if inv.UserID == uuid.Nil {
				return nil
			}
			user := getStudioUser{
				ID:      inv.UserID.String(),
				Name:    inv.User.FullName,
				IsAdmin: inv.User.IsAdmin,
				Nice:    inv.User.Nice,
			}
			return &user
		}(),
		Department: func() *getStudioDepartment {
			if inv.DepartmentID == nil {
				return nil
			}
			department := getStudioDepartment{
				ID:   inv.DepartmentID.String(),
				Name: inv.Department.Name,
			}
			return &department
		}(),
		Organization: func() *getStudioOrganization {
			if inv.OrganizationID == nil {
				return nil
			}
			organization := getStudioOrganization{
				ID:   inv.OrganizationID.String(),
				Name: inv.Organization.Name,
			}
			return &organization
		}(),
	}
}

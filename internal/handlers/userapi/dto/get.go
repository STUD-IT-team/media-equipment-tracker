package dto

import (
	"media-equipment-tracker/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetUserResponse struct {
	ID                        string                   `json:"id"`
	FullName                  string                   `json:"full_name"`
	Email                     string                   `json:"email"`
	Nice                      int                      `json:"nice"`
	IsAdmin                   bool                     `json:"is_admin"`
	Organizations             []getUserOrganizations   `json:"organizations"`
	Departments               []getUserDepartments     `json:"departments"`
	StudioInvocations         []getStudioInvocation    `json:"studio_invocations"`
	EquipmentInvocations      []getEquipmentInvocation `json:"equipment_invocations"`
	AdminStudioInvocations    []getStudioInvocation    `json:"admin_studio_invocations"`
	AdminEquipmentInvocations []getEquipmentInvocation `json:"admin_equipment_invocations"`
}

type getUserOrganizations struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getUserDepartments struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type getStudioInvocation struct {
	ID                  string `json:"id"`
	EventName           string `json:"eventName"`
	ShootingDescription string `json:"shootingDescription"`
	Status              string `json:"status"`
	StartTime           string `json:"startTime"`
	EndTime             string `json:"endTime"`
}

type getEquipmentInvocation struct {
	ID        string `json:"id"`
	EventName string `json:"eventName"`
	Status    string `json:"status"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func DeserializeGetUserRequest(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func SerializeGetUserResponse(_ *gin.Context, user *domain.User) any {
	resp := GetUserResponse{
		ID:                        user.ID.String(),
		FullName:                  user.FullName,
		Email:                     user.Email,
		Nice:                      user.Nice,
		IsAdmin:                   user.IsAdmin,
		Organizations:             make([]getUserOrganizations, 0, len(user.Organizations)),
		Departments:               make([]getUserDepartments, 0, len(user.Departments)),
		StudioInvocations:         make([]getStudioInvocation, 0, len(user.StudioInvocations)),
		EquipmentInvocations:      make([]getEquipmentInvocation, 0, len(user.EquipmentInvocations)),
		AdminStudioInvocations:    make([]getStudioInvocation, 0, len(user.AdminStudioInvocations)),
		AdminEquipmentInvocations: make([]getEquipmentInvocation, 0, len(user.AdminEquipmentInvocations)),
	}

	for _, org := range user.Organizations {
		resp.Organizations = append(resp.Organizations, getUserOrganizations{
			ID:   org.ID.String(),
			Name: org.Name,
		})
	}
	for _, department := range user.Departments {
		resp.Departments = append(resp.Departments, getUserDepartments{
			ID:   department.DepartmentID.String(),
			Name: department.Department.Name,
			Role: string(department.Role),
		})
	}
	for _, inv := range user.StudioInvocations {
		resp.StudioInvocations = append(resp.StudioInvocations, getStudioInvocation{
			ID:                  inv.ID.String(),
			EventName:           inv.EventName,
			ShootingDescription: inv.ShootingDescription,
			Status:              string(inv.Status),
			StartTime:           inv.StartTime.Format(time.RFC3339),
			EndTime:             inv.EndTime.Format(time.RFC3339),
		})
	}
	for _, inv := range user.EquipmentInvocations {
		resp.EquipmentInvocations = append(resp.EquipmentInvocations, getEquipmentInvocation{
			ID:        inv.ID.String(),
			EventName: inv.EventName,
			Status:    string(inv.Status),
			StartTime: inv.StartTime.Format(time.RFC3339),
			EndTime:   inv.EndTime.Format(time.RFC3339),
		})
	}
	for _, inv := range user.AdminStudioInvocations {
		resp.AdminStudioInvocations = append(resp.AdminStudioInvocations, getStudioInvocation{
			ID:                  inv.ID.String(),
			EventName:           inv.EventName,
			ShootingDescription: inv.ShootingDescription,
			Status:              string(inv.Status),
			StartTime:           inv.StartTime.Format(time.RFC3339),
			EndTime:             inv.EndTime.Format(time.RFC3339),
		})
	}
	for _, inv := range user.AdminEquipmentInvocations {
		resp.AdminEquipmentInvocations = append(resp.AdminEquipmentInvocations, getEquipmentInvocation{
			ID:        inv.ID.String(),
			EventName: inv.EventName,
			Status:    string(inv.Status),
			StartTime: inv.StartTime.Format(time.RFC3339),
			EndTime:   inv.EndTime.Format(time.RFC3339),
		})
	}

	return resp
}

package dto

import (
	"fmt"
	"time"

	"media-equipment-tracker/internal/application/userservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

type MeUserResponse struct {
	ID            string                `json:"id"`
	FullName      string                `json:"full_name"`
	Email         string                `json:"email"`
	Nice          int                   `json:"nice"`
	IsAdmin       bool                  `json:"is_admin"`
	Organizations []meUserOrganizations `json:"organizations"`
	Departments   []meUserDepartments   `json:"departments"`
}

type meUserOrganizations struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type meUserDepartments struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func SerializeMeUserResponse(_ *gin.Context, user *domain.User) any {
	resp := MeUserResponse{
		ID:            user.ID.String(),
		FullName:      user.FullName,
		Email:         user.Email,
		Nice:          user.Nice,
		IsAdmin:       user.IsAdmin,
		Organizations: make([]meUserOrganizations, 0, len(user.Organizations)),
		Departments:   make([]meUserDepartments, 0, len(user.Departments)),
	}
	if user.Organizations != nil {
		for _, org := range user.Organizations {
			resp.Organizations = append(resp.Organizations, meUserOrganizations{
				ID:   org.ID.String(),
				Name: org.Name,
			})
		}
	}

	if user.Departments != nil {
		for _, department := range user.Departments {
			resp.Departments = append(resp.Departments, meUserDepartments{
				ID:   department.DepartmentID.String(),
				Name: department.Department.Name,
				Role: string(department.Role),
			})
		}
	}

	return resp
}

type UpdateSelfRequest struct {
	FullName *string `json:"full_name,omitempty" example:"Иванов Иван Иванович"`
	Email    *string `json:"email,omitempty" example:"ivan@example.com"`
}

type UpdateSelfResponse struct {
	ShortUserItem
}

func DeserializeUpdateSelfRequest(c *gin.Context) (*userservice.UpdateSelfRequest, error) {
	var req UpdateSelfRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, err
	}
	return &userservice.UpdateSelfRequest{
		FullName: req.FullName,
		Email:    req.Email,
	}, nil
}

func SerializeUpdateSelfResponse(_ *gin.Context, user *domain.User) any {
	resp := UpdateSelfResponse{
		ShortUserItem: ShortFromUser(user),
	}
	return resp
}

type MyInvocationsRequest struct {
	Statuses []string `form:"statuses,omitempty"`
	Type     string   `form:"type,omitempty" binding:"required"`
}

type MyInvocationsResponse struct {
	EquipmentInvocations      []myEquipmentInvocation `json:"equipment_invocations"`
	StudioInvocations         []myStudioInvocation    `json:"studio_invocations"`
	AdminEquipmentInvocations []myEquipmentInvocation `json:"admin_equipment_invocations"`
	AdminStudioInvocations    []myStudioInvocation    `json:"admin_studio_invocations"`
}

type myEquipmentInvocation struct {
	ID        string `json:"id"`
	EventName string `json:"event_name"`
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type myStudioInvocation struct {
	ID                  string `json:"id"`
	EventName           string `json:"event_name"`
	ShootingDescription string `json:"shooting_description"`
	Status              string `json:"status"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
}

func DeserializeMyInvocationsRequest(c *gin.Context) (*userservice.MyInvocationsRequest, error) {
	var req MyInvocationsRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return nil, err
	}

	t, err := mapType(req.Type)
	if err != nil {
		return nil, err
	}
	var equipmentStatuses []domain.EquipmentInvocationStatus
	if t == userservice.All || t == userservice.Equipment {
		equipmentStatuses := make([]domain.EquipmentInvocationStatus, 0, len(req.Statuses))
		for _, status := range req.Statuses {
			s, err := mapEquipmentStatus(status)
			if err != nil {
				return nil, err
			}
			equipmentStatuses = append(equipmentStatuses, s)
		}
	}

	var studioStatuses []domain.StudioInvocationStatus
	if t == userservice.All || t == userservice.Studio {
		studioStatuses := make([]domain.StudioInvocationStatus, 0, len(req.Statuses))
		for _, status := range req.Statuses {
			s, err := mapStudioStatus(status)
			if err != nil {
				return nil, err
			}
			studioStatuses = append(studioStatuses, s)
		}
	}

	return &userservice.MyInvocationsRequest{
		StudioStatuses:    studioStatuses,
		EquipmentStatuses: equipmentStatuses,
		Type:              t,
	}, nil
}

func SerializeMyInvocationsResponse(_ *gin.Context, invs *userservice.MyInvocationsResponse) any {
	resp := MyInvocationsResponse{
		EquipmentInvocations:      make([]myEquipmentInvocation, 0, len(invs.EquipmentInvocations)),
		StudioInvocations:         make([]myStudioInvocation, 0, len(invs.StudioInvocations)),
		AdminEquipmentInvocations: make([]myEquipmentInvocation, 0, len(invs.AdminEquipmentInvocations)),
		AdminStudioInvocations:    make([]myStudioInvocation, 0, len(invs.AdminStudioInvocations)),
	}

	if invs.EquipmentInvocations != nil {
		for _, inv := range invs.EquipmentInvocations {
			resp.EquipmentInvocations = append(resp.EquipmentInvocations, myEquipmentInvocation{
				ID:        inv.ID.String(),
				EventName: inv.EventName,
				Status:    string(inv.Status),
				StartTime: inv.StartTime.Format(time.RFC3339),
				EndTime:   inv.EndTime.Format(time.RFC3339),
			})
		}
	}

	if invs.StudioInvocations != nil {
		for _, inv := range invs.StudioInvocations {
			resp.StudioInvocations = append(resp.StudioInvocations, myStudioInvocation{
				ID:                  inv.ID.String(),
				EventName:           inv.EventName,
				ShootingDescription: inv.ShootingDescription,
				Status:              string(inv.Status),
				StartTime:           inv.StartTime.Format(time.RFC3339),
				EndTime:             inv.EndTime.Format(time.RFC3339),
			})
		}
	}
	if invs.AdminEquipmentInvocations != nil {
		for _, inv := range invs.AdminEquipmentInvocations {
			resp.AdminEquipmentInvocations = append(resp.AdminEquipmentInvocations, myEquipmentInvocation{
				ID:        inv.ID.String(),
				EventName: inv.EventName,
				Status:    string(inv.Status),
				StartTime: inv.StartTime.Format(time.RFC3339),
				EndTime:   inv.EndTime.Format(time.RFC3339),
			})
		}
	}
	if invs.AdminStudioInvocations != nil {
		for _, inv := range invs.AdminStudioInvocations {
			resp.AdminStudioInvocations = append(resp.AdminStudioInvocations, myStudioInvocation{
				ID:                  inv.ID.String(),
				EventName:           inv.EventName,
				ShootingDescription: inv.ShootingDescription,
				Status:              string(inv.Status),
				StartTime:           inv.StartTime.Format(time.RFC3339),
				EndTime:             inv.EndTime.Format(time.RFC3339),
			})
		}
	}

	return resp
}

func mapStudioStatus(status string) (domain.StudioInvocationStatus, error) {
	switch status {
	case "created":
		return domain.StudioCreated, nil
	case "under_review":
		return domain.StudioUnderReview, nil
	case "changes_required":
		return domain.StudioChangesRequired, nil
	case "approved":
		return domain.StudioApproved, nil
	case "completed":
		return domain.StudioCompleted, nil
	case "cancelled":
		return domain.StudioCancelled, nil
	default:
		return "", fmt.Errorf("invalid status %s", status)
	}
}

func mapEquipmentStatus(status string) (domain.EquipmentInvocationStatus, error) {
	switch status {
	case "created":
		return domain.InvocationCreated, nil
	case "under_review":
		return domain.InvocationUnderReview, nil
	case "changes_required":
		return domain.InvocationChangesRequired, nil
	case "approved":
		return domain.InvocationApproved, nil
	case "equipment_issued":
		return domain.InvocationEquipmentIssued, nil
	case "equipment_returned":
		return domain.InvocationEquipmentReturned, nil
	case "completed":
		return domain.InvocationCompleted, nil
	case "cancelled":
		return domain.InvocationCancelled, nil
	default:
		return "", fmt.Errorf("invalid status %s", status)
	}
}

func mapType(t string) (userservice.InvocationType, error) {
	switch t {
	case "studio":
		return userservice.Studio, nil
	case "equipment":
		return userservice.Equipment, nil
	case "all":
		return userservice.All, nil
	default:
		return "", fmt.Errorf("invalid type %s", t)
	}
}

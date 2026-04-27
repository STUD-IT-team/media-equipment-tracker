package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetDepartmentRequest struct {
	ID string `uri:"id" binding:"required"`
}

type DepartmentDetailResponse struct {
	ID        string                   `json:"id"`
	Name      string                   `json:"name"`
	Users     []UserDepartmentResponse `json:"users,omitempty"`
	Equipment []EquipmentResponse      `json:"equipment,omitempty"`
}

type UserDepartmentResponse struct {
	User UserResponse `json:"user"`
	Role string       `json:"role"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type EquipmentResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ShortName       string `json:"short_name"`
	InventoryNumber string `json:"inventory_number"`
	Category        string `json:"category"`
	Status          string `json:"status"`
}

func DeserializeGetDepartmentRequest(c *gin.Context) (uuid.UUID, error) {
	var req GetDepartmentRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(req.ID)
}

func SerializeGetDepartmentResponse(_ *gin.Context, department *domain.Department) any {
	resp := DepartmentDetailResponse{
		ID:   department.ID.String(),
		Name: department.Name,
	}

	if department.Users != nil {
		users := make([]UserDepartmentResponse, 0, len(department.Users))
		for _, userDep := range department.Users {
			users = append(users, UserDepartmentResponse{
				User: UserResponse{
					ID:       userDep.User.ID.String(),
					Email:    userDep.User.Email,
					FullName: userDep.User.FullName,
				},
				Role: string(userDep.Role),
			})
		}
		resp.Users = users
	}

	if department.Equipment != nil {
		equipment := make([]EquipmentResponse, 0, len(department.Equipment))
		for _, eq := range department.Equipment {
			equipment = append(equipment, EquipmentResponse{
				ID:              eq.ID.String(),
				Name:            eq.Name,
				ShortName:       eq.ShortName,
				InventoryNumber: eq.InventoryNumber,
				Category:        eq.Category,
				Status:          string(eq.Status),
			})
		}
		resp.Equipment = equipment
	}

	return resp
}

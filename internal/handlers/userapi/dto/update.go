package dto

import (
	"fmt"

	"media-equipment-tracker/internal/application/userservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	FullName         *string                `json:"full_name,omitempty" example:"Иванов Иван Иванович"`
	Email            *string                `json:"email,omitempty" example:"ivan@example.com"`
	Nice             *int                   `json:"nice,omitempty" example:"150"`
	IsAdmin          *bool                  `json:"is_admin,omitempty"`
	OrganizationsIDs []string               `json:"organizations_ids,omitempty"`
	Departments      []updateUserDepartment `json:"departments,omitempty"`
}

type updateUserDepartment struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type UpdateUserResponse struct {
	ShortUserItem
}

func DeserializeUpdateUserRequest(c *gin.Context) (*userservice.UpdateUserRequest, error) {
	var req UpdateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return nil, err
	}

	organizationIDs := make([]uuid.UUID, 0, len(req.OrganizationsIDs))
	for _, id := range req.OrganizationsIDs {
		id, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		organizationIDs = append(organizationIDs, id)
	}

	departments := make([]userservice.UpdateUserDepartment, 0, len(req.Departments))
	for _, d := range req.Departments {
		id, err := uuid.Parse(d.ID)
		if err != nil {
			return nil, err
		}
		role, err := mapRole(d.Role)
		if err != nil {
			return nil, err
		}
		departments = append(departments, userservice.UpdateUserDepartment{
			ID:   id,
			Role: role,
		})
	}

	return &userservice.UpdateUserRequest{
		ID:              id,
		FullName:        req.FullName,
		Email:           req.Email,
		Nice:            req.Nice,
		IsAdmin:         req.IsAdmin,
		OrganizationIDs: organizationIDs,
		Departments:     departments,
	}, nil
}

func mapRole(role string) (domain.RoleInDepartment, error) {
	switch role {
	case "trainee":
		return domain.RoleTrainee, nil
	case "activist":
		return domain.RoleActivist, nil
	default:
		return "", fmt.Errorf("invalid role: %s", role)
	}
}

func SerializeUpdateUserResponse(_ *gin.Context, user *domain.User) any {
	resp := UpdateUserResponse{
		ShortUserItem: ShortFromUser(user),
	}
	return resp
}

type UpdateNiceRequest struct {
	Nice int `json:"nice" binding:"required"`
}

type UpdateNiceResponse struct {
	ShortUserItem
}

func DeserializeUpdateNiceRequest(c *gin.Context) (*userservice.UpdateNiceRequest, error) {
	var req UpdateNiceRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return nil, err
	}

	return &userservice.UpdateNiceRequest{
		ID:   id,
		Nice: req.Nice,
	}, nil
}

func SerializeUpdateNiceResponse(_ *gin.Context, user *domain.User) any {
	resp := UpdateNiceResponse{
		ShortUserItem: ShortFromUser(user),
	}
	return resp
}

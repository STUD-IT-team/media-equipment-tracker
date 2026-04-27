package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetOrganizationRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type OrganizationDetailResponse struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Users []UserResponse `json:"users,omitempty"`
}

func DeserializeGetOrganizationRequest(c *gin.Context) (uuid.UUID, error) {
	var req GetOrganizationRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(req.ID)
}

func SerializeGetOrganizationResponse(_ *gin.Context, organization *domain.Organization) any {
	resp := OrganizationDetailResponse{
		ID:   organization.ID.String(),
		Name: organization.Name,
	}

	if organization.Users != nil {
		users := make([]UserResponse, 0, len(organization.Users))
		for _, user := range organization.Users {
			users = append(users, UserResponse{
				ID:       user.ID.String(),
				Email:    user.Email,
				FullName: user.FullName,
			})
		}
		resp.Users = users
	}

	return resp
}

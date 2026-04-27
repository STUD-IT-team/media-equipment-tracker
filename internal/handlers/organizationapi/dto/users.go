package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetOrganizationUsersRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetOrganizationUsersResponse struct {
	Items []UserResponse `json:"items"`
}

func DeserializeGetOrganizationUsersRequest(c *gin.Context) (uuid.UUID, error) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func SerializeGetOrganizationUsersResponse(_ *gin.Context, users []*domain.User) any {
	items := make([]UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, UserResponse{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
		})
	}
	return GetOrganizationUsersResponse{
		Items: items,
	}
}

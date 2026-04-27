package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

type GetAllUsersResponse struct {
	Items []ShortUserItem `json:"items"`
}

func SerializeGetAllUsersResponse(_ *gin.Context, users []*domain.User) any {
	resp := GetAllUsersResponse{
		Items: make([]ShortUserItem, len(users)),
	}
	for i, user := range users {
		resp.Items[i] = ShortFromUser(user)
	}
	return resp
}

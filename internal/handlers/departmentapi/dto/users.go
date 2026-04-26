package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetDepartmentUsersRequest struct {
	ID   string  `uri:"id" binding:"required"`
	Role *string `form:"role"`
}

type GetDepartmentUsersResponse struct {
	Items []UserDepartmentResponse `json:"items"`
}

func DeserializeGetDepartmentUsersRequest(c *gin.Context) (uuid.UUID, *domain.RoleInDepartment, error) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, nil, err
	}

	var req GetDepartmentUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return uuid.Nil, nil, err
	}

	var role *domain.RoleInDepartment
	if req.Role != nil {
		r := domain.RoleInDepartment(*req.Role)
		role = &r
	}

	return id, role, nil
}

func SerializeGetDepartmentUsersResponse(_ *gin.Context, users []*domain.UserDepartment) any {
	items := make([]UserDepartmentResponse, 0, len(users))
	for _, userDep := range users {
		items = append(items, UserDepartmentResponse{
			User: UserResponse{
				ID:       userDep.User.ID.String(),
				Email:    userDep.User.Email,
				FullName: userDep.User.FullName,
			},
			Role: string(userDep.Role),
		})
	}
	return GetDepartmentUsersResponse{
		Items: items,
	}
}

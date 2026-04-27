package dto

import (
	"media-equipment-tracker/internal/application/departmentservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateDepartmentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DeserializeUpdateDepartmentRequest(c *gin.Context) (uuid.UUID, departmentservice.UpdateDepartmentRequest, error) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, departmentservice.UpdateDepartmentRequest{}, err
	}

	var req UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return uuid.Nil, departmentservice.UpdateDepartmentRequest{}, err
	}

	return id, departmentservice.UpdateDepartmentRequest{
		Name: req.Name,
	}, nil
}

func SerializeUpdateDepartmentResponse(_ *gin.Context, department *domain.Department) any {
	return UpdateDepartmentResponse{
		ID:   department.ID.String(),
		Name: department.Name,
	}
}

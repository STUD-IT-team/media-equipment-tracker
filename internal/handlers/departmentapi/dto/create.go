package dto

import (
	"media-equipment-tracker/internal/application/departmentservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

type CreateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateDepartmentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DeserializeCreateDepartmentRequest(c *gin.Context) (departmentservice.CreateDepartmentRequest, error) {
	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return departmentservice.CreateDepartmentRequest{}, err
	}

	return departmentservice.CreateDepartmentRequest{
		Name: req.Name,
	}, nil
}

func SerializeCreateDepartmentResponse(_ *gin.Context, department *domain.Department) any {
	return CreateDepartmentResponse{
		ID:   department.ID.String(),
		Name: department.Name,
	}
}

package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

type ListDepartmentsRequest struct {
	Search string `form:"search"`
}

type DepartmentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListDepartmentsResponse struct {
	Items []DepartmentResponse `json:"items"`
}

func DeserializeListDepartmentsRequest(c *gin.Context) (ListDepartmentsRequest, error) {
	var req ListDepartmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return ListDepartmentsRequest{}, err
	}
	return req, nil
}

func SerializeListDepartmentsResponse(_ *gin.Context, departments []*domain.Department) any {
	items := make([]DepartmentResponse, 0, len(departments))
	for _, dep := range departments {
		items = append(items, DepartmentResponse{
			ID:   dep.ID.String(),
			Name: dep.Name,
		})
	}
	return ListDepartmentsResponse{
		Items: items,
	}
}

package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

type ListOrganizationsRequest struct {
	Search string `form:"search"`
}

type OrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListOrganizationsResponse struct {
	Items []OrganizationResponse `json:"items"`
}

func DeserializeListOrganizationsRequest(c *gin.Context) (ListOrganizationsRequest, error) {
	var req ListOrganizationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return ListOrganizationsRequest{}, err
	}
	return req, nil
}

func SerializeListOrganizationsResponse(_ *gin.Context, organizations []*domain.Organization) any {
	items := make([]OrganizationResponse, 0, len(organizations))
	for _, org := range organizations {
		items = append(items, OrganizationResponse{
			ID:   org.ID.String(),
			Name: org.Name,
		})
	}
	return ListOrganizationsResponse{
		Items: items,
	}
}

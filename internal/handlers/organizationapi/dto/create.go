package dto

import (
	"media-equipment-tracker/internal/application/organizationservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateOrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DeserializeCreateOrganizationRequest(c *gin.Context) (organizationservice.CreateOrganizationRequest, error) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return organizationservice.CreateOrganizationRequest{}, err
	}

	return organizationservice.CreateOrganizationRequest{
		Name: req.Name,
	}, nil
}

func SerializeCreateOrganizationResponse(_ *gin.Context, organization *domain.Organization) any {
	return CreateOrganizationResponse{
		ID:   organization.ID.String(),
		Name: organization.Name,
	}
}

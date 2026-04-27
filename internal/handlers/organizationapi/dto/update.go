package dto

import (
	"media-equipment-tracker/internal/application/organizationservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateOrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DeserializeUpdateOrganizationRequest(c *gin.Context) (uuid.UUID, organizationservice.UpdateOrganizationRequest, error) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, organizationservice.UpdateOrganizationRequest{}, err
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return uuid.Nil, organizationservice.UpdateOrganizationRequest{}, err
	}

	return id, organizationservice.UpdateOrganizationRequest{
		Name: req.Name,
	}, nil
}

func SerializeUpdateOrganizationResponse(_ *gin.Context, organization *domain.Organization) any {
	return UpdateOrganizationResponse{
		ID:   organization.ID.String(),
		Name: organization.Name,
	}
}

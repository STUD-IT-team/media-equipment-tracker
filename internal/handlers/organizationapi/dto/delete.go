package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteOrganizationRequest struct {
	ID string `uri:"id" binding:"required"`
}

func DeserializeDeleteOrganizationRequest(c *gin.Context) (uuid.UUID, error) {
	var req DeleteOrganizationRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(req.ID)
}

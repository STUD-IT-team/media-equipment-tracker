package dto

import (
	"media-equipment-tracker/internal/application/studioservice"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateStudioStatusRequest struct {
	Status         string  `json:"status" binding:"required"`
	CuratorComment *string `json:"curator_comment,omitempty"`
}

func DeserializeUpdateStudioStatusRequest(c *gin.Context) (*studioservice.UpdateStudioInvocationStatusRequest, error) {
	var req UpdateStudioStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return nil, err
	}

	status, err := mapStudioStatus(req.Status)
	if err != nil {
		return nil, err
	}

	return &studioservice.UpdateStudioInvocationStatusRequest{
		ID:             id,
		Status:         status,
		CuratorComment: req.CuratorComment,
	}, nil
}

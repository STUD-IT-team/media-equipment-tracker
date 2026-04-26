package dto

import (
	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateEquipmentStatusRequest struct {
	Status         string  `json:"status" binding:"required"`
	CuratorComment *string `json:"curator_comment,omitempty"`
}

func DeserializeUpdateEquipmentStatusRequest(c *gin.Context) (*invocationservice.UpdateInvocationStatusRequest, error) {
	var req UpdateEquipmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return nil, err
	}

	status, err := mapInvocationStatus(req.Status)
	if err != nil {
		return nil, err
	}

	return &invocationservice.UpdateInvocationStatusRequest{
		ID:             id,
		Status:         status,
		CuratorComment: req.CuratorComment,
	}, nil
}

type UpdateInvocationStatusResponse struct {
	ShortInvocationItem
}

func SerializeUpdateInvocationStatusResponse(_ *gin.Context, inv *domain.EquipmentInvocation) any {
	return UpdateInvocationStatusResponse{
		ShortInvocationItem: ShortFromInvocation(inv),
	}
}

package dto

import (
	"time"

	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateInvocationRequest struct {
	EventName    *string  `json:"event_name" binding:"omitempty"`
	StartTime    *string  `json:"start_time" binding:"omitempty"`
	EndTime      *string  `json:"end_time" binding:"omitempty"`
	EquipmentIDs []string `json:"equipment_ids" binding:"required,min=1,dive"`
}

func DeserializeUpdateInvocationRequest(c *gin.Context) (*invocationservice.UpdateInvocationRequest, error) {
	var req UpdateInvocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return nil, err
	}

	var startTime *time.Time
	if req.StartTime != nil {
		t, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			return nil, err
		}
		startTime = &t
	}

	var endTime *time.Time
	if req.EndTime != nil {
		t, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return nil, err
		}
		endTime = &t
	}

	var equipmentIDs []uuid.UUID
	if req.EquipmentIDs != nil {
		equipmentIDs = make([]uuid.UUID, 0, len(req.EquipmentIDs))
		for _, id := range req.EquipmentIDs {
			equipmentID, err := uuid.Parse(id)
			if err != nil {
				return nil, err
			}
			equipmentIDs = append(equipmentIDs, equipmentID)
		}
	}

	return &invocationservice.UpdateInvocationRequest{
		ID:           id,
		EventName:    req.EventName,
		StartTime:    startTime,
		EndTime:      endTime,
		EquipmentIDs: equipmentIDs,
	}, nil
}

type UpdateInvocationResponse struct {
	ShortInvocationItem
}

func SerializeUpdateInvocationResponse(_ *gin.Context, inv *domain.EquipmentInvocation) any {
	return UpdateInvocationResponse{
		ShortInvocationItem: ShortFromInvocation(inv),
	}
}

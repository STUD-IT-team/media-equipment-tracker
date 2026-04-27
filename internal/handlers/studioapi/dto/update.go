package dto

import (
	"time"

	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateStudioRequest struct {
	EventName           *string `json:"event_name,omitempty"`
	ShootingDescription *string `json:"shooting_description,omitempty"`
	StartTime           *string `json:"start_time,omitempty"`
	EndTime             *string `json:"end_time,omitempty"`
	NeedsChromakey      *bool   `json:"needs_chromakey,omitempty"`
	NeedsCyclorama      *bool   `json:"needs_cyclorama,omitempty"`
	NeedsBlackFabric    *bool   `json:"needs_black_fabric,omitempty"`
}

func DeserializeUpdateStudioRequest(c *gin.Context) (*studioservice.UpdateStudioInvocationRequest, error) {
	var req UpdateStudioRequest
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

	return &studioservice.UpdateStudioInvocationRequest{
		ID:                  id,
		EventName:           req.EventName,
		ShootingDescription: req.ShootingDescription,
		StartTime:           startTime,
		EndTime:             endTime,
		NeedsChromakey:      req.NeedsChromakey,
		NeedsCyclorama:      req.NeedsCyclorama,
		NeedsBlackFabric:    req.NeedsBlackFabric,
	}, nil
}

type UpdateStudioResponse struct {
	ShortStudioItem
}

func SerializeUpdateStudioResponse(_ *gin.Context, inv *domain.StudioInvocation) any {
	return UpdateStudioResponse{
		ShortStudioItem: ShortFromStudio(inv),
	}
}

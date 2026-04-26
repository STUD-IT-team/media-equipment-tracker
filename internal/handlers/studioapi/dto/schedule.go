package dto

import (
	"fmt"
	"time"

	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

func DeserializeScheduleStudioRequest(c *gin.Context) (*studioservice.StudioScheduleRequest, error) {
	var startTime time.Time
	if val := c.Query("start_time"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		startTime = t
	} else {
		return nil, fmt.Errorf("start_time is required")
	}

	var endTime time.Time
	if val := c.Query("end_time"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		endTime = t
	} else {
		return nil, fmt.Errorf("end_time is required")
	}

	return &studioservice.StudioScheduleRequest{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

type ScheduleStudioResponse struct {
	Items []ShortStudioItem `json:"items"`
}

func SerializeScheduleStudioResponse(_ *gin.Context, items []*domain.StudioInvocation) any {
	response := ScheduleStudioResponse{
		Items: make([]ShortStudioItem, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, ShortFromStudio(item))
	}

	return response
}

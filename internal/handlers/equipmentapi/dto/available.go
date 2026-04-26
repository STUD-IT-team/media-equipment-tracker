package dto

import (
	"time"

	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeserializeAvailabilityEquipmentRequest(c *gin.Context) (equipmentservice.EquipmentAvailabilityRequest, error) {
	var id uuid.UUID
	var err error

	if val := c.Param("id"); val != "" {
		id, err = uuid.Parse(val)
		if err != nil {
			return equipmentservice.EquipmentAvailabilityRequest{}, err
		}
	}

	startTimeQuery := c.Query("start_time")
	endTimeQuery := c.Query("end_time")

	if startTimeQuery == "" || endTimeQuery == "" {
		return equipmentservice.EquipmentAvailabilityRequest{}, errs.NewValidationError("_time", "start_time and end_time are required")
	}

	startTime, err := time.Parse(time.RFC3339, startTimeQuery)
	if err != nil {
		return equipmentservice.EquipmentAvailabilityRequest{}, err
	}

	endTime, err := time.Parse(time.RFC3339, endTimeQuery)
	if err != nil {
		return equipmentservice.EquipmentAvailabilityRequest{}, err
	}
	return equipmentservice.EquipmentAvailabilityRequest{
		ID:        id,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

type AvailabilityEquipmentResponse struct {
	Available              bool                  `json:"available"`
	Equipment              ShortEquipmentItem    `json:"equipment"`
	ConflictingInvocations []ShortInvocationItem `json:"conflicting_invocations"`
}

type ShortInvocationItem struct {
	InvocationID string `json:"invocation_id"`
	EventName    string `json:"event_name"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Status       string `json:"status"`
}

func SerializeAvailabilityEquipmentResponse(_ *gin.Context, resp equipmentservice.EquipmentAvailabilityResponse) any {
	conflictingInvocations := make([]ShortInvocationItem, 0, len(resp.ConflictingInvocations))
	if resp.ConflictingInvocations != nil {
		for _, i := range resp.ConflictingInvocations {
			conflictingInvocations = append(conflictingInvocations, ShortInvocationItem{
				InvocationID: i.ID.String(),
				EventName:    i.EventName,
				StartTime:    i.StartTime.Format(time.RFC3339),
				EndTime:      i.EndTime.Format(time.RFC3339),
				Status:       string(i.Status),
			})
		}
	}
	return AvailabilityEquipmentResponse{
		Available:              resp.Available,
		Equipment:              FromEquipment(resp.Equipment),
		ConflictingInvocations: conflictingInvocations,
	}
}

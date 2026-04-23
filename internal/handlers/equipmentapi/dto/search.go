package dto

import (
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func mapEquipmentStatus(status string) (domain.EquipmentStatus, error) {
	switch status {
	case "available":
		return domain.EquipmentStatusAvailable, nil
	case "issued":
		return domain.EquipmentStatusIssued, nil
	case "under_maintenance":
		return domain.EquipmentStatusUnderMaintenance, nil
	case "unavailable":
		return domain.EquipmentStatusUnavailable, nil
	default:
		return "", nil
	}
}

func DeserializeSearchEquipmentRequest(c *gin.Context) (equipmentservice.SearchEquipmentRequest, error) {
	var req equipmentservice.SearchEquipmentRequest

	if val := c.Query("category"); val != "" {
		req.Categories = []string{val}
	}

	if val := c.Query("status"); val != "" {
		status, err := mapEquipmentStatus(val)
		if err != nil {
			return req, err
		}
		req.Statuses = []domain.EquipmentStatus{status}
	}

	if val := c.Query("available_to_trainee"); val != "" {
		var a bool = (val == "true")
		req.AvailableToTrainee = &a
	}

	if val := c.Query("department_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return req, err
		}
		req.DepartmentIDs = []uuid.UUID{id}
	}

	if val := c.Query("available_at"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return req, err
		}
		req.AvailableAt = &t
	}

	if val := c.Query("search"); val != "" {
		req.SearchString = &val
	}

	return req, nil
}

type SearchEquipmentResponse struct {
	Items []ShortEquipmentItem `json:"items"`
}

func SerializeSearchEquipmentResponse(c *gin.Context, items []*domain.Equipment) any {
	response := SearchEquipmentResponse{
		Items: make([]ShortEquipmentItem, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, FromEquipment(item))
	}

	return response
}

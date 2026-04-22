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
	Items []SearchEquipmentItem
}

type SearchEquipmentItem struct {
	ID                 string `json:"id"`
	InventoryNumber    string `json:"inventory_number"`
	Name               string `json:"name"`
	ShortName          string `json:"short_name"`
	Category           string `json:"category"`
	AvailableToTrainee bool   `json:"available_to_trainee"`
	Status             string `json:"status"`
}

func SerializeSearchEquipmentResponse(c *gin.Context, items []*domain.Equipment) any {
	response := SearchEquipmentResponse{
		Items: make([]SearchEquipmentItem, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, SearchEquipmentItem{
			ID:                 item.ID.String(),
			InventoryNumber:    item.InventoryNumber,
			Name:               item.Name,
			ShortName:          item.ShortName,
			Category:           item.Category,
			AvailableToTrainee: item.AvailableToTrainee,
			Status:             string(item.Status),
		})
	}

	return response
}

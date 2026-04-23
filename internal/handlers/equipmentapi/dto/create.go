package dto

import (
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateEquipmentRequest struct {
	InventoryNumber    string   `json:"inventory_number" binding:"required"`
	Name               string   `json:"name" binding:"required"`
	ShortName          string   `json:"short_name" binding:"required"`
	Category           string   `json:"category" binding:"required"`
	AvailableToTrainee bool     `json:"available_to_trainee" binding:"required"`
	Status             string   `json:"status" binding:"required"`
	Departments        []string `json:"department_ids,omitempty"`
}

type CreateEquipmentResponse struct {
	ShortEquipmentItem
}

func DeserializeCreateEquipmentRequest(c *gin.Context) (equipmentservice.CreateEquipmentRequest, error) {
	var req CreateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return equipmentservice.CreateEquipmentRequest{}, err
	}

	status, err := mapEquipmentStatus(req.Status)
	if err != nil {
		return equipmentservice.CreateEquipmentRequest{}, err
	}

	depIDs := make([]uuid.UUID, 0, len(req.Departments))
	for _, depID := range req.Departments {
		id, err := uuid.Parse(depID)
		if err != nil {
			return equipmentservice.CreateEquipmentRequest{}, err
		}
		depIDs = append(depIDs, id)
	}

	return equipmentservice.CreateEquipmentRequest{
		InventoryNumber:    req.InventoryNumber,
		Name:               req.Name,
		ShortName:          req.ShortName,
		Category:           req.Category,
		AvailableToTrainee: req.AvailableToTrainee,
		Status:             status,
		Departments:        depIDs,
	}, nil
}

func SerializeCreateEquipmentResponse(c *gin.Context, equipment *domain.Equipment) any {
	return CreateEquipmentResponse{
		ShortEquipmentItem: FromEquipment(equipment),
	}
}

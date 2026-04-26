package dto

import (
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateEquipmentRequestBody struct {
	InventoryNumber    *string  `json:"inventory_number,omitempty"`
	Name               *string  `json:"name,omitempty"`
	ShortName          *string  `json:"short_name,omitempty"`
	Category           *string  `json:"category,omitempty"`
	AvailableToTrainee *bool    `json:"available_to_trainee,omitempty"`
	Status             *string  `json:"status,omitempty"`
	Departments        []string `json:"department_ids,omitempty"`
}

func DeserializeUpdateEquipmentRequest(c *gin.Context) (equipmentservice.UpdateEquipmentRequest, error) {
	var body UpdateEquipmentRequestBody
	var id uuid.UUID
	var err error

	if val := c.Param("id"); val != "" {
		id, err = uuid.Parse(val)
		if err != nil {
			return equipmentservice.UpdateEquipmentRequest{}, err
		}
	}

	err = c.BindJSON(&body)
	if err != nil {
		return equipmentservice.UpdateEquipmentRequest{}, err
	}
	var status *domain.EquipmentStatus
	if body.Status != nil {
		s, err := mapEquipmentStatus(*body.Status)
		if err != nil {
			return equipmentservice.UpdateEquipmentRequest{}, err
		}
		status = &s
	}

	departmentIDs := make([]uuid.UUID, 0, len(body.Departments))
	for _, departmentID := range body.Departments {
		departmentID, err := uuid.Parse(departmentID)
		if err != nil {
			return equipmentservice.UpdateEquipmentRequest{}, err
		}
		departmentIDs = append(departmentIDs, departmentID)
	}

	return equipmentservice.UpdateEquipmentRequest{
		ID:                 id,
		InventoryNumber:    body.InventoryNumber,
		Name:               body.Name,
		ShortName:          body.ShortName,
		Category:           body.Category,
		AvailableToTrainee: body.AvailableToTrainee,
		Status:             status,
		Departments:        departmentIDs,
	}, nil
}

type UpdateEquipmentResponse struct {
	ShortEquipmentItem
}

func SerializeUpdateEquipmentResponse(_ *gin.Context, equipment *domain.Equipment) any {
	return UpdateEquipmentResponse{
		ShortEquipmentItem: FromEquipment(equipment),
	}
}

package dto

import (
	"fmt"
	"time"

	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeserializeGetEquipmentRequest(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func DeserializeGetEquipmentByInventoryNumberRequest(c *gin.Context) (string, error) {
	inv := c.Param("inv")
	if inv == "" {
		return "", fmt.Errorf("inventory number is required")
	}
	return inv, nil
}

type GetEquipmentResponse struct {
	ID                 uuid.UUID                `json:"id"`
	InventoryNumber    string                   `json:"inventory_number"`
	Name               string                   `json:"name"`
	ShortName          string                   `json:"short_name"`
	Category           string                   `json:"category"`
	AvailableToTrainee bool                     `json:"available_to_trainee"`
	Status             string                   `json:"status"`
	Departments        []getEquipmentDepartment `json:"departments"`
	CurrentInvocation  *getEquipmentInvocation  `json:"current_invocation,omitempty"`
}

type getEquipmentDepartment struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type getEquipmentInvocation struct {
	ID        uuid.UUID `json:"id"`
	EventName string    `json:"event_name"`
	UserID    uuid.UUID `json:"user_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
}

func SerializeGetEquipmentResponse(_ *gin.Context, equipment *domain.Equipment) any {
	departments := make([]getEquipmentDepartment, 0, len(equipment.Departments))
	if equipment.Departments != nil {
		for _, d := range equipment.Departments {
			departments = append(departments, getEquipmentDepartment{
				ID:   d.ID,
				Name: d.Name,
			})
		}
	}
	var currentInvocation *getEquipmentInvocation
	if equipment.CurrentInvocation != nil {
		currentInvocation = &getEquipmentInvocation{
			ID:        equipment.CurrentInvocation.ID,
			EventName: equipment.CurrentInvocation.EventName,
			UserID:    equipment.CurrentInvocation.UserID,
			StartTime: equipment.CurrentInvocation.StartTime.Format(time.RFC3339),
			EndTime:   equipment.CurrentInvocation.EndTime.Format(time.RFC3339),
		}
	}

	return GetEquipmentResponse{
		ID:                 equipment.ID,
		InventoryNumber:    equipment.InventoryNumber,
		Name:               equipment.Name,
		ShortName:          equipment.ShortName,
		Category:           equipment.Category,
		AvailableToTrainee: equipment.AvailableToTrainee,
		Status:             string(equipment.Status),
		Departments:        departments,
		CurrentInvocation:  currentInvocation,
	}
}

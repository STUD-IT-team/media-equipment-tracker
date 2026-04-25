package dto

import "media-equipment-tracker/internal/domain"

type ShortEquipmentItem struct {
	ID                 string `json:"id"`
	InventoryNumber    string `json:"inventory_number"`
	Name               string `json:"name"`
	ShortName          string `json:"short_name"`
	Category           string `json:"category"`
	AvailableToTrainee bool   `json:"available_to_trainee"`
	Status             string `json:"status"`
}

func FromEquipment(equipment *domain.Equipment) ShortEquipmentItem {
	return ShortEquipmentItem{
		ID:                 equipment.ID.String(),
		InventoryNumber:    equipment.InventoryNumber,
		Name:               equipment.Name,
		ShortName:          equipment.ShortName,
		Category:           equipment.Category,
		AvailableToTrainee: equipment.AvailableToTrainee,
		Status:             string(equipment.Status),
	}
}

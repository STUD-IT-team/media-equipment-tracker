package dto

import (
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetDepartmentEquipmentRequest struct {
	ID            string `uri:"id" binding:"required"`
	AvailableOnly bool   `form:"available_only"`
}

type GetDepartmentEquipmentResponse struct {
	Items []EquipmentResponse `json:"items"`
}

func DeserializeGetDepartmentEquipmentRequest(c *gin.Context) (uuid.UUID, bool, error) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, false, err
	}

	var req GetDepartmentEquipmentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return uuid.Nil, false, err
	}

	return id, req.AvailableOnly, nil
}

func SerializeGetDepartmentEquipmentResponse(_ *gin.Context, equipment []*domain.Equipment) any {
	items := make([]EquipmentResponse, 0, len(equipment))
	for _, eq := range equipment {
		items = append(items, EquipmentResponse{
			ID:              eq.ID.String(),
			Name:            eq.Name,
			ShortName:       eq.ShortName,
			InventoryNumber: eq.InventoryNumber,
			Category:        eq.Category,
			Status:          string(eq.Status),
		})
	}
	return GetDepartmentEquipmentResponse{
		Items: items,
	}
}

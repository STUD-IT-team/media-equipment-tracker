package pgequipmentininvocation_test

import (
	"media-equipment-tracker/internal/domain"
	"time"

	"github.com/google/uuid"
)

func newEquipment() *domain.Equipment {
	return &domain.Equipment{
		ID:              uuid.New(),
		InventoryNumber: uuid.NewString(),
		Name:            "Camera",
		ShortName:       "cam",
		Category:        "video",
		Status:          domain.EquipmentStatusAvailable,
	}
}

func newInvocation(depID uuid.UUID, userID uuid.UUID) *domain.EquipmentInvocation {
	return &domain.EquipmentInvocation{
		ID:           uuid.New(),
		EventName:    "event",
		StartTime:    time.Now().Add(-time.Hour * 96),
		EndTime:      time.Now().Add(time.Hour * 96),
		Status:       domain.InvocationCreated,
		DepartmentID: &depID,
		UserID:       userID,
	}
}

func newEqInInv(invID, eqID uuid.UUID) *domain.EquipmentInInvocation {
	return &domain.EquipmentInInvocation{
		InvocationID: invID,
		EquipmentID:  eqID,
		Status:       domain.EquipmentNotIssued,
	}
}

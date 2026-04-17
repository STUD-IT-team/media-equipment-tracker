//go:build integration

package pgequipment_test

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

func newDepartment() *domain.Department {
	return &domain.Department{
		ID:   uuid.New(),
		Name: "dept",
	}
}

func newUser() *domain.User {
	return &domain.User{
		ID:           uuid.New(),
		FullName:     "test",
		Email:        uuid.NewString() + "@mail.ru",
		HashPassword: "hash",
		Nice:         domain.DefaultNice,
	}
}

func newInvocation(depID uuid.UUID, userID uuid.UUID) *domain.EquipmentInvocation {
	return &domain.EquipmentInvocation{
		ID:                  uuid.New(),
		EventName:           "event",
		StartTime:           time.Now().Add(-time.Hour * 96),
		EndTime:             time.Now().Add(time.Hour * 96),
		EquipmentReturnTime: nil,
		SdCardReturnTime:    nil,
		Status:              domain.InvocationCreated,
		DepartmentID:        &depID,
		UserID:              userID,
	}
}

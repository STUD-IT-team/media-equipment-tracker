//go:build integration

package pgstudioinvocation_test

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
func newStudioInvocation(orgID *uuid.UUID, depID *uuid.UUID, userID *uuid.UUID, adminID *uuid.UUID) *domain.StudioInvocation {
	now := time.Now()

	return &domain.StudioInvocation{
		ID:                  uuid.New(),
		EventName:           "Event",
		ShootingDescription: "Desc",
		StartTime:           now,
		EndTime:             now.Add(time.Hour),
		Status:              domain.StudioCreated,

		OrganizationID: orgID,
		DepartmentID:   depID,
		AdminID:        adminID,
		UserID:         *userID,
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

func newOrg() *domain.Organization {
	return &domain.Organization{
		ID:   uuid.New(),
		Name: "org",
	}
}

func newDept() *domain.Department {
	return &domain.Department{
		ID:   uuid.New(),
		Name: "dept",
	}
}

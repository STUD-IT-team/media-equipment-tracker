package pgmessage_test

import (
	"media-equipment-tracker/internal/domain"
	"time"

	"github.com/google/uuid"
)

func newUser() *domain.User {
	return &domain.User{
		ID:           uuid.New(),
		FullName:     "test",
		Email:        uuid.NewString() + "@mail.ru",
		HashPassword: "hash",
		Nice:         domain.DefaultNice,
	}
}

func newDepartment() *domain.Department {
	return &domain.Department{
		ID:   uuid.New(),
		Name: "dept",
	}
}

func newEquipmentInvocation(depID uuid.UUID, userID uuid.UUID) *domain.EquipmentInvocation {
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

func newStudioInvocation(depID uuid.UUID, userID uuid.UUID) *domain.StudioInvocation {
	now := time.Now()

	return &domain.StudioInvocation{
		ID:                  uuid.New(),
		EventName:           "Event",
		ShootingDescription: "Desc",
		StartTime:           now,
		EndTime:             now.Add(time.Hour),
		Status:              domain.StudioCreated,

		DepartmentID: &depID,
		UserID:       userID,
	}
}

func newMessageEquipmentInvocation(invID, senderID, recipientID uuid.UUID) *domain.MessageEquipmentInvocation {
	return &domain.MessageEquipmentInvocation{
		ID:           uuid.New(),
		Content:      "test message",
		SentTime:     time.Now(),
		InvocationID: invID,
		SenderID:     senderID,
		RecipientID:  recipientID,
	}
}

func newMessageStudioInvocation(invID, senderID, recipientID uuid.UUID) *domain.MessageStudioInvocation {
	return &domain.MessageStudioInvocation{
		ID:           uuid.New(),
		Content:      "test message",
		SentTime:     time.Now(),
		InvocationID: invID,
		SenderID:     senderID,
		RecipientID:  recipientID,
	}
}

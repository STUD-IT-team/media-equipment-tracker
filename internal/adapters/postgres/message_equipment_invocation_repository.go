package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
)

type messageEquipmentInvocationRepository struct {
	db *gorm.DB
}

func NewMessageEquipmentInvocationRepository(db *gorm.DB) domain.MessageEquipmentInvocationRepository {
	return &messageEquipmentInvocationRepository{db: db}
}

func (r *messageEquipmentInvocationRepository) Create(ctx context.Context, message *domain.MessageEquipmentInvocation) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageEquipmentInvocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MessageEquipmentInvocation, error) {
	var message domain.MessageEquipmentInvocation
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageEquipmentInvocationRepository) GetAll(ctx context.Context) ([]*domain.MessageEquipmentInvocation, error) {
	var messages []*domain.MessageEquipmentInvocation
	err := r.db.WithContext(ctx).Find(&messages).Error
	return messages, err
}

func (r *messageEquipmentInvocationRepository) Update(ctx context.Context, message *domain.MessageEquipmentInvocation) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *messageEquipmentInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.MessageEquipmentInvocation{}, id).Error
}

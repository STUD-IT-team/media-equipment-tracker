package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
)

type messageStudioInvocationRepository struct {
	db *gorm.DB
}

func NewMessageStudioInvocationRepository(db *gorm.DB) domain.MessageStudioInvocationRepository {
	return &messageStudioInvocationRepository{db: db}
}

func (r *messageStudioInvocationRepository) Create(ctx context.Context, message *domain.MessageStudioInvocation) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageStudioInvocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MessageStudioInvocation, error) {
	var message domain.MessageStudioInvocation
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageStudioInvocationRepository) GetAll(ctx context.Context) ([]*domain.MessageStudioInvocation, error) {
	var messages []*domain.MessageStudioInvocation
	err := r.db.WithContext(ctx).Find(&messages).Error
	return messages, err
}

func (r *messageStudioInvocationRepository) Update(ctx context.Context, message *domain.MessageStudioInvocation) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *messageStudioInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.MessageStudioInvocation{}, id).Error
}

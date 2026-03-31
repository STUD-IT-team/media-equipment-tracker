package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"
)

type MessageEquipmentInvocationRepository struct {
	db *gorm.DB
}

func NewMessageEquipmentInvocationRepository(db *gorm.DB) domain.MessageEquipmentInvocationRepository {
	return &MessageEquipmentInvocationRepository{db: db}
}

func (r *MessageEquipmentInvocationRepository) applyOptions(opts []domain.MessageEquipmentInvocationOption) *gorm.DB {
	options := &domain.MessageEquipmentInvocationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *MessageEquipmentInvocationRepository) Get(id uuid.UUID, with ...domain.MessageEquipmentInvocationOption) (*domain.MessageEquipmentInvocation, error) {
	var msg domain.MessageEquipmentInvocation
	query := r.applyOptions(with)
	if err := query.First(&msg, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewEntityNotFoundError("MessageEquipmentInvocation", id)
		}
		return nil, errors.NewRepositoryError("get", err)
	}
	return &msg, nil
}

func (r *MessageEquipmentInvocationRepository) GetInvocation(invocationID uuid.UUID, with ...domain.MessageEquipmentInvocationOption) ([]*domain.MessageEquipmentInvocation, error) {
	var messages []*domain.MessageEquipmentInvocation
	query := r.applyOptions(with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&messages).Error; err != nil {
		return nil, errors.NewRepositoryError("get_invocation", err)
	}
	return messages, nil
}

func (r *MessageEquipmentInvocationRepository) Create(messageEquipmentInvocation *domain.MessageEquipmentInvocation) error {
	if err := r.db.Create(messageEquipmentInvocation).Error; err != nil {
		return errors.NewRepositoryError("create", err)
	}
	return nil
}

func (r *MessageEquipmentInvocationRepository) Update(messageEquipmentInvocation *domain.MessageEquipmentInvocation) error {
	if err := r.db.Save(messageEquipmentInvocation).Error; err != nil {
		return errors.NewRepositoryError("update", err)
	}
	return nil
}

func (r *MessageEquipmentInvocationRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.MessageEquipmentInvocation{}, "id = ?", id).Error; err != nil {
		return errors.NewRepositoryError("delete", err)
	}
	return nil
}

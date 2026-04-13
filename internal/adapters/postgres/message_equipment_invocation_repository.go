package postgresrepo

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
			return nil, errs.NewEntityNotFoundError("MessageEquipmentInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &msg, nil
}

func (r *MessageEquipmentInvocationRepository) GetInvocation(invocationID uuid.UUID, with ...domain.MessageEquipmentInvocationOption) ([]*domain.MessageEquipmentInvocation, error) {
	var messages []*domain.MessageEquipmentInvocation
	query := r.applyOptions(with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&messages).Error; err != nil {
		return nil, errs.NewRepositoryError("get_invocation", err)
	}
	return messages, nil
}

func (r *MessageEquipmentInvocationRepository) Create(messageEquipmentInvocation *domain.MessageEquipmentInvocation) error {
	if err := r.db.Create(messageEquipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *MessageEquipmentInvocationRepository) Update(messageEquipmentInvocation *domain.MessageEquipmentInvocation) error {
	if err := r.db.Save(messageEquipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *MessageEquipmentInvocationRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.MessageEquipmentInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

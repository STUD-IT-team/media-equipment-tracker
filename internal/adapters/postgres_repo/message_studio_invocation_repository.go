package postgresrepo

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageStudioInvocationRepository struct {
	db *gorm.DB
}

func NewMessageStudioInvocationRepository(db *gorm.DB) domain.MessageStudioInvocationRepository {
	return &MessageStudioInvocationRepository{db: db}
}

func (r *MessageStudioInvocationRepository) applyOptions(opts []domain.MessageStudioInvocationOption) *gorm.DB {
	options := &domain.MessageStudioInvocationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *MessageStudioInvocationRepository) Get(id uuid.UUID, with ...domain.MessageStudioInvocationOption) (*domain.MessageStudioInvocation, error) {
	var msg domain.MessageStudioInvocation
	query := r.applyOptions(with)
	if err := query.First(&msg, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("MessageStudioInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &msg, nil
}

func (r *MessageStudioInvocationRepository) GetInvocation(invocationID uuid.UUID, with ...domain.MessageStudioInvocationOption) ([]*domain.MessageStudioInvocation, error) {
	var messages []*domain.MessageStudioInvocation
	query := r.applyOptions(with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&messages).Error; err != nil {
		return nil, errs.NewRepositoryError("get_invocation", err)
	}
	return messages, nil
}

func (r *MessageStudioInvocationRepository) Create(messageStudioInvocation *domain.MessageStudioInvocation) error {
	if err := r.db.Create(messageStudioInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *MessageStudioInvocationRepository) Update(messageStudioInvocation *domain.MessageStudioInvocation) error {
	if err := r.db.Save(messageStudioInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *MessageStudioInvocationRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.MessageStudioInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

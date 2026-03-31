package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"
)

type studioInvocationRepository struct {
	db *gorm.DB
}

func NewStudioInvocationRepository(db *gorm.DB) domain.StudioInvocationRepository {
	return &studioInvocationRepository{db: db}
}

func (r *studioInvocationRepository) Create(ctx context.Context, studioInvocation *domain.StudioInvocation) error {
	return r.db.WithContext(ctx).Create(studioInvocation).Error
}

func (r *studioInvocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.StudioInvocation, error) {
	var studioInvocation domain.StudioInvocation
	err := r.db.WithContext(ctx).First(&studioInvocation, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &studioInvocation, nil
}

func (r *studioInvocationRepository) GetAll(ctx context.Context) ([]*domain.StudioInvocation, error) {
	var studioInvocations []*domain.StudioInvocation
	err := r.db.WithContext(ctx).Find(&studioInvocations).Error
	return studioInvocations, err
}

func (r *studioInvocationRepository) Update(ctx context.Context, studioInvocation *domain.StudioInvocation) error {
	return r.db.WithContext(ctx).Save(studioInvocation).Error
}

func (r *studioInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.StudioInvocation{}, id).Error
}

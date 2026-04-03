package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"
)

type StudioInvocationRepository struct {
	db *gorm.DB
}

func NewStudioInvocationRepository(db *gorm.DB) domain.StudioInvocationRepository {
	return &StudioInvocationRepository{db: db}
}

func (r *StudioInvocationRepository) applyOptions(opts []domain.StudioInvocationOption) *gorm.DB {
	options := &domain.StudioInvocationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *StudioInvocationRepository) Get(id uuid.UUID, with ...domain.StudioInvocationOption) (*domain.StudioInvocation, error) {
	var inv domain.StudioInvocation
	query := r.applyOptions(with)
	if err := query.First(&inv, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewEntityNotFoundError("StudioInvocation", id)
		}
		return nil, errors.NewRepositoryError("get", err)
	}
	return &inv, nil
}

func (r *StudioInvocationRepository) List(with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	var invocations []*domain.StudioInvocation
	query := r.applyOptions(with)
	if err := query.Find(&invocations).Error; err != nil {
		return nil, errors.NewRepositoryError("list", err)
	}
	return invocations, nil
}

func (r *StudioInvocationRepository) Reload(studioInvocation *domain.StudioInvocation, with ...domain.StudioInvocationOption) error {
	query := r.applyOptions(with)
	if err := query.First(studioInvocation, "id = ?", studioInvocation.ID).Error; err != nil {
		return errors.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *StudioInvocationRepository) Create(studioInvocation *domain.StudioInvocation) error {
	if err := r.db.Create(studioInvocation).Error; err != nil {
		return errors.NewRepositoryError("create", err)
	}
	return nil
}

func (r *StudioInvocationRepository) Update(studioInvocation *domain.StudioInvocation) error {
	if err := r.db.Save(studioInvocation).Error; err != nil {
		return errors.NewRepositoryError("update", err)
	}
	return nil
}

func (r *StudioInvocationRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.StudioInvocation{}, "id = ?", id).Error; err != nil {
		return errors.NewRepositoryError("delete", err)
	}
	return nil
}

package postgresstudioinvocation

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresStudioInvocationRepository struct {
	db *gorm.DB
}

func NewPostgresStudioInvocationRepository(db *gorm.DB) *PostgresStudioInvocationRepository {
	return &PostgresStudioInvocationRepository{db: db}
}

var _ domain.StudioInvocationRepository = (*PostgresStudioInvocationRepository)(nil)

func (r *PostgresStudioInvocationRepository) applyOptions(opts []domain.StudioInvocationOption) *gorm.DB {
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

func (r *PostgresStudioInvocationRepository) Get(id uuid.UUID, with ...domain.StudioInvocationOption) (*domain.StudioInvocation, error) {
	var inv domain.StudioInvocation
	query := r.applyOptions(with)
	if err := query.First(&inv, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("StudioInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &inv, nil
}

func (r *PostgresStudioInvocationRepository) List(with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	var invocations []*domain.StudioInvocation
	query := r.applyOptions(with)
	if err := query.Find(&invocations).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return invocations, nil
}

func (r *PostgresStudioInvocationRepository) Reload(studioInvocation *domain.StudioInvocation, with ...domain.StudioInvocationOption) error {
	query := r.applyOptions(with)
	if err := query.First(studioInvocation, "id = ?", studioInvocation.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresStudioInvocationRepository) Create(studioInvocation *domain.StudioInvocation) error {
	if err := r.db.Create(studioInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresStudioInvocationRepository) Update(studioInvocation *domain.StudioInvocation) error {
	if err := r.db.Save(studioInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresStudioInvocationRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.StudioInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

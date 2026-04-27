package pgstudioinvocation

import (
	"context"
	"errors"

	"media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresStudioInvocationRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresStudioInvocationRepository(db *gormtx.DBGetter) *PostgresStudioInvocationRepository {
	return &PostgresStudioInvocationRepository{db: db}
}

var _ domain.StudioInvocationRepository = (*PostgresStudioInvocationRepository)(nil)

func (r *PostgresStudioInvocationRepository) applyOptions(ctx context.Context, opts []domain.StudioInvocationOption) *gorm.DB {
	options := &domain.StudioInvocationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return nil
	}
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *PostgresStudioInvocationRepository) Get(ctx context.Context, id uuid.UUID, with ...domain.StudioInvocationOption) (*domain.StudioInvocation, error) {
	var inv domain.StudioInvocation
	query := r.applyOptions(ctx, with)
	if err := query.First(&inv, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewEntityNotFoundError("StudioInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &inv, nil
}

func (r *PostgresStudioInvocationRepository) List(ctx context.Context, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	var invocations []*domain.StudioInvocation
	query := r.applyOptions(ctx, with)
	if err := query.Find(&invocations).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return invocations, nil
}

func (r *PostgresStudioInvocationRepository) Reload(ctx context.Context, studioInvocation *domain.StudioInvocation, with ...domain.StudioInvocationOption) error {
	query := r.applyOptions(ctx, with)
	if err := query.First(studioInvocation, "id = ?", studioInvocation.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresStudioInvocationRepository) Create(ctx context.Context, studioInvocation *domain.StudioInvocation) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("create", err)
	}
	if err := db.Omit(clause.Associations).Create(studioInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresStudioInvocationRepository) Update(ctx context.Context, studioInvocation *domain.StudioInvocation) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("update", err)
	}
	if err := db.Omit(clause.Associations).Save(studioInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresStudioInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	if err := db.Delete(&domain.StudioInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

func (r *PostgresStudioInvocationRepository) Search(ctx context.Context, search *studiosearch.SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	var invocations []*domain.StudioInvocation
	query := r.applyOptions(ctx, with)
	if search.SearchString != nil {
		query = query.Where("event_name like ?", "%"+*search.SearchString+"%").Or("shooting_description like ?", "%"+*search.SearchString+"%")
	}

	// Если хотя бы кусочек события в промежутке, то подходит
	if search.StartTime != nil {
		query = query.Where("end_time >= ?", *search.StartTime)
	}
	if search.EndTime != nil {
		query = query.Where("start_time <= ?", *search.EndTime)
	}

	if search.AdminID != nil {
		query = query.Where("admin_id = ?", *search.AdminID)
	}

	if search.UserID != nil {
		query = query.Where("user_id = ?", *search.UserID)
	}

	if search.DepartmentID != nil {
		query = query.Where("department_id = ?", *search.DepartmentID)
	}

	if search.OrganizationID != nil {
		query = query.Where("organization_id = ?", *search.OrganizationID)
	}

	if err := query.Find(&invocations).Error; err != nil {
		return nil, errs.NewRepositoryError("search", err)
	}

	return invocations, nil
}

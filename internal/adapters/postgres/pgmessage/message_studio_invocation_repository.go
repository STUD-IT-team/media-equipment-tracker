package pgmessage

import (
	"context"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresMessageStudioInvocationRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresMessageStudioInvocationRepository(db *gormtx.DBGetter) *PostgresMessageStudioInvocationRepository {
	return &PostgresMessageStudioInvocationRepository{db: db}
}

var _ domain.MessageStudioInvocationRepository = (*PostgresMessageStudioInvocationRepository)(nil)

func (r *PostgresMessageStudioInvocationRepository) applyOptions(ctx context.Context, opts []domain.MessageStudioInvocationOption) *gorm.DB {
	options := &domain.MessageStudioInvocationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	db, _ := r.db.GetDB(ctx)
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *PostgresMessageStudioInvocationRepository) Get(ctx context.Context, id uuid.UUID, with ...domain.MessageStudioInvocationOption) (*domain.MessageStudioInvocation, error) {
	var msg domain.MessageStudioInvocation
	query := r.applyOptions(ctx, with)
	if err := query.First(&msg, "id = ?", id).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("MessageStudioInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &msg, nil
}

func (r *PostgresMessageStudioInvocationRepository) GetInvocation(ctx context.Context, invocationID uuid.UUID, with ...domain.MessageStudioInvocationOption) ([]*domain.MessageStudioInvocation, error) {
	var messages []*domain.MessageStudioInvocation
	query := r.applyOptions(ctx, with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&messages).Error; err != nil {
		return nil, errs.NewRepositoryError("get_invocation", err)
	}
	return messages, nil
}

func (r *PostgresMessageStudioInvocationRepository) Create(ctx context.Context, messageStudioInvocation *domain.MessageStudioInvocation) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Create(messageStudioInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresMessageStudioInvocationRepository) Update(ctx context.Context, messageStudioInvocation *domain.MessageStudioInvocation) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Save(messageStudioInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresMessageStudioInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Delete(&domain.MessageStudioInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

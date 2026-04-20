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

type PostgresMessageEquipmentInvocationRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresMessageEquipmentInvocationRepository(db *gormtx.DBGetter) *PostgresMessageEquipmentInvocationRepository {
	return &PostgresMessageEquipmentInvocationRepository{db: db}
}

var _ domain.MessageEquipmentInvocationRepository = (*PostgresMessageEquipmentInvocationRepository)(nil)

func (r *PostgresMessageEquipmentInvocationRepository) applyOptions(ctx context.Context, opts []domain.MessageEquipmentInvocationOption) *gorm.DB {
	options := &domain.MessageEquipmentInvocationOptions{}
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

func (r *PostgresMessageEquipmentInvocationRepository) Get(ctx context.Context, id uuid.UUID, with ...domain.MessageEquipmentInvocationOption) (*domain.MessageEquipmentInvocation, error) {
	var msg domain.MessageEquipmentInvocation
	query := r.applyOptions(ctx, with)
	if err := query.First(&msg, "id = ?", id).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("MessageEquipmentInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &msg, nil
}

func (r *PostgresMessageEquipmentInvocationRepository) GetInvocation(ctx context.Context, invocationID uuid.UUID, with ...domain.MessageEquipmentInvocationOption) ([]*domain.MessageEquipmentInvocation, error) {
	var messages []*domain.MessageEquipmentInvocation
	query := r.applyOptions(ctx, with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&messages).Error; err != nil {
		return nil, errs.NewRepositoryError("get_invocation", err)
	}
	return messages, nil
}

func (r *PostgresMessageEquipmentInvocationRepository) Create(ctx context.Context, messageEquipmentInvocation *domain.MessageEquipmentInvocation) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Create(messageEquipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresMessageEquipmentInvocationRepository) Update(ctx context.Context, messageEquipmentInvocation *domain.MessageEquipmentInvocation) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Save(messageEquipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresMessageEquipmentInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Delete(&domain.MessageEquipmentInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

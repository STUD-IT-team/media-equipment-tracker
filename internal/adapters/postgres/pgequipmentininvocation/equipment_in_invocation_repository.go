package pgequipmentininvocation

import (
	"context"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresEquipmentInInvocationRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresEquipmentInInvocationRepository(db *gormtx.DBGetter) *PostgresEquipmentInInvocationRepository {
	return &PostgresEquipmentInInvocationRepository{db: db}
}

var _ domain.EquipmentInInvocationRepository = (*PostgresEquipmentInInvocationRepository)(nil)

func (r *PostgresEquipmentInInvocationRepository) applyOptions(ctx context.Context, opts []domain.EquipmentInInvocationOption) *gorm.DB {
	options := &domain.EquipmentInInvocationOptions{}
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

func (r *PostgresEquipmentInInvocationRepository) Get(ctx context.Context, invocationID, equipmentID uuid.UUID, with ...domain.EquipmentInInvocationOption) (*domain.EquipmentInInvocation, error) {
	var eii domain.EquipmentInInvocation
	query := r.applyOptions(ctx, with)
	if err := query.First(&eii, "invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("EquipmentInInvocation", invocationID)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &eii, nil
}

func (r *PostgresEquipmentInInvocationRepository) GetByInvocation(ctx context.Context, invocationID uuid.UUID, with ...domain.EquipmentInInvocationOption) ([]*domain.EquipmentInInvocation, error) {
	var items []*domain.EquipmentInInvocation
	query := r.applyOptions(ctx, with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&items).Error; err != nil {
		return nil, errs.NewRepositoryError("get_by_invocation", err)
	}
	return items, nil
}

func (r *PostgresEquipmentInInvocationRepository) GetByEquipment(ctx context.Context, equipmentID uuid.UUID, with ...domain.EquipmentInInvocationOption) ([]*domain.EquipmentInInvocation, error) {
	var items []*domain.EquipmentInInvocation
	query := r.applyOptions(ctx, with)
	if err := query.Where("equipment_id = ?", equipmentID).Find(&items).Error; err != nil {
		return nil, errs.NewRepositoryError("get_by_equipment", err)
	}
	return items, nil
}

func (r *PostgresEquipmentInInvocationRepository) Create(ctx context.Context, equipmentInInvocation *domain.EquipmentInInvocation) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Create(equipmentInInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresEquipmentInInvocationRepository) Update(ctx context.Context, equipmentInInvocation *domain.EquipmentInInvocation) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Save(equipmentInInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresEquipmentInInvocationRepository) Delete(ctx context.Context, invocationID, equipmentID uuid.UUID) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Delete(&domain.EquipmentInInvocation{}, "invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

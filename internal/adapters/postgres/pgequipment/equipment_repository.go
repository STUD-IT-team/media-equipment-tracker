package pgequipment

import (
	"context"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresEquipmentRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresEquipmentRepository(db *gormtx.DBGetter) *PostgresEquipmentRepository {
	return &PostgresEquipmentRepository{db: db}
}

var _ domain.EquipmentRepository = (*PostgresEquipmentRepository)(nil)

func (r *PostgresEquipmentRepository) applyOptions(ctx context.Context, opts []domain.EquipmentOption) *gorm.DB {
	options := &domain.EquipmentOptions{}
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

func (r *PostgresEquipmentRepository) Get(ctx context.Context, id uuid.UUID, with ...domain.EquipmentOption) (*domain.Equipment, error) {
	var eq domain.Equipment
	query := r.applyOptions(ctx, with)
	if err := query.First(&eq, "id = ?", id).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("Equipment", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &eq, nil
}

func (r *PostgresEquipmentRepository) GetUnoccupied(ctx context.Context, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(ctx, with)
	if err := query.Where("current_invocation_id IS NULL").Find(&equipment).Error; err != nil {
		return nil, errs.NewRepositoryError("get_unoccupied", err)
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) List(ctx context.Context, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(ctx, with)
	if err := query.Find(&equipment).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) Reload(ctx context.Context, equipment *domain.Equipment, with ...domain.EquipmentOption) error {
	query := r.applyOptions(ctx, with)
	if err := query.First(equipment, "id = ?", equipment.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Create(ctx context.Context, equipment *domain.Equipment) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Create(equipment).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Update(ctx context.Context, equipment *domain.Equipment) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(clause.Associations).Save(equipment).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Delete(&domain.Equipment{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

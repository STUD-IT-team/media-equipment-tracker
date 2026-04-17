package pgequipment

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresEquipmentRepository struct {
	db *gorm.DB
}

func NewPostgresEquipmentRepository(db *gorm.DB) *PostgresEquipmentRepository {
	return &PostgresEquipmentRepository{db: db}
}

var _ domain.EquipmentRepository = (*PostgresEquipmentRepository)(nil)

func (r *PostgresEquipmentRepository) applyOptions(opts []domain.EquipmentOption) *gorm.DB {
	options := &domain.EquipmentOptions{}
	for _, opt := range opts {
		opt(options)
	}
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *PostgresEquipmentRepository) Get(id uuid.UUID, with ...domain.EquipmentOption) (*domain.Equipment, error) {
	var eq domain.Equipment
	query := r.applyOptions(with)
	if err := query.First(&eq, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("Equipment", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &eq, nil
}

func (r *PostgresEquipmentRepository) GetUnoccupied(with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(with)
	if err := query.Where("current_invocation_id IS NULL").Find(&equipment).Error; err != nil {
		return nil, errs.NewRepositoryError("get_unoccupied", err)
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) List(with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(with)
	if err := query.Find(&equipment).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) Reload(equipment *domain.Equipment, with ...domain.EquipmentOption) error {
	query := r.applyOptions(with)
	if err := query.First(equipment, "id = ?", equipment.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Create(equipment *domain.Equipment) error {
	if err := r.db.Omit(clause.Associations).Create(equipment).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Update(equipment *domain.Equipment) error {
	if err := r.db.Omit(clause.Associations).Save(equipment).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.Equipment{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

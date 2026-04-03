package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"
)

type EquipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) domain.EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) applyOptions(opts []domain.EquipmentOption) *gorm.DB {
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

func (r *EquipmentRepository) Get(id uuid.UUID, with ...domain.EquipmentOption) (*domain.Equipment, error) {
	var eq domain.Equipment
	query := r.applyOptions(with)
	if err := query.First(&eq, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewEntityNotFoundError("Equipment", id)
		}
		return nil, errors.NewRepositoryError("get", err)
	}
	return &eq, nil
}

func (r *EquipmentRepository) GetUnoccupied(with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(with)
	if err := query.Where("current_invocation_id IS NULL").Find(&equipment).Error; err != nil {
		return nil, errors.NewRepositoryError("get_unoccupied", err)
	}
	return equipment, nil
}

func (r *EquipmentRepository) List(with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(with)
	if err := query.Find(&equipment).Error; err != nil {
		return nil, errors.NewRepositoryError("list", err)
	}
	return equipment, nil
}

func (r *EquipmentRepository) Reload(equipment *domain.Equipment, with ...domain.EquipmentOption) error {
	query := r.applyOptions(with)
	if err := query.First(equipment, "id = ?", equipment.ID).Error; err != nil {
		return errors.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *EquipmentRepository) Create(equipment *domain.Equipment) error {
	if err := r.db.Create(equipment).Error; err != nil {
		return errors.NewRepositoryError("create", err)
	}
	return nil
}

func (r *EquipmentRepository) Update(equipment *domain.Equipment) error {
	if err := r.db.Save(equipment).Error; err != nil {
		return errors.NewRepositoryError("update", err)
	}
	return nil
}

func (r *EquipmentRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.Equipment{}, "id = ?", id).Error; err != nil {
		return errors.NewRepositoryError("delete", err)
	}
	return nil
}

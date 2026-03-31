package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
)

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) domain.EquipmentRepository {
	return &equipmentRepository{db: db}
}

func (r *equipmentRepository) Create(ctx context.Context, equipment *domain.Equipment) error {
	return r.db.WithContext(ctx).Create(equipment).Error
}

func (r *equipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Equipment, error) {
	var equipment domain.Equipment
	err := r.db.WithContext(ctx).First(&equipment, id).Error
	if err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (r *equipmentRepository) GetAll(ctx context.Context) ([]*domain.Equipment, error) {
	var equipments []*domain.Equipment
	err := r.db.WithContext(ctx).Find(&equipments).Error
	return equipments, err
}

func (r *equipmentRepository) Update(ctx context.Context, equipment *domain.Equipment) error {
	return r.db.WithContext(ctx).Save(equipment).Error
}

func (r *equipmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Equipment{}, id).Error
}

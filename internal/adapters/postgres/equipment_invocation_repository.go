package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
)

type equipmentInvocationRepository struct {
	db *gorm.DB
}

func NewEquipmentInvocationRepository(db *gorm.DB) domain.EquipmentInvocationRepository {
	return &equipmentInvocationRepository{db: db}
}

func (r *equipmentInvocationRepository) Create(ctx context.Context, equipmentInvocation *domain.EquipmentInvocation) error {
	return r.db.WithContext(ctx).Create(equipmentInvocation).Error
}

func (r *equipmentInvocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.EquipmentInvocation, error) {
	var equipmentInvocation domain.EquipmentInvocation
	err := r.db.WithContext(ctx).First(&equipmentInvocation, id).Error
	if err != nil {
		return nil, err
	}
	return &equipmentInvocation, nil
}

func (r *equipmentInvocationRepository) GetAll(ctx context.Context) ([]*domain.EquipmentInvocation, error) {
	var equipmentInvocations []*domain.EquipmentInvocation
	err := r.db.WithContext(ctx).Find(&equipmentInvocations).Error
	return equipmentInvocations, err
}

func (r *equipmentInvocationRepository) Update(ctx context.Context, equipmentInvocation *domain.EquipmentInvocation) error {
	return r.db.WithContext(ctx).Save(equipmentInvocation).Error
}

func (r *equipmentInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.EquipmentInvocation{}, id).Error
}

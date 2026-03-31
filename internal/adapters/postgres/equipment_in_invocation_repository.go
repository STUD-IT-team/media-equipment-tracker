package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
)

type equipmentInInvocationRepository struct {
	db *gorm.DB
}

func NewEquipmentInInvocationRepository(db *gorm.DB) domain.EquipmentInInvocationRepository {
	return &equipmentInInvocationRepository{db: db}
}

func (r *equipmentInInvocationRepository) Create(ctx context.Context, equipmentInInvocation *domain.EquipmentInInvocation) error {
	return r.db.WithContext(ctx).Create(equipmentInInvocation).Error
}

func (r *equipmentInInvocationRepository) GetByInvocationIDAndEquipmentID(ctx context.Context, invocationID, equipmentID uuid.UUID) (*domain.EquipmentInInvocation, error) {
	var equipmentInInvocation domain.EquipmentInInvocation
	err := r.db.WithContext(ctx).Where("invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).First(&equipmentInInvocation).Error
	if err != nil {
		return nil, err
	}
	return &equipmentInInvocation, nil
}

func (r *equipmentInInvocationRepository) GetAllByInvocationID(ctx context.Context, invocationID uuid.UUID) ([]*domain.EquipmentInInvocation, error) {
	var equipmentInInvocations []*domain.EquipmentInInvocation
	err := r.db.WithContext(ctx).Where("invocation_id = ?", invocationID).Find(&equipmentInInvocations).Error
	return equipmentInInvocations, err
}

func (r *equipmentInInvocationRepository) GetAll(ctx context.Context) ([]*domain.EquipmentInInvocation, error) {
	var equipmentInInvocations []*domain.EquipmentInInvocation
	err := r.db.WithContext(ctx).Find(&equipmentInInvocations).Error
	return equipmentInInvocations, err
}

func (r *equipmentInInvocationRepository) Update(ctx context.Context, equipmentInInvocation *domain.EquipmentInInvocation) error {
	return r.db.WithContext(ctx).Save(equipmentInInvocation).Error
}

func (r *equipmentInInvocationRepository) Delete(ctx context.Context, invocationID, equipmentID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).Delete(&domain.EquipmentInInvocation{}).Error
}

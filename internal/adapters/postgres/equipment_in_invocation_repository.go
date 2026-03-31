package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"
)

type EquipmentInInvocationRepository struct {
	db *gorm.DB
}

func NewEquipmentInInvocationRepository(db *gorm.DB) domain.EquipmentInInvocationRepository {
	return &EquipmentInInvocationRepository{db: db}
}

func (r *EquipmentInInvocationRepository) applyOptions(opts []domain.EquipmentInInvocationOption) *gorm.DB {
	options := &domain.EquipmentInInvocationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *EquipmentInInvocationRepository) Get(invocationID uuid.UUID, equipmentID uuid.UUID, with ...domain.EquipmentInInvocationOption) (*domain.EquipmentInInvocation, error) {
	var eii domain.EquipmentInInvocation
	query := r.applyOptions(with)
	if err := query.First(&eii, "invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewEntityNotFoundError("EquipmentInInvocation", invocationID)
		}
		return nil, errors.NewRepositoryError("get", err)
	}
	return &eii, nil
}

func (r *EquipmentInInvocationRepository) GetByInvocation(invocationID uuid.UUID, with ...domain.EquipmentInInvocationOption) ([]*domain.EquipmentInInvocation, error) {
	var items []*domain.EquipmentInInvocation
	query := r.applyOptions(with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&items).Error; err != nil {
		return nil, errors.NewRepositoryError("get_by_invocation", err)
	}
	return items, nil
}

func (r *EquipmentInInvocationRepository) GetByEquipment(equipmentID uuid.UUID, with ...domain.EquipmentInInvocationOption) ([]*domain.EquipmentInInvocation, error) {
	var items []*domain.EquipmentInInvocation
	query := r.applyOptions(with)
	if err := query.Where("equipment_id = ?", equipmentID).Find(&items).Error; err != nil {
		return nil, errors.NewRepositoryError("get_by_equipment", err)
	}
	return items, nil
}

func (r *EquipmentInInvocationRepository) Create(equipmentInInvocation *domain.EquipmentInInvocation) error {
	if err := r.db.Create(equipmentInInvocation).Error; err != nil {
		return errors.NewRepositoryError("create", err)
	}
	return nil
}

func (r *EquipmentInInvocationRepository) Update(equipmentInInvocation *domain.EquipmentInInvocation) error {
	if err := r.db.Save(equipmentInInvocation).Error; err != nil {
		return errors.NewRepositoryError("update", err)
	}
	return nil
}

func (r *EquipmentInInvocationRepository) Delete(invocationID uuid.UUID, equipmentID uuid.UUID) error {
	if err := r.db.Delete(&domain.EquipmentInInvocation{}, "invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).Error; err != nil {
		return errors.NewRepositoryError("delete", err)
	}
	return nil
}

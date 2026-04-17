package pgequipmentininvocation

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresEquipmentInInvocationRepository struct {
	db *gorm.DB
}

func NewPostgresEquipmentInInvocationRepository(db *gorm.DB) *PostgresEquipmentInInvocationRepository {
	return &PostgresEquipmentInInvocationRepository{db: db}
}

var _ domain.EquipmentInInvocationRepository = (*PostgresEquipmentInInvocationRepository)(nil)

func (r *PostgresEquipmentInInvocationRepository) applyOptions(opts []domain.EquipmentInInvocationOption) *gorm.DB {
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

func (r *PostgresEquipmentInInvocationRepository) Get(invocationID, equipmentID uuid.UUID, with ...domain.EquipmentInInvocationOption) (*domain.EquipmentInInvocation, error) {
	var eii domain.EquipmentInInvocation
	query := r.applyOptions(with)
	if err := query.First(&eii, "invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("EquipmentInInvocation", invocationID)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &eii, nil
}

func (r *PostgresEquipmentInInvocationRepository) GetByInvocation(invocationID uuid.UUID, with ...domain.EquipmentInInvocationOption) ([]*domain.EquipmentInInvocation, error) {
	var items []*domain.EquipmentInInvocation
	query := r.applyOptions(with)
	if err := query.Where("invocation_id = ?", invocationID).Find(&items).Error; err != nil {
		return nil, errs.NewRepositoryError("get_by_invocation", err)
	}
	return items, nil
}

func (r *PostgresEquipmentInInvocationRepository) GetByEquipment(equipmentID uuid.UUID, with ...domain.EquipmentInInvocationOption) ([]*domain.EquipmentInInvocation, error) {
	var items []*domain.EquipmentInInvocation
	query := r.applyOptions(with)
	if err := query.Where("equipment_id = ?", equipmentID).Find(&items).Error; err != nil {
		return nil, errs.NewRepositoryError("get_by_equipment", err)
	}
	return items, nil
}

func (r *PostgresEquipmentInInvocationRepository) Create(equipmentInInvocation *domain.EquipmentInInvocation) error {
	if err := r.db.Omit(clause.Associations).Create(equipmentInInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresEquipmentInInvocationRepository) Update(equipmentInInvocation *domain.EquipmentInInvocation) error {
	if err := r.db.Omit(clause.Associations).Save(equipmentInInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresEquipmentInInvocationRepository) Delete(invocationID, equipmentID uuid.UUID) error {
	if err := r.db.Delete(&domain.EquipmentInInvocation{}, "invocation_id = ? AND equipment_id = ?", invocationID, equipmentID).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

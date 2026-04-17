package pgequipmentinvocation

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresEquipmentInvocationRepository struct {
	db *gorm.DB
}

func NewPostgresEquipmentInvocationRepository(db *gorm.DB) *PostgresEquipmentInvocationRepository {
	return &PostgresEquipmentInvocationRepository{db: db}
}

var _ domain.EquipmentInvocationRepository = (*PostgresEquipmentInvocationRepository)(nil)

func (r *PostgresEquipmentInvocationRepository) applyOptions(opts []domain.EquipmentInvocationOption) *gorm.DB {
	options := &domain.EquipmentInvocationOptions{}
	for _, opt := range opts {
		opt(options)
	}
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *PostgresEquipmentInvocationRepository) Get(id uuid.UUID, with ...domain.EquipmentInvocationOption) (*domain.EquipmentInvocation, error) {
	var inv domain.EquipmentInvocation
	query := r.applyOptions(with)
	if err := query.First(&inv, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("EquipmentInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &inv, nil
}

func (r *PostgresEquipmentInvocationRepository) List(with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
	var invocations []*domain.EquipmentInvocation
	query := r.applyOptions(with)
	if err := query.Find(&invocations).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return invocations, nil
}

func (r *PostgresEquipmentInvocationRepository) Reload(equipmentInvocation *domain.EquipmentInvocation, with ...domain.EquipmentInvocationOption) error {
	query := r.applyOptions(with)
	if err := query.First(equipmentInvocation, "id = ?", equipmentInvocation.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresEquipmentInvocationRepository) Create(equipmentInvocation *domain.EquipmentInvocation) error {
	tx := r.db.Begin()
	defer tx.Rollback()
	if err := tx.Omit(clause.Associations).Create(equipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}

	if equipmentInvocation.Equipment != nil {
		tx.Exec("DELETE FROM equipment_in_invocation WHERE invocation_id = ?", equipmentInvocation.ID)
		for _, e := range equipmentInvocation.Equipment {
			if err := tx.Exec(
				"INSERT INTO equipment_in_invocation (invocation_id, equipment_id, status) VALUES (?, ?, ?)",
				equipmentInvocation.ID, e.EquipmentID, e.Status,
			).Error; err != nil {
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (r *PostgresEquipmentInvocationRepository) Update(equipmentInvocation *domain.EquipmentInvocation) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	if err := tx.Omit(clause.Associations).Save(equipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}

	if equipmentInvocation.Equipment != nil {
		tx.Exec("DELETE FROM equipment_in_invocation WHERE invocation_id = ?", equipmentInvocation.ID)
		for _, e := range equipmentInvocation.Equipment {
			if err := tx.Exec(
				"INSERT INTO equipment_in_invocation (invocation_id, equipment_id, status) VALUES (?, ?, ?)",
				equipmentInvocation.ID, e.EquipmentID, e.Status,
			).Error; err != nil {
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *PostgresEquipmentInvocationRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.EquipmentInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

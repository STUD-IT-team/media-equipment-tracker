package postgresrepo

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EquipmentInvocationRepository struct {
	db *gorm.DB
}

func NewEquipmentInvocationRepository(db *gorm.DB) domain.EquipmentInvocationRepository {
	return &EquipmentInvocationRepository{db: db}
}

func (r *EquipmentInvocationRepository) applyOptions(opts []domain.EquipmentInvocationOption) *gorm.DB {
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

func (r *EquipmentInvocationRepository) Get(id uuid.UUID, with ...domain.EquipmentInvocationOption) (*domain.EquipmentInvocation, error) {
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

func (r *EquipmentInvocationRepository) List(with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
	var invocations []*domain.EquipmentInvocation
	query := r.applyOptions(with)
	if err := query.Find(&invocations).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return invocations, nil
}

func (r *EquipmentInvocationRepository) Reload(equipmentInvocation *domain.EquipmentInvocation, with ...domain.EquipmentInvocationOption) error {
	query := r.applyOptions(with)
	if err := query.First(equipmentInvocation, "id = ?", equipmentInvocation.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *EquipmentInvocationRepository) Create(equipmentInvocation *domain.EquipmentInvocation) error {
	if err := r.db.Create(equipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *EquipmentInvocationRepository) Update(equipmentInvocation *domain.EquipmentInvocation) error {
	if err := r.db.Save(equipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *EquipmentInvocationRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.EquipmentInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

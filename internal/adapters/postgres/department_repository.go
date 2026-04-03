package postgres

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) applyOptions(opts []domain.DepartmentOption) *gorm.DB {
	options := &domain.DepartmentOptions{}
	for _, opt := range opts {
		opt(options)
	}
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *DepartmentRepository) Get(id uuid.UUID, with ...domain.DepartmentOption) (*domain.Department, error) {
	var department domain.Department
	query := r.applyOptions(with)
	if err := query.First(&department, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewEntityNotFoundError("Department", id)
		}
		return nil, errors.NewRepositoryError("get", err)
	}
	return &department, nil
}

func (r *DepartmentRepository) List(with ...domain.DepartmentOption) ([]*domain.Department, error) {
	var departments []*domain.Department
	query := r.applyOptions(with)
	if err := query.Find(&departments).Error; err != nil {
		return nil, errors.NewRepositoryError("list", err)
	}
	return departments, nil
}

func (r *DepartmentRepository) Reload(department *domain.Department, with ...domain.DepartmentOption) error {
	query := r.applyOptions(with)
	if err := query.First(department, "id = ?", department.ID).Error; err != nil {
		return errors.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *DepartmentRepository) Create(department *domain.Department) error {
	if err := r.db.Create(department).Error; err != nil {
		return errors.NewRepositoryError("create", err)
	}
	return nil
}

func (r *DepartmentRepository) Update(department *domain.Department) error {
	if err := r.db.Save(department).Error; err != nil {
		return errors.NewRepositoryError("update", err)
	}
	return nil
}

func (r *DepartmentRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.Department{}, "id = ?", id).Error; err != nil {
		return errors.NewRepositoryError("delete", err)
	}
	return nil
}

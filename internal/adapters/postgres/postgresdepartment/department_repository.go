package postgresdepartment

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresDepartmentRepository struct {
	db *gorm.DB
}

func NewPostgresDepartmentRepository(db *gorm.DB) *PostgresDepartmentRepository {
	return &PostgresDepartmentRepository{db: db}
}

var _ domain.DepartmentRepository = (*PostgresDepartmentRepository)(nil)

func (r *PostgresDepartmentRepository) applyOptions(opts []domain.DepartmentOption) *gorm.DB {
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

func (r *PostgresDepartmentRepository) Get(id uuid.UUID, with ...domain.DepartmentOption) (*domain.Department, error) {
	var department domain.Department
	query := r.applyOptions(with)
	if err := query.First(&department, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("Department", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &department, nil
}

func (r *PostgresDepartmentRepository) List(with ...domain.DepartmentOption) ([]*domain.Department, error) {
	var departments []*domain.Department
	query := r.applyOptions(with)
	if err := query.Find(&departments).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return departments, nil
}

func (r *PostgresDepartmentRepository) Reload(department *domain.Department, with ...domain.DepartmentOption) error {
	query := r.applyOptions(with)
	if err := query.First(department, "id = ?", department.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresDepartmentRepository) Create(department *domain.Department) error {
	if err := r.db.Create(department).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresDepartmentRepository) Update(department *domain.Department) error {
	if err := r.db.Save(department).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresDepartmentRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.Department{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

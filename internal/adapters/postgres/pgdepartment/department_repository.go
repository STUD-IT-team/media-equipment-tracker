package pgdepartment

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	tx := r.db.Begin()
	defer tx.Rollback()
	if err := tx.Omit(clause.Associations).Create(department).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}

	if department.Equipment != nil {
		tx.Exec("DELETE FROM equipment_department WHERE department_id = ?", department.ID)
		for _, e := range department.Equipment {
			if err := tx.Exec(
				"INSERT INTO equipment_department (department_id, equipment_id) VALUES (?, ?)",
				department.ID, e.ID,
			).Error; err != nil {
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (r *PostgresDepartmentRepository) Update(department *domain.Department) error {
	tx := r.db.Begin()
	defer tx.Rollback()
	if err := tx.Omit(clause.Associations).Save(department).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}

	if department.Equipment != nil {
		tx.Exec("DELETE FROM equipment_department WHERE department_id = ?", department.ID)
		for _, e := range department.Equipment {
			if err := tx.Exec(
				"INSERT INTO equipment_department (department_id, equipment_id) VALUES (?, ?)",
				department.ID, e.ID,
			).Error; err != nil {
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (r *PostgresDepartmentRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.Department{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

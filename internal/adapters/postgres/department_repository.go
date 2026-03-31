package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"
)

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, department *domain.Department) error {
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	var department domain.Department
	err := r.db.WithContext(ctx).First(&department, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) GetAll(ctx context.Context) ([]*domain.Department, error) {
	var departments []*domain.Department
	err := r.db.WithContext(ctx).Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) Update(ctx context.Context, department *domain.Department) error {
	return r.db.WithContext(ctx).Save(department).Error
}

func (r *departmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Department{}, id).Error
}

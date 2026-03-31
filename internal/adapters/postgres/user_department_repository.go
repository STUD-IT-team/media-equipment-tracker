package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
)

type userDepartmentRepository struct {
	db *gorm.DB
}

func NewUserDepartmentRepository(db *gorm.DB) domain.UserDepartmentRepository {
	return &userDepartmentRepository{db: db}
}

func (r *userDepartmentRepository) Create(ctx context.Context, userDepartment *domain.UserDepartment) error {
	return r.db.WithContext(ctx).Create(userDepartment).Error
}

func (r *userDepartmentRepository) GetByUserIDAndDepartmentID(ctx context.Context, userID, departmentID uuid.UUID) (*domain.UserDepartment, error) {
	var userDepartment domain.UserDepartment
	err := r.db.WithContext(ctx).Where("user_id = ? AND department_id = ?", userID, departmentID).First(&userDepartment).Error
	if err != nil {
		return nil, err
	}
	return &userDepartment, nil
}

func (r *userDepartmentRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserDepartment, error) {
	var userDepartments []*domain.UserDepartment
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userDepartments).Error
	return userDepartments, err
}

func (r *userDepartmentRepository) GetAllByDepartmentID(ctx context.Context, departmentID uuid.UUID) ([]*domain.UserDepartment, error) {
	var userDepartments []*domain.UserDepartment
	err := r.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(&userDepartments).Error
	return userDepartments, err
}

func (r *userDepartmentRepository) GetAll(ctx context.Context) ([]*domain.UserDepartment, error) {
	var userDepartments []*domain.UserDepartment
	err := r.db.WithContext(ctx).Find(&userDepartments).Error
	return userDepartments, err
}

func (r *userDepartmentRepository) Update(ctx context.Context, userDepartment *domain.UserDepartment) error {
	return r.db.WithContext(ctx).Save(userDepartment).Error
}

func (r *userDepartmentRepository) Delete(ctx context.Context, userID, departmentID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND department_id = ?", userID, departmentID).Delete(&domain.UserDepartment{}).Error
}

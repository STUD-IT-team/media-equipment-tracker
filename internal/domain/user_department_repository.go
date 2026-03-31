package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserDepartmentRepository interface {
	Create(ctx context.Context, userDepartment *UserDepartment) error
	GetByUserIDAndDepartmentID(ctx context.Context, userID, departmentID uuid.UUID) (*UserDepartment, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*UserDepartment, error)
	GetAllByDepartmentID(ctx context.Context, departmentID uuid.UUID) ([]*UserDepartment, error)
	GetAll(ctx context.Context) ([]*UserDepartment, error)
	Update(ctx context.Context, userDepartment *UserDepartment) error
	Delete(ctx context.Context, userID, departmentID uuid.UUID) error
}

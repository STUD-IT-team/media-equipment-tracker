package domain

import "github.com/google/uuid"

type RoleInDepartment string

const (
	RoleTrainee  RoleInDepartment = "trainee"
	RoleActivist RoleInDepartment = "activist"
)

type UserDepartment struct {
	UserID       uuid.UUID        `gorm:"type:uuid;primaryKey" json:"user_id"`
	DepartmentID uuid.UUID        `gorm:"type:uuid;primaryKey" json:"department_id"`
	Role         RoleInDepartment `gorm:"type:role_in_department" json:"role"`

	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

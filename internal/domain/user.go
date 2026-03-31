package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	FullName     string    `gorm:"column:full_name;type:varchar(255);not null"`
	Email        string    `gorm:"type:varchar(255);unique;not null;index:idx_user_email"`
	HashPassword string    `gorm:"column:hash_password;type:varchar(255);not null"`
	Nice         int       `gorm:"type:int;check:nice > 0;not null"`
	IsAdmin      bool      `gorm:"column:is_admin;type:boolean"`

	// Сохраняется (только ID)
	Organizations []*Organization `gorm:"many2many:user_organization;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:organization_id"`
	// Cохраняется (только сама связь, не создаёт и не обновляет Department)
	Departments []*UserDepartment `gorm:"foreignKey:UserID"`

	// Сохраняется (только ID)
	EquipmentInvocations []*EquipmentInvocation `gorm:"foreignKey:UserID"`
	// Сохраняется (только ID)
	AdminEquipmentInvocations []*EquipmentInvocation `gorm:"foreignKey:AdminID"`

	// Сохраняется (только ID)
	StudioInvocations []*StudioInvocation `gorm:"foreignKey:UserID"`
	// Сохраняется (только ID)
	AdminStudioInvocations []*StudioInvocation `gorm:"foreignKey:AdminID"`
}

func (User) TableName() string {
	return "user"
}

type RoleInDepartment string

const (
	RoleTrainee  RoleInDepartment = "trainee"
	RoleActivist RoleInDepartment = "activist"
)

type UserDepartment struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	User   *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`

	DepartmentID uuid.UUID   `gorm:"type:uuid;primaryKey" json:"department_id"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`

	Role RoleInDepartment `gorm:"type:role_in_department" json:"role"`
}

func (UserDepartment) TableName() string {
	return "user_department"
}

type UserOptions struct {
	relations                     []string
	withOrganizations             bool
	withDepartments               bool
	withEquipmentInvocations      bool
	withAdminEquipmentInvocations bool
	withStudioInvocations         bool
	withAdminStudioInvocations    bool
}

func (o *UserOptions) Relations() []string {
	return o.relations
}

type UserOption func(options *UserOptions)

func UserWithOrganizations() UserOption {
	return func(options *UserOptions) {
		options.withOrganizations = true
		options.relations = append(options.relations, "Organizations")
	}
}

func UserWithDepartments() UserOption {
	return func(options *UserOptions) {
		options.withDepartments = true
		options.relations = append(options.relations, "Departments")
	}
}

func UserWithEquipmentInvocations() UserOption {
	return func(options *UserOptions) {
		options.withEquipmentInvocations = true
		options.relations = append(options.relations, "EquipmentInvocations")
	}
}

func UserWithAdminEquipmentInvocations() UserOption {
	return func(options *UserOptions) {
		options.withAdminEquipmentInvocations = true
		options.relations = append(options.relations, "AdminEquipmentInvocations")
	}
}

func UserWithStudioInvocations() UserOption {
	return func(options *UserOptions) {
		options.withStudioInvocations = true
		options.relations = append(options.relations, "StudioInvocations")
	}
}

func UserWithAdminStudioInvocations() UserOption {
	return func(options *UserOptions) {
		options.withAdminStudioInvocations = true
		options.relations = append(options.relations, "AdminStudioInvocations")
	}
}

type UserRepository interface {
	Get(id uuid.UUID, with ...UserOption) (*User, error)
	List(with ...UserOption) ([]*User, error)
	Reload(user *User, with ...UserOption) error

	Create(user *User) error
	Update(user *User) error
	Delete(id uuid.UUID) error
}

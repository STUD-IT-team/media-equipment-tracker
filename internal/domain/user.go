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

	Organizations []*Organization `gorm:"many2many:user_organization;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:organization_id"`
	//Departments                 []*Department                 `gorm:"many2many:user_department;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:department_id"`
	UserDepartments []*UserDepartment `gorm:"foreignKey:UserID"`
	//EquipmentInvocations        []*EquipmentInvocation        `gorm:"foreignKey:UserID"`
	//AdminEquipmentInvocations   []*EquipmentInvocation        `gorm:"foreignKey:AdminID"`
	//EquipmentInvocationMessages []*MessageEquipmentInvocation `gorm:"foreignKey:SenderID"`
	//ReceivedEquipmentMessages   []*MessageEquipmentInvocation `gorm:"foreignKey:RecipientID"`
	//StudioInvocations           []*StudioInvocation           `gorm:"foreignKey:UserID"`
	//AdminStudioInvocations      []*StudioInvocation           `gorm:"foreignKey:AdminID"`
	//StudioInvocationMessages    []*MessageStudioInvocation    `gorm:"foreignKey:SenderID"`
	//ReceivedStudioMessages      []*MessageStudioInvocation    `gorm:"foreignKey:RecipientID"`
}

func (User) TableName() string {
	return "user"
}

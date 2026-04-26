package departmentservice

import (
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/txmanager"
)

type DepartmentService interface {
	CreateDepartmentService
	ListDepartmentService
	GetDepartmentService
	UpdateDepartmentService
	DeleteDepartmentService
	GetDepartmentUsersService
	GetDepartmentEquipmentService
}

type departmentService struct {
	CreateDepartmentService
	ListDepartmentService
	GetDepartmentService
	UpdateDepartmentService
	DeleteDepartmentService
	GetDepartmentUsersService
	GetDepartmentEquipmentService
}

var _ DepartmentService = (*departmentService)(nil)

func NewDepartmentService(
	auther authzservice.AuthZ,
	departmentRepository domain.DepartmentRepository,
	userRepository domain.UserRepository,
	equipmentRepository domain.EquipmentRepository,
	txManager txmanager.TxManager,
) DepartmentService {
	return &departmentService{
		CreateDepartmentService:       NewCreateDepartmentService(auther, departmentRepository, txManager),
		ListDepartmentService:         NewListDepartmentService(departmentRepository),
		GetDepartmentService:          NewGetDepartmentService(departmentRepository),
		UpdateDepartmentService:       NewUpdateDepartmentService(auther, departmentRepository, txManager),
		DeleteDepartmentService:       NewDeleteDepartmentService(auther, departmentRepository, txManager),
		GetDepartmentUsersService:     NewGetDepartmentUsersService(departmentRepository),
		GetDepartmentEquipmentService: NewGetDepartmentEquipmentService(departmentRepository),
	}
}

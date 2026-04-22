package equipmentservice

import "media-equipment-tracker/internal/domain"

type EquipmentService interface {
	SearchEquipmentService
}

type equipmentService struct {
	SearchEquipmentService
	equipmentRepository domain.EquipmentRepository
}

var _ EquipmentService = (*equipmentService)(nil)

func NewEquipmentService(equipmentRepository domain.EquipmentRepository, searchRepository SearchEquipmentRepository) EquipmentService {
	return &equipmentService{
		equipmentRepository:    equipmentRepository,
		SearchEquipmentService: NewSearchEquipmentService(searchRepository),
	}
}

package pgequipment

import (
	"context"
	"errors"

	"media-equipment-tracker/internal/adapters/postgres/pgutils"
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresEquipmentRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresEquipmentRepository(db *gormtx.DBGetter) *PostgresEquipmentRepository {
	return &PostgresEquipmentRepository{db: db}
}

var _ domain.EquipmentRepository = (*PostgresEquipmentRepository)(nil)

func (r *PostgresEquipmentRepository) applyOptions(ctx context.Context, opts []domain.EquipmentOption) *gorm.DB {
	options := &domain.EquipmentOptions{}
	for _, opt := range opts {
		opt(options)
	}
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return nil
	}
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}
	return query
}

func (r *PostgresEquipmentRepository) Get(ctx context.Context, id uuid.UUID, with ...domain.EquipmentOption) (*domain.Equipment, error) {
	var eq domain.Equipment
	query := r.applyOptions(ctx, with)
	if err := query.First(&eq, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewEntityNotFoundError("Equipment", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &eq, nil
}

func (r *PostgresEquipmentRepository) GetUnoccupied(ctx context.Context, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(ctx, with)
	if err := query.Where("current_invocation_id IS NULL").Find(&equipment).Error; err != nil {
		return nil, errs.NewRepositoryError("get_unoccupied", err)
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) GetByInventoryNumber(ctx context.Context, inventoryNumber string, with ...domain.EquipmentOption) (*domain.Equipment, error) {
	var equipment *domain.Equipment
	query := r.applyOptions(ctx, with)
	if err := query.Where("inventory_number = ?", inventoryNumber).First(&equipment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewEntityNotFoundError("Equipment", inventoryNumber)
		}
		return nil, errs.NewRepositoryError("get_by_inventory_number", err)
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) List(ctx context.Context, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	var equipment []*domain.Equipment
	query := r.applyOptions(ctx, with)
	if err := query.Find(&equipment).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) Reload(ctx context.Context, equipment *domain.Equipment, with ...domain.EquipmentOption) error {
	query := r.applyOptions(ctx, with)
	if err := query.First(equipment, "id = ?", equipment.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Create(ctx context.Context, equipment *domain.Equipment) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("create", err)
	}

	if err := db.Omit(clause.Associations).Create(equipment).Error; err != nil {
		if pgutils.IsUniqueViolationError(err) {
			return errs.NewEntityAlreadyExistsError("Equipment", equipment)
		}
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Update(ctx context.Context, equipment *domain.Equipment) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("update", err)
	}
	if err := db.Omit(clause.Associations).Save(equipment).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	if err := db.Delete(&domain.Equipment{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

func (r *PostgresEquipmentRepository) Search(ctx context.Context, search *equipmentservice.SearchEquipmentRequest, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	equipment := make([]*domain.Equipment, 0)
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return nil, errs.NewRepositoryError("search", err)
	}
	query := r.applyOptions(ctx, with)

	if search.SearchString != nil {
		query = query.Where("name LIKE ? OR short_name LIKE ?", "%"+*search.SearchString+"%", "%"+*search.SearchString+"%")
	}
	if search.Categories != nil {
		query = query.Where("category IN (?)", search.Categories)
	}
	if search.Statuses != nil {
		query = query.Where("status IN (?)", search.Statuses)
	}
	if search.AvailableToTrainee != nil {
		query = query.Where("available_to_trainee = ?", *search.AvailableToTrainee)
	}
	if search.DepartmentIDs != nil {
		subQuery := db.Table("equipment_department").
			Where("equipment_id = equipment.id").
			Where("department_id IN (?)", search.DepartmentIDs).
			Select("COUNT(DISTINCT department_id)")

		query = query.Where("(?) = ?", subQuery, len(search.DepartmentIDs))
	}
	if search.AvailableAt != nil {
		subQuery := db.Table("equipment_in_invocation").
			Select("1").
			Joins("equipment_invocation ON equipment_invocation.id = equipment_in_invocation.invocation_id").
			Where("equipment_in_invocation.equipment_id = equipment.id").
			Where("? BETWEEN equipment_invocation.start_time AND equipment_invocation.end_time", search.AvailableAt)

		query = query.Where("NOT EXISTS (?)", subQuery)
	}

	if err := query.Find(&equipment).Error; err != nil {
		return nil, errs.NewRepositoryError("search", err)
	}

	return equipment, nil
}

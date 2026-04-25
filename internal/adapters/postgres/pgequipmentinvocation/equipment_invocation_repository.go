package pgequipmentinvocation

import (
	"context"
	"errors"

	"media-equipment-tracker/internal/application/invocationservice/invocationsearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresEquipmentInvocationRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresEquipmentInvocationRepository(db *gormtx.DBGetter) *PostgresEquipmentInvocationRepository {
	return &PostgresEquipmentInvocationRepository{db: db}
}

var _ domain.EquipmentInvocationRepository = (*PostgresEquipmentInvocationRepository)(nil)

func (r *PostgresEquipmentInvocationRepository) applyOptions(ctx context.Context, opts []domain.EquipmentInvocationOption) *gorm.DB {
	options := &domain.EquipmentInvocationOptions{}
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

func (r *PostgresEquipmentInvocationRepository) Get(ctx context.Context, id uuid.UUID, with ...domain.EquipmentInvocationOption) (*domain.EquipmentInvocation, error) {
	var inv domain.EquipmentInvocation
	query := r.applyOptions(ctx, with)
	if err := query.First(&inv, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewEntityNotFoundError("EquipmentInvocation", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &inv, nil
}

func (r *PostgresEquipmentInvocationRepository) List(ctx context.Context, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
	var invocations []*domain.EquipmentInvocation
	query := r.applyOptions(ctx, with)
	if err := query.Find(&invocations).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return invocations, nil
}

func (r *PostgresEquipmentInvocationRepository) Reload(ctx context.Context, equipmentInvocation *domain.EquipmentInvocation, with ...domain.EquipmentInvocationOption) error {
	query := r.applyOptions(ctx, with)
	if err := query.First(equipmentInvocation, "id = ?", equipmentInvocation.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

func (r *PostgresEquipmentInvocationRepository) upsertOmitFields() []string {
	return []string{"Equipment.Equipment", "Equipment.Invocation", "Admin", "User", "Organization", "Department"}
}

func (r *PostgresEquipmentInvocationRepository) Create(ctx context.Context, equipmentInvocation *domain.EquipmentInvocation) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("create", err)
	}
	if err := db.Omit(r.upsertOmitFields()...).Create(equipmentInvocation).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}

	return nil
}

func (r *PostgresEquipmentInvocationRepository) Update(ctx context.Context, equipmentInvocation *domain.EquipmentInvocation) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("update", err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(r.upsertOmitFields()...).Save(equipmentInvocation).Error; err != nil {
			return errs.NewRepositoryError("update", err)
		}

		if equipmentInvocation.Equipment != nil {
			if err := tx.Model(equipmentInvocation).Association("Equipment").Replace(equipmentInvocation.Equipment); err != nil {
				return errs.NewRepositoryError("update", err)
			}
		}
		return nil
	})

	if err != nil {
		return errs.NewRepositoryError("update", err)
	}

	return nil
}

func (r *PostgresEquipmentInvocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	if err := db.Delete(&domain.EquipmentInvocation{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

func (r *PostgresEquipmentInvocationRepository) Search(ctx context.Context, req *invocationsearch.SearchInvocationRequest, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return nil, errs.NewRepositoryError("search", err)
	}

	var invocations []*domain.EquipmentInvocation
	query := r.applyOptions(ctx, with)
	if req.SearchString != nil {
		query = query.Where("event_name like ?", "%"+*req.SearchString+"%")
	}

	// Если хотя бы кусочек события в промежутке, то подходит
	if req.StartTime != nil {
		query = query.Where("end_time >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("start_time <= ?", *req.EndTime)
	}

	if req.EquipmentIDs != nil {
		subQuery := db.Model(&domain.EquipmentInInvocation{}).
			Where("invocation_id = equipment_invocation.id").
			Where("equipment_id IN ?", req.EquipmentIDs).
			Select("COUNT(DISTINCT equipment_id)")

		query = query.Where("(?) = ?", subQuery, len(req.EquipmentIDs))
	}

	if req.Statuses != nil {
		query = query.Where("status IN (?)", req.Statuses)
	}

	if req.AdminID != nil {
		query = query.Where("admin_id = ?", *req.AdminID)
	}

	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}

	if err := query.Find(&invocations).Error; err != nil {
		return nil, errs.NewRepositoryError("search", err)
	}

	return invocations, nil
}

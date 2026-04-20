package pgdepartment

import (
	"context"

	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresDepartmentRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresDepartmentRepository(db *gormtx.DBGetter) *PostgresDepartmentRepository {
	return &PostgresDepartmentRepository{db: db}
}

var _ domain.DepartmentRepository = (*PostgresDepartmentRepository)(nil)

func (r *PostgresDepartmentRepository) applyOptions(ctx context.Context, opts []domain.DepartmentOption) *gorm.DB {
	options := &domain.DepartmentOptions{}
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

func (r *PostgresDepartmentRepository) Get(ctx context.Context, id uuid.UUID, with ...domain.DepartmentOption) (*domain.Department, error) {
	var department domain.Department
	query := r.applyOptions(ctx, with)
	if err := query.First(&department, "id = ?", id).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("Department", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}
	return &department, nil
}

func (r *PostgresDepartmentRepository) List(ctx context.Context, with ...domain.DepartmentOption) ([]*domain.Department, error) {
	var departments []*domain.Department
	query := r.applyOptions(ctx, with)
	if err := query.Find(&departments).Error; err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}
	return departments, nil
}

func (r *PostgresDepartmentRepository) Reload(ctx context.Context, department *domain.Department, with ...domain.DepartmentOption) error {
	query := r.applyOptions(ctx, with)
	if err := query.First(department, "id = ?", department.ID).Error; err != nil {
		return errs.NewRepositoryError("reload", err)
	}
	return nil
}

//go:inline
func (r *PostgresDepartmentRepository) upsertOmitFields() []string {
	// * для того, чтобы во всех операциях сохраняласть только связь, в том числе: Create, Save, Replace
	return []string{"Equipment.*", "EquipmentInvocations", "StudioInvocations", "Users"}
}

func (r *PostgresDepartmentRepository) Create(ctx context.Context, department *domain.Department) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("create", err)
	}
	if err := db.Omit(r.upsertOmitFields()...).Create(department).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}

	return nil
}

func (r *PostgresDepartmentRepository) Update(ctx context.Context, department *domain.Department) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("update", err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(r.upsertOmitFields()...).Save(department).Error; err != nil {
			return errs.NewRepositoryError("update", err)
		}
		// Нужно, так как Save может только добавлять связи, но не удалять существующие в БД
		if department.Equipment != nil {
			if err := tx.Model(department).Association("Equipment").Replace(department.Equipment); err != nil {
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

func (r *PostgresDepartmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, err := r.db.GetDB(ctx)
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	if err := db.Delete(&domain.Department{}, "id = ?", id).Error; err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

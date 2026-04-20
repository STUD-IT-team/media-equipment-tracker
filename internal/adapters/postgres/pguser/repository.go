package pguser

import (
	"context"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresUserRepository(db *gormtx.DBGetter) domain.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Get(ctx context.Context, id uuid.UUID, opts ...domain.UserOption) (*domain.User, error) {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, _ := r.db.GetDB(ctx)
	var user domain.User
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.First(&user, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("User", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string, opts ...domain.UserOption) (*domain.User, error) {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, _ := r.db.GetDB(ctx)
	var user domain.User
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.First(&user, "email = ?", email).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("User", email)
		}
		return nil, errs.NewRepositoryError("get by email", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) List(ctx context.Context, opts ...domain.UserOption) ([]*domain.User, error) {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, _ := r.db.GetDB(ctx)
	var users []*domain.User
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}

	return users, nil
}

func (r *PostgresUserRepository) Reload(ctx context.Context, user *domain.User, opts ...domain.UserOption) error {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, _ := r.db.GetDB(ctx)
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(user, "id = ?", user.ID).Error
	if err != nil {
		return errs.NewRepositoryError("reload", err)
	}

	return nil
}

//go:inline
func (r *PostgresUserRepository) upsertOmitFields() []string {
	return []string{"Organizations.*", "Departments.User", "Departments.Department", "EquipmentInvocations", "AdminEquipmentInvocations", "StudioInvocations", "AdminStudioInvocations"}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	db, _ := r.db.GetDB(ctx)
	if err := db.Omit(r.upsertOmitFields()...).Create(user).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}

	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	db, _ := r.db.GetDB(ctx)

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(r.upsertOmitFields()...).Save(user).Error; err != nil {
			return errs.NewRepositoryError("update", err)
		}

		if user.Organizations != nil {
			if err := tx.Model(user).Association("Organizations").Replace(user.Organizations); err != nil {
				return errs.NewRepositoryError("update", err)
			}
		}

		if user.Departments != nil {
			if err := tx.Model(user).Association("Departments").Replace(user.Departments); err != nil {
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

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, _ := r.db.GetDB(ctx)
	err := db.Delete(&domain.User{}, "id = ?", id).Error
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

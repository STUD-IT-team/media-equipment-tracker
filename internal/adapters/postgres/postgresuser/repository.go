package postgresuser

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Get(id uuid.UUID, opts ...domain.UserOption) (*domain.User, error) {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var user domain.User
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("User", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string, opts ...domain.UserOption) (*domain.User, error) {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var user domain.User
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("User", email)
		}
		return nil, errs.NewRepositoryError("get by email", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) List(opts ...domain.UserOption) ([]*domain.User, error) {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var users []*domain.User
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}

	return users, nil
}

func (r *PostgresUserRepository) Reload(user *domain.User, opts ...domain.UserOption) error {
	options := &domain.UserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(user, "id = ?", user.ID).Error
	if err != nil {
		return errs.NewRepositoryError("reload", err)
	}

	return nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresUserRepository) Delete(id uuid.UUID) error {
	err := r.db.Delete(&domain.User{}, "id = ?", id).Error
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

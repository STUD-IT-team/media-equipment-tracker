package pguser

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	tx := r.db.Begin()
	defer tx.Rollback()
	if err := tx.Omit(clause.Associations).Create(user).Error; err != nil {
		return errs.NewRepositoryError("create", err)
	}

	if user.Organizations != nil {
		tx.Exec("DELETE FROM user_organization WHERE user_id = ?", user.ID)
		for _, o := range user.Organizations {
			if err := tx.Exec(
				"INSERT INTO user_organization (user_id, organization_id) VALUES (?, ?)",
				user.ID, o.ID,
			).Error; err != nil {
				return err
			}
		}
	}
	if user.Departments != nil {
		tx.Exec("DELETE FROM user_department WHERE user_id = ?", user.ID)
		for _, d := range user.Departments {
			if err := tx.Exec(
				"INSERT INTO user_department (user_id, department_id, role) VALUES (?, ?, ?)",
				user.ID, d.DepartmentID, d.Role,
			).Error; err != nil {
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	tx := r.db.Begin()
	defer tx.Rollback()
	if err := tx.Omit(clause.Associations).Updates(user).Error; err != nil {
		return errs.NewRepositoryError("update", err)
	}

	if user.Organizations != nil {
		tx.Exec("DELETE FROM user_organization WHERE user_id = ?", user.ID)
		for _, o := range user.Organizations {
			if err := tx.Exec(
				"INSERT INTO user_organization (user_id, organization_id) VALUES (?, ?)",
				user.ID, o.ID,
			).Error; err != nil {
				return err
			}
		}
	}
	if user.Departments != nil {
		tx.Exec("DELETE FROM user_department WHERE user_id = ?", user.ID)
		for _, d := range user.Departments {
			if err := tx.Exec(
				"INSERT INTO user_department (user_id, department_id, role) VALUES (?, ?, ?)",
				user.ID, d.DepartmentID, d.Role,
			).Error; err != nil {
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (r *PostgresUserRepository) Delete(id uuid.UUID) error {
	err := r.db.Delete(&domain.User{}, "id = ?", id).Error
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

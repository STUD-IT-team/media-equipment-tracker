package postgresorganization

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresOrganizationRepository struct {
	db *gorm.DB
}

func NewPostgresOrganizationRepository(db *gorm.DB) *PostgresOrganizationRepository {
	return &PostgresOrganizationRepository{db: db}
}

var _ domain.OrganizationRepository = (*PostgresOrganizationRepository)(nil)

func (r *PostgresOrganizationRepository) Get(id uuid.UUID, opts ...domain.OrganizationOption) (*domain.Organization, error) {
	options := &domain.OrganizationOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var organization domain.Organization
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.First(&organization, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewEntityNotFoundError("Organization", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}

	return &organization, nil
}

func (r *PostgresOrganizationRepository) List(opts ...domain.OrganizationOption) ([]*domain.Organization, error) {
	options := &domain.OrganizationOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var organizations []*domain.Organization
	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(&organizations).Error
	if err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}

	return organizations, nil
}

func (r *PostgresOrganizationRepository) Reload(organization *domain.Organization, opts ...domain.OrganizationOption) error {
	options := &domain.OrganizationOptions{}
	for _, opt := range opts {
		opt(options)
	}

	query := r.db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(organization, "id = ?", organization.ID).Error
	if err != nil {
		return errs.NewRepositoryError("reload", err)
	}

	return nil
}

func (r *PostgresOrganizationRepository) Create(organization *domain.Organization) error {
	err := r.db.Create(organization).Error
	if err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresOrganizationRepository) Update(organization *domain.Organization) error {
	err := r.db.Save(organization).Error
	if err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresOrganizationRepository) Delete(id uuid.UUID) error {
	err := r.db.Delete(&domain.Organization{}, "id = ?", id).Error
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

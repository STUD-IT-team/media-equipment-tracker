package postgres

import (
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) domain.OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Get(id uuid.UUID, opts ...domain.OrganizationOption) (*domain.Organization, error) {
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
			return nil, errors.NewEntityNotFoundError("Organization", id)
		}
		return nil, errors.NewRepositoryError("get", err)
	}

	return &organization, nil
}

func (r *OrganizationRepository) List(opts ...domain.OrganizationOption) ([]*domain.Organization, error) {
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
		return nil, errors.NewRepositoryError("list", err)
	}

	return organizations, nil
}

func (r *OrganizationRepository) Reload(organization *domain.Organization, opts ...domain.OrganizationOption) error {
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
		return errors.NewRepositoryError("reload", err)
	}

	return nil
}

func (r *OrganizationRepository) Create(organization *domain.Organization) error {
	err := r.db.Create(organization).Error
	if err != nil {
		return errors.NewRepositoryError("create", err)
	}
	return nil
}

func (r *OrganizationRepository) Update(organization *domain.Organization) error {
	err := r.db.Save(organization).Error
	if err != nil {
		return errors.NewRepositoryError("update", err)
	}
	return nil
}

func (r *OrganizationRepository) Delete(id uuid.UUID) error {
	err := r.db.Delete(&domain.Organization{}, "id = ?", id).Error
	if err != nil {
		return errors.NewRepositoryError("delete", err)
	}
	return nil
}

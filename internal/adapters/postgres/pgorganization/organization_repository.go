package pgorganization

import (
	"context"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager/gormtx"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type PostgresOrganizationRepository struct {
	db *gormtx.DBGetter
}

func NewPostgresOrganizationRepository(db *gormtx.DBGetter) *PostgresOrganizationRepository {
	return &PostgresOrganizationRepository{db: db}
}

var _ domain.OrganizationRepository = (*PostgresOrganizationRepository)(nil)

func (r *PostgresOrganizationRepository) Get(ctx context.Context, id uuid.UUID, opts ...domain.OrganizationOption) (*domain.Organization, error) {
	options := &domain.OrganizationOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, _ := r.db.GetDB(ctx)
	var organization domain.Organization
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.First(&organization, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errs.NewEntityNotFoundError("Organization", id)
		}
		return nil, errs.NewRepositoryError("get", err)
	}

	return &organization, nil
}

func (r *PostgresOrganizationRepository) List(ctx context.Context, opts ...domain.OrganizationOption) ([]*domain.Organization, error) {
	options := &domain.OrganizationOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, _ := r.db.GetDB(ctx)
	var organizations []*domain.Organization
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(&organizations).Error
	if err != nil {
		return nil, errs.NewRepositoryError("list", err)
	}

	return organizations, nil
}

func (r *PostgresOrganizationRepository) Reload(ctx context.Context, organization *domain.Organization, opts ...domain.OrganizationOption) error {
	options := &domain.OrganizationOptions{}
	for _, opt := range opts {
		opt(options)
	}

	db, _ := r.db.GetDB(ctx)
	query := db
	for _, rel := range options.Relations() {
		query = query.Preload(rel)
	}

	err := query.Find(organization, "id = ?", organization.ID).Error
	if err != nil {
		return errs.NewRepositoryError("reload", err)
	}

	return nil
}

func (r *PostgresOrganizationRepository) Create(ctx context.Context, organization *domain.Organization) error {
	db, _ := r.db.GetDB(ctx)
	err := db.Omit(clause.Associations).Create(organization).Error
	if err != nil {
		return errs.NewRepositoryError("create", err)
	}
	return nil
}

func (r *PostgresOrganizationRepository) Update(ctx context.Context, organization *domain.Organization) error {
	db, _ := r.db.GetDB(ctx)
	err := db.Omit(clause.Associations).Save(organization).Error
	if err != nil {
		return errs.NewRepositoryError("update", err)
	}
	return nil
}

func (r *PostgresOrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, _ := r.db.GetDB(ctx)
	err := db.Delete(&domain.Organization{}, "id = ?", id).Error
	if err != nil {
		return errs.NewRepositoryError("delete", err)
	}
	return nil
}

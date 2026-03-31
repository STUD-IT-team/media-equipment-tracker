package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/domain"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) domain.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, organization *domain.Organization) error {
	return r.db.WithContext(ctx).Create(organization).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	var organization domain.Organization
	err := r.db.WithContext(ctx).First(&organization, id).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *organizationRepository) GetAll(ctx context.Context) ([]*domain.Organization, error) {
	var organizations []*domain.Organization
	err := r.db.WithContext(ctx).Find(&organizations).Error
	return organizations, err
}

func (r *organizationRepository) Update(ctx context.Context, organization *domain.Organization) error {
	return r.db.WithContext(ctx).Save(organization).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Organization{}, id).Error
}

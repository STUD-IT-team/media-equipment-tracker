package organizationservice

import (
	"context"

	"media-equipment-tracker/internal/domain"

	"github.com/google/uuid"
)

type GetOrganizationService interface {
	GetOrganization(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
}

type getOrganizationService struct {
	organizationRepository domain.OrganizationRepository
}

var _ GetOrganizationService = (*getOrganizationService)(nil)

func NewGetOrganizationService(
	organizationRepository domain.OrganizationRepository,
) GetOrganizationService {
	return &getOrganizationService{
		organizationRepository: organizationRepository,
	}
}

func (s *getOrganizationService) GetOrganization(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	return s.organizationRepository.Get(ctx, id, domain.WithUsers())
}

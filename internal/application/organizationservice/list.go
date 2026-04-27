package organizationservice

import (
	"context"
	"strings"

	"media-equipment-tracker/internal/domain"
)

type ListOrganizationService interface {
	ListOrganizations(ctx context.Context, search string) ([]*domain.Organization, error)
}

type listOrganizationService struct {
	organizationRepository domain.OrganizationRepository
}

var _ ListOrganizationService = (*listOrganizationService)(nil)

func NewListOrganizationService(
	organizationRepository domain.OrganizationRepository,
) ListOrganizationService {
	return &listOrganizationService{
		organizationRepository: organizationRepository,
	}
}

func (s *listOrganizationService) ListOrganizations(ctx context.Context, search string) ([]*domain.Organization, error) {
	organizations, err := s.organizationRepository.List(ctx, domain.WithUsers())
	if err != nil {
		return nil, err
	}

	// Simple search by name (case insensitive)
	if search != "" {
		filtered := make([]*domain.Organization, 0)
		for _, org := range organizations {
			if strings.Contains(strings.ToLower(org.Name), strings.ToLower(search)) {
				filtered = append(filtered, org)
			}
		}
		return filtered, nil
	}

	return organizations, nil
}

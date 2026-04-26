package organizationservice

import (
	"context"

	"media-equipment-tracker/internal/domain"

	"github.com/google/uuid"
)

type GetOrganizationUsersService interface {
	GetOrganizationUsers(ctx context.Context, id uuid.UUID) ([]*domain.User, error)
}

type getOrganizationUsersService struct {
	organizationRepository domain.OrganizationRepository
}

var _ GetOrganizationUsersService = (*getOrganizationUsersService)(nil)

func NewGetOrganizationUsersService(
	organizationRepository domain.OrganizationRepository,
) GetOrganizationUsersService {
	return &getOrganizationUsersService{
		organizationRepository: organizationRepository,
	}
}

func (s *getOrganizationUsersService) GetOrganizationUsers(ctx context.Context, id uuid.UUID) ([]*domain.User, error) {
	// First check if organization exists
	_, err := s.organizationRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Since Users is many2many, we need to preload it
	organization, err := s.organizationRepository.Get(ctx, id, domain.WithUsers())
	if err != nil {
		return nil, err
	}

	if organization.Users == nil {
		return []*domain.User{}, nil
	}

	return organization.Users, nil
}
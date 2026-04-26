package organizationservice

import (
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/txmanager"
)

type OrganizationService interface {
	CreateOrganizationService
	ListOrganizationService
	GetOrganizationService
	UpdateOrganizationService
	DeleteOrganizationService
	GetOrganizationUsersService
}

type organizationService struct {
	CreateOrganizationService
	ListOrganizationService
	GetOrganizationService
	UpdateOrganizationService
	DeleteOrganizationService
	GetOrganizationUsersService
}

var _ OrganizationService = (*organizationService)(nil)

func NewOrganizationService(
	auther authzservice.AuthZ,
	organizationRepository domain.OrganizationRepository,
	txManager txmanager.TxManager,
) OrganizationService {
	return &organizationService{
		CreateOrganizationService:   NewCreateOrganizationService(auther, organizationRepository, txManager),
		ListOrganizationService:     NewListOrganizationService(organizationRepository),
		GetOrganizationService:      NewGetOrganizationService(organizationRepository),
		UpdateOrganizationService:   NewUpdateOrganizationService(auther, organizationRepository, txManager),
		DeleteOrganizationService:   NewDeleteOrganizationService(auther, organizationRepository, txManager),
		GetOrganizationUsersService: NewGetOrganizationUsersService(organizationRepository),
	}
}

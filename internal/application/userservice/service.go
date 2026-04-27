package userservice

import (
	"context"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type UserService interface {
	UpdateUserService
	MeUserService
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	UpdateUserService
	MeUserService
	userRep domain.UserRepository
	authz   authzservice.AuthZ
	txm     txmanager.TxManager
}

func NewUserService(
	userRep domain.UserRepository,
	authz authzservice.AuthZ,
	txm txmanager.TxManager,
) (UserService, error) {
	service := &userService{
		UpdateUserService: NewUpdateUserService(userRep, txm, authz),
		MeUserService:     NewMeUserService(userRep, authz, txm),
		userRep:           userRep,
		authz:             authz,
		txm:               txm,
	}
	return service, nil
}

func (s *userService) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	payload, err := s.authz.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	user, err := s.userRep.Get(ctx, id, domain.UserWithDepartments(), domain.UserWithOrganizations(), domain.UserWithEquipmentInvocations(), domain.UserWithAdminEquipmentInvocations(), domain.UserWithStudioInvocations(), domain.UserWithAdminStudioInvocations())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	payload, err := s.authz.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	users, err := s.userRep.List(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	payload, err := s.authz.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	if payload.UserID == id {
		return errs.NewValidationError("id", "can't delete current user")
	}

	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		user, err := s.userRep.Get(ctx, id, domain.UserWithStudioInvocations(), domain.UserWithEquipmentInvocations())
		if err != nil {
			return err
		}

		if len(user.EquipmentInvocations) > 0 || len(user.StudioInvocations) > 0 {
			return errs.NewValidationError("id", "can't delete user with equipment or studio invocations")
		}

		err = s.userRep.Delete(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

package userservice

import (
	"context"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	ID              uuid.UUID   `validate:"required"`
	FullName        *string     `validate:"omitempty"`
	Email           *string     `validate:"omitempty,email"`
	Nice            *int        `validate:"omitempty,min=0,max=100"`
	IsAdmin         *bool       `validate:"omitempty"`
	OrganizationIDs []uuid.UUID `validate:"omitempty"`
	Departments     []UpdateUserDepartment
}

type UpdateUserDepartment struct {
	ID   uuid.UUID               `validate:"required"`
	Role domain.RoleInDepartment `validate:"required"`
}

type UpdateNiceRequest struct {
	ID   uuid.UUID `validate:"required"`
	Nice int       `validate:"required,min=0,max=100"`
}

type UpdateUserService interface {
	Update(ctx context.Context, req *UpdateUserRequest) (*domain.User, error)
	UpdateNice(ctx context.Context, req *UpdateNiceRequest) (*domain.User, error)
}

type updateUserService struct {
	userRepo domain.UserRepository
	txm      txmanager.TxManager
	auther   authzservice.AuthZ
}

func NewUpdateUserService(
	userRepo domain.UserRepository,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) UpdateUserService {
	return &updateUserService{
		userRepo: userRepo,
		txm:      txm,
		auther:   auther,
	}
}

var _ UpdateUserService = (*updateUserService)(nil)

func (s *updateUserService) Update(ctx context.Context, req *UpdateUserRequest) (*domain.User, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateUserRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	var user *domain.User
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		user, err = s.userRepo.Get(ctx, req.ID, domain.UserWithDepartments(), domain.UserWithOrganizations())
		if err != nil {
			return err
		}

		if req.FullName != nil {
			user.FullName = *req.FullName
		}
		if req.Email != nil {
			user.Email = *req.Email
		}
		if req.Nice != nil {
			user.Nice = *req.Nice
		}
		if req.IsAdmin != nil {
			user.IsAdmin = *req.IsAdmin
		}
		if req.OrganizationIDs != nil {
			user.Organizations = make([]*domain.Organization, 0, len(req.OrganizationIDs))
			for _, id := range req.OrganizationIDs {
				user.Organizations = append(user.Organizations, &domain.Organization{
					ID: id,
				})
			}
		}

		if req.Departments != nil {
			user.Departments = make([]*domain.UserDepartment, 0, len(req.Departments))
			for _, d := range req.Departments {
				user.Departments = append(user.Departments, &domain.UserDepartment{
					UserID:       req.ID,
					DepartmentID: d.ID,
					Role:         d.Role,
				})
			}
		}

		err = s.userRepo.Update(ctx, user)
		if err != nil {
			return err
		}

		err = s.userRepo.Reload(ctx, user, domain.UserWithDepartments(), domain.UserWithOrganizations())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *updateUserService) UpdateNice(ctx context.Context, req *UpdateNiceRequest) (*domain.User, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateNiceRequest", err.Error())
	}

	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return nil, errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	var user *domain.User
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		user, err = s.userRepo.Get(ctx, req.ID, domain.UserWithDepartments(), domain.UserWithOrganizations())
		if err != nil {
			return err
		}

		user.Nice = req.Nice

		err = s.userRepo.Update(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

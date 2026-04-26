package userservice

import (
	"context"
	"fmt"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"
	"media-equipment-tracker/pkg/txmanager"
	"slices"
)

type MeUserService interface {
	Me(ctx context.Context) (*domain.User, error)
	Update(ctx context.Context, req *UpdateSelfRequest) (*domain.User, error)
	Invocations(ctx context.Context, req *MyInvocationsRequest) ([]*domain.EquipmentInvocation, error)
}

type meUserService struct {
	userRepo domain.UserRepository
	authz    authzservice.AuthZ
	txm      txmanager.TxManager
}

func NewMeUserService(
	userRepo domain.UserRepository,
	authz authzservice.AuthZ,
	txm txmanager.TxManager,
) MeUserService {
	return &meUserService{
		userRepo: userRepo,
		authz:    authz,
		txm:      txm,
	}
}

var _ MeUserService = (*meUserService)(nil)

func (s *meUserService) Me(ctx context.Context) (*domain.User, error) {
	payload, err := s.authz.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Get(ctx, payload.UserID, domain.UserWithDepartments(), domain.UserWithOrganizations())
}

type UpdateSelfRequest struct {
	FullName *string `validate:"omitempty"`
	Email    *string `validate:"omitempty,email"`
}

func (s *meUserService) Update(ctx context.Context, req *UpdateSelfRequest) (*domain.User, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("UpdateSelfRequest", err.Error())
	}

	payload, err := s.authz.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var user *domain.User
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		user, err = s.userRepo.Get(ctx, payload.UserID, domain.UserWithDepartments(), domain.UserWithOrganizations())
		if err != nil {
			return err
		}

		if req.FullName != nil {
			user.FullName = *req.FullName
		}
		if req.Email != nil {
			user.Email = *req.Email
		}

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

type InvocationType string

const (
	Studio    InvocationType = "studio"
	Equipment InvocationType = "equipment"
	All       InvocationType = "all"
)

type MyInvocationsRequest struct {
	StudioStatuses    []domain.StudioInvocationStatus    `validate:"omitempty"`
	EquipmentStatuses []domain.EquipmentInvocationStatus `validate:"omitempty"`
	Type              InvocationType                     `validate:"required"`
}

type MyInvocationsResponse struct {
	EquipmentInvocations      []*domain.EquipmentInvocation
	StudioInvocations         []*domain.StudioInvocation
	AdminEquipmentInvocations []*domain.EquipmentInvocation
	AdminStudioInvocations    []*domain.StudioInvocation
}

func (s *meUserService) Invocations(ctx context.Context, req *MyInvocationsRequest) ([]*domain.EquipmentInvocation, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("MyInvocationsRequest", err.Error())
	}

	payload, err := s.authz.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	opts := []domain.UserOption{}
	switch req.Type {
	case Studio:
		opts = append(opts, domain.UserWithStudioInvocations(), domain.UserWithAdminStudioInvocations())
	case Equipment:
		opts = append(opts, domain.UserWithEquipmentInvocations(), domain.UserWithAdminEquipmentInvocations())
	case All:
		opts = append(opts, domain.UserWithStudioInvocations(), domain.UserWithEquipmentInvocations(), domain.UserWithAdminStudioInvocations(), domain.UserWithAdminEquipmentInvocations())
	default:
		return nil, errs.NewValidationError("req.Type", fmt.Sprintf("invalid type: %s", req.Type))
	}

	user, err := s.userRepo.Get(ctx, payload.UserID, opts...)
	if err != nil {
		return nil, err
	}

	resp := &MyInvocationsResponse{
		EquipmentInvocations:      make([]*domain.EquipmentInvocation, 0, len(user.EquipmentInvocations)),
		StudioInvocations:         make([]*domain.StudioInvocation, 0, len(user.StudioInvocations)),
		AdminEquipmentInvocations: make([]*domain.EquipmentInvocation, 0, len(user.AdminEquipmentInvocations)),
		AdminStudioInvocations:    make([]*domain.StudioInvocation, 0, len(user.AdminStudioInvocations)),
	}

	if req.Type == Equipment || req.Type == All {
		for _, inv := range user.EquipmentInvocations {
			if len(req.EquipmentStatuses) > 0 && slices.Contains(req.EquipmentStatuses, inv.Status) {
				resp.EquipmentInvocations = append(resp.EquipmentInvocations, inv)
			}
		}
		for _, inv := range user.AdminEquipmentInvocations {
			if len(req.EquipmentStatuses) > 0 && slices.Contains(req.EquipmentStatuses, inv.Status) {
				resp.AdminEquipmentInvocations = append(resp.AdminEquipmentInvocations, inv)
			}
		}
	}
	if req.Type == Studio || req.Type == All {
		for _, inv := range user.StudioInvocations {
			if len(req.StudioStatuses) > 0 && slices.Contains(req.StudioStatuses, inv.Status) {
				resp.StudioInvocations = append(resp.StudioInvocations, inv)
			}
		}
		for _, inv := range user.AdminStudioInvocations {
			if len(req.StudioStatuses) > 0 && slices.Contains(req.StudioStatuses, inv.Status) {
				resp.AdminStudioInvocations = append(resp.AdminStudioInvocations, inv)
			}
		}
	}

	return resp.EquipmentInvocations, nil
}

package studioservice

import (
	"context"
	"fmt"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type CuratorStudioService interface {
	Become(ctx context.Context, id uuid.UUID) error
	Approve(ctx context.Context, id uuid.UUID) error
	RequestChanges(ctx context.Context, id uuid.UUID) error
	Complete(ctx context.Context, id uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID) error
}

type curatorStudioService struct {
	studioRepo domain.StudioInvocationRepository
	txm        txmanager.TxManager
	auther     authzservice.AuthZ
}

func NewCuratorStudioService(
	studioRepo domain.StudioInvocationRepository,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) CuratorStudioService {
	return &curatorStudioService{
		studioRepo: studioRepo,
		txm:        txm,
		auther:     auther,
	}
}

var _ CuratorStudioService = (*curatorStudioService)(nil)

func (s *curatorStudioService) Become(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		inv, err := s.studioRepo.Get(ctx, id, domain.StudioInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.Status == domain.StudioCompleted || inv.Status == domain.StudioCancelled {
			return errs.NewValidationError("Status", "can't become curator when studio is completed or cancelled")
		}

		if inv.AdminID != nil {
			return errs.NewValidationError("EquipmentInvocation", "Invocation already has admin")
		}

		inv.AdminID = &payload.UserID
		inv.Status = domain.StudioUnderReview

		err = s.studioRepo.Update(ctx, inv)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *curatorStudioService) Approve(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.StudioInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.studioRepo.Get(ctx, id, domain.StudioInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}
		if inv.Status != domain.StudioUnderReview {
			return errs.NewValidationError("Status", fmt.Sprintf("can't approve studio in status %s", inv.Status))
		}

		inv.Status = domain.StudioApproved

		err = s.studioRepo.Update(ctx, inv)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *curatorStudioService) RequestChanges(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.StudioInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.studioRepo.Get(ctx, id, domain.StudioInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status != domain.StudioUnderReview && inv.Status != domain.StudioApproved {
			return errs.NewValidationError("Status", fmt.Sprintf("can't request changes in status %s", inv.Status))
		}

		inv.Status = domain.StudioChangesRequired

		err = s.studioRepo.Update(ctx, inv)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *curatorStudioService) Complete(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.StudioInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.studioRepo.Get(ctx, id, domain.StudioInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status != domain.StudioApproved {
			return errs.NewValidationError("Status", fmt.Sprintf("can't complete studio in status %s", inv.Status))
		}

		inv.Status = domain.StudioCompleted

		err = s.studioRepo.Update(ctx, inv)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *curatorStudioService) Cancel(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.StudioInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.studioRepo.Get(ctx, id, domain.StudioInvocationWithAdmin(), domain.StudioInvocationWithUser())
		if err != nil {
			return err
		}

		if inv.UserID != payload.UserID && (inv.AdminID == nil || *inv.AdminID != payload.UserID) {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status == domain.StudioCompleted || inv.Status == domain.StudioCancelled {
			return errs.NewValidationError("Status", "can't cancel studio when it is completed or cancelled")
		}

		inv.Status = domain.StudioCancelled

		err = s.studioRepo.Update(ctx, inv)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

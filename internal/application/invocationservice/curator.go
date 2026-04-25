package invocationservice

import (
	"context"
	"fmt"
	"slices"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/google/uuid"
)

type CuratorInvocationService interface {
	Become(ctx context.Context, id uuid.UUID) error
	Approve(ctx context.Context, id uuid.UUID) error
	RequestReview(ctx context.Context, id uuid.UUID) error
	Complete(ctx context.Context, id uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID) error
	Issue(ctx context.Context, id, equipmentID uuid.UUID) error
	Return(ctx context.Context, id, equipmentID uuid.UUID) error
}

type curatorInvocationService struct {
	invocationRepo         domain.EquipmentInvocationRepository
	updateEquipmentService equipmentservice.UpdateEquipmentService
	txm                    txmanager.TxManager
	auther                 authzservice.AuthZ
}

func NewCuratorInvocationService(
	invocationRepo domain.EquipmentInvocationRepository,
	updateEquipmentService equipmentservice.UpdateEquipmentService,
	txm txmanager.TxManager,
	auther authzservice.AuthZ,
) CuratorInvocationService {
	return &curatorInvocationService{
		invocationRepo:         invocationRepo,
		updateEquipmentService: updateEquipmentService,
		txm:                    txm,
		auther:                 auther,
	}
}

var _ CuratorInvocationService = (*curatorInvocationService)(nil)

func (s *curatorInvocationService) Become(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}

	if !slices.Contains(payload.Roles, domain.AdminRole) {
		return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
	}

	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		inv, err := s.invocationRepo.Get(ctx, id)
		if err != nil {
			return err
		}

		inv.AdminID = &payload.UserID
		inv.Status = domain.InvocationUnderReview

		err = s.invocationRepo.Update(ctx, inv)
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

func (s *curatorInvocationService) Approve(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, id, domain.EquipmentInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		inv.Status = domain.InvocationApproved

		err = s.invocationRepo.Update(ctx, inv)
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

func (s *curatorInvocationService) RequestReview(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, id, domain.EquipmentInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		inv.Status = domain.InvocationUnderReview

		err = s.invocationRepo.Update(ctx, inv)
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

func (s *curatorInvocationService) Issue(ctx context.Context, id, equipmentID uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, id, domain.EquipmentInvocationWithAdmin(), domain.EquipmentInvocationWithEquipment())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status != domain.InvocationEquipmentIssued && inv.Status != domain.InvocationApproved {
			return errs.NewValidationError("Status", fmt.Sprintf("can't issue equipment in status %s", inv.Status))
		}

		var eqToIssue *domain.EquipmentInInvocation
		for _, eq := range inv.Equipment {
			if eq.EquipmentID == equipmentID {
				eqToIssue = eq
				break
			}
		}

		if eqToIssue == nil {
			return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s not found in invocation %s", equipmentID, id))
		}

		if eqToIssue.Status != domain.EquipmentNotIssued {
			return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s is already issued (status %s)", eqToIssue.EquipmentID, eqToIssue.Status))
		}

		if eqToIssue.Equipment.Status != domain.EquipmentStatusAvailable {
			return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s is not available (status %s)", eqToIssue.EquipmentID, eqToIssue.Equipment.Status))
		}

		inv.Status = domain.InvocationEquipmentIssued
		eqToIssue.Status = domain.EquipmentIssued
		eqToIssue.Equipment.Status = domain.EquipmentStatusIssued

		_, err = s.updateEquipmentService.UpdateEquipment(ctx, &equipmentservice.UpdateEquipmentRequest{
			ID:     eqToIssue.EquipmentID,
			Status: &eqToIssue.Equipment.Status,
		})
		if err != nil {
			return err
		}

		err = s.invocationRepo.Update(ctx, inv)
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

func (s *curatorInvocationService) Return(ctx context.Context, id, equipmentID uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}

	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, id, domain.EquipmentInvocationWithAdmin(), domain.EquipmentInvocationWithEquipment())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status != domain.InvocationEquipmentIssued {
			return errs.NewValidationError("Status", fmt.Sprintf("can't return equipment in status %s", inv.Status))
		}

		var eqToReturn *domain.EquipmentInInvocation
		for _, eq := range inv.Equipment {
			if eq.EquipmentID == equipmentID {
				eqToReturn = eq
				break
			}
		}

		if eqToReturn == nil {
			return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s not found in invocation %s", equipmentID, id))
		}

		if eqToReturn.Status != domain.EquipmentIssued {
			return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s is not issued (status %s)", eqToReturn.EquipmentID, eqToReturn.Status))
		}

		if eqToReturn.Equipment.Status != domain.EquipmentStatusIssued {
			return errs.NewValidationError("Equipment", fmt.Sprintf("equipment %s is not issued (status %s)", eqToReturn.EquipmentID, eqToReturn.Equipment.Status))
		}

		eqToReturn.Status = domain.EquipmentReturned
		eqToReturn.Equipment.Status = domain.EquipmentStatusAvailable

		var allReturned = true
		for _, eq := range inv.Equipment {
			if eq.Status != domain.EquipmentReturned {
				allReturned = false
				break
			}
		}
		if allReturned {
			inv.Status = domain.InvocationEquipmentReturned
		}
		_, err = s.updateEquipmentService.UpdateEquipment(ctx, &equipmentservice.UpdateEquipmentRequest{
			ID:     eqToReturn.EquipmentID,
			Status: &eqToReturn.Equipment.Status,
		})
		if err != nil {
			return err
		}

		err = s.invocationRepo.Update(ctx, inv)
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

func (s *curatorInvocationService) Complete(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, id, domain.EquipmentInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status != domain.InvocationEquipmentReturned {
			return errs.NewValidationError("Status", fmt.Sprintf("can't complete invocation in status %s", inv.Status))
		}

		inv.Status = domain.InvocationCompleted

		err = s.invocationRepo.Update(ctx, inv)
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

func (s *curatorInvocationService) Cancel(ctx context.Context, id uuid.UUID) error {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	var inv *domain.EquipmentInvocation
	err = s.txm.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		inv, err = s.invocationRepo.Get(ctx, id, domain.EquipmentInvocationWithAdmin())
		if err != nil {
			return err
		}

		if inv.AdminID == nil || *inv.AdminID != payload.UserID {
			return errs.NewRoleAuthError([]domain.RoleAuth{domain.AdminRole}, payload.Roles)
		}

		if inv.Status == domain.InvocationEquipmentIssued {
			return errs.NewValidationError("Status", "can't cancel invocation when equipment is issued")
		}

		inv.Status = domain.InvocationCancelled

		err = s.invocationRepo.Update(ctx, inv)
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

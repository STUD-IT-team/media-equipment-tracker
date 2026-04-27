package accessservice

import (
	"context"
	"fmt"

	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"

	"github.com/google/uuid"
)

type HaveEquipmentAccessRequest struct {
	EquipmentID  uuid.UUID
	Department   *domain.Department
	Organization *domain.Organization
}

type HaveStudioAccessRequest struct {
	Department   *domain.Department
	Organization *domain.Organization
}

type AccessService interface {
	HaveAccessToEquipment(ctx context.Context, req *HaveEquipmentAccessRequest) (bool, error)
	HaveAccessToStudio(ctx context.Context, req *HaveStudioAccessRequest) (bool, error)
}

type accessService struct {
	auther authzservice.AuthZ
}

var _ AccessService = (*accessService)(nil)

func NewAccessService(
	auther authzservice.AuthZ,
) AccessService {
	return &accessService{
		auther: auther,
	}
}

func (s *accessService) HaveAccessToEquipment(ctx context.Context, req *HaveEquipmentAccessRequest) (bool, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return false, err
	}
	userID := payload.UserID
	eqID := req.EquipmentID
	department := req.Department
	organization := req.Organization

	if req.Department != nil {
		if department.Users == nil {
			return false, errs.NewValidationError("HaveEquipmentAccessRequest", "department's users must be loaded")
		}
		if department.Equipment == nil {
			return false, errs.NewValidationError("HaveEquipmentAccessRequest", "department's equipment must be loaded")
		}
		var userDep *domain.UserDepartment
		for _, u := range department.Users {
			if u.UserID == userID {
				userDep = u
				break
			}
		}
		if userDep == nil {
			return false, errs.NewEquipmentAccessError("user not in specified department")
		}

		var equipment *domain.Equipment
		for _, e := range department.Equipment {
			if e.ID == eqID {
				equipment = e
				break
			}
		}

		// Отделы имеют доступ только к своему оборудованию
		if equipment == nil {
			return false, errs.NewEquipmentAccessError("department do not own specified equipment")
		}

		// Стажеры только к стажёрскому оборудованию
		if userDep.Role == domain.RoleTrainee && !equipment.AvailableToTrainee {
			return false, errs.NewEquipmentAccessError("user do not have access to activist's equipment (role is trainee)")
		}

		// Если активсит не добри, то он может получить только стажёрское оборудование
		if userDep.Role == domain.RoleActivist && !equipment.AvailableToTrainee && userDep.User.Nice < domain.EquipmentNiceThreshold {
			return false, errs.NewEquipmentAccessError(fmt.Sprintf("user do not have access to activist's equipment (nice < %d)", domain.EquipmentNiceThreshold))
		}

		return true, nil
	} else if req.Organization != nil {
		if organization.Users == nil {
			return false, errs.NewValidationError("HaveEquipmentAccessRequest", "organization's users must be loaded")
		}

		var user *domain.User
		for _, u := range organization.Users {
			if u.ID == userID {
				user = u
				break
			}
		}

		if user == nil {
			return false, errs.NewEquipmentAccessError("user not in specified organization")
		}

		// Организация имеет доступ к любому оборудовнию

		if user.Nice < domain.EquipmentNiceThreshold {
			return false, errs.NewEquipmentAccessError(fmt.Sprintf("user do not have access to equipment (nice < %d)", domain.EquipmentNiceThreshold))
		}

		return true, nil
	}
	return false, errs.NewValidationError("HaveEquipmentAccessRequest", "need either department or organization specified")
}

func (s *accessService) HaveAccessToStudio(ctx context.Context, req *HaveStudioAccessRequest) (bool, error) {
	payload, err := s.auther.TokenPayloadFromContext(ctx)
	if err != nil {
		return false, err
	}
	userID := payload.UserID

	if req.Department != nil {
		if req.Department.Users == nil {
			return false, errs.NewValidationError("HaveStudioAccessRequest", "department's users must be loaded")
		}

		var userDep *domain.UserDepartment
		for _, u := range req.Department.Users {
			if u.UserID == userID {
				userDep = u
				break
			}
		}

		if userDep == nil {
			return false, errs.NewStudioAccessError("user not in specified department")
		}

		if userDep.Role == domain.RoleTrainee {
			return false, errs.NewStudioAccessError("user do not have access to studio (role is trainee)")
		}

		if userDep.User.Nice < domain.StudioNiceThreshold {
			return false, errs.NewStudioAccessError(fmt.Sprintf("user do not have access to studio (nice < %d)", domain.StudioNiceThreshold))
		}

		return true, nil
	} else if req.Organization != nil {
		if req.Organization.Users == nil {
			return false, errs.NewValidationError("HaveStudioAccessRequest", "organization's users must be loaded")
		}

		var user *domain.User
		for _, u := range req.Organization.Users {
			if u.ID == userID {
				user = u
				break
			}
		}

		if user.Nice < domain.StudioNiceThreshold {
			return false, errs.NewStudioAccessError(fmt.Sprintf("user do not have access to studio (nice < %d)", domain.StudioNiceThreshold))
		}

		return true, nil
	}

	return false, errs.NewValidationError("HaveStudioAccessRequest", "need either department or organization specified")
}

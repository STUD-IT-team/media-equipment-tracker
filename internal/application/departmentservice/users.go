package departmentservice

import (
	"context"

	"media-equipment-tracker/internal/domain"

	"github.com/google/uuid"
)

type GetDepartmentUsersService interface {
	GetDepartmentUsers(ctx context.Context, id uuid.UUID, role *domain.RoleInDepartment) ([]*domain.UserDepartment, error)
}

type getDepartmentUsersService struct {
	departmentRepository domain.DepartmentRepository
}

var _ GetDepartmentUsersService = (*getDepartmentUsersService)(nil)

func NewGetDepartmentUsersService(
	departmentRepository domain.DepartmentRepository,
) GetDepartmentUsersService {
	return &getDepartmentUsersService{
		departmentRepository: departmentRepository,
	}
}

func (s *getDepartmentUsersService) GetDepartmentUsers(ctx context.Context, id uuid.UUID, role *domain.RoleInDepartment) ([]*domain.UserDepartment, error) {
	department, err := s.departmentRepository.Get(ctx, id, domain.DepartmentWithUsers())
	if err != nil {
		return nil, err
	}

	users := department.Users
	if role != nil {
		filtered := make([]*domain.UserDepartment, 0)
		for _, user := range users {
			if user.Role == *role {
				filtered = append(filtered, user)
			}
		}
		return filtered, nil
	}

	return users, nil
}

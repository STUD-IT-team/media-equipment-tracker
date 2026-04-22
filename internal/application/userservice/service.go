package userservice

import (
	"context"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"

	"github.com/google/uuid"
)

type UserService interface {
	GetById(ctx context.Context, id uuid.UUID, with ...domain.UserOption) (*domain.User, error)
	GetCurrent(ctx context.Context, with ...domain.UserOption) (*domain.User, error)
	GetAll(ctx context.Context, with ...domain.UserOption) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRep domain.UserRepository
	authz   authzservice.AuthZ
}

func NewUserService(userRep domain.UserRepository, authz authzservice.AuthZ) (UserService, error) {
	service := &userService{
		userRep: userRep,
		authz:   authz,
	}
	return service, nil
}

func (s *userService) GetById(ctx context.Context, id uuid.UUID, with ...domain.UserOption) (*domain.User, error) {
	user, err := s.userRep.Get(ctx, id, with...)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetCurrent(ctx context.Context, with ...domain.UserOption) (*domain.User, error) {
	tokenPayload, err := s.authz.TokenPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetById(ctx, tokenPayload.UserID, with...)
}

func (s *userService) GetAll(ctx context.Context, with ...domain.UserOption) ([]*domain.User, error) {
	users, err := s.userRep.List(ctx, with...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := s.userRep.Update(ctx, user); err != nil {
		return nil, err
	}
	updatedUser, err := s.userRep.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *userService) DeleteById(ctx context.Context, id uuid.UUID) error {
	return s.userRep.Delete(ctx, id)
}

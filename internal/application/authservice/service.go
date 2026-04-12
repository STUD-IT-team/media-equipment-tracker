package authservice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"media-equipment-tracker/internal/adapters/inmem"
	"media-equipment-tracker/internal/domain"
)

var (
	ErrCreateAuthUserService = errors.New("error create auth user service")
)

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,min=4,max=50" example:"test@mail.ru"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token" example:"7249ede8-5083-4bd6-ad09-0b5fa3c5f2de"`
}

type RegisterUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=4,max=50" example:"Ivan Ivanov Ivanovich"`
	Email    string `json:"email" binding:"required,min=4,max=50" example:"test@mail.ru"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type AuthUserService interface {
	LoginUser(ctx context.Context, lur LoginUserRequest) (string, error)
	RegisterUser(ctx context.Context, rur RegisterUserRequest) error
	VerifyByToken(token string, needRoles []domain.RoleAuth) (*domain.TokenPayload, error)
	LogoutUser(ctx context.Context, token string)
}

type authUserService struct {
	tokenMaker          TokenMaker
	hasher              Hasher
	accessTokenDuration time.Duration
	userRep             domain.UserRepository
	tokenRep            inmem.TokenRepository
}

func NewAuthUser(
	tokenMaker TokenMaker, hasher Hasher, accessTokenDuration time.Duration, urep domain.UserRepository,
	tokenRep inmem.TokenRepository,
) (AuthUserService, error) {
	if tokenMaker == nil || hasher == nil || urep == nil || accessTokenDuration <= 0 {
		return nil, ErrCreateAuthUserService
	}
	server := &authUserService{
		tokenMaker:          tokenMaker,
		hasher:              hasher,
		accessTokenDuration: accessTokenDuration,
		userRep:             urep,
		tokenRep:            tokenRep,
	}
	return server, nil
}

func (s *authUserService) LoginUser(_ context.Context, lur LoginUserRequest) (string, error) {
	user, err := s.userRep.GetByEmail(lur.Email)
	if err != nil {
		return "", err
	}
	if err := s.hasher.CheckPassword(lur.Password, user.HashPassword); err != nil {
		return "", err
	}
	roles := make([]domain.RoleAuth, 0)
	if user.IsAdmin {
		roles = append(roles, domain.AdminRole)
	}
	accessToken, err := s.tokenMaker.CreateToken(user.ID, roles, s.accessTokenDuration)
	if err != nil {
		return "", err
	}
	if ok := s.tokenRep.Add(accessToken); !ok {
		return "", errors.New("failed to add access token")
	}
	return accessToken, nil
}

func (s *authUserService) RegisterUser(_ context.Context, rur RegisterUserRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return err
	}
	return s.userRep.Create(&domain.User{
		ID:           uuid.New(),
		FullName:     rur.FullName,
		Email:        rur.Email,
		HashPassword: hashedPassword,
		Nice:         domain.DefaultNice,
		IsAdmin:      false,
	})
}

func (s *authUserService) VerifyByToken(tokenStr string, needRoles []domain.RoleAuth) (*domain.TokenPayload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, needRoles)
}

func (s *authUserService) LogoutUser(_ context.Context, token string) {
	s.tokenRep.Delete(token)
}

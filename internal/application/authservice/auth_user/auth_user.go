package auth_user

import (
	"context"
	"media-equipment-tracker/internal/application/authservice/hasher"
	tokenmaker "media-equipment-tracker/internal/application/authservice/token_maker"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/handlers/auth-api/dto"
	"time"

	"github.com/google/uuid"
)

type AuthUserService interface {
	LoginUser(ctx context.Context, lur dto.LoginUserRequest) (string, error)
	RegisterUser(ctx context.Context, rur dto.RegisterUserRequest) error
	VerifyByToken(token string) (*tokenmaker.Payload, error)
}

type authUserService struct {
	tokenMaker          tokenmaker.TokenMaker
	hasher              hasher.Hasher
	accessTokenDuration time.Duration
	userrep             domain.UserRepository
}

func NewAuthUser(
	tokenMaker tokenmaker.TokenMaker, hasher hasher.Hasher, accessTokenDuration time.Duration, urep domain.UserRepository,
) AuthUserService {
	server := &authUserService{
		tokenMaker:          tokenMaker,
		hasher:              hasher,
		accessTokenDuration: accessTokenDuration,
		userrep:             urep,
	}
	return server
}

func (s *authUserService) LoginUser(ctx context.Context, lur dto.LoginUserRequest) (string, error) {
	user, err := s.userrep.GetByEmail(lur.Email)
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(lur.Password, user.HashPassword)
	if err != nil {
		return "", err
	}
	accessToken, err := s.tokenMaker.CreateToken(
		user.ID,
		tokenmaker.UserRole,
		s.accessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authUserService) RegisterUser(ctx context.Context, rur dto.RegisterUserRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return err
	}
	return s.userrep.Create(&domain.User{
		ID:           uuid.New(),
		FullName:     rur.FullName,
		Email:        rur.Email,
		HashPassword: hashedPassword,
		Nice:         domain.DefaultNice,
		IsAdmin:      false,
	})
}

func (s *authUserService) VerifyByToken(tokenStr string) (*tokenmaker.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, tokenmaker.UserRole)
}

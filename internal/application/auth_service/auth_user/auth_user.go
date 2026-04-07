package auth_user

import (
	"context"
	"errors"
	"media-equipment-tracker/internal/adapters/in_mem"
	"media-equipment-tracker/internal/application/auth_service/hasher"
	tokenmaker "media-equipment-tracker/internal/application/auth_service/token_maker"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/handlers/auth-api/dto"
	"time"

	"github.com/google/uuid"
)

type AuthUserService interface {
	LoginUser(ctx context.Context, lur dto.LoginUserRequest) (string, error)
	RegisterUser(ctx context.Context, rur dto.RegisterUserRequest) error
	VerifyByToken(token string, needRoles []domain.RoleAuth) (*domain.TokenPayload, error)
	LogoutUser(ctx context.Context, token string)
}

type authUserService struct {
	tokenMaker          tokenmaker.TokenMaker
	hasher              hasher.Hasher
	accessTokenDuration time.Duration
	userrep             domain.UserRepository
	tokenRep            in_mem.TokenRepository
}

func NewAuthUser(
	tokenMaker tokenmaker.TokenMaker, hasher hasher.Hasher, accessTokenDuration time.Duration, urep domain.UserRepository,
) AuthUserService {
	server := &authUserService{
		tokenMaker:          tokenMaker,
		hasher:              hasher,
		accessTokenDuration: accessTokenDuration,
		userrep:             urep,
		tokenRep:            in_mem.GetTokenRepository(),
	}
	return server
}

func (s *authUserService) LoginUser(ctx context.Context, lur dto.LoginUserRequest) (string, error) {
	user, err := s.userrep.GetByEmail(lur.Email)
	if err != nil {
		return "", err
	}
	if err = s.hasher.CheckPassword(lur.Password, user.HashPassword); err != nil {
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

func (s *authUserService) VerifyByToken(tokenStr string, needRoles []domain.RoleAuth) (*domain.TokenPayload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, needRoles)
}

func (s *authUserService) LogoutUser(ctx context.Context, token string) {
	s.tokenRep.Delete(token)
}

package tokenmaker

import (
	"errors"
	"media-equipment-tracker/internal/domain"
	"time"

	"github.com/google/uuid"
)

type TokenMaker interface {
	// duration - корректный срок действия, return - одписанную строку токена или ошибку
	CreateToken(id uuid.UUID, roles []domain.RoleAuth, duration time.Duration) (string, error)
	VerifyToken(token string, roles []domain.RoleAuth) (*domain.TokenPayload, error)
}

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrExpiredToken  = errors.New("token has expired")
	ErrIncorrectRole = errors.New("incorrect role")
)

func NewTokenMaker(symmetricKey string) (TokenMaker, error) {
	return NewJWTMaker(symmetricKey)
}

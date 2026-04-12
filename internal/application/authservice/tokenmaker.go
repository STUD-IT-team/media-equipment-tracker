package authservice

import (
	"errors"

	"media-equipment-tracker/internal/domain"
)

type TokenMaker interface {
	CreateToken(*domain.TokenPayload) (string, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
}

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrExpiredToken  = errors.New("token has expired")
	ErrIncorrectRole = errors.New("incorrect role")
)

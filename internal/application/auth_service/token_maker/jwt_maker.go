package tokenmaker

import (
	"errors"
	"fmt"
	"media-equipment-tracker/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// алгоритм с симметричным ключом для подписи токенов

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32 // TODO to config

func NewJWTMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(userID uuid.UUID, roles []domain.RoleAuth, duration time.Duration) (string, error) {
	payload, err := domain.NewTokenPayload(userID, roles, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string, roles []domain.RoleAuth) (*domain.TokenPayload, error) {
	// проверить заголовок токена, и убедиться, что алгоритм подписи соответствует тому, который используется для подписи токенов.
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC) // потому что используем HS256
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &domain.TokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*domain.TokenPayload)
	if !ok {
		return nil, ErrInvalidToken
	}
	for _, role := range roles {
		for _, expectedRole := range payload.Roles {
			if expectedRole == role {
				return nil, ErrIncorrectRole
			}
		}
	}
	return payload, nil
}

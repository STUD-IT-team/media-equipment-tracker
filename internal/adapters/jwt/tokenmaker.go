package tokenmaker

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	"media-equipment-tracker/internal/application/authservice"
	"media-equipment-tracker/internal/domain"
)

// алгоритм с симметричным ключом для подписи токенов

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32 // TODO to config

func NewJWTMaker(secretKey string) (authservice.TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(payload *domain.TokenPayload) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*domain.TokenPayload, error) {
	// проверить заголовок токена, и убедиться, что алгоритм подписи соответствует тому, который используется для подписи токенов.
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC) // потому что используем HS256
		if !ok {
			return nil, authservice.ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &domain.TokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, authservice.ErrExpiredToken) {
			return nil, authservice.ErrExpiredToken
		}
		return nil, authservice.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*domain.TokenPayload)
	if !ok {
		return nil, authservice.ErrInvalidToken
	}
	return payload, nil
}

package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RoleAuth string

const (
	AdminRole RoleAuth = "ADMIN"
)

var ErrExpiredToken = errors.New("token has expired")

type TokenPayload struct {
	UserID    uuid.UUID  `json:"user_id"`
	Roles     []RoleAuth `json:"role"`
	ExpiredAt time.Time  `json:"expired_at"`
}

func NewTokenPayload(personID uuid.UUID, roles []RoleAuth, duration time.Duration) (*TokenPayload, error) {
	payload := &TokenPayload{
		UserID:    personID,
		Roles:     roles,
		ExpiredAt: time.Now().UTC().Add(duration),
	}
	return payload, nil
}

func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

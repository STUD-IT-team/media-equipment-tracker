package authzservice

import (
	"context"
	"errors"
	"media-equipment-tracker/internal/domain"
)

type authZContextKey int

const (
	AuthZContextKey authZContextKey = iota
)

var (
	ErrNotAuthZ = errors.New("not authorized")
)

type AuthZ interface {
	Authorize(ctx context.Context, payload domain.TokenPayload) context.Context
	TokenPayloadFromContext(ctx context.Context) (domain.TokenPayload, error)
}

func NewAuthZ() AuthZ {
	return &authZ{}
}

type authZ struct {
}

func (a *authZ) Authorize(ctx context.Context, payload domain.TokenPayload) context.Context {
	return context.WithValue(ctx, AuthZContextKey, payload)
}

func (a *authZ) TokenPayloadFromContext(ctx context.Context) (domain.TokenPayload, error) {
	payload, ok := ctx.Value(AuthZContextKey).(domain.TokenPayload)
	if !ok {
		return domain.TokenPayload{}, ErrNotAuthZ
	}
	return payload, nil
}

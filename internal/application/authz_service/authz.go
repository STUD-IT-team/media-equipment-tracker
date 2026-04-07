package authz_service

import (
	"context"
	"errors"
	"media-equipment-tracker/internal/domain"
	"sync"
)

var authz AuthZ
var authzMutex sync.Mutex

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

func GetAuthZ() AuthZ {
	authzMutex.Lock()
	defer authzMutex.Unlock()
	if authz != nil {
		return authz
	}
	authz = &authZ{}
	return authz
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

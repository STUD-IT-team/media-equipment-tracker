//go:build unit

package userservice_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"media-equipment-tracker/internal/domain"
)

type UserRepoMock struct{ mock.Mock }

func (m *UserRepoMock) Get(ctx context.Context, id uuid.UUID, with ...domain.UserOption) (*domain.User, error) {
	args := m.Called(ctx, id, with)
	if v := args.Get(0); v != nil {
		return v.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *UserRepoMock) GetByEmail(ctx context.Context, email string, with ...domain.UserOption) (*domain.User, error) {
	args := m.Called(ctx, email, with)
	if v := args.Get(0); v != nil {
		return v.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *UserRepoMock) List(ctx context.Context, with ...domain.UserOption) ([]*domain.User, error) {
	args := m.Called(ctx, with)
	return args.Get(0).([]*domain.User), args.Error(1)
}
func (m *UserRepoMock) Reload(ctx context.Context, user *domain.User, with ...domain.UserOption) error {
	return m.Called(ctx, user, with).Error(0)
}
func (m *UserRepoMock) Create(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *UserRepoMock) Update(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *UserRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type AuthZMock struct{ mock.Mock }

func (m *AuthZMock) Authorize(ctx context.Context, payload domain.TokenPayload) context.Context {
	return m.Called(ctx, payload).Get(0).(context.Context)
}
func (m *AuthZMock) TokenPayloadFromContext(ctx context.Context) (domain.TokenPayload, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.TokenPayload), args.Error(1)
}

type TxManagerMock struct{ mock.Mock }

func (m *TxManagerMock) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	err := fn(ctx)
	if err != nil {
		return err
	}
	return nil
}
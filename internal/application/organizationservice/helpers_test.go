//go:build unit

package organizationservice_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"media-equipment-tracker/internal/domain"
)

type OrganizationRepoMock struct{ mock.Mock }

func (m *OrganizationRepoMock) Get(ctx context.Context, id uuid.UUID, with ...domain.OrganizationOption) (*domain.Organization, error) {
	args := m.Called(ctx, id, with)
	if v := args.Get(0); v != nil {
		return v.(*domain.Organization), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrganizationRepoMock) List(ctx context.Context, with ...domain.OrganizationOption) ([]*domain.Organization, error) {
	args := m.Called(ctx, with)
	if v := args.Get(0); v != nil {
		return v.([]*domain.Organization), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrganizationRepoMock) Reload(ctx context.Context, org *domain.Organization, with ...domain.OrganizationOption) error {
	return m.Called(ctx, org, with).Error(0)
}

func (m *OrganizationRepoMock) Create(ctx context.Context, org *domain.Organization) error {
	return m.Called(ctx, org).Error(0)
}

func (m *OrganizationRepoMock) Update(ctx context.Context, org *domain.Organization) error {
	return m.Called(ctx, org).Error(0)
}

func (m *OrganizationRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
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
	return m.Called(ctx, fn).Error(0)
}

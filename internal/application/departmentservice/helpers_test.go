//go:build unit

package departmentservice_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"media-equipment-tracker/internal/domain"
)

type DepartmentRepoMock struct{ mock.Mock }

func (m *DepartmentRepoMock) Get(ctx context.Context, id uuid.UUID, with ...domain.DepartmentOption) (*domain.Department, error) {
	args := m.Called(ctx, id, with)
	if v := args.Get(0); v != nil {
		return v.(*domain.Department), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *DepartmentRepoMock) List(ctx context.Context, with ...domain.DepartmentOption) ([]*domain.Department, error) {
	args := m.Called(ctx, with)
	if v := args.Get(0); v != nil {
		return v.([]*domain.Department), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *DepartmentRepoMock) Reload(ctx context.Context, dep *domain.Department, with ...domain.DepartmentOption) error {
	return m.Called(ctx, dep, with).Error(0)
}

func (m *DepartmentRepoMock) Create(ctx context.Context, dep *domain.Department) error {
	return m.Called(ctx, dep).Error(0)
}

func (m *DepartmentRepoMock) Update(ctx context.Context, dep *domain.Department) error {
	return m.Called(ctx, dep).Error(0)
}

func (m *DepartmentRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type UserRepoMock struct{ mock.Mock }

func (m *UserRepoMock) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if v := args.Get(0); v != nil {
		return v.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if v := args.Get(0); v != nil {
		return v.([]*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
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

type EquipmentRepoMock struct{ mock.Mock }

func (m *EquipmentRepoMock) Get(ctx context.Context, id uuid.UUID, with ...domain.EquipmentOption) (*domain.Equipment, error) {
	args := m.Called(ctx, id, with)
	if v := args.Get(0); v != nil {
		return v.(*domain.Equipment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *EquipmentRepoMock) GetByInventoryNumber(ctx context.Context, inventoryNumber string, with ...domain.EquipmentOption) (*domain.Equipment, error) {
	args := m.Called(ctx, inventoryNumber, with)
	if v := args.Get(0); v != nil {
		return v.(*domain.Equipment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *EquipmentRepoMock) GetUnoccupied(ctx context.Context, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	return nil, nil
}

func (m *EquipmentRepoMock) List(ctx context.Context, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	return nil, nil
}

func (m *EquipmentRepoMock) Reload(ctx context.Context, equipment *domain.Equipment, with ...domain.EquipmentOption) error {
	return nil
}

func (m *EquipmentRepoMock) Create(ctx context.Context, equipment *domain.Equipment) error {
	return m.Called(ctx, equipment).Error(0)
}

func (m *EquipmentRepoMock) Update(ctx context.Context, equipment *domain.Equipment) error {
	return m.Called(ctx, equipment).Error(0)
}

func (m *EquipmentRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
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
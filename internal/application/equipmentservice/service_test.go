//go:build unit

package equipmentservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/application/invocationservice/invocationsearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
)

type EquipmentSuite struct {
	suite.Suite

	svc                  equipmentservice.EquipmentService
	equipmentRepo        *EquipmentRepoMock
	auther               *AuthZMock
	txManager            *TxManagerMock
	searchRepo           *SearchEquipmentRepoMock
	searchInvocationRepo *SearchInvocationRepoMock
	departmentRepo       *DepartmentRepoMock
}

func (s *EquipmentSuite) SetupTest() {
	s.equipmentRepo = new(EquipmentRepoMock)
	s.auther = new(AuthZMock)
	s.txManager = new(TxManagerMock)
	s.searchRepo = new(SearchEquipmentRepoMock)
	s.searchInvocationRepo = new(SearchInvocationRepoMock)
	s.departmentRepo = new(DepartmentRepoMock)

	s.svc = equipmentservice.NewEquipmentService(
		s.auther,
		s.equipmentRepo,
		s.searchRepo,
		s.searchInvocationRepo,
		s.departmentRepo,
		s.txManager,
	)
}

func newEquipment() *domain.Equipment {
	return &domain.Equipment{
		ID:              uuid.New(),
		InventoryNumber: "INV001",
		Name:            "Test Equipment",
		Category:        "Test",
		Status:          domain.EquipmentStatusAvailable,
	}
}

func (s *EquipmentSuite) TestGet_Success() {
	eq := newEquipment()

	s.equipmentRepo.On("Get", mock.Anything, eq.ID, mock.Anything, mock.Anything, mock.Anything).Return(eq, nil)

	result, err := s.svc.Get(context.Background(), eq.ID)

	s.NoError(err)
	s.Equal(eq, result)
}

func (s *EquipmentSuite) TestGet_Failure() {
	id := uuid.New()
	expectedErr := errors.New("not found")

	s.equipmentRepo.On("Get", mock.Anything, id, mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	_, err := s.svc.Get(context.Background(), id)

	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *EquipmentSuite) TestGetByInventoryNumber_Success() {
	eq := newEquipment()

	s.equipmentRepo.On("GetByInventoryNumber", mock.Anything, eq.InventoryNumber, mock.Anything, mock.Anything, mock.Anything).Return(eq, nil)

	result, err := s.svc.GetByInventoryNumber(context.Background(), eq.InventoryNumber)

	s.NoError(err)
	s.Equal(eq, result)
}

func (s *EquipmentSuite) TestDelete_Success() {
	eq := newEquipment()
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		eq.Invocations = []*domain.EquipmentInInvocation{} // no invocations
		s.equipmentRepo.On("Get", mock.Anything, eq.ID, mock.Anything).Return(eq, nil)
		s.equipmentRepo.On("Delete", mock.Anything, eq.ID).Return(nil)
		err := fn(ctx)
		s.NoError(err)
	})

	err := s.svc.Delete(ctx, eq.ID)

	s.NoError(err)
}

func (s *EquipmentSuite) TestDelete_NoAdmin() {
	eq := newEquipment()
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)

	err := s.svc.Delete(ctx, eq.ID)

	s.Error(err)
	s.True(errs.IsRoleAuthError(err))
}

func (s *EquipmentSuite) TestDelete_HasInvocations() {
	eq := newEquipment()
	payload := domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}

	ctx := context.Background()
	s.auther.On("TokenPayloadFromContext", ctx).Return(payload, nil)
	s.txManager.On("WithinTx", ctx, mock.AnythingOfType("func(context.Context) error")).Return(equipmentservice.ErrEquipmentHasInvocations).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		eq.Invocations = []*domain.EquipmentInInvocation{{}} // has invocations
		s.equipmentRepo.On("Get", mock.Anything, eq.ID, mock.Anything).Return(eq, nil)
		err := fn(ctx)
		s.Error(err)
	})

	err := s.svc.Delete(ctx, eq.ID)

	s.Error(err)
	s.Equal(equipmentservice.ErrEquipmentHasInvocations, err)
}

// Mocks

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

type SearchEquipmentRepoMock struct{ mock.Mock }

func (m *SearchEquipmentRepoMock) Search(ctx context.Context, search *equipmentservice.SearchEquipmentRequest, with ...domain.EquipmentOption) ([]*domain.Equipment, error) {
	args := m.Called(ctx, search, with)
	return args.Get(0).([]*domain.Equipment), args.Error(1)
}

type SearchInvocationRepoMock struct{ mock.Mock }

func (m *SearchInvocationRepoMock) Search(ctx context.Context, search *invocationsearch.SearchInvocationRequest, with ...domain.EquipmentInvocationOption) ([]*domain.EquipmentInvocation, error) {
	args := m.Called(ctx, search, with)
	return args.Get(0).([]*domain.EquipmentInvocation), args.Error(1)
}

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
	return args.Get(0).([]*domain.Department), args.Error(1)
}
func (m *DepartmentRepoMock) Reload(ctx context.Context, dep *domain.Department, with ...domain.DepartmentOption) error {
	return m.Called(ctx, dep, with).Error(0)
}
func (m *DepartmentRepoMock) Create(ctx context.Context, department *domain.Department) error {
	return m.Called(ctx, department).Error(0)
}
func (m *DepartmentRepoMock) Update(ctx context.Context, department *domain.Department) error {
	return m.Called(ctx, department).Error(0)
}
func (m *DepartmentRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func TestEquipmentSuite(t *testing.T) {
	suite.Run(t, new(EquipmentSuite))
}

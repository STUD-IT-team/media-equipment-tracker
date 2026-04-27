//go:build unit

package studioservice_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"media-equipment-tracker/internal/application/accessservice"
	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/domain"
)

type StudioInvocationRepoMock struct{ mock.Mock }

func (m *StudioInvocationRepoMock) Get(ctx context.Context, id uuid.UUID, with ...domain.StudioInvocationOption) (*domain.StudioInvocation, error) {
	args := m.Called(ctx, id, with)
	if v := args.Get(0); v != nil {
		return v.(*domain.StudioInvocation), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *StudioInvocationRepoMock) List(ctx context.Context, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	args := m.Called(ctx, with)
	return args.Get(0).([]*domain.StudioInvocation), args.Error(1)
}
func (m *StudioInvocationRepoMock) Reload(ctx context.Context, studioInvocation *domain.StudioInvocation, with ...domain.StudioInvocationOption) error {
	return m.Called(ctx, studioInvocation, with).Error(0)
}
func (m *StudioInvocationRepoMock) Create(ctx context.Context, studioInvocation *domain.StudioInvocation) error {
	return m.Called(ctx, studioInvocation).Error(0)
}
func (m *StudioInvocationRepoMock) Update(ctx context.Context, studioInvocation *domain.StudioInvocation) error {
	return m.Called(ctx, studioInvocation).Error(0)
}
func (m *StudioInvocationRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type SearchStudioServiceMock struct{ mock.Mock }

func (m *SearchStudioServiceMock) Search(ctx context.Context, search *studiosearch.SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	args := m.Called(ctx, search, with)
	return args.Get(0).([]*domain.StudioInvocation), args.Error(1)
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

// Additional mocks needed for embedded services
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
	return args.Get(0).([]*domain.Organization), args.Error(1)
}
func (m *OrganizationRepoMock) Reload(ctx context.Context, org *domain.Organization, with ...domain.OrganizationOption) error {
	return m.Called(ctx, org, with).Error(0)
}
func (m *OrganizationRepoMock) Create(ctx context.Context, organization *domain.Organization) error {
	return m.Called(ctx, organization).Error(0)
}
func (m *OrganizationRepoMock) Update(ctx context.Context, organization *domain.Organization) error {
	return m.Called(ctx, organization).Error(0)
}
func (m *OrganizationRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type AccessServiceMock struct{ mock.Mock }

func (m *AccessServiceMock) HaveAccessToEquipment(ctx context.Context, req *accessservice.HaveEquipmentAccessRequest) (bool, error) {
	args := m.Called(ctx, req)
	return args.Bool(0), args.Error(1)
}
func (m *AccessServiceMock) HaveAccessToStudio(ctx context.Context, req *accessservice.HaveStudioAccessRequest) (bool, error) {
	args := m.Called(ctx, req)
	return args.Bool(0), args.Error(1)
}

type SearchStudioInvocationRepoMock struct{ mock.Mock }

func (m *SearchStudioInvocationRepoMock) Search(ctx context.Context, search *studiosearch.SearchStudioInvocationRequest, with ...domain.StudioInvocationOption) ([]*domain.StudioInvocation, error) {
	args := m.Called(ctx, search, with)
	return args.Get(0).([]*domain.StudioInvocation), args.Error(1)
}

type AvailabilityStudioServiceMock struct{ mock.Mock }

func (m *AvailabilityStudioServiceMock) Availability(ctx context.Context, req studioservice.AvailabilityStudioInvocationRequest) (*studioservice.StudioInvocationAvailabilityResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*studioservice.StudioInvocationAvailabilityResponse), args.Error(1)
}

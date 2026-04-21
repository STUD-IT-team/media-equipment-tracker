//go:build integration

package pgequipmentinvocation_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgorganization"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
)

type EquipmentInvocationRepositorySuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.EquipmentInvocationRepository
	user domain.UserRepository
	org  domain.OrganizationRepository
}

func (s *EquipmentInvocationRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *EquipmentInvocationRepositorySuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *EquipmentInvocationRepositorySuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))
	s.user = pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
	s.org = pgorganization.NewPostgresOrganizationRepository(gormtx.NewDBGetter(gdb))
}

func (s *EquipmentInvocationRepositorySuite) TestCreate_Get() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	inv := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(ctx, user))
	s.NoError(s.org.Create(ctx, org))

	s.Require().NoError(s.repo.Create(ctx, inv))

	got, err := s.repo.Get(ctx, inv.ID)
	s.NoError(err)
	s.Equal(inv.EventName, got.EventName)
}

func (s *EquipmentInvocationRepositorySuite) TestGet_NotFound() {
	ctx := context.Background()
	_, err := s.repo.Get(ctx, uuid.New())
	s.Error(err)
}

func (s *EquipmentInvocationRepositorySuite) TestList() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	inv1 := newInvocation(&org.ID, nil, user.ID, nil)
	inv2 := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(ctx, user))
	s.NoError(s.org.Create(ctx, org))

	s.NoError(s.repo.Create(ctx, inv1))
	s.NoError(s.repo.Create(ctx, inv2))

	list, err := s.repo.List(ctx)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *EquipmentInvocationRepositorySuite) TestUpdate() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	inv := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(ctx, user))
	s.NoError(s.org.Create(ctx, org))
	s.NoError(s.repo.Create(ctx, inv))

	inv.EventName = "updated"
	s.NoError(s.repo.Update(ctx, inv))

	got, _ := s.repo.Get(ctx, inv.ID)
	s.Equal("updated", got.EventName)
}

func (s *EquipmentInvocationRepositorySuite) TestDelete() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	inv := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(ctx, user))
	s.NoError(s.org.Create(ctx, org))
	s.NoError(s.repo.Create(ctx, inv))

	s.NoError(s.repo.Delete(ctx, inv.ID))

	_, err := s.repo.Get(ctx, inv.ID)
	s.Error(err)
}

func TestEquipmentInvocationRepositorySuite(t *testing.T) {
	suite.Run(t, new(EquipmentInvocationRepositorySuite))
}

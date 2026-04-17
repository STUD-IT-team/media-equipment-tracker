//go:build integration

package pgequipmentinvocation_test

import (
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
	s.repo = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gdb)
	s.user = pguser.NewPostgresUserRepository(gdb)
	s.org = pgorganization.NewPostgresOrganizationRepository(gdb)
}

func (s *EquipmentInvocationRepositorySuite) TestCreate_Get() {
	org := newOrg()
	user := newUser()
	inv := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(user))
	s.NoError(s.org.Create(org))

	s.Require().NoError(s.repo.Create(inv))

	got, err := s.repo.Get(inv.ID)
	s.NoError(err)
	s.Equal(inv.EventName, got.EventName)
}

func (s *EquipmentInvocationRepositorySuite) TestGet_NotFound() {
	_, err := s.repo.Get(uuid.New())
	s.Error(err)
}

func (s *EquipmentInvocationRepositorySuite) TestList() {
	org := newOrg()
	user := newUser()
	inv1 := newInvocation(&org.ID, nil, user.ID, nil)
	inv2 := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(user))
	s.NoError(s.org.Create(org))

	s.NoError(s.repo.Create(inv1))
	s.NoError(s.repo.Create(inv2))

	list, err := s.repo.List()
	s.NoError(err)
	s.Len(list, 2)
}

func (s *EquipmentInvocationRepositorySuite) TestUpdate() {
	org := newOrg()
	user := newUser()
	inv := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(user))
	s.NoError(s.org.Create(org))
	s.NoError(s.repo.Create(inv))

	inv.EventName = "updated"
	s.NoError(s.repo.Update(inv))

	got, _ := s.repo.Get(inv.ID)
	s.Equal("updated", got.EventName)
}

func (s *EquipmentInvocationRepositorySuite) TestDelete() {
	org := newOrg()
	user := newUser()
	inv := newInvocation(&org.ID, nil, user.ID, nil)
	s.NoError(s.user.Create(user))
	s.NoError(s.org.Create(org))
	s.NoError(s.repo.Create(inv))

	s.NoError(s.repo.Delete(inv.ID))

	_, err := s.repo.Get(inv.ID)
	s.Error(err)
}

func TestEquipmentInvocationRepositorySuite(t *testing.T) {
	suite.Run(t, new(EquipmentInvocationRepositorySuite))
}

package pgequipmentininvocation_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgequipmentininvocation"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type EquipmentInInvocationRepositorySuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.EquipmentInInvocationRepository
}

func (s *EquipmentInInvocationRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *EquipmentInInvocationRepositorySuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *EquipmentInInvocationRepositorySuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgequipmentininvocation.NewPostgresEquipmentInInvocationRepository(gdb)
}

func (s *EquipmentInInvocationRepositorySuite) TestCreate_Get() {
	inv := newInvocation(uuid.New(), uuid.New())
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	// Note: assuming inv and eq are created elsewhere, here just test the junction
	s.Require().NoError(s.repo.Create(eii))

	got, err := s.repo.Get(inv.ID, eq.ID)
	s.NoError(err)
	s.Equal(eii.Status, got.Status)
}

func (s *EquipmentInInvocationRepositorySuite) TestGet_NotFound() {
	_, err := s.repo.Get(uuid.New(), uuid.New())
	s.Error(err)
}

func (s *EquipmentInInvocationRepositorySuite) TestUpdate() {
	inv := newInvocation(uuid.New(), uuid.New())
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)
	s.Require().NoError(s.repo.Create(eii))

	eii.Status = domain.EquipmentIssued
	s.NoError(s.repo.Update(eii))

	got, _ := s.repo.Get(inv.ID, eq.ID)
	s.Equal(domain.EquipmentIssued, got.Status)
}

func (s *EquipmentInInvocationRepositorySuite) TestDelete() {
	inv := newInvocation(uuid.New(), uuid.New())
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)
	s.Require().NoError(s.repo.Create(eii))

	s.NoError(s.repo.Delete(inv.ID, eq.ID))

	_, err := s.repo.Get(inv.ID, eq.ID)
	s.Error(err)
}

func TestEquipmentInInvocationRepositorySuite(t *testing.T) {
	suite.Run(t, new(EquipmentInInvocationRepositorySuite))
}

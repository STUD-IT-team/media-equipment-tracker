//go:build integration

package pgequipmentininvocation_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentininvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type EquipmentInInvocationRepositorySuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.EquipmentInInvocationRepository

	userID uuid.UUID
	deptID uuid.UUID
	eqID   uuid.UUID
	invID  uuid.UUID
}

func (s *EquipmentInInvocationRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	userRepo := pguser.NewPostgresUserRepository(gdb)
	deptRepo := pgdepartment.NewPostgresDepartmentRepository(gdb)
	eqRepo := pgequipment.NewPostgresEquipmentRepository(gdb)
	invRepo := pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gdb)

	user := newUser()
	dept := newDept()
	eq := newEquipment()
	inv := newInvocation(dept.ID, user.ID)

	s.userID = user.ID
	s.deptID = dept.ID
	s.eqID = eq.ID
	s.invID = inv.ID

	s.Require().NoError(userRepo.Create(user))
	s.Require().NoError(deptRepo.Create(dept))
	s.Require().NoError(eqRepo.Create(eq))
	s.Require().NoError(invRepo.Create(inv))

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
	eii := newEqInInv(s.invID, s.eqID)

	s.Require().NoError(s.repo.Create(eii))

	got, err := s.repo.Get(s.invID, s.eqID)
	s.NoError(err)
	s.Equal(eii.Status, got.Status)
}

func (s *EquipmentInInvocationRepositorySuite) TestGet_NotFound() {
	_, err := s.repo.Get(uuid.New(), uuid.New())
	s.Error(err)
}

func (s *EquipmentInInvocationRepositorySuite) TestUpdate() {
	eii := newEqInInv(s.invID, s.eqID)
	s.Require().NoError(s.repo.Create(eii))

	eii.Status = domain.EquipmentIssued
	s.NoError(s.repo.Update(eii))

	got, _ := s.repo.Get(s.invID, s.eqID)
	s.Equal(domain.EquipmentIssued, got.Status)
}

func (s *EquipmentInInvocationRepositorySuite) TestDelete() {
	eii := newEqInInv(s.invID, s.eqID)
	s.Require().NoError(s.repo.Create(eii))

	s.NoError(s.repo.Delete(s.invID, s.eqID))

	_, err := s.repo.Get(s.invID, s.eqID)
	s.Error(err)
}

func TestEquipmentInInvocationRepositorySuite(t *testing.T) {
	suite.Run(t, new(EquipmentInInvocationRepositorySuite))
}

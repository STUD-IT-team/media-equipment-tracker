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

type EquipmentInInvocationRelationsSuite struct {
	suite.Suite

	db      *gorm.DB
	pg      *pgtest.PgTestDatabase
	repo    domain.EquipmentInInvocationRepository
	invRepo domain.EquipmentInvocationRepository
	eqRepo  domain.EquipmentRepository
	UserID  uuid.UUID
	DeptID  uuid.UUID
}

func (s *EquipmentInInvocationRelationsSuite) SetupSuite() {
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

	user := newUser()
	dept := newDept()
	s.UserID = user.ID
	s.DeptID = dept.ID

	s.Require().NoError(userRepo.Create(user))
	s.Require().NoError(deptRepo.Create(dept))

	s.Require().NoError(pg.CreateTemplate())
}

func (s *EquipmentInInvocationRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *EquipmentInInvocationRelationsSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgequipmentininvocation.NewPostgresEquipmentInInvocationRepository(gdb)
	s.invRepo = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gdb)
	s.eqRepo = pgequipment.NewPostgresEquipmentRepository(gdb)
}

// --- Invocation (belongs to) ---

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_AutoCreate() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.eqRepo.Create(eq))

	// FK ошибка
	s.Require().Error(s.repo.Create(eii))
	_, err := s.invRepo.Get(inv.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_AutoCreateOnUpdate() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.eqRepo.Create(eq))

	inv2 := newInvocation(s.DeptID, s.UserID)
	eii.InvocationID = inv2.ID
	eii.Invocation = inv2
	// FK ошибка
	s.Require().Error(s.repo.Update(eii))

	_, err := s.invRepo.Get(inv2.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_Preload() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(inv))
	s.Require().NoError(s.eqRepo.Create(eq))
	s.Require().NoError(s.repo.Create(eii))

	// без preload
	raw, _ := s.repo.Get(inv.ID, eq.ID)
	s.Nil(raw.Invocation)

	// с preload
	with, _ := s.repo.Get(inv.ID, eq.ID, domain.WithInvocation())
	s.NotNil(with.Invocation)
	s.Equal(inv.ID, with.Invocation.ID)
}

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_NoUpdateThroughEquipmentInInvocation() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(inv))
	s.Require().NoError(s.eqRepo.Create(eq))
	s.Require().NoError(s.repo.Create(eii))

	// preload
	with, _ := s.repo.Get(inv.ID, eq.ID, domain.WithInvocation())

	with.Invocation.EventName = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.invRepo.Get(inv.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.EventName)
}

// --- Equipment (belongs to) ---

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_AutoCreate() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(inv))

	// FK ошибка
	s.Require().Error(s.repo.Create(eii))
	_, err := s.eqRepo.Get(eq.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_AutoCreateOnUpdate() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(inv))

	eq2 := newEquipment()
	eii.EquipmentID = eq2.ID
	eii.Equipment = eq2
	// FK ошибка
	s.Require().Error(s.repo.Update(eii))

	_, err := s.eqRepo.Get(eq2.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_Preload() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(inv))
	s.Require().NoError(s.eqRepo.Create(eq))
	s.Require().NoError(s.repo.Create(eii))

	// без preload
	raw, _ := s.repo.Get(inv.ID, eq.ID)
	s.Nil(raw.Equipment)

	// с preload
	with, _ := s.repo.Get(inv.ID, eq.ID, domain.WithEquipment())
	s.NotNil(with.Equipment)
	s.Equal(eq.ID, with.Equipment.ID)
}

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_NoUpdateThroughEquipmentInInvocation() {
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(inv))
	s.Require().NoError(s.eqRepo.Create(eq))
	s.Require().NoError(s.repo.Create(eii))

	// preload
	with, _ := s.repo.Get(inv.ID, eq.ID, domain.WithEquipment())

	with.Equipment.Name = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.eqRepo.Get(eq.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.Name)
}

// ===== run =====

func TestEquipmentInInvocationRelationsSuite(t *testing.T) {
	suite.Run(t, new(EquipmentInInvocationRelationsSuite))
}

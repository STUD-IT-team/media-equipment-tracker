//go:build integration

package pgequipmentininvocation_test

import (
	"context"
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
	"media-equipment-tracker/pkg/txmanager/gormtx"
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

	ctx := context.Background()
	userRepo := pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
	deptRepo := pgdepartment.NewPostgresDepartmentRepository(gormtx.NewDBGetter(gdb))

	user := newUser()
	dept := newDept()
	s.UserID = user.ID
	s.DeptID = dept.ID

	s.Require().NoError(userRepo.Create(ctx, user))
	s.Require().NoError(deptRepo.Create(ctx, dept))

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
	s.repo = pgequipmentininvocation.NewPostgresEquipmentInInvocationRepository(gormtx.NewDBGetter(gdb))
	s.invRepo = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))
	s.eqRepo = pgequipment.NewPostgresEquipmentRepository(gormtx.NewDBGetter(gdb))
}

// --- Invocation (belongs to) ---

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_AutoCreate() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.eqRepo.Create(ctx, eq))

	// FK ошибка
	s.Require().Error(s.repo.Create(ctx, eii))
	_, err := s.invRepo.Get(ctx, inv.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_AutoCreateOnUpdate() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.eqRepo.Create(ctx, eq))

	inv2 := newInvocation(s.DeptID, s.UserID)
	eii.InvocationID = inv2.ID
	eii.Invocation = inv2
	// FK ошибка
	s.Require().Error(s.repo.Update(ctx, eii))

	_, err := s.invRepo.Get(ctx, inv2.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_Preload() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(ctx, inv))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))
	s.Require().NoError(s.repo.Create(ctx, eii))

	// без preload
	raw, _ := s.repo.Get(ctx, inv.ID, eq.ID)
	s.Nil(raw.Invocation)

	// с preload
	with, _ := s.repo.Get(ctx, inv.ID, eq.ID, domain.WithInvocation())
	s.NotNil(with.Invocation)
	s.Equal(inv.ID, with.Invocation.ID)
}

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_NoUpdateThroughEquipmentInInvocation() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(ctx, inv))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))
	s.Require().NoError(s.repo.Create(ctx, eii))

	// preload
	with, _ := s.repo.Get(ctx, inv.ID, eq.ID, domain.WithInvocation())

	with.Invocation.EventName = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.invRepo.Get(ctx, inv.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.EventName)
}

// --- Equipment (belongs to) ---

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_AutoCreate() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(ctx, inv))

	// FK ошибка
	s.Require().Error(s.repo.Create(ctx, eii))
	_, err := s.eqRepo.Get(ctx, eq.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_AutoCreateOnUpdate() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(ctx, inv))

	eq2 := newEquipment()
	eii.EquipmentID = eq2.ID
	eii.Equipment = eq2
	// FK ошибка
	s.Require().Error(s.repo.Update(ctx, eii))

	_, err := s.eqRepo.Get(ctx, eq2.ID)
	s.Error(err)
}

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_Preload() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(ctx, inv))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))
	s.Require().NoError(s.repo.Create(ctx, eii))

	// без preload
	raw, _ := s.repo.Get(ctx, inv.ID, eq.ID)
	s.Nil(raw.Equipment)

	// с preload
	with, _ := s.repo.Get(ctx, inv.ID, eq.ID, domain.WithEquipment())
	s.NotNil(with.Equipment)
	s.Equal(eq.ID, with.Equipment.ID)
}

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_NoUpdateThroughEquipmentInInvocation() {
	ctx := context.Background()
	inv := newInvocation(s.DeptID, s.UserID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.invRepo.Create(ctx, inv))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))
	s.Require().NoError(s.repo.Create(ctx, eii))

	// preload
	with, _ := s.repo.Get(ctx, inv.ID, eq.ID, domain.WithEquipment())

	with.Equipment.Name = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.eqRepo.Get(ctx, eq.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.Name)
}

// ===== run =====

func TestEquipmentInInvocationRelationsSuite(t *testing.T) {
	suite.Run(t, new(EquipmentInInvocationRelationsSuite))
}

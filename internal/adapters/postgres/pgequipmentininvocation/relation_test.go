package pgequipmentininvocation_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentininvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

func newUser() *domain.User {
	return &domain.User{
		ID:           uuid.New(),
		FullName:     "test",
		Email:        uuid.NewString() + "@mail.ru",
		HashPassword: "hash",
		Nice:         domain.DefaultNice,
	}
}

type EquipmentInInvocationRelationsSuite struct {
	suite.Suite

	db       *gorm.DB
	pg       *pgtest.PgTestDatabase
	repo     domain.EquipmentInInvocationRepository
	invRepo  domain.EquipmentInvocationRepository
	eqRepo   domain.EquipmentRepository
	userRepo domain.UserRepository
}

func (s *EquipmentInInvocationRelationsSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
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
	s.userRepo = pguser.NewPostgresUserRepository(gdb)
}

// --- Invocation (belongs to) ---

func (s *EquipmentInInvocationRelationsSuite) TestInvocation_Preload() {
	user := newUser()
	inv := newInvocation(uuid.New(), user.ID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.userRepo.Create(user))
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
	user := newUser()
	inv := newInvocation(uuid.New(), user.ID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.userRepo.Create(user))
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

func (s *EquipmentInInvocationRelationsSuite) TestEquipment_Preload() {
	user := newUser()
	inv := newInvocation(uuid.New(), user.ID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.userRepo.Create(user))
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
	user := newUser()
	inv := newInvocation(uuid.New(), user.ID)
	eq := newEquipment()
	eii := newEqInInv(inv.ID, eq.ID)

	s.Require().NoError(s.userRepo.Create(user))
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

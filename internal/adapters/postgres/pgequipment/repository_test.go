package pgequipment_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type EquipmentRepositorySuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.EquipmentRepository
}

func (s *EquipmentRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *EquipmentRepositorySuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *EquipmentRepositorySuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgequipment.NewPostgresEquipmentRepository(gdb)
}

// --- tests ---

func (s *EquipmentRepositorySuite) TestCreate_Get() {
	eq := newEquipment()

	s.Require().NoError(s.repo.Create(eq))

	got, err := s.repo.Get(eq.ID)
	s.NoError(err)
	s.Equal(eq.InventoryNumber, got.InventoryNumber)
}

func (s *EquipmentRepositorySuite) TestGet_NotFound() {
	_, err := s.repo.Get(uuid.New())
	s.Error(err)
}

func (s *EquipmentRepositorySuite) TestList() {
	e1, e2 := newEquipment(), newEquipment()

	err := s.repo.Create(e1)
	s.NoError(err)
	err = s.repo.Create(e2)
	s.NoError(err)

	list, err := s.repo.List()
	s.NoError(err)
	s.Len(list, 2)
}

func (s *EquipmentRepositorySuite) TestUpdate() {
	eq := newEquipment()
	err := s.repo.Create(eq)
	s.NoError(err)

	eq.Name = "updated"
	s.NoError(s.repo.Update(eq))

	got, err := s.repo.Get(eq.ID)
	s.NoError(err)
	s.Equal("updated", got.Name)
}

func (s *EquipmentRepositorySuite) TestDelete() {
	eq := newEquipment()
	err := s.repo.Create(eq)
	s.NoError(err)

	s.NoError(s.repo.Delete(eq.ID))

	_, err = s.repo.Get(eq.ID)
	s.Error(err)
}

func (s *EquipmentRepositorySuite) TestReload() {
	eq := newEquipment()
	err := s.repo.Create(eq)
	s.NoError(err)

	eq.Name = "updated"
	err = s.repo.Update(eq)
	s.NoError(err)

	eq.Name = "stale"
	s.NoError(s.repo.Reload(eq))

	s.Equal("updated", eq.Name)
}

// ===== run =====

func TestEquipmentRepositorySuite(t *testing.T) {
	suite.Run(t, new(EquipmentRepositorySuite))
}

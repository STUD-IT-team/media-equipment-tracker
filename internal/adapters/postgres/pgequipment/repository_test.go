//go:build integration

package pgequipment_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
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
	s.repo = pgequipment.NewPostgresEquipmentRepository(gormtx.NewDBGetter(gdb))
}

// --- tests ---

func (s *EquipmentRepositorySuite) TestCreate_Get() {
	ctx := context.Background()
	eq := newEquipment()

	s.Require().NoError(s.repo.Create(ctx, eq))

	got, err := s.repo.Get(ctx, eq.ID)
	s.NoError(err)
	s.Equal(eq.InventoryNumber, got.InventoryNumber)
}

func (s *EquipmentRepositorySuite) TestGet_NotFound() {
	ctx := context.Background()
	_, err := s.repo.Get(ctx, uuid.New())
	s.Error(err)
}

func (s *EquipmentRepositorySuite) TestList() {
	ctx := context.Background()
	e1, e2 := newEquipment(), newEquipment()

	err := s.repo.Create(ctx, e1)
	s.NoError(err)
	err = s.repo.Create(ctx, e2)
	s.NoError(err)

	list, err := s.repo.List(ctx)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *EquipmentRepositorySuite) TestUpdate() {
	ctx := context.Background()
	eq := newEquipment()
	err := s.repo.Create(ctx, eq)
	s.NoError(err)

	eq.Name = "updated"
	s.NoError(s.repo.Update(ctx, eq))

	got, err := s.repo.Get(ctx, eq.ID)
	s.NoError(err)
	s.Equal("updated", got.Name)
}

func (s *EquipmentRepositorySuite) TestDelete() {
	ctx := context.Background()
	eq := newEquipment()
	err := s.repo.Create(ctx, eq)
	s.NoError(err)

	s.NoError(s.repo.Delete(ctx, eq.ID))

	_, err = s.repo.Get(ctx, eq.ID)
	s.Error(err)
}

func (s *EquipmentRepositorySuite) TestReload() {
	ctx := context.Background()
	eq := newEquipment()
	err := s.repo.Create(ctx, eq)
	s.NoError(err)

	eq.Name = "updated"
	err = s.repo.Update(ctx, eq)
	s.NoError(err)

	eq.Name = "stale"
	s.NoError(s.repo.Reload(ctx, eq))

	s.Equal("updated", eq.Name)
}

// ===== run =====

func TestEquipmentRepositorySuite(t *testing.T) {
	suite.Run(t, new(EquipmentRepositorySuite))
}

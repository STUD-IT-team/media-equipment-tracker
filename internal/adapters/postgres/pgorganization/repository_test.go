//go:build integration

package pgorganization_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgorganization"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
)

type OrganizationRepositorySuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.OrganizationRepository
}

func (s *OrganizationRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *OrganizationRepositorySuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *OrganizationRepositorySuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgorganization.NewPostgresOrganizationRepository(gormtx.NewDBGetter(gdb))
}

func (s *OrganizationRepositorySuite) TestCreate_Get() {
	ctx := context.Background()
	org := newOrg()

	s.Require().NoError(s.repo.Create(ctx, org))

	got, err := s.repo.Get(ctx, org.ID)
	s.NoError(err)
	s.Equal(org.Name, got.Name)
}

func (s *OrganizationRepositorySuite) TestGet_NotFound() {
	ctx := context.Background()
	_, err := s.repo.Get(ctx, uuid.New())
	s.Error(err)
}

func (s *OrganizationRepositorySuite) TestList() {
	ctx := context.Background()
	o1, o2 := newOrg(), newOrg()

	s.NoError(s.repo.Create(ctx, o1))
	s.NoError(s.repo.Create(ctx, o2))

	list, err := s.repo.List(ctx)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *OrganizationRepositorySuite) TestUpdate() {
	ctx := context.Background()
	org := newOrg()
	s.NoError(s.repo.Create(ctx, org))

	org.Name = "updated"
	s.NoError(s.repo.Update(ctx, org))

	got, _ := s.repo.Get(ctx, org.ID)
	s.Equal("updated", got.Name)
}

func (s *OrganizationRepositorySuite) TestDelete() {
	ctx := context.Background()
	org := newOrg()
	s.NoError(s.repo.Create(ctx, org))

	s.NoError(s.repo.Delete(ctx, org.ID))

	_, err := s.repo.Get(ctx, org.ID)
	s.Error(err)
}

func TestOrganizationRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrganizationRepositorySuite))
}

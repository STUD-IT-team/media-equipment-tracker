//go:build integration

package pguser_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
)

type UserSuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.UserRepository
}

func (s *UserSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())

}

func (s *UserSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *UserSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())
	db := stdlib.OpenDBFromPool(s.pg.Pool())

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	s.Require().NoError(err)

	s.db = gdb
	s.repo = pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
}

func (s *UserSuite) TestCreate_Get() {
	ctx := context.Background()
	u := newUser()

	s.Require().NoError(s.repo.Create(ctx, u))

	got, err := s.repo.Get(ctx, u.ID)
	s.NoError(err)
	s.Equal(u.Email, got.Email)
}

func (s *UserSuite) TestGet_NotFound() {
	ctx := context.Background()
	_, err := s.repo.Get(ctx, uuid.New())
	s.Error(err)
}

func (s *UserSuite) TestGetByEmail() {
	ctx := context.Background()
	u := newUser()
	_ = s.repo.Create(ctx, u)

	got, err := s.repo.GetByEmail(ctx, u.Email)
	s.NoError(err)
	s.Equal(u.ID, got.ID)
}

func (s *UserSuite) TestGetByEmail_NotFound() {
	ctx := context.Background()
	_, err := s.repo.GetByEmail(ctx, "nope@mail.ru")
	s.Error(err)
}

func (s *UserSuite) TestList() {
	ctx := context.Background()
	u1, u2 := newUser(), newUser()
	_ = s.repo.Create(ctx, u1)
	_ = s.repo.Create(ctx, u2)

	list, err := s.repo.List(ctx)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *UserSuite) TestUpdate() {
	ctx := context.Background()
	u := newUser()
	_ = s.repo.Create(ctx, u)

	u.FullName = "updated"
	s.NoError(s.repo.Update(ctx, u))

	got, _ := s.repo.Get(ctx, u.ID)
	s.Equal("updated", got.FullName)
}

func (s *UserSuite) TestDelete() {
	ctx := context.Background()
	u := newUser()
	_ = s.repo.Create(ctx, u)

	s.NoError(s.repo.Delete(ctx, u.ID))

	_, err := s.repo.Get(ctx, u.ID)
	s.Error(err)
}

func (s *UserSuite) TestReload() {
	ctx := context.Background()
	u := newUser()
	_ = s.repo.Create(ctx, u)

	u.FullName = "updated"
	_ = s.repo.Update(ctx, u)

	u.FullName = "stale"
	s.NoError(s.repo.Reload(ctx, u))

	s.Equal("updated", u.FullName)
}

// ===== run =====

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

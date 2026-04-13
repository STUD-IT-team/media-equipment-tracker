//go:build integration

package postgresuser_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/postgresuser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
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
	s.repo = postgresuser.NewPostgresUserRepository(gdb)
}

func (s *UserSuite) TestCreate_Get() {
	u := newUser()

	s.Require().NoError(s.repo.Create(u))

	got, err := s.repo.Get(u.ID)
	s.NoError(err)
	s.Equal(u.Email, got.Email)
}

func (s *UserSuite) TestGet_NotFound() {
	_, err := s.repo.Get(uuid.New())
	s.Error(err)
}

func (s *UserSuite) TestGetByEmail() {
	u := newUser()
	_ = s.repo.Create(u)

	got, err := s.repo.GetByEmail(u.Email)
	s.NoError(err)
	s.Equal(u.ID, got.ID)
}

func (s *UserSuite) TestGetByEmail_NotFound() {
	_, err := s.repo.GetByEmail("nope@mail.ru")
	s.Error(err)
}

func (s *UserSuite) TestList() {
	u1, u2 := newUser(), newUser()
	_ = s.repo.Create(u1)
	_ = s.repo.Create(u2)

	list, err := s.repo.List()
	s.NoError(err)
	s.Len(list, 2)
}

func (s *UserSuite) TestUpdate() {
	u := newUser()
	_ = s.repo.Create(u)

	u.FullName = "updated"
	s.NoError(s.repo.Update(u))

	got, _ := s.repo.Get(u.ID)
	s.Equal("updated", got.FullName)
}

func (s *UserSuite) TestDelete() {
	u := newUser()
	_ = s.repo.Create(u)

	s.NoError(s.repo.Delete(u.ID))

	_, err := s.repo.Get(u.ID)
	s.Error(err)
}

func (s *UserSuite) TestReload() {
	u := newUser()
	_ = s.repo.Create(u)

	u.FullName = "updated"
	_ = s.repo.Update(u)

	u.FullName = "stale"
	s.NoError(s.repo.Reload(u))

	s.Equal("updated", u.FullName)
}

// ===== run =====

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

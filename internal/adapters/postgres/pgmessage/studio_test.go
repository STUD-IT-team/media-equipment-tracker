//go:build integration

package pgmessage_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgmessage"
	"media-equipment-tracker/internal/adapters/postgres/pgstudioinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type MessageStudioInvocationRepositorySuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.MessageStudioInvocationRepository

	senderID    uuid.UUID
	recipientID uuid.UUID
	invID       uuid.UUID
}

func (s *MessageStudioInvocationRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	userRepo := pguser.NewPostgresUserRepository(gdb)
	deptRepo := pgdepartment.NewPostgresDepartmentRepository(gdb)
	invRepo := pgstudioinvocation.NewPostgresStudioInvocationRepository(gdb)

	sender := newUser()
	recipient := newUser()
	dept := newDepartment()
	inv := newStudioInvocation(dept.ID, sender.ID)

	s.senderID = sender.ID
	s.recipientID = recipient.ID
	s.invID = inv.ID

	s.Require().NoError(userRepo.Create(sender))
	s.Require().NoError(userRepo.Create(recipient))
	s.Require().NoError(deptRepo.Create(dept))
	s.Require().NoError(invRepo.Create(inv))

	s.Require().NoError(pg.CreateTemplate())
}

func (s *MessageStudioInvocationRepositorySuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *MessageStudioInvocationRepositorySuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgmessage.NewPostgresMessageStudioInvocationRepository(gdb)
}

func (s *MessageStudioInvocationRepositorySuite) TestCreate_Get() {
	msg := newMessageStudioInvocation(s.invID, s.senderID, s.recipientID)

	s.Require().NoError(s.repo.Create(msg))

	got, err := s.repo.Get(msg.ID)
	s.NoError(err)
	s.Equal(msg.Content, got.Content)
}

func (s *MessageStudioInvocationRepositorySuite) TestGet_NotFound() {
	_, err := s.repo.Get(uuid.New())
	s.Error(err)
}

func (s *MessageStudioInvocationRepositorySuite) TestGetInvocation() {
	msg1 := newMessageStudioInvocation(s.invID, s.senderID, s.recipientID)
	msg2 := newMessageStudioInvocation(s.invID, s.senderID, s.recipientID)

	s.Require().NoError(s.repo.Create(msg1))
	s.Require().NoError(s.repo.Create(msg2))

	list, err := s.repo.GetInvocation(s.invID)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *MessageStudioInvocationRepositorySuite) TestUpdate() {
	msg := newMessageStudioInvocation(s.invID, s.senderID, s.recipientID)
	s.NoError(s.repo.Create(msg))

	msg.Content = "updated"
	s.NoError(s.repo.Update(msg))

	got, _ := s.repo.Get(msg.ID)
	s.Equal("updated", got.Content)
}

func (s *MessageStudioInvocationRepositorySuite) TestDelete() {
	msg := newMessageStudioInvocation(s.invID, s.senderID, s.recipientID)
	s.NoError(s.repo.Create(msg))

	s.NoError(s.repo.Delete(msg.ID))

	_, err := s.repo.Get(msg.ID)
	s.Error(err)
}

func TestMessageStudioInvocationRepositorySuite(t *testing.T) {
	suite.Run(t, new(MessageStudioInvocationRepositorySuite))
}

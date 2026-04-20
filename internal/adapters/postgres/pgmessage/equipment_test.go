//go:build integration

package pgmessage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgmessage"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
)

type MessageEquipmentInvocationRepositorySuite struct {
	suite.Suite

	db   *gorm.DB
	pg   *pgtest.PgTestDatabase
	repo domain.MessageEquipmentInvocationRepository

	senderID    uuid.UUID
	recipientID uuid.UUID
	invID       uuid.UUID
}

func (s *MessageEquipmentInvocationRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	ctx := context.Background()
	userRepo := pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
	deptRepo := pgdepartment.NewPostgresDepartmentRepository(gormtx.NewDBGetter(gdb))
	invRepo := pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))

	sender := newUser()
	recipient := newUser()
	dept := newDepartment()
	inv := newEquipmentInvocation(dept.ID, sender.ID)

	s.senderID = sender.ID
	s.recipientID = recipient.ID
	s.invID = inv.ID

	s.Require().NoError(userRepo.Create(ctx, sender))
	s.Require().NoError(userRepo.Create(ctx, recipient))
	s.Require().NoError(deptRepo.Create(ctx, dept))
	s.Require().NoError(invRepo.Create(ctx, inv))

	s.Require().NoError(pg.CreateTemplate())
}

func (s *MessageEquipmentInvocationRepositorySuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *MessageEquipmentInvocationRepositorySuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgmessage.NewPostgresMessageEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))
}

func (s *MessageEquipmentInvocationRepositorySuite) TestCreate_Get() {
	ctx := context.Background()
	msg := newMessageEquipmentInvocation(s.invID, s.senderID, s.recipientID)

	s.Require().NoError(s.repo.Create(ctx, msg))

	got, err := s.repo.Get(ctx, msg.ID)
	s.NoError(err)
	s.Equal(msg.Content, got.Content)
}

func (s *MessageEquipmentInvocationRepositorySuite) TestGet_NotFound() {
	ctx := context.Background()
	_, err := s.repo.Get(ctx, uuid.New())
	s.Error(err)
}

func (s *MessageEquipmentInvocationRepositorySuite) TestGetInvocation() {
	ctx := context.Background()
	msg1 := newMessageEquipmentInvocation(s.invID, s.senderID, s.recipientID)
	msg2 := newMessageEquipmentInvocation(s.invID, s.senderID, s.recipientID)

	s.Require().NoError(s.repo.Create(ctx, msg1))
	s.Require().NoError(s.repo.Create(ctx, msg2))

	list, err := s.repo.GetInvocation(ctx, s.invID)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *MessageEquipmentInvocationRepositorySuite) TestUpdate() {
	ctx := context.Background()
	msg := newMessageEquipmentInvocation(s.invID, s.senderID, s.recipientID)
	s.NoError(s.repo.Create(ctx, msg))

	msg.Content = "updated"
	s.NoError(s.repo.Update(ctx, msg))

	got, _ := s.repo.Get(ctx, msg.ID)
	s.Equal("updated", got.Content)
}

func (s *MessageEquipmentInvocationRepositorySuite) TestDelete() {
	ctx := context.Background()
	msg := newMessageEquipmentInvocation(s.invID, s.senderID, s.recipientID)
	s.NoError(s.repo.Create(ctx, msg))

	s.NoError(s.repo.Delete(ctx, msg.ID))

	_, err := s.repo.Get(ctx, msg.ID)
	s.Error(err)
}

func TestMessageEquipmentInvocationRepositorySuite(t *testing.T) {
	suite.Run(t, new(MessageEquipmentInvocationRepositorySuite))
}

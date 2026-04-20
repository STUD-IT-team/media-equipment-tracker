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

type MessageEquipmentInvocationRelationsSuite struct {
	suite.Suite

	db       *gorm.DB
	pg       *pgtest.PgTestDatabase
	repo     domain.MessageEquipmentInvocationRepository
	userRepo domain.UserRepository

	invID uuid.UUID
}

func (s *MessageEquipmentInvocationRelationsSuite) SetupSuite() {
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
	invRepo := pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))
	deptRepo := pgdepartment.NewPostgresDepartmentRepository(gormtx.NewDBGetter(gdb))
	userRepo := pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))

	dept := newDepartment()
	user := newUser()
	inv := newEquipmentInvocation(dept.ID, user.ID)

	s.invID = inv.ID

	s.Require().NoError(deptRepo.Create(ctx, dept))
	s.Require().NoError(userRepo.Create(ctx, user))
	s.Require().NoError(invRepo.Create(ctx, inv))

	s.Require().NoError(pg.CreateTemplate())
}

func (s *MessageEquipmentInvocationRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *MessageEquipmentInvocationRelationsSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgmessage.NewPostgresMessageEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))
	s.userRepo = pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
}

// --- Sender (belongs to) ---

func (s *MessageEquipmentInvocationRelationsSuite) TestSender_AutoCreate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().Error(s.repo.Create(ctx, msg))

	_, err := s.userRepo.Get(ctx, sender.ID)
	s.Error(err)
}

func (s *MessageEquipmentInvocationRelationsSuite) TestSender_AutoCreateOnUpdate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	sender2 := newUser()
	sender2.FullName = "Wawa"
	msg.SenderID = sender2.ID
	msg.Sender = sender2
	// FK ошибка
	s.Require().Error(s.repo.Update(ctx, msg))

	_, err := s.userRepo.Get(ctx, sender2.ID)
	s.Error(err)
}

func (s *MessageEquipmentInvocationRelationsSuite) TestSender_Preload() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// без preload
	raw, _ := s.repo.Get(ctx, msg.ID)
	s.Nil(raw.Sender)

	// с preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageEquipmentInvocationWithSender())
	s.NotNil(with.Sender)
	s.Equal(sender.ID, with.Sender.ID)
}

func (s *MessageEquipmentInvocationRelationsSuite) TestSender_NoUpdateThroughMessage() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageEquipmentInvocationWithSender())

	with.Sender.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.userRepo.Get(ctx, sender.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// --- Recipient (belongs to) ---

func (s *MessageEquipmentInvocationRelationsSuite) TestRecipient_AutoCreate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().Error(s.repo.Create(ctx, msg))

	_, err := s.userRepo.Get(ctx, recipient.ID)
	s.Error(err)
}

func (s *MessageEquipmentInvocationRelationsSuite) TestRecipient_AutoCreateOnUpdate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	recipient2 := newUser()
	recipient2.FullName = "Wawa"
	msg.RecipientID = recipient2.ID
	msg.Recipient = recipient2
	// FK ошибка
	s.Require().Error(s.repo.Update(ctx, msg))

	_, err := s.userRepo.Get(ctx, recipient2.ID)
	s.Error(err)
}

func (s *MessageEquipmentInvocationRelationsSuite) TestRecipient_Preload() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// без preload
	raw, _ := s.repo.Get(ctx, msg.ID)
	s.Nil(raw.Recipient)

	// с preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageEquipmentInvocationWithRecipient())
	s.NotNil(with.Recipient)
	s.Equal(recipient.ID, with.Recipient.ID)
}

func (s *MessageEquipmentInvocationRelationsSuite) TestRecipient_NoUpdateThroughMessage() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageEquipmentInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageEquipmentInvocationWithRecipient())

	with.Recipient.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.userRepo.Get(ctx, recipient.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// ===== run =====

func TestMessageEquipmentInvocationRelationsSuite(t *testing.T) {
	suite.Run(t, new(MessageEquipmentInvocationRelationsSuite))
}

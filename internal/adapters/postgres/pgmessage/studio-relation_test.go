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
	"media-equipment-tracker/internal/adapters/postgres/pgmessage"
	"media-equipment-tracker/internal/adapters/postgres/pgstudioinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
)

type MessageStudioInvocationRelationsSuite struct {
	suite.Suite

	db       *gorm.DB
	pg       *pgtest.PgTestDatabase
	repo     domain.MessageStudioInvocationRepository
	userRepo domain.UserRepository

	invID uuid.UUID
}

func (s *MessageStudioInvocationRelationsSuite) SetupSuite() {
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
	invRepo := pgstudioinvocation.NewPostgresStudioInvocationRepository(gormtx.NewDBGetter(gdb))

	dept := newDepartment()
	user := newUser()
	inv := newStudioInvocation(dept.ID, user.ID)

	s.invID = inv.ID

	s.Require().NoError(userRepo.Create(ctx, user))
	s.Require().NoError(deptRepo.Create(ctx, dept))
	s.Require().NoError(invRepo.Create(ctx, inv))

	s.Require().NoError(pg.CreateTemplate())
}

func (s *MessageStudioInvocationRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *MessageStudioInvocationRelationsSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.repo = pgmessage.NewPostgresMessageStudioInvocationRepository(gormtx.NewDBGetter(gdb))
	s.userRepo = pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
}

// --- Sender (belongs to) ---

func (s *MessageStudioInvocationRelationsSuite) TestSender_AutoCreate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().Error(s.repo.Create(ctx, msg))

	_, err := s.userRepo.Get(ctx, sender.ID)
	s.Error(err)
}

func (s *MessageStudioInvocationRelationsSuite) TestSender_AutoCreateOnUpdate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

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

func (s *MessageStudioInvocationRelationsSuite) TestSender_Preload() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// без preload
	raw, _ := s.repo.Get(ctx, msg.ID)
	s.Nil(raw.Sender)

	// с preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageStudioInvocationWithSender())
	s.NotNil(with.Sender)
	s.Equal(sender.ID, with.Sender.ID)
}

func (s *MessageStudioInvocationRelationsSuite) TestSender_NoUpdateThroughMessage() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageStudioInvocationWithSender())

	with.Sender.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.userRepo.Get(ctx, sender.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// --- Recipient (belongs to) ---

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_AutoCreate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().Error(s.repo.Create(ctx, msg))

	_, err := s.userRepo.Get(ctx, recipient.ID)
	s.Error(err)
}

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_AutoCreateOnUpdate() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

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

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_Preload() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// без preload
	raw, _ := s.repo.Get(ctx, msg.ID)
	s.Nil(raw.Recipient)

	// с preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageStudioInvocationWithRecipient())
	s.NotNil(with.Recipient)
	s.Equal(recipient.ID, with.Recipient.ID)
}

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_NoUpdateThroughMessage() {
	ctx := context.Background()
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(ctx, sender))
	s.Require().NoError(s.userRepo.Create(ctx, recipient))
	s.Require().NoError(s.repo.Create(ctx, msg))

	// preload
	with, _ := s.repo.Get(ctx, msg.ID, domain.MessageStudioInvocationWithRecipient())

	with.Recipient.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.userRepo.Get(ctx, recipient.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// ===== run =====

func TestMessageStudioInvocationRelationsSuite(t *testing.T) {
	suite.Run(t, new(MessageStudioInvocationRelationsSuite))
}

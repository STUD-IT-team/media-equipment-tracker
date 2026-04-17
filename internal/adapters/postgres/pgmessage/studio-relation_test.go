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

	userRepo := pguser.NewPostgresUserRepository(gdb)
	deptRepo := pgdepartment.NewPostgresDepartmentRepository(gdb)
	invRepo := pgstudioinvocation.NewPostgresStudioInvocationRepository(gdb)

	dept := newDepartment()
	user := newUser()
	inv := newStudioInvocation(dept.ID, user.ID)

	s.invID = inv.ID

	s.Require().NoError(userRepo.Create(user))
	s.Require().NoError(deptRepo.Create(dept))
	s.Require().NoError(invRepo.Create(inv))

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
	s.repo = pgmessage.NewPostgresMessageStudioInvocationRepository(gdb)
	s.userRepo = pguser.NewPostgresUserRepository(gdb)
}

// --- Sender (belongs to) ---

func (s *MessageStudioInvocationRelationsSuite) TestSender_AutoCreate() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(recipient))
	s.Require().Error(s.repo.Create(msg))

	_, err := s.userRepo.Get(sender.ID)
	s.Error(err)
}

func (s *MessageStudioInvocationRelationsSuite) TestSender_AutoCreateOnUpdate() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(sender))
	s.Require().NoError(s.userRepo.Create(recipient))
	s.Require().NoError(s.repo.Create(msg))

	sender2 := newUser()
	sender2.FullName = "Wawa"
	msg.SenderID = sender2.ID
	msg.Sender = sender2
	// FK ошибка
	s.Require().Error(s.repo.Update(msg))

	_, err := s.userRepo.Get(sender2.ID)
	s.Error(err)
}

func (s *MessageStudioInvocationRelationsSuite) TestSender_Preload() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(sender))
	s.Require().NoError(s.userRepo.Create(recipient))
	s.Require().NoError(s.repo.Create(msg))

	// без preload
	raw, _ := s.repo.Get(msg.ID)
	s.Nil(raw.Sender)

	// с preload
	with, _ := s.repo.Get(msg.ID, domain.MessageStudioInvocationWithSender())
	s.NotNil(with.Sender)
	s.Equal(sender.ID, with.Sender.ID)
}

func (s *MessageStudioInvocationRelationsSuite) TestSender_NoUpdateThroughMessage() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(sender))
	s.Require().NoError(s.userRepo.Create(recipient))
	s.Require().NoError(s.repo.Create(msg))

	// preload
	with, _ := s.repo.Get(msg.ID, domain.MessageStudioInvocationWithSender())

	with.Sender.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.userRepo.Get(sender.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// --- Recipient (belongs to) ---

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_AutoCreate() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(sender))
	s.Require().Error(s.repo.Create(msg))

	_, err := s.userRepo.Get(recipient.ID)
	s.Error(err)
}

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_AutoCreateOnUpdate() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(sender))
	s.Require().NoError(s.userRepo.Create(recipient))
	s.Require().NoError(s.repo.Create(msg))

	recipient2 := newUser()
	recipient2.FullName = "Wawa"
	msg.RecipientID = recipient2.ID
	msg.Recipient = recipient2
	// FK ошибка
	s.Require().Error(s.repo.Update(msg))

	_, err := s.userRepo.Get(recipient2.ID)
	s.Error(err)
}

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_Preload() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(sender))
	s.Require().NoError(s.userRepo.Create(recipient))
	s.Require().NoError(s.repo.Create(msg))

	// без preload
	raw, _ := s.repo.Get(msg.ID)
	s.Nil(raw.Recipient)

	// с preload
	with, _ := s.repo.Get(msg.ID, domain.MessageStudioInvocationWithRecipient())
	s.NotNil(with.Recipient)
	s.Equal(recipient.ID, with.Recipient.ID)
}

func (s *MessageStudioInvocationRelationsSuite) TestRecipient_NoUpdateThroughMessage() {
	sender := newUser()
	recipient := newUser()
	msg := newMessageStudioInvocation(s.invID, sender.ID, recipient.ID)

	s.Require().NoError(s.userRepo.Create(sender))
	s.Require().NoError(s.userRepo.Create(recipient))
	s.Require().NoError(s.repo.Create(msg))

	// preload
	with, _ := s.repo.Get(msg.ID, domain.MessageStudioInvocationWithRecipient())

	with.Recipient.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.userRepo.Get(recipient.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// ===== run =====

func TestMessageStudioInvocationRelationsSuite(t *testing.T) {
	suite.Run(t, new(MessageStudioInvocationRelationsSuite))
}

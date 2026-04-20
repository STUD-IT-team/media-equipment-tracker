//go:build integration

package pgorganization_test

import (
	"testing"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgorganization"
	"media-equipment-tracker/internal/adapters/postgres/pgstudioinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type OrganizationRelationsSuite struct {
	suite.Suite

	db                   *gorm.DB
	pg                   *pgtest.PgTestDatabase
	orgRepo              domain.OrganizationRepository
	userRepo             domain.UserRepository
	eqInvocationRepo     domain.EquipmentInvocationRepository
	studioInvocationRepo domain.StudioInvocationRepository
	deptRepo             domain.DepartmentRepository
}

func (s *OrganizationRelationsSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *OrganizationRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *OrganizationRelationsSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.orgRepo = pgorganization.NewPostgresOrganizationRepository(gdb)
	s.userRepo = pguser.NewPostgresUserRepository(gdb)
	s.eqInvocationRepo = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gdb)
	s.studioInvocationRepo = pgstudioinvocation.NewPostgresStudioInvocationRepository(gdb)
	s.deptRepo = pgdepartment.NewPostgresDepartmentRepository(gdb)
}

// --- Users (many2many) ---

func (s *OrganizationRelationsSuite) TestUsers_Preload() {
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))

	user.Organizations = []*domain.Organization{org}
	s.Require().NoError(s.userRepo.Update(user))

	user, err := s.userRepo.Get(user.ID, domain.UserWithOrganizations())
	s.NoError(err)
	s.Len(user.Organizations, 1)
	s.Equal(org.ID, user.Organizations[0].ID)

	// без preload
	raw, _ := s.orgRepo.Get(org.ID)
	s.Empty(raw.Users)

	// с preload
	with, _ := s.orgRepo.Get(org.ID, domain.WithUsers())
	s.Len(with.Users, 1)
	s.Equal(user.ID, with.Users[0].ID)
}

func (s *OrganizationRelationsSuite) TestUsers_PreloadAfterDelete() {
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))

	user.Organizations = []*domain.Organization{org}
	s.Require().NoError(s.userRepo.Update(user))
	user.Organizations = []*domain.Organization{}
	s.Require().NoError(s.userRepo.Update(user))

	// без preload
	raw, _ := s.orgRepo.Get(org.ID)
	s.Empty(raw.Users)

	// с preload
	with, _ := s.orgRepo.Get(org.ID, domain.WithUsers())
	s.Empty(with.Users)
}

func (s *OrganizationRelationsSuite) TestUsers_NoAutoCreate() {
	org := newOrg()
	user := newUser()

	org.Users = []*domain.User{user}

	// пользователь не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().NoError(s.orgRepo.Create(org))

	_, err := s.userRepo.Get(user.ID)
	s.Error(err)
}

func (s *OrganizationRelationsSuite) TestUsers_NoAutoCreateOnUpdate() {
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.orgRepo.Create(org))

	org.Users = []*domain.User{user}
	s.Require().NoError(s.userRepo.Update(user))

	with, _ := s.orgRepo.Get(org.ID, domain.WithUsers())
	s.Empty(with.Users)
}

// --- EquipmentInvocations (has many) ---

func (s *OrganizationRelationsSuite) TestEquipmentInvocations_Preload() {
	org := newOrg()
	user := newUser()
	invocation := newInvocation(nil, user.ID)

	invocation.OrganizationID = &org.ID

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.eqInvocationRepo.Create(invocation))

	// без preload
	raw, _ := s.orgRepo.Get(org.ID)
	s.Empty(raw.EquipmentInvocations)

	// с preload
	with, _ := s.orgRepo.Get(org.ID, domain.WithEquipmentInvocations())
	s.Len(with.EquipmentInvocations, 1)
	s.Equal(invocation.ID, with.EquipmentInvocations[0].ID)
}

func (s *OrganizationRelationsSuite) TestEquipmentInvocations_NoAutoCreate() {
	org := newOrg()
	user := newUser()
	invocation := newInvocation(nil, user.ID)

	invocation.OrganizationID = &org.ID

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))

	org.EquipmentInvocations = []*domain.EquipmentInvocation{invocation}

	// invocation не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().NoError(s.orgRepo.Update(org))

	_, err := s.eqInvocationRepo.Get(invocation.ID)
	s.Error(err)
}

func (s *OrganizationRelationsSuite) TestEquipmentInvocations_NoAutoCreateOnUpdate() {
	org := newOrg()
	user := newUser()
	inv := newInvocation(nil, user.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))

	org.EquipmentInvocations = []*domain.EquipmentInvocation{inv}
	s.Require().NoError(s.userRepo.Update(user))

	with, _ := s.orgRepo.Get(org.ID, domain.WithEquipmentInvocations())
	s.Empty(with.EquipmentInvocations)
}

// --- StudioInvocations (has many) ---

func (s *OrganizationRelationsSuite) TestStudioInvocations_Preload() {
	org := newOrg()
	user := newUser()
	studioInvocation := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.studioInvocationRepo.Create(studioInvocation))

	// без preload
	raw, _ := s.orgRepo.Get(org.ID)
	s.Empty(raw.StudioInvocations)

	// с preload
	with, _ := s.orgRepo.Get(org.ID, domain.WithStudioInvocations())
	s.Len(with.StudioInvocations, 1)
	s.Equal(studioInvocation.ID, with.StudioInvocations[0].ID)
}

func (s *OrganizationRelationsSuite) TestStudioInvocations_NoAutoCreate() {
	org := newOrg()
	dept := newDept()
	user := newUser()
	studioInvocation := newStudioInvocation(&org.ID, &dept.ID, &user.ID, &user.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.deptRepo.Create(dept))
	s.Require().NoError(s.userRepo.Create(user))

	org.StudioInvocations = []*domain.StudioInvocation{studioInvocation}

	// studioInvocation не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().NoError(s.orgRepo.Update(org))

	_, err := s.studioInvocationRepo.Get(studioInvocation.ID)
	s.Error(err)
}

func (s *OrganizationRelationsSuite) TestStudioInvocations_NoAutoCreateOnUpdate() {
	org := newOrg()
	user := newUser()
	studioInvocation := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))

	org.StudioInvocations = []*domain.StudioInvocation{studioInvocation}
	s.Require().NoError(s.userRepo.Update(user))

	with, _ := s.orgRepo.Get(org.ID, domain.WithStudioInvocations())
	s.Empty(with.StudioInvocations)
}

// ===== run =====

func TestOrganizationRelationsSuite(t *testing.T) {
	suite.Run(t, new(OrganizationRelationsSuite))
}

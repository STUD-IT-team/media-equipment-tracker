//go:build integration

package pgorganization_test

import (
	"context"
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
	"media-equipment-tracker/pkg/txmanager/gormtx"
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
	s.orgRepo = pgorganization.NewPostgresOrganizationRepository(gormtx.NewDBGetter(gdb))
	s.userRepo = pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
	s.eqInvocationRepo = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))
	s.studioInvocationRepo = pgstudioinvocation.NewPostgresStudioInvocationRepository(gormtx.NewDBGetter(gdb))
	s.deptRepo = pgdepartment.NewPostgresDepartmentRepository(gormtx.NewDBGetter(gdb))
}

// --- Users (many2many) ---

func (s *OrganizationRelationsSuite) TestUsers_Preload() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	user.Organizations = []*domain.Organization{org}
	s.Require().NoError(s.userRepo.Update(ctx, user))

	user, err := s.userRepo.Get(ctx, user.ID, domain.UserWithOrganizations())
	s.NoError(err)
	s.Len(user.Organizations, 1)
	s.Equal(org.ID, user.Organizations[0].ID)

	// без preload
	raw, _ := s.orgRepo.Get(ctx, org.ID)
	s.Empty(raw.Users)

	// с preload
	with, _ := s.orgRepo.Get(ctx, org.ID, domain.WithUsers())
	s.Len(with.Users, 1)
	s.Equal(user.ID, with.Users[0].ID)
}

func (s *OrganizationRelationsSuite) TestUsers_PreloadAfterDelete() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	user.Organizations = []*domain.Organization{org}
	s.Require().NoError(s.userRepo.Update(ctx, user))
	user.Organizations = []*domain.Organization{}
	s.Require().NoError(s.userRepo.Update(ctx, user))

	// без preload
	raw, _ := s.orgRepo.Get(ctx, org.ID)
	s.Empty(raw.Users)

	// с preload
	with, _ := s.orgRepo.Get(ctx, org.ID, domain.WithUsers())
	s.Empty(with.Users)
}

func (s *OrganizationRelationsSuite) TestUsers_NoAutoCreate() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()

	org.Users = []*domain.User{user}

	// пользователь не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().NoError(s.orgRepo.Create(ctx, org))

	_, err := s.userRepo.Get(ctx, user.ID)
	s.Error(err)
}

func (s *OrganizationRelationsSuite) TestUsers_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.orgRepo.Create(ctx, org))

	org.Users = []*domain.User{user}
	s.Require().NoError(s.userRepo.Update(ctx, user))

	with, _ := s.orgRepo.Get(ctx, org.ID, domain.WithUsers())
	s.Empty(with.Users)
}

// --- EquipmentInvocations (has many) ---

func (s *OrganizationRelationsSuite) TestEquipmentInvocations_Preload() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	invocation := newInvocation(nil, user.ID)

	invocation.OrganizationID = &org.ID

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.userRepo.Create(ctx, user))
	s.Require().NoError(s.eqInvocationRepo.Create(ctx, invocation))

	// без preload
	raw, _ := s.orgRepo.Get(ctx, org.ID)
	s.Empty(raw.EquipmentInvocations)

	// с preload
	with, _ := s.orgRepo.Get(ctx, org.ID, domain.WithEquipmentInvocations())
	s.Len(with.EquipmentInvocations, 1)
	s.Equal(invocation.ID, with.EquipmentInvocations[0].ID)
}

func (s *OrganizationRelationsSuite) TestEquipmentInvocations_NoAutoCreate() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	invocation := newInvocation(nil, user.ID)

	invocation.OrganizationID = &org.ID

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	org.EquipmentInvocations = []*domain.EquipmentInvocation{invocation}

	// invocation не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().NoError(s.orgRepo.Update(ctx, org))

	_, err := s.eqInvocationRepo.Get(ctx, invocation.ID)
	s.Error(err)
}

func (s *OrganizationRelationsSuite) TestEquipmentInvocations_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	inv := newInvocation(nil, user.ID)

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	org.EquipmentInvocations = []*domain.EquipmentInvocation{inv}
	s.Require().NoError(s.userRepo.Update(ctx, user))

	with, _ := s.orgRepo.Get(ctx, org.ID, domain.WithEquipmentInvocations())
	s.Empty(with.EquipmentInvocations)
}

// --- StudioInvocations (has many) ---

func (s *OrganizationRelationsSuite) TestStudioInvocations_Preload() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	studioInvocation := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.userRepo.Create(ctx, user))
	s.Require().NoError(s.studioInvocationRepo.Create(ctx, studioInvocation))

	// без preload
	raw, _ := s.orgRepo.Get(ctx, org.ID)
	s.Empty(raw.StudioInvocations)

	// с preload
	with, _ := s.orgRepo.Get(ctx, org.ID, domain.WithStudioInvocations())
	s.Len(with.StudioInvocations, 1)
	s.Equal(studioInvocation.ID, with.StudioInvocations[0].ID)
}

func (s *OrganizationRelationsSuite) TestStudioInvocations_NoAutoCreate() {
	ctx := context.Background()
	org := newOrg()
	dept := newDept()
	user := newUser()
	studioInvocation := newStudioInvocation(&org.ID, &dept.ID, &user.ID, &user.ID)

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.deptRepo.Create(ctx, dept))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	org.StudioInvocations = []*domain.StudioInvocation{studioInvocation}

	// studioInvocation не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().NoError(s.orgRepo.Update(ctx, org))

	_, err := s.studioInvocationRepo.Get(ctx, studioInvocation.ID)
	s.Error(err)
}

func (s *OrganizationRelationsSuite) TestStudioInvocations_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()
	studioInvocation := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.orgRepo.Create(ctx, org))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	org.StudioInvocations = []*domain.StudioInvocation{studioInvocation}
	s.Require().NoError(s.userRepo.Update(ctx, user))

	with, _ := s.orgRepo.Get(ctx, org.ID, domain.WithStudioInvocations())
	s.Empty(with.StudioInvocations)
}

// ===== run =====

func TestOrganizationRelationsSuite(t *testing.T) {
	suite.Run(t, new(OrganizationRelationsSuite))
}

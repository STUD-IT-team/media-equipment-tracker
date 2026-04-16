//go:build integration

package pguser_test

import (
	"testing"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgorganization"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type UserRelationsSuite struct {
	suite.Suite

	db       *gorm.DB
	pg       *pgtest.PgTestDatabase
	userRepo domain.UserRepository
	orgRepo  domain.OrganizationRepository
	deptRepo domain.DepartmentRepository
}

func (s *UserRelationsSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *UserRelationsSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())
	db := stdlib.OpenDBFromPool(s.pg.Pool())

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	s.Require().NoError(err)

	s.db = gdb
	s.userRepo = pguser.NewPostgresUserRepository(gdb)
	s.orgRepo = pgorganization.NewPostgresOrganizationRepository(gdb)
	s.deptRepo = pgdepartment.NewPostgresDepartmentRepository(gdb)
}

func (s *UserRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

// --- Organizations (many2many) ---

func (s *UserRelationsSuite) TestOrganizations_Preload() {
	u := newUser()
	org := newOrg()

	s.Require().NoError(s.orgRepo.Create(org))

	u.Organizations = []*domain.Organization{org}
	s.Require().NoError(s.userRepo.Create(u))

	// без preload
	raw, _ := s.userRepo.Get(u.ID)
	s.Empty(raw.Organizations)

	// с preload
	with, _ := s.userRepo.Get(u.ID, domain.UserWithOrganizations())
	s.Len(with.Organizations, 1)
	s.Equal(org.ID, with.Organizations[0].ID)
}

func (s *UserRelationsSuite) TestOrganizations_NoAutoCreate() {
	u := newUser()
	org := newOrg()

	u.Organizations = []*domain.Organization{org}

	// ошибка FK
	s.Require().Error(s.userRepo.Create(u))
}

func (s *UserRelationsSuite) TestOrganizations_NoUpdateThroughUser() {
	u := newUser()
	org := newOrg()

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(u))

	u.Organizations = []*domain.Organization{org}
	s.Require().NoError(s.userRepo.Update(u))

	with, _ := s.userRepo.Get(u.ID, domain.UserWithOrganizations())
	s.NotNil(with)

	with.Organizations[0].Name = "HACKED"

	s.Require().NoError(s.userRepo.Update(with))
	got, _ := s.orgRepo.Get(org.ID)
	s.NotEqual("HACKED", got.Name)
}

func (s *UserRelationsSuite) TestOrganizations_NoAutoCreateOnUpdate() {
	u := newUser()
	org := newOrg()

	s.Require().NoError(s.userRepo.Create(u))

	u.Organizations = []*domain.Organization{org}
	// ошибка FK
	s.Require().Error(s.userRepo.Update(u))

	with, _ := s.userRepo.Get(u.ID, domain.UserWithOrganizations())
	s.NotNil(with)
	s.Empty(with.Organizations)
}

// --- Departments (join entity) ---

func (s *UserRelationsSuite) TestDepartments_Relation_And_Reload() {
	u := newUser()
	dep := newDept()

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(u))

	u.Departments = []*domain.UserDepartment{{
		UserID:       u.ID,
		DepartmentID: dep.ID,
		Role:         domain.RoleActivist,
	}}
	s.Require().NoError(s.userRepo.Update(u))

	// без preload
	raw, _ := s.userRepo.Get(u.ID)
	s.Empty(raw.Departments)

	// с preload
	with, _ := s.userRepo.Get(u.ID, domain.UserWithDepartments())
	s.Len(with.Departments, 1)
	s.Equal(dep.ID, with.Departments[0].DepartmentID)
	s.Equal(domain.RoleActivist, with.Departments[0].Role)
}

func (s *UserRelationsSuite) TestDepartments_NoAutoCreate() {
	u := newUser()
	dep := newDept()

	u.Departments = []*domain.UserDepartment{
		{
			UserID:       u.ID,
			DepartmentID: dep.ID,
			Role:         domain.RoleTrainee,
		},
	}

	// FK ошибка
	err := s.userRepo.Create(u)
	s.Error(err)

	_, err = s.deptRepo.Get(dep.ID)
	s.Error(err)
}

func (s *UserRelationsSuite) TestDepartments_NoUpdateThroughUser() {
	u := newUser()
	dep := newDept()

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(u))

	u.Departments = []*domain.UserDepartment{{
		UserID:       u.ID,
		DepartmentID: dep.ID,
		Role:         domain.RoleTrainee,
	}}
	s.Require().NoError(s.userRepo.Update(u))

	with, _ := s.userRepo.Get(u.ID, domain.UserWithDepartments())
	s.NotNil(with)

	with.Departments[0].Department.Name = "HACKED"

	s.Require().NoError(s.userRepo.Update(with))
	got, _ := s.deptRepo.Get(dep.ID)
	s.NotEqual("HACKED", got.Name)
}

// ===== run =====

func TestUserRelationsSuite(t *testing.T) {
	suite.Run(t, new(UserRelationsSuite))
}

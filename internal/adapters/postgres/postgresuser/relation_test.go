//go:build integration

package postgresuser_test

import (
	"testing"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postgresrepo "media-equipment-tracker/internal/adapters/postgres"
	"media-equipment-tracker/internal/adapters/postgres/postgresuser"
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
	s.userRepo = postgresuser.NewPostgresUserRepository(gdb)
	s.orgRepo = postgresrepo.NewOrganizationRepository(gdb)
	s.deptRepo = postgresrepo.NewDepartmentRepository(gdb)
}

func (s *UserRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

// --- Organizations (many2many) ---

func (s *UserRelationsSuite) TestOrganizations_Preload_And_SaveOnlyIDs() {
	u := newUser()
	org := newOrg()

	s.Require().NoError(s.orgRepo.Create(org))

	u.Organizations = []*domain.Organization{org}
	s.Require().NoError(s.userRepo.Create(u))

	// без preload — пусто
	raw, _ := s.userRepo.Get(u.ID)
	s.Empty(raw.Organizations)

	// с preload — есть
	with, _ := s.userRepo.Get(u.ID, domain.UserWithOrganizations())
	s.Len(with.Organizations, 1)
	s.Equal(org.ID, with.Organizations[0].ID)
}

// --- Departments (join entity) ---

func (s *UserRelationsSuite) TestDepartments_Relation_And_Reload() {
	u := newUser()
	dept := newDept()

	s.Require().NoError(s.deptRepo.Create(dept))
	s.Require().NoError(s.userRepo.Create(u))

	link := &domain.UserDepartment{
		UserID:       u.ID,
		DepartmentID: dept.ID,
		Role:         domain.RoleTrainee,
	}

	// preload
	got, _ := s.userRepo.Get(u.ID, domain.UserWithDepartments())
	s.Len(got.Departments, 0)

	// add link
	got.Departments = append(got.Departments, link)
	s.userRepo.Update(got)

	// reload
	u.Departments = nil
	s.Require().NoError(s.userRepo.Reload(u, domain.UserWithDepartments()))
	s.Len(u.Departments, 1)
}

func (s *UserRelationsSuite) TestDepartments_NoAutoCreate() {
	u := newUser()
	dept := newDept()

	u.Departments = []*domain.UserDepartment{
		{
			UserID:       u.ID,
			DepartmentID: dept.ID,
			Role:         domain.RoleTrainee,
		},
	}

	// создаём пользователя — департамент не должен создаться
	s.Require().Error(s.userRepo.Create(u))

	_, err := s.deptRepo.Get(dept.ID)
	s.Error(err) // департамент не появился
}

// ===== run =====

func TestUserRelationsSuite(t *testing.T) {
	suite.Run(t, new(UserRelationsSuite))
}

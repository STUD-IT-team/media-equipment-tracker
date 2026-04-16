package pgstudioinvocation_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgorganization"
	"media-equipment-tracker/internal/adapters/postgres/pgstudioinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type StudioInvocationSuite struct {
	suite.Suite

	db *gorm.DB
	pg *pgtest.PgTestDatabase

	repo domain.StudioInvocationRepository
	org  domain.OrganizationRepository
	dep  domain.DepartmentRepository
	user domain.UserRepository
}

func (s *StudioInvocationSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *StudioInvocationSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *StudioInvocationSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	s.Require().NoError(err)

	s.db = gdb

	s.repo = pgstudioinvocation.NewPostgresStudioInvocationRepository(gdb)
	s.org = pgorganization.NewPostgresOrganizationRepository(gdb)
	s.dep = pgdepartment.NewPostgresDepartmentRepository(gdb)
	s.user = pguser.NewPostgresUserRepository(gdb)
}

func (s *StudioInvocationSuite) TestCreate_Get() {
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.org.Create(org))
	s.Require().NoError(s.user.Create(user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(inv))

	got, err := s.repo.Get(inv.ID)
	s.NoError(err)
	s.Equal(inv.EventName, got.EventName)
}

func (s *StudioInvocationSuite) TestGet_NotFound() {
	_, err := s.repo.Get(uuid.New())
	s.Error(err)
}

func (s *StudioInvocationSuite) TestList() {
	org := newOrg()
	dep := newDept()
	user := newUser()

	s.Require().NoError(s.org.Create(org))
	s.Require().NoError(s.dep.Create(dep))
	s.Require().NoError(s.user.Create(user))

	inv1 := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	inv2 := newStudioInvocation(nil, &dep.ID, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(inv1))
	s.Require().NoError(s.repo.Create(inv2))

	list, err := s.repo.List()
	s.NoError(err)
	s.Len(list, 2)
}

func (s *StudioInvocationSuite) TestUpdate() {
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.org.Create(org))
	s.Require().NoError(s.user.Create(user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(inv))

	inv.EventName = "UPDATED"
	s.Require().NoError(s.repo.Update(inv))

	got, err := s.repo.Get(inv.ID)
	s.NoError(err)
	s.Equal("UPDATED", got.EventName)
}

func (s *StudioInvocationSuite) TestDelete() {
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.org.Create(org))
	s.Require().NoError(s.user.Create(user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(inv))
	s.Require().NoError(s.repo.Delete(inv.ID))

	_, err := s.repo.Get(inv.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestPreload_AllRelations() {
	org := newOrg()
	admin := newUser()
	user := newUser()

	s.Require().NoError(s.org.Create(org))
	s.Require().NoError(s.user.Create(admin))
	s.Require().NoError(s.user.Create(user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &admin.ID)

	s.Require().NoError(s.repo.Create(inv))

	with, err := s.repo.Get(inv.ID,
		domain.StudioInvocationWithOrganization(),
		domain.StudioInvocationWithAdmin(),
		domain.StudioInvocationWithDepartment(),
		domain.StudioInvocationWithUser(),
	)

	s.NoError(err)

	s.Equal(org.ID, with.Organization.ID)
	s.Nil(with.Department)
	s.Equal(admin.ID, with.Admin.ID)
	s.Equal(user.ID, with.User.ID)
}

func (s *StudioInvocationSuite) TestOrganization_NoAutoCreate() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	org := newOrg()

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, nil)

	// organization не существует → FK ошибка
	s.Require().Error(s.repo.Create(inv))

	_, err := s.org.Get(org.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestDepartment_NoAutoCreate() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	dep := newDept()

	inv := newStudioInvocation(nil, &dep.ID, &usr.ID, nil)

	s.Require().Error(s.repo.Create(inv))

	_, err := s.dep.Get(dep.ID)
	s.Error(err)
}
func (s *StudioInvocationSuite) TestUser_NoAutoCreate() {
	org := newOrg()
	s.NoError(s.org.Create(org))

	user := newUser()

	inv := newStudioInvocation(&org.ID, nil, &user.ID, nil)

	s.Require().Error(s.repo.Create(inv))

	_, err := s.user.Get(user.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestAdmin_NoAutoCreate() {
	org := newOrg()
	s.NoError(s.org.Create(org))

	usr := newUser()
	s.NoError(s.user.Create(usr))
	admin := newUser()

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, &admin.ID)

	s.Require().Error(s.repo.Create(inv))

	_, err := s.user.Get(admin.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestOrganization_NoAutoCreateOnUpdate() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(inv))

	org := newOrg()
	inv.OrganizationID = &org.ID
	s.Require().Error(s.repo.Update(inv))

	_, err := s.org.Get(org.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestDepartment_NoAutoCreateOnUpdate() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(inv))

	dep := newDept()
	inv.DepartmentID = &dep.ID
	inv.OrganizationID = nil
	s.Require().Error(s.repo.Update(inv))

	_, err := s.dep.Get(dep.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestUser_NoAutoCreateOnUpdate() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(inv))

	user := newUser()
	inv.UserID = user.ID
	s.Require().Error(s.repo.Update(inv))

	_, err := s.user.Get(user.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestAdmin_NoAutoCreateOnUpdate() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(inv))

	user := newUser()
	inv.AdminID = &user.ID
	s.Require().Error(s.repo.Update(inv))

	_, err := s.user.Get(user.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_OrganizationNotChanged() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	org := newOrg()
	s.Require().NoError(s.org.Create(org))

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(inv))

	with, err := s.repo.Get(inv.ID,
		domain.StudioInvocationWithOrganization(),
	)
	s.Require().NoError(err)

	with.Organization.Name = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.org.Get(org.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.Name)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_DepartmentNotChanged() {
	usr := newUser()
	s.NoError(s.user.Create(usr))

	dep := newDept()
	s.Require().NoError(s.dep.Create(dep))

	inv := newStudioInvocation(nil, &dep.ID, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(inv))

	with, err := s.repo.Get(inv.ID,
		domain.StudioInvocationWithDepartment(),
	)
	s.Require().NoError(err)

	with.Department.Name = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.dep.Get(dep.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.Name)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_UserNotChanged() {
	org := newOrg()
	s.NoError(s.org.Create(org))
	usr := newUser()
	s.NoError(s.user.Create(usr))

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(inv))

	with, err := s.repo.Get(inv.ID,
		domain.StudioInvocationWithUser(),
	)
	s.Require().NoError(err)

	with.User.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.user.Get(usr.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.FullName)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_AdminNotChanged() {
	org := newOrg()
	s.NoError(s.org.Create(org))
	usr := newUser()
	s.NoError(s.user.Create(usr))
	admin := newUser()
	s.Require().NoError(s.user.Create(admin))

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, &admin.ID)

	s.Require().NoError(s.repo.Create(inv))

	with, err := s.repo.Get(inv.ID,
		domain.StudioInvocationWithAdmin(),
	)
	s.Require().NoError(err)

	with.Admin.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(with))

	got, err := s.user.Get(admin.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.FullName)
}

func TestStudioInvocationSuite(t *testing.T) {
	suite.Run(t, new(StudioInvocationSuite))
}

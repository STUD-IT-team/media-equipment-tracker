//go:build integration

package pgstudioinvocation_test

import (
	"context"
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
	"media-equipment-tracker/pkg/txmanager/gormtx"
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

	s.repo = pgstudioinvocation.NewPostgresStudioInvocationRepository(gormtx.NewDBGetter(gdb))
	s.org = pgorganization.NewPostgresOrganizationRepository(gormtx.NewDBGetter(gdb))
	s.dep = pgdepartment.NewPostgresDepartmentRepository(gormtx.NewDBGetter(gdb))
	s.user = pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
}

func (s *StudioInvocationSuite) TestCreate_Get() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.org.Create(ctx, org))
	s.Require().NoError(s.user.Create(ctx, user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(ctx, inv))

	got, err := s.repo.Get(ctx, inv.ID)
	s.NoError(err)
	s.Equal(inv.EventName, got.EventName)
}

func (s *StudioInvocationSuite) TestGet_NotFound() {
	ctx := context.Background()
	_, err := s.repo.Get(ctx, uuid.New())
	s.Error(err)
}

func (s *StudioInvocationSuite) TestList() {
	ctx := context.Background()
	org := newOrg()
	dep := newDept()
	user := newUser()

	s.Require().NoError(s.org.Create(ctx, org))
	s.Require().NoError(s.dep.Create(ctx, dep))
	s.Require().NoError(s.user.Create(ctx, user))

	inv1 := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	inv2 := newStudioInvocation(nil, &dep.ID, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(ctx, inv1))
	s.Require().NoError(s.repo.Create(ctx, inv2))

	list, err := s.repo.List(ctx)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *StudioInvocationSuite) TestUpdate() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.org.Create(ctx, org))
	s.Require().NoError(s.user.Create(ctx, user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(ctx, inv))

	inv.EventName = "UPDATED"
	s.Require().NoError(s.repo.Update(ctx, inv))

	got, err := s.repo.Get(ctx, inv.ID)
	s.NoError(err)
	s.Equal("UPDATED", got.EventName)
}

func (s *StudioInvocationSuite) TestDelete() {
	ctx := context.Background()
	org := newOrg()
	user := newUser()

	s.Require().NoError(s.org.Create(ctx, org))
	s.Require().NoError(s.user.Create(ctx, user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &user.ID)

	s.Require().NoError(s.repo.Create(ctx, inv))
	s.Require().NoError(s.repo.Delete(ctx, inv.ID))

	_, err := s.repo.Get(ctx, inv.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestPreload_AllRelations() {
	ctx := context.Background()
	org := newOrg()
	admin := newUser()
	user := newUser()

	s.Require().NoError(s.org.Create(ctx, org))
	s.Require().NoError(s.user.Create(ctx, admin))
	s.Require().NoError(s.user.Create(ctx, user))

	inv := newStudioInvocation(&org.ID, nil, &user.ID, &admin.ID)

	s.Require().NoError(s.repo.Create(ctx, inv))

	with, err := s.repo.Get(ctx, inv.ID,
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
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	org := newOrg()

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, nil)

	// organization не существует → FK ошибка
	s.Require().Error(s.repo.Create(ctx, inv))

	_, err := s.org.Get(ctx, org.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestDepartment_NoAutoCreate() {
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	dep := newDept()

	inv := newStudioInvocation(nil, &dep.ID, &usr.ID, nil)

	s.Require().Error(s.repo.Create(ctx, inv))

	_, err := s.dep.Get(ctx, dep.ID)
	s.Error(err)
}
func (s *StudioInvocationSuite) TestUser_NoAutoCreate() {
	ctx := context.Background()
	org := newOrg()
	s.NoError(s.org.Create(ctx, org))

	user := newUser()

	inv := newStudioInvocation(&org.ID, nil, &user.ID, nil)

	s.Require().Error(s.repo.Create(ctx, inv))

	_, err := s.user.Get(ctx, user.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestAdmin_NoAutoCreate() {
	ctx := context.Background()
	org := newOrg()
	s.NoError(s.org.Create(ctx, org))

	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))
	admin := newUser()

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, &admin.ID)

	s.Require().Error(s.repo.Create(ctx, inv))

	_, err := s.user.Get(ctx, admin.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestOrganization_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(ctx, someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(ctx, inv))

	org := newOrg()
	inv.OrganizationID = &org.ID
	s.Require().Error(s.repo.Update(ctx, inv))

	_, err := s.org.Get(ctx, org.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestDepartment_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(ctx, someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(ctx, inv))

	dep := newDept()
	inv.DepartmentID = &dep.ID
	inv.OrganizationID = nil
	s.Require().Error(s.repo.Update(ctx, inv))

	_, err := s.dep.Get(ctx, dep.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestUser_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(ctx, someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(ctx, inv))

	user := newUser()
	inv.UserID = user.ID
	s.Require().Error(s.repo.Update(ctx, inv))

	_, err := s.user.Get(ctx, user.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestAdmin_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	someOrg := newOrg()
	someOrg.Name = "some"
	s.Require().NoError(s.org.Create(ctx, someOrg))

	inv := newStudioInvocation(&someOrg.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(ctx, inv))

	user := newUser()
	inv.AdminID = &user.ID
	s.Require().Error(s.repo.Update(ctx, inv))

	_, err := s.user.Get(ctx, user.ID)
	s.Error(err)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_OrganizationNotChanged() {
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	org := newOrg()
	s.Require().NoError(s.org.Create(ctx, org))

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(ctx, inv))

	with, err := s.repo.Get(ctx, inv.ID,
		domain.StudioInvocationWithOrganization(),
	)
	s.Require().NoError(err)

	with.Organization.Name = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.org.Get(ctx, org.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.Name)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_DepartmentNotChanged() {
	ctx := context.Background()
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	dep := newDept()
	s.Require().NoError(s.dep.Create(ctx, dep))

	inv := newStudioInvocation(nil, &dep.ID, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(ctx, inv))

	with, err := s.repo.Get(ctx, inv.ID,
		domain.StudioInvocationWithDepartment(),
	)
	s.Require().NoError(err)

	with.Department.Name = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.dep.Get(ctx, dep.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.Name)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_UserNotChanged() {
	ctx := context.Background()
	org := newOrg()
	s.NoError(s.org.Create(ctx, org))
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, nil)

	s.Require().NoError(s.repo.Create(ctx, inv))

	with, err := s.repo.Get(ctx, inv.ID,
		domain.StudioInvocationWithUser(),
	)
	s.Require().NoError(err)

	with.User.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.user.Get(ctx, usr.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.FullName)
}

func (s *StudioInvocationSuite) TestUpdateThroughStudioInvocation_AdminNotChanged() {
	ctx := context.Background()
	org := newOrg()
	s.NoError(s.org.Create(ctx, org))
	usr := newUser()
	s.NoError(s.user.Create(ctx, usr))
	admin := newUser()
	s.Require().NoError(s.user.Create(ctx, admin))

	inv := newStudioInvocation(&org.ID, nil, &usr.ID, &admin.ID)

	s.Require().NoError(s.repo.Create(ctx, inv))

	with, err := s.repo.Get(ctx, inv.ID,
		domain.StudioInvocationWithAdmin(),
	)
	s.Require().NoError(err)

	with.Admin.FullName = "HACKED"

	s.Require().NoError(s.repo.Update(ctx, with))

	got, err := s.user.Get(ctx, admin.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.FullName)
}

func TestStudioInvocationSuite(t *testing.T) {
	suite.Run(t, new(StudioInvocationSuite))
}

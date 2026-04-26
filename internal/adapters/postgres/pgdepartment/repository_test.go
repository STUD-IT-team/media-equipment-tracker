//go:build integration

package pgdepartment_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgstudioinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
)

type DepartmentRepositorySuite struct {
	suite.Suite

	db       *gorm.DB
	pg       *pgtest.PgTestDatabase
	deptRepo domain.DepartmentRepository
	eqRepo   domain.EquipmentRepository
	userRepo domain.UserRepository

	invocationRepository domain.EquipmentInvocationRepository
	studioInvocationRepo domain.StudioInvocationRepository
}

func (s *DepartmentRepositorySuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *DepartmentRepositorySuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *DepartmentRepositorySuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.deptRepo = pgdepartment.NewPostgresDepartmentRepository(gormtx.NewDBGetter(gdb))
	s.eqRepo = pgequipment.NewPostgresEquipmentRepository(gormtx.NewDBGetter(gdb))
	s.userRepo = pguser.NewPostgresUserRepository(gormtx.NewDBGetter(gdb))
	s.invocationRepository = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gormtx.NewDBGetter(gdb))
	s.studioInvocationRepo = pgstudioinvocation.NewPostgresStudioInvocationRepository(gormtx.NewDBGetter(gdb))
}

func (s *DepartmentRepositorySuite) TestCreate_Get() {
	ctx := context.Background()
	dep := newDepartment()

	s.Require().NoError(s.deptRepo.Create(ctx, dep))

	got, err := s.deptRepo.Get(ctx, dep.ID)
	s.NoError(err)
	s.Equal(dep.Name, got.Name)
}

func (s *DepartmentRepositorySuite) TestGet_NotFound() {
	ctx := context.Background()
	_, err := s.deptRepo.Get(ctx, uuid.New())
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestList() {
	ctx := context.Background()
	d1, d2 := newDepartment(), newDepartment()

	s.NoError(s.deptRepo.Create(ctx, d1))
	s.NoError(s.deptRepo.Create(ctx, d2))

	list, err := s.deptRepo.List(ctx)
	s.NoError(err)
	s.Len(list, 2)
}

func (s *DepartmentRepositorySuite) TestUpdate() {
	ctx := context.Background()
	dep := newDepartment()
	s.NoError(s.deptRepo.Create(ctx, dep))

	dep.Name = "updated"
	s.NoError(s.deptRepo.Update(ctx, dep))

	got, _ := s.deptRepo.Get(ctx, dep.ID)
	s.Equal("updated", got.Name)
}

func (s *DepartmentRepositorySuite) TestDelete() {
	ctx := context.Background()
	dep := newDepartment()
	s.NoError(s.deptRepo.Create(ctx, dep))

	s.NoError(s.deptRepo.Delete(ctx, dep.ID))

	_, err := s.deptRepo.Get(ctx, dep.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestEquipment_Preload() {
	ctx := context.Background()
	dep := newDepartment()
	eq := newEquipment()

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))

	dep.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	// без preload
	raw, _ := s.deptRepo.Get(ctx, dep.ID)
	s.Empty(raw.Equipment)

	// с preload
	with, _ := s.deptRepo.Get(ctx, dep.ID, domain.DepartmentWithEquipment())
	s.Len(with.Equipment, 1)
	s.Equal(eq.ID, with.Equipment[0].ID)
}

func (s *DepartmentRepositorySuite) TestEquipment_PreloadAfterDelete() {
	ctx := context.Background()
	dep := newDepartment()
	eq := newEquipment()

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))

	dep.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	dep.Equipment = []*domain.Equipment{}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	with, _ := s.deptRepo.Get(ctx, dep.ID, domain.DepartmentWithEquipment())
	s.Empty(with.Equipment)
}

func (s *DepartmentRepositorySuite) TestEquipment_PreloadUpdate() {
	ctx := context.Background()
	dep := newDepartment()
	eq := newEquipment()
	eq2 := newEquipment()
	eq2.Name = "another"

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))
	s.Require().NoError(s.eqRepo.Create(ctx, eq2))

	dep.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	dep.Equipment = []*domain.Equipment{eq, eq2}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	with, _ := s.deptRepo.Get(ctx, dep.ID, domain.DepartmentWithEquipment())
	s.Len(with.Equipment, 2)
}

func (s *DepartmentRepositorySuite) TestEquipment_NoAutoCreate() {
	ctx := context.Background()
	dep := newDepartment()
	eq := newEquipment()

	dep.Equipment = []*domain.Equipment{eq}

	// equipment не существует → должна быть ошибка FK
	s.Require().Error(s.deptRepo.Create(ctx, dep))

	_, err := s.eqRepo.Get(ctx, eq.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestEquipment_NoAutoCreateOnUpdate() {
	ctx := context.Background()
	dep := newDepartment()
	s.Require().NoError(s.deptRepo.Create(ctx, dep))

	eq := newEquipment()
	dep.Equipment = []*domain.Equipment{eq}

	// equipment не существует → FK ошибка
	s.Require().Error(s.deptRepo.Update(ctx, dep))

	eq, err := s.eqRepo.Get(ctx, eq.ID)
	s.Error(err)
	s.Nil(eq)
}

func (s *DepartmentRepositorySuite) TestEquipment_Replace() {
	ctx := context.Background()
	dep := newDepartment()
	eq1 := newEquipment()
	eq2 := newEquipment()

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.eqRepo.Create(ctx, eq1))
	s.Require().NoError(s.eqRepo.Create(ctx, eq2))

	dep.Equipment = []*domain.Equipment{eq1}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	dep.Equipment = []*domain.Equipment{eq2}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	with, _ := s.deptRepo.Get(ctx, dep.ID, domain.DepartmentWithEquipment())
	s.Len(with.Equipment, 1)
	s.Equal(eq2.ID, with.Equipment[0].ID)
}

func (s *DepartmentRepositorySuite) TestEquipment_NoUpdateThroughDepartment() {
	ctx := context.Background()
	dep := newDepartment()
	eq := newEquipment()
	srcName := eq.Name

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.eqRepo.Create(ctx, eq))

	dep.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	dep.Equipment[0].Name = "HACKED"

	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	got, err := s.eqRepo.Get(ctx, eq.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.Name)
	s.Equal(srcName, got.Name)
}

func (s *DepartmentRepositorySuite) TestUsers_Preload() {
	ctx := context.Background()
	dep := newDepartment()
	user := newUser()

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	// руками создаём связь (так как repo её не сохраняет)
	err := s.db.Exec(
		"INSERT INTO user_department (user_id, department_id) VALUES (?, ?)",
		user.ID, dep.ID,
	).Error
	s.Require().NoError(err)

	// без preload
	raw, _ := s.deptRepo.Get(ctx, dep.ID)
	s.Empty(raw.Users)

	// с preload
	with, _ := s.deptRepo.Get(ctx, dep.ID, domain.DepartmentWithUsers())
	s.Len(with.Users, 1)
	s.Equal(user.ID, with.Users[0].UserID)
}

func (s *DepartmentRepositorySuite) TestUsers_NoAutoCreate() {
	ctx := context.Background()
	dep := newDepartment()
	user := newUser()

	dep.Users = []*domain.UserDepartment{
		{UserID: user.ID, User: user, DepartmentID: dep.ID, Department: dep, Role: domain.RoleTrainee},
	}

	// не должен создавать user и не должен падать
	s.Require().NoError(s.deptRepo.Create(ctx, dep))

	_, err := s.userRepo.Get(ctx, user.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestEquipmentInvocations_Preload() {
	ctx := context.Background()
	dep := newDepartment()
	user := newUser()
	inv := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.userRepo.Create(ctx, user))
	s.Require().NoError(s.invocationRepository.Create(ctx, inv))

	// без preload
	raw, _ := s.deptRepo.Get(ctx, dep.ID)
	s.Empty(raw.EquipmentInvocations)

	// с preload
	with, _ := s.deptRepo.Get(ctx, dep.ID, domain.DepartmentWithEquipmentInvocations())
	s.Len(with.EquipmentInvocations, 1)
	s.Equal(inv.ID, with.EquipmentInvocations[0].ID)
}

func (s *DepartmentRepositorySuite) TestEquipmentInvocations_NoAutoCreate() {
	ctx := context.Background()
	dep := newDepartment()
	user := newUser()
	inv := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	dep.EquipmentInvocations = []*domain.EquipmentInvocation{inv}

	// invocation не должен создаться
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	_, err := s.invocationRepository.Get(ctx, inv.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestStudioInvocations_Preload() {
	ctx := context.Background()
	dep := newDepartment()
	user := newUser()
	inv := newStudioInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.userRepo.Create(ctx, user))
	s.Require().NoError(s.studioInvocationRepo.Create(ctx, inv))

	// без preload
	raw, _ := s.deptRepo.Get(ctx, dep.ID)
	s.Empty(raw.StudioInvocations)

	// с preload
	with, _ := s.deptRepo.Get(ctx, dep.ID, domain.DepartmentWithStudioInvocations())
	s.Len(with.StudioInvocations, 1)
	s.Equal(inv.ID, with.StudioInvocations[0].ID)
}

func (s *DepartmentRepositorySuite) TestStudioInvocations_NoAutoCreate() {
	ctx := context.Background()
	dep := newDepartment()
	user := newUser()
	inv := newStudioInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(ctx, dep))
	s.Require().NoError(s.userRepo.Create(ctx, user))

	dep.StudioInvocations = []*domain.StudioInvocation{inv}

	// не должен создавать invocation
	s.Require().NoError(s.deptRepo.Update(ctx, dep))

	_, err := s.studioInvocationRepo.Get(ctx, inv.ID)
	s.Error(err)
}

func TestDepartmentRepositorySuite(t *testing.T) {
	suite.Run(t, new(DepartmentRepositorySuite))
}

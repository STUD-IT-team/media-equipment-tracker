//go:build integration

package pgdepartment_test

import (
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
	s.deptRepo = pgdepartment.NewPostgresDepartmentRepository(gdb)
	s.eqRepo = pgequipment.NewPostgresEquipmentRepository(gdb)
	s.userRepo = pguser.NewPostgresUserRepository(gdb)
	s.invocationRepository = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gdb)
	s.studioInvocationRepo = pgstudioinvocation.NewPostgresStudioInvocationRepository(gdb)
}

func (s *DepartmentRepositorySuite) TestCreate_Get() {
	dep := newDepartment()

	s.Require().NoError(s.deptRepo.Create(dep))

	got, err := s.deptRepo.Get(dep.ID)
	s.NoError(err)
	s.Equal(dep.Name, got.Name)
}

func (s *DepartmentRepositorySuite) TestGet_NotFound() {
	_, err := s.deptRepo.Get(uuid.New())
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestList() {
	d1, d2 := newDepartment(), newDepartment()

	s.NoError(s.deptRepo.Create(d1))
	s.NoError(s.deptRepo.Create(d2))

	list, err := s.deptRepo.List()
	s.NoError(err)
	s.Len(list, 2)
}

func (s *DepartmentRepositorySuite) TestUpdate() {
	dep := newDepartment()
	s.NoError(s.deptRepo.Create(dep))

	dep.Name = "updated"
	s.NoError(s.deptRepo.Update(dep))

	got, _ := s.deptRepo.Get(dep.ID)
	s.Equal("updated", got.Name)
}

func (s *DepartmentRepositorySuite) TestDelete() {
	dep := newDepartment()
	s.NoError(s.deptRepo.Create(dep))

	s.NoError(s.deptRepo.Delete(dep.ID))

	_, err := s.deptRepo.Get(dep.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestEquipment_Preload() {
	dep := newDepartment()
	eq := newEquipment()

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.eqRepo.Create(eq))

	dep.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(dep))

	// без preload
	raw, _ := s.deptRepo.Get(dep.ID)
	s.Empty(raw.Equipment)

	// с preload
	with, _ := s.deptRepo.Get(dep.ID, domain.DepartmentWithEquipment())
	s.Len(with.Equipment, 1)
	s.Equal(eq.ID, with.Equipment[0].ID)
}

func (s *DepartmentRepositorySuite) TestEquipment_PreloadAfterDelete() {
	dep := newDepartment()
	eq := newEquipment()

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.eqRepo.Create(eq))

	dep.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(dep))

	dep.Equipment = []*domain.Equipment{}
	s.Require().NoError(s.deptRepo.Update(dep))

	with, _ := s.deptRepo.Get(dep.ID, domain.DepartmentWithEquipment())
	s.Empty(with.Equipment)
}

func (s *DepartmentRepositorySuite) TestEquipment_NoAutoCreate() {
	dep := newDepartment()
	eq := newEquipment()

	dep.Equipment = []*domain.Equipment{eq}

	// equipment не существует → должна быть ошибка FK
	s.Require().Error(s.deptRepo.Create(dep))

	_, err := s.eqRepo.Get(eq.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestEquipment_NoAutoCreateOnUpdate() {
	dep := newDepartment()
	s.Require().NoError(s.deptRepo.Create(dep))

	eq := newEquipment()
	dep.Equipment = []*domain.Equipment{eq}

	// equipment не существует → FK ошибка
	s.Require().Error(s.deptRepo.Update(dep))
}

func (s *DepartmentRepositorySuite) TestEquipment_Replace() {
	dep := newDepartment()
	eq1 := newEquipment()
	eq2 := newEquipment()

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.eqRepo.Create(eq1))
	s.Require().NoError(s.eqRepo.Create(eq2))

	dep.Equipment = []*domain.Equipment{eq1}
	s.Require().NoError(s.deptRepo.Update(dep))

	dep.Equipment = []*domain.Equipment{eq2}
	s.Require().NoError(s.deptRepo.Update(dep))

	with, _ := s.deptRepo.Get(dep.ID, domain.DepartmentWithEquipment())
	s.Len(with.Equipment, 1)
	s.Equal(eq2.ID, with.Equipment[0].ID)
}

func (s *DepartmentRepositorySuite) TestEquipment_NoUpdateThroughDepartment() {
	dep := newDepartment()
	eq := newEquipment()
	srcName := eq.Name

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.eqRepo.Create(eq))

	dep.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(dep))

	dep.Equipment[0].Name = "HACKED"

	s.Require().NoError(s.deptRepo.Update(dep))

	got, err := s.eqRepo.Get(eq.ID)
	s.Require().NoError(err)

	s.NotEqual("HACKED", got.Name)
	s.Equal(srcName, got.Name)
}

func (s *DepartmentRepositorySuite) TestUsers_Preload() {
	dep := newDepartment()
	user := newUser()

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))

	// руками создаём связь (так как repo её не сохраняет)
	err := s.db.Exec(
		"INSERT INTO user_department (user_id, department_id) VALUES (?, ?)",
		user.ID, dep.ID,
	).Error
	s.Require().NoError(err)

	// без preload
	raw, _ := s.deptRepo.Get(dep.ID)
	s.Empty(raw.Users)

	// с preload
	with, _ := s.deptRepo.Get(dep.ID, domain.DepartmentWithUsers())
	s.Len(with.Users, 1)
	s.Equal(user.ID, with.Users[0].ID)
}

func (s *DepartmentRepositorySuite) TestUsers_NoAutoCreate() {
	dep := newDepartment()
	user := newUser()

	dep.Users = []*domain.User{user}

	// не должен создавать user и не должен падать
	s.Require().NoError(s.deptRepo.Create(dep))

	_, err := s.userRepo.Get(user.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestEquipmentInvocations_Preload() {
	dep := newDepartment()
	user := newUser()
	inv := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.invocationRepository.Create(inv))

	// без preload
	raw, _ := s.deptRepo.Get(dep.ID)
	s.Empty(raw.EquipmentInvocations)

	// с preload
	with, _ := s.deptRepo.Get(dep.ID, domain.DepartmentWithEquipmentInvocations())
	s.Len(with.EquipmentInvocations, 1)
	s.Equal(inv.ID, with.EquipmentInvocations[0].ID)
}

func (s *DepartmentRepositorySuite) TestEquipmentInvocations_NoAutoCreate() {
	dep := newDepartment()
	user := newUser()
	inv := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))

	dep.EquipmentInvocations = []*domain.EquipmentInvocation{inv}

	// invocation не должен создаться
	s.Require().NoError(s.deptRepo.Update(dep))

	_, err := s.invocationRepository.Get(inv.ID)
	s.Error(err)
}

func (s *DepartmentRepositorySuite) TestStudioInvocations_Preload() {
	dep := newDepartment()
	user := newUser()
	inv := newStudioInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.studioInvocationRepo.Create(inv))

	// без preload
	raw, _ := s.deptRepo.Get(dep.ID)
	s.Empty(raw.StudioInvocations)

	// с preload
	with, _ := s.deptRepo.Get(dep.ID, domain.DepartmentWithStudioInvocations())
	s.Len(with.StudioInvocations, 1)
	s.Equal(inv.ID, with.StudioInvocations[0].ID)
}

func (s *DepartmentRepositorySuite) TestStudioInvocations_NoAutoCreate() {
	dep := newDepartment()
	user := newUser()
	inv := newStudioInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))

	dep.StudioInvocations = []*domain.StudioInvocation{inv}

	// не должен создавать invocation
	s.Require().NoError(s.deptRepo.Update(dep))

	_, err := s.studioInvocationRepo.Get(inv.ID)
	s.Error(err)
}

func TestDepartmentRepositorySuite(t *testing.T) {
	suite.Run(t, new(DepartmentRepositorySuite))
}

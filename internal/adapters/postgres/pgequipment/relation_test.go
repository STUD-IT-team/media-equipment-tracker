package pgequipment_test

import (
	"testing"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type EquipmentRelationsSuite struct {
	suite.Suite

	db                   *gorm.DB
	pg                   *pgtest.PgTestDatabase
	eqRepo               domain.EquipmentRepository
	deptRepo             domain.DepartmentRepository
	invocationRepository domain.EquipmentInvocationRepository
	userRepo             domain.UserRepository
}

func (s *EquipmentRelationsSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *EquipmentRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *EquipmentRelationsSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.eqRepo = pgequipment.NewPostgresEquipmentRepository(gdb)
	s.deptRepo = pgdepartment.NewPostgresDepartmentRepository(gdb)
	s.invocationRepository = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gdb)
	s.userRepo = pguser.NewPostgresUserRepository(gdb)
}

// --- Departments (many2many) ---

func (s *EquipmentRelationsSuite) TestDepartments_Preload() {
	eq := newEquipment()
	dept := newDepartment()

	s.Require().NoError(s.deptRepo.Create(dept))
	s.Require().NoError(s.eqRepo.Create(eq))

	dept.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(dept))

	// без preload
	raw, _ := s.eqRepo.Get(eq.ID)
	s.Empty(raw.Departments)

	// с preload
	with, _ := s.eqRepo.Get(eq.ID, domain.EquipmentWithDepartments())
	s.Len(with.Departments, 1)
	s.Equal(dept.ID, with.Departments[0].ID)
}

func (s *EquipmentRelationsSuite) TestDepartments_PreloadAfterDelete() {
	eq := newEquipment()
	dept := newDepartment()

	s.Require().NoError(s.deptRepo.Create(dept))
	s.Require().NoError(s.eqRepo.Create(eq))

	dept.Equipment = []*domain.Equipment{eq}
	s.Require().NoError(s.deptRepo.Update(dept))
	dept.Equipment = []*domain.Equipment{}
	s.Require().NoError(s.deptRepo.Update(dept))

	// без preload
	raw, _ := s.eqRepo.Get(eq.ID)
	s.Empty(raw.Departments)

	// с preload
	with, _ := s.eqRepo.Get(eq.ID, domain.EquipmentWithDepartments())
	s.Empty(with.Departments)
}

func (s *EquipmentRelationsSuite) TestDepartments_NoAutoCreate() {
	eq := newEquipment()
	dept := newDepartment()

	eq.Departments = []*domain.Department{dept}

	// департамент не должен создаться автоматически\
	// так как модель не управляет их связью
	s.Require().NoError(s.eqRepo.Create(eq))

	_, err := s.deptRepo.Get(dept.ID)
	s.Error(err)
}

// --- Invocations (many2many) ---

func (s *EquipmentRelationsSuite) TestInvocations_Preload() {
	eq := newEquipment()
	user := newUser()
	dep := newDepartment()
	invocation := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.eqRepo.Create(eq))

	invocation.Equipment = []*domain.EquipmentInInvocation{
		{
			Invocation:   invocation,
			Equipment:    eq,
			InvocationID: invocation.ID,
			EquipmentID:  eq.ID,
			Status:       domain.EquipmentNotIssued,
		},
	}

	s.Require().NoError(s.invocationRepository.Create(invocation))

	// без preload
	raw, _ := s.eqRepo.Get(eq.ID)
	s.Empty(raw.Invocations)

	// с preload
	with, _ := s.eqRepo.Get(eq.ID, domain.EquipmentWithEquipmentInInvocations())
	s.Len(with.Invocations, 1)
	s.Equal(invocation.ID, with.Invocations[0].InvocationID)
}

func (s *EquipmentRelationsSuite) TestInvocations_NoAutoCreate() {
	eq := newEquipment()
	user := newUser()
	dep := newDepartment()
	invocation := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))

	eq.Invocations = []*domain.EquipmentInInvocation{
		{
			Invocation:   invocation,
			Equipment:    eq,
			InvocationID: invocation.ID,
			EquipmentID:  eq.ID,
			Status:       domain.EquipmentNotIssued,
		},
	}

	// invocation не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().NoError(s.eqRepo.Create(eq))

	_, err := s.invocationRepository.Get(invocation.ID)
	s.Error(err)
}

func (s *EquipmentRelationsSuite) TestCurrentInvocation_Preload() {
	eq := newEquipment()
	user := newUser()
	dep := newDepartment()
	invocation := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.invocationRepository.Create(invocation))

	eq.CurrentInvocationID = &invocation.ID
	s.Require().NoError(s.eqRepo.Create(eq))

	// без preload
	raw, _ := s.eqRepo.Get(eq.ID)
	s.Nil(raw.CurrentInvocation)

	// с preload
	with, _ := s.eqRepo.Get(eq.ID, domain.EquipmentWithCurrentInvocation())
	s.NotNil(with.CurrentInvocation)
	s.Equal(invocation.ID, with.CurrentInvocation.ID)
}

func (s *EquipmentRelationsSuite) TestCurrentInvocation_PreloadAfterUnset() {
	eq := newEquipment()
	user := newUser()
	dep := newDepartment()
	invocation := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.invocationRepository.Create(invocation))

	eq.CurrentInvocationID = &invocation.ID
	s.Require().NoError(s.eqRepo.Create(eq))

	eq.CurrentInvocationID = nil
	s.Require().NoError(s.eqRepo.Update(eq))

	with, _ := s.eqRepo.Get(eq.ID, domain.EquipmentWithCurrentInvocation())
	s.Nil(with.CurrentInvocation)
}

func (s *EquipmentRelationsSuite) TestCurrentInvocation_NoAutoCreate() {
	eq := newEquipment()
	user := newUser()
	dep := newDepartment()
	invocation := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))

	eq.CurrentInvocationID = &invocation.ID
	eq.CurrentInvocation = invocation

	// invocation не должен создаться автоматически
	// так как модель не управляет их связью
	// Но ошибка из-за FK CurrentInvocationID
	s.Require().Error(s.eqRepo.Create(eq))

	_, err := s.invocationRepository.Get(invocation.ID)
	s.Error(err)
}

func (s *EquipmentRelationsSuite) TestCurrentInvocation_NoUpdateThroughEquipment() {
	eq := newEquipment()
	user := newUser()
	dep := newDepartment()
	invocation := newInvocation(dep.ID, user.ID)

	s.Require().NoError(s.deptRepo.Create(dep))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.invocationRepository.Create(invocation))

	eq.CurrentInvocationID = &invocation.ID
	s.Require().NoError(s.eqRepo.Create(eq))

	// preload
	with, _ := s.eqRepo.Get(eq.ID, domain.EquipmentWithCurrentInvocation())

	with.CurrentInvocation.EventName = "HACKED"

	s.Require().NoError(s.eqRepo.Update(with))

	got, err := s.invocationRepository.Get(invocation.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.EventName)
}

// ===== run =====

func TestEquipmentRelationsSuite(t *testing.T) {
	suite.Run(t, new(EquipmentRelationsSuite))
}

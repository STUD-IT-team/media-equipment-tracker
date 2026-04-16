package pgequipmentinvocation_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentininvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgorganization"
	"media-equipment-tracker/internal/adapters/postgres/pguser"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/pkg/pgtest"
)

type EquipmentInvocationRelationsSuite struct {
	suite.Suite

	db                        *gorm.DB
	pg                        *pgtest.PgTestDatabase
	invocationRepo            domain.EquipmentInvocationRepository
	orgRepo                   domain.OrganizationRepository
	deptRepo                  domain.DepartmentRepository
	userRepo                  domain.UserRepository
	equipmentRepo             domain.EquipmentRepository
	equipmentInInvocationRepo domain.EquipmentInInvocationRepository
}

func (s *EquipmentInvocationRelationsSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())
	s.Require().NoError(pg.CreateTemplate())
}

func (s *EquipmentInvocationRelationsSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *EquipmentInvocationRelationsSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	s.Require().NoError(err)

	s.db = gdb
	s.invocationRepo = pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(gdb)
	s.orgRepo = pgorganization.NewPostgresOrganizationRepository(gdb)
	s.deptRepo = pgdepartment.NewPostgresDepartmentRepository(gdb)
	s.userRepo = pguser.NewPostgresUserRepository(gdb)
	s.equipmentRepo = pgequipment.NewPostgresEquipmentRepository(gdb)
	s.equipmentInInvocationRepo = pgequipmentininvocation.NewPostgresEquipmentInInvocationRepository(gdb)
}

// --- Organization (belongs to) ---

func (s *EquipmentInvocationRelationsSuite) TestOrganization_Preload() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// без preload
	raw, _ := s.invocationRepo.Get(invocation.ID)
	s.Nil(raw.Organization)

	// с preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithOrganization())
	s.NotNil(with.Organization)
	s.Equal(org.ID, with.Organization.ID)
}

func (s *EquipmentInvocationRelationsSuite) TestOrganization_NoAutoCreate() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)

	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))

	invocation.OrganizationID = &org.ID
	invocation.Organization = org

	// organization не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().Error(s.invocationRepo.Create(invocation))

	_, err := s.orgRepo.Get(org.ID)
	s.Error(err)
}

func (s *EquipmentInvocationRelationsSuite) TestOrganization_NoUpdateThroughInvocation() {
	org := newOrg()
	user := newUser()
	admin := newUser()

	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithOrganization())

	with.Organization.Name = "HACKED"

	s.Require().NoError(s.invocationRepo.Update(with))

	got, err := s.orgRepo.Get(org.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.Name)
}

// --- Department (belongs to) ---

func (s *EquipmentInvocationRelationsSuite) TestDepartment_Preload() {
	dept := newDept()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(nil, &dept.ID, user.ID, &admin.ID)

	s.Require().NoError(s.deptRepo.Create(dept))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// без preload
	raw, _ := s.invocationRepo.Get(invocation.ID)
	s.Nil(raw.Department)

	// с preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithDepartment())
	s.NotNil(with.Department)
	s.Equal(dept.ID, with.Department.ID)
}

func (s *EquipmentInvocationRelationsSuite) TestDepartment_NoAutoCreate() {
	dept := newDept()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(nil, &dept.ID, user.ID, &admin.ID)

	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))

	invocation.DepartmentID = &dept.ID
	invocation.Department = dept

	// department не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().Error(s.invocationRepo.Create(invocation))

	_, err := s.deptRepo.Get(dept.ID)
	s.Error(err)
}

func (s *EquipmentInvocationRelationsSuite) TestDepartment_NoUpdateThroughInvocation() {
	dept := newDept()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(nil, &dept.ID, user.ID, &admin.ID)

	s.Require().NoError(s.deptRepo.Create(dept))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithDepartment())

	with.Department.Name = "HACKED"

	s.Require().NoError(s.invocationRepo.Update(with))

	got, err := s.deptRepo.Get(dept.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.Name)
}

// --- User (belongs to) ---

func (s *EquipmentInvocationRelationsSuite) TestUser_Preload() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// без preload
	raw, _ := s.invocationRepo.Get(invocation.ID)
	s.Nil(raw.User)

	// с preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithUser())
	s.NotNil(with.User)
	s.Equal(user.ID, with.User.ID)
}

func (s *EquipmentInvocationRelationsSuite) TestUser_NoAutoCreate() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)

	s.Require().NoError(s.orgRepo.Create(org))

	invocation.User = user

	// user не должен создаться автоматически
	// так как модель не управляет их связью
	// Но ошибка из-за отсутствия UserID
	s.Require().Error(s.invocationRepo.Create(invocation))

	_, err := s.userRepo.Get(user.ID)
	s.Error(err)
}

func (s *EquipmentInvocationRelationsSuite) TestUser_NoUpdateThroughInvocation() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithUser())

	with.User.FullName = "HACKED"

	s.Require().NoError(s.invocationRepo.Update(with))

	got, err := s.userRepo.Get(user.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// --- Admin (belongs to) ---

func (s *EquipmentInvocationRelationsSuite) TestAdmin_Preload() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)
	invocation.AdminID = &admin.ID

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// без preload
	raw, _ := s.invocationRepo.Get(invocation.ID)
	s.Nil(raw.Admin)

	// с preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithAdmin())
	s.NotNil(with.Admin)
	s.Equal(admin.ID, with.Admin.ID)
}

func (s *EquipmentInvocationRelationsSuite) TestAdmin_NoAutoCreate() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)
	invocation.AdminID = &admin.ID

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))

	invocation.Admin = admin

	// admin не должен создаться автоматически
	// так как модель не управляет их связью
	s.Require().Error(s.invocationRepo.Create(invocation))

	_, err := s.userRepo.Get(admin.ID)
	s.Error(err)
}

func (s *EquipmentInvocationRelationsSuite) TestAdmin_NoUpdateThroughInvocation() {
	org := newOrg()
	user := newUser()
	admin := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, &admin.ID)
	invocation.AdminID = &admin.ID

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.userRepo.Create(admin))
	s.Require().NoError(s.invocationRepo.Create(invocation))

	// preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithAdmin())

	with.Admin.FullName = "HACKED"

	s.Require().NoError(s.invocationRepo.Update(with))

	got, err := s.userRepo.Get(admin.ID)
	s.Require().NoError(err)
	s.NotEqual("HACKED", got.FullName)
}

// --- EquipmentInInvocation (has many) ---

func (s *EquipmentInvocationRelationsSuite) TestEquipmentInInvocation_Preload() {
	org := newOrg()
	user := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, nil)
	eq := newEquipment()
	eqInInv := &domain.EquipmentInInvocation{
		InvocationID: invocation.ID,
		EquipmentID:  eq.ID,
		Status:       domain.EquipmentNotIssued,
	}

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.invocationRepo.Create(invocation))
	s.Require().NoError(s.equipmentRepo.Create(eq))
	s.Require().NoError(s.equipmentInInvocationRepo.Create(eqInInv))

	// без preload
	raw, _ := s.invocationRepo.Get(invocation.ID)
	s.Empty(raw.Equipment)

	// с preload
	with, _ := s.invocationRepo.Get(invocation.ID, domain.EquipmentInvocationWithEquipment())
	s.Len(with.Equipment, 1)
	s.Equal(invocation.ID, with.Equipment[0].InvocationID)
	s.Equal(eqInInv.EquipmentID, with.Equipment[0].EquipmentID)
}

func (s *EquipmentInvocationRelationsSuite) TestEquipmentInInvocation_Create() {
	org := newOrg()
	user := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, nil)
	eq := newEquipment()
	eqInInv := &domain.EquipmentInInvocation{
		InvocationID: invocation.ID,
		EquipmentID:  eq.ID,
		Status:       domain.EquipmentNotIssued,
	}

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.equipmentRepo.Create(eq))
	invocation.Equipment = []*domain.EquipmentInInvocation{eqInInv}

	// Может создаться со связью, так equipment создан
	s.Require().NoError(s.invocationRepo.Create(invocation))

	_, err := s.equipmentInInvocationRepo.Get(invocation.ID, eqInInv.EquipmentID)
	s.NoError(err)
}

func (s *EquipmentInvocationRelationsSuite) TestEquipmentInInvocation_Update() {
	org := newOrg()
	user := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, nil)
	eq := newEquipment()
	eqInInv := &domain.EquipmentInInvocation{
		InvocationID: invocation.ID,
		EquipmentID:  eq.ID,
		Status:       domain.EquipmentNotIssued,
	}
	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	s.Require().NoError(s.equipmentRepo.Create(eq))

	s.Require().NoError(s.invocationRepo.Create(invocation))

	invocation.Equipment = []*domain.EquipmentInInvocation{eqInInv}
	// Может создаться со связью, так equipment создан
	s.Require().NoError(s.equipmentInInvocationRepo.Update(eqInInv))

	got, _ := s.equipmentInInvocationRepo.Get(invocation.ID, eqInInv.EquipmentID)
	s.Equal(domain.EquipmentNotIssued, got.Status)
}

func (s *EquipmentInvocationRelationsSuite) TestEquipmentInInvocation_NoAutoCreate() {
	org := newOrg()
	user := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, nil)
	eqInInv := &domain.EquipmentInInvocation{
		InvocationID: invocation.ID,
		EquipmentID:  uuid.New(),
		Status:       domain.EquipmentNotIssued,
	}

	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))
	invocation.Equipment = []*domain.EquipmentInInvocation{eqInInv}

	// FK ошибка
	s.Require().Error(s.invocationRepo.Create(invocation))

	_, err := s.equipmentInInvocationRepo.Get(invocation.ID, eqInInv.EquipmentID)
	s.Error(err)
}

func (s *EquipmentInvocationRelationsSuite) TestEquipmentInInvocation_NoAutoCreateOnUpdate() {
	org := newOrg()
	user := newUser()
	invocation := newInvocation(&org.ID, nil, user.ID, nil)
	eqInInv := &domain.EquipmentInInvocation{
		InvocationID: invocation.ID,
		EquipmentID:  uuid.New(),
		Status:       domain.EquipmentNotIssued,
	}
	s.Require().NoError(s.orgRepo.Create(org))
	s.Require().NoError(s.userRepo.Create(user))

	s.Require().NoError(s.invocationRepo.Create(invocation))

	invocation.Equipment = []*domain.EquipmentInInvocation{eqInInv}
	// FK ошибка
	s.Require().Error(s.invocationRepo.Update(invocation))

	_, err := s.equipmentInInvocationRepo.Get(invocation.ID, eqInInv.EquipmentID)
	s.Error(err)
}

// ===== run =====

func TestEquipmentInvocationRelationsSuite(t *testing.T) {
	suite.Run(t, new(EquipmentInvocationRelationsSuite))
}

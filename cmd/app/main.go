package main

import (
	"fmt"

	"media-equipment-tracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"media-equipment-tracker/pkg/txmanager/gormtx"

	"media-equipment-tracker/cmd/app/config"

	"media-equipment-tracker/internal/adapters/bcrypthasher"
	"media-equipment-tracker/internal/adapters/inmem"
	"media-equipment-tracker/internal/adapters/postgres/pgdepartment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipment"
	"media-equipment-tracker/internal/adapters/postgres/pgequipmentinvocation"
	"media-equipment-tracker/internal/adapters/postgres/pgorganization"
	"media-equipment-tracker/internal/adapters/postgres/pguser"

	jwt "media-equipment-tracker/internal/adapters/jwt"
	"media-equipment-tracker/internal/application/accessservice"
	authuser "media-equipment-tracker/internal/application/authservice"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/application/departmentservice"
	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/application/organizationservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/handlers"
	"media-equipment-tracker/internal/handlers/authapi"
	"media-equipment-tracker/internal/handlers/departmentapi"
	"media-equipment-tracker/internal/handlers/equipmentapi"
	"media-equipment-tracker/internal/handlers/invocationapi"
	"media-equipment-tracker/internal/handlers/organizationapi"
	"media-equipment-tracker/internal/middleware"
)

func main() {
	engine := gin.New()

	logger.InitLogger()
	engine.Use(middleware.LoggerMiddleware())

	dbGetter, txManager, err := gormtx.New(
		gormtx.WithHost(config.PostgresHost),
		gormtx.WithPort(uint16(config.PostgresPort)),
		gormtx.WithUser(config.PostgresUser),
		gormtx.WithPassword(config.PostgresPassword),
		gormtx.WithDatabase(config.PostgresDatabase),
	)
	if err != nil {
		logrus.Fatal("Failed to connect to database:", err)
	}

	// Repository
	userRepo := pguser.NewPostgresUserRepository(dbGetter)
	equipmentRepo := pgequipment.NewPostgresEquipmentRepository(dbGetter)
	departmentRepo := pgdepartment.NewPostgresDepartmentRepository(dbGetter)
	organizationRepo := pgorganization.NewPostgresOrganizationRepository(dbGetter)
	invocationRepo := pgequipmentinvocation.NewPostgresEquipmentInvocationRepository(dbGetter)

	// Auth
	authZ := authzservice.NewAuthZ()
	tokenMaker, err := jwt.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		panic(err.Error())
	}
	hasher, err := bcrypthasher.NewBcryptHasher()
	if err != nil {
		panic(err.Error())
	}
	tokenRep := inmem.NewTokenRepository()
	authUserServ, err := authuser.NewAuthUser(tokenMaker, hasher, config.AccessTokenDuration, userRepo, tokenRep)
	if err != nil {
		panic(err.Error())
	}

	// Services
	accessService := accessservice.NewAccessService(authZ)
	equipmentService := equipmentservice.NewEquipmentService(authZ, equipmentRepo, equipmentRepo, invocationRepo, departmentRepo, txManager)
	departmentService := departmentservice.NewDepartmentService(authZ, departmentRepo, userRepo, equipmentRepo, txManager)
	organizationService := organizationservice.NewOrganizationService(authZ, organizationRepo, txManager)
	invocationService := invocationservice.NewInvocationService(invocationRepo, invocationRepo, departmentRepo, organizationRepo, equipmentService, equipmentService, accessService, txManager, authZ)

	// Groups
	healthRouter := handlers.NewHealthRouter(engine.Group("/"))
	_ = healthRouter

	apiGroup := engine.Group(config.APIVersion)
	usersGroup := apiGroup.Group("/")
	usersGroup.Use(middleware.AuthMiddleware(authZ, tokenRep, authUserServ, []domain.RoleAuth{}))

	adminsGroup := apiGroup.Group("/")
	adminsGroup.Use(middleware.AuthMiddleware(authZ, tokenRep, authUserServ, []domain.RoleAuth{domain.AdminRole}))

	// Routers
	authUserRouter := authapi.NewRouter(apiGroup, authUserServ)
	_ = authUserRouter

	// Departments
	departmentRouter := departmentapi.NewRouter(usersGroup, departmentService)
	_ = departmentRouter

	// Organizations
	organizationRouter := organizationapi.NewRouter(usersGroup, organizationService)
	_ = organizationRouter

	// Equipment
	equipmentRouter := equipmentapi.NewRouter(usersGroup, equipmentService)
	_ = equipmentRouter

	// Invocation
	invocationRouter := invocationapi.NewRouter(usersGroup, invocationService)
	_ = invocationRouter

	if err := engine.Run(fmt.Sprintf(":%d", config.AppPort)); err != nil {
		panic(err.Error())
	}
}

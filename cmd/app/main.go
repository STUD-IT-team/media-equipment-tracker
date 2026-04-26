package main

import (
	"fmt"
	"log"
	"media-equipment-tracker/internal/application/userservice"
	"media-equipment-tracker/internal/handlers/userapi"

	"github.com/gin-gonic/gin"

	"media-equipment-tracker/pkg/txmanager/gormtx"

	"media-equipment-tracker/cmd/app/config"

	"media-equipment-tracker/internal/adapters/bcrypthasher"
	"media-equipment-tracker/internal/adapters/inmem"
	"media-equipment-tracker/internal/adapters/postgres/pguser"

	jwt "media-equipment-tracker/internal/adapters/jwt"
	authuser "media-equipment-tracker/internal/application/authservice"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/handlers"
	"media-equipment-tracker/internal/handlers/authapi"
	"media-equipment-tracker/internal/middleware"
)

func main() {
	engine := gin.New()

	dbGetter, _, err := gormtx.New(
		gormtx.WithHost(config.PostgresHost),
		gormtx.WithPort(uint16(config.PostgresPort)),
		gormtx.WithUser(config.PostgresUser),
		gormtx.WithPassword(config.PostgresPassword),
		gormtx.WithDatabase(config.PostgresDatabase),
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Repository
	userRepo := pguser.NewPostgresUserRepository(dbGetter)

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
	userServ, err := userservice.NewUserService(userRepo, authZ)
	if err != nil {
		panic(err.Error())
	}

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
	userRouter := userapi.NewUserRouter(usersGroup, userServ)
	_ = userRouter
	userRouterForAdmin := userapi.NewUserRouterForAdmin(adminsGroup, userServ)
	_ = userRouterForAdmin

	if err := engine.Run(fmt.Sprintf(":%d", config.AppPort)); err != nil {
		panic(err.Error())
	}
}

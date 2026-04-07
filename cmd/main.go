package main

import (
	"fmt"
	"log"
	"media-equipment-tracker/internal/adapters/postgres_repo"
	"media-equipment-tracker/internal/application/auth_service/auth_user"
	"media-equipment-tracker/internal/application/auth_service/hasher"
	tokenmaker "media-equipment-tracker/internal/application/auth_service/token_maker"
	"media-equipment-tracker/internal/config"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/handlers"
	auth_api "media-equipment-tracker/internal/handlers/auth-api"
	"media-equipment-tracker/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	engine := gin.New()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.PostgresHost, config.PostgresUser, config.PostgresPassword, config.PostgresDatabase, config.PostgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Repository
	userRepo := postgres_repo.NewUserRepository(db)

	// Auth
	tokenMaker, err := tokenmaker.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		panic(err.Error())
	}
	hasher, err := hasher.NewHasher()
	if err != nil {
		panic(err.Error())
	}
	authUserServ := auth_user.NewAuthUser(tokenMaker, hasher, config.AccessTokenDuration, userRepo)

	// Groups
	healthRouter := handlers.NewHealthRouter(engine.Group("/"))
	_ = healthRouter

	apiGroup := engine.Group(config.Api_version)
	usersGroup := apiGroup.Group("/")
	usersGroup.Use(middleware.AuthMiddleware(authUserServ, []domain.RoleAuth{}))

	adminsGroup := apiGroup.Group("/")
	adminsGroup.Use(middleware.AuthMiddleware(authUserServ, []domain.RoleAuth{domain.AdminRole}))

	// Routers
	authUserRouter := auth_api.NewAuthUserRouter(apiGroup, authUserServ)
	_ = authUserRouter

	if err := engine.Run(fmt.Sprintf(":%d", config.AppPort)); err != nil {
		panic(err.Error())
	}
}

package main

import (
	"fmt"
	"log"
	"media-equipment-tracker/internal/adapters/postgres_repo"
	"media-equipment-tracker/internal/application/authservice/auth_user"
	"media-equipment-tracker/internal/application/authservice/hasher"
	tokenmaker "media-equipment-tracker/internal/application/authservice/token_maker"
	auth_api "media-equipment-tracker/internal/handlers/auth-api"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	postgresHost        = "localhost"
	postgresPort        = 5432
	postgresUser        = "uuser"
	postgresPassword    = "ppassword"
	postgresDatabase    = "eqtracker"
	tokenSymmetricKey   = "12345678901234567890123456789012"
	accessTokenDuration = 24 * time.Hour
	api_version         = "/api/v1"
	appPort             = 8080
)

func main() {
	engine := gin.New()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", postgresHost, postgresUser, postgresPassword, postgresDatabase, postgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Repository
	userRepo := postgres_repo.NewUserRepository(db)

	// Auth
	tokenMaker, err := tokenmaker.NewTokenMaker(tokenSymmetricKey)
	if err != nil {
		panic(err.Error())
	}
	hasher, err := hasher.NewHasher()
	if err != nil {
		panic(err.Error())
	}
	authUserServ := auth_user.NewAuthUser(tokenMaker, hasher, accessTokenDuration, userRepo)

	// Groups

	apiGroup := engine.Group(api_version)

	// Routers
	authUserRouter := auth_api.NewAuthUserRouter(apiGroup, authUserServ)
	_ = authUserRouter

	if err := engine.Run(fmt.Sprintf(":%d", appPort)); err != nil {
		panic(err.Error())
	}
}

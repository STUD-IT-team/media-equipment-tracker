package main

import (
	"fmt"
	"log"
	"media-equipment-tracker/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getUsers(db *gorm.DB) {
	var users []domain.User

	// Вариант 1: Получить всех пользователей
	result := db.Find(&users)
	if result.Error != nil {
		log.Printf("Error getting users: %v", result.Error)
		return
	}

	fmt.Printf("\n=== Все пользователи (найдено: %d) ===\n", result.RowsAffected)
	for _, user := range users {
		fmt.Printf("ID: %s, Name: %s, Email: %s, IsAdmin: %v\n",
			user.ID, user.FullName, user.Email, user.IsAdmin)
	}
}

func main() {
	dsn := "host=localhost user=uuser password=ppassword dbname=eqtracker port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Включаем логирование SQL запросов
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database!")

	getUsers(db)
	//getUserByEmail(db, "john@example.com")
	//getEquipmentWithDepartments(db)
}

package database

import (
	"fmt"
	"log"
	"test-constructor/config"
	"test-constructor/internal/models"
	"test-constructor/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.Load()
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBName,
		cfg.DBPassword,
	)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных", err)
	}

	DB = connection

	err = DB.AutoMigrate(
		&models.User{},
		&models.Test{},
		&models.Question{},
		&models.Answer{},
		&models.Role{},
		&models.Attempt{},
		&models.EventConfig{},
		&models.ExtraThreshold{},
		&models.UserEvent{},
	)
	if err != nil {
		log.Fatal("Ошибка миграции базы данных", err)
	}

	if err := migrations.SeedRoles(DB); err != nil {
		log.Println("Не удалось заполнить роли:", err)
	}

	if err := migrations.SeedAdmin(DB); err != nil {
		log.Println("Не удалось создать админа", err)
	}

	log.Println("База данных подключена")
}

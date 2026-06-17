package database

import (
	"fmt"
	"log"
	"test-constructor/config"
	"test-constructor/internal/domain"
	"test-constructor/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	cfg := config.Load()
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBName,
		cfg.DBPassword,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных", err)
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Test{},
		&domain.Question{},
		&domain.Answer{},
		&domain.Role{},
		&domain.Attempt{},
		&domain.EventConfig{},
		&domain.ExtraThreshold{},
		&domain.UserEvent{},
	)
	if err != nil {
		log.Fatal("Ошибка миграции базы данных", err)
	}

	if err := migrations.SeedRoles(db); err != nil {
		log.Println("Не удалось заполнить роли:", err)
	}

	if err := migrations.SeedAdmin(db); err != nil {
		log.Println("Не удалось создать админа", err)
	}

	log.Println("База данных подключена")
	return db
}

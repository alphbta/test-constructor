package database

import (
	"fmt"
	"log"
	"test-constructor/config"
	"test-constructor/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.Load()
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=testconstructor port=5432 sslmode=disable", cfg.DBPassword)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных", err)
	}

	DB = connection

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Ошибка миграции базы данных", err)
	}

	log.Println("База данных подключена")
}

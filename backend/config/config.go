package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPassword    string
	AdminEmail    string
	AdminPassword string
	JWTSecret     string
	JWTTTL        int64
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@example.com"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "17a3229b-e5c6-4ab0-ba86-3d87cb7f23fe"
	}

	var jwtTTL int64 = 24

	return &Config{
		dbPassword,
		adminEmail,
		adminPassword,
		jwtSecret,
		jwtTTL,
	}
}

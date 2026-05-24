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
	CRMService    string
	CRMToken      string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
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
	crmService := os.Getenv("CRM_SERVICE")
	if crmService == "" {
		crmService = "http://127.0.0.1:8000"
	}
	crmToken := os.Getenv("CRM_TOKEN")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "testconstructor"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	var jwtTTL int64 = 24

	return &Config{
		dbPassword,
		adminEmail,
		adminPassword,
		jwtSecret,
		jwtTTL,
		crmService,
		crmToken,
		dbHost,
		dbPort,
		dbName,
		dbUser,
	}
}

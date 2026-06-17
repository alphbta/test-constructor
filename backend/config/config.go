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

	return &Config{
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin"),
		JWTSecret:     getEnv("JWT_SECRET", "17a3229b-e5c6-4ab0-ba86-3d87cb7f23fe"),
		CRMService:    getEnv("CRM_SERVICE", "http://127.0.0.1:8000"),
		CRMToken:      getEnv("CRM_TOKEN", ""),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "testconstructor"),
		DBUser:        getEnv("DB_USER", "postgres"),
		JWTTTL:        24,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

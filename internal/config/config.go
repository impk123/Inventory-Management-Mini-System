package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	DatabaseURL        string
	RedisURL           string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	FrontendURL        string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/inventory_db"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3100"),
	}, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

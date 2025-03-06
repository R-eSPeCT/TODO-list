package config

import (
	"os"
)

// Config содержит конфигурацию приложения.
type Config struct {
	DatabaseURL string
}

// LoadConfig загружает конфигурацию из переменных среды
func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "postgres://tododb:tododb@localhost:5432/tododb?sslmode=disable"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

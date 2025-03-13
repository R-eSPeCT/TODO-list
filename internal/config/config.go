package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config содержит конфигурацию приложения.
type Config struct {
	Port            string
	DatabaseURL     string
	RedisURL        string
	JWTSecret       string
	JWTExpiration   time.Duration
	AllowedOrigins  []string
	RateLimitMax    int
	RateLimitWindow time.Duration
}

// LoadConfig загружает конфигурацию из переменных среды
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:            getEnvOrDefault("PORT", "3000"),
		DatabaseURL:     getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todo?sslmode=disable"),
		RedisURL:        getEnvOrDefault("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:       getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		JWTExpiration:   getDurationEnvOrDefault("JWT_EXPIRATION", 24*time.Hour),
		AllowedOrigins:  getStringSliceEnvOrDefault("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		RateLimitMax:    getIntEnvOrDefault("RATE_LIMIT_MAX", 100),
		RateLimitWindow: getDurationEnvOrDefault("RATE_LIMIT_WINDOW", time.Hour),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Port == "" {
		return errors.New("port is required")
	}

	if c.DatabaseURL == "" {
		return errors.New("database URL is required")
	}

	if c.RedisURL == "" {
		return errors.New("redis URL is required")
	}

	if c.JWTSecret == "" {
		return errors.New("JWT secret is required")
	}

	if c.JWTExpiration <= 0 {
		return errors.New("JWT expiration must be positive")
	}

	if len(c.AllowedOrigins) == 0 {
		return errors.New("at least one allowed origin is required")
	}

	if c.RateLimitMax <= 0 {
		return errors.New("rate limit max must be positive")
	}

	if c.RateLimitWindow <= 0 {
		return errors.New("rate limit window must be positive")
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getStringSliceEnvOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

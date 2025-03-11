package config

import (
	"fmt"
	"os"
	"strconv"
)

// RedisConfig содержит настройки подключения к Redis
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisConfig создает новую конфигурацию Redis из переменных окружения
func NewRedisConfig() *RedisConfig {
	port, _ := strconv.Atoi(getEnvOrDefault("REDIS_PORT", "6379"))
	db, _ := strconv.Atoi(getEnvOrDefault("REDIS_DB", "0"))

	return &RedisConfig{
		Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
		Port:     port,
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       db,
	}
}

// GetAddr возвращает адрес подключения к Redis
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

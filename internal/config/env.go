package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnvOrDefault получает значение переменной окружения или возвращает значение по умолчанию
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDurationEnvOrDefault получает значение переменной окружения как time.Duration или возвращает значение по умолчанию
func GetDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetIntEnvOrDefault получает значение переменной окружения как int или возвращает значение по умолчанию
func GetIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetStringSliceEnvOrDefault получает значение переменной окружения как []string или возвращает значение по умолчанию
func GetStringSliceEnvOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

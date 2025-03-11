package config

import (
	"strconv"
	"time"
)

// GRPCConfig содержит настройки gRPC сервера
type GRPCConfig struct {
	Port             int
	JWTSecretKey     string
	TokenDuration    time.Duration
	MaxRequestSize   int
	KeepAlive        time.Duration
	KeepAliveTimeout time.Duration
}

// NewGRPCConfig создает новую конфигурацию gRPC из переменных окружения
func NewGRPCConfig() *GRPCConfig {
	port, _ := strconv.Atoi(getEnvOrDefault("GRPC_PORT", "50051"))
	maxRequestSize, _ := strconv.Atoi(getEnvOrDefault("GRPC_MAX_REQUEST_SIZE", "4194304")) // 4MB
	keepAlive, _ := time.ParseDuration(getEnvOrDefault("GRPC_KEEP_ALIVE", "60s"))
	keepAliveTimeout, _ := time.ParseDuration(getEnvOrDefault("GRPC_KEEP_ALIVE_TIMEOUT", "20s"))
	tokenDuration, _ := time.ParseDuration(getEnvOrDefault("JWT_TOKEN_DURATION", "15m"))

	return &GRPCConfig{
		Port:             port,
		JWTSecretKey:     getEnvOrDefault("JWT_SECRET_KEY", "your-secret-key"),
		TokenDuration:    tokenDuration,
		MaxRequestSize:   maxRequestSize,
		KeepAlive:        keepAlive,
		KeepAliveTimeout: keepAliveTimeout,
	}
}

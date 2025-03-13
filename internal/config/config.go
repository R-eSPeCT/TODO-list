package config

import (
	"errors"
	"time"

	"github.com/yourusername/todo-list/pkg/env"
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
	GRPC            *GRPCConfig
	HTTP            *HTTPConfig
}

// GRPCConfig содержит настройки gRPC сервера
type GRPCConfig struct {
	Port             int
	JWTSecretKey     string
	TokenDuration    time.Duration
	MaxRequestSize   int
	KeepAlive        time.Duration
	KeepAliveTimeout time.Duration
}

// HTTPConfig содержит настройки HTTP сервера
type HTTPConfig struct {
	Port            string
	AllowedOrigins  []string
	RateLimitMax    int
	RateLimitWindow time.Duration
}

// LoadConfig загружает конфигурацию из переменных среды
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:            env.GetEnvOrDefault("PORT", "3000"),
		DatabaseURL:     env.GetEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todo?sslmode=disable"),
		RedisURL:        env.GetEnvOrDefault("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:       env.GetEnvOrDefault("JWT_SECRET", "your-secret-key"),
		JWTExpiration:   env.GetDurationEnvOrDefault("JWT_EXPIRATION", 24*time.Hour),
		AllowedOrigins:  env.GetStringSliceEnvOrDefault("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		RateLimitMax:    env.GetIntEnvOrDefault("RATE_LIMIT_MAX", 100),
		RateLimitWindow: env.GetDurationEnvOrDefault("RATE_LIMIT_WINDOW", time.Hour),
	}

	// Загрузка gRPC конфигурации
	grpcConfig, err := NewGRPCConfig()
	if err != nil {
		return nil, err
	}
	cfg.GRPC = grpcConfig

	// Загрузка HTTP конфигурации
	httpConfig, err := NewHTTPConfig()
	if err != nil {
		return nil, err
	}
	cfg.HTTP = httpConfig

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// NewGRPCConfig создает новую конфигурацию gRPC из переменных окружения
func NewGRPCConfig() (*GRPCConfig, error) {
	port := env.GetIntEnvOrDefault("GRPC_PORT", 50051)
	maxRequestSize := env.GetIntEnvOrDefault("GRPC_MAX_REQUEST_SIZE", 4194304) // 4MB
	keepAlive := env.GetDurationEnvOrDefault("GRPC_KEEP_ALIVE", 60*time.Second)
	keepAliveTimeout := env.GetDurationEnvOrDefault("GRPC_KEEP_ALIVE_TIMEOUT", 20*time.Second)
	tokenDuration := env.GetDurationEnvOrDefault("JWT_TOKEN_DURATION", 15*time.Minute)

	g := &GRPCConfig{
		Port:             port,
		JWTSecretKey:     env.GetEnvOrDefault("JWT_SECRET_KEY", "your-secret-key"),
		TokenDuration:    tokenDuration,
		MaxRequestSize:   maxRequestSize,
		KeepAlive:        keepAlive,
		KeepAliveTimeout: keepAliveTimeout,
	}

	if err := g.validate(); err != nil {
		return nil, err
	}
	return g, nil
}

// NewHTTPConfig создает новую конфигурацию HTTP из переменных окружения
func NewHTTPConfig() (*HTTPConfig, error) {
	h := &HTTPConfig{
		Port:            env.GetEnvOrDefault("PORT", "3000"),
		AllowedOrigins:  env.GetStringSliceEnvOrDefault("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		RateLimitMax:    env.GetIntEnvOrDefault("RATE_LIMIT_MAX", 100),
		RateLimitWindow: env.GetDurationEnvOrDefault("RATE_LIMIT_WINDOW", time.Hour),
	}

	if err := h.validate(); err != nil {
		return nil, err
	}
	return h, nil
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

func (g *GRPCConfig) validate() error {
	if g.Port <= 0 {
		return errors.New("gRPC port must be positive")
	}
	if g.MaxRequestSize <= 0 {
		return errors.New("gRPC max request size must be positive")
	}
	if g.KeepAlive <= 0 {
		return errors.New("gRPC keep alive must be positive")
	}
	if g.KeepAliveTimeout <= 0 {
		return errors.New("gRPC keep alive timeout must be positive")
	}
	if g.TokenDuration <= 0 {
		return errors.New("gRPC token duration must be positive")
	}
	if g.JWTSecretKey == "" {
		return errors.New("gRPC JWT secret key is required")
	}
	return nil
}

func (h *HTTPConfig) validate() error {
	if h.Port == "" {
		return errors.New("HTTP port is required")
	}
	if len(h.AllowedOrigins) == 0 {
		return errors.New("at least one allowed origin is required")
	}
	if h.RateLimitMax <= 0 {
		return errors.New("HTTP rate limit max must be positive")
	}
	if h.RateLimitWindow <= 0 {
		return errors.New("HTTP rate limit window must be positive")
	}
	return nil
}

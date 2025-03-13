package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	// Сохраняем текущие значения переменных окружения
	originalEnv := make(map[string]string)
	for _, env := range []string{
		"PORT",
		"DATABASE_URL",
		"REDIS_URL",
		"JWT_SECRET_KEY",
		"JWT_TOKEN_DURATION",
		"GRPC_PORT",
		"GRPC_MAX_REQUEST_SIZE",
		"GRPC_KEEP_ALIVE",
		"GRPC_KEEP_ALIVE_TIMEOUT",
	} {
		if value := os.Getenv(env); value != "" {
			originalEnv[env] = value
		}
	}

	// Восстанавливаем значения после теста
	defer func() {
		for env, value := range originalEnv {
			if value != "" {
				os.Setenv(env, value)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "valid config",
			env: map[string]string{
				"PORT":                    "3000",
				"DATABASE_URL":            "postgres://postgres:Salaamnder0101@localhost:5432/todo_list?sslmode=disable",
				"REDIS_URL":               "redis://localhost:6379/0",
				"JWT_SECRET_KEY":          "test-secret-key",
				"JWT_TOKEN_DURATION":      "24h",
				"GRPC_PORT":               "50051",
				"GRPC_MAX_REQUEST_SIZE":   "4194304",
				"GRPC_KEEP_ALIVE":         "60s",
				"GRPC_KEEP_ALIVE_TIMEOUT": "20s",
			},
			wantErr: false,
		},
		{
			name: "missing required env",
			env: map[string]string{
				"PORT": "3000",
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			env: map[string]string{
				"PORT":                    "invalid",
				"DATABASE_URL":            "postgres://postgres:Salamander0101@localhost:5432/tododb?sslmode=disable",
				"REDIS_URL":               "redis://localhost:6379/0",
				"JWT_SECRET_KEY":          "test-secret-key",
				"JWT_TOKEN_DURATION":      "24h",
				"GRPC_PORT":               "50051",
				"GRPC_MAX_REQUEST_SIZE":   "4194304",
				"GRPC_KEEP_ALIVE":         "60s",
				"GRPC_KEEP_ALIVE_TIMEOUT": "20s",
			},
			wantErr: true,
		},
		{
			name: "invalid duration",
			env: map[string]string{
				"PORT":                    "3000",
				"DATABASE_URL":            "postgres://postgres:Salamander0101@localhost:5432/tododb?sslmode=disable",
				"REDIS_URL":               "redis://localhost:6379/0",
				"JWT_SECRET_KEY":          "test-secret-key",
				"JWT_TOKEN_DURATION":      "invalid",
				"GRPC_PORT":               "50051",
				"GRPC_MAX_REQUEST_SIZE":   "4194304",
				"GRPC_KEEP_ALIVE":         "60s",
				"GRPC_KEEP_ALIVE_TIMEOUT": "20s",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменные окружения для теста
			for env, value := range tt.env {
				os.Setenv(env, value)
			}

			config, err := LoadConfig()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, config)

			// Проверяем значения конфигурации
			assert.Equal(t, tt.env["PORT"], config.Port)
			assert.Equal(t, tt.env["DATABASE_URL"], config.DatabaseURL)
			assert.Equal(t, tt.env["REDIS_URL"], config.RedisURL)
			assert.Equal(t, tt.env["JWT_SECRET_KEY"], config.JWTSecretKey)
			assert.Equal(t, tt.env["JWT_TOKEN_DURATION"], config.JWTTokenDuration)
			assert.Equal(t, tt.env["GRPC_PORT"], config.GRPCConfig.Port)
			assert.Equal(t, tt.env["GRPC_MAX_REQUEST_SIZE"], config.GRPCConfig.MaxRequestSize)
			assert.Equal(t, tt.env["GRPC_KEEP_ALIVE"], config.GRPCConfig.KeepAlive)
			assert.Equal(t, tt.env["GRPC_KEEP_ALIVE_TIMEOUT"], config.GRPCConfig.KeepAliveTimeout)
		})
	}
}

func TestNewGRPCConfig(t *testing.T) {
	// Сохраняем текущие значения переменных окружения
	originalEnv := make(map[string]string)
	for _, env := range []string{
		"GRPC_PORT",
		"GRPC_MAX_REQUEST_SIZE",
		"GRPC_KEEP_ALIVE",
		"GRPC_KEEP_ALIVE_TIMEOUT",
	} {
		if value := os.Getenv(env); value != "" {
			originalEnv[env] = value
		}
	}

	// Восстанавливаем значения после теста
	defer func() {
		for env, value := range originalEnv {
			if value != "" {
				os.Setenv(env, value)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "valid config",
			env: map[string]string{
				"GRPC_PORT":               "50051",
				"GRPC_MAX_REQUEST_SIZE":   "4194304",
				"GRPC_KEEP_ALIVE":         "60s",
				"GRPC_KEEP_ALIVE_TIMEOUT": "20s",
			},
			wantErr: false,
		},
		{
			name: "missing required env",
			env: map[string]string{
				"GRPC_PORT": "50051",
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			env: map[string]string{
				"GRPC_PORT":               "invalid",
				"GRPC_MAX_REQUEST_SIZE":   "4194304",
				"GRPC_KEEP_ALIVE":         "60s",
				"GRPC_KEEP_ALIVE_TIMEOUT": "20s",
			},
			wantErr: true,
		},
		{
			name: "invalid duration",
			env: map[string]string{
				"GRPC_PORT":               "50051",
				"GRPC_MAX_REQUEST_SIZE":   "4194304",
				"GRPC_KEEP_ALIVE":         "invalid",
				"GRPC_KEEP_ALIVE_TIMEOUT": "20s",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменные окружения для теста
			for env, value := range tt.env {
				os.Setenv(env, value)
			}

			config, err := NewGRPCConfig()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, config)

			// Проверяем значения конфигурации
			assert.Equal(t, tt.env["GRPC_PORT"], config.Port)
			assert.Equal(t, tt.env["GRPC_MAX_REQUEST_SIZE"], config.MaxRequestSize)
			assert.Equal(t, tt.env["GRPC_KEEP_ALIVE"], config.KeepAlive)
			assert.Equal(t, tt.env["GRPC_KEEP_ALIVE_TIMEOUT"], config.KeepAliveTimeout)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Сохраняем текущие значения переменных окружения
	originalEnv := make(map[string]string)
	envVars := []string{
		"PORT",
		"DATABASE_URL",
		"REDIS_URL",
		"JWT_SECRET",
		"JWT_EXPIRATION",
		"ALLOWED_ORIGINS",
		"RATE_LIMIT_MAX",
		"RATE_LIMIT_WINDOW",
		"GRPC_PORT",
		"GRPC_MAX_REQUEST_SIZE",
		"GRPC_KEEP_ALIVE",
		"GRPC_KEEP_ALIVE_TIMEOUT",
		"JWT_SECRET_KEY",
		"JWT_TOKEN_DURATION",
	}

	for _, env := range envVars {
		if value := os.Getenv(env); value != "" {
			originalEnv[env] = value
		}
	}

	// Восстанавливаем переменные окружения после теста
	defer func() {
		for _, env := range envVars {
			if originalValue, exists := originalEnv[env]; exists {
				os.Setenv(env, originalValue)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		wantErr  bool
		validate func(*testing.T, *Config)
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"PORT":                    "8080",
				"DATABASE_URL":            "postgres://user:pass@localhost:5432/db",
				"REDIS_URL":               "redis://localhost:6379/0",
				"JWT_SECRET":              "test-secret",
				"JWT_EXPIRATION":          "24h",
				"ALLOWED_ORIGINS":         "http://localhost:3000,http://localhost:8080",
				"RATE_LIMIT_MAX":          "100",
				"RATE_LIMIT_WINDOW":       "1h",
				"GRPC_PORT":               "50051",
				"GRPC_MAX_REQUEST_SIZE":   "4194304",
				"GRPC_KEEP_ALIVE":         "60s",
				"GRPC_KEEP_ALIVE_TIMEOUT": "20s",
				"JWT_SECRET_KEY":          "grpc-secret-key",
				"JWT_TOKEN_DURATION":      "15m",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				// Проверка основных полей
				assert.Equal(t, "8080", cfg.Port)
				assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
				assert.Equal(t, "redis://localhost:6379/0", cfg.RedisURL)
				assert.Equal(t, "test-secret", cfg.JWTSecret)
				assert.Equal(t, 24*time.Hour, cfg.JWTExpiration)
				assert.Equal(t, []string{"http://localhost:3000", "http://localhost:8080"}, cfg.AllowedOrigins)
				assert.Equal(t, 100, cfg.RateLimitMax)
				assert.Equal(t, time.Hour, cfg.RateLimitWindow)

				// Проверка HTTP конфигурации
				assert.Equal(t, "8080", cfg.HTTP.Port)
				assert.Equal(t, []string{"http://localhost:3000", "http://localhost:8080"}, cfg.HTTP.AllowedOrigins)
				assert.Equal(t, 100, cfg.HTTP.RateLimitMax)
				assert.Equal(t, time.Hour, cfg.HTTP.RateLimitWindow)

				// Проверка gRPC конфигурации
				assert.Equal(t, 50051, cfg.GRPC.Port)
				assert.Equal(t, "grpc-secret-key", cfg.GRPC.JWTSecretKey)
				assert.Equal(t, 15*time.Minute, cfg.GRPC.TokenDuration)
				assert.Equal(t, 4194304, cfg.GRPC.MaxRequestSize)
				assert.Equal(t, 60*time.Second, cfg.GRPC.KeepAlive)
				assert.Equal(t, 20*time.Second, cfg.GRPC.KeepAliveTimeout)
			},
		},
		{
			name:    "default values",
			envVars: map[string]string{},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				// Проверка основных полей по умолчанию
				assert.Equal(t, "3000", cfg.Port)
				assert.Equal(t, "postgres://postgres:postgres@localhost:5432/todo?sslmode=disable", cfg.DatabaseURL)
				assert.Equal(t, "redis://localhost:6379/0", cfg.RedisURL)
				assert.Equal(t, "your-secret-key", cfg.JWTSecret)
				assert.Equal(t, 24*time.Hour, cfg.JWTExpiration)
				assert.Equal(t, []string{"http://localhost:3000"}, cfg.AllowedOrigins)
				assert.Equal(t, 100, cfg.RateLimitMax)
				assert.Equal(t, time.Hour, cfg.RateLimitWindow)

				// Проверка HTTP конфигурации по умолчанию
				assert.Equal(t, "3000", cfg.HTTP.Port)
				assert.Equal(t, []string{"http://localhost:3000"}, cfg.HTTP.AllowedOrigins)
				assert.Equal(t, 100, cfg.HTTP.RateLimitMax)
				assert.Equal(t, time.Hour, cfg.HTTP.RateLimitWindow)

				// Проверка gRPC конфигурации по умолчанию
				assert.Equal(t, 50051, cfg.GRPC.Port)
				assert.Equal(t, "your-secret-key", cfg.GRPC.JWTSecretKey)
				assert.Equal(t, 15*time.Minute, cfg.GRPC.TokenDuration)
				assert.Equal(t, 4194304, cfg.GRPC.MaxRequestSize)
				assert.Equal(t, 60*time.Second, cfg.GRPC.KeepAlive)
				assert.Equal(t, 20*time.Second, cfg.GRPC.KeepAliveTimeout)
			},
		},
		{
			name: "invalid JWT expiration",
			envVars: map[string]string{
				"JWT_EXPIRATION": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid rate limit max",
			envVars: map[string]string{
				"RATE_LIMIT_MAX": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid rate limit window",
			envVars: map[string]string{
				"RATE_LIMIT_WINDOW": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid gRPC port",
			envVars: map[string]string{
				"GRPC_PORT": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid gRPC keep alive",
			envVars: map[string]string{
				"GRPC_KEEP_ALIVE": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем переменные окружения перед тестом
			for _, env := range envVars {
				os.Unsetenv(env)
			}

			// Устанавливаем тестовые переменные окружения
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := LoadConfig()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestConfig_validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Port:            "8080",
				DatabaseURL:     "postgres://user:pass@localhost:5432/db",
				RedisURL:        "redis://localhost:6379/0",
				JWTSecret:       "test-secret",
				JWTExpiration:   24 * time.Hour,
				AllowedOrigins:  []string{"http://localhost:3000"},
				RateLimitMax:    100,
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "empty port",
			config: &Config{
				DatabaseURL:     "postgres://user:pass@localhost:5432/db",
				RedisURL:        "redis://localhost:6379/0",
				JWTSecret:       "test-secret",
				JWTExpiration:   24 * time.Hour,
				AllowedOrigins:  []string{"http://localhost:3000"},
				RateLimitMax:    100,
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "empty database URL",
			config: &Config{
				Port:            "8080",
				RedisURL:        "redis://localhost:6379/0",
				JWTSecret:       "test-secret",
				JWTExpiration:   24 * time.Hour,
				AllowedOrigins:  []string{"http://localhost:3000"},
				RateLimitMax:    100,
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "empty Redis URL",
			config: &Config{
				Port:            "8080",
				DatabaseURL:     "postgres://user:pass@localhost:5432/db",
				JWTSecret:       "test-secret",
				JWTExpiration:   24 * time.Hour,
				AllowedOrigins:  []string{"http://localhost:3000"},
				RateLimitMax:    100,
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "empty JWT secret",
			config: &Config{
				Port:            "8080",
				DatabaseURL:     "postgres://user:pass@localhost:5432/db",
				RedisURL:        "redis://localhost:6379/0",
				JWTExpiration:   24 * time.Hour,
				AllowedOrigins:  []string{"http://localhost:3000"},
				RateLimitMax:    100,
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "zero JWT expiration",
			config: &Config{
				Port:            "8080",
				DatabaseURL:     "postgres://user:pass@localhost:5432/db",
				RedisURL:        "redis://localhost:6379/0",
				JWTSecret:       "test-secret",
				AllowedOrigins:  []string{"http://localhost:3000"},
				RateLimitMax:    100,
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "empty allowed origins",
			config: &Config{
				Port:            "8080",
				DatabaseURL:     "postgres://user:pass@localhost:5432/db",
				RedisURL:        "redis://localhost:6379/0",
				JWTSecret:       "test-secret",
				JWTExpiration:   24 * time.Hour,
				RateLimitMax:    100,
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "zero rate limit max",
			config: &Config{
				Port:            "8080",
				DatabaseURL:     "postgres://user:pass@localhost:5432/db",
				RedisURL:        "redis://localhost:6379/0",
				JWTSecret:       "test-secret",
				JWTExpiration:   24 * time.Hour,
				AllowedOrigins:  []string{"http://localhost:3000"},
				RateLimitWindow: time.Hour,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "zero rate limit window",
			config: &Config{
				Port:           "8080",
				DatabaseURL:    "postgres://user:pass@localhost:5432/db",
				RedisURL:       "redis://localhost:6379/0",
				JWTSecret:      "test-secret",
				JWTExpiration:  24 * time.Hour,
				AllowedOrigins: []string{"http://localhost:3000"},
				RateLimitMax:   100,
				HTTP: &HTTPConfig{
					Port:            "8080",
					AllowedOrigins:  []string{"http://localhost:3000"},
					RateLimitMax:    100,
					RateLimitWindow: time.Hour,
				},
				GRPC: &GRPCConfig{
					Port:             50051,
					JWTSecretKey:     "test-secret",
					TokenDuration:    15 * time.Minute,
					MaxRequestSize:   4194304,
					KeepAlive:        60 * time.Second,
					KeepAliveTimeout: 20 * time.Second,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

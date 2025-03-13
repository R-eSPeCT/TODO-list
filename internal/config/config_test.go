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
				"DATABASE_URL":            "postgres://postgres:postgres@localhost:5432/todo_list?sslmode=disable",
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
				"DATABASE_URL":            "postgres://postgres:postgres@localhost:5432/todo_list?sslmode=disable",
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
				"DATABASE_URL":            "postgres://postgres:postgres@localhost:5432/todo_list?sslmode=disable",
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

			config, err := NewConfig()
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
package config

// Config содержит конфигурацию приложения.
type Config struct {
	DatabaseURL string
}

// LoadConfig загружает конфигурацию из переменных среды
func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "postgres://postgres:Salamander0101@localhost:5432/tododb?sslmode=disable"),
	}
}

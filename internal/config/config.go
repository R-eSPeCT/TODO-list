package config

import "os"

// Config содержит конфигурацию приложения.
type Config struct {
	DatabaseURL string
	ServerPort  string
}

// LoadConfig загружает конфигурацию из переменных среды
func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:Salamander0101@localhost:5432/tododb"
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "3000"
	}

	return &Config{
		DatabaseURL: dbURL,
		ServerPort:  serverPort,
	}
}

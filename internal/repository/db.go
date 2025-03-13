package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Connect создает подключение к базе данных PostgreSQL
func Connect(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Выполняем миграции
	if err := runMigrations(pool); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return pool, nil
}

// runMigrations выполняет SQL-миграции
func runMigrations(pool *pgxpool.Pool) error {
	// Создаем таблицу для отслеживания миграций, если её нет
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Получаем список уже примененных миграций
	rows, err := pool.Query(context.Background(), "SELECT name FROM migrations ORDER BY id")
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration name: %w", err)
		}
		appliedMigrations[name] = true
	}

	// Получаем список файлов миграций
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Сортируем файлы по имени
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Выполняем каждую миграцию
	for _, file := range migrationFiles {
		if appliedMigrations[file] {
			continue
		}

		// Читаем содержимое файла миграции
		content, err := ioutil.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Начинаем транзакцию
		tx, err := pool.Begin(context.Background())
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Выполняем SQL-скрипт
		if _, err := tx.Exec(context.Background(), string(content)); err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		// Отмечаем миграцию как выполненную
		if _, err := tx.Exec(context.Background(), "INSERT INTO migrations (name) VALUES ($1)", file); err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("failed to mark migration %s as applied: %w", file, err)
		}

		// Подтверждаем транзакцию
		if err := tx.Commit(context.Background()); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		log.Printf("Applied migration: %s", file)
	}

	return nil
}

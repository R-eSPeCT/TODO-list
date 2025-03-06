package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// createSchema гарантирует существование необходимых таблиц
func createSchema(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now())`,
	)
	if err != nil {
		log.Printf("Failed to create todos table: %v", err)
		return err
	}
	return nil
}

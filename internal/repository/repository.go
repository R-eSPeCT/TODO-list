package repository

import (
	"context"
	"github.com/google/uuid"
	"github
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/google/uuid"
	"vacancy/internal/models"
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

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

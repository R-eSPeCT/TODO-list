package repository

import (
	"TODO-lis
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Repositories struct {
	User UserRepository
	Todo TodoRepository
}

func NewRepositories(db *pgxpool.Pool) (*Repositories, error) {
	if err := createSchema(db); err != nil {
		return nil, err
	}

	return &Repositories{
		User: NewUserRepository(db),
		Todo: NewTodoRepository(db),
	}, nil
}

// createSchema гарантирует существование необходимых таблиц
func createSchema(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now()
		)`)
	if err != nil {
		log.Printf("Failed to create users table: %v", err)
		return err
	}

	_, err = db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
			priority TEXT CHECK (priority IN ('low', 'medium', 'high')) DEFAULT 'medium',
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now()
		)`)
	if err != nil {
		log.Printf("Failed to create tasks table: %v", err)
		return err
	}
	return nil
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetGroupedTodos(ctx context.Context, userID uuid.UUID) ([]models.TodoGroup, error)
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/yourusername/todo-list/internal/models"
	"strings"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	_, err := r.db.Exec(ctx,
		"INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return errors.New("user with this email or username already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	err := r.db.QueryRow(ctx,
		"SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	var user models.User
	err := r.db.QueryRow(ctx,
		"SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if user.ID == uuid.Nil {
		return errors.New("invalid user ID")
	}

	_, err := r.db.Exec(ctx,
		"UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = $4 WHERE id = $5",
		user.Username, user.Email, user.PasswordHash, user.UpdatedAt, user.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return errors.New("user with this email or username already exists")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid user ID")
	}

	result, err := r.db.Exec(ctx,
		"DELETE FROM users WHERE id = $1",
		id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

// isUniqueViolation проверяет, является ли ошибка нарушением уникального ограничения
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// Проверяем код ошибки PostgreSQL для нарушения уникального ограничения
	return strings.Contains(err.Error(), "unique constraint")
}

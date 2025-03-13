package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/google/uuid"
)

type todoRepository struct {
	db *sql.DB
}

func (r *todoRepository) Create(ctx context.Context, todo *models.Todo) error {
	query := `
		INSERT INTO todos (id, title, description, status, priority, due_date, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		todo.ID, todo.Title, todo.Description, todo.Status,
		todo.Priority, todo.DueDate, todo.UserID,
		todo.CreatedAt, todo.UpdatedAt,
	)
	return err
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `
		SELECT id, title, description, status, priority, due_date, user_id, created_at, updated_at
		FROM todos WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Status,
		&todo.Priority, &todo.DueDate, &todo.UserID,
		&todo.CreatedAt, &todo.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("todo not found")
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Todo, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, user_id, created_at, updated_at
		FROM todos WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Status,
			&todo.Priority, &todo.DueDate, &todo.UserID,
			&todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *models.Todo) error {
	query := `
		UPDATE todos
		SET title = $1, description = $2, status = $3, priority = $4,
			due_date = $5, updated_at = $6
		WHERE id = $7 AND user_id = $8
	`
	result, err := r.db.ExecContext(ctx, query,
		todo.Title, todo.Description, todo.Status, todo.Priority,
		todo.DueDate, todo.UpdatedAt, todo.ID, todo.UserID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("todo not found or unauthorized")
	}
	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM todos WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("todo not found")
	}
	return nil
}

func (r *todoRepository) GetGroupedByStatus(ctx context.Context, userID uuid.UUID) (map[string][]models.Todo, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, user_id, created_at, updated_at
		FROM todos WHERE user_id = $1
		ORDER BY status, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouped := make(map[string][]models.Todo)
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Status,
			&todo.Priority, &todo.DueDate, &todo.UserID,
			&todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		grouped[todo.Status] = append(grouped[todo.Status], todo)
	}
	return grouped, nil
}

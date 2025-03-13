package repository

import (
	"TODO-list/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type todoRepo struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) TodoRepository {
	return &todoRepo{db: db}
}

func (r *todoRepo) Create(ctx context.Context, todo *models.Todo) error {
	if todo == nil {
		return errors.New("todo is nil")
	}

	if todo.UserID == uuid.Nil {
		return errors.New("invalid user ID")
	}

	if todo.ID == uuid.Nil {
		todo.ID = uuid.New()
	}
	if todo.CreatedAt.IsZero() {
		now := time.Now()
		todo.CreatedAt = now
		todo.UpdatedAt = now
	}
	if todo.Status == "" {
		todo.Status = "new"
	}
	if todo.Priority == "" {
		todo.Priority = "medium"
	}

	_, err := r.db.Exec(ctx,
		"INSERT INTO tasks (id, user_id, title, description, status, priority, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Priority, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

func (r *todoRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid todo ID")
	}

	var todo models.Todo
	err := r.db.QueryRow(ctx,
		"SELECT id, user_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE id = $1",
		id,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Priority, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	return &todo, nil
}

func (r *todoRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	rows, err := r.db.Query(ctx,
		"SELECT id, user_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE user_id = $1 ORDER BY created_at DESC",
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Priority, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, &todo)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todos: %w", err)
	}
	return todos, nil
}

func (r *todoRepo) Update(ctx context.Context, todo *models.Todo) error {
	if todo == nil {
		return errors.New("todo is nil")
	}

	if todo.ID == uuid.Nil {
		return errors.New("invalid todo ID")
	}

	if todo.UserID == uuid.Nil {
		return errors.New("invalid user ID")
	}

	todo.UpdatedAt = time.Now()
	result, err := r.db.Exec(ctx,
		"UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, updated_at = $5 WHERE id = $6 AND user_id = $7",
		todo.Title, todo.Description, todo.Status, todo.Priority, todo.UpdatedAt, todo.ID, todo.UserID)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func (r *todoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid todo ID")
	}

	result, err := r.db.Exec(ctx,
		"DELETE FROM tasks WHERE id = $1",
		id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func (r *todoRepo) GetGroupedTodos(ctx context.Context, userID uuid.UUID) ([]models.TodoGroup, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	rows, err := r.db.Query(ctx, `
		SELECT status, priority, COUNT(*) as count,
		ARRAY_AGG(ROW_TO_JSON(t)) as tasks
		FROM (
			SELECT id, user_id, title, description, status, priority, created_at, updated_at
			FROM tasks
			WHERE user_id = $1
			ORDER BY created_at DESC
		) t
		GROUP BY status, priority
		ORDER BY priority DESC, status
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query grouped todos: %w", err)
	}
	defer rows.Close()

	var groups []models.TodoGroup
	for rows.Next() {
		var group models.TodoGroup
		var tasksJson []byte
		if err := rows.Scan(&group.Status, &group.Priority, &group.Count, &tasksJson); err != nil {
			return nil, fmt.Errorf("failed to scan todo group: %w", err)
		}
		if err := json.Unmarshal(tasksJson, &group.Tasks); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
		}
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todo groups: %w", err)
	}
	return groups, nil
}

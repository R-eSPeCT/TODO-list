package repository

import (
	"TODO-list/internal/models"
	"context"
	"encoding/json"
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
	return err
}

func (r *todoRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.QueryRow(ctx,
		"SELECT id, user_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE id = $1",
		id,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Priority, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, user_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE user_id = $1 ORDER BY created_at DESC",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Priority, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (r *todoRepo) Update(ctx context.Context, todo *models.Todo) error {
	todo.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		"UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, updated_at = $5 WHERE id = $6 AND user_id = $7",
		todo.Title, todo.Description, todo.Status, todo.Priority, todo.UpdatedAt, todo.ID, todo.UserID)
	return err
}

func (r *todoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM tasks WHERE id = $1",
		id)
	return err
}

func (r *todoRepo) GetGroupedTodos(ctx context.Context, userID uuid.UUID) ([]models.TodoGroup, error) {
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
		return nil, err
	}
	defer rows.Close()

	var groups []models.TodoGroup
	for rows.Next() {
		var group models.TodoGroup
		var tasksJson []byte
		if err := rows.Scan(&group.Status, &group.Priority, &group.Count, &tasksJson); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(tasksJson, &group.Tasks); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

package repository

import (
	"TODO-list/internal/models"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetTodos(userID int) ([]models.Todo, error) {
	rows, err := r.db.Query(context.Background(),
		"SELECT id, title, status, updated_at, created_at, description, priority FROM tasks WHERE user_id = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.UpdatedAt, &todo.CreatedAt, &todo.Description, &todo.Priority); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) CreateTodo(todo models.Todo, userID int) error {
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO tasks (title, status, description, priority, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())",
		todo.Title, todo.Status, todo.Description, todo.Priority, userID)
	return err
}

func (r *TodoRepository) UpdateTodo(id int, todo models.Todo, userID int) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE tasks SET title = $1, status = $2, description = $3, priority = $4 WHERE id = $5 AND user_id = $6",
		todo.Title, todo.Status, todo.Description, todo.Priority, id, userID)
	return err
}

func (r *TodoRepository) DeleteTodo(id int, userID int) error {
	_, err := r.db.Exec(context.Background(),
		"DELETE FROM tasks WHERE id = $1 AND user_id = $2",
		id, userID)
	return err
}

func (r *TodoRepository) GetTodoByID(id int, userID int) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.QueryRow(context.Background(),
		"SELECT id, title, status, updated_at, created_at, description, priority FROM tasks WHERE id = $1 AND user_id = $2",
		id, userID,
	).Scan(&todo.ID, &todo.Title, &todo.Status, &todo.UpdatedAt, &todo.CreatedAt, &todo.Description, &todo.Priority)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) GetGroupedTodos(userID int) ([]models.TodoGroup, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT status, priority, COUNT(*) as count,
		ARRAY_AGG(ROW_TO_JSON(t)) as tasks
		FROM (
			SELECT id, title, description, status, priority, created_at, updated_at
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

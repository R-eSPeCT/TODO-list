package repository

import (
	"TODO-list/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetTodos() ([]models.Todo, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, title, status, updated_at,created_at,description FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.UpdatedAt, &todo.CreatedAt, &todo.Description); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) CreateTodo(todo models.Todo) error {
	_, err := r.db.Exec(context.Background(), "INSERT INTO tasks (title, status, description) VALUES ($1, $2, $3)", todo.Title, todo.Status, todo.Description)
	return err
}

func (r *TodoRepository) UpdateTodo(id int, todo models.Todo) error {
	_, err := r.db.Exec(context.Background(), "UPDATE tasks SET title = $1, status = $2 WHERE id = $3", todo.Title, todo.Status, id)
	return err
}

func (r *TodoRepository) DeleteTodo(id int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
	return err
}

package models

import (
	"time"

	"github.com/google/uuid"
)

// Todo представляет задачу в системе
type Todo struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	Priority    string    `json:"priority" db:"priority"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTodoRequest представляет запрос на создание задачи
type CreateTodoRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
}

// UpdateTodoRequest представляет запрос на обновление задачи
type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// TodoGroup представляет группировку задач
type TodoGroup struct {
	Status   string  `json:"status"`
	Priority string  `json:"priority"`
	Count    int     `json:"count"`
	Tasks    []*Todo `json:"tasks"`
}

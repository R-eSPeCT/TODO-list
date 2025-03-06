package models

import "time"

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    string    `json:"priority" validate:"required,oneof=low medium high"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoGroup представляет группировку задач
type TodoGroup struct {
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Count    int    `json:"count"`
	Tasks    []Todo `json:"tasks"`
}

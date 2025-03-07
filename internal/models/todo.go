package models

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoGroup представляет группировку задач
type TodoGroup struct {
	Status   string  `json:"status"`
	Priority string  `json:"priority"`
	Count    int     `json:"count"`
	Tasks    []*Todo `json:"tasks"`
}

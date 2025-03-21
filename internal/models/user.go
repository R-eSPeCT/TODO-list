package models

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя в системе
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest представляет запрос на регистрацию пользователя
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest представляет запрос на вход пользователя
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse представляет ответ на запрос входа
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

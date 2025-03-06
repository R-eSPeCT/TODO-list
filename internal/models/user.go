package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username" validate:"required,min=3,max=50"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"-"` // Пароль не отправляется в JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserCreate представляет данные для создания пользователя
type UserCreate struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLogin представляет данные для входа
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

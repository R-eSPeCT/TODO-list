package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIResponse представляет стандартный формат ответа API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewSuccessResponse создает успешный ответ API
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse создает ответ API с ошибкой
func NewErrorResponse(err string) *APIResponse {
	return &APIResponse{
		Success: false,
		Error:   err,
	}
}

// WithTimeout создает контекст с таймаутом для операций с базой данных
func WithTimeout(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Context(), 5*time.Second)
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationErrorResponse создает ответ API с ошибками валидации
func NewValidationErrorResponse(errors []ValidationError) *APIResponse {
	return &APIResponse{
		Success: false,
		Data:    errors,
		Error:   "Validation failed",
	}
}

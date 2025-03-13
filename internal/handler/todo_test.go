package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) *fiber.App {
	app := fiber.New()
	return app
}

func TestTodoHandler_Create(t *testing.T) {
	app := setupTestApp(t)
	handler := NewTodoHandler(nil) // Мокаем репозиторий

	app.Post("/todos", handler.Create)

	token := "test-token"

	tests := []struct {
		name       string
		token      string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name:  "valid todo",
			token: token,
			payload: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
				"status":      "pending",
				"due_date":    time.Now().Add(24 * time.Hour),
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:  "missing title",
			token: token,
			payload: map[string]interface{}{
				"description": "Test Description",
				"status":      "pending",
				"due_date":    time.Now().Add(24 * time.Hour),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			token:      "",
			payload:    map[string]interface{}{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestTodoHandler_GetAll(t *testing.T) {
	app := setupTestApp(t)
	handler := NewTodoHandler(nil) // Мокаем репозиторий

	app.Get("/todos", handler.GetAll)

	token := "test-token"

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "valid request",
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestTodoHandler_GetByID(t *testing.T) {
	app := setupTestApp(t)
	handler := NewTodoHandler(nil) // Мокаем репозиторий

	app.Get("/todos/:id", handler.GetByID)

	token := "test-token"
	todoID := uuid.New()

	tests := []struct {
		name       string
		token      string
		todoID     string
		wantStatus int
	}{
		{
			name:       "valid request",
			token:      token,
			todoID:     todoID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid todo id",
			token:      token,
			todoID:     "invalid-id",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			token:      "",
			todoID:     todoID.String(),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/todos/"+tt.todoID, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestTodoHandler_Update(t *testing.T) {
	app := setupTestApp(t)
	handler := NewTodoHandler(nil) // Мокаем репозиторий

	app.Put("/todos/:id", handler.Update)

	token := "test-token"
	todoID := uuid.New()

	tests := []struct {
		name       string
		token      string
		todoID     string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name:   "valid update",
			token:  token,
			todoID: todoID.String(),
			payload: map[string]interface{}{
				"title":       "Updated Todo",
				"description": "Updated Description",
				"status":      "completed",
				"due_date":    time.Now().Add(48 * time.Hour),
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "invalid todo id",
			token:  token,
			todoID: "invalid-id",
			payload: map[string]interface{}{
				"title": "Updated Todo",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			token:      "",
			todoID:     todoID.String(),
			payload:    map[string]interface{}{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/todos/"+tt.todoID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestTodoHandler_Delete(t *testing.T) {
	app := setupTestApp(t)
	handler := NewTodoHandler(nil) // Мокаем репозиторий

	app.Delete("/todos/:id", handler.Delete)

	token := "test-token"
	todoID := uuid.New()

	tests := []struct {
		name       string
		token      string
		todoID     string
		wantStatus int
	}{
		{
			name:       "valid delete",
			token:      token,
			todoID:     todoID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid todo id",
			token:      token,
			todoID:     "invalid-id",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			token:      "",
			todoID:     todoID.String(),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/todos/"+tt.todoID, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
} 
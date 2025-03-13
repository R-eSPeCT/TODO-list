package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) *fiber.App {
	app := fiber.New()
	return app
}

func TestUserHandler_Register(t *testing.T) {
	app := setupTestApp(t)
	handler := NewUserHandler(nil) // Мокаем репозиторий

	app.Post("/register", handler.Register)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid registration",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			payload: map[string]interface{}{
				"email":    "invalid-email",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	app := setupTestApp(t)
	handler := NewUserHandler(nil) // Мокаем репозиторий

	app.Post("/login", handler.Login)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid login",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing fields",
			payload: map[string]interface{}{
				"email": "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	app := setupTestApp(t)
	handler := NewUserHandler(nil) // Мокаем репозиторий

	app.Get("/profile", handler.GetProfile)

	// Создаем тестовый токен
	token := "test-token"

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "valid token",
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	app := setupTestApp(t)
	handler := NewUserHandler(nil) // Мокаем репозиторий

	app.Put("/profile", handler.UpdateProfile)

	token := "test-token"

	tests := []struct {
		name       string
		token      string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name:  "valid update",
			token: token,
			payload: map[string]interface{}{
				"email": "updated@example.com",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "invalid email",
			token: token,
			payload: map[string]interface{}{
				"email": "invalid-email",
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

			req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader(body))
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
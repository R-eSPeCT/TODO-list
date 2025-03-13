package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id uuid.UUID) (*Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockTodoRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Todo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Todo), args.Error(1)
}

func (m *MockTodoRepository) Update(ctx context.Context, todo *Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestApp(t *testing.T) *fiber.App {
	app := fiber.New()
	return app
}

func TestTodoHandler_Create(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockTodoRepository)
	handler := NewTodoHandler(mockRepo)

	app.Post("/todos", handler.Create)

	// Генерируем валидный токен
	userID := uuid.New()
	jwtManager := NewJWTManager("test-secret-key", "1h")
	token, err := jwtManager.Generate(userID)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		payload    map[string]interface{}
		wantStatus int
		setupMock  func()
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
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*Todo")).Return(nil)
			},
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
			if tt.setupMock != nil {
				tt.setupMock()
			}

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

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_GetAll(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockTodoRepository)
	handler := NewTodoHandler(mockRepo)

	app.Get("/todos", handler.GetAll)

	// Генерируем валидный токен
	userID := uuid.New()
	jwtManager := NewJWTManager("test-secret-key", "1h")
	token, err := jwtManager.Generate(userID)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name:       "valid request",
			token:      token,
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetByUserID", mock.Anything, userID).Return([]*Todo{}, nil)
			},
		},
		{
			name:       "missing token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_GetByID(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockTodoRepository)
	handler := NewTodoHandler(mockRepo)

	app.Get("/todos/:id", handler.GetByID)

	// Генерируем валидный токен
	userID := uuid.New()
	jwtManager := NewJWTManager("test-secret-key", "1h")
	token, err := jwtManager.Generate(userID)
	require.NoError(t, err)

	todoID := uuid.New()

	tests := []struct {
		name       string
		token      string
		todoID     string
		wantStatus int
		setupMock  func()
	}{
		{
			name:       "valid request",
			token:      token,
			todoID:     todoID.String(),
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetByID", mock.Anything, todoID).Return(&Todo{}, nil)
			},
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
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/todos/"+tt.todoID, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_Update(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockTodoRepository)
	handler := NewTodoHandler(mockRepo)

	app.Put("/todos/:id", handler.Update)

	// Генерируем валидный токен
	userID := uuid.New()
	jwtManager := NewJWTManager("test-secret-key", "1h")
	token, err := jwtManager.Generate(userID)
	require.NoError(t, err)

	todoID := uuid.New()

	tests := []struct {
		name       string
		token      string
		todoID     string
		payload    map[string]interface{}
		wantStatus int
		setupMock  func()
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
			setupMock: func() {
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*Todo")).Return(nil)
			},
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
			if tt.setupMock != nil {
				tt.setupMock()
			}

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

			mockRepo.AssertExpectations(t)
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

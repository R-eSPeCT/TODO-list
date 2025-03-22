package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/R-eSPeCT/todo-list/internal/auth"
	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoRepository) Update(ctx context.Context, todo *models.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTodoHandler_Create(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	h := NewTodoHandler(mockRepo, jwtManager)

	app := fiber.New()
	app.Post("/todos", h.Create)

	userID := uuid.New()
	testUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	token, err := jwtManager.Generate(testUser)
	require.NoError(t, err)

	tests := []struct {
		name       string
		input      map[string]interface{}
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name: "успешное создание",
			input: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
			},
			token:      token,
			wantStatus: http.StatusCreated,
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil)
			},
		},
		{
			name: "missing title",
			input: map[string]interface{}{
				"description": "Test Description",
			},
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			input:      map[string]interface{}{},
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			jsonBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/todos", bytes.NewReader(jsonBody))
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
	mockRepo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	h := NewTodoHandler(mockRepo, jwtManager)

	app := fiber.New()
	app.Get("/todos", h.GetAll)

	userID := uuid.New()
	testUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	token, err := jwtManager.Generate(testUser)
	require.NoError(t, err)

	todos := []*models.Todo{
		{
			ID:          uuid.New(),
			UserID:      userID,
			Title:       "Test Todo 1",
			Description: "Test Description 1",
		},
		{
			ID:          uuid.New(),
			UserID:      userID,
			Title:       "Test Todo 2",
			Description: "Test Description 2",
		},
	}

	tests := []struct {
		name       string
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name:       "успешное получение",
			token:      token,
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetAll", mock.Anything, userID).Return(todos, nil)
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

			req := httptest.NewRequest("GET", "/todos", nil)
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
	mockRepo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	h := NewTodoHandler(mockRepo, jwtManager)

	app := fiber.New()
	app.Get("/todos/:id", h.GetByID)

	userID := uuid.New()
	testUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	token, err := jwtManager.Generate(testUser)
	require.NoError(t, err)

	todoID := uuid.New()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Test Todo",
		Description: "Test Description",
	}

	tests := []struct {
		name       string
		todoID     string
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name:       "успешное получение",
			todoID:     todoID.String(),
			token:      token,
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetByID", mock.Anything, todoID).Return(todo, nil)
			},
		},
		{
			name:       "invalid todo id",
			todoID:     "invalid-id",
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			todoID:     todoID.String(),
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest("GET", "/todos/"+tt.todoID, nil)
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
	mockRepo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	h := NewTodoHandler(mockRepo, jwtManager)

	app := fiber.New()
	app.Put("/todos/:id", h.Update)

	userID := uuid.New()
	testUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	token, err := jwtManager.Generate(testUser)
	require.NoError(t, err)

	todoID := uuid.New()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Test Todo",
		Description: "Test Description",
	}

	tests := []struct {
		name       string
		todoID     string
		input      map[string]interface{}
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name:   "успешное обновление",
			todoID: todoID.String(),
			input: map[string]interface{}{
				"title":       "Updated Todo",
				"description": "Updated Description",
			},
			token:      token,
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetByID", mock.Anything, todoID).Return(todo, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil)
			},
		},
		{
			name:   "invalid todo id",
			todoID: "invalid-id",
			input: map[string]interface{}{
				"title": "Updated Todo",
			},
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			todoID:     todoID.String(),
			input:      map[string]interface{}{},
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			jsonBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("PUT", "/todos/"+tt.todoID, bytes.NewReader(jsonBody))
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
	mockRepo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	h := NewTodoHandler(mockRepo, jwtManager)

	app := fiber.New()
	app.Delete("/todos/:id", h.Delete)

	userID := uuid.New()
	testUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	token, err := jwtManager.Generate(testUser)
	require.NoError(t, err)

	todoID := uuid.New()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Test Todo",
		Description: "Test Description",
	}

	tests := []struct {
		name       string
		todoID     string
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name:       "успешное удаление",
			todoID:     todoID.String(),
			token:      token,
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetByID", mock.Anything, todoID).Return(todo, nil)
				mockRepo.On("Delete", mock.Anything, todoID).Return(nil)
			},
		},
		{
			name:       "invalid todo id",
			todoID:     "invalid-id",
			token:      token,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing token",
			todoID:     todoID.String(),
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest("DELETE", "/todos/"+tt.todoID, nil)
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

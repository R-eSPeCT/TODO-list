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

func (m *MockTodoRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetGroupedTodos(ctx context.Context, userID uuid.UUID) ([]models.TodoGroup, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.TodoGroup), args.Error(1)
}

func TestTodoHandler_Create(t *testing.T) {
	repo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test-key"))
	handler := NewTodoHandler(repo, jwtManager)
	app := fiber.New()
	app.Post("/api/todos", handler.CreateTodo)

	tests := []struct {
		name       string
		input      map[string]interface{}
		token      string
		setupMock  func()
		wantStatus int
	}{
		{
			name: "успешное создание",
			input: map[string]interface{}{
				"title":       "Тестовая задача",
				"description": "Описание тестовой задачи",
			},
			token: createTestToken(t, jwtManager),
			setupMock: func() {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "неверный формат запроса",
			input:      map[string]interface{}{},
			token:      createTestToken(t, jwtManager),
			setupMock:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "отсутствует токен",
			input:      map[string]interface{}{},
			token:      "",
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/todos", bytes.NewReader(body))
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			repo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_GetAll(t *testing.T) {
	repo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test-key"))
	handler := NewTodoHandler(repo, jwtManager)
	app := fiber.New()
	app.Get("/api/todos", handler.GetAll)

	userID := uuid.New()
	todos := []*models.Todo{
		{
			ID:          uuid.New(),
			Title:       "Задача 1",
			Description: "Описание 1",
			UserID:      userID,
		},
	}

	tests := []struct {
		name       string
		token      string
		setupMock  func()
		wantStatus int
	}{
		{
			name:  "успешное получение",
			token: createTestToken(t, jwtManager),
			setupMock: func() {
				repo.On("GetAll", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(todos, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "отсутствует токен",
			token:      "",
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			repo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_GetByID(t *testing.T) {
	repo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test-key"))
	handler := NewTodoHandler(repo, jwtManager)
	app := fiber.New()
	app.Get("/api/todos/:id", handler.GetByID)

	todoID := uuid.New()
	userID := uuid.New()
	todo := &models.Todo{
		ID:          todoID,
		Title:       "Тестовая задача",
		Description: "Описание задачи",
		UserID:      userID,
	}

	tests := []struct {
		name       string
		todoID     string
		token      string
		setupMock  func()
		wantStatus int
	}{
		{
			name:   "успешное получение",
			todoID: todoID.String(),
			token:  createTestToken(t, jwtManager),
			setupMock: func() {
				repo.On("GetByID", mock.Anything, todoID).Return(todo, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "неверный ID",
			todoID:     "invalid-id",
			token:      createTestToken(t, jwtManager),
			setupMock:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "отсутствует токен",
			todoID:     todoID.String(),
			token:      "",
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/api/todos/"+tt.todoID, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			repo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_Update(t *testing.T) {
	repo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test-key"))
	handler := NewTodoHandler(repo, jwtManager)
	app := fiber.New()
	app.Put("/api/todos/:id", handler.UpdateTodo)

	todoID := uuid.New()
	userID := uuid.New()
	todo := &models.Todo{
		ID:          todoID,
		Title:       "Тестовая задача",
		Description: "Описание задачи",
		UserID:      userID,
	}

	tests := []struct {
		name       string
		todoID     string
		input      map[string]interface{}
		token      string
		setupMock  func()
		wantStatus int
	}{
		{
			name:   "успешное обновление",
			todoID: todoID.String(),
			input: map[string]interface{}{
				"title":       "Обновленная задача",
				"description": "Обновленное описание",
			},
			token: createTestToken(t, jwtManager),
			setupMock: func() {
				repo.On("GetByID", mock.Anything, todoID).Return(todo, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "неверный ID",
			todoID:     "invalid-id",
			input:      map[string]interface{}{},
			token:      createTestToken(t, jwtManager),
			setupMock:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "отсутствует токен",
			todoID:     todoID.String(),
			input:      map[string]interface{}{},
			token:      "",
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPut, "/api/todos/"+tt.todoID, bytes.NewReader(body))
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			repo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_Delete(t *testing.T) {
	repo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test-key"))
	handler := NewTodoHandler(repo, jwtManager)
	app := fiber.New()
	app.Delete("/api/todos/:id", handler.DeleteTodo)

	todoID := uuid.New()
	userID := uuid.New()
	todo := &models.Todo{
		ID:          todoID,
		Title:       "Тестовая задача",
		Description: "Описание задачи",
		UserID:      userID,
	}

	tests := []struct {
		name       string
		todoID     string
		token      string
		setupMock  func()
		wantStatus int
	}{
		{
			name:   "успешное удаление",
			todoID: todoID.String(),
			token:  createTestToken(t, jwtManager),
			setupMock: func() {
				repo.On("GetByID", mock.Anything, todoID).Return(todo, nil)
				repo.On("Delete", mock.Anything, todoID).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "неверный ID",
			todoID:     "invalid-id",
			token:      createTestToken(t, jwtManager),
			setupMock:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "отсутствует токен",
			todoID:     todoID.String(),
			token:      "",
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodDelete, "/api/todos/"+tt.todoID, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			repo.AssertExpectations(t)
		})
	}
}

func TestTodoHandler_GetGroupedTodos(t *testing.T) {
	repo := new(MockTodoRepository)
	jwtManager := auth.NewJWTManager([]byte("test-key"))
	handler := NewTodoHandler(repo, jwtManager)
	app := fiber.New()
	app.Get("/api/todos/grouped", handler.GetGroupedTodos)

	userID := uuid.New()
	groupedTodos := []models.TodoGroup{
		{
			Status: "pending",
			Tasks: []*models.Todo{
				{
					ID:          uuid.New(),
					Title:       "Задача 1",
					Description: "Описание 1",
					UserID:      userID,
					Status:      "pending",
				},
			},
		},
		{
			Status: "completed",
			Tasks: []*models.Todo{
				{
					ID:          uuid.New(),
					Title:       "Задача 2",
					Description: "Описание 2",
					UserID:      userID,
					Status:      "completed",
				},
			},
		},
	}

	tests := []struct {
		name       string
		token      string
		setupMock  func()
		wantStatus int
	}{
		{
			name:  "успешное получение",
			token: createTestToken(t, jwtManager),
			setupMock: func() {
				repo.On("GetGroupedTodos", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(groupedTodos, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "отсутствует токен",
			token:      "",
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/api/todos/grouped", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			repo.AssertExpectations(t)
		})
	}
}

func createTestToken(t *testing.T, jwtManager *auth.JWTManager) string {
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}
	token, err := jwtManager.Generate(user)
	assert.NoError(t, err)
	return token
}

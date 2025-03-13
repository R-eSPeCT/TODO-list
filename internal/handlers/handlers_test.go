package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/R-eSPeCT/todo-list/internal/repository/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupTestServer создает тестовый сервер с моком репозитория
func setupTestServer(t *testing.T) (*fiber.App, *mocks.TodoRepository) {
	app := fiber.New()
	mockRepo := new(mocks.TodoRepository)
	todoHandler := NewTodoHandler(mockRepo)
	setupTodoRoutes(app, todoHandler)
	return app, mockRepo
}

func TestGetTodos(t *testing.T) {
	app, mockRepo := setupTestServer(t)

	// Тестирование получения всех задач
	t.Run("Get all todos", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.On("GetByUserID", mock.Anything, userID).Return([]*models.Todo{
			{ID: uuid.New(), Title: "Task 1", Description: "Description 1"},
			{ID: uuid.New(), Title: "Task 2", Description: "Description 2"},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var todos []*models.Todo
		err = json.NewDecoder(resp.Body).Decode(&todos)
		assert.NoError(t, err)
		assert.Len(t, todos, 2)
	})

	// Тестирование получения задачи по ID
	t.Run("Get todo by ID", func(t *testing.T) {
		todoID := uuid.New()
		userID := uuid.New()
		mockRepo.On("GetByID", mock.Anything, todoID).Return(&models.Todo{
			ID:          todoID,
			UserID:      userID,
			Title:       "Task 1",
			Description: "Description 1",
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/todos/"+todoID.String(), nil)
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var todo *models.Todo
		err = json.NewDecoder(resp.Body).Decode(&todo)
		assert.NoError(t, err)
		assert.Equal(t, todoID, todo.ID)
	})

	// Тестирование добавления новой задачи
	t.Run("Add new todo", func(t *testing.T) {
		userID := uuid.New()
		todo := &models.Todo{
			Title:       "New Task",
			Description: "New Description",
			UserID:      userID,
		}
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		body, _ := json.Marshal(todo)
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdTodo *models.Todo
		err = json.NewDecoder(resp.Body).Decode(&createdTodo)
		assert.NoError(t, err)
		assert.Equal(t, "New Task", createdTodo.Title)
	})

	// Тестирование обновления задачи
	t.Run("Update todo", func(t *testing.T) {
		todoID := uuid.New()
		userID := uuid.New()
		updatedTodo := &models.Todo{
			ID:          todoID,
			UserID:      userID,
			Title:       "Updated Task",
			Description: "Updated Description",
		}
		mockRepo.On("GetByID", mock.Anything, todoID).Return(updatedTodo, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		body, _ := json.Marshal(updatedTodo)
		req := httptest.NewRequest(http.MethodPut, "/todos/"+todoID.String(), bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedTodoResponse *models.Todo
		err = json.NewDecoder(resp.Body).Decode(&updatedTodoResponse)
		assert.NoError(t, err)
		assert.Equal(t, todoID, updatedTodoResponse.ID)
		assert.Equal(t, "Updated Task", updatedTodoResponse.Title)
	})

	// Тестирование удаления задачи
	t.Run("Delete todo", func(t *testing.T) {
		todoID := uuid.New()
		userID := uuid.New()
		mockRepo.On("GetByID", mock.Anything, todoID).Return(&models.Todo{
			ID:     todoID,
			UserID: userID,
		}, nil)
		mockRepo.On("Delete", mock.Anything, todoID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/todos/"+todoID.String(), nil)
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestGetTodo(t *testing.T) {
	app, mockRepo := setupTestServer(t)
	// ... existing code ...
}

func TestCreateTodo(t *testing.T) {
	app, mockRepo := setupTestServer(t)
	// ... existing code ...
}

func TestUpdateTodo(t *testing.T) {
	app, mockRepo := setupTestServer(t)
	// ... existing code ...
}

func TestDeleteTodo(t *testing.T) {
	app, mockRepo := setupTestServer(t)
	// ... existing code ...
}

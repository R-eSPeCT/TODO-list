package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/todo-list/internal/models"
	"github.com/yourusername/todo-list/internal/repository/mocks"
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
		mockRepo.On("GetAll").Return([]*Todo{
			{ID: 1, Title: "Task 1", Description: "Description 1"},
			{ID: 2, Title: "Task 2", Description: "Description 2"},
		}, nil)

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/todos", nil))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var todos []*Todo
		err = json.Unmarshal(resp.Body.Bytes(), &todos)
		assert.NoError(t, err)
		assert.Len(t, todos, 2)
	})

	// Тестирование получения задачи по ID
	t.Run("Get todo by ID", func(t *testing.T) {
		mockRepo.On("Get", 1).Return(&Todo{ID: 1, Title: "Task 1", Description: "Description 1"}, nil)

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/todos/1", nil))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var todo *Todo
		err = json.Unmarshal(resp.Body.Bytes(), &todo)
		assert.NoError(t, err)
		assert.Equal(t, 1, todo.ID)
	})

	// Тестирование добавления новой задачи
	t.Run("Add new todo", func(t *testing.T) {
		todo := &Todo{Title: "New Task", Description: "New Description"}
		mockRepo.On("Add", mock.Anything).Return(todo, nil)

		body, _ := json.Marshal(todo)
		resp, err := app.Test(httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body)))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdTodo *Todo
		err = json.Unmarshal(resp.Body.Bytes(), &createdTodo)
		assert.NoError(t, err)
		assert.Equal(t, "New Task", createdTodo.Title)
	})

	// Тестирование обновления задачи
	t.Run("Update todo", func(t *testing.T) {
		updatedTodo := &Todo{ID: 1, Title: "Updated Task", Description: "Updated Description"}
		mockRepo.On("Update", 1, mock.Anything).Return(updatedTodo, nil)

		body, _ := json.Marshal(updatedTodo)
		resp, err := app.Test(httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewBuffer(body)))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedTodoResponse *Todo
		err = json.Unmarshal(resp.Body.Bytes(), &updatedTodoResponse)
		assert.NoError(t, err)
		assert.Equal(t, 1, updatedTodoResponse.ID)
		assert.Equal(t, "Updated Task", updatedTodoResponse.Title)
	})

	// Тестирование удаления задачи
	t.Run("Delete todo", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil)

		resp, err := app.Test(httptest.NewRequest(http.MethodDelete, "/todos/1", nil))
		assert.NoError(t, err)
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

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourusername/todo-list/internal/models"
	"github.com/yourusername/todo-list/internal/repository"
)

// TodoHandler представляет собой обработчик HTTP-запросов для работы с задачами (Todo).
type TodoHandler struct {
	repo repository.TodoRepository
}

// NewTodoHandler создает новый экземпляр TodoHandler с использованием репозитория.
func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

// GetTodos обрабатывает GET-запрос для получения списка всех задач пользователя.
// Возвращает JSON-ответ с списком задач или ошибку, если что-то пошло не так.
func (h *TodoHandler) GetTodos(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	ctx := context.Background()
	todos, err := h.repo.GetByUserID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get todos",
		})
	}
	return c.JSON(todos)
}

// CreateTodo обрабатывает POST-запрос для создания новой задачи.
// Принимает JSON с данными задачи, валидирует их и сохраняет в базе данных.
// Возвращает статус 201 в случае успешного создания задачи.
func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if todo.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	todo.ID = uuid.New()
	todo.UserID = userID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = todo.CreatedAt

	if todo.Status == "" {
		todo.Status = "new"
	} else if !isValidStatus(todo.Status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be one of: new, in_progress, done",
		})
	}

	if todo.Priority == "" {
		todo.Priority = "medium"
	} else if !isValidPriority(todo.Priority) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid priority. Must be one of: low, medium, high",
		})
	}

	ctx := context.Background()
	if err := h.repo.Create(ctx, &todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create todo",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(todo)
}

// isValidStatus проверяет, является ли статус допустимым
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"new":         true,
		"in_progress": true,
		"done":        true,
	}
	return validStatuses[status]
}

// isValidPriority проверяет, является ли приоритет допустимым
func isValidPriority(priority string) bool {
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}
	return validPriorities[priority]
}

// UpdateTodo обрабатывает PUT-запрос для обновления существующей задачи.
// Принимает ID задачи из параметров запроса и JSON с новыми данными задачи.
// Обновляет задачу в базе данных и возвращает статус 200 в случае успеха.
func (h TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(400).SendString("Invalid user ID")
	}

	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid todo ID")
	}

	// Проверяем существование и принадлежность задачи пользователю
	existingTodo, err := h.repo.GetByID(c.Context(), todoID)
	if err != nil {
		return c.Status(404).SendString("Todo not found")
	}
	if existingTodo.UserID != userID {
		return c.Status(403).SendString("Access denied")
	}

	var updateTodo models.Todo
	if err := c.BodyParser(&updateTodo); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Обновляем только разрешенные поля
	existingTodo.Title = updateTodo.Title
	existingTodo.Description = updateTodo.Description
	existingTodo.Status = updateTodo.Status
	existingTodo.Priority = updateTodo.Priority
	existingTodo.UpdatedAt = time.Now()

	if !isValidStatus(existingTodo.Status) {
		return c.Status(400).SendString("Invalid status")
	}
	if !isValidPriority(existingTodo.Priority) {
		return c.Status(400).SendString("Invalid priority")
	}

	if err := h.repo.Update(c.Context(), existingTodo); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(existingTodo)
}

// DeleteTodo обрабатывает DELETE-запрос для удаления задачи.
// Принимает ID задачи из параметров запроса и удаляет задачу из базы данных.
// Возвращает статус 204 в случае успешного удаления.
func (h TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(400).SendString("Invalid user ID")
	}

	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid todo ID")
	}

	// Проверяем существование и принадлежность задачи пользователю
	todo, err := h.repo.GetByID(c.Context(), todoID)
	if err != nil {
		return c.Status(404).SendString("Todo not found")
	}
	if todo.UserID != userID {
		return c.Status(403).SendString("Access denied")
	}

	if err := h.repo.Delete(c.Context(), todoID); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(204)
}

// GetTodoByID обрабатывает GET-запрос для получения задачи по её ID.
// Принимает ID задачи из параметров запроса и возвращает JSON с данными задачи.
func (h TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(400).SendString("Invalid user ID")
	}

	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid todo ID")
	}

	todo, err := h.repo.GetByID(c.Context(), todoID)
	if err != nil {
		return c.Status(404).SendString("Todo not found")
	}

	if todo.UserID != userID {
		return c.Status(403).SendString("Access denied")
	}

	return c.JSON(todo)
}

// GetGroupedTodos обрабатывает GET-запрос для получения сгруппированных задач.
func (h TodoHandler) GetGroupedTodos(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(400).SendString("Invalid user ID")
	}

	groups, err := h.repo.GetGroupedTodos(c.Context(), userID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(groups)
}

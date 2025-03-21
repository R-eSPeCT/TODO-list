package handler

import (
	"context"
	"time"

	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/R-eSPeCT/todo-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

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

	if !isValidStatus(input.Status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	if !isValidPriority(input.Priority) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid priority",
		})
	}

	todo := &models.Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Status:      input.Status,
		Priority:    input.Priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Create(c.Context(), todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create todo",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

// isValidStatus проверяет, является ли статус допустимым
func isValidStatus(status string) bool {
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// isValidPriority проверяет, является ли приоритет допустимым
func isValidPriority(priority string) bool {
	validPriorities := []string{"low", "medium", "high"}
	for _, p := range validPriorities {
		if p == priority {
			return true
		}
	}
	return false
}

// UpdateTodo обрабатывает PUT-запрос для обновления существующей задачи.
func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID format",
		})
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	todo, err := h.repo.GetByID(c.Context(), todoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	todo.Title = input.Title
	todo.Description = input.Description
	todo.DueDate = input.DueDate
	todo.Status = input.Status
	todo.Priority = input.Priority
	todo.UpdatedAt = time.Now()

	if err := h.repo.Update(c.Context(), todo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update todo",
		})
	}

	return c.JSON(todo)
}

// DeleteTodo обрабатывает DELETE-запрос для удаления задачи.
func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID format",
		})
	}

	if err := h.repo.Delete(c.Context(), todoID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete todo",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetTodoByID обрабатывает GET-запрос для получения задачи по её ID.
func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID format",
		})
	}

	todo, err := h.repo.GetByID(c.Context(), todoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	return c.JSON(todo)
}

// GetGroupedTodos обрабатывает GET-запрос для получения сгруппированных задач.
func (h *TodoHandler) GetGroupedTodos(c *fiber.Ctx) error {
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

	todos, err := h.repo.GetByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get todos",
		})
	}

	// Группируем задачи по статусу
	grouped := make(map[string][]*models.Todo)
	for _, todo := range todos {
		grouped[todo.Status] = append(grouped[todo.Status], todo)
	}

	return c.JSON(grouped)
}

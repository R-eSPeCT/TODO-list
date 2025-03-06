package handlers

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TodoHandler представляет собой обработчик HTTP-запросов для работы с задачами (Todo).
type TodoHandler struct {
	repo *repository.TodoRepository // Репозиторий для взаимодействия с базой данных
}

// NewTodoHandler создает новый экземпляр TodoHandler с использованием пула соединений к базе данных.
func NewTodoHandler(db *pgxpool.Pool) *TodoHandler {
	return &TodoHandler{repo: repository.NewTodoRepository(db)}
}

// GetTodos обрабатывает GET-запрос для получения списка всех задач пользователя.
// Возвращает JSON-ответ с списком задач или ошибку, если что-то пошло не так.
func (h TodoHandler) GetTodos(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	todos, err := h.repo.GetTodos(userID)
	if err != nil {
		return c.Status(500).SendString(err.Error()) // В случае ошибки возвращаем статус 500 и текст ошибки
	}
	return c.JSON(todos) // Возвращаем список задач в формате JSON
}

// CreateTodo обрабатывает POST-запрос для создания новой задачи.
// Принимает JSON с данными задачи, валидирует их и сохраняет в базе данных.
// Возвращает статус 201 в случае успешного создания задачи.
func (h TodoHandler) CreateTodo(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Валидация обязательных полей
	if todo.Title == "" {
		return c.Status(400).SendString("Title is required")
	}
	if todo.Description == "" {
		return c.Status(400).SendString("Description is required")
	}
	if todo.Status == "" {
		todo.Status = "pending" // Устанавливаем статус по умолчанию
	} else if !isValidStatus(todo.Status) {
		return c.Status(400).SendString("Invalid status. Must be one of: pending, in_progress, completed")
	}
	if todo.Priority == "" {
		todo.Priority = "medium"
	} else if !isValidPriority(todo.Priority) {
		return c.Status(400).SendString("Invalid priority. Must be one of: low, medium, high")
	}

	if err := h.repo.CreateTodo(todo, userID); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(201)
}

// isValidStatus проверяет, является ли статус допустимым
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
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
	userID := c.Locals("userID").(int)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid ID") // В случае неверного ID возвращаем статус 400
	}

	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).SendString(err.Error()) // В случае ошибки парсинга возвращаем статус 400 и текст ошибки
	}

	if err := h.repo.UpdateTodo(id, todo, userID); err != nil {
		return c.Status(500).SendString(err.Error()) // В случае ошибки обновления возвращаем статус 500 и текст ошибки
	}
	return c.SendStatus(200) // Возвращаем статус 200 (OK)
}

// DeleteTodo обрабатывает DELETE-запрос для удаления задачи.
// Принимает ID задачи из параметров запроса и удаляет задачу из базы данных.
// Возвращает статус 204 в случае успешного удаления.
func (h TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid ID") // В случае неверного ID возвращаем статус 400
	}

	if err := h.repo.DeleteTodo(id, userID); err != nil {
		return c.Status(500).SendString(err.Error()) // В случае ошибки удаления возвращаем статус 500 и текст ошибки
	}
	return c.SendStatus(204) // Возвращаем статус 204 (No Content)
}

// GetTodoByID обрабатывает GET-запрос для получения задачи по её ID.
// Принимает ID задачи из параметров запроса и возвращает JSON с данными задачи.
func (h TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	todo, err := h.repo.GetTodoByID(id, userID)
	if err != nil {
		return c.Status(404).SendString("Task not found")
	}

	return c.JSON(todo)
}

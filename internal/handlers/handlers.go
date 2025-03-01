package handlers

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TodoHandler struct {
	repo *repository.TodoRepository
}

func NewTodoHandler(db *pgxpool.Pool) *TodoHandler {
	return &TodoHandler{repo: repository.NewTodoRepository(db)}
}

func (h TodoHandler) GetTodos(c *fiber.Ctx) error {
	todos, err := h.repo.GetTodos()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(todos)
}

func (h TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if err := h.repo.CreateTodo(todo); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(201)
}

func (h TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := h.repo.UpdateTodo(id, todo); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(200)
}

func (h TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	if err := h.repo.DeleteTodo(id); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendStatus(204)
}

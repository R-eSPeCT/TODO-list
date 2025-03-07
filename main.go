package main

import (
	"TODO-list/internal/handlers"
	"TODO-list/internal/repository"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

func main() {
	// Инициализация подключения к базе данных
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Инициализация репозиториев
	repos, err := repository.NewRepositories(dbpool)
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v\n", err)
	}

	// Инициализация обработчиков
	todoHandler := handlers.NewTodoHandler(repos.Todo)
	userHandler := handlers.NewUserHandler(repos.User)

	// Создание экземпляра Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Маршруты для пользователей
	users := app.Group("/api/users")
	users.Post("/register", userHandler.Register)
	users.Post("/login", userHandler.Login)

	// Маршруты для задач (требуют аутентификации)
	todos := app.Group("/api/todos")
	todos.Get("/", todoHandler.GetTodos)
	todos.Post("/", todoHandler.CreateTodo)
	todos.Get("/grouped", todoHandler.GetGroupedTodos)
	todos.Get("/:id", todoHandler.GetTodoByID)
	todos.Put("/:id", todoHandler.UpdateTodo)
	todos.Delete("/:id", todoHandler.DeleteTodo)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

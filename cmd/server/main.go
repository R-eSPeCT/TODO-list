package main

import (
	"TODO-list/internal/config"
	"TODO-list/internal/handlers"
	"TODO-list/internal/middleware"
	"TODO-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	//загружаю конфиг
	cfg := config.LoadConfig()

	// Подключение к базе данных
	db, err := repository.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация Fiber
	app := fiber.New()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Инициализация обработчиков
	todoHandler := handlers.NewTodoHandler(todoRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	// Публичные маршруты (без авторизации)
	app.Post("/users/register", userHandler.Register)
	app.Post("/users/login", userHandler.Login)

	// Защищенные маршруты (требуют авторизации)
	api := app.Group("/api", middleware.AuthMiddleware())

	// Маршруты для задач
	api.Get("/tasks", todoHandler.GetTodos)
	api.Post("/tasks", todoHandler.CreateTodo)
	api.Put("/tasks/:id", todoHandler.UpdateTodo)
	api.Delete("/tasks/:id", todoHandler.DeleteTodo)
	api.Get("/tasks/:id", todoHandler.GetTodoByID)
	api.Get("/tasks/grouped", todoHandler.GetGroupedTodos)

	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}

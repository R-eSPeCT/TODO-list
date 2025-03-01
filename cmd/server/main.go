package main

import (
	"TODO-list/internal/config"
	"TODO-list/internal/handlers"
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

	// Инициализация обработчиков
	todoHandler := handlers.NewTodoHandler(db)

	// Маршруты
	app.Get("/tasks", todoHandler.GetTodos)
	app.Post("/tasks", todoHandler.CreateTodo)
	app.Put("/tasks/:id", todoHandler.UpdateTodo)
	app.Delete("/tasks/:id", todoHandler.DeleteTodo)

	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}

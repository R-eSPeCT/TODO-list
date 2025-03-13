package main

import (
	"TODO-list/internal/config"
	"TODO-list/internal/handlers"
	"TODO-list/internal/middleware"
	"TODO-list/internal/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация Fiber
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 5,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

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

	// Обработка graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gracefully...")
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Запуск сервера
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

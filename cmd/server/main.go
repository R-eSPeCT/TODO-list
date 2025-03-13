package main

import (
	"context"
	"github.com/yourusername/todo-list/internal/auth"

	"github.com/yourusername/todo-list/internal/config"
	"github.com/yourusername/todo-list/internal/handlers"
	"github.com/yourusername/todo-list/internal/middleware"
	"github.com/yourusername/todo-list/internal/repository"
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

	// Инициализация JWT менеджера
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)

	// Инициализация Fiber
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 5,
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
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}))

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Инициализация обработчиков
	todoHandler := handlers.NewTodoHandler(todoRepo)
	userHandler := handlers.NewUserHandler(userRepo, jwtManager)

	// Публичные маршруты (без авторизации)
	app.Post("/users/register", userHandler.Register)
	app.Post("/users/login", userHandler.Login)

	// Защищенные маршруты (требуют авторизации)
	api := app.Group("/api", middleware.AuthMiddleware(jwtManager))

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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Запуск сервера
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

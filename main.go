package main

import (
	"TODO-list/internal/config"
	"TODO-list/internal/handlers"
	"TODO-list/internal/middleware"
	"TODO-list/internal/repository"
	"TODO-list/internal/services"
	"TODO-list/pkg/cache"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

func main() {
	// Загрузка конфигурации Redis
	redisConfig := config.NewRedisConfig()

	// Инициализация Redis
	redisCache, err := cache.NewRedisCache(redisConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()

	// Подключение к базе данных
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Инициализация репозиториев
	repos := repository.NewRepositories(dbpool)

	// Инициализация сервисов
	services := services.NewServices(repos)

	// Инициализация обработчиков
	handlers := handlers.NewHandler(services)

	// Создание Fiber приложения
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Rate limiting для API endpoints
	apiLimiter := middleware.RateLimit(redisCache, middleware.RateLimitConfig{
		Max:       100,       // 100 запросов
		Duration:  time.Hour, // за 1 час
		KeyPrefix: "rate_limit_api",
	})

	// Rate limiting для аутентификации
	authLimiter := middleware.RateLimit(redisCache, middleware.RateLimitConfig{
		Max:       5,                // 5 попыток
		Duration:  15 * time.Minute, // за 15 минут
		KeyPrefix: "rate_limit_auth",
	})

	// Роуты для пользователей
	users := app.Group("/api/users")
	users.Post("/register", handlers.Register)
	users.Post("/login", authLimiter, handlers.Login)

	// Роуты для задач с rate limiting
	todos := app.Group("/api/todos", apiLimiter)
	todos.Get("/", handlers.GetTodos)
	todos.Post("/", handlers.CreateTodo)
	todos.Get("/grouped", handlers.GetGroupedTodos)
	todos.Get("/:id", handlers.GetTodoByID)
	todos.Put("/:id", handlers.UpdateTodo)
	todos.Delete("/:id", handlers.DeleteTodo)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

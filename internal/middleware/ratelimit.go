package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/todo-list/pkg/cache"
)

// RateLimitConfig конфигурация для rate limiting
type RateLimitConfig struct {
	// Максимальное количество запросов за период
	Max int
	// Период времени для подсчета запросов
	Duration time.Duration
	// Префикс для ключей в Redis
	KeyPrefix string
}

// RateLimit создает middleware для ограничения частоты запросов
func RateLimit(cache cache.Cache, config RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем IP адрес клиента
		ip := c.IP()
		// Формируем ключ для Redis
		key := fmt.Sprintf("%s:%s", config.KeyPrefix, ip)

		// Проверяем существование ключа
		exists, err := cache.Exists(c.Context(), key)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limit check failed",
			})
		}

		if !exists {
			// Если ключ не существует, создаем его со значением 1
			err = cache.Set(c.Context(), key, 1, config.Duration)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Rate limit initialization failed",
				})
			}
		} else {
			// Увеличиваем счетчик
			count, err := cache.Increment(c.Context(), key)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Rate limit increment failed",
				})
			}

			// Если превышен лимит, возвращаем ошибку
			if count > int64(config.Max) {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
				})
			}
		}

		return c.Next()
	}
}

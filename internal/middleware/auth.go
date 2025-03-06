package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

// AuthMiddleware проверяет наличие и валидность токена авторизации
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем токен из заголовка
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).SendString("Unauthorized")
		}

		// Проверяем формат токена
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).SendString("Invalid authorization header")
		}

		// TODO: Здесь должна быть проверка JWT токена
		// Пока просто извлекаем userID из токена
		userID := 1 // Временное решение, нужно заменить на реальную проверку JWT

		// Сохраняем userID в контексте
		c.Locals("userID", userID)

		return c.Next()
	}
}

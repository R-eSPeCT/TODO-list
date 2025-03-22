package middleware

import (
	"strings"

	"github.com/R-eSPeCT/todo-list/internal/auth"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware создает middleware для проверки JWT токена в заголовке Authorization.
// Извлекает токен из заголовка, проверяет его валидность и добавляет ID пользователя в контекст.
func AuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Требуется заголовок Authorization",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Неверный формат заголовка Authorization",
			})
		}

		token := parts[1]
		claims, err := jwtManager.Validate(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Недействительный токен",
			})
		}

		// Добавляем ID пользователя в контекст для использования в следующих обработчиках
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}

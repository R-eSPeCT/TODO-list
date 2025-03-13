package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/todo-list/internal/auth"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]
		claims, err := jwtManager.Validate(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}

// validateToken проверяет JWT токен
func validateToken(tokenString string) (*jwt.Token, error) {
	// TODO: Реализовать проверку JWT токена
	return nil, nil
}

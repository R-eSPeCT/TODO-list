package handlers

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{repo: repository.NewUserRepository(db)}
}

// Register обрабатывает регистрацию нового пользователя
func (h UserHandler) Register(c *fiber.Ctx) error {
	var user models.UserCreate
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Проверяем обязательные поля
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).SendString("All fields are required")
	}

	// Проверяем длину пароля
	if len(user.Password) < 6 {
		return c.Status(400).SendString("Password must be at least 6 characters long")
	}

	newUser, err := h.repo.CreateUser(user)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(newUser)
}

// Login обрабатывает вход пользователя
func (h UserHandler) Login(c *fiber.Ctx) error {
	var login models.UserLogin
	if err := c.BodyParser(&login); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	user, err := h.repo.GetUserByEmail(login.Email)
	if err != nil {
		return c.Status(401).SendString("Invalid credentials")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password)); err != nil {
		return c.Status(401).SendString("Invalid credentials")
	}

	// В реальном приложении здесь нужно создать JWT токен
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

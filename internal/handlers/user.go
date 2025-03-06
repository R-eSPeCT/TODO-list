package handlers

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"regexp"
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
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валидация полей
	if err := validateUserCreate(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newUser, err := h.repo.CreateUser(user)
	if err != nil {
		// Проверяем на уникальные ограничения
		if isUniqueViolation(err) {
			return c.Status(400).JSON(fiber.Map{
				"error": "User with this email or username already exists",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(201).JSON(newUser)
}

// Login обрабатывает вход пользователя
func (h UserHandler) Login(c *fiber.Ctx) error {
	var login models.UserLogin
	if err := c.BodyParser(&login); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валидация email
	if !isValidEmail(login.Email) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	user, err := h.repo.GetUserByEmail(login.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// validateUserCreate проверяет корректность данных пользователя
func validateUserCreate(user models.UserCreate) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return fiber.NewError(400, "All fields are required")
	}

	if len(user.Username) < 3 || len(user.Username) > 50 {
		return fiber.NewError(400, "Username must be between 3 and 50 characters")
	}

	if !isValidEmail(user.Email) {
		return fiber.NewError(400, "Invalid email format")
	}

	if len(user.Password) < 6 {
		return fiber.NewError(400, "Password must be at least 6 characters long")
	}

	return nil
}

// isValidEmail проверяет корректность email
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isUniqueViolation проверяет, является ли ошибка нарушением уникального ограничения
func isUniqueViolation(err error) bool {
	// TODO: Добавить проверку на конкретный код ошибки PostgreSQL
	return err != nil
}

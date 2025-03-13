package handlers

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Register обрабатывает регистрацию нового пользователя
func (h UserHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем, существует ли пользователь
	existingUser, err := h.repo.GetByEmail(c.Context(), input.Email)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Создаем нового пользователя
	user := &models.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.Create(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// Login обрабатывает вход пользователя
func (h UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Получаем пользователя
	user, err := h.repo.GetByEmail(c.Context(), input.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Генерируем JWT токен
	token, err := generateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// generateJWT генерирует JWT токен для пользователя
func generateJWT(userID uuid.UUID) (string, error) {
	// TODO: Реализовать генерацию JWT токена
	return "", nil
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

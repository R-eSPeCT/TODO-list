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
	var userCreate models.UserCreate
	if err := c.BodyParser(&userCreate); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validateUserCreate(userCreate); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Проверяем, не существует ли уже пользователь с таким email
	existingUser, err := h.repo.GetByEmail(c.Context(), userCreate.Email)
	if err != nil && err != pgx.ErrNoRows {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to check email uniqueness",
		})
	}
	if existingUser != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// Создаем нового пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to process password",
		})
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     userCreate.Username,
		Email:        userCreate.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.Create(c.Context(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(201).JSON(&models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// Login обрабатывает вход пользователя
func (h UserHandler) Login(c *fiber.Ctx) error {
	var login models.UserLogin
	if err := c.BodyParser(&login); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if !isValidEmail(login.Email) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	user, err := h.repo.GetByEmail(c.Context(), login.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.JSON(&models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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

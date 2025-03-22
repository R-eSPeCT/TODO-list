package handler

import (
	"github.com/R-eSPeCT/todo-list/internal/auth"
	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/R-eSPeCT/todo-list/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

// UserHandler представляет собой обработчик HTTP-запросов для работы с пользователями.
type UserHandler struct {
	repo       repository.UserRepository
	jwtManager *auth.JWTManager
}

// NewUserHandler создает новый экземпляр UserHandler.
func NewUserHandler(repo repository.UserRepository, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Register обрабатывает регистрацию нового пользователя.
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем длину пароля
	if len(input.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters long",
		})
	}

	// Проверяем корректность email
	if !isValidEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
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
		ID:        uuid.New(),
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.Create(c.Context(), user); err != nil {
		if isUniqueViolation(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login обрабатывает вход пользователя.
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	// Проверяем существование пользователя
	user, err := h.repo.GetByEmail(c.Context(), input.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при поиске пользователя",
		})
	}
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Неверный email или пароль",
		})
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Неверный email или пароль",
		})
	}

	// Генерируем JWT токен
	token, err := h.jwtManager.Generate(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании токена",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

// GetProfile обрабатывает получение профиля пользователя
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат ID пользователя",
		})
	}

	user, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении профиля",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// UpdateProfile обрабатывает обновление профиля пользователя
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат ID пользователя",
		})
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	if input.Email != "" && !isValidEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат email",
		})
	}

	user, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении профиля",
		})
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при хешировании пароля",
			})
		}
		user.Password = string(hashedPassword)
	}

	user.UpdatedAt = time.Now()

	if err := h.repo.Update(c.Context(), user); err != nil {
		if isUniqueViolation(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email уже используется",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при обновлении профиля",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// isValidEmail проверяет корректность email.
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// isUniqueViolation проверяет, является ли ошибка нарушением уникального ограничения.
func isUniqueViolation(err error) bool {
	// Реализация зависит от используемой базы данных
	return false
}

package services

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, userCreate *models.UserCreate) error {
	// Валидация данных создания пользователя
	if err := validateUserCreate(userCreate); err != nil {
		return err
	}

	// Проверка существования пользователя с таким email
	existingUser, _ := s.repo.GetByEmail(ctx, userCreate.Email)
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Создание нового пользователя
	user := &models.User{
		ID:       uuid.New(),
		Username: userCreate.Username,
		Email:    userCreate.Email,
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	// Установка времени создания и обновления
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	// Проверка существования пользователя
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Обновление времени изменения
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Вспомогательные функции для валидации
func validateUserCreate(user *models.UserCreate) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func validateUser(user *models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

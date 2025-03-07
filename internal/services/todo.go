package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
	"vacancy/internal/models"
	"vacancy/internal/repository"
)

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{
		repo: repo,
	}
}

func (s *todoService) Create(ctx context.Context, todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title is required")
	}

	// Генерация UUID если не установлен
	if todo.ID == uuid.Nil {
		todo.ID = uuid.New()
	}

	// Установка времени создания и обновления
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	// Установка статуса по умолчанию
	if todo.Status == "" {
		todo.Status = "new"
	}

	return s.repo.Create(ctx, todo)
}

func (s *todoService) GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *todoService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *todoService) Update(ctx context.Context, todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title is required")
	}

	// Проверка валидности статуса
	if !isValidStatus(todo.Status) {
		return errors.New("invalid status")
	}

	// Обновление времени изменения
	todo.UpdatedAt = time.Now()

	return s.repo.Update(ctx, todo)
}

func (s *todoService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Вспомогательная функция для проверки статуса
func isValidStatus(status string) bool {
	validStatuses := []string{"new", "in_progress", "done"}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

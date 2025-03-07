package services

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
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

	// Установка приоритета по умолчанию
	if todo.Priority == "" {
		todo.Priority = "medium"
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

	// Проверка валидности статуса и приоритета
	if !isValidStatus(todo.Status) {
		return errors.New("invalid status")
	}
	if !isValidPriority(todo.Priority) {
		return errors.New("invalid priority")
	}

	// Обновление времени изменения
	todo.UpdatedAt = time.Now()

	return s.repo.Update(ctx, todo)
}

func (s *todoService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *todoService) GetGroupedTodos(ctx context.Context, userID uuid.UUID) ([]models.TodoGroup, error) {
	return s.repo.GetGroupedTodos(ctx, userID)
}

// Вспомогательные функции для валидации
func isValidStatus(status string) bool {
	validStatuses := []string{"new", "in_progress", "done"}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

func isValidPriority(priority string) bool {
	validPriorities := []string{"low", "medium", "high"}
	for _, p := range validPriorities {
		if priority == p {
			return true
		}
	}
	return false
}

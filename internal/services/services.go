package services

import (
	"TODO-list/internal/models"
	"TODO-list/internal/repository"
	"context"
	"github.com/google/uuid"
)

type Services struct {
	User UserService
	Todo TodoService
}

type UserService interface {
	Create(ctx context.Context, userCreate *models.UserCreate) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TodoService interface {
	Create(ctx context.Context, todo *models.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Todo, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetGroupedTodos(ctx context.Context, userID uuid.UUID) ([]models.TodoGroup, error)
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		User: NewUserService(repos.User),
		Todo: NewTodoService(repos.Todo),
	}
}

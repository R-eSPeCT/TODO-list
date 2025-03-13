package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	// Подключаемся к тестовой базе данных
	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/todo_list_test?sslmode=disable")
	require.NoError(t, err)
	require.NotNil(t, pool)

	// Очищаем таблицы перед каждым тестом
	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE todos CASCADE")
	require.NoError(t, err)

	return pool
}

func TestTodoRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewTodoRepository(pool)

	// Создаем тестового пользователя
	userID := uuid.New()
	_, err := pool.Exec(context.Background(), "INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		userID, "test@example.com", "password123", time.Now(), time.Now())
	require.NoError(t, err)

	tests := []struct {
		name    string
		todo    *Todo
		wantErr bool
	}{
		{
			name: "valid todo",
			todo: &Todo{
				ID:          uuid.New(),
				UserID:      userID,
				Title:       "Test Todo",
				Description: "Test Description",
				Status:      "pending",
				DueDate:     time.Now().Add(24 * time.Hour),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		 {
			name: "invalid user id",
			todo: &Todo{
				ID:          uuid.New(),
				UserID:      uuid.New(), // Несуществующий пользователь
				Title:       "Test Todo",
				Description: "Test Description",
				Status:      "pending",
				DueDate:     time.Now().Add(24 * time.Hour),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.todo)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Проверяем, что задача создана
			todo, err := repo.GetByID(context.Background(), tt.todo.ID)
			assert.NoError(t, err)
			assert.Equal(t, tt.todo.Title, todo.Title)
		})
	}
}

func TestTodoRepository_GetByUserID(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewTodoRepository(pool)

	// Создаем тестового пользователя
	userID := uuid.New()
	_, err := pool.Exec(context.Background(), "INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		userID, "test@example.com", "password123", time.Now(), time.Now())
	require.NoError(t, err)

	// Создаем тестовые задачи
	todo1 := &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       "Todo 1",
		Description: "Description 1",
		Status:      "pending",
		DueDate:     time.Now().Add(24 * time.Hour),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	todo2 := &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       "Todo 2",
		Description: "Description 2",
		Status:      "completed",
		DueDate:     time.Now().Add(48 * time.Hour),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Create(context.Background(), todo1)
	require.NoError(t, err)
	err = repo.Create(context.Background(), todo2)
	require.NoError(t, err)

	todos, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, todos, 2)

	// Проверяем содержимое задач
	todoTitles := make(map[string]bool)
	for _, todo := range todos {
		todoTitles[todo.Title] = true
	}
	assert.True(t, todoTitles["Todo 1"])
	assert.True(t, todoTitles["Todo 2"])
}

func TestTodoRepository_Update(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewTodoRepository(pool)

	// Создаем тестового пользователя
	userID := uuid.New()
	_, err := pool.Exec(context.Background(), "INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		userID, "test@example.com", "password123", time.Now(), time.Now())
	require.NoError(t, err)

	// Создаем тестовую задачу
	todo := &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       "Test Todo",
		Description: "Test Description",
		Status:      "pending",
		DueDate:     time.Now().Add(24 * time.Hour),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.Create(context.Background(), todo)
	require.NoError(t, err)

	// Обновляем задачу
	todo.Title = "Updated Todo"
	todo.Status = "completed"
	err = repo.Update(context.Background(), todo)
	assert.NoError(t, err)

	// Проверяем обновление
	updatedTodo, err := repo.GetByID(context.Background(), todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Todo", updatedTodo.Title)
	assert.Equal(t, "completed", updatedTodo.Status)
}

func TestTodoRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewTodoRepository(pool)

	// Создаем тестового пользователя
	userID := uuid.New()
	_, err := pool.Exec(context.Background(), "INSERT INTO users (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		userID, "test@example.com", "password123", time.Now(), time.Now())
	require.NoError(t, err)

	// Создаем тестовую задачу
	todo := &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       "Test Todo",
		Description: "Test Description",
		Status:      "pending",
		DueDate:     time.Now().Add(24 * time.Hour),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.Create(context.Background(), todo)
	require.NoError(t, err)

	// Удаляем задачу
	err = repo.Delete(context.Background(), todo.ID)
	assert.NoError(t, err)

	// Проверяем, что задача удалена
	_, err = repo.GetByID(context.Background(), todo.ID)
	assert.Error(t, err)
} 
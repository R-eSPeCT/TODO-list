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
	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE users CASCADE")
	require.NoError(t, err)

	return pool
}

func TestUserRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &User{
				ID:        uuid.New(),
				Email:     "test@example.com",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &User{
				ID:        uuid.New(),
				Email:     "test@example.com", // Тот же email
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Проверяем, что пользователь создан
			user, err := repo.GetByID(context.Background(), tt.user.ID)
			assert.NoError(t, err)
			assert.Equal(t, tt.user.Email, user.Email)
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	// Создаем тестового пользователя
	testUser := &User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "non-existing user",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByEmail(context.Background(), tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, testUser.Email, user.Email)
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	// Создаем тестового пользователя
	testUser := &User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Обновляем пользователя
	testUser.Email = "updated@example.com"
	err = repo.Update(context.Background(), testUser)
	assert.NoError(t, err)

	// Проверяем обновление
	updatedUser, err := repo.GetByID(context.Background(), testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated@example.com", updatedUser.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	// Создаем тестового пользователя
	testUser := &User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Удаляем пользователя
	err = repo.Delete(context.Background(), testUser.ID)
	assert.NoError(t, err)

	// Проверяем, что пользователь удален
	_, err = repo.GetByID(context.Background(), testUser.ID)
	assert.Error(t, err)
} 
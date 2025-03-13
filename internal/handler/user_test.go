package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestApp(t *testing.T) *fiber.App {
	app := fiber.New()
	return app
}

func TestUserHandler_Register(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	app.Post("/register", handler.Register)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
		setupMock  func()
	}{
		{
			name: "valid registration",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*User")).Return(nil)
			},
		},
		{
			name: "invalid email",
			payload: map[string]interface{}{
				"email":    "invalid-email",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	app.Post("/login", handler.Login)

	// Создаем тестового пользователя
	userID := uuid.New()
	email := "test@example.com"
	password := "password123"

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
		setupMock  func()
	}{
		{
			name: "valid login",
			payload: map[string]interface{}{
				"email":    email,
				"password": password,
			},
			wantStatus: http.StatusOK,
			setupMock: func() {
				user := &User{
					ID:       userID,
					Email:    email,
					Password: password,
				}
				mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
			},
		},
		{
			name: "invalid email",
			payload: map[string]interface{}{
				"email":    "invalid@example.com",
				"password": password,
			},
			wantStatus: http.StatusUnauthorized,
			setupMock: func() {
				mockRepo.On("GetByEmail", mock.Anything, "invalid@example.com").Return(nil, nil)
			},
		},
		{
			name: "invalid password",
			payload: map[string]interface{}{
				"email":    email,
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
			setupMock: func() {
				user := &User{
					ID:       userID,
					Email:    email,
					Password: password,
				}
				mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	app.Get("/profile", handler.GetProfile)

	// Генерируем валидный токен
	userID := uuid.New()
	jwtManager := NewJWTManager("test-secret-key", "1h")
	token, err := jwtManager.Generate(userID)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name:       "valid request",
			token:      token,
			wantStatus: http.StatusOK,
			setupMock: func() {
				user := &User{
					ID:    userID,
					Email: "test@example.com",
				}
				mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
			},
		},
		{
			name:       "missing token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	app := setupTestApp(t)
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	app.Put("/profile", handler.UpdateProfile)

	// Генерируем валидный токен
	userID := uuid.New()
	jwtManager := NewJWTManager("test-secret-key", "1h")
	token, err := jwtManager.Generate(userID)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		payload    map[string]interface{}
		wantStatus int
		setupMock  func()
	}{
		{
			name:  "valid update",
			token: token,
			payload: map[string]interface{}{
				"email": "updated@example.com",
			},
			wantStatus: http.StatusOK,
			setupMock: func() {
				user := &User{
					ID:    userID,
					Email: "test@example.com",
				}
				mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*User")).Return(nil)
			},
		},
		{
			name:       "missing token",
			token:      "",
			payload:    map[string]interface{}{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "invalid email",
			token: token,
			payload: map[string]interface{}{
				"email": "invalid-email",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

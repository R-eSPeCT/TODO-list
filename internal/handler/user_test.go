package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/R-eSPeCT/todo-list/internal/auth"
	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestApp(t *testing.T, repo *MockUserRepository) *fiber.App {
	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	h := NewUserHandler(repo, jwtManager)

	app := fiber.New()
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
	app.Get("/profile", h.GetProfile)
	app.Put("/profile", h.UpdateProfile)
	return app
}

func TestUserHandler_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	app := setupTestApp(t, mockRepo)

	tests := []struct {
		name       string
		input      map[string]string
		wantStatus int
		setupMock  func()
	}{
		{
			name: "успешная регистрация",
			input: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			},
		},
		{
			name: "invalid email",
			input: map[string]string{
				"email":    "invalid-email",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			input: map[string]string{
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

			jsonBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	app := setupTestApp(t, mockRepo)

	userID := uuid.New()
	email := "test@example.com"
	password := "password123"
	hashedPassword := "$2a$10$..." // Replace with actual hashed password

	tests := []struct {
		name       string
		input      map[string]string
		wantStatus int
		setupMock  func()
	}{
		{
			name: "успешный вход",
			input: map[string]string{
				"email":    email,
				"password": password,
			},
			wantStatus: http.StatusOK,
			setupMock: func() {
				user := &models.User{
					ID:       userID,
					Email:    email,
					Password: hashedPassword,
				}
				mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
			},
		},
		{
			name: "invalid email",
			input: map[string]string{
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
			input: map[string]string{
				"email":    email,
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
			setupMock: func() {
				user := &models.User{
					ID:       userID,
					Email:    email,
					Password: hashedPassword,
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

			jsonBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	app := setupTestApp(t, mockRepo)

	userID := uuid.New()
	testUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	token, err := jwtManager.Generate(testUser)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		wantStatus int
		setupMock  func()
	}{
		{
			name:       "успешное получение профиля",
			token:      token,
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetByID", mock.Anything, userID).Return(testUser, nil)
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

			req := httptest.NewRequest("GET", "/profile", nil)
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
	mockRepo := new(MockUserRepository)
	app := setupTestApp(t, mockRepo)

	userID := uuid.New()
	testUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	jwtManager := auth.NewJWTManager([]byte("test_secret"))
	token, err := jwtManager.Generate(testUser)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		input      map[string]string
		wantStatus int
		setupMock  func()
	}{
		{
			name:  "успешное обновление профиля",
			token: token,
			input: map[string]string{
				"email": "newemail@example.com",
			},
			wantStatus: http.StatusOK,
			setupMock: func() {
				mockRepo.On("GetByID", mock.Anything, userID).Return(testUser, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			},
		},
		{
			name:       "missing token",
			token:      "",
			input:      map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "invalid email",
			token: token,
			input: map[string]string{
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

			jsonBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("PUT", "/profile", bytes.NewReader(jsonBody))
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

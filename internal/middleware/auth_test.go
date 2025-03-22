package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/R-eSPeCT/todo-list/internal/auth"
	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	app := fiber.New()
	jwtManager := auth.NewJWTManager([]byte("test-secret-key"))
	app.Use(AuthMiddleware(jwtManager))

	// Генерируем валидный токен
	userID := uuid.New()
	user := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}
	validToken, err := jwtManager.Generate(user)
	require.NoError(t, err)

	// Тестовый обработчик
	app.Get("/test", func(c *fiber.Ctx) error {
		ctxUserID := c.Locals("userID")
		assert.NotNil(t, ctxUserID)
		assert.IsType(t, uuid.UUID{}, ctxUserID)
		assert.Equal(t, userID, ctxUserID)
		return c.SendStatus(http.StatusOK)
	})

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "valid token",
			token:      validToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token format",
			token:      "invalid-token-format",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestCorsMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(CorsMiddleware())

	// Тестовый обработчик
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	tests := []struct {
		name             string
		origin           string
		method           string
		wantAllowOrigin  string
		wantAllowMethods string
		wantAllowHeaders string
		wantStatus       int
	}{
		{
			name:             "allowed origin",
			origin:           "http://localhost:3000",
			method:           http.MethodGet,
			wantAllowOrigin:  "http://localhost:3000",
			wantAllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			wantAllowHeaders: "Origin,Content-Type,Accept,Authorization",
			wantStatus:       http.StatusOK,
		},
		{
			name:             "disallowed origin",
			origin:           "http://malicious.com",
			method:           http.MethodGet,
			wantAllowOrigin:  "",
			wantAllowMethods: "",
			wantAllowHeaders: "",
			wantStatus:       http.StatusOK,
		},
		{
			name:             "preflight request",
			origin:           "http://localhost:3000",
			method:           http.MethodOptions,
			wantAllowOrigin:  "http://localhost:3000",
			wantAllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			wantAllowHeaders: "Origin,Content-Type,Accept,Authorization",
			wantStatus:       http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, tt.wantAllowOrigin, resp.Header.Get("Access-Control-Allow-Origin"))
			assert.Equal(t, tt.wantAllowMethods, resp.Header.Get("Access-Control-Allow-Methods"))
			assert.Equal(t, tt.wantAllowHeaders, resp.Header.Get("Access-Control-Allow-Headers"))
		})
	}
}

func TestLoggerMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	// Тестовый обработчик
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверяем, что логгер не падает при ошибке
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(http.StatusInternalServerError, "test error")
	})

	req = httptest.NewRequest(http.MethodGet, "/error", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

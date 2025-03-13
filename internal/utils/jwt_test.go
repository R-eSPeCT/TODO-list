package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Проверяем, что токен можно разобрать
	claims, err := ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestParseToken(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, userID, claims.UserID)
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID)
	require.NoError(t, err)

	// Проверяем, что токен действителен
	claims, err := ParseToken(token)
	require.NoError(t, err)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
}

func TestTokenClaims(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID)
	require.NoError(t, err)

	claims, err := ParseToken(token)
	require.NoError(t, err)

	// Проверяем все поля claims
	assert.Equal(t, userID, claims.UserID)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
	assert.True(t, claims.IssuedAt <= time.Now().Unix())
} 
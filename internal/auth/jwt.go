package auth

import (
	"fmt"
	"time"

	"github.com/R-eSPeCT/todo-list/internal/models"
	"github.com/golang-jwt/jwt"
)

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey []byte
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey []byte) *JWTManager {
	return &JWTManager{secretKey: secretKey}
}

// Claims представляет структуру JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// Generate создает новый JWT токен для пользователя
func (m *JWTManager) Generate(user *models.User) (string, error) {
	claims := Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// Validate проверяет JWT токен и возвращает claims
func (m *JWTManager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

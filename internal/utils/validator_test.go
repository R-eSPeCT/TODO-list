package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "test@sub.example.com",
			wantErr: false,
		},
		{
			name:    "invalid email without @",
			email:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email without domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "Password123!",
			wantErr:  false,
		},
		{
			name:     "password without uppercase",
			password: "password123!",
			wantErr:  true,
		},
		{
			name:     "password without lowercase",
			password: "PASSWORD123!",
			wantErr:  true,
		},
		{
			name:     "password without numbers",
			password: "Password!",
			wantErr:  true,
		},
		{
			name:     "password without special characters",
			password: "Password123",
			wantErr:  true,
		},
		{
			name:     "short password",
			password: "Pass1!",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateTodo(t *testing.T) {
	tests := []struct {
		name    string
		todo    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid todo",
			todo: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
				"status":      "pending",
				"due_date":    time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "missing title",
			todo: map[string]interface{}{
				"description": "Test Description",
				"status":      "pending",
				"due_date":    time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "empty title",
			todo: map[string]interface{}{
				"title":       "",
				"description": "Test Description",
				"status":      "pending",
				"due_date":    time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			todo: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
				"status":      "invalid",
				"due_date":    time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "past due date",
			todo: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
				"status":      "pending",
				"due_date":    time.Now().Add(-24 * time.Hour),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTodo(tt.todo)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
} 
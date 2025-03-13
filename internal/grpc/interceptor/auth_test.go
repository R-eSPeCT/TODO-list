package interceptor

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type testServer struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *testServer) Context() context.Context {
	return s.ctx
}

func TestAuthInterceptor_Unary(t *testing.T) {
	jwtManager := NewJWTManager("test-secret-key", "1h")
	interceptor := NewAuthInterceptor(jwtManager)

	tests := []struct {
		name    string
		token   string
		method  string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "valid-token",
			method:  "/todo.TodoService/CreateTodo",
			wantErr: false,
		},
		{
			name:    "public method",
			token:   "",
			method:  "/todo.TodoService/Register",
			wantErr: false,
		},
		{
			name:    "missing token",
			token:   "",
			method:  "/todo.TodoService/CreateTodo",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "invalid-token",
			method:  "/todo.TodoService/CreateTodo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.token != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
					"authorization": "Bearer " + tt.token,
				}))
			}

			info := &grpc.UnaryServerInfo{
				FullMethod: tt.method,
			}

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			}

			_, err := interceptor.Unary()(ctx, nil, info, handler)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, codes.Unauthenticated, status.Code(err))
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAuthInterceptor_Stream(t *testing.T) {
	jwtManager := NewJWTManager("test-secret-key", "1h")
	interceptor := NewAuthInterceptor(jwtManager)

	tests := []struct {
		name    string
		token   string
		method  string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "valid-token",
			method:  "/todo.TodoService/StreamTodos",
			wantErr: false,
		},
		{
			name:    "public method",
			token:   "",
			method:  "/todo.TodoService/Register",
			wantErr: false,
		},
		{
			name:    "missing token",
			token:   "",
			method:  "/todo.TodoService/StreamTodos",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "invalid-token",
			method:  "/todo.TodoService/StreamTodos",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.token != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
					"authorization": "Bearer " + tt.token,
				}))
			}

			info := &grpc.StreamServerInfo{
				FullMethod: tt.method,
			}

			handler := func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			}

			stream := &testServer{ctx: ctx}
			err := interceptor.Stream()(nil, stream, info, handler)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, codes.Unauthenticated, status.Code(err))
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAuthInterceptor_Authorize(t *testing.T) {
	jwtManager := NewJWTManager("test-secret-key", "1h")
	interceptor := NewAuthInterceptor(jwtManager)

	userID := uuid.New()
	token, err := jwtManager.Generate(userID)
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
			ctx := context.Background()
			if tt.token != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
					"authorization": "Bearer " + tt.token,
				}))
			}

			userID, err := interceptor.authorize(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, userID)
		})
	}
} 
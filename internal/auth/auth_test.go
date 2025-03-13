package auth

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/todo-list/internal/models"
	"github.com/yourusername/todo-list/internal/repository/mocks"
	pb "github.com/yourusername/todo-list/pkg/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type testServer struct {
	t       *testing.T
	lis     *bufconn.Listener
	server  *GRPCServer
	conn    *grpc.ClientConn
	cleanup func()
}

func NewTestServer(t *testing.T) *testServer {
	lis := bufconn.Listen(bufSize)
	server := NewGRPCServer(nil, nil) // Здесь можно добавить моки репозиториев если нужно
	go func() {
		if err := server.Serve(lis); err != nil {
			t.Errorf("server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	cleanup := func() {
		conn.Close()
		lis.Close()
		server.Stop()
	}

	return &testServer{
		t:       t,
		lis:     lis,
		server:  server,
		conn:    conn,
		cleanup: cleanup,
	}
}

func (s *testServer) Cleanup() {
	s.cleanup()
}

func (s *testServer) ClientConn() *grpc.ClientConn {
	return s.conn
}

func TestGRPCServer_Register(t *testing.T) {
	// Создаем тестовый сервер
	grpcTestServer := NewTestServer(t)
	defer grpcTestServer.Cleanup()

	// Создаем тестовый клиент
	conn := grpcTestServer.ClientConn()
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	tests := []struct {
		name    string
		req     *pb.RegisterRequest
		wantErr bool
	}{
		{
			name: "valid registration",
			req: &pb.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			req: &pb.RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			req: &pb.RegisterRequest{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			resp, err := client.Register(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.NotEmpty(t, resp.Token)
		})
	}
}

func TestGRPCServer_Login(t *testing.T) {
	// Создаем тестовый сервер
	grpcTestServer := NewTestServer(t)
	defer grpcTestServer.Cleanup()

	// Создаем тестовый клиент
	conn := grpcTestServer.ClientConn()
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	tests := []struct {
		name    string
		req     *pb.LoginRequest
		wantErr bool
	}{
		{
			name: "valid login",
			req: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			req: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "non-existent user",
			req: &pb.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			resp, err := client.Login(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.NotEmpty(t, resp.Token)
		})
	}
}

func TestGRPCServer_ValidateToken(t *testing.T) {
	// Создаем тестовый сервер
	grpcTestServer := NewTestServer(t)
	defer grpcTestServer.Cleanup()

	// Создаем тестовый клиент
	conn := grpcTestServer.ClientConn()
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "valid-token", // Здесь нужно будет добавить реальный токен
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			resp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{
				Token: tt.token,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.True(t, resp.Valid)
		})
	}
}

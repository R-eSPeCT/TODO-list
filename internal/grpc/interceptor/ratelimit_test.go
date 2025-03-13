package interceptor

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
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

func setupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// Очищаем тестовую базу данных
	err := client.FlushDB(context.Background()).Err()
	require.NoError(t, err)

	return client
}

func TestRateLimitInterceptor_Unary(t *testing.T) {
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	config := &RateLimitConfig{
		Cache:           redisClient,
		MaxRequestCount: 2,
		Duration:        time.Second * 5,
		KeyPrefix:       "test:ratelimit:",
	}
	interceptor := NewRateLimitInterceptor(config)

	tests := []struct {
		name       string
		method     string
		clientAddr string
		requests   int
		wantErr    bool
	}{
		{
			name:       "within limit",
			method:     "/todo.TodoService/CreateTodo",
			clientAddr: "127.0.0.1",
			requests:   1,
			wantErr:    false,
		},
		{
			name:       "exceed limit",
			method:     "/todo.TodoService/CreateTodo",
			clientAddr: "127.0.0.1",
			requests:   3,
			wantErr:    true,
		},
		{
			name:       "different methods",
			method:     "/todo.TodoService/GetTodo",
			clientAddr: "127.0.0.1",
			requests:   1,
			wantErr:    false,
		},
		{
			name:       "different clients",
			method:     "/todo.TodoService/CreateTodo",
			clientAddr: "127.0.0.2",
			requests:   1,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
				"x-forwarded-for": tt.clientAddr,
			}))

			info := &grpc.UnaryServerInfo{
				FullMethod: tt.method,
			}

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			}

			for i := 0; i < tt.requests; i++ {
				_, err := interceptor.Unary()(ctx, nil, info, handler)
				if i == tt.requests-1 && tt.wantErr {
					assert.Error(t, err)
					assert.Equal(t, codes.ResourceExhausted, status.Code(err))
					return
				}
				assert.NoError(t, err)
				time.Sleep(time.Millisecond * 100) // Небольшая задержка между запросами
			}
		})
	}
}

func TestRateLimitInterceptor_Stream(t *testing.T) {
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	config := &RateLimitConfig{
		Cache:           redisClient,
		MaxRequestCount: 2,
		Duration:        time.Second * 5,
		KeyPrefix:       "test:ratelimit:",
	}
	interceptor := NewRateLimitInterceptor(config)

	tests := []struct {
		name       string
		method     string
		clientAddr string
		requests   int
		wantErr    bool
	}{
		{
			name:       "within limit",
			method:     "/todo.TodoService/StreamTodos",
			clientAddr: "127.0.0.1",
			requests:   1,
			wantErr:    false,
		},
		{
			name:       "exceed limit",
			method:     "/todo.TodoService/StreamTodos",
			clientAddr: "127.0.0.1",
			requests:   3,
			wantErr:    true,
		},
		{
			name:       "different methods",
			method:     "/todo.TodoService/StreamCompleted",
			clientAddr: "127.0.0.1",
			requests:   1,
			wantErr:    false,
		},
		{
			name:       "different clients",
			method:     "/todo.TodoService/StreamTodos",
			clientAddr: "127.0.0.2",
			requests:   1,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
				"x-forwarded-for": tt.clientAddr,
			}))

			info := &grpc.StreamServerInfo{
				FullMethod: tt.method,
			}

			handler := func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			}

			for i := 0; i < tt.requests; i++ {
				stream := &testServer{ctx: ctx}
				err := interceptor.Stream()(nil, stream, info, handler)
				if i == tt.requests-1 && tt.wantErr {
					assert.Error(t, err)
					assert.Equal(t, codes.ResourceExhausted, status.Code(err))
					return
				}
				assert.NoError(t, err)
				time.Sleep(time.Millisecond * 100) // Небольшая задержка между запросами
			}
		})
	}
}

func TestRateLimitInterceptor_KeyGeneration(t *testing.T) {
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	config := &RateLimitConfig{
		Cache:           redisClient,
		MaxRequestCount: 2,
		Duration:        time.Second * 5,
		KeyPrefix:       "test:ratelimit:",
	}
	interceptor := NewRateLimitInterceptor(config)

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
		"x-forwarded-for": "127.0.0.1",
	}))

	info := &grpc.UnaryServerInfo{
		FullMethod: "/todo.TodoService/CreateTodo",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	// Первый запрос
	_, err := interceptor.Unary()(ctx, nil, info, handler)
	assert.NoError(t, err)

	// Проверяем, что ключ создан в Redis
	key := config.KeyPrefix + "127.0.0.1:/todo.TodoService/CreateTodo"
	exists, err := redisClient.Exists(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	// Проверяем значение счетчика
	count, err := redisClient.Get(ctx, key).Int()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Проверяем TTL
	ttl, err := redisClient.TTL(ctx, key).Result()
	assert.NoError(t, err)
	assert.True(t, ttl > 0 && ttl <= config.Duration)
}

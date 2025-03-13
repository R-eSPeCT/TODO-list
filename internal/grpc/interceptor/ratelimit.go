package interceptor

import (
	"context"
	"fmt"
	"github.com/R-eSPeCT/todo-list/pkg/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

// RateLimitConfig конфигурация для rate limiting
type RateLimitConfig struct {
	Cache     cache.Cache
	Max       int           // Максимальное количество запросов
	Duration  time.Duration // Период времени
	KeyPrefix string        // Префикс для ключей в Redis
}

// RateLimitInterceptor представляет интерцептор для rate limiting
type RateLimitInterceptor struct {
	config RateLimitConfig
}

// NewRateLimitInterceptor создает новый интерцептор rate limiting
func NewRateLimitInterceptor(config RateLimitConfig) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		config: config,
	}
}

// Unary возвращает унарный интерцептор для rate limiting
func (i *RateLimitInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Получаем IP адрес клиента
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "cannot get peer info")
		}

		// Формируем ключ для Redis
		key := fmt.Sprintf("%s:%s:%s", i.config.KeyPrefix, info.FullMethod, p.Addr.String())

		// Проверяем существование ключа
		exists, err := i.config.Cache.Exists(ctx, key)
		if err != nil {
			return nil, status.Error(codes.Internal, "rate limit check failed")
		}

		if !exists {
			// Если ключ не существует, создаем его со значением 1
			err = i.config.Cache.Set(ctx, key, 1, i.config.Duration)
			if err != nil {
				return nil, status.Error(codes.Internal, "rate limit initialization failed")
			}
		} else {
			// Увеличиваем счетчик
			count, err := i.config.Cache.Increment(ctx, key)
			if err != nil {
				return nil, status.Error(codes.Internal, "rate limit increment failed")
			}

			// Если превышен лимит, возвращаем ошибку
			if count > int64(i.config.Max) {
				return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
			}
		}

		return handler(ctx, req)
	}
}

// Stream возвращает стрим интерцептор для rate limiting
func (i *RateLimitInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Получаем IP адрес клиента
		p, ok := peer.FromContext(stream.Context())
		if !ok {
			return status.Error(codes.Internal, "cannot get peer info")
		}

		// Формируем ключ для Redis
		key := fmt.Sprintf("%s:%s:%s", i.config.KeyPrefix, info.FullMethod, p.Addr.String())

		// Проверяем существование ключа
		exists, err := i.config.Cache.Exists(stream.Context(), key)
		if err != nil {
			return status.Error(codes.Internal, "rate limit check failed")
		}

		if !exists {
			// Если ключ не существует, создаем его со значением 1
			err = i.config.Cache.Set(stream.Context(), key, 1, i.config.Duration)
			if err != nil {
				return status.Error(codes.Internal, "rate limit initialization failed")
			}
		} else {
			// Увеличиваем счетчик
			count, err := i.config.Cache.Increment(stream.Context(), key)
			if err != nil {
				return status.Error(codes.Internal, "rate limit increment failed")
			}

			// Если превышен лимит, возвращаем ошибку
			if count > int64(i.config.Max) {
				return status.Error(codes.ResourceExhausted, "rate limit exceeded")
			}
		}

		return handler(srv, stream)
	}
}

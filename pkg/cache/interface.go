package cache

import (
	"context"
	"time"
)

// Cache определяет интерфейс для работы с кэшем
type Cache interface {
	// Get получает значение из кэша по ключу
	Get(ctx context.Context, key string) ([]byte, error)

	// Set устанавливает значение в кэш с опциональным временем жизни
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete удаляет значение из кэша по ключу
	Delete(ctx context.Context, key string) error

	// Exists проверяет существование ключа в кэше
	Exists(ctx context.Context, key string) (bool, error)

	// Increment увеличивает значение счетчика
	Increment(ctx context.Context, key string) (int64, error)

	// SetNX устанавливает значение, только если ключ не существует
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)

	// Close закрывает соединение с кэшем
	Close() error
}

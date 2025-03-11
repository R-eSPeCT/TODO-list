package cache

import (
	"TODO-list/internal/config"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisCache реализует интерфейс Cache с использованием Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache создает новый экземпляр RedisCache
func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var bytes []byte
	var err error

	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		bytes, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	return c.client.Set(ctx, key, bytes, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *RedisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	var bytes []byte
	var err error

	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		bytes, err = json.Marshal(value)
		if err != nil {
			return false, err
		}
	}

	return c.client.SetNX(ctx, key, bytes, ttl).Result()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

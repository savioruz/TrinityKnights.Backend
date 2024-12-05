package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type ImplCache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *ImplCache {
	return &ImplCache{
		client: client,
	}
}

var (
	ErrCacheMiss   = errors.New("cache: key not found")
	ErrCacheFailed = errors.New("cache: failed to get key")
	ErrUnmarshal   = errors.New("cache: failed to unmarshal data")
	ErrMarshal     = errors.New("cache: failed to marshal data")
)

func (c *ImplCache) Get(key string, value interface{}) error {
	data, err := c.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	} else if err != nil {
		return ErrCacheFailed
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return ErrUnmarshal
	}

	return nil
}

func (c *ImplCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrMarshal
	}

	if _, err := c.client.Set(context.Background(), key, data, expiration).Result(); err != nil {
		return ErrCacheFailed
	}

	return nil
}

func (c *ImplCache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *ImplCache) DeletePattern(pattern string) error {
	keys, err := c.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(context.Background(), keys...).Err()
	}

	return nil
}

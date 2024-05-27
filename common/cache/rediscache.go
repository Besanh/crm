package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cache *redis.Client
}

var ctx = context.Background()

func NewRedisCache(client *redis.Client) IRedisCache {
	return &RedisCache{
		cache: client,
	}
}

func valueToString(value any) (string, error) {
	tmp, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(tmp), nil
}

func (c *RedisCache) Set(key string, value any) error {
	_, err := c.cache.Set(ctx, key, value, redis.KeepTTL).Result()
	return err
}

func (c *RedisCache) SetTTL(key string, value any, ttl time.Duration) error {
	valueStr, err := valueToString(value)
	if err != nil {
		return err
	}
	_, err = c.cache.Set(ctx, key, valueStr, ttl).Result()
	return err
}

func (c *RedisCache) Get(key string) (string, error) {
	value, err := c.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

func (c *RedisCache) Close() {
	c.cache.Close()
}

func (c *RedisCache) Del(key string) error {
	_, err := c.cache.Del(ctx, key).Result()
	return err
}
func (c *RedisCache) HSet(key string, field string, value any) error {
	tmp, err := json.Marshal(value)
	if err != nil {
		return err
	}
	data := []any{field, string(tmp)}
	_, err = c.cache.HMSet(ctx, key, data).Result()
	return err
}

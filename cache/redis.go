package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache represents a Redis cache client.
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new instance of RedisCache.
func NewRedisCache(addr, password string, db int) *RedisCache {
	// Initialize Redis client options
	opt := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}

	// Create a new Redis client
	client := redis.NewClient(opt)

	// Ping the Redis server to check connectivity
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

// Set adds a new key-value pair to the Redis cache.
func (rc *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	err := rc.client.Set(rc.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s in Redis: %w", key, err)
	}
	return nil
}

// Get retrieves a value from the Redis cache based on the key.
func (rc *RedisCache) Get(key string) (interface{}, error) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s not found in Redis", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key %s from Redis: %w", key, err)
	}
	return val, nil
}

// Delete removes a key-value pair from the Redis cache.
func (rc *RedisCache) Delete(key string) error {
	err := rc.client.Del(rc.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s from Redis: %w", key, err)
	}
	return nil
}

// Close closes the Redis client connection.
func (rc *RedisCache) Close() error {
	err := rc.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	return nil
}

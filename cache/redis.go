/*
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
func NewRedisCache(addr string, db int) *RedisCache { //password string

		// Initialize Redis client options
		opt := &redis.Options{
			Addr: addr,
			//Password: password,
			DB: db,
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
	 previously used
*/

/* working
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
func NewRedisCache(addr string, db int) *RedisCache {
	// Initialize Redis client options
	opt := &redis.Options{
		Addr: addr,
		DB:   db,
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
func (rc *RedisCache) Get(key string) (string, error) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found in Redis", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s from Redis: %w", key, err)
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
*/

/* --- working one
package cache

import (
	"context"
	"fmt"
	"time"

	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisCache represents a Redis cache client.
/* .......working
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new instance of RedisCache.
func NewRedisCache(addr string, db int) *RedisCache { //password string
	// Initialize Redis client options
	opt := &redis.Options{
		Addr: addr,
		//Password: password,
		DB: db,
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
func (rc *RedisCache) Get(key string) (string, error) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found in Redis", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s from Redis: %w", key, err)
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
.....*/
// redis.go

// RedisCache represents a Redis cache client.

// redis.go

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

// RedisCache represents a Redis cache client.
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.RWMutex
}

// NewRedisCache creates a new instance of RedisCache.
func NewRedisCache(addr string, db int) *RedisCache {
	opt := &redis.Options{
		Addr: addr,
		DB:   db,
	}

	client := redis.NewClient(opt) //creates new redis client

	ctx := context.Background()
	_, err := client.Ping(ctx).Result() //checks whether it connected to redis server
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
	rc.mu.Lock()         //
	defer rc.mu.Unlock() //
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value to JSON: %w", err)
	}

	err = rc.client.Set(rc.ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s in Redis: %w", key, err)
	}
	return nil
}

// Get retrieves a value from the Redis cache based on the key.
func (rc *RedisCache) Get(key string) (string, error) {
	rc.mu.Lock()         //
	defer rc.mu.Unlock() //
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found in Redis", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s from Redis: %w", key, err)
	}
	return val, nil
}

// Delete removes a key-value pair from the Redis cache.
func (rc *RedisCache) Delete(key string) error {
	rc.mu.Lock()         //
	defer rc.mu.Unlock() //
	err := rc.client.Del(rc.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s from Redis: %w", key, err)
	}
	return nil
}

// Close closes the Redis client connection.
func (rc *RedisCache) Close() error {
	rc.mu.Lock()         //
	defer rc.mu.Unlock() //
	err := rc.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	return nil
}

// GetAllKeys retrieves all keys from the Redis cache.
func (rc *RedisCache) GetAllKeys() ([]string, error) {
	rc.mu.Lock()         //
	defer rc.mu.Unlock() //
	keys, err := rc.client.Keys(rc.ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all keys from Redis: %w", err)
	}
	return keys, nil
}

// DeleteAllKeys deletes all keys from the Redis cache.
func (rc *RedisCache) DeleteAllKeys() error {
	rc.mu.Lock()         //
	defer rc.mu.Unlock() //
	keys, err := rc.client.Keys(rc.ctx, "*").Result()
	if err != nil {
		return fmt.Errorf("failed to get all keys from Redis: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	err = rc.client.Del(rc.ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete all keys from Redis: %w", err)
	}
	return nil
}

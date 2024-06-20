package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}
type InMemoryCache struct {
	cache map[string]*cacheItem
}
type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]*cacheItem),
	}
}
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	expiration := time.Now().Add(ttl)
	c.cache[key] = &cacheItem{
		value:      value,
		expiration: expiration,
	}
	return nil
}
func (c *InMemoryCache) Get(key string) (interface{}, error) {
	item, found := c.cache[key]
	if !found {
		return nil, fmt.Errorf("key not found in cache")
	}
	if item.expiration.Before(time.Now()) {
		return nil, fmt.Errorf(" key got expired")
	}
	return item.value, nil
}
func (c *InMemoryCache) Delete(key string) error {
	delete(c.cache, key)
	return nil
}

type Rediscache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*Rediscache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)

	}
	return &Rediscache{
		client: client,
	}, nil
}

func (c *Rediscache) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s in RedisCache %v", key, err)
	}
	return nil
}
func (c *Rediscache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s not found in Redis", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch key %s from Redis: %v", key, err)
	}
	return val, nil
}

func (c *Rediscache) Delete(key string) error {
	ctx := context.Background()
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("Failed to Delete the key %s from redis: %v", key, err)
	}
	return nil

}

type MemcachedCache struct {
	client *memcache.Client
}

func NewMemcachedCache(addr string) *MemcachedCache {
	return &MemcachedCache{
		client: memcache.New(addr),
	}
}
func (c *MemcachedCache) Set(key string, value interface{}, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(fmt.Sprintf("%v", value)),
		Expiration: int32(ttl.Seconds()),
	}
	err := c.client.Set(item)
	if err != nil {
		return fmt.Errorf("failed to set key %s in Memcached: %v", key, err)
	}
	return nil
}
func (c *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := c.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, fmt.Errorf("key %s not found in Memcached", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key %s from Memcached: %v", key, err)
	}
	return string(item.Value), nil
}

// Delete removes a key-value pair from the Memcached cache.
func (c *MemcachedCache) Delete(key string) error {
	err := c.client.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete key %s from Memcached: %v", key, err)
	}
	return nil
}

// NewCache creates a new instance of Cache based on the provided configuration.
func NewCache(cacheType string, config interface{}) (Cache, error) {
	switch cacheType {
	case "inmemory":
		return NewInMemoryCache(), nil
	case "redis":
		redisConfig, ok := config.(RedisConfig)
		if !ok {
			return nil, fmt.Errorf("invalid Redis configuration")
		}
		return NewRedisCache(redisConfig.Addr, redisConfig.Password, redisConfig.DB)
	case "memcached":
		memcachedConfig, ok := config.(MemcachedConfig)
		if !ok {
			return nil, fmt.Errorf("invalid Memcached configuration")
		}
		return NewMemcachedCache(memcachedConfig.Addr), nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}

// RedisConfig holds configuration options for connecting to Redis.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// MemcachedConfig holds configuration options for connecting to Memcached.
type MemcachedConfig struct {
	Addr string
}

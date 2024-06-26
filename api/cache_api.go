package cache

/*
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
*/
//it defines the cache_interface with commom methods like get,set,getallkeys,setallkeys,deleteallkeys...
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
	GetAllKeys() ([]string, error)
	DeleteAllKeys() error
}

type InMemoryCache struct {
	cache map[string]*cacheItem //sores cached items in memory using a map
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
		delete(c.cache, key)
		return nil, fmt.Errorf("key expired")
	}
	return item.value, nil
}

func (c *InMemoryCache) Delete(key string) error {
	delete(c.cache, key)
	return nil
}

func (c *InMemoryCache) GetAllKeys() ([]string, error) {
	keys := make([]string, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}
	return keys, nil
}

func (c *InMemoryCache) DeleteAllKeys() error {
	c.cache = make(map[string]*cacheItem)
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
		return fmt.Errorf("failed to set key %s in RedisCache: %v", key, err)
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
		return fmt.Errorf("failed to delete key %s from Redis: %v", key, err)
	}
	return nil
}

func (c *Rediscache) GetAllKeys() ([]string, error) {
	ctx := context.Background()
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all keys from Redis: %v", err)
	}
	return keys, nil
}

func (c *Rediscache) DeleteAllKeys() error {
	ctx := context.Background()
	keys, err := c.GetAllKeys()
	if err != nil {
		return fmt.Errorf("failed to fetch keys for deletion: %v", err)
	}
	for _, key := range keys {
		err = c.client.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("failed to delete key %s from Redis: %v", key, err)
		}
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

func (c *MemcachedCache) Delete(key string) error {
	err := c.client.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete key %s from Memcached: %v", key, err)
	}
	return nil
}

func (c *MemcachedCache) GetAllKeys() ([]string, error) {
	// Memcached does not support fetching all keys out of the box.
	return nil, fmt.Errorf("GetAllKeys not supported for Memcached")
}

func (c *MemcachedCache) DeleteAllKeys() error {
	// Memcached does not support deleting all keys out of the box.
	return fmt.Errorf("DeleteAllKeys not supported for Memcached")
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

/*
import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	GetAllKeys() ([]string, error)
	DeleteAllKeys() error
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
		delete(c.cache, key)
		return nil, fmt.Errorf("key expired")
	}
	return item.value, nil
}

func (c *InMemoryCache) Delete(key string) error {
	delete(c.cache, key)
	return nil
}

func (c *InMemoryCache) GetAllKeys() ([]string, error) {
	keys := make([]string, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}
	return keys, nil
}

func (c *InMemoryCache) DeleteAllKeys() error {
	c.cache = make(map[string]*cacheItem)
	return nil
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
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
	return &RedisCache{
		client: client,
	}, nil
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s in RedisCache: %v", key, err)
	}
	return nil
}

func (c *RedisCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s not found in Redis", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch key %s from Redis: %v", key, err)
	}
	return val, nil
}

func (c *RedisCache) Delete(key string) error {
	ctx := context.Background()
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s from Redis: %v", key, err)
	}
	return nil
}

func (c *RedisCache) GetAllKeys() ([]string, error) {
	ctx := context.Background()
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all keys from Redis: %v", err)
	}
	return keys, nil
}

func (c *RedisCache) DeleteAllKeys() error {
	ctx := context.Background()
	keys, err := c.GetAllKeys()
	if err != nil {
		return fmt.Errorf("failed to fetch keys for deletion: %v", err)
	}
	for _, key := range keys {
		err = c.client.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("failed to delete key %s from Redis: %v", key, err)
		}
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

func (c *MemcachedCache) Delete(key string) error {
	err := c.client.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete key %s from Memcached: %v", key, err)
	}
	return nil
}

func (c *MemcachedCache) GetAllKeys() ([]string, error) {
	// Memcached does not support fetching all keys out of the box.
	return nil, fmt.Errorf("GetAllKeys not supported for Memcached")
}

func (c *MemcachedCache) DeleteAllKeys() error {
	// Memcached does not support deleting all keys out of the box.
	return fmt.Errorf("DeleteAllKeys not supported for Memcached")
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

var cache Cache

func main() {
	// Choose the cache type and initialize it
	cacheType := "inmemory" // Change to "redis" or "memcached" as needed
	var err error
	switch cacheType {
	case "redis":
		cache, err = NewCache(cacheType, RedisConfig{Addr: "localhost:6379", Password: "", DB: 0})
	case "memcached":
		cache, err = NewCache(cacheType, MemcachedConfig{Addr: "localhost:11211"})
	default:
		cache, err = NewCache(cacheType, nil)
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize cache: %v", err))
	}

	r := gin.Default()

	r.POST("/caching-library/cache/user", setCache)
	r.GET("/caching-library/cache/user/:key", getCache)
	r.DELETE("/caching-library/cache/user/:key", deleteCache)
	r.GET("/caching-library/cache/keys", getAllKeys)
	r.DELETE("/caching-library/cache/keys", deleteAllKeys)

	r.Run(":8080")
}

func setCache(c *gin.Context) {
	var entry CacheEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := entry.Username
	err := cache.Set(key, entry, time.Hour) // Setting cache entry with 1 hour TTL
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getCache(c *gin.Context) {
	key := c.Param("key")
	value, err := cache.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": value})
}

func deleteCache(c *gin.Context) {
	key := c.Param("key")
	err := cache.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getAllKeys(c *gin.Context) {
	keys, err := cache.GetAllKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func deleteAllKeys(c *gin.Context) {
	err := cache.DeleteAllKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type CacheEntry struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
*/

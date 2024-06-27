package handlers

//uses gin framework to handle http request for caching data across multipple caching system
import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ragavcr7/caching-library/cache"
)

type CacheHandler struct {
	memcachedCache *cache.MemcachedCache
	redisCache     *cache.RedisCache
	inMemoryCache  *cache.InMemoryCache
	lruCache       *cache.LRUCache
}

func NewCacheHandler(memcachedCache *cache.MemcachedCache, redisCache *cache.RedisCache, inMemoryCache *cache.InMemoryCache, lruCache *cache.LRUCache) *CacheHandler {
	return &CacheHandler{
		memcachedCache: memcachedCache,
		redisCache:     redisCache,
		inMemoryCache:  inMemoryCache,
		lruCache:       lruCache,
	}
}

func (ch *CacheHandler) SetupRoutes(router *gin.Engine) {
	router.POST("/cache/:key", ch.SetCache)
	router.GET("/cache/:key", ch.GetCache)
	router.DELETE("/cache/:key", ch.DeleteCache)
	router.GET("/cache", ch.GetAllCache)
	router.DELETE("/cache", ch.DeleteAllCache) // d added for testing
}

func (ch *CacheHandler) SetCache(c *gin.Context) {
	key := c.Param("key") //retrieves key parameter from url

	var value interface{}
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	// Cache in Memcached
	expirationMem := 5 * time.Minute
	if err := ch.memcachedCache.Set(key, value, expirationMem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to set key %s in Memcached: %s", key, err.Error())})
		return
	}

	// Cache in Redis
	expirationRedis := 5 * time.Minute
	valueJSON, err := json.Marshal(value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to marshal value to JSON: %s", err.Error())})
		return
	}
	if err := ch.redisCache.Set(key, string(valueJSON), expirationRedis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to set key %s in Redis: %s", key, err.Error())})
		return
	}

	// Cache in InMemory
	expirationInMem := 5 * time.Minute
	ch.inMemoryCache.Set(key, value, expirationInMem)

	// Cache in LRUCache
	ch.lruCache.Set(key, value, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully cached value with key %s", key)})
}

func (ch *CacheHandler) GetCache(c *gin.Context) {
	key := c.Param("key")

	// Try retrieving from LRUCache first
	if value, found := ch.lruCache.Get(key); found {
		c.JSON(http.StatusOK, gin.H{"value": value})
		return
	}

	// If not found in LRUCache, try Memcached
	memcachedValue, err := ch.memcachedCache.Get(key)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"value": memcachedValue})
		return
	}

	// If not found in Memcached, try Redis
	redisValue, err := ch.redisCache.Get(key)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"value": redisValue})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Key %s not found in cache", key)})
}

func (ch *CacheHandler) DeleteCache(c *gin.Context) {
	key := c.Param("key")

	// Delete from Memcached
	if err := ch.memcachedCache.Delete(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete key %s from Memcached: %s", key, err.Error())})
		return
	}

	// Delete from Redis
	if err := ch.redisCache.Delete(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete key %s from Redis: %s", key, err.Error())})
		return
	}

	// Delete from InMemory
	ch.inMemoryCache.Delete(key)

	// Delete from LRUCache
	ch.lruCache.Remove(key)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully deleted key %s from all caches", key)})
}

func (ch *CacheHandler) GetAllCache(c *gin.Context) {
	// Fetch all from InMemory
	inMemValues := ch.inMemoryCache.GetAllKeys()

	// Fetch all from LRUCache
	lruValues := ch.lruCache.GetAll()

	c.JSON(http.StatusOK, gin.H{
		"inMemoryCache": inMemValues,
		"lruCache":      lruValues,
	})
}

func (ch *CacheHandler) DeleteAllCache(c *gin.Context) {
	// Delete all from Memcached
	if err := ch.memcachedCache.DeleteAllKeys(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete all keys from Memcached: %s", err.Error())})
		return
	}

	// Delete all from Redis
	if err := ch.redisCache.DeleteAllKeys(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete all keys from Redis: %s", err.Error())})
		return
	}

	// Clear InMemory cache
	ch.inMemoryCache.DeleteAllKeys()

	// Clear LRUCache
	ch.lruCache.Clear()

	c.JSON(http.StatusOK, gin.H{"message": "Successfully cleared all caches"})
}

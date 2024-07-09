package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ragavcr7/caching-library/cache"
)

// Example user(data) structure for demonstration
/*
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
*/
func main() {
	// Capacity for InMemoryCache
	capacity := 3
	inMemoryCache := cache.NewInMemoryCache(capacity)
	memcachedAddr := "localhost:11211"
	// Initialize LRU cache
	lruCapacity := 3 // Adjust capacity as needed
	lruCache := cache.NewLRUCacheWithMemcached(lruCapacity, memcachedAddr)

	// Create the Gin router
	router := gin.Default()

	// Initialize API handlers
	cacheHandler := NewCacheHandler(inMemoryCache, lruCache)

	// Routes for caching endpoints
	cacheHandler.SetupRoutes(router)
	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

type CacheHandler struct {
	inMemoryCache *cache.InMemoryCache
	lruCache      *cache.LRUCache
}

func NewCacheHandler(inMemoryCache *cache.InMemoryCache, lruCache *cache.LRUCache) *CacheHandler {
	return &CacheHandler{
		inMemoryCache: inMemoryCache,
		lruCache:      lruCache,
	}
}

func (ch *CacheHandler) SetupRoutes(router *gin.Engine) {
	router.POST("/cache/:key", ch.SetCache)
	router.GET("/cache/:key", ch.GetCache)
	router.DELETE("/cache/:key", ch.DeleteCache)
	router.GET("/cache", ch.GetAllCache)
	router.DELETE("/cache", ch.DeleteAllCache)
}

func (ch *CacheHandler) SetCache(c *gin.Context) {
	key := c.Param("key")

	var value interface{}
	if err := c.ShouldBindJSON(&value); err != nil { //binding json payload to value coz the request send by user will be in the form of json so we have to map it to go's struct to access it.
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
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

	// If not found in LRUCache, try InMemoryCache
	if value, found := ch.inMemoryCache.Get(key); found {
		c.JSON(http.StatusOK, gin.H{"value": value})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("key %s not found in caches", key)})
}

func (ch *CacheHandler) DeleteCache(c *gin.Context) {
	key := c.Param("key")

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
	// Clear InMemory cache
	ch.inMemoryCache.DeleteAllKeys()

	// Clear LRUCache
	ch.lruCache.Clear()

	c.JSON(http.StatusOK, gin.H{"message": "Successfully cleared all caches"})
}

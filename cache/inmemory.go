// Package cache implements an in-memory cache with LRU eviction policy.
package cache

import (
	"container/list"
	"sync"
	"time"
)

// CacheItem represents an item stored in the cache.
type CacheItem struct {
	key        string
	value      interface{}
	expiration *time.Time
}

// InMemoryCache represents an in-memory cache with LRU eviction.
type InMemoryCache struct {
	cache         map[string]*list.Element // Map for fast access to cache items
	evictionList  *list.List               // Doubly linked list to manage LRU eviction
	maxSize       int                      // Maximum number of items the cache can hold
	mutex         sync.Mutex               // Mutex for thread safety
	expiration    time.Duration            // Default expiration time for cache items
	cleanupTicker *time.Ticker             // Ticker for periodic cleanup of expired items
	stopCleanup   chan bool                // Channel to stop the cleanup routine
}

// NewInMemoryCache creates a new instance of InMemoryCache.
func NewInMemoryCache(maxSize int, expiration time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		cache:        make(map[string]*list.Element),
		evictionList: list.New(),
		maxSize:      maxSize,
		expiration:   expiration,
		stopCleanup:  make(chan bool),
	}

	// Start a goroutine to periodically clean up expired cache items
	cache.startCleanupRoutine()

	return cache
}

// startCleanupRoutine starts a goroutine to periodically clean up expired cache items.
func (c *InMemoryCache) startCleanupRoutine() {
	c.cleanupTicker = time.NewTicker(c.expiration)
	go func() {
		for {
			select {
			case <-c.cleanupTicker.C:
				c.cleanup()
			case <-c.stopCleanup:
				c.cleanupTicker.Stop()
				return
			}
		}
	}()
}

// cleanup removes expired cache items from the cache.
func (c *InMemoryCache) cleanup() {
	now := time.Now()
	c.mutex.Lock()         // To lock the resource
	defer c.mutex.Unlock() // will get executed at the end of this function since we have dont using the resource

	for {
		element := c.evictionList.Back()
		if element == nil { //denotes empty linked list
			break
		}

		item := element.Value.(*CacheItem)
		if item.expiration != nil && item.expiration.Before(now) {
			c.evictionList.Remove(element)
			delete(c.cache, item.key)
		} else {
			break
		}
	}
}

// Set adds a new key-value pair to the cache or updates the value if the key already exists.
func (c *InMemoryCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the key already exists AND IF YES IT UPDATES ITS VALUE
	if element, exists := c.cache[key]; exists {
		c.evictionList.MoveToFront(element)
		element.Value.(*CacheItem).value = value //VALUE UPDATING
		return
	}

	// If cache is full, evict the least recently used item
	if len(c.cache) >= c.maxSize {
		c.evictLRU()
	}

	// Add new item to cache and eviction list
	expiration := time.Now().Add(c.expiration)
	item := &CacheItem{
		key:        key,
		value:      value,
		expiration: &expiration,
	}
	element := c.evictionList.PushFront(item)
	c.cache[key] = element
}

// Get retrieves a value from the cache based on the key.
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.cache[key]; exists {
		c.evictionList.MoveToFront(element)
		return element.Value.(*CacheItem).value, true
	}
	return nil, false
}

// Delete removes a key-value pair from the cache.
func (c *InMemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.cache[key]; exists {
		c.evictionList.Remove(element)
		delete(c.cache, key)
	}
}

// evictLRU removes the least recently used item from the cache.
func (c *InMemoryCache) evictLRU() {
	if element := c.evictionList.Back(); element != nil {
		item := c.evictionList.Remove(element).(*CacheItem)
		delete(c.cache, item.key)
	}
}

// StopCleanup stops the periodic cleanup of expired cache items.
func (c *InMemoryCache) StopCleanup() {
	c.stopCleanup <- true
}

// Len returns the current number of items in the cache.
func (c *InMemoryCache) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.cache)
}

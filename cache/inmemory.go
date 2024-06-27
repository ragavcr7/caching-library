// Package cache implements an in-memory cache with LRU eviction policy.
package cache

/*
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
	expiration    time.Time            // Default expiration time for cache items
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
func (c *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) {
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
	expiration = c.expiration
	item := &CacheItem{
		key:        key,
		value:      value,
		expiration: &expiration,
	}
	element := c.evictionList.PushFront(item)
	c.cache[key] = element
}

// Get retrieves a value from the cache based on the key.
func (c *InMemoryCache) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.cache[key]; exists {
		c.evictionList.MoveToFront(element)
		return element.Value.(*CacheItem).value
	}
	return nil
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
*/
// inmemory.go
import (
	"container/list"
	"sync"
	"time"
)

// entry represents a cache entry stored in the InMemoryCache.
type entry struct {
	key        string
	value      interface{}
	expiration time.Time
}

// InMemoryCache represents an in-memory cache with LRU eviction and expiration handling.
type InMemoryCache struct {
	capacity  int
	cache     map[string]*list.Element
	evictList *list.List
	mu        sync.Mutex
}

// NewInMemoryCache creates a new instance of InMemoryCache with the specified capacity.
func NewInMemoryCache(capacity int) *InMemoryCache {
	return &InMemoryCache{
		capacity:  capacity,
		cache:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Set adds a new key-value pair to the in-memory cache with optional expiration.
func (imc *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	// Check if the key exists in the cache
	if elem, ok := imc.cache[key]; ok {
		// If key exists, update the value, expiration, and move the element to front (MRU position)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).expiration = time.Now().Add(expiration)
		imc.evictList.MoveToFront(elem)
	} else {
		// If key does not exist, add the new entry to the cache
		if len(imc.cache) >= imc.capacity {
			// Evict least recently used element if cache is full
			imc.evictOldest()
		}

		// Add new entry to the cache with expiration time
		newElem := imc.evictList.PushFront(&entry{
			key:        key,
			value:      value,
			expiration: time.Now().Add(expiration),
		})
		imc.cache[key] = newElem
	}

	// Start a goroutine to periodically check and evict expired entries
	go imc.startEvictionTicker()
}

// Get retrieves a value from the in-memory cache based on the key.
func (imc *InMemoryCache) Get(key string) (interface{}, bool) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		// Check if the entry is expired
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			// If expired, delete the entry and return not found
			imc.removeElement(elem)
			return nil, false
		}

		// Move the accessed element to front (MRU position)
		imc.evictList.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

// Delete removes a key-value pair from the in-memory cache.
func (imc *InMemoryCache) Delete(key string) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		imc.removeElement(elem)
	}
}

// GetAllKeys retrieves all keys from the in-memory cache.
func (imc *InMemoryCache) GetAllKeys() []string {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	keys := make([]string, 0, len(imc.cache))
	for key := range imc.cache {
		keys = append(keys, key)
	}
	return keys
}

// DeleteAllKeys deletes all keys from the in-memory cache.
func (imc *InMemoryCache) DeleteAllKeys() error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	imc.cache = make(map[string]*list.Element)
	imc.evictList.Init()
	return nil
}

// evictOldest evicts the least recently used entry from the cache.
func (imc *InMemoryCache) evictOldest() {
	if imc.evictList.Len() > 0 {
		oldest := imc.evictList.Back()
		if oldest != nil {
			imc.removeElement(oldest)
		}
	}
}

// removeElement removes an element from the cache and evictList.
func (imc *InMemoryCache) removeElement(e *list.Element) {
	imc.evictList.Remove(e)
	delete(imc.cache, e.Value.(*entry).key)
}

// startEvictionTicker starts a ticker to periodically check and evict expired entries.
func (imc *InMemoryCache) startEvictionTicker() {
	ticker := time.NewTicker(5 * time.Minute) // Adjust tick interval as needed
	defer ticker.Stop()

	for range ticker.C {
		imc.evictExpiredEntries()
	}
}

// evictExpiredEntries evicts all expired entries from the cache.
func (imc *InMemoryCache) evictExpiredEntries() {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	for _, elem := range imc.cache { //error key unused oocured here
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			imc.removeElement(elem)
		}
	}
}

/*
import (
	"container/list"
	"sync"
	"time"
	"fmt"
)

// entry represents a cache entry stored in the InMemoryCache.
type entry struct {
	key        string
	value      interface{}
	expiration time.Time
}

// InMemoryCache represents an in-memory cache with LRU eviction and expiration handling.
type InMemoryCache struct {
	capacity  int
	cache     map[string]*list.Element
	evictList *list.List
	mu        sync.Mutex
}

// NewInMemoryCache creates a new instance of InMemoryCache with the specified capacity.
func NewInMemoryCache(capacity int) *InMemoryCache {
	return &InMemoryCache{
		capacity:  capacity,
		cache:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Set adds a new key-value pair to the in-memory cache with optional expiration.
func (imc *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	// Check if the key exists in the cache
	if elem, ok := imc.cache[key]; ok {
		// If key exists, update the value, expiration, and move the element to front (MRU position)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).expiration = time.Now().Add(expiration)
		imc.evictList.MoveToFront(elem)
	} else {
		// If key does not exist, add the new entry to the cache
		if len(imc.cache) >= imc.capacity {
			// Evict least recently used element if cache is full
			imc.evictOldest()
		}

		// Add new entry to the cache with expiration time
		newElem := imc.evictList.PushFront(&entry{
			key:        key,
			value:      value,
			expiration: time.Now().Add(expiration),
		})
		imc.cache[key] = newElem
	}

	// Start a goroutine to periodically check and evict expired entries
	go imc.startEvictionTicker()

	return nil
}

// Get retrieves a value from the in-memory cache based on the key.
func (imc *InMemoryCache) Get(key string) (interface{}, error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		// Check if the entry is expired
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			// If expired, delete the entry and return not found
			imc.removeElement(elem)
			return nil, fmt.Errorf("key expired")
		}

		// Move the accessed element to front (MRU position)
		imc.evictList.MoveToFront(elem)
		return elem.Value.(*entry).value, nil
	}
	return nil, fmt.Errorf("key not found")
}

// Delete removes a key-value pair from the in-memory cache.
func (imc *InMemoryCache) Delete(key string) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		imc.removeElement(elem)
		return nil
	}
	return fmt.Errorf("key not found")
}

// GetAllKeys retrieves all keys from the in-memory cache.
func (imc *InMemoryCache) GetAllKeys() ([]string, error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	keys := make([]string, 0, len(imc.cache))
	for key := range imc.cache {
		keys = append(keys, key)
	}
	return keys, nil
}

// DeleteAllKeys deletes all keys from the in-memory cache.
func (imc *InMemoryCache) DeleteAllKeys() error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	imc.cache = make(map[string]*list.Element)
	imc.evictList.Init()
	return nil
}

// evictOldest evicts the least recently used entry from the cache.
func (imc *InMemoryCache) evictOldest() {
	if imc.evictList.Len() > 0 {
		oldest := imc.evictList.Back()
		if oldest != nil {
			imc.removeElement(oldest)
		}
	}
}

// removeElement removes an element from the cache and evictList.
func (imc *InMemoryCache) removeElement(e *list.Element) {
	imc.evictList.Remove(e)
	delete(imc.cache, e.Value.(*entry).key)
}

// startEvictionTicker starts a ticker to periodically check and evict expired entries.
func (imc *InMemoryCache) startEvictionTicker() {
	ticker := time.NewTicker(5 * time.Minute) // Adjust tick interval as needed
	defer ticker.Stop()

	for range ticker.C {
		imc.evictExpiredEntries()
	}
}

// evictExpiredEntries evicts all expired entries from the cache.
func (imc *InMemoryCache) evictExpiredEntries() {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	for _, elem := range imc.cache {
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			imc.removeElement(elem)
		}
	}
}
*/

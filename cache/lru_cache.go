package cache

import (
	"container/list"
	"sync"
	"time"
)

// LRUCacheEntry represents an entry in the LRUCache.
type LRUCacheEntry struct {
	key        string
	value      interface{}
	expiration time.Time
}

// LRUCache represents a Least Recently Used (LRU) cache.
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	eviction *list.List
	mutex    sync.Mutex
}

// NewLRUCache creates a new instance of LRUCache with the specified capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		eviction: list.New(),
	}
}

// Set adds a new key-value pair to the LRUCache with an expiration duration.
func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if key already exists
	if elem, ok := c.cache[key]; ok {
		c.eviction.MoveToFront(elem)
		entry := elem.Value.(*LRUCacheEntry)
		entry.value = value
		entry.expiration = time.Now().Add(expiration)
		return
	}

	// Evict least recently used if capacity exceeded
	if len(c.cache) >= c.capacity {
		c.evict()
	}

	// Add new entry
	entry := &LRUCacheEntry{
		key:        key,
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	elem := c.eviction.PushFront(entry)
	c.cache[key] = elem
}

// Get retrieves a value from the LRUCache by key.
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		entry := elem.Value.(*LRUCacheEntry)
		// Check expiration
		if entry.expiration.After(time.Now()) {
			c.eviction.MoveToFront(elem)
			return entry.value, true
		}
		// Expired entry, evict it
		c.eviction.Remove(elem)
		delete(c.cache, key)
	}

	return nil, false
}

// Remove removes a key-value pair from the LRUCache.
func (c *LRUCache) Remove(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.eviction.Remove(elem)
		delete(c.cache, key)
	}
}

// GetAll retrieves all key-value pairs from the LRUCache.
func (c *LRUCache) GetAll() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := make(map[string]interface{})
	for key, elem := range c.cache {
		entry := elem.Value.(*LRUCacheEntry)
		if entry.expiration.After(time.Now()) {
			result[key] = entry.value
		} else {
			// Expired entry, evict it
			c.eviction.Remove(elem)
			delete(c.cache, key)
		}
	}
	return result
}

// Clear clears all key-value pairs from the LRUCache.
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*list.Element)
	c.eviction = list.New()
}

// evict removes the least recently used entry from the LRUCache.
func (c *LRUCache) evict() {
	if len(c.cache) == 0 {
		return
	}

	// Remove least recently used entry
	last := c.eviction.Back()
	if last != nil {
		entry := last.Value.(*LRUCacheEntry)
		delete(c.cache, entry.key)
		c.eviction.Remove(last)
	}
}

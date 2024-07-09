package cache

import (
	"container/list"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// LRUCacheEntry represents an entry in the LRUCache.
type LRUCacheEntry struct {
	key        string
	value      interface{}
	expiration time.Time
}

// LRUCache represents a Least Recently Used (LRU) cache with Memcached integration.
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	eviction *list.List
	mutex    sync.Mutex
	memcache *memcache.Client
}

// NewLRUCacheWithMemcached creates a new instance of LRUCache with Memcached integration.
func NewLRUCacheWithMemcached(capacity int, memcachedAddr string) *LRUCache {
	client := memcache.New(memcachedAddr)
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		eviction: list.New(),
		memcache: client,
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

		// Update Memcached value
		if err := c.setMemcachedValue(key, value, expiration); err != nil {
			// Handle error
			fmt.Printf("Failed to update Memcached: %v\n", err)
		}
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

	// Set Memcached value
	if err := c.setMemcachedValue(key, value, expiration); err != nil {
		// Handle error
		fmt.Printf("Failed to set Memcached: %v\n", err)
	}
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

			// Retrieve from Memcached
			if value, err := c.getMemcachedValue(key); err == nil {
				return value, true
			}
		} else {
			// Expired entry, evict it
			c.eviction.Remove(elem)
			delete(c.cache, key)

			// Delete from Memcached
			c.deleteMemcachedValue(key)
		}
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

		// Delete from Memcached
		c.deleteMemcachedValue(key)
	}
}

// Clear clears all key-value pairs from the LRUCache.
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key := range c.cache {
		c.deleteMemcachedValue(key)
	}
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

		// Delete from Memcached
		c.deleteMemcachedValue(entry.key)
	}
}

// setMemcachedValue sets a value in Memcached with the given expiration.
func (c *LRUCache) setMemcachedValue(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        key,
		Value:      jsonValue,
		Expiration: int32(expiration.Seconds()),
	}
	return c.memcache.Set(item)
}

// getMemcachedValue retrieves a value from Memcached by key.
func (c *LRUCache) getMemcachedValue(key string) (interface{}, error) {
	item, err := c.memcache.Get(key)
	if err != nil {
		return nil, err
	}

	var value interface{}
	if err := json.Unmarshal(item.Value, &value); err != nil {
		return nil, err
	}
	return value, nil
}

// deleteMemcachedValue deletes a value from Memcached by key.
func (c *LRUCache) deleteMemcachedValue(key string) error {
	return c.memcache.Delete(key)
}

// GetAll returns all key-value pairs currently in the LRUCache.
func (c *LRUCache) GetAll() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	allValues := make(map[string]interface{})
	for key, elem := range c.cache {
		entry := elem.Value.(*LRUCacheEntry)
		if entry.expiration.After(time.Now()) {
			allValues[key] = entry.value
		}
	}
	return allValues
}

// DeleteAll removes all key-value pairs from the LRUCache.
func (c *LRUCache) DeleteAll() {
	c.Clear()
}

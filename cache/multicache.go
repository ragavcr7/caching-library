package cache

import (
	"sync"
	"time"
)

// Multicache represents a multicache system with both InMemoryCache and LRUCacheWithMemcached.
type Multicache struct {
	inMemoryCache   *InMemoryCache
	lruCacheWithMem *LRUCache
	mutex           sync.Mutex
}

// NewMulticache creates a new instance of Multicache with specified capacities and Memcached address.
func NewMulticache(inMemoryCapacity, lruCapacity int, memcachedAddr string) *Multicache {
	return &Multicache{
		inMemoryCache:   NewInMemoryCache(inMemoryCapacity),
		lruCacheWithMem: NewLRUCacheWithMemcached(lruCapacity, memcachedAddr),
	}
}

// Set adds a key-value pair to both InMemoryCache and LRUCacheWithMemcached.
func (mc *Multicache) Set(key string, value interface{}, expiration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.inMemoryCache.Set(key, value, expiration)
	mc.lruCacheWithMem.Set(key, value, expiration)
}

// Get retrieves a value from either InMemoryCache or LRUCacheWithMemcached based on availability.
func (mc *Multicache) Get(key string) (interface{}, bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// Try to get from InMemoryCache first
	if value, ok := mc.inMemoryCache.Get(key); ok {
		return value, true
	}

	// If not found in InMemoryCache, try LRUCacheWithMemcached
	return mc.lruCacheWithMem.Get(key)
}

// Remove removes a key-value pair from both InMemoryCache and LRUCacheWithMemcached.
func (mc *Multicache) Remove(key string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.inMemoryCache.Delete(key)
	mc.lruCacheWithMem.Remove(key)
}

// GetAllKeys retrieves all keys from both InMemoryCache and LRUCacheWithMemcached.
func (mc *Multicache) GetAllKeys() []string {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// Get keys from InMemoryCache
	keys := mc.inMemoryCache.GetAllKeys()

	// Append keys from LRUCacheWithMemcached (assuming no duplicates for simplicity)
	lruKeys := mc.lruCacheWithMem.GetAll()
	for _, key := range lruKeys {
		found := false
		for _, k := range keys {
			if k == key {
				found = true
				break
			}
		}
		if !found {
			keys = append(keys, key.(string))
		}
	}

	return keys
}

// DeleteAllKeys deletes all keys from both InMemoryCache and LRUCacheWithMemcached.
func (mc *Multicache) DeleteAllKeys() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.inMemoryCache.DeleteAllKeys()
	mc.lruCacheWithMem.DeleteAll()
}

// Clear clears all key-value pairs from both InMemoryCache and LRUCacheWithMemcached.
func (mc *Multicache) Clear() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.inMemoryCache.DeleteAllKeys()
	mc.lruCacheWithMem.Clear()
}

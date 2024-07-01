package cache

/*
import (
	"time"
)

// MultiCache implements Cache interface with both in-memory and memcached caches.
type MultiCache struct {
	inMemoryCache  *InMemoryCache
	memcachedCache *MemcachedCache
}

// NewMultiCache creates a new instance of MultiCache.
func NewMultiCache(memcachedServer string, capacity int) *MultiCache {
	return &MultiCache{
		inMemoryCache:  NewInMemoryCache(capacity),
		memcachedCache: NewMemcachedCache(memcachedServer),
	}
}

// Set adds a new key-value pair to both in-memory and memcached caches.
func (mc *MultiCache) Set(key string, value interface{}, expiration time.Duration) error{
	// Set in both caches
	if err := mc.inMemoryCache.Set(key, value, expiration); err != nil {
		return err
	}
	return mc.memcachedCache.Set(key, value, expiration)
}

// Get retrieves a value from both in-memory and memcached caches.
func (mc *MultiCache) Get(key string) (string, error) {
	// Try to get from in-memory cache first
	value, err := mc.inMemoryCache.Get(key)
	if !err {
		return value.(string), nil
	}
	// If not found in in-memory cache, try memcached
	return mc.memcachedCache.Get(key)
}

// Delete removes a key-value pair from both in-memory and memcached caches.
func (mc *MultiCache) Delete(key string) error {
	// Delete from both caches
	if err := mc.inMemoryCache.Delete(key); err != nil {
		return err
	}
	return mc.memcachedCache.Delete(key)
}

// GetAllKeys retrieves all keys from both in-memory and memcached caches.
func (mc *MultiCache) GetAllKeys() ([]string, error) {
	// Combine keys from both caches
	inMemoryKeys := mc.inMemoryCache.GetAllKeys()
	memcachedKeys, err := mc.memcachedCache.GetAllKeys()
	if err != nil {
		return nil, err
	}
	keys := append(inMemoryKeys, memcachedKeys...)
	return keys, nil
}

// DeleteAllKeys deletes all keys from both in-memory and memcached caches.
func (mc *MultiCache) DeleteAllKeys() error {
	// Delete from both caches
	if err := mc.inMemoryCache.DeleteAllKeys(); err != nil {
		return err
	}
	return mc.memcachedCache.DeleteAllKeys()
}
*/

import (
	"log"
	"time"
)

// MultiCache implements Cache interface with both in-memory and memcached caches.
type MultiCache struct {
	inMemoryCache  *InMemoryCache
	memcachedCache *MemcachedCache
}

// NewMultiCache creates a new instance of MultiCache.
func NewMultiCache(memcachedServer string, capacity int) *MultiCache {
	return &MultiCache{
		inMemoryCache:  NewInMemoryCache(capacity),
		memcachedCache: NewMemcachedCache(memcachedServer),
	}
}

// Set adds a new key-value pair to both in-memory and memcached caches.
func (mc *MultiCache) Set(key string, value interface{}, expiration time.Duration) error {
	// Set in both caches
	err := mc.inMemoryCache.Set(key, value, expiration)
	if err != nil {
		log.Printf("Error setting value in in-memory cache: %v", err)
		return err
	}
	log.Printf("Successfully set key: %s in in-memory cache", key)

	err = mc.memcachedCache.Set(key, value, expiration)
	if err != nil {
		log.Printf("Error setting value in memcached: %v", err)
		return err
	}
	log.Printf("Successfully set key: %s in memcached", key)

	return nil
}

// Get retrieves a value from both in-memory and memcached caches.
func (mc *MultiCache) Get(key string) (interface{}, error) {
	// Try to get from in-memory cache first
	value, found := mc.inMemoryCache.Get(key)
	if found {
		log.Printf("Successfully retrieved key: %s from in-memory cache", key)
		return value, nil
	}
	log.Printf("Key: %s not found in in-memory cache, trying memcached", key)

	// If not found in in-memory cache, try memcached
	value, err := mc.memcachedCache.Get(key)
	if err != nil {
		log.Printf("Error getting value from memcached for key: %s, %v", key, err)
		return nil, err
	}
	log.Printf("Successfully retrieved key: %s from memcached", key)

	// Cache the retrieved value back to in-memory cache
	err = mc.inMemoryCache.Set(key, value, 10*time.Minute)
	if err != nil {
		log.Printf("Error setting value in in-memory cache after retrieving from memcached: %v", err)
	}

	return value, nil
}

// Delete removes a key-value pair from both in-memory and memcached caches.
func (mc *MultiCache) Delete(key string) error {
	// Delete from both caches
	err := mc.inMemoryCache.Delete(key)
	if err != nil {
		log.Printf("Error deleting key: %s from in-memory cache: %v", key, err)
		return err
	}
	log.Printf("Successfully deleted key: %s from in-memory cache", key)

	err = mc.memcachedCache.Delete(key)
	if err != nil {
		log.Printf("Error deleting key: %s from memcached: %v", key, err)
		return err
	}
	log.Printf("Successfully deleted key: %s from memcached", key)

	return nil
}

// GetAllKeys retrieves all keys from both in-memory and memcached caches.
func (mc *MultiCache) GetAllKeys() ([]string, error) {
	// Combine keys from both caches
	inMemoryKeys := mc.inMemoryCache.GetAllKeys()
	log.Printf("Successfully retrieved keys from in-memory cache: %v", inMemoryKeys)

	memcachedKeys, err := mc.memcachedCache.GetAllKeys()
	if err != nil {
		log.Printf("Error getting all keys from memcached: %v", err)
		return nil, err
	}
	log.Printf("Successfully retrieved keys from memcached: %v", memcachedKeys)

	keys := append(inMemoryKeys, memcachedKeys...)
	return keys, nil
}

// DeleteAllKeys deletes all keys from both in-memory and memcached caches.
func (mc *MultiCache) DeleteAllKeys() error {
	// Delete from both caches
	err := mc.inMemoryCache.DeleteAllKeys()
	if err != nil {
		log.Printf("Error deleting all keys from in-memory cache: %v", err)
		return err
	}
	log.Printf("Successfully deleted all keys from in-memory cache")

	err = mc.memcachedCache.DeleteAllKeys()
	if err != nil {
		log.Printf("Error deleting all keys from memcached: %v", err)
		return err
	}
	log.Printf("Successfully deleted all keys from memcached")

	return nil
}

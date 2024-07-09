// memcached.go
package cache

/* this package only implements memcache not lru with memcache for that refer lru_cache.go file
import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedCache represents a Memcached cache client.
type MemcachedCache struct {
	client  *memcache.Client //ptr to memcached client
	mu      sync.RWMutex     //rw muted to handle concurrent access to cache
	keyList map[string]bool  //new
}

// CONSTRUCTOR -- NewMemcachedCache creates a new instance of MemcachedCache.
func NewMemcachedCache(serverAddr string) *MemcachedCache {
	client := memcache.New(serverAddr)
	return &MemcachedCache{
		client:  client,
		keyList: make(map[string]bool), //new
	}
}

// Set adds a new key-value pair to the Memcached cache.
func (mc *MemcachedCache) Set(key string, value interface{}, expiration time.Duration) error {
	mc.mu.Lock()         //for thread safety purpose
	defer mc.mu.Unlock() //
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value to JSON: %w", err)
	}

	item := &memcache.Item{
		Key:        key,
		Value:      jsonValue,
		Expiration: int32(expiration.Seconds()),
	}
	if err := mc.client.Set(item); err != nil {
		return fmt.Errorf("failed to set key %s in Memcached: %w", key, err)
	}
	mc.keyList[key] = true // new
	return nil
}

// Get retrieves a value from the Memcached cache based on the key.
func (mc *MemcachedCache) Get(key string) (string, error) {
	mc.mu.Lock()         //lock the mutex for duration of this method
	defer mc.mu.Unlock() //
	item, err := mc.client.Get(key)
	if err == memcache.ErrCacheMiss { //not found
		return "", fmt.Errorf("key %s not found in Memcached", key)
	} else if err != nil { //other rtrieval errors
		return "", fmt.Errorf("failed to get key %s from Memcached: %w", key, err)
	}
	//return string(item.Value), nil //item.value is genraly a json so here we are converting it to string so we can able to return it
	var value string
	if err := json.Unmarshal(item.Value, &value); err != nil {
		return "", fmt.Errorf("failed to unmarshal value: %w", err)
	}
	return value, nil
}

// Delete removes a key-value pair from the Memcached cache.
func (mc *MemcachedCache) Delete(key string) error {
	mc.mu.Lock()         //
	defer mc.mu.Unlock() //
	if err := mc.client.Delete(key); err != nil {
		return fmt.Errorf("failed to delete key %s from Memcached: %w", key, err)
	}
	delete(mc.keyList, key) //new
	return nil
}

// GetAllKeys retrieves all keys from the Memcached cache.
func (mc *MemcachedCache) GetAllKeys() ([]string, error) {
	mc.mu.Lock()                               //
	defer mc.mu.Unlock()                       //
	keys := make([]string, 0, len(mc.keyList)) // new
	for key := range mc.keyList {
		keys = append(keys, key)
	}
	// Memcached does not support listing all keys, so this operation is not possible with basic memcache clients.
	return keys, nil
}

// DeleteAllKeys deletes all keys from the Memcached cache.
func (mc *MemcachedCache) DeleteAllKeys() error {
	mc.mu.Lock()         //
	defer mc.mu.Unlock() //
	// Memcached does not support bulk deletion, so this operation is not possible with basic memcache clients.
	//return fmt.Errorf("deleting all keys is not supported in Memcached")
	mc.keyList = make(map[string]bool) // Clear keyList --- new
	return nil                         // ---new
}

// MemcachedCache is a cache implementation using Memcached.
*/

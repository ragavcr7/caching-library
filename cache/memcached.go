/* implementation of memcached cache for manual inputs
package cache

import (

	"fmt"
	"time"








	"github.com/bradfitz/gomemcache/memcache"

)

type MemcachedCache struct { //it represents a client

		client *memcache.Client
	}

func NewMemcachedCache(serverAddr string) *MemcachedCache { //using pointer here to avoid the copy

		client := memcache.New(serverAddr)
		return &MemcachedCache{
			client: client,
		}
	}

// this method is used to add a new key value pair to cache and using the capital S for set to make it as public
func (mc *MemcachedCache) Set(key string, value interface{}, expiration time.Duration) error { //we all know key will always be string and value can be of any comparable type so we sets interface for value

		item := &memcache.Item{
			Key:        key,
			Value:      []byte(fmt.Sprintf("%v", value)), //byte dt since all the information from client is gonna be in byte format
			Expiration: int32(expiration.Seconds()),
		}
		return mc.client.Set(item) //manages the connection to memcached server so we can perform set,get,delete operations
	}

// this method is used to fetch the value based on the passed key from the memcached cache

	func (mc *MemcachedCache) Get(key string) (interface{}, error) {
		item, err := mc.client.Get(key)
		if err == memcache.ErrCacheMiss {
			return nil, fmt.Errorf("%v not found in the Memcached cache", key)
		} else if err != nil {
			return nil, fmt.Errorf("failed to get key %v from Memcached: %v", key, err)
		}
		return string(item.Value), nil
	}

	func (mc *MemcachedCache) Delete(key string) error {
		return mc.client.Delete(key)
	}
*/
/*
//working one
package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedCache struct {
	client *memcache.Client
}

func NewMemcachedCache(serverAddr string) *MemcachedCache {
	client := memcache.New(serverAddr)
	return &MemcachedCache{
		client: client,
	}
}

func (mc *MemcachedCache) Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	item := &memcache.Item{
		Key:        key,
		Value:      jsonValue,
		Expiration: int32(expiration.Seconds()),
	}
	return mc.client.Set(item)
}

func (mc *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := mc.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, fmt.Errorf("%v not found in the Memcached cache", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key %v from Memcached: %v", key, err)
	}
	return item.Value, nil
}

func (mc *MemcachedCache) Delete(key string) error {
	return mc.client.Delete(key)
}
*/
// memcached.go
package cache

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"sync"
	"time"
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

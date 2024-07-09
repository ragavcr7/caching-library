package tests

import (
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

// Adjusted integration test with error handling and verification
func TestCacheIntegration(t *testing.T) {
	// Initialize caches
	inMemoryCapacity := 5
	lruCapacity := 3
	memcachedAddr := "localhost:11211"

	inMemoryCache := cache.NewInMemoryCache(inMemoryCapacity)
	lruCache := cache.NewLRUCacheWithMemcached(lruCapacity, memcachedAddr)
	multicache := cache.NewMulticache(inMemoryCapacity, lruCapacity, memcachedAddr)

	defer func() {
		inMemoryCache.DeleteAllKeys()
		lruCache.Clear()
		multicache.Clear()
	}()

	// Test scenario 1: Set in InMemoryCache and get from Multicache
	key1 := "key1"
	value1 := "value1"
	expiration1 := 1 * time.Minute

	inMemoryCache.Set(key1, value1, expiration1)
	time.Sleep(100 * time.Millisecond) // Allow caches to sync

	if val, found := multicache.Get(key1); !found || val.(string) != value1 {
		t.Errorf("Expected value %s for key %s from Multicache, got %v", value1, key1, val)
	}

	// Test scenario 2: Set in LRUCache and get from Multicache
	key2 := "key2"
	value2 := 12345
	expiration2 := 5 * time.Second

	lruCache.Set(key2, value2, expiration2)
	time.Sleep(100 * time.Millisecond) // Allow caches to sync

	if val, found := multicache.Get(key2); !found || val.(int) != value2 {
		t.Errorf("Expected value %d for key %s from Multicache, got %v", value2, key2, val)
	}

	// Test scenario 3: Set in Multicache and verify individual caches
	key3 := "key3"
	value3 := true
	expiration3 := 10 * time.Minute

	multicache.Set(key3, value3, expiration3)
	time.Sleep(100 * time.Millisecond) // Allow caches to sync

	// Verify in InMemoryCache
	if val, found := inMemoryCache.Get(key3); !found || val.(bool) != value3 {
		t.Errorf("Expected value %v for key %s in InMemoryCache, got %v", value3, key3, val)
	}

	// Verify in LRUCache
	if val, found := lruCache.Get(key3); !found || val.(bool) != value3 {
		t.Errorf("Expected value %v for key %s in LRUCache, got %v", value3, key3, val)
	}

	// Test scenario 4: Delete from Multicache and verify
	multicache.Remove(key1)
	time.Sleep(100 * time.Millisecond) // Allow caches to sync

	if _, found := multicache.Get(key1); found {
		t.Errorf("Expected key %s to be deleted from Multicache, but it still exists", key1)
	}

	// Verify in InMemoryCache
	if _, found := inMemoryCache.Get(key1); found {
		t.Errorf("Expected key %s to be deleted from InMemoryCache, but it still exists", key1)
	}

	// Verify in LRUCache
	if _, found := lruCache.Get(key1); found {
		t.Errorf("Expected key %s to be deleted from LRUCache, but it still exists", key1)
	}

	// Test scenario 5: Delete all keys from Multicache and verify
	multicache.DeleteAllKeys()
	time.Sleep(100 * time.Millisecond) // Allow caches to sync

	if keys := multicache.GetAllKeys(); len(keys) != 0 {
		t.Errorf("Expected 0 keys in Multicache after DeleteAllKeys, got %d", len(keys))
	}

	// Verify in InMemoryCache
	if keys := inMemoryCache.GetAllKeys(); len(keys) != 0 {
		t.Errorf("Expected 0 keys in InMemoryCache after DeleteAllKeys, got %d", len(keys))
	}

	// Verify in LRUCache
	if keys := lruCache.GetAll(); len(keys) != 0 {
		t.Errorf("Expected 0 keys in LRUCache after DeleteAllKeys, got %d", len(keys))
	}
}

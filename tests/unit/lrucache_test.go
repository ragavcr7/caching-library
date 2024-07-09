package tests

import (
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func TestLRUCacheOperations(t *testing.T) {
	// Initialize LRUCache with Memcached
	cache := cache.NewLRUCacheWithMemcached(3, "127.0.0.1:11211")

	// Test Set operation
	cache.Set("key1", "value1", 5*time.Second)
	cache.Set("key2", 12345, 10*time.Second)
	cache.Set("key3", true, 15*time.Second)

	// Test Get operation
	testGetOperation(cache, "key1", "value1", t)
	testGetOperation(cache, "key2", 12345, t)
	testGetOperation(cache, "key3", true, t)

	// Test GetAll operation
	allValues := cache.GetAll()
	if len(allValues) != 3 {
		t.Errorf("Expected 3 items in cache, got %d", len(allValues))
	}
	checkCacheContent(allValues, "key1", "value1", t)
	checkCacheContent(allValues, "key2", 12345, t)
	checkCacheContent(allValues, "key3", true, t)

	// Test DeleteAll operation
	cache.DeleteAll()
	allValues = cache.GetAll()
	if len(allValues) != 0 {
		t.Errorf("Expected cache to be empty after DeleteAll, got %d items", len(allValues))
	}

	// Test Get after DeleteAll
	_, found := cache.Get("key1")
	if found {
		t.Errorf("Expected key1 to not exist after DeleteAll, but it was found")
	}
}

func testGetOperation(cache *cache.LRUCache, key string, expectedValue interface{}, t *testing.T) {
	value, found := cache.Get(key)
	if !found {
		t.Errorf("Expected key '%s' to exist in cache, but it was not found", key)
	} else {
		if value != expectedValue {
			t.Errorf("Expected key '%s' value to be %v, got %v", key, expectedValue, value)
		}
	}
}

func checkCacheContent(allValues map[string]interface{}, key string, expectedValue interface{}, t *testing.T) {
	if value, ok := allValues[key]; !ok {
		t.Errorf("Expected key '%s' to be in cache, but it was not found", key)
	} else {
		if value != expectedValue {
			t.Errorf("Expected key '%s' value to be %v, got %v", key, expectedValue, value)
		}
	}
}

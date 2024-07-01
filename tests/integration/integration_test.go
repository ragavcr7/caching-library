package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
	"github.com/ragavcr7/caching-library/cache_interface"
)

const (
	memcachedServerAddr = "localhost:11211"
	cacheCapacity       = 10
)

// setupInMemoryCache sets up a new instance of InMemoryCache for testing.
func setupInMemoryCache() cache_interface.Cache {
	return cache.NewInMemoryCache(cacheCapacity)
}

// setupMemcachedCache sets up a new instance of MemcachedCache for testing.
func setupMemcachedCache() cache_interface.Cache {
	return cache.NewMemcachedCache(memcachedServerAddr)
}

// TestInMemoryCache tests the InMemoryCache implementation.
func TestInMemoryCache(t *testing.T) {
	cacheInterface := setupInMemoryCache()
	runIntegrationTests(t, cacheInterface)
}

// TestMemcachedCache tests the MemcachedCache implementation.
func TestMemcachedCache(t *testing.T) {
	cacheInterface := setupMemcachedCache()
	runIntegrationTests(t, cacheInterface)
}

// runIntegrationTests runs a series of integration tests on the provided cacheInterface.
func runIntegrationTests(t *testing.T, cacheInterface cache_interface.Cache) {
	t.Helper()

	t.Run("SetAndGet", func(t *testing.T) {
		testSetAndGet(t, cacheInterface)
	})

	t.Run("Delete", func(t *testing.T) {
		testDelete(t, cacheInterface)
	})

	t.Run("GetAllKeys", func(t *testing.T) {
		testGetAllKeys(t, cacheInterface)
	})

	t.Run("DeleteAllKeys", func(t *testing.T) {
		testDeleteAllKeys(t, cacheInterface)
	})
}

func testSetAndGet(t *testing.T, cacheInterface cache_interface.Cache) {
	err := cacheInterface.Set("key1", "value1", 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	val, found := cacheInterface.Get("key1")
	if found != nil {
		t.Error("Expected key1 to be found in cache, but it wasn't")
	}
	if val != "value1" {
		t.Errorf("Expected value: value1, got: %v", val)
	}

	// Test expiration
	time.Sleep(2 * time.Minute)
	_, found = cacheInterface.Get("key1")
	if found != nil {
		t.Error("Expected key1 to be expired and not found, but it was found")
	}
}

func testDelete(t *testing.T, cacheInterface cache_interface.Cache) {
	err := cacheInterface.Set("key2", "value2", 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	err = cacheInterface.Delete("key2")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	_, found := cacheInterface.Get("key2")
	if found != nil {
		t.Error("Expected key2 to be deleted and not found, but it was found")
	}
}

func testGetAllKeys(t *testing.T, cacheInterface cache_interface.Cache) {
	keys := []string{"key3", "key4", "key5"}
	for _, key := range keys {
		err := cacheInterface.Set(key, fmt.Sprintf("value_%s", key), 1*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set value for key %s: %v", key, err)
		}
	}

	allKeys, err := cacheInterface.GetAllKeys()
	if err != nil {
		t.Fatalf("Failed to get all keys: %v", err)
	}

	keySet := make(map[string]bool)
	for _, key := range allKeys {
		keySet[key] = true
	}

	for _, key := range keys {
		if !keySet[key] {
			t.Errorf("Expected key %s in GetAllKeys result, but it was missing", key)
		}
	}
}

func testDeleteAllKeys(t *testing.T, cacheInterface cache_interface.Cache) {
	keys := []string{"key6", "key7", "key8"}
	for _, key := range keys {
		err := cacheInterface.Set(key, fmt.Sprintf("value_%s", key), 1*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set value for key %s: %v", key, err)
		}
	}

	err := cacheInterface.DeleteAllKeys()
	if err != nil {
		t.Fatalf("Failed to delete all keys: %v", err)
	}

	allKeys, err := cacheInterface.GetAllKeys()
	if err != nil {
		t.Fatalf("Failed to get all keys after DeleteAllKeys: %v", err)
	}

	if len(allKeys) > 0 {
		t.Errorf("Expected all keys to be deleted, but found %d keys remaining", len(allKeys))
	}
}

// multicache_test.go

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func TestMulticacheSetAndGet(t *testing.T) {
	mc := cache.NewMulticache(10, 10, "localhost:11211")
	defer mc.Clear()

	key := "testKey"
	value := "testValue"
	expiration := 5 * time.Minute

	mc.Set(key, value, expiration)

	if val, found := mc.Get(key); !found || val.(string) != value {
		t.Errorf("Expected value %s for key %s, got %v", value, key, val)
	}
}

func TestMulticacheDelete(t *testing.T) {
	mc := cache.NewMulticache(10, 10, "localhost:11211")
	defer mc.Clear()

	key := "testKey"
	value := "testValue"
	expiration := 5 * time.Minute

	mc.Set(key, value, expiration)
	mc.Remove(key)

	if _, found := mc.Get(key); found {
		t.Errorf("Expected key %s to be deleted, but it still exists", key)
	}
}

func TestMulticacheGetAllKeys(t *testing.T) {
	mc := cache.NewMulticache(10, 10, "localhost:11211")
	defer mc.Clear()

	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := 5 * time.Minute
		mc.Set(key, value, expiration)
	}

	keys := mc.GetAllKeys()
	expectedKeys := []string{"key0", "key1", "key2"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}

	for _, key := range expectedKeys {
		found := false
		for _, k := range keys {
			if k == key {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected key %s not found in GetAllKeys result", key)
		}
	}
}

func TestMulticacheDeleteAllKeys(t *testing.T) {
	mc := cache.NewMulticache(10, 10, "localhost:11211")
	defer mc.Clear()

	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := 5 * time.Minute
		mc.Set(key, value, expiration)
	}

	mc.DeleteAllKeys()
	keys := mc.GetAllKeys()

	if len(keys) != 0 {
		t.Errorf("Expected 0 keys after DeleteAllKeys, got %d", len(keys))
	}
}

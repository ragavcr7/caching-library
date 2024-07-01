package benchmarks

import (
	"fmt"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func BenchmarkMultiCache_Set(b *testing.B) {
	mc := cache.NewMultiCache("localhost:11211", 100)

	for n := 0; n < b.N; n++ {
		key := fmt.Sprintf("test_key_%d", n)
		value := fmt.Sprintf("test_value_%d", n)
		err := mc.Set(key, value, 5*time.Minute)
		if err != nil {
			b.Errorf("Error setting key-value pair: %v", err)
		}
	}
}

func BenchmarkMultiCache_Get(b *testing.B) {
	mc := cache.NewMultiCache("localhost:11211", 100)
	key := "test_key"
	value := "test_value"
	mc.Set(key, value, 5*time.Minute)

	for n := 0; n < b.N; n++ {
		_, err := mc.Get(key)
		if err != nil {
			b.Errorf("Error getting value for key %s: %v", key, err)
		}
	}
}

func BenchmarkMultiCache_GetAllKeys(b *testing.B) {
	mc := cache.NewMultiCache("localhost:11211", 100)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("test_key_%d", i)
		value := fmt.Sprintf("test_value_%d", i)
		mc.Set(key, value, 5*time.Minute)
	}

	for n := 0; n < b.N; n++ {
		_, err := mc.GetAllKeys()
		if err != nil {
			b.Errorf("Error getting all keys: %v", err)
		}
	}
}

func BenchmarkMultiCache_Delete(b *testing.B) {
	mc := cache.NewMultiCache("localhost:11211", 100)
	key := "test_key"
	value := "test_value"
	mc.Set(key, value, 5*time.Minute)

	for n := 0; n < b.N; n++ {
		err := mc.Delete(key)
		if err != nil {
			b.Errorf("Error deleting key %s: %v", key, err)
		}
	}
}

func BenchmarkMultiCache_DeleteAllKeys(b *testing.B) {
	mc := cache.NewMultiCache("localhost:11211", 100)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("test_key_%d", i)
		value := fmt.Sprintf("test_value_%d", i)
		mc.Set(key, value, 5*time.Minute)
	}

	for n := 0; n < b.N; n++ {
		err := mc.DeleteAllKeys()
		if err != nil {
			b.Errorf("Error deleting all keys: %v", err)
		}
	}
}

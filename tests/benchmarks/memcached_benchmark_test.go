// benchmark.go
package benchmarks

import (
	"fmt"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

const (
	memcachedServerAddr = "localhost:11211"
	cacheKeyPrefix      = "benchmark_key_"
	cacheValue          = "test_value"
	cacheExpiration     = 60 * time.Second
	cacheSize           = 1000
)

func BenchmarkMemcachedCacheSet(b *testing.B) {
	mc := cache.NewMemcachedCache(memcachedServerAddr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%s%d", cacheKeyPrefix, i)
		if err := mc.Set(key, cacheValue, cacheExpiration); err != nil {
			b.Fatalf("failed to set key %s: %v", key, err)
		}
	}
}

func BenchmarkMemcachedCacheGet(b *testing.B) {
	mc := cache.NewMemcachedCache(memcachedServerAddr)
	for i := 0; i < cacheSize; i++ {
		key := fmt.Sprintf("%s%d", cacheKeyPrefix, i)
		if err := mc.Set(key, cacheValue, cacheExpiration); err != nil {
			b.Fatalf("failed to set key %s: %v", key, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%s%d", cacheKeyPrefix, i%cacheSize)
		if _, err := mc.Get(key); err != nil {
			b.Fatalf("failed to get key %s: %v", key, err)
		}
	}
}

func BenchmarkMemcachedCacheDelete(b *testing.B) {
	mc := cache.NewMemcachedCache(memcachedServerAddr)
	for i := 0; i < cacheSize; i++ {
		key := fmt.Sprintf("%s%d", cacheKeyPrefix, i)
		if err := mc.Set(key, cacheValue, cacheExpiration); err != nil {
			b.Fatalf("failed to set key %s: %v", key, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%s%d", cacheKeyPrefix, i%cacheSize)
		if err := mc.Delete(key); err != nil {
			b.Fatalf("failed to delete key %s: %v", key, err)
		}
	}
}

func BenchmarkMemcachedCacheGetAllKeys(b *testing.B) {
	mc := cache.NewMemcachedCache(memcachedServerAddr)
	for i := 0; i < cacheSize; i++ {
		key := fmt.Sprintf("%s%d", cacheKeyPrefix, i)
		if err := mc.Set(key, cacheValue, cacheExpiration); err != nil {
			b.Fatalf("failed to set key %s: %v", key, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := mc.GetAllKeys(); err != nil {
			b.Fatalf("failed to get all keys: %v", err)
		}
	}
}

func BenchmarkMemcachedCacheDeleteAllKeys(b *testing.B) {
	mc := cache.NewMemcachedCache(memcachedServerAddr)
	for i := 0; i < cacheSize; i++ {
		key := fmt.Sprintf("%s%d", cacheKeyPrefix, i)
		if err := mc.Set(key, cacheValue, cacheExpiration); err != nil {
			b.Fatalf("failed to set key %s: %v", key, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := mc.DeleteAllKeys(); err != nil {
			b.Fatalf("failed to delete all keys: %v", err)
		}
	}
}

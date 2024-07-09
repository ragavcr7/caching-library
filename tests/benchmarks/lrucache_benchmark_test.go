// lrucache_benchmark_test.go

package benchmarks

import (
	"fmt"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func BenchmarkLRUCacheSet(b *testing.B) {
	lru := cache.NewLRUCacheWithMemcached(1000, "localhost:11211")
	defer lru.Clear()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		lru.Set(key, value, expiration)
	}
}

func BenchmarkLRUCacheGet(b *testing.B) {
	lru := cache.NewLRUCacheWithMemcached(1000, "localhost:11211")
	defer lru.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		lru.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		lru.Get(key)
	}
}

func BenchmarkLRUCacheDelete(b *testing.B) {
	lru := cache.NewLRUCacheWithMemcached(1000, "localhost:11211")
	defer lru.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		lru.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		lru.Remove(key)
	}
}

func BenchmarkLRUCacheGetAll(b *testing.B) {
	lru := cache.NewLRUCacheWithMemcached(1000, "localhost:11211")
	defer lru.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		lru.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.GetAll()
	}
}

func BenchmarkLRUCacheDeleteAll(b *testing.B) {
	lru := cache.NewLRUCacheWithMemcached(1000, "localhost:11211")
	defer lru.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		lru.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.DeleteAll()
	}
}

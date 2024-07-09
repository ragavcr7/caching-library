package benchmarks

//becnhmark file for inmemory_cache
import (
	"fmt"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func benchmarkSet(b *testing.B, cache *cache.InMemoryCache, key string, value interface{}, expiration time.Duration) {
	for i := 0; i < b.N; i++ {
		cache.Set(key, value, expiration)
	}
}

func benchmarkGet(b *testing.B, cache *cache.InMemoryCache, key string) {
	// Assuming Get also involves setting the value first for testing purposes
	cache.Set(key, "value", 5*time.Minute)
	for i := 0; i < b.N; i++ {
		cache.Get(key)
	}
}

func benchmarkDelete(b *testing.B, cache *cache.InMemoryCache, key string) {
	// Assuming Delete also involves setting the value first for testing purposes
	cache.Set(key, "value", 5*time.Minute)
	for i := 0; i < b.N; i++ {
		cache.Delete(key)
	}
}

func benchmarkGetAllKeys(b *testing.B, cache *cache.InMemoryCache) {
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value", 2*time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetAllKeys()
	}
}

func benchmarkDeleteAllKeys(b *testing.B, cache *cache.InMemoryCache) {
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value", 2*time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.DeleteAllKeys()
	}
}

func BenchmarkInMemoryCache_Set(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	key := "name"
	value := "ragav"
	expiration := 5 * time.Minute
	b.ResetTimer()
	benchmarkSet(b, cache, key, value, expiration)
}

func BenchmarkInMemoryCache_Get(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	key := "name"
	b.ResetTimer()
	benchmarkGet(b, cache, key)
}

func BenchmarkInMemoryCache_Delete(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	key := "name"
	b.ResetTimer()
	benchmarkDelete(b, cache, key)
}

func BenchmarkInMemoryCache_GetAllKeys(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	b.ResetTimer()
	benchmarkGetAllKeys(b, cache)
}

func BenchmarkInMemoryCache_DeleteAllKeys(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	b.ResetTimer()
	benchmarkDeleteAllKeys(b, cache)
}

package tests

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
	for i := 0; i < b.N; i++ {
		cache.Get(key)
	}
}
func benchmarkDelete(b *testing.B, cache *cache.InMemoryCache, key string) {
	for i := 0; i < b.N; i++ {
		cache.Delete(key)
	}
}
func BenchmakrinMmeoryCache_Set(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	key := "name"
	value := "ragav"
	expiration := 5 * time.Minute
	b.ResetTimer()
	benchmarkSet(b, cache, key, value, expiration)

}
func BenchmakrinMemeryCache_Get(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	key := "name"
	value := "ragav"
	expiration := 5 * time.Minute
	cache.Set(key, value, expiration)
	b.ResetTimer()
	benchmarkGet(b, cache, key)
}
func BenchmarkinMemoryCache_Delete(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	key := "name"
	value := "ragav"
	expiration := 5 * time.Minute
	cache.Set(key, value, expiration)
	b.ResetTimer()
	benchmarkDelete(b, cache, key)
}
func BenchmarkinMemoryCache_GetAllKeys(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value", 2*time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetAllKeys()
	}
}
func BenchmakrinMemeryCache_DeleteAllKeys(b *testing.B) {
	cache := cache.NewInMemoryCache(1000)
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value", 2*time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.DeleteAllKeys()
	}
}

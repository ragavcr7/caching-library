// multicache_benchmark_test.go

package benchmarks

import (
	"fmt"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func BenchmarkMulticacheSet(b *testing.B) {
	mc := cache.NewMulticache(1000, 1000, "localhost:11211")
	defer mc.Clear()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		mc.Set(key, value, expiration)
	}
}

func BenchmarkMulticacheGet(b *testing.B) {
	mc := cache.NewMulticache(1000, 1000, "localhost:11211")
	defer mc.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		mc.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		mc.Get(key)
	}
}

func BenchmarkMulticacheRemove(b *testing.B) {
	mc := cache.NewMulticache(1000, 1000, "localhost:11211")
	defer mc.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		mc.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		mc.Remove(key)
	}
}

func BenchmarkMulticacheGetAllKeys(b *testing.B) {
	mc := cache.NewMulticache(1000, 1000, "localhost:11211")
	defer mc.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		mc.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.GetAllKeys()
	}
}

func BenchmarkMulticacheDeleteAllKeys(b *testing.B) {
	mc := cache.NewMulticache(1000, 1000, "localhost:11211")
	defer mc.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		mc.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.DeleteAllKeys()
	}
}

func BenchmarkMulticacheClear(b *testing.B) {
	mc := cache.NewMulticache(1000, 1000, "localhost:11211")
	defer mc.Clear()

	// Populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expiration := time.Minute
		mc.Set(key, value, expiration)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Clear()
	}
}

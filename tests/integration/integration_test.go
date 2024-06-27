package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
	"github.com/ragavcr7/caching-library/cache_interface"
)

const (
	memcachedServerAddr = "localhost:11211"
)

func setupInMemoryCache() *cache.InMemoryCache {
	return cache_interface.NewInMemoryCache(5)
}

func setupMemcachedCache() *cache.MemcachedCache {
	return cache_interface.NewMemcachedCache(memcachedServerAddr)
}

func TestIntegration_InmemoryCache(t *testing.T) {
	inMemoryCache := setupInMemoryCache()
	runIntegrationTests(t, inMemoryCache)
}

func TestIntegration_MemcachedCache(t *testing.T) {
	MemcachedCache := setupMemcachedCache()
	runIntegrationTests(t, MemcachedCache)
}
func runIntegrationTests(t *testing.T, cacheInterface cache.Cache) {
	t.Run("SetAndGet", func(t *testing.T) {
		testSetAndGet(t, cacheInterface)
	})

	t.Run("Delete", func(t *testing.T) {
		testDelete(t, cacheInterface)
	})
	t.Run("ConcurrentAccess", func(t *testing.T) {
		testConcurrentAccess(t, cacheInterface)
	})
}

func testSetAndGet(t *testing.T, cacheInterface cache.Cache) {
	err := CacheInterface.Set("name", "Kumar", 2*time.Minute)
	if err != nil {
		t.Fatalf("failed to set value: %v", err)
	}
	val, err := CacheInterface.Get("name")
	if err != nil {
		t.Fatalf("Key not found :%v", err)
	}
	if val != "Kumar" {
		t.Fatalf("Expeceted value is: kumar but Got: %v", val)
	}

	//Test expiration
	time.Sleep(2 * time.Minute)
	_, err = CacheInterface.Get("name")
	if err == nil {
		t.Errorf("Key is expected to be expired but found...")
	}
}

func testDelete(t *Testing.T, cacheInterface cache.Cache) {
	err := CacheInterface.Set("name", "Thor", 2*time.Minute)
	if err != nil {
		t.Fatalf("Failed to set the value : %v", err)
	}
	CacheInterface.Delete("name")
	_, found := CacheInterface.Get("name")
	if found == nil {
		t.Fatalf("Key: name is expected to be delete but Found : %v", err)
	}

}
func testConcurrentAccess(t *testing.T, CacheInterface cache.Cache) {
	var wg sync.WaitGroup
	const numGoRoutines = 10
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		func(i int) {
			defer wg.Done()
			err := CacheInterface.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
			if err != nil {
				t.Fatalf("failed to set key : %d and got :%v", i, err)
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		func(i int) {
			defer wg.Done()
			val, err := CacheInterface.Get(fmt.Sprintf("key%d"), i)
			if err != nil {
				t.Fatalf("failed to get key : %d and got : %v", i, err)
			}
			expectedValue := fmt.Sprintf("key%d", i)
			if expectedValue != val {
				t.Errorf("Value expected is: %v but Got: %v", expectedValue, val)
			}
		}(i)
	}
	wg.Wait()
}

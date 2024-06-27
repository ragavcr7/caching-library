/* this file is to test the memcache.go implementations... like set ,get,getall,deleteall methods*/
package tests

/*
import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

const serverAddr = "localhost:11211"

func TestMemcachedCache_SetandGet(T *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	err := cache.Set("name", "pichumani", 1*time.Minute)
	if err != nil {
		T.Errorf("failed to set key : %v", err)
	}
	val, err := cache.Get("name")
	if err != nil {
		T.Fatalf("can't able to fetch the key: %v", err)
	}
	if val != "pichumani" {
		T.Errorf("Expected value:pichumani but Got : %v", val)
	}
	//check for expiration
	time.Sleep(1 * time.Minute)
	_, err = cache.Get("name")
	if err != nil {
		T.Errorf("key: name is expected to be expired but found")
	}
}

func TestMemCache_Delete(T *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	err := cache.Set("name", "mahi", 1*time.Minute)
	if err != nil {
		T.Errorf("failed to set key:%v", err)

	}
	err = cache.Delete("name")
	if err != nil {
		T.Errorf("Failed to delete the key:%v", err)
	}
	_, err = cache.Get("name")
	if err == nil {
		T.Errorf("DANGER: Deleted key is still accessible means key didnt get deleted")
	}
}
func TestMemCache_GetAllKeys(T *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	_, err := cache.GetAllKeys()
	if err != nil {
		T.Errorf("expected to get error, but got none")
	}
}
func TestMemCache_DeleteAllKeys(T *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	err := cache.DeleteAllKeys()
	if err != nil {
		T.Errorf("expected to get not able to delete all keys at once error , but got none")
	}
}

// this method is used to check concurrent acces to cache to make sure thread-saftey and synchronisation.(like whether the CRUD operations are synchronised)
func TestInMemoryCache_SetGetConcurrent(t *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	var wg sync.WaitGroup
	const numGoRoutines = 10
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 1*time.Minute)
			if err != nil {
				t.Errorf("failed to set key %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			val, err := cache.Get(fmt.Sprintf("key%d", i))
			if err != nil {
				t.Errorf("Unable to fetch the value for the key: %v", err)
			}
			expectedValue := fmt.Sprintf("value%d", i)
			if val != expectedValue {
				t.Errorf("Expected value: %v but Got value: %v", expectedValue, val)
			}
		}(i)
	}
	wg.Wait()
}
*/

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

const serverAddr = "localhost:11211"

func TestMemcachedCache_SetandGet(t *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	err := cache.Set("name", "pichumani", 1*time.Minute)
	if err != nil {
		t.Errorf("failed to set key: %v", err)
	}
	val, err := cache.Get("name")
	if err != nil {
		t.Fatalf("can't able to fetch the key: %v", err)
	}
	if val != string("pichumani") {
		t.Errorf("Expected value: pichumani but Got: %v", val)
	}
	//check for expiration
	time.Sleep(2 * time.Minute)
	_, err = cache.Get("name")
	if err == nil {
		t.Errorf("key: name is expected to be expired but found")
	}
}

func TestMemCache_Delete(t *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	err := cache.Set("name", "mahi", 1*time.Minute)
	if err != nil {
		t.Errorf("failed to set key: %v", err)
	}
	err = cache.Delete("name")
	if err != nil {
		t.Errorf("Failed to delete the key: %v", err)
	}
	_, err = cache.Get("name")
	if err == nil {
		t.Errorf("DANGER: Deleted key is still accessible means keys doesn't get deleted")
	}
}

func TestMemCache_GetAllKeys(t *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	_, err := cache.GetAllKeys()
	if err == nil {
		t.Errorf("expected to get error, but got none")
	}
}

func TestMemCache_DeleteAllKeys(t *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	err := cache.DeleteAllKeys()
	if err == nil {
		t.Errorf("expected to get not able to delete all keys at once error, but got none")
	}
}

func TestInMemoryCache_SetGetConcurrent(t *testing.T) {
	cache := cache.NewMemcachedCache(serverAddr)
	var wg sync.WaitGroup
	const numGoRoutines = 10
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 1*time.Minute)
			if err != nil {
				t.Errorf("failed to set key %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			val, err := cache.Get(fmt.Sprintf("key%d", i))
			if err != nil {
				t.Errorf("Unable to fetch the value for the key: %v", err)
			}
			expectedValue := fmt.Sprintf("value%d", i)
			if val != expectedValue {
				t.Errorf("Expected value: %v but Got value: %v", expectedValue, val)
			}
		}(i)
	}
	wg.Wait()
}

/*
what are the cases we need to consider
1.invdiudal set ---- what are the possible cases ------1. setting gets passed (valid input), setting gets failed(invalid input).
2.indvidual get ----- what are the possible cases ---- 1. get gets failed(when key not present),passed(when key is present)
3.indvidual delete----- what are the possible cases ----1. valid key ...2.invalid key
4.group of delete  ----- delete allll..... necessarily dont have to be tested...
*/
package tests

import (
	"testing"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func TestInMemoryCache_setAndGet(t *testing.T) {
	cache := cache.NewInMemoryCache(5)
	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 1*time.Minute)
	val, found := cache.Get("key1")
	if !found || val != "value1" {
		t.Errorf("expected value for key1 is : value1 but Got: %V", val)
	}
	val, found = cache.Get("key2")
	if !found || val != "value2" {
		t.Errorf("expected value for key2 is :value2 but Got: %v", val)
	}

	//test expiration
	time.Sleep(5 * time.Minute)
	_, found = cache.Get("key1")
	if found {
		t.Errorf("Expected key1 to be expired,but it was found in cache")
	}
}

func TestInMemoryCache_Delete(t *testing.T) {
	cache := cache.NewInMemoryCache(5)
	cache.Set("name", "ragav", 1*time.Minute)
	cache.Delete("name")
	_, found := cache.Get("name")
	if found {
		t.Errorf("Key name is expected to be deleted but found")
	}
}
func TestInMemoryCache_GetAllKeys(t *testing.T) {
	cache := cache.NewInMemoryCache(5)
	cache.Set("name1", "undertaker", 1*time.Minute)
	cache.Set("name2", "JhonCena", 1*time.Minute)
	cache.Set("name3", "mrlesnar", 1*time.Minute)
	keys := cache.GetAllKeys()
	key_slice := []string{"name1", "name2", "name3"}
	if len(keys) != len(key_slice) {
		t.Errorf("Expected %d keys, Got %d keys", len(key_slice), len(keys))
	}
	for _, key1 := range key_slice {
		found := false
		for _, key2 := range keys {
			if key1 != key2 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected key %s in keys, but not found", key1)
		}
	}
}
func TestInMemoryCache_DeleteAllKeys(t *testing.T) {
	cache := cache.NewInMemoryCache(5)
	cache.Set("name1", "ronaldo", 1*time.Minute)
	cache.Set("name2", "messi", 1*time.Minute)
	cache.Set("name3", "neymar", 1*time.Minute)
	cache.DeleteAllKeys()
	keys := cache.GetAllKeys()
	if len(keys) != 0 {
		t.Errorf("Expected length of cache is to be empty i.e 0 but found %d: keys ", len(keys))
	}
}

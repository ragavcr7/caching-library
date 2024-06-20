package cache

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedCache struct { //it represents a client
	client *memcache.Client
}

func NewMemcachedCache(serverAddr string) *MemcachedCache { //using pointer here to avoid the copy
	client := memcache.New(serverAddr)
	return &MemcachedCache{
		client: client,
	}
}

// this method is used to add a new key value pair to cache and using the capital S for set to make it as public
func (mc *MemcachedCache) Set(key string, value interface{}, expiration time.Duration) error { //we all know key will always be string and value can be of any comparable type so we sets interface for value
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(fmt.Sprintf("%v", value)), //byte dt since all the information from client is gonna be in byte format
		Expiration: int32(expiration.Seconds()),
	}
	return mc.client.Set(item) //manages the connection to memcached server so we can perform set,get,delete operations
}

// this method is used to fetch the value based on the passed key from the memcached cache
func (mc *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := mc.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, fmt.Errorf("%v not found in the Memcached cache", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key %v from Memcached: %v", key, err)
	}
	return string(item.Value), nil
}

func (mc *MemcachedCache) Delete(key string) error {
	return mc.client.Delete(key)
}

package cache

// inmemory.go
import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// entry represents a cache entry stored in the InMemoryCache.
type entry struct {
	key        string
	value      interface{}
	expiration time.Time
}

// InMemoryCache represents an in-memory cache with LRU eviction and expiration handling.
type InMemoryCache struct {
	capacity  int
	cache     map[string]*list.Element
	evictList *list.List
	mu        sync.Mutex
}

// NewInMemoryCache creates a new instance of InMemoryCache with the specified capacity.
func NewInMemoryCache(capacity int) *InMemoryCache {
	return &InMemoryCache{
		capacity:  capacity,
		cache:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Set adds a new key-value pair to the in-memory cache with optional expiration.
func (imc *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	// Check if the key exists in the cache
	if elem, ok := imc.cache[key]; ok {
		// If key exists, update the value, expiration, and move the element to front (MRU position)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).expiration = time.Now().Add(expiration)
		imc.evictList.MoveToFront(elem)
	} else {
		// If key does not exist, add the new entry to the cache
		if len(imc.cache) >= imc.capacity {
			// Evict least recently used element if cache is full
			imc.evictOldest()
		}

		// Add new entry to the cache with expiration time
		newElem := imc.evictList.PushFront(&entry{
			key:        key,
			value:      value,
			expiration: time.Now().Add(expiration),
		})
		imc.cache[key] = newElem
	}

	// Start a goroutine to periodically check and evict expired entries
	go imc.startEvictionTicker()
	return fmt.Errorf("passed succesfully")
}

// Get retrieves a value from the in-memory cache based on the key.
func (imc *InMemoryCache) Get(key string) (interface{}, bool) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		// Check if the entry is expired
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			// If expired, delete the entry and return not found
			imc.removeElement(elem)
			return nil, false
		}

		// Move the accessed element to front (MRU position)
		imc.evictList.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

// Delete removes a key-value pair from the in-memory cache.
func (imc *InMemoryCache) Delete(key string) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		imc.removeElement(elem)
	}
	return fmt.Errorf("DELETE AINT DELETING")
}

// GetAllKeys retrieves all keys from the in-memory cache.
func (imc *InMemoryCache) GetAllKeys() []string {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	keys := make([]string, 0, len(imc.cache))
	for key := range imc.cache {
		keys = append(keys, key)
	}
	return keys
}

// DeleteAllKeys deletes all keys from the in-memory cache.
func (imc *InMemoryCache) DeleteAllKeys() error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	imc.cache = make(map[string]*list.Element)
	imc.evictList.Init()
	return nil
}

// evictOldest evicts the least recently used entry from the cache.
func (imc *InMemoryCache) evictOldest() {
	if imc.evictList.Len() > 0 {
		oldest := imc.evictList.Back()
		if oldest != nil {
			imc.removeElement(oldest)
		}
	}
}

// removeElement removes an element from the cache and evictList.
func (imc *InMemoryCache) removeElement(e *list.Element) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	imc.evictList.Remove(e)
	delete(imc.cache, e.Value.(*entry).key)
}

// startEvictionTicker starts a ticker to periodically check and evict expired entries.
func (imc *InMemoryCache) startEvictionTicker() {
	ticker := time.NewTicker(5 * time.Minute) // Adjust tick interval as needed
	defer ticker.Stop()

	for range ticker.C {
		imc.evictExpiredEntries()
	}
}

// evictExpiredEntries evicts all expired entries from the cache.
func (imc *InMemoryCache) evictExpiredEntries() {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	for _, elem := range imc.cache { //error key unused oocured here
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			imc.removeElement(elem)
		}
	}
}

/*
import (
	"container/list"
	"sync"
	"time"
	"fmt"
)

// entry represents a cache entry stored in the InMemoryCache.
type entry struct {
	key        string
	value      interface{}
	expiration time.Time
}

// InMemoryCache represents an in-memory cache with LRU eviction and expiration handling.
type InMemoryCache struct {
	capacity  int
	cache     map[string]*list.Element
	evictList *list.List
	mu        sync.Mutex
}

// NewInMemoryCache creates a new instance of InMemoryCache with the specified capacity.
func NewInMemoryCache(capacity int) *InMemoryCache {
	return &InMemoryCache{
		capacity:  capacity,
		cache:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Set adds a new key-value pair to the in-memory cache with optional expiration.
func (imc *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	// Check if the key exists in the cache
	if elem, ok := imc.cache[key]; ok {
		// If key exists, update the value, expiration, and move the element to front (MRU position)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).expiration = time.Now().Add(expiration)
		imc.evictList.MoveToFront(elem)
	} else {
		// If key does not exist, add the new entry to the cache
		if len(imc.cache) >= imc.capacity {
			// Evict least recently used element if cache is full
			imc.evictOldest()
		}

		// Add new entry to the cache with expiration time
		newElem := imc.evictList.PushFront(&entry{
			key:        key,
			value:      value,
			expiration: time.Now().Add(expiration),
		})
		imc.cache[key] = newElem
	}

	// Start a goroutine to periodically check and evict expired entries
	go imc.startEvictionTicker()

	return nil
}

// Get retrieves a value from the in-memory cache based on the key.
func (imc *InMemoryCache) Get(key string) (interface{}, error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		// Check if the entry is expired
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			// If expired, delete the entry and return not found
			imc.removeElement(elem)
			return nil, fmt.Errorf("key expired")
		}

		// Move the accessed element to front (MRU position)
		imc.evictList.MoveToFront(elem)
		return elem.Value.(*entry).value, nil
	}
	return nil, fmt.Errorf("key not found")
}

// Delete removes a key-value pair from the in-memory cache.
func (imc *InMemoryCache) Delete(key string) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		imc.removeElement(elem)
		return nil
	}
	return fmt.Errorf("key not found")
}

// GetAllKeys retrieves all keys from the in-memory cache.
func (imc *InMemoryCache) GetAllKeys() ([]string, error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	keys := make([]string, 0, len(imc.cache))
	for key := range imc.cache {
		keys = append(keys, key)
	}
	return keys, nil
}

// DeleteAllKeys deletes all keys from the in-memory cache.
func (imc *InMemoryCache) DeleteAllKeys() error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	imc.cache = make(map[string]*list.Element)
	imc.evictList.Init()
	return nil
}

// evictOldest evicts the least recently used entry from the cache.
func (imc *InMemoryCache) evictOldest() {
	if imc.evictList.Len() > 0 {
		oldest := imc.evictList.Back()
		if oldest != nil {
			imc.removeElement(oldest)
		}
	}
}

// removeElement removes an element from the cache and evictList.
func (imc *InMemoryCache) removeElement(e *list.Element) {
	imc.evictList.Remove(e)
	delete(imc.cache, e.Value.(*entry).key)
}

// startEvictionTicker starts a ticker to periodically check and evict expired entries.
func (imc *InMemoryCache) startEvictionTicker() {
	ticker := time.NewTicker(5 * time.Minute) // Adjust tick interval as needed
	defer ticker.Stop()

	for range ticker.C {
		imc.evictExpiredEntries()
	}
}

// evictExpiredEntries evicts all expired entries from the cache.
func (imc *InMemoryCache) evictExpiredEntries() {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	for _, elem := range imc.cache {
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			imc.removeElement(elem)
		}
	}
}
*/

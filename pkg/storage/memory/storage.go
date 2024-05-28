package memory

import (
	"sync"
	"time"
)

// CacheItem represents an item stored in the cache with its associated TTL.
type CacheItem struct {
	value  string
	expiry time.Time // TTL for a key
}

// Cache represents an in-memory key-value store with expiry support.
type Cache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

// NewCache creates and initializes a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

// Set adds or updates a key-value pair in the cache with the given TTL.
func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheItem{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

// Get retrieves the value associated with the given key from the cache.
// It also checks for expiry and removes expired items.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.data[key]
	if !ok {
		return "", false
	}
	// item found - check for expiry
	if item.expiry.Before(time.Now()) {
		// remove entry from cache if time is beyond the expiry
		delete(c.data, key)
		return "", false
	}
	return item.value, true
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear removes all key-value pairs from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]CacheItem)
}

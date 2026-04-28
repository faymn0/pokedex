package pokeapi

import (
	"sync"
	"time"
)

type Cache struct {
	data  map[string]cacheEntry
	mutex sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		value:     value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cacheEntry, ok := c.data[key]
	return cacheEntry.value, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.Sub(entry.createdAt) > interval {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}

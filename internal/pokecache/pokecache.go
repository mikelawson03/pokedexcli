package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entry    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct{
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entry[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, entry := range c.entry {
		if now.Sub(entry.createdAt) > c.interval {  // calculates difference between createdAt and now and returns true if > interval
			delete(c.entry, key)
		}
			
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entry:    make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()

	return c
}
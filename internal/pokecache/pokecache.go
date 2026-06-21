package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
	stop    chan struct{}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{entries: make(map[string]cacheEntry), stop: make(chan struct{})}

	go c.reapLoop(interval)

	return c
}

func (c *Cache) Close() {
	close(c.stop)
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			c.mu.Lock()
			for k, v := range c.entries {
				if v.createdAt.Before(t.Add(-interval)) {
					delete(c.entries, k)
				}
			}
			c.mu.Unlock()

		case <-c.stop:
			return
		}
	}
}

package services

import (
	"fmt"
	"sync"
	"time"
)

type InMemoryUrlCache struct {
	cache map[string]InMemoryUrlCacheEntry
	mut   *sync.Mutex
}
type InMemoryUrlCacheEntry struct {
	Url string
	TTL time.Time
}

func InitInMemoryUrlCache() *InMemoryUrlCache {
	return &InMemoryUrlCache{
		cache: make(map[string]InMemoryUrlCacheEntry),
		mut:   &sync.Mutex{},
	}
}

func (c *InMemoryUrlCache) Set(objectName string, url string, ttl time.Duration) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	if c.cache == nil {
		c.cache = make(map[string]InMemoryUrlCacheEntry)
	}
	c.cache[objectName] = InMemoryUrlCacheEntry{Url: url, TTL: time.Now().Add(ttl)}
	return nil
}

func (c *InMemoryUrlCache) Get(objectName string) (string, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	entry, exists := c.cache[objectName]
	if !exists {
		return "", fmt.Errorf("cache miss for object: %s", objectName)
	}
	if entry.TTL.Before(time.Now()) {
		delete(c.cache, objectName)
		return "", fmt.Errorf("cache expired for object: %s", objectName)
	}
	return entry.Url, nil
}

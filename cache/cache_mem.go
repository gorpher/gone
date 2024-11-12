package cache

import (
	"encoding/json"
	"sync"
	"time"
)

type MemoryCache struct {
	mutex sync.RWMutex
	cache map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: map[string]string{},
	}
}
func (c *MemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return []byte(c.cache[k]), nil
}

func (c *MemoryCache) Set(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	switch v := value.(type) {
	case string:
		c.cache[key] = v
	case []byte:
		c.cache[key] = string(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		c.cache[key] = string(data)
	}
	return nil
}

func (c *MemoryCache) Del(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, key)
	return nil
}

func (c *MemoryCache) SetWithTTL(key string, value string, duration time.Duration) error {
	_ = time.AfterFunc(duration, func() {
		c.mutex.Lock()
		delete(c.cache, key)
		c.mutex.Unlock()
	})
	return c.Set(key, value)
}

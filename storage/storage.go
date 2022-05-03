package storage

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("value not found")
var ErrCannotBeEmpty = errors.New("key cannot be empty")

type Cache interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type MemCache struct {
	cache map[string]string
	mu    sync.RWMutex
}

func NewMemCache() *MemCache {
	return &MemCache{
		cache: make(map[string]string, 0),
		mu:    sync.RWMutex{},
	}
}

func (c *MemCache) Set(key, value string) error {
	if len(key) == 0 {
		return ErrCannotBeEmpty
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = value

	return nil
}

func (c *MemCache) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrCannotBeEmpty
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.cache[key]

	if !ok {
		return "", ErrNotFound
	}

	return v, nil
}

func (c *MemCache) Delete(key string) error {
	if len(key) == 0 {
		return ErrCannotBeEmpty
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	//if _, ok := c.cache[key]; !ok {
	//	return ErrNotFound
	//}

	delete(c.cache, key)

	return nil
}

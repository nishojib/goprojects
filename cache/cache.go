package cache

import (
	"slices"
	"sync"
	"time"
)

// Cache is a key-value storage.
type Cache[K comparable, V any] struct {
	ttl               time.Duration
	mu                sync.Mutex
	data              map[K]entryWithTimeout[V]
	maxSize           int
	chronologicalKeys []K
}

type entryWithTimeout[V any] struct {
	value   V
	expires time.Time
}

// New creates a usable Cache with an initialized data.
func New[K comparable, V any](maxSize int, ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		ttl:               ttl,
		data:              make(map[K]entryWithTimeout[V]),
		maxSize:           maxSize,
		chronologicalKeys: make([]K, 0, maxSize),
	}
}

// Read returns the associated value for a key,
// and a boolean to false if the key is absent.
func (c *Cache[K, V]) Read(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zeroV V

	v, ok := c.data[key]
	switch {
	case !ok:
		return zeroV, false
	case v.expires.Before(time.Now()):
		c.deleteKeyValue(key)
		return zeroV, false
	default:
		return v.value, true
	}
}

// Upsert overwrites the value for a given key.
func (c *Cache[K, V]) Upsert(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, alreadyPresent := c.data[key]
	switch {
	case alreadyPresent:
		c.deleteKeyValue(key)
	case len(c.data) == c.maxSize:
		c.deleteKeyValue(c.chronologicalKeys[0])
	}

	c.addKeyValue(key, value)
}

// Delete removes the entry for the given key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.deleteKeyValue(key)
}

// addKeyValue inserts a key and its value into the cache.
func (c *Cache[K, V]) addKeyValue(key K, value V) {
	c.data[key] = entryWithTimeout[V]{
		value:   value,
		expires: time.Now().Add(c.ttl),
	}
	c.chronologicalKeys = append(c.chronologicalKeys, key)
}

// deleteKeyValue removes a key and its associated value from the cache.
func (c *Cache[K, V]) deleteKeyValue(key K) {
	c.chronologicalKeys = slices.DeleteFunc(c.chronologicalKeys, func(k K) bool { return k == key })
	delete(c.data, key)
}

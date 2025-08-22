package build

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
)

type parseCacheEntry[V any] struct {
	value V
	mu    sync.Mutex
	dirty bool
}

type parseCache[K comparable, V any] struct {
	entries collections.SyncMap[K, *parseCacheEntry[V]]
}

func (c *parseCache[K, V]) LoadOrStoreNew(key K, parse func() V) (V, bool) {
	return c.LoadOrStoreNewIf(key, func() (V, bool) {
		return parse(), true
	})
}

func (c *parseCache[K, V]) LoadOrStoreNewIf(key K, parse func() (V, bool)) (V, bool) {
	newEntry := &parseCacheEntry[V]{}
	newEntry.mu.Lock()
	defer newEntry.mu.Unlock()
	if entry, loaded := c.entries.LoadOrStore(key, newEntry); loaded && !entry.dirty {
		// Ensure it was parsed before returning
		entry.mu.Lock()
		defer entry.mu.Unlock()
		return entry.value, true
	}
	value, ok := parse()
	if ok {
		newEntry.value = value
	} else {
		// Dont use the cache entry
		newEntry.dirty = true
		c.entries.Delete(key)
	}
	return value, false
}

func (c *parseCache[K, V]) Load(key K) (V, bool) {
	if entry, ok := c.entries.Load(key); ok && !entry.dirty {
		entry.mu.Lock()
		defer entry.mu.Unlock()
		return entry.value, true
	}
	var zero V
	return zero, false
}

func (c *parseCache[K, V]) Store(key K, value V) {
	c.entries.Store(key, &parseCacheEntry[V]{value: value})
}

func (c *parseCache[K, V]) Delete(key K) {
	c.entries.Delete(key)
}

func (c *parseCache[K, V]) Reset() {
	c.entries = collections.SyncMap[K, *parseCacheEntry[V]]{}
}

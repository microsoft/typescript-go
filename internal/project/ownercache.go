package project

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
)

type ownerCacheEntry[V any] struct {
	mu     sync.Mutex
	value  V
	owners map[uint64]struct{}
}

type OwnerCache[K comparable, V any, LoadArgs any] struct {
	entries collections.SyncMap[K, *ownerCacheEntry[V]]

	isExpired func(K, V, LoadArgs) bool
	parse     func(K, LoadArgs) V
}

func NewOwnerCache[K comparable, V any, LoadArgs any](
	parse func(K, LoadArgs) V,
	isExpired func(K, V, LoadArgs) bool,
) *OwnerCache[K, V, LoadArgs] {
	return &OwnerCache[K, V, LoadArgs]{
		isExpired: isExpired,
		parse:     parse,
	}
}

func (c *OwnerCache[K, V, LoadArgs]) LoadAndAcquire(identity K, owner uint64, loadArgs LoadArgs) V {
	entry, loaded := c.loadOrStoreLockedEntry(identity)
	defer entry.mu.Unlock()
	if !loaded || c.isExpired != nil && c.isExpired(identity, entry.value, loadArgs) {
		entry.value = c.parse(identity, loadArgs)
	}
	entry.owners[owner] = struct{}{}
	return entry.value
}

func (c *OwnerCache[K, V, LoadArgs]) Acquire(identity K, owner uint64, value V) {
	entry, loaded := c.loadOrStoreLockedEntry(identity)
	defer entry.mu.Unlock()
	if !loaded {
		entry.value = value
	}
	entry.owners[owner] = struct{}{}
}

// TryAcquire adds an owner to an existing live entry.
// It returns false instead of creating a new entry when the value has already
// been fully released.
func (c *OwnerCache[K, V, LoadArgs]) TryAcquire(identity K, owner uint64) bool {
	entry, ok := c.entries.Load(identity)
	if !ok {
		return false
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	if len(entry.owners) == 0 {
		return false
	}
	entry.owners[owner] = struct{}{}
	return true
}

func (c *OwnerCache[K, V, LoadArgs]) Has(identity K) bool {
	_, ok := c.entries.Load(identity)
	return ok
}

func (c *OwnerCache[K, V, LoadArgs]) Release(identity K, owner uint64) {
	entry, ok := c.entries.Load(identity)
	if !ok {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	delete(entry.owners, owner)
	if len(entry.owners) == 0 {
		c.entries.Delete(identity)
	}
}

func (c *OwnerCache[K, V, LoadArgs]) loadOrStoreLockedEntry(key K) (*ownerCacheEntry[V], bool) {
	entry := &ownerCacheEntry[V]{
		owners: make(map[uint64]struct{}),
	}
	entry.mu.Lock()
	existing, loaded := c.entries.LoadOrStore(key, entry)
	if loaded {
		entry.mu.Unlock()
		existing.mu.Lock()
		if len(existing.owners) == 0 {
			existing.mu.Unlock()
			return c.loadOrStoreLockedEntry(key)
		}
		return existing, true
	}
	return entry, false
}

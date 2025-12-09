package project

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/zeebo/xxh3"
)

// VersionedCacheKey combines an identity key with a content hash.
// The identity identifies the logical resource (e.g., file path + parse options),
// while the hash identifies a specific version of that resource's content.
type VersionedCacheKey[K comparable] struct {
	Identity K
	Hash     xxh3.Uint128
}

// versionedCacheEntry holds a cached value and its reference count.
// Entries are immutable once created - the value never changes.
type versionedCacheEntry[V any] struct {
	mu       sync.Mutex
	value    V
	refCount int
}

// VersionedCacheOptions configures cache behavior.
type VersionedCacheOptions struct {
	// DisableDeletion prevents entries from being removed from the cache.
	// Used for testing.
	DisableDeletion bool
}

// VersionedCache is a thread-safe, ref-counted cache that supports multiple
// versions of the same logical resource. Entries are keyed by both an identity
// (e.g., filename + options) and a content hash, so different versions of the
// same file can coexist in the cache.
//
// This design supports snapshot-based systems where multiple snapshots may
// reference different versions of the same file simultaneously (e.g., speculative
// edits alongside mainline server state).
//
// Type parameters:
//   - K: The identity key type (e.g., file path, parse options struct)
//   - V: The cached value type (e.g., *ast.SourceFile)
type VersionedCache[K comparable, V any, ParseArgs any] struct {
	Options VersionedCacheOptions
	entries collections.SyncMap[VersionedCacheKey[K], *versionedCacheEntry[V]]
	parse   func(K, ParseArgs) V
}

// Acquire retrieves or creates a cache entry for the given identity and hash.
// If an entry exists with matching identity and hash, its refcount is incremented
// and the cached value is returned. Otherwise, parse() is called to create the
// value, which is stored and returned with refcount 1.
//
// The caller is responsible for calling Deref when done with the value.
func (c *VersionedCache[K, V, ParseArgs]) Acquire(identity K, hash xxh3.Uint128, parseArgs ParseArgs) V {
	key := VersionedCacheKey[K]{Identity: identity, Hash: hash}
	entry, loaded := c.loadOrStoreNewLockedEntry(key)
	defer entry.mu.Unlock()
	if !loaded {
		// New entry - parse the value
		entry.value = c.parse(identity, parseArgs)
	}
	return entry.value
}

// Ref increments the reference count for an existing entry.
// Panics if the entry does not exist.
func (c *VersionedCache[K, V, ParseArgs]) Ref(identity K, hash xxh3.Uint128) {
	key := VersionedCacheKey[K]{Identity: identity, Hash: hash}
	entry, ok := c.entries.Load(key)
	if !ok {
		panic("versioned cache entry not found")
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	if entry.refCount <= 0 && !c.Options.DisableDeletion {
		// Entry was deleted while we were acquiring the lock
		newEntry, loaded := c.loadOrStoreNewLockedEntry(key)
		defer newEntry.mu.Unlock()
		if !loaded {
			newEntry.value = entry.value
		}
		return
	}
	entry.refCount++
}

// Deref decrements the reference count for an entry.
// When the refcount reaches zero, the entry is removed from the cache
// (unless DisableDeletion is set).
func (c *VersionedCache[K, V, ParseArgs]) Deref(identity K, hash xxh3.Uint128) {
	key := VersionedCacheKey[K]{Identity: identity, Hash: hash}
	entry, ok := c.entries.Load(key)
	if !ok {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	entry.refCount--
	if entry.refCount <= 0 && !c.Options.DisableDeletion {
		c.entries.Delete(key)
	}
}

// Has returns true if an entry exists for the given identity and hash.
func (c *VersionedCache[K, V, ParseArgs]) Has(identity K, hash xxh3.Uint128) bool {
	key := VersionedCacheKey[K]{Identity: identity, Hash: hash}
	_, ok := c.entries.Load(key)
	return ok
}

// loadOrStoreNewLockedEntry loads an existing entry or creates a new one.
// The returned entry's mutex is locked and its refCount is incremented
// (or initialized to 1 in the case of a new entry).
func (c *VersionedCache[K, V, ParseArgs]) loadOrStoreNewLockedEntry(key VersionedCacheKey[K]) (*versionedCacheEntry[V], bool) {
	entry := &versionedCacheEntry[V]{refCount: 1}
	entry.mu.Lock()
	existing, loaded := c.entries.LoadOrStore(key, entry)
	if loaded {
		existing.mu.Lock()
		if existing.refCount <= 0 && !c.Options.DisableDeletion {
			// Existing entry was deleted while we were acquiring the lock
			existing.mu.Unlock()
			return c.loadOrStoreNewLockedEntry(key)
		}
		existing.refCount++
		return existing, true
	}
	return entry, false
}

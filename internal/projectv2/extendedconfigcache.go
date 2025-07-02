package projectv2

import (
	"crypto/sha256"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type extendedConfigCache struct {
	entries collections.SyncMap[tspath.Path, *extendedConfigCacheEntry]
}

type extendedConfigCacheEntry struct {
	mu       sync.Mutex
	entry    *tsoptions.ExtendedConfigCacheEntry
	hash     [sha256.Size]byte
	refCount int
}

func (c *extendedConfigCache) acquire(fh fileHandle, path tspath.Path, parse func() *tsoptions.ExtendedConfigCacheEntry) *tsoptions.ExtendedConfigCacheEntry {
	entry, loaded := c.loadOrStoreNewLockedEntry(fh, path)
	defer entry.mu.Unlock()
	if !loaded || entry.hash != fh.Hash() {
		// Reparse the config if the hash has changed, or parse for the first time.
		entry.entry = parse()
		entry.hash = fh.Hash()
	}
	return entry.entry
}

func (c *extendedConfigCache) release(path tspath.Path) {
	if entry, ok := c.entries.Load(path); ok {
		entry.mu.Lock()
		entry.refCount--
		remove := entry.refCount <= 0
		entry.mu.Unlock()
		if remove {
			c.entries.Delete(path)
		}
	}
}

// loadOrStoreNewLockedEntry loads an existing entry or creates a new one. The returned
// entry's mutex is locked and its refCount is incremented (or initialized to 1
// in the case of a new entry).
func (c *extendedConfigCache) loadOrStoreNewLockedEntry(
	fh fileHandle,
	path tspath.Path,
) (*extendedConfigCacheEntry, bool) {
	entry := &extendedConfigCacheEntry{refCount: 1}
	entry.mu.Lock()
	existing, loaded := c.entries.LoadOrStore(path, entry)
	if loaded {
		existing.mu.Lock()
		existing.refCount++
		return existing, true
	}
	return entry, false
}

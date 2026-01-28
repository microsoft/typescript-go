package build

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/microsoft/typescript-go/internal/collections"
)

type parseCacheEntry[V comparable] struct {
	value V
	reuse int
	mu    sync.Mutex
}

type parseCache[K comparable, V comparable] struct {
	entries      collections.SyncMap[K, *parseCacheEntry[V]]
	entriesCount atomic.Int64
	reusedCount  atomic.Int64
	calledCount  atomic.Int64
}

func (c *parseCache[K, V]) loadOrStore(
	key K,
	parse func(K) V,
) (value V, loaded bool) {
	c.calledCount.Add(1)
	newEntry := &parseCacheEntry[V]{}
	newEntry.mu.Lock()
	defer newEntry.mu.Unlock()
	if entry, loaded := c.entries.LoadOrStore(key, newEntry); loaded {
		entry.mu.Lock()
		defer entry.mu.Unlock()
		c.reusedCount.Add(1)
		entry.reuse++
		return entry.value, true
	} else {
		c.entriesCount.Add(1)
	}
	newEntry.value = parse(key)
	return newEntry.value, false
}

func (c *parseCache[K, V]) store(key K, value V) {
	c.entries.Store(key, &parseCacheEntry[V]{value: value})
}

func (c *parseCache[K, V]) delete(key K) {
	c.entries.Delete(key)
}

func (c *parseCache[K, V]) reset() {
	c.entries = collections.SyncMap[K, *parseCacheEntry[V]]{}
}

func (c *parseCache[K, V]) compact(w io.Writer) {
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "	Called count: %d\n", c.calledCount.Load())
	fmt.Fprintf(w, "	Reused count: %d\n", c.reusedCount.Load())
	fmt.Fprintf(w, "	Entries count: %d\n", c.entriesCount.Load())
	start := time.Now()
	c.entries.Range(func(key K, entry *parseCacheEntry[V]) bool {
		entry.mu.Lock()
		reuse := entry.reuse
		value := entry.value
		entry.mu.Unlock()
		var zero V
		if reuse == 0 || value == zero {
			c.delete(key)
			c.entriesCount.Add(-1)
		}
		return true
	})
	fmt.Fprintf(w, "	Parse cache compacted to %d entries in %v\n", c.entriesCount.Load(), time.Since(start))
	fmt.Fprintln(w, "")
}

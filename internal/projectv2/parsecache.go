package projectv2

import (
	"crypto/sha256"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type parseCacheKey struct {
	ast.SourceFileParseOptions
	scriptKind core.ScriptKind
}

func newParseCacheKey(
	options ast.SourceFileParseOptions,
	scriptKind core.ScriptKind,
) parseCacheKey {
	return parseCacheKey{
		SourceFileParseOptions: options,
		scriptKind:             scriptKind,
	}
}

type parseCacheEntry struct {
	mu         sync.Mutex
	sourceFile *ast.SourceFile
	hash       [sha256.Size]byte
	refCount   int
}

type parseCache struct {
	options tspath.ComparePathsOptions
	entries collections.SyncMap[parseCacheKey, *parseCacheEntry]
}

func (c *parseCache) Acquire(
	fh FileHandle,
	opts ast.SourceFileParseOptions,
	scriptKind core.ScriptKind,
) *ast.SourceFile {
	key := newParseCacheKey(opts, scriptKind)
	entry, loaded := c.loadOrStoreNewLockedEntry(key)
	defer entry.mu.Unlock()
	if !loaded || entry.hash != fh.Hash() {
		// Reparse the file if the hash has changed, or parse for the first time.
		entry.sourceFile = parser.ParseSourceFile(opts, fh.Content(), scriptKind)
		entry.hash = fh.Hash()
	}
	return entry.sourceFile
}

func (c *parseCache) Ref(file *ast.SourceFile) {
	key := newParseCacheKey(file.ParseOptions(), file.ScriptKind)
	if entry, ok := c.entries.Load(key); ok {
		entry.mu.Lock()
		entry.refCount++
		entry.mu.Unlock()
	} else {
		panic("parse cache entry not found")
	}
}

func (c *parseCache) Release(file *ast.SourceFile) {
	key := newParseCacheKey(file.ParseOptions(), file.ScriptKind)
	if entry, ok := c.entries.Load(key); ok {
		entry.mu.Lock()
		entry.refCount--
		remove := entry.refCount <= 0
		entry.mu.Unlock()
		if remove {
			c.entries.Delete(key)
		}
	}
}

// loadOrStoreNewLockedEntry loads an existing entry or creates a new one. The returned
// entry's mutex is locked and its refCount is incremented (or initialized to 1 in the
// case of a new entry).
func (c *parseCache) loadOrStoreNewLockedEntry(key parseCacheKey) (*parseCacheEntry, bool) {
	entry := &parseCacheEntry{refCount: 1}
	entry.mu.Lock()
	existing, loaded := c.entries.LoadOrStore(key, entry)
	if loaded {
		existing.mu.Lock()
		existing.refCount++
		return existing, true
	}
	return entry, false
}

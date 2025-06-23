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
	sourceFile *ast.SourceFile
	hash       [sha256.Size]byte
	refCount   int
	mu         sync.Mutex
}

type parseCache struct {
	options tspath.ComparePathsOptions
	entries collections.SyncMap[parseCacheKey, *parseCacheEntry]
}

func (c *parseCache) acquireDocument(
	fh fileHandle,
	opts ast.SourceFileParseOptions,
	scriptKind core.ScriptKind,
) *ast.SourceFile {
	key := newParseCacheKey(opts, scriptKind)
	entry, loaded := c.loadOrStoreNewEntry(key)
	if loaded {
		// Existing entry found, increment ref count and check hash
		entry.mu.Lock()
		entry.refCount++
		if entry.hash != fh.Hash() {
			// Reparse the file if the hash has changed
			entry.sourceFile = parser.ParseSourceFile(opts, fh.Content(), scriptKind)
			entry.hash = fh.Hash()
		}
		entry.mu.Unlock()
		return entry.sourceFile
	}

	// New entry created (still holding lock)
	entry.sourceFile = parser.ParseSourceFile(opts, fh.Content(), scriptKind)
	entry.hash = fh.Hash()
	entry.mu.Unlock()
	return entry.sourceFile
}

func (c *parseCache) releaseDocument(file *ast.SourceFile) {
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

func (c *parseCache) loadOrStoreNewEntry(key parseCacheKey) (*parseCacheEntry, bool) {
	entry := &parseCacheEntry{refCount: 1}
	entry.mu.Lock()
	existing, loaded := c.entries.LoadOrStore(key, entry)
	if loaded {
		return existing, true
	}
	return entry, false
}

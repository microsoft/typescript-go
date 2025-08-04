package execute

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// extendedConfigCache is a minimal implementation of tsoptions.ExtendedConfigCache.
// It is concurrency-safe, but stores cached entries permanently. This implementation
// should not be used for long-running processes where configuration changes over the
// course of multiple compilations.
type extendedConfigCache struct {
	mu sync.Mutex
	m  map[tspath.Path]*tsoptions.ExtendedConfigCacheEntry
}

var _ tsoptions.ExtendedConfigCache = (*extendedConfigCache)(nil)

// GetExtendedConfig implements tsoptions.ExtendedConfigCache.
func (e *extendedConfigCache) GetExtendedConfig(fileName string, path tspath.Path, parse func() *tsoptions.ExtendedConfigCacheEntry) *tsoptions.ExtendedConfigCacheEntry {
	e.mu.Lock()
	defer e.mu.Unlock()
	if entry, ok := e.m[path]; ok {
		return entry
	}
	entry := parse()
	if e.m == nil {
		e.m = make(map[tspath.Path]*tsoptions.ExtendedConfigCacheEntry)
	}
	e.m[path] = entry
	return entry
}

package vfswatch_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"github.com/microsoft/typescript-go/internal/vfs/vfswatch"
)

// countingFS wraps a vfs.FS and counts calls to GetAccessibleEntries.
type countingFS struct {
	vfs.FS
	n atomic.Int64
}

func (c *countingFS) GetAccessibleEntries(path string) vfs.Entries {
	c.n.Add(1)
	return c.FS.GetAccessibleEntries(path)
}

// TestHasChangesNoRedundantGetAccessibleEntries verifies that HasChangesFromWatchState
// calls GetAccessibleEntries exactly once per directory tracked via a recursive wildcard,
// not twice. Previously a redundant WalkDir pass in hasChanges doubled the calls.
//
// With /src as a recursive wildcard dir, snapshotPaths adds /src and /src/sub to
// watchState with ChildrenHash. A single HasChangesFromWatchState call should therefore
// call GetAccessibleEntries exactly 2 times — once for /src and once for /src/sub.
func TestHasChangesNoRedundantGetAccessibleEntries(t *testing.T) {
	t.Parallel()

	inner := vfstest.FromMap(map[string]string{
		"/src/a.ts":      "const a = 1;",
		"/src/b.ts":      "const b = 2;",
		"/src/sub/c.ts":  "const c = 3;",
		"/tsconfig.json": "{}",
	}, true)
	cfs := &countingFS{FS: inner}

	fw := vfswatch.NewFileWatcher(cfs, 10*time.Millisecond, true, func() {})
	fw.UpdateWatchState(
		[]string{"/src/a.ts", "/src/b.ts", "/src/sub/c.ts", "/tsconfig.json"},
		map[string]bool{"/src": true},
	)

	cfs.n.Store(0) // reset counter after baseline snapshot

	fw.HasChangesFromWatchState()

	// /src and /src/sub are the two dirs tracked with ChildrenHash.
	// Each should be hashed exactly once; the old code hashed each twice.
	if got := cfs.n.Load(); got != 2 {
		t.Errorf("GetAccessibleEntries called %d times, want 2 (once per tracked dir)", got)
	}
}

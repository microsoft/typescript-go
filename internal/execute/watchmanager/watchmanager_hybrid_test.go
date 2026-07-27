package watchmanager

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/fswatch"
)

// fakeWatch is a no-op io.Closer standing in for a registered watch.
type fakeWatch struct{ dir string }

func (fakeWatch) Close() error { return nil }

// fakeBackend records watched directories. A directory fails with failWith when
// failDir reports true for it (or, when failDir is nil, for every directory).
// WatchDirectories emulates the real all-or-nothing batch semantics: if any
// request would fail, the whole batch is rolled back and the error returned.
type fakeBackend struct {
	mu        sync.Mutex
	name      string
	failWith  error
	failDir   func(string) bool
	watched   []string
	recursive map[string]bool
}

func (b *fakeBackend) shouldFail(dir string) bool {
	if b.failWith == nil {
		return false
	}
	return b.failDir == nil || b.failDir(dir)
}

func (b *fakeBackend) record(dir string, recursive bool) io.Closer {
	b.watched = append(b.watched, dir)
	if b.recursive == nil {
		b.recursive = map[string]bool{}
	}
	b.recursive[dir] = recursive
	return fakeWatch{dir: dir}
}

func (b *fakeBackend) WatchDirectory(dir string, fn fswatch.WatchCallback, recursive bool, ignore func(string) bool) (io.Closer, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.shouldFail(dir) {
		return nil, fmt.Errorf("%s: %w", b.name, b.failWith)
	}
	return b.record(dir, recursive), nil
}

func (b *fakeBackend) WatchDirectories(requests []WatchDirectoryRequest) ([]io.Closer, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range requests {
		if b.shouldFail(r.Dir) {
			return nil, fmt.Errorf("%s: %w", b.name, b.failWith)
		}
	}
	closers := make([]io.Closer, len(requests))
	for i, r := range requests {
		closers[i] = b.record(r.Dir, r.Recursive)
	}
	return closers, nil
}

func (b *fakeBackend) watchedDirs() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := append([]string(nil), b.watched...)
	sort.Strings(out)
	return out
}

func newHybridManager(primary, secondary WatchBackend) *WatchManager {
	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	if secondary != nil {
		wm.secondaryFactory = func() WatchBackend { return secondary }
	}
	return wm
}

func reconcile(t *testing.T, wm *WatchManager, desired map[string]bool) error {
	t.Helper()
	wm.Lock()
	defer wm.Unlock()
	return wm.ReconcileWatches(desired)
}

// TestReconcileHybridRoutesPerDirectory verifies the core hybrid behavior: when
// the primary (fanotify) backend can watch some directories but not others
// (cross-mount), each directory is routed to the backend that can watch it —
// the supported ones stay on the primary, only the unsupported ones use the
// secondary (inotify). Addresses the cross-mount concern on PR #4661.
func TestReconcileHybridRoutesPerDirectory(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{
		name:     "fanotify",
		failWith: fswatch.ErrFilesystemUnsupported,
		failDir:  func(dir string) bool { return strings.HasPrefix(dir, "/mnt/fuse") },
	}
	secondary := &fakeBackend{name: "inotify"}
	wm := newHybridManager(primary, secondary)

	desired := map[string]bool{
		"/project":         true,  // ext4 -> primary
		"/project/src":     false, // ext4 -> primary
		"/mnt/fuse/deps":   true,  // FUSE -> secondary
		"/mnt/fuse/deps/a": false, // FUSE -> secondary
	}
	if err := reconcile(t, wm, desired); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	if got, want := primary.watchedDirs(), []string{"/project", "/project/src"}; !equalStrings(got, want) {
		t.Errorf("primary watched %v, want %v", got, want)
	}
	if got, want := secondary.watchedDirs(), []string{"/mnt/fuse/deps", "/mnt/fuse/deps/a"}; !equalStrings(got, want) {
		t.Errorf("secondary watched %v, want %v", got, want)
	}
	// Every desired directory ends up watched, and fallback flags are correct.
	for dir := range desired {
		wd, ok := wm.watchedDirs[dir]
		if !ok {
			t.Errorf("%s not watched", dir)
			continue
		}
		wantFallback := strings.HasPrefix(dir, "/mnt/fuse")
		if wd.usesFallback != wantFallback {
			t.Errorf("%s usesFallback=%v, want %v", dir, wd.usesFallback, wantFallback)
		}
	}
	// The primary backend is never replaced.
	if wm.Backend() != primary {
		t.Errorf("primary backend should remain in place")
	}
}

// TestReconcileHybridAllUnsupported covers the pure Docker/NTFS case where the
// primary can watch nothing: every directory routes to the secondary, but the
// primary backend is still not replaced (so later dirs on a supported mount
// could still use it).
func TestReconcileHybridAllUnsupported(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}
	secondary := &fakeBackend{name: "inotify"}
	wm := newHybridManager(primary, secondary)

	desired := map[string]bool{"/w": false, "/w/src": true}
	if err := reconcile(t, wm, desired); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if got, want := secondary.watchedDirs(), []string{"/w", "/w/src"}; !equalStrings(got, want) {
		t.Errorf("secondary watched %v, want %v", got, want)
	}
	if len(primary.watchedDirs()) != 0 {
		t.Errorf("primary should have watched nothing, got %v", primary.watchedDirs())
	}
	if secondary.recursive["/w/src"] != true || secondary.recursive["/w"] != false {
		t.Errorf("recursive flags not preserved through routing: %v", secondary.recursive)
	}
	if wm.Backend() != primary {
		t.Errorf("primary backend should remain in place")
	}
}

// TestReconcileHybridNoSecondary verifies that without a secondary backend, an
// unsupported-filesystem error is surfaced rather than silently swallowed.
func TestReconcileHybridNoSecondary(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}
	wm := newHybridManager(primary, nil)

	if err := reconcile(t, wm, map[string]bool{"/w": false}); err == nil {
		t.Fatal("expected error when no secondary backend is available")
	}
}

// TestReconcileHybridSecondaryFactoryNil verifies that a secondary factory that
// yields no backend (e.g. inotify unavailable) also surfaces the error.
func TestReconcileHybridSecondaryFactoryNil(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}
	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	wm.secondaryFactory = func() WatchBackend { return nil }

	if err := reconcile(t, wm, map[string]bool{"/w": false}); err == nil {
		t.Fatal("expected error when secondary factory returns nil")
	}
}

// TestReconcileHybridUnrelatedError verifies that a non-filesystem error does
// not trigger secondary routing.
func TestReconcileHybridUnrelatedError(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrUnavailable}
	secondary := &fakeBackend{name: "inotify"}
	wm := newHybridManager(primary, secondary)

	if err := reconcile(t, wm, map[string]bool{"/w": false}); err == nil {
		t.Fatal("expected error to be surfaced")
	}
	if len(secondary.watchedDirs()) != 0 {
		t.Errorf("secondary should not have been used for a non-filesystem error")
	}
}

// TestReconcileHybridHappyPath verifies that when the primary can watch
// everything, the secondary is never instantiated or used.
func TestReconcileHybridHappyPath(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify"} // never fails
	secondary := &fakeBackend{name: "inotify"}
	wm := newHybridManager(primary, secondary)

	if err := reconcile(t, wm, map[string]bool{"/w": false, "/w/src": true}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if got, want := primary.watchedDirs(), []string{"/w", "/w/src"}; !equalStrings(got, want) {
		t.Errorf("primary watched %v, want %v", got, want)
	}
	if len(secondary.watchedDirs()) != 0 {
		t.Errorf("secondary should be unused on the happy path")
	}
	if wm.secondary != nil {
		t.Errorf("secondary backend should not be instantiated on the happy path")
	}
}

// TestEnsureDefaultBackendWiresFanotifyFallback verifies the production wiring:
// EnsureDefaultBackend must set secondaryFactory exactly when the auto-selected
// backend is fanotify. Environment-aware so it is valid regardless of which
// backend fswatch.Default() picks on the test host.
func TestEnsureDefaultBackendWiresFanotifyFallback(t *testing.T) {
	t.Parallel()

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.EnsureDefaultBackend()

	fsb, ok := wm.backend.(*FSWatchBackend)
	if !ok {
		t.Fatalf("expected default backend to be *FSWatchBackend, got %T", wm.backend)
	}
	if fsb.Inner.Name() == "fanotify" {
		if wm.secondaryFactory == nil {
			t.Fatal("fanotify default must set a secondary (inotify) factory")
		}
	} else if wm.secondaryFactory != nil {
		t.Fatalf("non-fanotify default (%s) must not set a secondary factory", fsb.Inner.Name())
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

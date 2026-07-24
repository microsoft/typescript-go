package watchmanager

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/fswatch"
)

// fakeWatch is a no-op io.Closer standing in for a registered watch.
type fakeWatch struct{ dir string }

func (fakeWatch) Close() error { return nil }

// fakeBackend records watched directories. When failWith is non-nil it fails
// every WatchDirectories call with that error, simulating a backend that
// cannot operate on the target filesystem.
type fakeBackend struct {
	mu        sync.Mutex
	name      string
	failWith  error
	watched   []string
	recursive map[string]bool
	calls     int
}

func (b *fakeBackend) WatchDirectory(dir string, fn fswatch.WatchCallback, recursive bool, ignore func(string) bool) (io.Closer, error) {
	closers, err := b.WatchDirectories([]WatchDirectoryRequest{{Dir: dir, Callback: fn, Recursive: recursive, Ignore: ignore}})
	if err != nil {
		return nil, err
	}
	return closers[0], nil
}

func (b *fakeBackend) WatchDirectories(requests []WatchDirectoryRequest) ([]io.Closer, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.calls++
	if b.failWith != nil {
		return nil, fmt.Errorf("%s: %w", b.name, b.failWith)
	}
	closers := make([]io.Closer, len(requests))
	for i, r := range requests {
		b.watched = append(b.watched, r.Dir)
		if b.recursive == nil {
			b.recursive = map[string]bool{}
		}
		b.recursive[r.Dir] = r.Recursive
		closers[i] = fakeWatch{dir: r.Dir}
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

// TestReconcileWatchesFallsBackOnUnsupportedFilesystem verifies that when the
// primary backend reports fswatch.ErrFilesystemUnsupported, the manager swaps
// to the fallback backend and re-arms every desired directory. This models the
// Docker bind-mount case where fanotify's name_to_handle_at is unsupported and
// tsc --watch must fall back to inotify (issue #63646).
func TestReconcileWatchesFallsBackOnUnsupportedFilesystem(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}
	fallback := &fakeBackend{name: "inotify"}

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	wm.fallbackBackend = func() WatchBackend { return fallback }

	desired := map[string]bool{"/workspace/src": true, "/workspace": false}
	wm.Lock()
	err := wm.ReconcileWatches(desired)
	wm.Unlock()
	if err != nil {
		t.Fatalf("ReconcileWatches returned error after fallback: %v", err)
	}

	if wm.Backend() != fallback {
		t.Fatalf("expected backend to be swapped to the fallback backend")
	}
	got := fallback.watchedDirs()
	want := []string{"/workspace", "/workspace/src"}
	if len(got) != len(want) {
		t.Fatalf("fallback watched %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("fallback watched %v, want %v", got, want)
		}
	}
}

// TestReconcileWatchesFallbackOnlyOnce verifies that if the fallback backend
// also reports an unsupported filesystem, the manager does not loop and instead
// surfaces the error.
func TestReconcileWatchesFallbackOnlyOnce(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}
	fallback := &fakeBackend{name: "inotify", failWith: fswatch.ErrFilesystemUnsupported}

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	wm.fallbackBackend = func() WatchBackend { return fallback }

	wm.Lock()
	err := wm.ReconcileWatches(map[string]bool{"/workspace": false})
	wm.Unlock()
	if err == nil {
		t.Fatalf("expected error when both backends are unsupported")
	}
	if primary.calls != 1 || fallback.calls != 1 {
		t.Fatalf("expected exactly one call to each backend, got primary=%d fallback=%d", primary.calls, fallback.calls)
	}
}

// TestReconcileWatchesNoFallbackForOtherErrors verifies that unrelated watch
// errors do not trigger a backend swap.
func TestReconcileWatchesNoFallbackForOtherErrors(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrUnavailable}
	fallback := &fakeBackend{name: "inotify"}

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	wm.fallbackBackend = func() WatchBackend { return fallback }

	wm.Lock()
	err := wm.ReconcileWatches(map[string]bool{"/workspace": false})
	wm.Unlock()
	if err == nil {
		t.Fatalf("expected error to be surfaced")
	}
	if wm.Backend() != primary {
		t.Fatalf("backend should not have been swapped for a non-filesystem error")
	}
	if fallback.calls != 0 {
		t.Fatalf("fallback backend should not have been used, got %d calls", fallback.calls)
	}
}

// TestReconcileWatchesFallbackPreservesRecursive verifies the recursive flag of
// each desired directory is preserved when the watches are re-armed on the
// fallback backend.
func TestReconcileWatchesFallbackPreservesRecursive(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}
	fallback := &fakeBackend{name: "inotify"}

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	wm.fallbackBackend = func() WatchBackend { return fallback }

	desired := map[string]bool{"/w/src": true, "/w": false}
	wm.Lock()
	err := wm.ReconcileWatches(desired)
	wm.Unlock()
	if err != nil {
		t.Fatalf("ReconcileWatches returned error: %v", err)
	}

	fallback.mu.Lock()
	defer fallback.mu.Unlock()
	if got := fallback.recursive["/w/src"]; got != true {
		t.Errorf("/w/src recursive = %v, want true", got)
	}
	if got := fallback.recursive["/w"]; got != false {
		t.Errorf("/w recursive = %v, want false", got)
	}
}

// TestReconcileWatchesAfterFallbackUsesFallbackBackend verifies that once a
// fallback has occurred, later reconciles arm new watches on the fallback
// backend and do not re-invoke the primary or re-trigger fallback.
func TestReconcileWatchesAfterFallbackUsesFallbackBackend(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}
	fallback := &fakeBackend{name: "inotify"}

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	wm.fallbackBackend = func() WatchBackend { return fallback }

	wm.Lock()
	if err := wm.ReconcileWatches(map[string]bool{"/w/a": false}); err != nil {
		t.Fatalf("first reconcile: %v", err)
	}
	// Second reconcile adds another directory.
	if err := wm.ReconcileWatches(map[string]bool{"/w/a": false, "/w/b": false}); err != nil {
		t.Fatalf("second reconcile: %v", err)
	}
	wm.Unlock()

	primaryCalls := primary.calls
	if primaryCalls != 1 {
		t.Errorf("primary backend called %d times, want 1 (only the initial failing attempt)", primaryCalls)
	}
	if wm.Backend() != fallback {
		t.Errorf("backend should still be the fallback")
	}
	got := fallback.watchedDirs()
	want := []string{"/w/a", "/w/b"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("fallback watched %v, want %v", got, want)
	}
}

// TestReconcileWatchesFallbackUnavailable verifies that when no fallback backend
// is available (factory returns nil), the manager surfaces the error, keeps the
// original backend, and does not retry indefinitely.
func TestReconcileWatchesFallbackUnavailable(t *testing.T) {
	t.Parallel()

	primary := &fakeBackend{name: "fanotify", failWith: fswatch.ErrFilesystemUnsupported}

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.backend = primary
	wm.fallbackBackend = func() WatchBackend { return nil } // e.g. inotify unavailable

	wm.Lock()
	err := wm.ReconcileWatches(map[string]bool{"/w": false})
	wm.Unlock()
	if err == nil {
		t.Fatalf("expected error to be surfaced when fallback is unavailable")
	}
	if wm.Backend() != primary {
		t.Errorf("backend should be unchanged when fallback is unavailable")
	}
	if primary.calls != 1 {
		t.Errorf("primary called %d times, want 1", primary.calls)
	}
}

// TestEnsureDefaultBackendWiresFanotifyFallback verifies the production wiring:
// EnsureDefaultBackend must populate fallbackBackend exactly when the
// auto-selected backend is fanotify (and leave it nil otherwise). Without this,
// the whole fallback path would silently never arm in production even though the
// isolated fallback tests pass. Environment-aware so it is valid regardless of
// which backend fswatch.Default() selects on the test host.
func TestEnsureDefaultBackendWiresFanotifyFallback(t *testing.T) {
	t.Parallel()

	wm := NewWatchManager(io.Discard, func(string) bool { return true })
	wm.EnsureDefaultBackend()

	fsb, ok := wm.backend.(*FSWatchBackend)
	if !ok {
		t.Fatalf("expected default backend to be *FSWatchBackend, got %T", wm.backend)
	}
	if fsb.Inner.Name() == "fanotify" {
		if wm.fallbackBackend == nil {
			t.Fatal("fanotify default must set an inotify fallbackBackend")
		}
	} else if wm.fallbackBackend != nil {
		t.Fatalf("non-fanotify default (%s) must not set a fallbackBackend", fsb.Inner.Name())
	}
}

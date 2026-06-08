package lspwatcher

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/fswatch"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func waitFor(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", msg)
}

func TestWatcher_CreateChangeDelete(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	var (
		mu      sync.Mutex
		batches [][]*lsproto.FileEvent
	)
	w := New(bundled.WrapFS(osvfs.FS()), func(changes []*lsproto.FileEvent) {
		mu.Lock()
		defer mu.Unlock()
		batches = append(batches, changes)
	}, logging.NewLogger(os.Stderr))
	t.Cleanup(w.Close)

	pattern := tspath.NormalizeSlashes(dir) + "/**/*"
	kind := lsproto.WatchKindCreate | lsproto.WatchKindChange | lsproto.WatchKindDelete
	if err := w.WatchFiles("test", []*lsproto.FileSystemWatcher{{
		GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &pattern},
		Kind:        &kind,
	}}); err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)

	file := filepath.Join(dir, "a.ts")
	if err := os.WriteFile(file, []byte("export {}"), 0o644); err != nil {
		t.Fatal(err)
	}

	collected := func() []*lsproto.FileEvent {
		mu.Lock()
		defer mu.Unlock()
		var all []*lsproto.FileEvent
		for _, b := range batches {
			all = append(all, b...)
		}
		return all
	}

	waitFor(t, func() bool {
		for _, e := range collected() {
			if e.Type == lsproto.FileChangeTypeChanged {
				return true
			}
		}
		return false
	}, "update event")

	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		for _, e := range collected() {
			if e.Type == lsproto.FileChangeTypeDeleted {
				return true
			}
		}
		return false
	}, "delete event")

	if err := w.UnwatchFiles("test"); err != nil {
		t.Fatal(err)
	}
}

func TestWatcher_KindFilter(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dirNorm := tspath.NormalizeSlashes(dir)

	var (
		mu  sync.Mutex
		got []*lsproto.FileEvent
	)
	backend := newFakeBackend()
	w := newWithBackend(bundled.WrapFS(osvfs.FS()), backend, func(changes []*lsproto.FileEvent) {
		mu.Lock()
		defer mu.Unlock()
		got = append(got, changes...)
	}, logging.NewLogger(os.Stderr))
	t.Cleanup(w.Close)

	pattern := dirNorm + "/**/*"
	kind := lsproto.WatchKindDelete
	if err := w.WatchFiles("test", []*lsproto.FileSystemWatcher{{
		GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &pattern},
		Kind:        &kind,
	}}); err != nil {
		t.Fatal(err)
	}
	backend.emitAll([]fswatch.Event{
		{Kind: fswatch.EventUpdate, Path: filepath.FromSlash(filepath.Join(dirNorm, "x.ts"))},
		{Kind: fswatch.EventDelete, Path: filepath.FromSlash(filepath.Join(dirNorm, "x.ts"))},
	}, nil)

	waitFor(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		for _, e := range got {
			if e.Type == lsproto.FileChangeTypeDeleted {
				return true
			}
		}
		return false
	}, "delete event")

	mu.Lock()
	for _, e := range got {
		if e.Type != lsproto.FileChangeTypeDeleted {
			t.Errorf("unexpected non-delete event: %+v", e)
		}
	}
	mu.Unlock()
}

func TestRootFromGlob(t *testing.T) {
	t.Parallel()
	cases := []struct {
		pattern string
		want    string
	}{
		{"/abs/path/**/*", "/abs/path"},
		{"/abs/path/", "/abs/path"},
		{"/abs/path/?.ts", "/abs/path"},
		{"/abs/path/{a,b}/*", "/abs/path"},
	}
	for _, c := range cases {
		if got := rootFromGlob(c.pattern); got != c.want {
			t.Errorf("rootFromGlob(%q) = %q, want %q", c.pattern, got, c.want)
		}
	}
}

type fakeBackend struct {
	mu       sync.Mutex
	byDir    map[string]fswatch.WatchCallback
	closed   map[string]int
	optCount map[string]int
}

func newFakeBackend() *fakeBackend {
	return &fakeBackend{
		byDir:    make(map[string]fswatch.WatchCallback),
		closed:   make(map[string]int),
		optCount: make(map[string]int),
	}
}

func (f *fakeBackend) WatchDirectory(dir string, fn fswatch.WatchCallback, opts ...fswatch.WatchOption) (watch, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byDir[dir] = fn
	f.optCount[dir] = len(opts)
	return fakeWatch{closeFn: func() error {
		f.mu.Lock()
		defer f.mu.Unlock()
		delete(f.byDir, dir)
		f.closed[dir]++
		return nil
	}}, nil
}

func (f *fakeBackend) emit(dir string, events []fswatch.Event, err error) {
	f.mu.Lock()
	cb := f.byDir[dir]
	f.mu.Unlock()
	if cb != nil {
		cb(events, err)
	}
}

func (f *fakeBackend) emitAll(events []fswatch.Event, err error) {
	f.mu.Lock()
	cbs := make([]fswatch.WatchCallback, 0, len(f.byDir))
	for _, cb := range f.byDir {
		cbs = append(cbs, cb)
	}
	f.mu.Unlock()
	for _, cb := range cbs {
		cb(events, err)
	}
}

type fakeWatch struct{ closeFn func() error }

func (w fakeWatch) Close() error { return w.closeFn() }

func TestWatcher_BookkeepingAndOverflow(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dirNorm := tspath.NormalizeSlashes(dir)
	pattern := dirNorm + "/**/*"

	fs := bundled.WrapFS(osvfs.FS())
	backend := newFakeBackend()
	var (
		mu  sync.Mutex
		got []*lsproto.FileEvent
	)
	w := newWithBackend(fs, backend, func(changes []*lsproto.FileEvent) {
		mu.Lock()
		defer mu.Unlock()
		got = append(got, changes...)
	}, logging.NewLogger(os.Stderr))

	if err := w.WatchFiles("id", []*lsproto.FileSystemWatcher{{
		GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &pattern},
	}}); err != nil {
		t.Fatal(err)
	}
	if err := w.WatchFiles("id", []*lsproto.FileSystemWatcher{{
		GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &pattern},
	}}); err == nil {
		t.Fatal("expected duplicate-id error")
	}

	backend.emitAll([]fswatch.Event{
		{Kind: fswatch.EventUpdate, Path: filepath.FromSlash(filepath.Join(dirNorm, "a.ts"))},
	}, fswatch.ErrOverflow)
	waitFor(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(got) > 0
	}, "events after overflow")

	if err := w.UnwatchFiles("missing"); err == nil {
		t.Fatal("expected unknown-id error")
	}
	if err := w.UnwatchFiles("id"); err != nil {
		t.Fatal(err)
	}
	if err := w.WatchFiles("id2", nil); err != nil {
		t.Fatal(err)
	}
	w.Close()
	if err := w.WatchFiles("id3", nil); err == nil {
		t.Fatal("expected closed error")
	}
}

func TestWatcher_NonRecursiveGlobIsNotRecursive(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dirNorm := tspath.NormalizeSlashes(dir)
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	subNorm := tspath.NormalizeSlashes(filepath.Join(dir, "sub"))

	fs := bundled.WrapFS(osvfs.FS())
	backend := newFakeBackend()
	w := newWithBackend(fs, backend, func([]*lsproto.FileEvent) {}, logging.NewLogger(os.Stderr))
	t.Cleanup(w.Close)

	recursive := dirNorm + "/**/*"
	nonRecursive := subNorm + "/*"
	if err := w.WatchFiles("id", []*lsproto.FileSystemWatcher{
		{GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &recursive}},
		{GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &nonRecursive}},
	}); err != nil {
		t.Fatal(err)
	}

	realDir := tspath.NormalizeSlashes(fs.Realpath(dirNorm))
	realSub := tspath.NormalizeSlashes(fs.Realpath(subNorm))

	backend.mu.Lock()
	defer backend.mu.Unlock()
	if got := backend.optCount[realDir]; got != 1 {
		t.Errorf("recursive glob %q: expected 1 watch option (WithRecursive), got %d", recursive, got)
	}
	if got := backend.optCount[realSub]; got != 0 {
		t.Errorf("non-recursive glob %q: expected 0 watch options, got %d", nonRecursive, got)
	}
}

func TestWatcher_MissingDirectoryIsSkipped(t *testing.T) {
	t.Parallel()

	fs := bundled.WrapFS(osvfs.FS())
	backend := newFakeBackend()
	w := newWithBackend(fs, backend, func([]*lsproto.FileEvent) {}, logging.NewLogger(os.Stderr))

	pattern := tspath.NormalizeSlashes(filepath.Join(t.TempDir(), "does-not-exist")) + "/**/*"
	if err := w.WatchFiles("id", []*lsproto.FileSystemWatcher{{
		GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &pattern},
	}}); err != nil {
		t.Fatal(err)
	}
	if err := w.UnwatchFiles("id"); err != nil {
		t.Fatal(err)
	}
}

func TestWatcher_WatchTerminatedDoesNotDropEvents(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dirNorm := tspath.NormalizeSlashes(dir)
	fs := bundled.WrapFS(osvfs.FS())
	backend := newFakeBackend()
	var got []*lsproto.FileEvent
	var mu sync.Mutex
	w := newWithBackend(fs, backend, func(changes []*lsproto.FileEvent) {
		mu.Lock()
		got = append(got, changes...)
		mu.Unlock()
	}, logging.NewLogger(os.Stderr))

	pattern := dirNorm + "/**/*"
	if err := w.WatchFiles("id", []*lsproto.FileSystemWatcher{{
		GlobPattern: lsproto.PatternOrRelativePattern{Pattern: &pattern},
	}}); err != nil {
		t.Fatal(err)
	}

	backend.emitAll([]fswatch.Event{
		{Kind: fswatch.EventUpdate, Path: filepath.FromSlash(filepath.Join(dirNorm, "b.ts"))},
	}, errors.Join(fswatch.ErrWatchTerminated, errors.New("simulated")))

	waitFor(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(got) > 0
	}, "events with watch-terminated error")
}

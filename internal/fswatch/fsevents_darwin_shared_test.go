//go:build darwin && (amd64 || arm64)

package fswatch

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestFSEventsWatcher(impl **fsEventsBackend) Watcher {
	return &watcher{
		name: "fsevents",
		factory: func() watcherImpl {
			*impl = newFSEventsBackend()
			return *impl
		},
	}
}

func TestFSEventsSharedStreamAcrossWatches(t *testing.T) {
	t.Parallel()

	var impl *fsEventsBackend
	watcherImpl := newTestFSEventsWatcher(&impl)
	root := newTmpDir(t)

	var subs []Watch
	for i := range 5 {
		dir := filepath.Join(root, fmt.Sprintf("dir%d", i))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		sub, err := watcherImpl.WatchDirectory(dir, func([]Event, error) {})
		if err != nil {
			t.Fatal(err)
		}
		subs = append(subs, sub)
	}
	t.Cleanup(func() {
		for _, sub := range subs {
			_ = sub.Close()
		}
	})

	impl.mu.Lock()
	streamCount := len(impl.streams)
	watchCount := len(impl.watches)
	impl.mu.Unlock()
	if streamCount != 1 {
		t.Fatalf("expected one shared FSEvents stream, got %d", streamCount)
	}
	if watchCount != len(subs) {
		t.Fatalf("expected %d logical watches, got %d", len(subs), watchCount)
	}
}

func TestFSEventsSharedStreamRoutesEvents(t *testing.T) {
	t.Parallel()

	var impl *fsEventsBackend
	watcherImpl := newTestFSEventsWatcher(&impl)
	root := newTmpDir(t)
	dirA := filepath.Join(root, "a")
	dirB := filepath.Join(root, "b")
	if err := os.MkdirAll(dirA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dirB, 0o755); err != nil {
		t.Fatal(err)
	}

	time.Sleep(preSubscribeSleep(watcherImpl))
	recA := newRecorder(t)
	recA.watcher = watcherImpl
	subA, err := watcherImpl.WatchDirectory(dirA, recA.callback)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = subA.Close() })

	recB := newRecorder(t)
	recB.watcher = watcherImpl
	subB, err := watcherImpl.WatchDirectory(dirB, recB.callback)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = subB.Close() })
	time.Sleep(settleSleep(watcherImpl))

	fileA := filepath.Join(dirA, "file.ts")
	if err := os.WriteFile(fileA, []byte("export {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	expectContains(t, recA, EventUpdate, fileA)
	assertNoEventsForPath(t, recB.drainQuiet(500*time.Millisecond), fileA, "sibling watch saw event")

	fileB := filepath.Join(dirB, "file.ts")
	if err := os.WriteFile(fileB, []byte("export {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	expectContains(t, recB, EventUpdate, fileB)
	assertNoEventsForPath(t, recA.drainQuiet(500*time.Millisecond), fileB, "sibling watch saw event")
}

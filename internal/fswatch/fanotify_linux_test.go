//go:build linux

package fswatch

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestLinuxFanotifyShutdownBeforeStart(t *testing.T) {
	t.Parallel()
	newFanotifyBackend().shutdown()
}

func TestLinuxFanotifyBackendSelection(t *testing.T) {
	t.Parallel()
	if !fanotifyAvailable() {
		t.Skip("fanotify not available")
	}
	impl, err := fanotifyWatcher.getImpl()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := impl.(*fanotifyBackend); !ok {
		t.Fatalf("fanotify watcher = %T, want *fanotifyBackend", impl)
	}
}

func TestLinuxFanotifySubscribeCleansUpAfterMarkFailure(t *testing.T) {
	t.Parallel()
	dir := newTmpDir(t)
	w := newDirectWatcher(t, dir)
	b := newFanotifyBackend()

	err := b.subscribe(w)
	var werr *dirWatchError
	if !errors.As(err, &werr) {
		t.Fatalf("subscribe error = %v, want *dirWatchError", err)
	}
	if werr.dirWatch != w {
		t.Fatalf("dirWatchError dirWatch = %p, want %p", werr.dirWatch, w)
	}
	if len(b.subscriptions) != 0 {
		t.Fatalf("subscriptions not cleaned up: %d remaining", len(b.subscriptions))
	}
}

func TestLinuxFanotifyParseDfidNameRoundTrip(t *testing.T) {
	t.Parallel()
	dir := newTmpDir(t)
	handle, _, err := unix.NameToHandleAt(unix.AT_FDCWD, dir, 0)
	if err != nil {
		t.Skipf("NameToHandleAt not supported: %v", err)
	}
	var st unix.Statfs_t
	if err = unix.Statfs(dir, &st); err != nil {
		t.Fatal(err)
	}
	key := makeFanotifyHandleKey(st.Fsid.Val, handle.Type(), handle.Bytes())
	if key.handle == "" {
		t.Fatal("empty handle bytes")
	}
	handle2, _, err := unix.NameToHandleAt(unix.AT_FDCWD, dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	key2 := makeFanotifyHandleKey(st.Fsid.Val, handle2.Type(), handle2.Bytes())
	if key != key2 {
		t.Fatalf("handle keys differ for same path:\n  1: %+v\n  2: %+v", key, key2)
	}
}

func TestFanotifyNoRenameFallback(t *testing.T) { //nolint:tparallel // subtests share fanotifyTestNoRename global toggle
	t.Parallel()
	if !fanotifyAvailable() {
		t.Skip("fanotify not available")
	}

	fanotifyTestNoRename.Store(true)
	t.Cleanup(func() {
		fanotifyTestNoRename.Store(false)
	})

	t.Run("FileRename", func(t *testing.T) { //nolint:paralleltest // parent uses shared global toggle
		dir := newTmpDir(t)
		f1 := subPath(dir)
		f2 := subPath(dir)
		if err := os.WriteFile(f1, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		r, _ := subscribeFor(t, dir, Fanotify())
		if err := os.Rename(f1, f2); err != nil {
			t.Fatal(err)
		}
		got := r.gather(r.deadline(), 100*time.Millisecond)
		assertEventSet(t, got, []wantEvent{
			{EventDelete, f1},
			{EventUpdate, f2},
		})
	})

	t.Run("FileCreate", func(t *testing.T) { //nolint:paralleltest // parent uses shared global toggle
		dir := newTmpDir(t)
		r, _ := subscribeFor(t, dir, Fanotify())
		f := subPath(dir)
		if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
			t.Fatal(err)
		}
		got := r.next(r.deadline())
		assertEventSequence(t, got, []wantEvent{{EventUpdate, f}})
	})

	// On the FAN_MOVED_FROM fallback path, moving a nested directory
	// subtree out of the watched root must drop every descendant
	// subscription, not just the renamed directory's own entry. With
	// the bug, subs for inner watched directories kept their old
	// (stale) s.path; subsequent modifications at the new location
	// would be reported against the no-longer-existing original path.
	// FAN_RENAME's handleRenameEvent already had the prefix cleanup;
	// this regression test pins the fallback path to the same shape.
	t.Run("RenameDirOutDropsDescendants", func(t *testing.T) { //nolint:paralleltest // parent uses shared global toggle
		watched := newTmpDir(t)
		outside := newTmpDir(t)

		sub := filepath.Join(watched, "sub")
		inner := filepath.Join(sub, "inner")
		if err := os.MkdirAll(inner, 0o755); err != nil {
			t.Fatal(err)
		}
		leaf := filepath.Join(inner, "leaf.txt")
		if err := os.WriteFile(leaf, []byte("v1"), 0o644); err != nil {
			t.Fatal(err)
		}

		r, _ := subscribeForOpts(t, watched, Fanotify(), WithRecursive())

		dest := filepath.Join(outside, "moved")
		if err := os.Rename(sub, dest); err != nil {
			t.Fatal(err)
		}
		_ = r.gatherUntilQuiet(r.deadline(), 500*time.Millisecond)

		// Modify the file at its new location. The kernel mark on
		// inner's inode is still active; without descendant cleanup
		// the modify event would surface against watched/sub/inner/leaf.txt.
		movedLeaf := filepath.Join(dest, "inner", "leaf.txt")
		if err := os.WriteFile(movedLeaf, []byte("v2-longer"), 0o644); err != nil {
			t.Fatal(err)
		}
		// Also create a new file in the moved tree to provoke FAN_CREATE
		// on inner's mark.
		if err := os.WriteFile(filepath.Join(dest, "inner", "new.txt"), []byte("n"), 0o644); err != nil {
			t.Fatal(err)
		}

		extra := r.drainQuiet(800 * time.Millisecond)
		oldPrefix := sub + string(filepath.Separator)
		for _, e := range extra {
			if e.Path == sub || strings.HasPrefix(e.Path, oldPrefix) {
				t.Fatalf("stale event for moved-out path %s: %+v\nall extras: %v",
					e.Path, e, toWantEvents(extra))
			}
		}
	})
}

func TestFanotifyCrossWatcherSameFs(t *testing.T) {
	t.Parallel()
	if !fanotifyAvailable() {
		t.Skip("fanotify not available")
	}

	t.Run("Modify", func(t *testing.T) {
		t.Parallel()
		dirA, dirB := newTmpDir(t), newTmpDir(t)
		pathA := filepath.Join(dirA, "child")
		pathB := filepath.Join(dirB, "child")
		for _, p := range []string{pathA, pathB} {
			if err := os.WriteFile(p, []byte("initial"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		rA, _ := subscribeFor(t, dirA, Fanotify())
		rB, _ := subscribeFor(t, dirB, Fanotify())

		if err := os.WriteFile(pathA, []byte("changed"), 0o644); err != nil {
			t.Fatal(err)
		}
		gotA := rA.gather(rA.deadline(), 200*time.Millisecond)
		assertEventSet(t, gotA, []wantEvent{{EventUpdate, pathA}})

		if gotB := rB.drainQuiet(200 * time.Millisecond); len(gotB) != 0 {
			t.Fatalf("watcher B got phantom events: %v", toWantEvents(gotB))
		}
	})
}

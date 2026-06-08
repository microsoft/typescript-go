package tsctests

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/fswatch"
	"github.com/microsoft/typescript-go/internal/testutil/fsbaselineutil"
)

// MockWatchBackend implements execute.WatchBackend for testing. It
// records all WatchDirectory/WatchFile calls so tests can verify that
// the correct watches are registered.  Events can be delivered through
// SendEvents, which routes them only through watches whose paths
// match, enforcing that tests fail if the wrong watches are set up.
type MockWatchBackend struct {
	mu              sync.Mutex
	Dirs            map[string]*MockWatch
	Files           map[string]*MockWatch
	DirectoryExists func(string) bool // if set, WatchDirectory fails for non-existent dirs
}

var _ execute.WatchBackend = (*MockWatchBackend)(nil)

// NewMockWatchBackend creates a ready-to-use mock backend.
func NewMockWatchBackend() *MockWatchBackend {
	return &MockWatchBackend{
		Dirs:  make(map[string]*MockWatch),
		Files: make(map[string]*MockWatch),
	}
}

// HasWatches reports whether any watches have been registered.
func (m *MockWatchBackend) HasWatches() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Dirs) > 0 || len(m.Files) > 0
}

// MockWatch records a single registered watch.
type MockWatch struct {
	Path      string
	Callback  fswatch.WatchCallback
	Recursive bool
	Ignore    func(string) bool
	Closed    bool
}

func (w *MockWatch) Close() error {
	w.Closed = true
	return nil
}

func (m *MockWatchBackend) WatchDirectory(dir string, fn fswatch.WatchCallback, recursive bool, ignore func(string) bool) (io.Closer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.DirectoryExists != nil && !m.DirectoryExists(dir) {
		return nil, fmt.Errorf("directory does not exist: %s", dir)
	}
	w := &MockWatch{Path: dir, Callback: fn, Recursive: recursive, Ignore: ignore}
	m.Dirs[dir] = w
	return w, nil
}

func (m *MockWatchBackend) WatchFile(path string, fn fswatch.WatchCallback) (io.Closer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w := &MockWatch{Path: path, Callback: fn}
	m.Files[path] = w
	return w, nil
}

// SendEvents routes events through the registered watch callbacks
// that match each event's path. Directory watches match if the event
// path is a child (or recursive descendant) of the watched directory.
// File watches match on exact path. Events that match no watch are
// silently dropped — this is by design so that tests fail when the
// production code doesn't register the needed watches.
func (m *MockWatchBackend) SendEvents(events []fswatch.Event) {
	// Snapshot callbacks under the lock, then invoke outside the lock
	// to avoid deadlock if the callback re-enters the mock.
	m.mu.Lock()
	type target struct {
		cb     fswatch.WatchCallback
		events []fswatch.Event
	}
	targets := make(map[*MockWatch]*target)

	for _, e := range events {
		// Check file watches (exact match).
		if w, ok := m.Files[e.Path]; ok && !w.Closed {
			if t, ok := targets[w]; ok {
				t.events = append(t.events, e)
			} else {
				targets[w] = &target{cb: w.Callback, events: []fswatch.Event{e}}
			}
		}
		// Check directory watches.
		for _, w := range m.Dirs {
			if w.Closed {
				continue
			}
			if w.Ignore != nil && w.Ignore(e.Path) {
				continue
			}
			if !pathIsUnder(e.Path, w.Path, w.Recursive) {
				continue
			}
			if t, ok := targets[w]; ok {
				t.events = append(t.events, e)
			} else {
				targets[w] = &target{cb: w.Callback, events: []fswatch.Event{e}}
			}
		}
	}
	m.mu.Unlock()

	for _, t := range targets {
		t.cb(t.events, nil)
	}
}

// SendChangedPaths converts a list of file changes into fswatch
// events with appropriate event kinds and routes them through
// registered watches via SendEvents.
func (m *MockWatchBackend) SendChangedPaths(changes []fsbaselineutil.FileChange) {
	events := make([]fswatch.Event, len(changes))
	for i, c := range changes {
		kind := fswatch.EventUpdate
		if c.Deleted {
			kind = fswatch.EventDelete
		}
		events[i] = fswatch.Event{Kind: kind, Path: c.Path}
	}
	m.SendEvents(events)
}

// pathIsUnder reports whether eventPath is inside dir. If recursive is
// false, only direct children match.
func pathIsUnder(eventPath, dir string, recursive bool) bool {
	if !strings.HasPrefix(eventPath, dir) {
		return false
	}
	rest := eventPath[len(dir):]
	if len(rest) == 0 {
		return false // exact match = the dir itself, not a child
	}
	if rest[0] != '/' {
		return false // e.g. dir="/foo", path="/foobar"
	}
	if !recursive {
		// Direct child only: no further '/' after the separator.
		return !strings.Contains(rest[1:], "/")
	}
	return true
}

// WatchState returns a deterministic, human-readable summary of all
// active watches. This is intended to be included in test baselines
// so that watch registration correctness is verified via snapshot diffs.
func (m *MockWatchBackend) WatchState() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var b strings.Builder
	b.WriteString("Watch Registrations::\n")

	// Directory watches, sorted by path.
	var dirs []string
	for path, w := range m.Dirs {
		if !w.Closed {
			dirs = append(dirs, path)
		}
	}
	sort.Strings(dirs)

	b.WriteString("Directory watches::\n")
	if len(dirs) == 0 {
		b.WriteString("  (none)\n")
	}
	for _, d := range dirs {
		w := m.Dirs[d]
		if w.Recursive {
			fmt.Fprintf(&b, "  %s (recursive)\n", d)
		} else {
			fmt.Fprintf(&b, "  %s\n", d)
		}
	}

	// File watches, sorted by path.
	var files []string
	for path, w := range m.Files {
		if !w.Closed {
			files = append(files, path)
		}
	}
	sort.Strings(files)

	b.WriteString("File watches::\n")
	if len(files) == 0 {
		b.WriteString("  (none)\n")
	}
	for _, f := range files {
		fmt.Fprintf(&b, "  %s\n", f)
	}

	return b.String()
}

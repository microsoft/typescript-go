// This package implements a polling-based file watcher designed
// for use by both the CLI watcher and the language server.
package vfswatch

import (
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
)

const DebounceWait = 250 * time.Millisecond

type WatchEntry struct {
	ModTime    time.Time
	Exists     bool
	ChildCount int // -1 if not tracked
}

type FileWatcher struct {
	fs                  vfs.FS
	pollInterval        time.Duration
	testing             bool
	callback            func()
	watchState          map[string]WatchEntry
	wildcardDirectories map[string]bool
	mu                  sync.Mutex
}

func NewFileWatcher(fs vfs.FS, pollInterval time.Duration, testing bool, callback func()) *FileWatcher {
	return &FileWatcher{
		fs:           fs,
		pollInterval: pollInterval,
		testing:      testing,
		callback:     callback,
	}
}

func (fw *FileWatcher) SetPollInterval(d time.Duration) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.pollInterval = d
}

func (fw *FileWatcher) WatchStateEntry(path string) (WatchEntry, bool) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	e, ok := fw.watchState[path]
	return e, ok
}

func (fw *FileWatcher) WatchStateIsEmpty() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.watchState == nil
}

func (fw *FileWatcher) UpdateWatchState(paths []string, wildcardDirs map[string]bool) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.watchState = make(map[string]WatchEntry, len(paths))
	for _, fn := range paths {
		if s := fw.fs.Stat(fn); s != nil {
			fw.watchState[fn] = WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: -1}
		} else {
			fw.watchState[fn] = WatchEntry{Exists: false, ChildCount: -1}
		}
	}
	fw.wildcardDirectories = wildcardDirs
	for dir, recursive := range wildcardDirs {
		if !recursive {
			continue
		}
		_ = fw.fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			entries := fw.fs.GetAccessibleEntries(path)
			count := len(entries.Files) + len(entries.Directories)
			if existing, ok := fw.watchState[path]; ok {
				existing.ChildCount = count
				fw.watchState[path] = existing
			} else {
				if s := fw.fs.Stat(path); s != nil {
					fw.watchState[path] = WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: count}
				}
			}
			return nil
		})
	}
}

func (fw *FileWatcher) WaitForSettled(now func() time.Time) {
	if fw.testing {
		return
	}
	fw.mu.Lock()
	wildcardDirs := fw.wildcardDirectories
	pollInterval := fw.pollInterval
	fw.mu.Unlock()
	current := fw.currentState()
	settledAt := now()
	tick := min(pollInterval, DebounceWait)
	for now().Sub(settledAt) < DebounceWait {
		time.Sleep(tick)
		if fw.hasChanges(current, wildcardDirs) {
			current = fw.currentState()
			settledAt = now()
		}
	}
}

func (fw *FileWatcher) currentState() map[string]WatchEntry {
	fw.mu.Lock()
	watchState := fw.watchState
	wildcardDirs := fw.wildcardDirectories
	fw.mu.Unlock()
	state := make(map[string]WatchEntry, len(watchState))
	for path := range watchState {
		if s := fw.fs.Stat(path); s != nil {
			state[path] = WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: -1}
		} else {
			state[path] = WatchEntry{Exists: false, ChildCount: -1}
		}
	}
	for dir, recursive := range wildcardDirs {
		if !recursive {
			continue
		}
		_ = fw.fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			entries := fw.fs.GetAccessibleEntries(path)
			count := len(entries.Files) + len(entries.Directories)
			if existing, ok := state[path]; ok {
				existing.ChildCount = count
				state[path] = existing
			} else {
				if s := fw.fs.Stat(path); s != nil {
					state[path] = WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: count}
				}
			}
			return nil
		})
	}
	return state
}

func (fw *FileWatcher) hasChanges(baseline map[string]WatchEntry, wildcardDirs map[string]bool) bool {
	for path, old := range baseline {
		s := fw.fs.Stat(path)
		if !old.Exists {
			if s != nil {
				return true
			}
		} else {
			if s == nil || !s.ModTime().Equal(old.ModTime) {
				return true
			}
		}
	}
	for dir, recursive := range wildcardDirs {
		if !recursive {
			continue
		}
		found := false
		_ = fw.fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			entry, ok := baseline[path]
			if !ok {
				found = true
				return vfs.SkipAll
			}
			if entry.ChildCount >= 0 {
				entries := fw.fs.GetAccessibleEntries(path)
				if len(entries.Files)+len(entries.Directories) != entry.ChildCount {
					found = true
					return vfs.SkipAll
				}
			}
			return nil
		})
		if found {
			return true
		}
	}
	return false
}

func (fw *FileWatcher) HasChangesFromWatchState() bool {
	fw.mu.Lock()
	ws := fw.watchState
	wildcardDirs := fw.wildcardDirectories
	fw.mu.Unlock()
	return fw.hasChanges(ws, wildcardDirs)
}

func (fw *FileWatcher) Run(now func() time.Time) {
	for {
		fw.mu.Lock()
		interval := fw.pollInterval
		ws := fw.watchState
		wildcardDirs := fw.wildcardDirectories
		fw.mu.Unlock()
		time.Sleep(interval)
		if ws == nil || fw.hasChanges(ws, wildcardDirs) {
			fw.WaitForSettled(now)
			fw.callback()
		}
	}
}

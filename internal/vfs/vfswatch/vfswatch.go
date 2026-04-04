// This package implements a polling-based file watcher designed
// for use by both the CLI watcher and the language server.
package vfswatch

import (
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/trackingvfs"
)

const DebounceWait = 250 * time.Millisecond

type FileWatcher struct {
	fs                  vfs.FS
	PollInterval        time.Duration
	testing             bool
	callback            func()
	WatchState          map[string]trackingvfs.WatchEntry
	WildcardDirectories map[string]bool
}

func NewFileWatcher(fs vfs.FS, pollInterval time.Duration, testing bool, callback func()) *FileWatcher {
	return &FileWatcher{
		fs:           fs,
		PollInterval: pollInterval,
		testing:      testing,
		callback:     callback,
	}
}

func (fw *FileWatcher) UpdateWatchedFiles(tfs *trackingvfs.FS) {
	fw.WatchState = make(map[string]trackingvfs.WatchEntry)
	tfs.SeenFiles.Range(func(fn string) bool {
		if s := fw.fs.Stat(fn); s != nil {
			fw.WatchState[fn] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: -1}
		} else {
			fw.WatchState[fn] = trackingvfs.WatchEntry{Exists: false, ChildCount: -1}
		}
		return true
	})
	for dir, recursive := range fw.WildcardDirectories {
		if !recursive {
			continue
		}
		_ = fw.fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			entries := fw.fs.GetAccessibleEntries(path)
			count := len(entries.Files) + len(entries.Directories)
			if existing, ok := fw.WatchState[path]; ok {
				existing.ChildCount = count
				fw.WatchState[path] = existing
			} else {
				if s := fw.fs.Stat(path); s != nil {
					fw.WatchState[path] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: count}
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
	current := fw.currentState()
	settledAt := now()
	for now().Sub(settledAt) < DebounceWait {
		time.Sleep(fw.PollInterval)
		if fw.HasChanges(current) {
			current = fw.currentState()
			settledAt = now()
		}
	}
}

func (fw *FileWatcher) currentState() map[string]trackingvfs.WatchEntry {
	state := make(map[string]trackingvfs.WatchEntry, len(fw.WatchState))
	for path := range fw.WatchState {
		if s := fw.fs.Stat(path); s != nil {
			state[path] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: -1}
		} else {
			state[path] = trackingvfs.WatchEntry{Exists: false, ChildCount: -1}
		}
	}
	for dir, recursive := range fw.WildcardDirectories {
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
					state[path] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: count}
				}
			}
			return nil
		})
	}
	return state
}

func (fw *FileWatcher) HasChanges(baseline map[string]trackingvfs.WatchEntry) bool {
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
	for dir, recursive := range fw.WildcardDirectories {
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

func (fw *FileWatcher) Run(now func() time.Time) {
	for {
		time.Sleep(fw.PollInterval)
		if fw.WatchState == nil || fw.HasChanges(fw.WatchState) {
			fw.WaitForSettled(now)
			fw.callback()
		}
	}
}

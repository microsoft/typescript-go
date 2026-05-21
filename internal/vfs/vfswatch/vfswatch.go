// This package tracks filesystem state and detects changes
// by comparing current state against a stored baseline.
package vfswatch

import (
	"slices"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/zeebo/xxh3"
)

type WatchEntry struct {
	ModTime      time.Time
	Exists       bool
	ChildrenHash uint64 // 0 if not tracked
}

type FileWatcher struct {
	fs         vfs.FS
	watchState map[string]WatchEntry
	mu         sync.Mutex
}

func NewFileWatcher(fs vfs.FS) *FileWatcher {
	return &FileWatcher{
		fs: fs,
	}
}

func (fw *FileWatcher) WatchStateEntry(path string) (WatchEntry, bool) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	e, ok := fw.watchState[path]
	return e, ok
}

func (fw *FileWatcher) WatchStateUninitialized() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.watchState == nil
}

func (fw *FileWatcher) UpdateWatchState(paths []string, wildcardDirs map[string]bool) {
	state := snapshotPaths(fw.fs, paths, wildcardDirs)
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.watchState = state
}

func snapshotPaths(fs vfs.FS, paths []string, wildcardDirs map[string]bool) map[string]WatchEntry {
	state := make(map[string]WatchEntry, len(paths))
	for _, fn := range paths {
		if s := fs.Stat(fn); s != nil {
			state[fn] = WatchEntry{ModTime: s.ModTime(), Exists: true}
		} else {
			state[fn] = WatchEntry{Exists: false}
		}
	}
	for dir, recursive := range wildcardDirs {
		if !recursive {
			snapshotDirEntry(fs, state, dir)
			continue
		}
		_ = fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			snapshotDirEntry(fs, state, path)
			return nil
		})
	}
	return state
}

func snapshotDirEntry(fs vfs.FS, state map[string]WatchEntry, dir string) {
	entries := fs.GetAccessibleEntries(dir)
	h := hashEntries(entries)
	if existing, ok := state[dir]; ok {
		existing.ChildrenHash = h
		state[dir] = existing
	} else {
		if s := fs.Stat(dir); s != nil {
			state[dir] = WatchEntry{ModTime: s.ModTime(), Exists: true, ChildrenHash: h}
		}
	}
}

func hashEntries(entries vfs.Entries) uint64 {
	dirs := slices.Clone(entries.Directories)
	files := slices.Clone(entries.Files)
	slices.Sort(dirs)
	slices.Sort(files)
	var h xxh3.Hasher
	for _, name := range dirs {
		_, _ = h.WriteString("d:")
		_, _ = h.WriteString(name)
		_, _ = h.Write([]byte{0})
	}
	for _, name := range files {
		_, _ = h.WriteString("f:")
		_, _ = h.WriteString(name)
		_, _ = h.Write([]byte{0})
	}
	return h.Sum64()
}

// hasChanges compares the current filesystem state against baseline.
//
// Tracked entries fall into two categories:
//
//   - Explicit paths (files the compiler depends on, plus directory paths
//     accessed via DirectoryExists/Stat/etc. during compilation). For these
//     we only need to know whether the path exists and, if it does, whether
//     its mtime has changed. We never depend on *what's inside* a directory
//     in this category — any specific file we care about is tracked
//     independently in this same map.
//
//   - Wildcard tree directories. snapshotPaths walks every directory under
//     each recursive wildcard root and stores it with a ChildrenHash that
//     covers the directory's listing. Re-hashing here detects any new,
//     deleted, or renamed file or subdirectory in those trees.
//
// Iterating baseline once therefore covers both: a single fs.Stat per entry,
// plus a fs.GetAccessibleEntries only for entries with ChildrenHash != 0
// (i.e. wildcard tree members).
func (fw *FileWatcher) hasChanges(baseline map[string]WatchEntry) bool {
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
			if old.ChildrenHash != 0 {
				entries := fw.fs.GetAccessibleEntries(path)
				if hashEntries(entries) != old.ChildrenHash {
					return true
				}
			}
		}
	}
	return false
}

// HasChangesFromWatchState compares the current filesystem against the
// stored watch state. Safe for concurrent use.
func (fw *FileWatcher) HasChangesFromWatchState() bool {
	fw.mu.Lock()
	ws := fw.watchState
	fw.mu.Unlock()
	return fw.hasChanges(ws)
}

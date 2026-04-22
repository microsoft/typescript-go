package vfswatch

import (
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/zeebo/xxh3"
)

const debounceWait = 250 * time.Millisecond

type WatchEvent struct {
	Created []string
	Deleted []string
	Changed []string
}

func (e WatchEvent) HasChanges() bool {
	return len(e.Created) > 0 || len(e.Deleted) > 0 || len(e.Changed) > 0
}

type watchEntry struct {
	modTime      time.Time
	childrenHash uint64 // 0 if not tracked
	isDir        bool
}

type FileWatcher struct {
	fs            vfs.FS
	pollInterval  time.Duration
	testing       bool
	callback      func(WatchEvent)
	directories   map[string]bool
	watchState    map[string]watchEntry
	watchStateGen uint64
	mu            sync.Mutex
	done          chan struct{}
	stopOnce      sync.Once
}

func NewFileWatcher(fs vfs.FS, pollInterval time.Duration, testing bool, callback func(WatchEvent)) *FileWatcher {
	return &FileWatcher{
		fs:           fs,
		pollInterval: pollInterval,
		testing:      testing,
		callback:     callback,
		done:         make(chan struct{}),
	}
}

func (fw *FileWatcher) Stop() {
	fw.stopOnce.Do(func() {
		close(fw.done)
	})
}

func (fw *FileWatcher) SetPollInterval(d time.Duration) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.pollInterval = d
}

func (fw *FileWatcher) UpdateWatchedDirectories(dirs map[string]bool) {
	state := snapshotDirectories(fw.fs, dirs)
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.directories = dirs
	fw.watchState = state
	fw.watchStateGen++
}

func (fw *FileWatcher) sleepOrDone(d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-fw.done:
		return true
	case <-t.C:
		return false
	}
}

func (fw *FileWatcher) WaitForSettled(now func() time.Time) {
	if fw.testing {
		return
	}
	fw.mu.Lock()
	dirs := fw.directories
	pollInterval := fw.pollInterval
	fw.mu.Unlock()
	current := snapshotDirectories(fw.fs, dirs)
	settledAt := now()
	tick := min(pollInterval, debounceWait)
	for now().Sub(settledAt) < debounceWait {
		if fw.sleepOrDone(tick) {
			return
		}
		next := snapshotDirectories(fw.fs, dirs)
		if diffSnapshots(current, next).HasChanges() {
			current = next
			settledAt = now()
		}
	}
}

func snapshotDirectories(fs vfs.FS, dirs map[string]bool) map[string]watchEntry {
	state := make(map[string]watchEntry)
	for dir, recursive := range dirs {
		if dir == "" {
			continue
		}
		snapshotDir(fs, state, dir, recursive)
	}
	return state
}

func joinWatchPath(dir, name string) string {
	if strings.HasSuffix(dir, "/") {
		return dir + name
	}
	return dir + "/" + name
}

func snapshotDir(fs vfs.FS, state map[string]watchEntry, dir string, recursive bool) {
	if !recursive {
		entries := fs.GetAccessibleEntries(dir)
		h := hashEntries(entries)
		if s := fs.Stat(dir); s != nil {
			state[dir] = watchEntry{modTime: s.ModTime(), childrenHash: h}
		}
		for _, file := range entries.Files {
			path := joinWatchPath(dir, file)
			if s := fs.Stat(path); s != nil {
				state[path] = watchEntry{modTime: s.ModTime()}
			}
		}
		for _, subdir := range entries.Directories {
			path := joinWatchPath(dir, subdir)
			if s := fs.Stat(path); s != nil {
				state[path] = watchEntry{modTime: s.ModTime(), isDir: true}
			}
		}
		return
	}
	_ = fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			snapshotDirEntry(fs, state, path)
		} else {
			if s := fs.Stat(path); s != nil {
				state[path] = watchEntry{modTime: s.ModTime()}
			}
		}
		return nil
	})
}

func snapshotDirEntry(fs vfs.FS, state map[string]watchEntry, dir string) {
	entries := fs.GetAccessibleEntries(dir)
	h := hashEntries(entries)
	if existing, ok := state[dir]; ok {
		existing.childrenHash = h
		state[dir] = existing
	} else {
		if s := fs.Stat(dir); s != nil {
			state[dir] = watchEntry{modTime: s.ModTime(), childrenHash: h}
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

func diffSnapshots(baseline, current map[string]watchEntry) WatchEvent {
	var event WatchEvent

	for path, old := range baseline {
		cur, exists := current[path]
		if !exists {
			event.Deleted = append(event.Deleted, path)
			continue
		}
		if !cur.modTime.Equal(old.modTime) {
			if old.childrenHash == 0 && !old.isDir {
				event.Changed = append(event.Changed, path)
			}
			continue
		}
		if old.childrenHash != 0 && cur.childrenHash != old.childrenHash {
			continue
		}
	}

	for path := range current {
		if _, inBaseline := baseline[path]; !inBaseline {
			event.Created = append(event.Created, path)
		}
	}

	return event
}

func (fw *FileWatcher) ScanForChanges() WatchEvent {
	fw.mu.Lock()
	ws := fw.watchState
	dirs := fw.directories
	fw.mu.Unlock()
	current := snapshotDirectories(fw.fs, dirs)
	return diffSnapshots(ws, current)
}

func (fw *FileWatcher) Run(now func() time.Time) {
	for {
		fw.mu.Lock()
		interval := fw.pollInterval
		ws := fw.watchState
		dirs := fw.directories
		fw.mu.Unlock()

		if fw.sleepOrDone(interval) {
			return
		}

		if ws == nil {
			continue
		}
		current := snapshotDirectories(fw.fs, dirs)
		event := diffSnapshots(ws, current)
		if event.HasChanges() {
			fw.WaitForSettled(now)
			fw.mu.Lock()
			gen := fw.watchStateGen
			ws = fw.watchState
			dirs = fw.directories
			fw.mu.Unlock()
			current = snapshotDirectories(fw.fs, dirs)
			event = diffSnapshots(ws, current)
			if event.HasChanges() {
				fw.callback(event)
			}
			fw.mu.Lock()
			if fw.watchStateGen == gen {
				fw.watchState = current
			}
			fw.mu.Unlock()
		}
	}
}

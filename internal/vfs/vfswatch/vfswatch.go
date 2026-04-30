package vfswatch

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/zeebo/xxh3"
)

// DebounceWait is the duration the polling watcher waits to coalesce rapid
// filesystem changes before reporting them as a single watch event.
const DebounceWait = 250 * time.Millisecond

type WatchLogger interface {
	Log(msg ...any)
	Logf(format string, args ...any)
}

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
	logger        WatchLogger
}

func NewFileWatcher(fs vfs.FS, pollInterval time.Duration, testing bool, callback func(WatchEvent)) *FileWatcher {
	return &FileWatcher{
		fs:           fs,
		pollInterval: pollInterval,
		testing:      testing,
		callback:     callback,
	}
}

func (fw *FileWatcher) SetLogger(logger WatchLogger) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.logger = logger
}

func (fw *FileWatcher) SetPollInterval(d time.Duration) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.pollInterval = d
}

func (fw *FileWatcher) UpdateWatchedDirectories(dirs map[string]bool) {
	start := time.Now()
	state := snapshotDirectories(fw.fs, dirs)
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.directories = dirs
	fw.watchState = state
	fw.watchStateGen++
	if fw.logger != nil {
		fw.logger.Logf("Polling watcher: watching %d directories (%d entries) in %v", len(dirs), len(state), time.Since(start))
	}
}

// sleepOrDone sleeps for the given duration, returning true if the context
// was cancelled during the sleep.
func sleepOrDone(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return true
	case <-t.C:
		return false
	}
}

func (fw *FileWatcher) WaitForSettled(ctx context.Context) {
	if fw.testing {
		return
	}
	fw.mu.Lock()
	dirs := fw.directories
	pollInterval := fw.pollInterval
	fw.mu.Unlock()
	current := snapshotDirectories(fw.fs, dirs)
	settledAt := time.Now()
	tick := min(pollInterval, DebounceWait)
	for time.Since(settledAt) < DebounceWait {
		if sleepOrDone(ctx, tick) {
			return
		}
		next := snapshotDirectories(fw.fs, dirs)
		if diffSnapshots(current, next).HasChanges() {
			current = next
			settledAt = time.Now()
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
	logger := fw.logger
	fw.mu.Unlock()
	scanStart := time.Now()
	current := snapshotDirectories(fw.fs, dirs)
	event := diffSnapshots(ws, current)
	if logger != nil {
		scanDuration := time.Since(scanStart)
		if event.HasChanges() {
			logger.Logf("Polling watcher scan: %d entries in %v, %s", len(current), scanDuration, formatEvent(event))
		} else if scanDuration > 100*time.Millisecond {
			logger.Logf("Polling watcher scan: %d entries in %v, no changes", len(current), scanDuration)
		}
	}
	return event
}

func (fw *FileWatcher) Run(ctx context.Context) {
	for {
		fw.mu.Lock()
		interval := fw.pollInterval
		ws := fw.watchState
		dirs := fw.directories
		logger := fw.logger
		fw.mu.Unlock()

		if sleepOrDone(ctx, interval) {
			return
		}

		if ws == nil {
			continue
		}
		scanStart := time.Now()
		current := snapshotDirectories(fw.fs, dirs)
		scanDuration := time.Since(scanStart)
		event := diffSnapshots(ws, current)
		if event.HasChanges() {
			if logger != nil {
				logger.Logf("Polling watcher: changes detected (scanned %d entries in %v)", len(current), scanDuration)
				for _, p := range event.Created {
					logger.Logf("  created: %s", p)
				}
				for _, p := range event.Changed {
					logger.Logf("  changed: %s", p)
				}
				for _, p := range event.Deleted {
					logger.Logf("  deleted: %s", p)
				}
			}
			fw.WaitForSettled(ctx)
			fw.mu.Lock()
			gen := fw.watchStateGen
			ws = fw.watchState
			dirs = fw.directories
			fw.mu.Unlock()
			current = snapshotDirectories(fw.fs, dirs)
			event = diffSnapshots(ws, current)
			if event.HasChanges() {
				if logger != nil {
					logger.Logf("Polling watcher: settled with %s", formatEvent(event))
				}
				fw.callback(event)
			}
			fw.mu.Lock()
			if fw.watchStateGen == gen {
				fw.watchState = current
			}
			fw.mu.Unlock()
		} else if logger != nil && scanDuration > 100*time.Millisecond {
			logger.Logf("Polling watcher: scan took %v (%d entries), no changes", scanDuration, len(current))
		}
	}
}

func formatEvent(e WatchEvent) string {
	return fmt.Sprintf("created=%d deleted=%d changed=%d", len(e.Created), len(e.Deleted), len(e.Changed))
}

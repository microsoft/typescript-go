package watchmanager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fswatch"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type watchedDir struct {
	closer    io.Closer
	recursive bool
}

// WatchManager manages fswatch directory watches, event accumulation,
// and DoCycle signaling. It is shared by the CLI watcher and the build
// mode orchestrator.
//
// Locking contract:
//   - Call Lock/Unlock around the entire DoCycle body.
//   - ReconcileWatches must be called under Lock.
//   - CloseAllWatches and handleWatchTerminated manage their own locking.
type WatchManager struct {
	mu          sync.Mutex
	backend     WatchBackend
	watchedDirs map[string]*watchedDir
	doCycleCh   chan struct{}

	// DebugLog receives verbose watch diagnostics when non-nil
	DebugLog io.Writer

	warnWriter io.Writer
	dirExists  func(string) bool

	changedMu       sync.Mutex
	changedPaths    map[string]fswatch.EventKind
	changedOverflow bool
}

const recursiveCoalesceThreshold = 10

func NewWatchManager(warnWriter io.Writer, dirExists func(string) bool) *WatchManager {
	return &WatchManager{
		watchedDirs: make(map[string]*watchedDir),
		doCycleCh:   make(chan struct{}, 1),
		warnWriter:  warnWriter,
		dirExists:   dirExists,
	}
}

func (wm *WatchManager) SetBackend(b WatchBackend) { wm.backend = b }

func (wm *WatchManager) Backend() WatchBackend { return wm.backend }

func (wm *WatchManager) EnsureDefaultBackend() {
	if wm.backend == nil {
		fsw := fswatch.Default()
		wm.backend = &FSWatchBackend{Inner: fsw}
		if wm.DebugLog != nil {
			fmt.Fprintf(wm.DebugLog, "[watch] using %s backend\n", fsw.Name())
		}
	}
}

func (wm *WatchManager) Lock() { wm.mu.Lock() }

func (wm *WatchManager) Unlock() { wm.mu.Unlock() }

func (wm *WatchManager) DoCycleCh() <-chan struct{} { return wm.doCycleCh }

func (wm *WatchManager) DrainEvents() (changed map[string]fswatch.EventKind, overflow bool) {
	wm.changedMu.Lock()
	changed = wm.changedPaths
	overflow = wm.changedOverflow
	wm.changedPaths = nil
	wm.changedOverflow = false
	wm.changedMu.Unlock()
	return
}

func (wm *WatchManager) ForceOverflow() {
	wm.changedMu.Lock()
	wm.changedOverflow = true
	wm.changedMu.Unlock()
}

func (wm *WatchManager) signalDoCycle() {
	select {
	case wm.doCycleCh <- struct{}{}:
		// Signal sent; the DoCycle loop will pick it up.
	default:
		// A signal is already pending; coalesced.
	}
}

func (wm *WatchManager) onWatchEvents(events []fswatch.Event, err error) {
	if err != nil {
		if errors.Is(err, fswatch.ErrOverflow) {
			if wm.DebugLog != nil {
				fmt.Fprintf(wm.DebugLog, "[watch] event overflow, triggering rebuild\n")
			}
			wm.changedMu.Lock()
			wm.changedOverflow = true
			wm.changedMu.Unlock()
			wm.signalDoCycle()
			return
		}
		fmt.Fprintf(wm.warnWriter, "Warning: File watch error: %v\n", err)
		return
	}

	if len(events) > 0 {
		if wm.DebugLog != nil {
			fmt.Fprintf(wm.DebugLog, "[watch] %d event(s): ", len(events))
			for i, e := range events {
				if i > 0 {
					fmt.Fprint(wm.DebugLog, ", ")
				}
				if i >= 5 {
					fmt.Fprintf(wm.DebugLog, "... and %d more", len(events)-i)
					break
				}
				fmt.Fprintf(wm.DebugLog, "%s %s", e.Kind, e.Path)
			}
			fmt.Fprintln(wm.DebugLog)
		}
		wm.changedMu.Lock()
		if wm.changedPaths == nil {
			wm.changedPaths = make(map[string]fswatch.EventKind, len(events))
		}
		for _, e := range events {
			wm.changedPaths[e.Path] = e.Kind
		}
		wm.changedMu.Unlock()
		wm.signalDoCycle()
	}
}

func (wm *WatchManager) handleWatchTerminated(dir string, identity *watchedDir) {
	if wm.DebugLog != nil {
		fmt.Fprintf(wm.DebugLog, "[watch] watch terminated: %s\n", dir)
	}
	var staleCloser io.Closer
	wm.mu.Lock()
	if wd, ok := wm.watchedDirs[dir]; ok && wd == identity {
		staleCloser = wd.closer
		delete(wm.watchedDirs, dir)
	}
	wm.mu.Unlock()
	if staleCloser != nil {
		staleCloser.Close()
	}
	wm.changedMu.Lock()
	wm.changedOverflow = true
	wm.changedMu.Unlock()
	wm.signalDoCycle()
}

func (wm *WatchManager) CloseAllWatches() {
	wm.mu.Lock()
	closers := make([]io.Closer, 0, len(wm.watchedDirs))
	for dir, wd := range wm.watchedDirs {
		closers = append(closers, wd.closer)
		delete(wm.watchedDirs, dir)
	}
	wm.mu.Unlock()
	for _, c := range closers {
		c.Close()
	}
}

func (wm *WatchManager) createDirWatch(dir string, recursive bool) error {
	entry := &watchedDir{recursive: recursive}
	cb := func(events []fswatch.Event, err error) {
		if err != nil && errors.Is(err, fswatch.ErrWatchTerminated) {
			wm.handleWatchTerminated(dir, entry)
			return
		}
		wm.onWatchEvents(events, err)
	}
	watch, err := wm.backend.WatchDirectory(dir, cb, recursive, ShouldIgnoreWatchPath)
	if err != nil {
		if wm.DebugLog != nil {
			fmt.Fprintf(wm.DebugLog, "[watch] failed to watch directory %s: %v\n", dir, err)
		}
		return fmt.Errorf("failed to watch directory %s: %w", dir, err)
	}
	entry.closer = watch
	wm.watchedDirs[dir] = entry
	return nil
}

func (wm *WatchManager) ResolveDesiredDirs(desiredDirs map[string]bool) map[string]bool {
	resolved := make(map[string]bool, len(desiredDirs))
	for dir, recursive := range desiredDirs {
		watchDir := dir
		watchRecursive := recursive
		for !wm.dirExists(watchDir) {
			parent := tspath.GetDirectoryPath(watchDir)
			if parent == watchDir {
				break
			}
			watchDir = parent
			watchRecursive = false // ancestor fallbacks are always non-recursive
		}
		if !wm.dirExists(watchDir) || !CanWatchDirectory(watchDir) {
			if wm.DebugLog != nil {
				fmt.Fprintf(wm.DebugLog, "[watch] no watchable ancestor for %s\n", dir)
			}
			continue
		}
		if watchDir != dir && wm.DebugLog != nil {
			fmt.Fprintf(wm.DebugLog, "[watch] resolved %s to ancestor %s\n", dir, watchDir)
		}
		if existing, has := resolved[watchDir]; has {
			resolved[watchDir] = existing || watchRecursive
		} else {
			resolved[watchDir] = watchRecursive
		}
	}
	return resolved
}

func (wm *WatchManager) CoalesceDesiredDirs(desiredDirs map[string]bool, opts tspath.ComparePathsOptions) map[string]bool {
	if wm.backend == nil || !wm.backend.HasFastRecursiveBackend() || len(desiredDirs) < recursiveCoalesceThreshold {
		return desiredDirs
	}

	pruned := removeDirsCoveredByRecursiveWatches(desiredDirs, opts)
	coalesced := coalesceAncestorDirs(pruned, opts)
	if len(coalesced) == len(desiredDirs) {
		return desiredDirs
	}
	return coalesced
}

func removeDirsCoveredByRecursiveWatches(desiredDirs map[string]bool, opts tspath.ComparePathsOptions) map[string]bool {
	result := make(map[string]bool, len(desiredDirs))
	for dir, recursive := range desiredDirs {
		covered := false
		for parent, parentRecursive := range desiredDirs {
			if parent == dir || !parentRecursive {
				continue
			}
			if tspath.ContainsPath(parent, dir, opts) {
				covered = true
				break
			}
		}
		if !covered {
			result[dir] = recursive
		}
	}
	return result
}

func coalesceAncestorDirs(desiredDirs map[string]bool, opts tspath.ComparePathsOptions) map[string]bool {
	type candidate struct {
		dir   string
		depth int
	}
	counts := make(map[string]int)
	for dir := range desiredDirs {
		parent := tspath.GetDirectoryPath(dir)
		for parent != "" && parent != dir {
			if CanWatchDirectory(parent) {
				counts[parent]++
			}
			next := tspath.GetDirectoryPath(parent)
			if next == parent {
				break
			}
			parent = next
		}
	}

	candidates := make([]candidate, 0, len(counts))
	for dir, count := range counts {
		if count >= recursiveCoalesceThreshold {
			candidates = append(candidates, candidate{dir: dir, depth: len(tspath.GetPathComponents(dir, opts.CurrentDirectory))})
		}
	}
	slices.SortFunc(candidates, func(a, b candidate) int {
		if a.depth != b.depth {
			return b.depth - a.depth
		}
		return opts.GetComparer()(a.dir, b.dir)
	})

	remaining := make(map[string]bool, len(desiredDirs))
	maps.Copy(remaining, desiredDirs)
	selected := make(map[string]struct{})
	for _, candidate := range candidates {
		var covered []string
		for dir := range remaining {
			if tspath.ContainsPath(candidate.dir, dir, opts) {
				covered = append(covered, dir)
			}
		}
		if len(covered) < recursiveCoalesceThreshold {
			continue
		}
		selected[candidate.dir] = struct{}{}
		for _, dir := range covered {
			delete(remaining, dir)
		}
	}

	if len(selected) == 0 {
		return desiredDirs
	}

	result := make(map[string]bool, len(remaining)+len(selected))
	maps.Copy(result, remaining)
	for dir := range selected {
		result[dir] = true
	}
	return result
}

func (wm *WatchManager) ReconcileWatches(desiredDirs map[string]bool) error {
	if wm.backend == nil {
		return nil
	}

	var watchErr error
	core.DiffMapsFunc(
		wm.watchedDirs,
		desiredDirs,
		func(wd *watchedDir, recursive bool) bool { return wd.recursive == recursive },
		func(dir string, recursive bool) {
			if wm.DebugLog != nil {
				fmt.Fprintf(wm.DebugLog, "[watch] watching directory %s (recursive=%v)\n", dir, recursive)
			}
			if err := wm.createDirWatch(dir, recursive); err != nil && watchErr == nil {
				watchErr = err
			}
		},
		func(dir string, wd *watchedDir) {
			if wm.DebugLog != nil {
				fmt.Fprintf(wm.DebugLog, "[watch] closing stale dir watch: %s\n", dir)
			}
			wd.closer.Close()
			delete(wm.watchedDirs, dir)
		},
		func(dir string, wd *watchedDir, recursive bool) {
			if wm.DebugLog != nil {
				fmt.Fprintf(wm.DebugLog, "[watch] recreating dir watch %s (recursive %v→%v)\n", dir, wd.recursive, recursive)
			}
			wd.closer.Close()
			delete(wm.watchedDirs, dir)
			if err := wm.createDirWatch(dir, recursive); err != nil && watchErr == nil {
				watchErr = err
			}
		},
	)
	return watchErr
}

func IsDirCoveredByWatch(dirs map[string]bool, dir string, opts tspath.ComparePathsOptions) bool {
	for wdir, recursive := range dirs {
		if recursive {
			if tspath.ContainsPath(wdir, dir, opts) {
				return true
			}
		} else if dir == wdir {
			return true
		}
	}
	return false
}

func (wm *WatchManager) IsPathUnderWatch(path string, opts tspath.ComparePathsOptions) bool {
	for dir := range wm.watchedDirs {
		if tspath.ContainsPath(dir, path, opts) {
			return true
		}
	}
	return false
}

func (wm *WatchManager) RunLoop(ctx context.Context, doCycle func()) {
	for {
		select {
		case <-ctx.Done():
			wm.CloseAllWatches()
			return
		case <-wm.doCycleCh:
			doCycle()
		}
	}
}

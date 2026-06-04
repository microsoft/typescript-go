// Package lspwatcher implements an in-process file watcher used as a
// drop-in replacement for LSP-based file watching when the client either
// does not support dynamic registration of file watchers or has
// requested the server-side watcher explicitly via the
// `useBuiltinWatcher` initialization option.
package lspwatcher

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/fswatch"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// throttleWindow mirrors VS Code's parcel watcher integration: give the
// first batch a short grace window so adjacent filesystem bursts coalesce.
const throttleWindow = 75 * time.Millisecond

type watch interface {
	Close() error
}

type watcherBackend interface {
	WatchDirectory(dir string, fn fswatch.WatchCallback, opts ...fswatch.WatchOption) (watch, error)
}

type defaultWatcherBackend struct {
	w fswatch.Watcher
}

func (d defaultWatcherBackend) WatchDirectory(dir string, fn fswatch.WatchCallback, opts ...fswatch.WatchOption) (watch, error) {
	return d.w.WatchDirectory(dir, fn, opts...)
}

// Watcher manages a set of file system subscriptions identified by
// WatcherID strings (matching the LSP server's project.WatcherID type).
// Events are delivered to onChanges in batches as `*lsproto.FileEvent`,
// shaped exactly like a `workspace/didChangeWatchedFiles` notification.
type Watcher struct {
	fs        vfs.FS
	backend   watcherBackend
	onChanges func(changes []*lsproto.FileEvent)
	logger    logging.Logger

	mu     sync.Mutex
	subs   map[string]*idSubscription
	closed bool

	// Pending batch state, protected by mu.
	pending    map[string]*lsproto.FileEvent
	flushTimer *time.Timer
}

// idSubscription holds the fswatch subscriptions associated with a single
// LSP WatcherID — there may be more than one because each
// FileSystemWatcher in the registration becomes its own subscription
// (different roots and kinds).
type idSubscription struct {
	subs []*singleWatch
}

// All path fields below are tspath-style (forward-slash) absolute paths.
type singleWatch struct {
	watch    watch
	rootReal string // canonicalized (symlink-resolved) root dir
	rootOrig string // root as requested by the LSP layer
	kind     lsproto.WatchKind
}

// New constructs a Watcher backed by internal/fswatch's platform-default
// watcher implementation.
func New(fs vfs.FS, onChanges func(changes []*lsproto.FileEvent), logger logging.Logger) *Watcher {
	return newWithBackend(fs, defaultWatcherBackend{w: fswatch.Default()}, onChanges, logger)
}

func newWithBackend(fs vfs.FS, backend watcherBackend, onChanges func(changes []*lsproto.FileEvent), logger logging.Logger) *Watcher {
	return &Watcher{
		fs:        fs,
		backend:   backend,
		onChanges: onChanges,
		logger:    logger,
		subs:      make(map[string]*idSubscription),
	}
}

// WatchFiles subscribes to each FileSystemWatcher under the given id.
// Errors registering individual watchers are logged but not returned —
// the rest of the batch should still be wired up so the session has at
// least partial coverage.
func (w *Watcher) WatchFiles(id string, watchers []*lsproto.FileSystemWatcher) error {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return errors.New("lspwatcher: closed")
	}
	if _, exists := w.subs[id]; exists {
		w.mu.Unlock()
		return fmt.Errorf("lspwatcher: watcher %q already exists", id)
	}
	entry := &idSubscription{}
	w.subs[id] = entry
	w.mu.Unlock()

	for _, fsw := range watchers {
		root, ok := watchRoot(fsw)
		if !ok || root == "" {
			w.logger.Logf("lspwatcher: skipping watcher %q: unrecognized pattern %q", id, watchPatternString(fsw))
			continue
		}
		kind := effectiveKind(fsw)
		realRoot := w.fs.Realpath(root)
		if !w.fs.DirectoryExists(realRoot) {
			w.logger.Logf("lspwatcher: cannot watch %q: directory does not exist", realRoot)
			continue
		}
		sw := &singleWatch{rootReal: realRoot, rootOrig: root, kind: kind}
		sub, err := w.backend.WatchDirectory(realRoot, w.makeCallback(sw), fswatch.WithRecursive())
		if err != nil {
			w.logger.Logf("lspwatcher: failed to subscribe to %q: %v", realRoot, err)
			continue
		}
		sw.watch = sub
		w.mu.Lock()
		entry.subs = append(entry.subs, sw)
		w.mu.Unlock()
	}
	return nil
}

// UnwatchFiles tears down all subscriptions associated with id.
func (w *Watcher) UnwatchFiles(id string) error {
	w.mu.Lock()
	entry, ok := w.subs[id]
	if !ok {
		w.mu.Unlock()
		return fmt.Errorf("lspwatcher: no watcher with id %q", id)
	}
	delete(w.subs, id)
	w.mu.Unlock()
	for _, sw := range entry.subs {
		if err := sw.watch.Close(); err != nil {
			w.logger.Logf("lspwatcher: unsubscribe %q: %v", sw.rootReal, err)
		}
	}
	return nil
}

// Close removes every subscription. Safe to call multiple times.
func (w *Watcher) Close() {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	w.closed = true
	all := w.subs
	w.subs = nil
	if w.flushTimer != nil {
		w.flushTimer.Stop()
		w.flushTimer = nil
	}
	w.pending = nil
	w.mu.Unlock()
	for _, entry := range all {
		for _, sw := range entry.subs {
			if err := sw.watch.Close(); err != nil {
				w.logger.Logf("lspwatcher: unsubscribe %q: %v", sw.rootReal, err)
			}
		}
	}
}

func (w *Watcher) makeCallback(sw *singleWatch) fswatch.WatchCallback {
	return func(events []fswatch.Event, err error) {
		if err != nil {
			if errors.Is(err, fswatch.ErrOverflow) {
				w.logger.Logf("lspwatcher: watch overflow in %q (some events may have been dropped): %v", sw.rootReal, err)
			} else {
				w.logger.Logf("lspwatcher: watch error in %q: %v", sw.rootReal, err)
			}
		}
		if len(events) == 0 {
			return
		}
		w.handleEvents(sw, events)
	}
}

func (w *Watcher) handleEvents(sw *singleWatch, events []fswatch.Event) {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	if w.pending == nil {
		w.pending = make(map[string]*lsproto.FileEvent, len(events))
	}
	for _, e := range events {
		var changeType lsproto.FileChangeType
		switch e.Kind {
		case fswatch.EventUpdate:
			// fswatch intentionally doesn't distinguish create vs update.
			// For LSP consumers this is fine: callers infer create/update
			// from their own cache and both should invalidate stale state.
			if sw.kind&(lsproto.WatchKindCreate|lsproto.WatchKindChange) == 0 {
				continue
			}
			changeType = lsproto.FileChangeTypeChanged
		case fswatch.EventDelete:
			if sw.kind&lsproto.WatchKindDelete == 0 {
				continue
			}
			changeType = lsproto.FileChangeTypeDeleted
		default:
			continue
		}

		path := tspath.NormalizeSlashes(e.Path)
		if sw.rootReal != sw.rootOrig {
			opts := tspath.ComparePathsOptions{UseCaseSensitiveFileNames: w.fs.UseCaseSensitiveFileNames()}
			if tspath.ContainsPath(sw.rootReal, path, opts) {
				rel := tspath.GetRelativePathFromDirectory(sw.rootReal, path, opts)
				if rel == "" {
					path = sw.rootOrig
				} else {
					path = tspath.CombinePaths(sw.rootOrig, rel)
				}
			}
		}

		uri := lsconv.FileNameToDocumentURI(path)
		w.pending[string(uri)] = &lsproto.FileEvent{
			Uri:  uri,
			Type: changeType,
		}
	}

	if w.flushTimer == nil {
		w.flushTimer = time.AfterFunc(throttleWindow, w.flush)
	}
	w.mu.Unlock()
}

func (w *Watcher) flush() {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	pending := w.pending
	w.pending = nil
	w.flushTimer = nil
	w.mu.Unlock()

	if len(pending) == 0 {
		return
	}
	changes := make([]*lsproto.FileEvent, 0, len(pending))
	for _, ev := range pending {
		changes = append(changes, ev)
	}
	w.onChanges(changes)
}

// watchRoot extracts the directory the fswatch subscription should be
// rooted at from a FileSystemWatcher. The patterns the project layer
// produces are always of the form `<dir>/**/*` (either as a Pattern
// with a fully-qualified directory or as a RelativePattern with a
// file:// BaseUri and a `**/*` pattern), so the heuristic of
// "everything before the first glob meta character" is reliable.
//
// Returned roots are tspath-normalized (forward-slash) absolute paths.
func watchRoot(fsw *lsproto.FileSystemWatcher) (string, bool) {
	if fsw.GlobPattern.Pattern != nil {
		return rootFromGlob(*fsw.GlobPattern.Pattern), true
	}
	if rp := fsw.GlobPattern.RelativePattern; rp != nil {
		var base string
		if rp.BaseUri.URI != nil {
			base = lsproto.DocumentUri(*rp.BaseUri.URI).FileName()
		} else {
			return "", false
		}
		joined := tspath.CombinePaths(base, rp.Pattern)
		return rootFromGlob(joined), true
	}
	return "", false
}

func rootFromGlob(pattern string) string {
	pattern = tspath.NormalizeSlashes(pattern)
	idx := -1
	for i := range len(pattern) {
		switch pattern[i] {
		case '*', '?', '[', '{':
			idx = i
		}
		if idx != -1 {
			break
		}
	}
	if idx == -1 {
		return tspath.NormalizePath(strings.TrimRight(pattern, "/"))
	}
	dir := pattern[:idx]
	dir = strings.TrimRight(dir, "/")
	if dir == "" {
		return ""
	}
	return tspath.NormalizePath(dir)
}

func watchPatternString(fsw *lsproto.FileSystemWatcher) string {
	if fsw.GlobPattern.Pattern != nil {
		return *fsw.GlobPattern.Pattern
	}
	if rp := fsw.GlobPattern.RelativePattern; rp != nil {
		var base string
		if rp.BaseUri.URI != nil {
			base = string(*rp.BaseUri.URI)
		}
		return base + "/" + rp.Pattern
	}
	return ""
}

func effectiveKind(fsw *lsproto.FileSystemWatcher) lsproto.WatchKind {
	if fsw.Kind != nil {
		return *fsw.Kind
	}
	return lsproto.WatchKindCreate | lsproto.WatchKindChange | lsproto.WatchKindDelete
}

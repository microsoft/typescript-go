package project

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

const (
	minWatchLocationDepth = 2
)

type fileSystemWatcherKey struct {
	pattern string
	kind    lsproto.WatchKind
}

type fileSystemWatcherValue struct {
	count int
	id    WatcherID
}

// watchRegistry tracks the current watch globs and how many individual
// WatchedFiles reference each glob. It provides ref-count helpers so callers
// don't manipulate the map directly.
//
// All methods are safe for concurrent use; locking is handled internally.
type watchRegistry struct {
	mu      sync.Mutex
	entries map[fileSystemWatcherKey]*fileSystemWatcherValue
	pending map[WatcherID]struct{}
}

func newWatchRegistry() *watchRegistry {
	return &watchRegistry{
		entries: make(map[fileSystemWatcherKey]*fileSystemWatcherValue),
		pending: make(map[WatcherID]struct{}),
	}
}

// Acquire increments the ref count for a watcher. If this is the first
// reference (count goes from 0 to 1), it returns true so the caller knows
// to register the watcher with the client.
func (r *watchRegistry) Acquire(watcher *lsproto.FileSystemWatcher, id WatcherID) (isNew bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := toFileSystemWatcherKey(watcher)
	value := r.entries[key]
	if value == nil {
		value = &fileSystemWatcherValue{id: id}
		r.entries[key] = value
	}
	value.count++
	return value.count == 1
}

// Release decrements the ref count for a watcher. If no references remain,
// the entry is removed and the function returns the WatcherID and true so
// the caller knows to unregister the watcher from the client.
func (r *watchRegistry) Release(watcher *lsproto.FileSystemWatcher) (id WatcherID, removed bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := toFileSystemWatcherKey(watcher)
	value := r.entries[key]
	if value == nil {
		return "", false
	}
	if value.count <= 1 {
		delete(r.entries, key)
		return value.id, true
	}
	value.count--
	return "", false
}

// MarkPending records that a watcher's registration failed and needs retry.
func (r *watchRegistry) MarkPending(id WatcherID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pending[id] = struct{}{}
}

// ClearPending removes a watcher from the pending set after successful registration.
func (r *watchRegistry) ClearPending(id WatcherID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.pending, id)
}

// IsPending returns true if the watcher needs retry due to a previous failure.
func (r *watchRegistry) IsPending(id WatcherID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.pending[id]
	return ok
}

type PatternsAndIgnored struct {
	directoriesOutsideWorkspace []string
	patternsInsideWorkspace     []string
	ignored                     map[string]struct{}
}

// toFileSystemWatcherKey produces a deduplication key for a file system watcher.
// Note: this key is a simple string concatenation of the base and pattern, so
// structurally different watchers (Pattern vs RelativePattern, URI vs WorkspaceFolder)
// could theoretically collide. In practice, workspace watchers use plain Pattern
// with filesystem paths while outside-workspace watchers use RelativePattern with
// file:// URIs, so collisions don't occur.
func toFileSystemWatcherKey(w *lsproto.FileSystemWatcher) fileSystemWatcherKey {
	kind := w.Kind
	if kind == nil {
		kind = new(lsproto.WatchKindCreate | lsproto.WatchKindChange | lsproto.WatchKindDelete)
	}
	var pattern string
	if w.GlobPattern.Pattern != nil {
		pattern = *w.GlobPattern.Pattern
	} else if w.GlobPattern.RelativePattern != nil {
		var base string
		if w.GlobPattern.RelativePattern.BaseUri.URI != nil {
			base = string(*w.GlobPattern.RelativePattern.BaseUri.URI)
		} else if w.GlobPattern.RelativePattern.BaseUri.WorkspaceFolder != nil {
			panic("workspace folder-based relative patterns not implemented")
		}
		pattern = base + "/" + w.GlobPattern.RelativePattern.Pattern
	}
	return fileSystemWatcherKey{pattern: pattern, kind: *kind}
}

func fileSystemWatcherGlobString(w *lsproto.FileSystemWatcher) string {
	if w.GlobPattern.Pattern != nil {
		return *w.GlobPattern.Pattern
	}
	if w.GlobPattern.RelativePattern != nil {
		var base string
		if w.GlobPattern.RelativePattern.BaseUri.URI != nil {
			base = string(*w.GlobPattern.RelativePattern.BaseUri.URI)
		} else if w.GlobPattern.RelativePattern.BaseUri.WorkspaceFolder != nil {
			panic("workspace folder-based relative patterns not implemented")
		}
		return base + "/" + w.GlobPattern.RelativePattern.Pattern
	}
	return ""
}

type WatcherID string

var watcherID atomic.Uint64

type WatchedFiles[T any] struct {
	name                         string
	watchKind                    lsproto.WatchKind
	hasRelativePatternCapability bool
	granularWatches              bool
	// computeGlobPatterns computes the set of watcher glob patterns from the
	// input. The provided fs must be the live session file system, not a
	// snapshot-cached FS: it is used to probe which directories currently
	// exist on disk so that only safe-to-watch directories are registered.
	computeGlobPatterns func(input T, fs vfs.FS) PatternsAndIgnored

	mu                       sync.RWMutex
	input                    T
	computeWatchersOnce      sync.Once
	workspaceWatchers        []*lsproto.FileSystemWatcher
	outsideWorkspaceWatchers []*lsproto.FileSystemWatcher
	ignored                  map[string]struct{}
	id                       uint64
}

func NewWatchedFiles[T any](name string, watchKind lsproto.WatchKind, hasRelativePatternCapability bool, granularWatches bool, computeGlobPatterns func(input T, fs vfs.FS) PatternsAndIgnored) *WatchedFiles[T] {
	return &WatchedFiles[T]{
		id:                           watcherID.Add(1),
		name:                         name,
		watchKind:                    watchKind,
		hasRelativePatternCapability: hasRelativePatternCapability,
		granularWatches:              granularWatches,
		computeGlobPatterns:          computeGlobPatterns,
	}
}

type Watchers struct {
	WatcherID                WatcherID
	WorkspaceWatchers        []*lsproto.FileSystemWatcher
	OutsideWorkspaceWatchers []*lsproto.FileSystemWatcher
	IgnoredPaths             map[string]struct{}
}

// Watchers computes the watcher set for these watched files. The provided fs
// must be the live session file system, not a snapshot-cached FS: it is used to
// probe which directories currently exist on disk so that only safe-to-watch
// directories are registered.
func (w *WatchedFiles[T]) Watchers(fs vfs.FS) Watchers {
	w.computeWatchersOnce.Do(func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		result := w.computeGlobPatterns(w.input, fs)
		globs := slices.Compact(slices.Sorted(slices.Values(result.patternsInsideWorkspace)))

		ignored := result.ignored
		// ignored is only used for logging and doesn't affect watcher identity
		w.ignored = ignored
		changed := false
		if !slices.EqualFunc(w.workspaceWatchers, globs, func(a *lsproto.FileSystemWatcher, b string) bool {
			return *a.GlobPattern.Pattern == b
		}) {
			w.workspaceWatchers = core.Map(globs, func(glob string) *lsproto.FileSystemWatcher {
				return &lsproto.FileSystemWatcher{
					GlobPattern: lsproto.PatternOrRelativePattern{
						Pattern: &glob,
					},
					Kind: &w.watchKind,
				}
			})
			changed = true
		}
		dirsOutside := slices.Compact(slices.Sorted(slices.Values(result.directoriesOutsideWorkspace)))
		if !slices.EqualFunc(w.outsideWorkspaceWatchers, dirsOutside, func(a *lsproto.FileSystemWatcher, b string) bool {
			return fileSystemWatcherGlobString(a) == recursiveDirectoryGlobPattern(b, w.hasRelativePatternCapability)
		}) {
			w.outsideWorkspaceWatchers = core.Map(dirsOutside, func(dir string) *lsproto.FileSystemWatcher {
				return newRecursiveDirectoryWatcher(dir, w.watchKind, w.hasRelativePatternCapability)
			})
			changed = true
		}
		if changed {
			w.id = watcherID.Add(1)
		}
	})

	w.mu.RLock()
	defer w.mu.RUnlock()
	return Watchers{
		WatcherID:                WatcherID(fmt.Sprintf("%s watcher %d", w.name, w.id)),
		WorkspaceWatchers:        w.workspaceWatchers,
		OutsideWorkspaceWatchers: w.outsideWorkspaceWatchers,
		IgnoredPaths:             w.ignored,
	}
}

// ID returns the identity of the current watcher set. Because computing the
// watchers can change the identity (when the resulting globs change), the
// provided fs must be the live session file system, not a snapshot-cached FS.
func (w *WatchedFiles[T]) ID(fs vfs.FS) WatcherID {
	if w == nil {
		return ""
	}
	return w.Watchers(fs).WatcherID
}

func (w *WatchedFiles[T]) Name() string {
	return w.name
}

func (w *WatchedFiles[T]) WatchKind() lsproto.WatchKind {
	return w.watchKind
}

func (w *WatchedFiles[T]) Clone(input T) *WatchedFiles[T] {
	if w == nil {
		return nil
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return &WatchedFiles[T]{
		name:                         w.name,
		watchKind:                    w.watchKind,
		hasRelativePatternCapability: w.hasRelativePatternCapability,
		granularWatches:              w.granularWatches,
		computeGlobPatterns:          w.computeGlobPatterns,
		workspaceWatchers:            w.workspaceWatchers,
		outsideWorkspaceWatchers:     w.outsideWorkspaceWatchers,
		input:                        input,
	}
}

func createResolutionLookupGlobMapper(
	workspaceDirectory string,
	libDirectory string,
	currentDirectory string,
	useCaseSensitiveFileNames bool,
	granularWatches bool,
) func(data *collections.SyncSet[tspath.Path], fs vfs.FS) PatternsAndIgnored {
	workspaceDirectoryPath := tspath.ToPath(workspaceDirectory, currentDirectory, useCaseSensitiveFileNames)
	currentDirectoryPath := tspath.ToPath(currentDirectory, currentDirectory, useCaseSensitiveFileNames)
	libDirectoryPath := tspath.ToPath(libDirectory, currentDirectory, useCaseSensitiveFileNames)

	// fs must be the live session file system, not a snapshot-cached FS: it is
	// used to probe which directories currently exist on disk.
	return func(data *collections.SyncSet[tspath.Path], fs vfs.FS) PatternsAndIgnored {
		var ignored map[string]struct{}
		var seenDirs collections.Set[tspath.Path]
		var includeWorkspace, includeRoot, includeLib bool
		var workspaceDirectories collections.Set[tspath.Path]
		var rootDirectories collections.Set[tspath.Path]
		var libDirectories collections.Set[tspath.Path]
		var nodeModulesDirectories collections.Set[tspath.Path]
		var externalDirectories collections.Set[tspath.Path]

		if data != nil {
			data.Range(func(path tspath.Path) bool {
				if tspath.IsDynamicFileName(string(path)) {
					return true
				}
				// Assuming all of the input paths are file paths, we can avoid
				// duplicate work by only taking one file per dir, since their outputs
				// will always be the same.
				if !seenDirs.AddIfAbsent(path.GetDirectoryPath()) {
					return true
				}

				if workspaceDirectoryPath.ContainsPath(path) {
					if granularWatches {
						workspaceDirectories.Add(path.GetDirectoryPath())
					} else {
						includeWorkspace = true
					}
				} else if currentDirectoryPath.ContainsPath(path) {
					if granularWatches {
						rootDirectories.Add(path.GetDirectoryPath())
					} else {
						includeRoot = true
					}
				} else if libDirectoryPath.ContainsPath(path) {
					if granularWatches {
						libDirectories.Add(path.GetDirectoryPath())
					} else {
						includeLib = true
					}
				} else if idx := strings.Index(string(path), "/node_modules/"); idx != -1 {
					if granularWatches {
						nodeModulesDirectories.Add(path.GetDirectoryPath())
					} else {
						nodeModulesDirectories.Add(path[:idx+len("/node_modules")])
					}
				} else {
					externalDirectories.Add(path.GetDirectoryPath())
				}
				return true
			})
		}

		var globs []string
		if granularWatches {
			// Granular watches register the specific directory of each probed
			// (often missing) file, non-recursively. We deliberately do not
			// consolidate to recursive watches on the workspace/root/node_modules
			// directories: that would collapse back into the broad watches this
			// mode exists to avoid. Directories that don't exist yet are emitted
			// as-is; the watcher backend is responsible for watching into
			// not-yet-existing trees.
			globs = appendDirectoryGlobs(globs, workspaceDirectories)
			globs = appendDirectoryGlobs(globs, rootDirectories)
			globs = appendDirectoryGlobs(globs, libDirectories)
			globs = appendDirectoryGlobs(globs, nodeModulesDirectories)
		} else {
			if includeWorkspace {
				globs = append(globs, getRecursiveGlobPattern(string(workspaceDirectoryPath)))
			}
			if includeRoot {
				globs = append(globs, getRecursiveGlobPattern(string(currentDirectoryPath)))
			}
			if includeLib {
				globs = append(globs, getRecursiveGlobPattern(string(libDirectoryPath)))
			}
			if nodeModulesDirectories.Len() > 0 {
				nodeModulesGlobs := make([]string, 0, nodeModulesDirectories.Len())
				for dir := range nodeModulesDirectories.Keys() {
					nodeModulesGlobs = append(nodeModulesGlobs, getRecursiveGlobPattern(string(dir)))
				}
				slices.Sort(nodeModulesGlobs)
				globs = append(globs, nodeModulesGlobs...)
			}
		}

		slices.Sort(globs)
		globs = slices.Compact(globs)

		var outsideDirs []string
		if externalDirectories.Len() > 0 {
			externalDirStrings := make([]string, 0, externalDirectories.Len())
			for dir := range externalDirectories.Keys() {
				externalDirStrings = append(externalDirStrings, string(dir))
			}
			if granularWatches {
				var externalDirs collections.Set[tspath.Path]
				for dir := range externalDirectories.Keys() {
					externalDirs.Add(dir)
				}
				outsideDirs = appendResolvedDirectories(outsideDirs, externalDirs, fs.DirectoryExists)
				slices.Sort(outsideDirs)
				outsideDirs = slices.Compact(outsideDirs)
			} else {
				externalDirectoryParents, ignoredExternalDirs := tspath.GetCommonParents(
					externalDirStrings,
					minWatchLocationDepth,
					getPathComponentsForWatching,
					tspath.ComparePathsOptions{UseCaseSensitiveFileNames: true}, // Already using tspath.Path
				)
				slices.Sort(externalDirectoryParents)
				ignored = ignoredExternalDirs
				outsideDirs = externalDirectoryParents
			}
		}

		return PatternsAndIgnored{
			directoriesOutsideWorkspace: outsideDirs,
			patternsInsideWorkspace:     globs,
			ignored:                     ignored,
		}
	}
}

// appendDirectoryGlobs appends a non-recursive glob (`<dir>/*`) for each
// directory. This is used by granular watch mode so that each watch covers only
// the immediate directory of a probed file, rather than its entire subtree.
func appendDirectoryGlobs(globs []string, directories collections.Set[tspath.Path]) []string {
	for dir := range directories.Keys() {
		globs = append(globs, getDirectoryGlobPattern(string(dir)))
	}
	return globs
}

// autoImportWatchGlobs computes the watch globs for the auto-import node_modules
// watcher. In broad mode each node_modules directory is watched recursively
// (`<nm>/**/*`); in granular mode it is watched non-recursively (`<nm>/*`) to
// avoid registering a recursive watch over the entire dependency tree.
func autoImportWatchGlobs(nodeModulesDirs map[tspath.Path]string, granularWatches bool) PatternsAndIgnored {
	patterns := make([]string, 0, len(nodeModulesDirs))
	for _, dir := range nodeModulesDirs {
		if granularWatches {
			patterns = append(patterns, getDirectoryGlobPattern(dir))
		} else {
			patterns = append(patterns, getRecursiveGlobPattern(dir))
		}
	}
	slices.Sort(patterns)
	return PatternsAndIgnored{
		patternsInsideWorkspace: patterns,
	}
}

func appendResolvedDirectories(dirs []string, directories collections.Set[tspath.Path], directoryExists func(path string) bool) []string {
	for dir := range directories.Keys() {
		if resolvedDir, ok := nearestExistingWatchDirectory(dir, directoryExists); ok {
			dirs = append(dirs, string(resolvedDir))
		}
	}
	return dirs
}

func nearestExistingWatchDirectory(dir tspath.Path, directoryExists func(path string) bool) (tspath.Path, bool) {
	if directoryExists == nil {
		return dir, true
	}
	for {
		if directoryExists(string(dir)) {
			return dir, true
		}
		parent := dir.GetDirectoryPath()
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func getTypingsLocationsGlobs(
	typingsFiles []string,
	typingsLocation string,
	workspaceDirectory string,
	currentDirectory string,
	useCaseSensitiveFileNames bool,
) PatternsAndIgnored {
	var includeTypingsLocation, includeWorkspace bool
	externalDirectories := make(map[tspath.Path]string)
	globs := make(map[tspath.Path]string)
	comparePathsOptions := tspath.ComparePathsOptions{
		CurrentDirectory:          currentDirectory,
		UseCaseSensitiveFileNames: useCaseSensitiveFileNames,
	}
	for _, file := range typingsFiles {
		if tspath.ContainsPath(typingsLocation, file, comparePathsOptions) {
			includeTypingsLocation = true
		} else if !tspath.ContainsPath(workspaceDirectory, file, comparePathsOptions) {
			directory := tspath.GetDirectoryPath(file)
			externalDirectories[tspath.ToPath(directory, currentDirectory, useCaseSensitiveFileNames)] = directory
		} else {
			includeWorkspace = true
		}
	}
	externalDirectoryParents, ignored := tspath.GetCommonParents(
		slices.Collect(maps.Values(externalDirectories)),
		minWatchLocationDepth,
		getPathComponentsForWatching,
		comparePathsOptions,
	)
	slices.Sort(externalDirectoryParents)
	if includeWorkspace {
		globs[tspath.ToPath(workspaceDirectory, currentDirectory, useCaseSensitiveFileNames)] = getRecursiveGlobPattern(workspaceDirectory)
	}
	if includeTypingsLocation {
		globs[tspath.ToPath(typingsLocation, currentDirectory, useCaseSensitiveFileNames)] = getRecursiveGlobPattern(typingsLocation)
	}
	return PatternsAndIgnored{
		directoriesOutsideWorkspace: externalDirectoryParents,
		patternsInsideWorkspace:     slices.Collect(maps.Values(globs)),
		ignored:                     ignored,
	}
}

func getPathComponentsForWatching(path string, currentDirectory string) []string {
	components := tspath.GetPathComponents(path, currentDirectory)
	rootLength := perceivedOsRootLengthForWatching(components)
	if rootLength <= 1 {
		return components
	}
	newRoot := tspath.CombinePaths(components[0], components[1:rootLength]...)
	return append([]string{newRoot}, components[rootLength:]...)
}

func perceivedOsRootLengthForWatching(pathComponents []string) int {
	length := len(pathComponents)
	if length <= 1 {
		return length
	}
	if strings.HasPrefix(pathComponents[0], "//") {
		// Group UNC roots (//server/share) into a single component
		return 2
	}
	if len(pathComponents[0]) == 3 && tspath.IsVolumeCharacter(pathComponents[0][0]) && pathComponents[0][1] == ':' && pathComponents[0][2] == '/' {
		// Windows-style volume
		if strings.EqualFold(pathComponents[1], "users") {
			// Group C:/Users/username into a single component
			return min(3, length)
		}
		return 1
	}
	if pathComponents[1] == "home" {
		// Group /home/username into a single component
		return min(3, length)
	}
	return 1
}

func getRecursiveGlobPattern(directory string) string {
	return fmt.Sprintf("%s/%s", tspath.RemoveTrailingDirectorySeparator(directory), "**/*")
}

// getDirectoryGlobPattern returns a non-recursive glob matching only the
// immediate children of directory (`<dir>/*`).
func getDirectoryGlobPattern(directory string) string {
	return fmt.Sprintf("%s/%s", tspath.RemoveTrailingDirectorySeparator(directory), "*")
}

// recursiveDirectoryGlobPattern returns the string form of a recursive watcher
// for the given directory that would be produced by newRecursiveDirectoryWatcher.
func recursiveDirectoryGlobPattern(directory string, useRelativePattern bool) string {
	if useRelativePattern {
		return string(lsconv.FileNameToDocumentURI(directory)) + "/**/*"
	}
	return getRecursiveGlobPattern(directory)
}

// newRecursiveDirectoryWatcher creates a FileSystemWatcher for recursively
// watching a directory. When useRelativePattern is true, a RelativePattern with
// a file:// base URI is used; otherwise a plain glob Pattern is used.
func newRecursiveDirectoryWatcher(directory string, kind lsproto.WatchKind, useRelativePattern bool) *lsproto.FileSystemWatcher {
	if useRelativePattern {
		baseUri := lsproto.URI(lsconv.FileNameToDocumentURI(directory))
		return &lsproto.FileSystemWatcher{
			GlobPattern: lsproto.PatternOrRelativePattern{
				RelativePattern: &lsproto.RelativePattern{
					BaseUri: lsproto.WorkspaceFolderOrURI{
						URI: &baseUri,
					},
					Pattern: "**/*",
				},
			},
			Kind: &kind,
		}
	}
	glob := getRecursiveGlobPattern(directory)
	return &lsproto.FileSystemWatcher{
		GlobPattern: lsproto.PatternOrRelativePattern{
			Pattern: &glob,
		},
		Kind: &kind,
	}
}

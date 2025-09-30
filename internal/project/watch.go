package project

import (
	"fmt"
	"maps"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/glob"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tspath"
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

type patternsAndIgnored struct {
	patterns []string
	ignored  map[string]struct{}
}

func toFileSystemWatcherKey(w *lsproto.FileSystemWatcher) fileSystemWatcherKey {
	if w.GlobPattern.RelativePattern != nil {
		panic("relative globs not implemented")
	}
	kind := w.Kind
	if kind == nil {
		kind = ptrTo(lsproto.WatchKindCreate | lsproto.WatchKindChange | lsproto.WatchKindDelete)
	}
	return fileSystemWatcherKey{pattern: *w.GlobPattern.Pattern, kind: *kind}
}

type WatcherID string

var watcherID atomic.Uint64

type WatchedFiles[T any] struct {
	name                string
	watchKind           lsproto.WatchKind
	computeGlobPatterns func(input T) patternsAndIgnored

	input                  T
	computeWatchersOnce    sync.Once
	watchers               []*lsproto.FileSystemWatcher
	ignored                map[string]struct{}
	computeParsedGlobsOnce sync.Once
	parsedGlobs            []*glob.Glob
	id                     uint64
}

func NewWatchedFiles[T any](name string, watchKind lsproto.WatchKind, computeGlobPatterns func(input T) patternsAndIgnored) *WatchedFiles[T] {
	return &WatchedFiles[T]{
		id:                  watcherID.Add(1),
		name:                name,
		watchKind:           watchKind,
		computeGlobPatterns: computeGlobPatterns,
	}
}

func (w *WatchedFiles[T]) Watchers() (WatcherID, []*lsproto.FileSystemWatcher) {
	w.computeWatchersOnce.Do(func() {
		result := w.computeGlobPatterns(w.input)
		globs := result.patterns
		ignored := result.ignored
		// ignored is only used for logging and doesn't affect watcher identity
		w.ignored = ignored
		if !slices.EqualFunc(w.watchers, globs, func(a *lsproto.FileSystemWatcher, b string) bool {
			return *a.GlobPattern.Pattern == b
		}) {
			w.watchers = core.Map(globs, func(glob string) *lsproto.FileSystemWatcher {
				return &lsproto.FileSystemWatcher{
					GlobPattern: lsproto.PatternOrRelativePattern{
						Pattern: &glob,
					},
					Kind: &w.watchKind,
				}
			})
			w.id = watcherID.Add(1)
		}
	})
	return WatcherID(fmt.Sprintf("%s watcher %d", w.name, w.id)), w.watchers
}

func (w *WatchedFiles[T]) ID() WatcherID {
	if w == nil {
		return ""
	}
	id, _ := w.Watchers()
	return id
}

func (w *WatchedFiles[T]) Name() string {
	return w.name
}

func (w *WatchedFiles[T]) WatchKind() lsproto.WatchKind {
	return w.watchKind
}

func (w *WatchedFiles[T]) ParsedGlobs() []*glob.Glob {
	w.computeParsedGlobsOnce.Do(func() {
		_, watchers := w.Watchers()
		w.parsedGlobs = make([]*glob.Glob, 0, len(watchers))
		for _, watcher := range watchers {
			if g, err := glob.Parse(*watcher.GlobPattern.Pattern); err == nil {
				w.parsedGlobs = append(w.parsedGlobs, g)
			} else {
				panic("failed to parse glob pattern: " + *watcher.GlobPattern.Pattern)
			}
		}
	})
	return w.parsedGlobs
}

func (w *WatchedFiles[T]) Clone(input T) *WatchedFiles[T] {
	return &WatchedFiles[T]{
		name:                w.name,
		watchKind:           w.watchKind,
		computeGlobPatterns: w.computeGlobPatterns,
		watchers:            w.watchers,
		input:               input,
		parsedGlobs:         w.parsedGlobs,
	}
}

func createResolutionLookupGlobMapper(workspaceDirectory string, currentDirectory string, useCaseSensitiveFileNames bool) func(data map[tspath.Path]string) patternsAndIgnored {
	rootPath := tspath.ToPath(currentDirectory, "", useCaseSensitiveFileNames)
	rootPathComponents := tspath.GetPathComponents(string(rootPath), "")
	comparePathsOptions := tspath.ComparePathsOptions{
		CurrentDirectory:          currentDirectory,
		UseCaseSensitiveFileNames: useCaseSensitiveFileNames,
	}

	return func(data map[tspath.Path]string) patternsAndIgnored {
		var ignored map[string]struct{}
		var seenDirs collections.Set[string]
		var includeWorkspace, includeRoot bool
		var externalDirectories map[tspath.Path]string

		for path, fileName := range data {
			// Assuming all of the input paths are filenames, we can avoid
			// duplicate work by only taking one file per dir, since their outputs
			// will always be the same.
			if !seenDirs.AddIfAbsent(tspath.GetDirectoryPath(string(path))) {
				continue
			}

			if tspath.ContainsPath(workspaceDirectory, fileName, comparePathsOptions) {
				includeWorkspace = true
				continue
			} else if tspath.ContainsPath(rootPathComponents[0], fileName, comparePathsOptions) {
				includeRoot = true
				continue
			} else {
				if externalDirectories == nil {
					externalDirectories = make(map[tspath.Path]string)
				}
				externalDirectories[path.GetDirectoryPath()] = tspath.GetDirectoryPath(fileName)
			}
		}

		var globs []string
		if includeWorkspace {
			globs = append(globs, getRecursiveGlobPattern(workspaceDirectory))
		}
		if includeRoot {
			globs = append(globs, getRecursiveGlobPattern(currentDirectory))
		}
		if len(externalDirectories) > 0 {
			externalDirectoryParents, ignoredExternalDirs := tspath.GetCommonParents(slices.Collect(maps.Values(externalDirectories)), minWatchLocationDepth, comparePathsOptions)
			slices.Sort(externalDirectoryParents)
			ignored = ignoredExternalDirs
			for _, dir := range externalDirectoryParents {
				globs = append(globs, getRecursiveGlobPattern(dir))
			}
		}

		return patternsAndIgnored{
			patterns: globs,
			ignored:  ignored,
		}
	}
}

func getTypingsLocationsGlobs(
	typingsFiles []string,
	typingsLocation string,
	workspaceDirectory string,
	currentDirectory string,
	useCaseSensitiveFileNames bool,
) patternsAndIgnored {
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
	externalDirectoryParents, ignored := tspath.GetCommonParents(slices.Collect(maps.Values(externalDirectories)), minWatchLocationDepth, comparePathsOptions)
	slices.Sort(externalDirectoryParents)
	if includeWorkspace {
		globs[tspath.ToPath(workspaceDirectory, currentDirectory, useCaseSensitiveFileNames)] = getRecursiveGlobPattern(workspaceDirectory)
	}
	if includeTypingsLocation {
		globs[tspath.ToPath(typingsLocation, currentDirectory, useCaseSensitiveFileNames)] = getRecursiveGlobPattern(typingsLocation)
	}
	for _, dir := range externalDirectoryParents {
		globs[tspath.ToPath(dir, currentDirectory, useCaseSensitiveFileNames)] = getRecursiveGlobPattern(dir)
	}
	return patternsAndIgnored{
		patterns: slices.Collect(maps.Values(globs)),
		ignored:  ignored,
	}
}

func ptrTo[T any](v T) *T {
	return &v
}

type resolutionWithLookupLocations interface {
	GetLookupLocations() *module.LookupLocations
}

func extractLookups[T resolutionWithLookupLocations](
	projectToPath func(string) tspath.Path,
	failedLookups map[tspath.Path]string,
	affectingLocations map[tspath.Path]string,
	cache map[tspath.Path]module.ModeAwareCache[T],
) {
	for _, resolvedModulesInFile := range cache {
		for _, resolvedModule := range resolvedModulesInFile {
			for _, failedLookupLocation := range resolvedModule.GetLookupLocations().FailedLookupLocations {
				path := projectToPath(failedLookupLocation)
				if _, ok := failedLookups[path]; !ok {
					failedLookups[path] = failedLookupLocation
				}
			}
			for _, affectingLocation := range resolvedModule.GetLookupLocations().AffectingLocations {
				path := projectToPath(affectingLocation)
				if _, ok := affectingLocations[path]; !ok {
					affectingLocations[path] = affectingLocation
				}
			}
		}
	}
}

func getNonRootFileGlobs(workspaceDir string, sourceFiles []*ast.SourceFile, rootFiles map[tspath.Path]string, comparePathsOptions tspath.ComparePathsOptions) patternsAndIgnored {
	var globs []string
	var includeWorkspace bool
	var ignored map[string]struct{}
	externalDirectories := make([]string, 0, max(0, len(sourceFiles)-len(rootFiles)))
	for _, sourceFile := range sourceFiles {
		if _, ok := rootFiles[sourceFile.Path()]; !ok {
			if tspath.ContainsPath(workspaceDir, sourceFile.FileName(), comparePathsOptions) {
				includeWorkspace = true
				continue
			}
			externalDirectories = append(externalDirectories, tspath.GetDirectoryPath(sourceFile.FileName()))
		}
	}

	if includeWorkspace {
		globs = append(globs, getRecursiveGlobPattern(workspaceDir))
	}
	if len(externalDirectories) > 0 {
		commonParents, ignoredDirs := tspath.GetCommonParents(externalDirectories, minWatchLocationDepth, comparePathsOptions)
		globs = append(globs, core.Map(commonParents, func(dir string) string {
			return getRecursiveGlobPattern(dir)
		})...)
		ignored = ignoredDirs
	}
	return patternsAndIgnored{
		patterns: globs,
		ignored:  ignored,
	}
}

func getRecursiveGlobPattern(directory string) string {
	return fmt.Sprintf("%s/%s", tspath.RemoveTrailingDirectorySeparator(directory), "**/*.{js,jsx,mjs,cjs,ts,tsx,mts,cts,json}")
}

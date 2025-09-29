package project

import (
	"fmt"
	"maps"
	"slices"
	"strings"
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
		newWatchers := core.Map(globs, func(glob string) *lsproto.FileSystemWatcher {
			return &lsproto.FileSystemWatcher{
				GlobPattern: lsproto.PatternOrRelativePattern{
					Pattern: &glob,
				},
				Kind: &w.watchKind,
			}
		})
		if !slices.EqualFunc(w.watchers, newWatchers, func(a, b *lsproto.FileSystemWatcher) bool {
			return *a.GlobPattern.Pattern == *b.GlobPattern.Pattern
		}) {
			w.watchers = newWatchers
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
		input:               input,
		parsedGlobs:         w.parsedGlobs,
	}
}

func createResolutionLookupGlobMapper(workspaceDirectory string, currentDirectory string, useCaseSensitiveFileNames bool) func(data map[tspath.Path]string) patternsAndIgnored {
	isWorkspaceWatchable := canWatchDirectoryOrFile(tspath.GetPathComponents(workspaceDirectory, ""))
	rootPath := tspath.ToPath(currentDirectory, "", useCaseSensitiveFileNames)
	rootPathComponents := tspath.GetPathComponents(string(rootPath), "")
	isRootWatchable := canWatchDirectoryOrFile(rootPathComponents)
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

			if isWorkspaceWatchable && tspath.ContainsPath(workspaceDirectory, fileName, comparePathsOptions) {
				includeWorkspace = true
				continue
			} else if isRootWatchable && tspath.ContainsPath(rootPathComponents[0], fileName, comparePathsOptions) {
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
			externalDirectoryParents, ignoredExternalDirs := tspath.GetCommonParents(slices.Collect(maps.Values(externalDirectories)), minWatchLocationDepth, getPathComponentsForWatching, comparePathsOptions)
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

func getPathComponentsForWatching(path string, currentDirectory string) []string {
	components := tspath.GetPathComponents(path, currentDirectory)
	rootLength := perceivedOsRootLengthForWatching(components)
	newRoot := tspath.CombinePaths(components[0], components[1:rootLength]...)
	return append([]string{newRoot}, components[rootLength:]...)
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
	isWorkspaceWatchable := canWatchDirectoryOrFile(tspath.GetPathComponents(workspaceDirectory, ""))
	globs := make(map[tspath.Path]string)
	comparePathsOptions := tspath.ComparePathsOptions{
		CurrentDirectory:          currentDirectory,
		UseCaseSensitiveFileNames: useCaseSensitiveFileNames,
	}
	for _, file := range typingsFiles {
		if tspath.ContainsPath(typingsLocation, file, comparePathsOptions) {
			includeTypingsLocation = true
		} else if !isWorkspaceWatchable || !tspath.ContainsPath(workspaceDirectory, file, comparePathsOptions) {
			directory := tspath.GetDirectoryPath(file)
			externalDirectories[tspath.ToPath(directory, currentDirectory, useCaseSensitiveFileNames)] = directory
		} else {
			includeWorkspace = true
		}
	}
	externalDirectoryParents, ignored := tspath.GetCommonParents(slices.Collect(maps.Values(externalDirectories)), minWatchLocationDepth, getPathComponentsForWatching, comparePathsOptions)
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

func perceivedOsRootLengthForWatching(pathComponents []string) int {
	length := len(pathComponents)
	// Ignore "/", "c:/"
	if length <= 1 {
		return 1
	}
	indexAfterOsRoot := 1
	firstComponent := pathComponents[0]
	isDosStyle := len(firstComponent) >= 2 && tspath.IsVolumeCharacter(firstComponent[0]) && firstComponent[1] == ':'
	if firstComponent != "/" && !isDosStyle && isDosStyleNextPart(pathComponents[1]) {
		// ignore "//vda1cs4850/c$/folderAtRoot"
		if length == 2 {
			return 2
		}
		indexAfterOsRoot = 2
		isDosStyle = true
	}

	afterOsRoot := pathComponents[indexAfterOsRoot]
	if isDosStyle && !strings.EqualFold(afterOsRoot, "users") {
		// Paths like c:/notUsers
		return indexAfterOsRoot
	}

	if strings.EqualFold(afterOsRoot, "workspaces") {
		// Paths like: /workspaces as codespaces hoist the repos in /workspaces so we have to exempt these from "2" level from root rule
		return indexAfterOsRoot + 1
	}

	// Paths like: c:/users/username or /home/username
	return indexAfterOsRoot + 2
}

func canWatchDirectoryOrFile(pathComponents []string) bool {
	length := len(pathComponents)
	// Ignore "/", "c:/"
	// ignore "/user", "c:/users" or "c:/folderAtRoot"
	if length < minWatchLocationDepth {
		return false
	}
	perceivedOsRootLength := perceivedOsRootLengthForWatching(pathComponents)
	return (length - perceivedOsRootLength) >= minWatchLocationDepth
}

func isDosStyleNextPart(part string) bool {
	return len(part) == 2 && tspath.IsVolumeCharacter(part[0]) && part[1] == '$'
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
	canWatchWorkspace := canWatchDirectoryOrFile(tspath.GetPathComponents(workspaceDir, ""))
	externalDirectories := make([]string, 0, max(0, len(sourceFiles)-len(rootFiles)))
	for _, sourceFile := range sourceFiles {
		if _, ok := rootFiles[sourceFile.Path()]; !ok {
			if canWatchWorkspace && tspath.ContainsPath(workspaceDir, sourceFile.FileName(), comparePathsOptions) {
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
		commonParents, ignoredDirs := tspath.GetCommonParents(externalDirectories, minWatchLocationDepth, getPathComponentsForWatching, comparePathsOptions)
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
	return fmt.Sprintf("%s/%s", directory, "**/*.{js,jsx,mjs,cjs,ts,tsx,mts,cts,json}")
}

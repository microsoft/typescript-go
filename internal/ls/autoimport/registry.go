package autoimport

import (
	"cmp"
	"context"
	"maps"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/project/dirty"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type newProgramStructure int

const (
	newProgramStructureFalse newProgramStructure = iota
	newProgramStructureSameFileNames
	newProgramStructureDifferentFileNames
)

// BucketState represents the dirty state of a bucket.
// In general, a bucket can be used for an auto-imports request if it is clean
// or if the only edited file is the one that was requested for auto-imports.
// Most edits within a file will not change the imports available to that file.
// However, two exceptions cause the bucket to be rebuilt after a change to a
// single file:
//
//  1. Local files are newly added to the project by a manual import
//  2. A node_modules dependency normally filtered out by package.json dependencies
//     is added to the project by a manual import
//
// Both of these cases take a bit of work to determine, but can only happen after
// a full (non-clone) program update. When this happens, the `newProgramStructure`
// flag is set until the next time the bucket is rebuilt, when those conditions
// will be checked.
type BucketState struct {
	// dirtyFile is the file that was edited last, if any. It does not necessarily
	// indicate that no other files have been edited, so it should be ignored if
	// `multipleFilesDirty` is set.
	dirtyFile           tspath.Path
	multipleFilesDirty  bool
	newProgramStructure newProgramStructure
	// fileExcludePatterns is the value of the corresponding user preference when
	// the bucket was built. If changed, the bucket should be rebuilt.
	fileExcludePatterns []string
}

func (b BucketState) Dirty() bool {
	return b.multipleFilesDirty || b.dirtyFile != "" || b.newProgramStructure > 0
}

func (b BucketState) DirtyFile() tspath.Path {
	if b.multipleFilesDirty {
		return ""
	}
	return b.dirtyFile
}

func (b BucketState) possiblyNeedsRebuildForFile(file tspath.Path, preferences *lsutil.UserPreferences) bool {
	return b.newProgramStructure > 0 || b.hasDirtyFileBesides(file) || !core.UnorderedEqual(b.fileExcludePatterns, preferences.AutoImportFileExcludePatterns)
}

func (b BucketState) hasDirtyFileBesides(file tspath.Path) bool {
	return b.multipleFilesDirty || b.dirtyFile != "" && b.dirtyFile != file
}

type RegistryBucket struct {
	state BucketState

	Paths collections.Set[tspath.Path]
	// IgnoredPackageNames is only defined for project buckets. It is the set of
	// package names that were present in the project's program, and not included
	// in a node_modules bucket, and ultimately not included in the project bucket
	// because they were only imported transitively. If an updated program's
	// ResolvedPackageNames contains one of these, the bucket should be rebuilt
	// because that package will be included.
	IgnoredPackageNames *collections.Set[string]
	// PackageNames is only defined for node_modules buckets. It is the full set of
	// package directory names in the node_modules directory (but not necessarily
	// inclued in the bucket).
	PackageNames *collections.Set[string]
	// DependencyNames is only defined for node_modules buckets. It is the set of
	// package names that will be included in the bucket if present in the directory,
	// computed from package.json dependencies. If nil, all packages are included
	// because at least one open file has access to this node_modules directory without
	// being filtered by a package.json.
	DependencyNames *collections.Set[string]
	// AmbientModuleNames is only defined for node_modules buckets. It is the set of
	// ambient module names found while extracting exports in the bucket.
	AmbientModuleNames map[string][]string
	// Entrypoints is only defined for node_modules buckets. Keys are package entrypoint
	// file paths, and values describe the ways of importing the package that would resolve
	// to that file.
	Entrypoints map[tspath.Path][]*module.ResolvedEntrypoint
	Index       *Index[*Export]
}

func newRegistryBucket() *RegistryBucket {
	return &RegistryBucket{
		state: BucketState{
			multipleFilesDirty:  true,
			newProgramStructure: newProgramStructureDifferentFileNames,
		},
	}
}

func (b *RegistryBucket) Clone() *RegistryBucket {
	return &RegistryBucket{
		state:               b.state,
		Paths:               b.Paths,
		IgnoredPackageNames: b.IgnoredPackageNames,
		PackageNames:        b.PackageNames,
		DependencyNames:     b.DependencyNames,
		AmbientModuleNames:  b.AmbientModuleNames,
		Entrypoints:         b.Entrypoints,
		Index:               b.Index,
	}
}

// markFileDirty should only be called within a Change call on the dirty map.
// Buckets are considered immutable once in a finalized registry.
func (b *RegistryBucket) markFileDirty(file tspath.Path) {
	if b.state.hasDirtyFileBesides(file) {
		b.state.multipleFilesDirty = true
	} else {
		b.state.dirtyFile = file
	}
}

type directory struct {
	name           string
	packageJson    *packagejson.InfoCacheEntry
	hasNodeModules bool
}

func (d *directory) Clone() *directory {
	return &directory{
		name:           d.name,
		packageJson:    d.packageJson,
		hasNodeModules: d.hasNodeModules,
	}
}

type Registry struct {
	toPath          func(fileName string) tspath.Path
	userPreferences *lsutil.UserPreferences

	// exports      map[tspath.Path][]*RawExport
	directories map[tspath.Path]*directory

	nodeModules map[tspath.Path]*RegistryBucket
	projects    map[tspath.Path]*RegistryBucket

	// specifierCache maps from importing file to target file to specifier.
	specifierCache map[tspath.Path]*collections.SyncMap[tspath.Path, string]
}

func NewRegistry(toPath func(fileName string) tspath.Path) *Registry {
	return &Registry{
		toPath:      toPath,
		directories: make(map[tspath.Path]*directory),
	}
}

func (r *Registry) IsPreparedForImportingFile(fileName string, projectPath tspath.Path, preferences *lsutil.UserPreferences) bool {
	if r == nil {
		return false
	}
	projectBucket, ok := r.projects[projectPath]
	if !ok {
		panic("project bucket missing")
	}
	path := r.toPath(fileName)
	if projectBucket.state.possiblyNeedsRebuildForFile(path, preferences) {
		return false
	}

	dirPath := path.GetDirectoryPath()
	for {
		if dirBucket, ok := r.nodeModules[dirPath]; ok {
			if dirBucket.state.possiblyNeedsRebuildForFile(path, preferences) {
				return false
			}
		}
		parent := dirPath.GetDirectoryPath()
		if parent == dirPath {
			break
		}
		dirPath = parent
	}
	return true
}

func (r *Registry) NodeModulesDirectories() map[tspath.Path]string {
	dirs := make(map[tspath.Path]string)
	for dirPath, dir := range r.directories {
		if dir.hasNodeModules {
			dirs[tspath.Path(tspath.CombinePaths(string(dirPath), "node_modules"))] = tspath.CombinePaths(dir.name, "node_modules")
		}
	}
	return dirs
}

func (r *Registry) Clone(ctx context.Context, change RegistryChange, host RegistryCloneHost, logger *logging.LogTree) (*Registry, error) {
	start := time.Now()
	if logger != nil {
		logger = logger.Fork("Building autoimport registry")
	}
	builder := newRegistryBuilder(r, host)
	if change.UserPreferences != nil {
		builder.userPreferences = change.UserPreferences
		if !core.UnorderedEqual(builder.userPreferences.AutoImportSpecifierExcludeRegexes, r.userPreferences.AutoImportSpecifierExcludeRegexes) {
			builder.specifierCache.Clear()
		}
	}
	builder.updateBucketAndDirectoryExistence(change, logger)
	builder.markBucketsDirty(change, logger)
	if change.RequestedFile != "" {
		builder.updateIndexes(ctx, change, logger)
	}
	if logger != nil {
		logger.Logf("Built autoimport registry in %v", time.Since(start))
	}
	registry := builder.Build()
	builder.host.Dispose()
	return registry, nil
}

type BucketStats struct {
	Path            tspath.Path
	ExportCount     int
	FileCount       int
	State           BucketState
	DependencyNames *collections.Set[string]
	PackageNames    *collections.Set[string]
}

type CacheStats struct {
	ProjectBuckets     []BucketStats
	NodeModulesBuckets []BucketStats
}

func (r *Registry) GetCacheStats() *CacheStats {
	stats := &CacheStats{}

	for path, bucket := range r.projects {
		exportCount := 0
		if bucket.Index != nil {
			exportCount = len(bucket.Index.entries)
		}
		stats.ProjectBuckets = append(stats.ProjectBuckets, BucketStats{
			Path:            path,
			ExportCount:     exportCount,
			FileCount:       bucket.Paths.Len(),
			State:           bucket.state,
			DependencyNames: bucket.DependencyNames,
			PackageNames:    bucket.PackageNames,
		})
	}

	for path, bucket := range r.nodeModules {
		exportCount := 0
		if bucket.Index != nil {
			exportCount = len(bucket.Index.entries)
		}
		stats.NodeModulesBuckets = append(stats.NodeModulesBuckets, BucketStats{
			Path:            path,
			ExportCount:     exportCount,
			FileCount:       bucket.Paths.Len(),
			State:           bucket.state,
			DependencyNames: bucket.DependencyNames,
			PackageNames:    bucket.PackageNames,
		})
	}

	slices.SortFunc(stats.ProjectBuckets, func(a, b BucketStats) int {
		return cmp.Compare(a.Path, b.Path)
	})
	slices.SortFunc(stats.NodeModulesBuckets, func(a, b BucketStats) int {
		return cmp.Compare(a.Path, b.Path)
	})

	return stats
}

type RegistryChange struct {
	RequestedFile tspath.Path
	OpenFiles     map[tspath.Path]string
	Changed       collections.Set[lsproto.DocumentUri]
	Created       collections.Set[lsproto.DocumentUri]
	Deleted       collections.Set[lsproto.DocumentUri]
	// RebuiltPrograms maps from project path to:
	//   - true: the program was rebuilt with a different set of file names
	//   - false: the program was rebuilt but the set of file names is unchanged
	RebuiltPrograms map[tspath.Path]bool
	UserPreferences *lsutil.UserPreferences
}

type RegistryCloneHost interface {
	module.ResolutionHost
	FS() vfs.FS
	GetDefaultProject(path tspath.Path) (tspath.Path, *compiler.Program)
	GetProgramForProject(projectPath tspath.Path) *compiler.Program
	GetPackageJson(fileName string) *packagejson.InfoCacheEntry
	GetSourceFile(fileName string, path tspath.Path) *ast.SourceFile
	Dispose()
}

type registryBuilder struct {
	host     RegistryCloneHost
	resolver *module.Resolver
	base     *Registry

	userPreferences *lsutil.UserPreferences
	directories     *dirty.Map[tspath.Path, *directory]
	nodeModules     *dirty.Map[tspath.Path, *RegistryBucket]
	projects        *dirty.Map[tspath.Path, *RegistryBucket]
	specifierCache  *dirty.MapBuilder[tspath.Path, *collections.SyncMap[tspath.Path, string], *collections.SyncMap[tspath.Path, string]]
}

func newRegistryBuilder(registry *Registry, host RegistryCloneHost) *registryBuilder {
	return &registryBuilder{
		host:     host,
		resolver: module.NewResolver(host, core.EmptyCompilerOptions, "", ""),
		base:     registry,

		userPreferences: registry.userPreferences.OrDefault(),
		directories:     dirty.NewMap(registry.directories),
		nodeModules:     dirty.NewMap(registry.nodeModules),
		projects:        dirty.NewMap(registry.projects),
		specifierCache:  dirty.NewMapBuilder(registry.specifierCache, core.Identity, core.Identity),
	}
}

func (b *registryBuilder) Build() *Registry {
	return &Registry{
		toPath:          b.base.toPath,
		userPreferences: b.userPreferences,
		directories:     core.FirstResult(b.directories.Finalize()),
		nodeModules:     core.FirstResult(b.nodeModules.Finalize()),
		projects:        core.FirstResult(b.projects.Finalize()),
		specifierCache:  core.FirstResult(b.specifierCache.Build()),
	}
}

func (b *registryBuilder) updateBucketAndDirectoryExistence(change RegistryChange, logger *logging.LogTree) {
	start := time.Now()
	neededProjects := make(map[tspath.Path]struct{})
	neededDirectories := make(map[tspath.Path]string)
	for path, fileName := range change.OpenFiles {
		neededProjects[core.FirstResult(b.host.GetDefaultProject(path))] = struct{}{}
		if strings.HasPrefix(fileName, "^/") {
			continue
		}
		dir := fileName
		dirPath := path
		for {
			dir = tspath.GetDirectoryPath(dir)
			lastDirPath := dirPath
			dirPath = dirPath.GetDirectoryPath()
			if dirPath == lastDirPath {
				break
			}
			if _, ok := neededDirectories[dirPath]; ok {
				break
			}
			neededDirectories[dirPath] = dir
		}

		if !b.specifierCache.Has(path) {
			b.specifierCache.Set(path, &collections.SyncMap[tspath.Path, string]{})
		}
	}

	for path := range b.base.specifierCache {
		if _, ok := change.OpenFiles[path]; !ok {
			b.specifierCache.Delete(path)
		}
	}

	var addedProjects, removedProjects []tspath.Path
	core.DiffMapsFunc(
		b.base.projects,
		neededProjects,
		func(_ *RegistryBucket, _ struct{}) bool {
			panic("never called because onChanged is nil")
		},
		func(projectPath tspath.Path, _ struct{}) {
			// Need and don't have
			b.projects.Add(projectPath, newRegistryBucket())
			addedProjects = append(addedProjects, projectPath)
		},
		func(projectPath tspath.Path, _ *RegistryBucket) {
			// Have and don't need
			b.projects.Delete(projectPath)
			removedProjects = append(removedProjects, projectPath)
		},
		nil,
	)
	if logger != nil {
		for _, projectPath := range addedProjects {
			logger.Logf("Added project: %s", projectPath)
		}
		for _, projectPath := range removedProjects {
			logger.Logf("Removed project: %s", projectPath)
		}
	}

	updateDirectory := func(dirPath tspath.Path, dirName string, packageJsonChanged bool) {
		packageJsonFileName := tspath.CombinePaths(dirName, "package.json")
		hasNodeModules := b.host.FS().DirectoryExists(tspath.CombinePaths(dirName, "node_modules"))
		if entry, ok := b.directories.Get(dirPath); ok {
			entry.ChangeIf(func(dir *directory) bool {
				return packageJsonChanged || dir.hasNodeModules != hasNodeModules
			}, func(dir *directory) {
				dir.packageJson = b.host.GetPackageJson(packageJsonFileName)
				dir.hasNodeModules = hasNodeModules
			})
		} else {
			b.directories.Add(dirPath, &directory{
				name:           dirName,
				packageJson:    b.host.GetPackageJson(packageJsonFileName),
				hasNodeModules: hasNodeModules,
			})
		}

		if packageJsonChanged {
			// package.json changes affecting node_modules are handled by comparing dependencies in updateIndexes
			return
		}

		if hasNodeModules {
			if _, ok := b.nodeModules.Get(dirPath); !ok {
				b.nodeModules.Add(dirPath, newRegistryBucket())
			}
		} else {
			b.nodeModules.TryDelete(dirPath)
		}
	}

	var addedNodeModulesDirs, removedNodeModulesDirs []tspath.Path
	core.DiffMapsFunc(
		b.base.directories,
		neededDirectories,
		func(dir *directory, dirName string) bool {
			packageJsonUri := lsconv.FileNameToDocumentURI(tspath.CombinePaths(dirName, "package.json"))
			return !change.Changed.Has(packageJsonUri) && !change.Deleted.Has(packageJsonUri) && !change.Created.Has(packageJsonUri)
		},
		func(dirPath tspath.Path, dirName string) {
			// Need and don't have
			hadNodeModules := b.base.nodeModules[dirPath] != nil
			updateDirectory(dirPath, dirName, false)
			if logger != nil {
				logger.Logf("Added directory: %s", dirPath)
			}
			if _, hasNow := b.nodeModules.Get(dirPath); hasNow && !hadNodeModules {
				addedNodeModulesDirs = append(addedNodeModulesDirs, dirPath)
			}
		},
		func(dirPath tspath.Path, dir *directory) {
			// Have and don't need
			hadNodeModules := b.base.nodeModules[dirPath] != nil
			b.directories.Delete(dirPath)
			b.nodeModules.TryDelete(dirPath)
			if logger != nil {
				logger.Logf("Removed directory: %s", dirPath)
			}
			if hadNodeModules {
				removedNodeModulesDirs = append(removedNodeModulesDirs, dirPath)
			}
		},
		func(dirPath tspath.Path, dir *directory, dirName string) {
			// package.json may have changed
			updateDirectory(dirPath, dirName, true)
			if logger != nil {
				logger.Logf("Changed directory: %s", dirPath)
			}
		},
	)
	if logger != nil {
		for _, dirPath := range addedNodeModulesDirs {
			logger.Logf("Added node_modules bucket: %s", dirPath)
		}
		for _, dirPath := range removedNodeModulesDirs {
			logger.Logf("Removed node_modules bucket: %s", dirPath)
		}
		logger.Logf("Updated buckets and directories in %v", time.Since(start))
	}
}

func (b *registryBuilder) markBucketsDirty(change RegistryChange, logger *logging.LogTree) {
	// Mark new program structures
	for projectPath, newFileNames := range change.RebuiltPrograms {
		if bucket, ok := b.projects.Get(projectPath); ok {
			bucket.Change(func(bucket *RegistryBucket) {
				bucket.state.newProgramStructure = core.IfElse(newFileNames, newProgramStructureDifferentFileNames, newProgramStructureSameFileNames)
			})
		}
	}

	// Mark files dirty, bailing out if all buckets already have multiple files dirty
	cleanNodeModulesBuckets := make(map[tspath.Path]struct{})
	cleanProjectBuckets := make(map[tspath.Path]struct{})
	b.nodeModules.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
		if !entry.Value().state.multipleFilesDirty {
			cleanNodeModulesBuckets[entry.Key()] = struct{}{}
		}
		return true
	})
	b.projects.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
		if !entry.Value().state.multipleFilesDirty {
			cleanProjectBuckets[entry.Key()] = struct{}{}
		}
		return true
	})

	markFilesDirty := func(uris map[lsproto.DocumentUri]struct{}) {
		if len(cleanNodeModulesBuckets) == 0 && len(cleanProjectBuckets) == 0 {
			return
		}
		for uri := range uris {
			path := b.base.toPath(uri.FileName())
			if len(cleanNodeModulesBuckets) > 0 {
				// For node_modules, mark the bucket dirty if anything changes in the directory
				if nodeModulesIndex := strings.Index(string(path), "/node_modules/"); nodeModulesIndex != -1 {
					dirPath := path[:nodeModulesIndex]
					if _, ok := cleanNodeModulesBuckets[dirPath]; ok {
						entry := core.FirstResult(b.nodeModules.Get(dirPath))
						entry.Change(func(bucket *RegistryBucket) { bucket.markFileDirty(path) })
						if !entry.Value().state.multipleFilesDirty {
							delete(cleanNodeModulesBuckets, dirPath)
						}
					}
				}
			}
			// For projects, mark the bucket dirty if the bucket contains the file directly.
			// Any other significant change, like a created failed lookup location, is
			// handled by newProgramStructure.
			for projectDirPath := range cleanProjectBuckets {
				entry, _ := b.projects.Get(projectDirPath)
				if entry.Value().Paths.Has(path) {
					entry.Change(func(bucket *RegistryBucket) { bucket.markFileDirty(path) })
					if !entry.Value().state.multipleFilesDirty {
						delete(cleanProjectBuckets, projectDirPath)
					}
				}
			}
		}
	}

	markFilesDirty(change.Created.Keys())
	markFilesDirty(change.Deleted.Keys())
	markFilesDirty(change.Changed.Keys())
}

func (b *registryBuilder) updateIndexes(ctx context.Context, change RegistryChange, logger *logging.LogTree) {
	type task struct {
		entry           *dirty.MapEntry[tspath.Path, *RegistryBucket]
		dependencyNames *collections.Set[string]
		result          *bucketBuildResult
		err             error
	}

	projectPath, _ := b.host.GetDefaultProject(change.RequestedFile)
	if projectPath == "" {
		return
	}

	var tasks []*task
	var wg sync.WaitGroup

	tspath.ForEachAncestorDirectoryPath(change.RequestedFile, func(dirPath tspath.Path) (any, bool) {
		if nodeModulesBucket, ok := b.nodeModules.Get(dirPath); ok {
			dirName := core.FirstResult(b.directories.Get(dirPath)).Value().name
			dependencies := b.computeDependenciesForNodeModulesDirectory(change, dirName, dirPath)
			if nodeModulesBucket.Value().state.hasDirtyFileBesides(change.RequestedFile) || !nodeModulesBucket.Value().DependencyNames.Equals(dependencies) {
				task := &task{entry: nodeModulesBucket, dependencyNames: dependencies}
				tasks = append(tasks, task)
				wg.Go(func() {
					result, err := b.buildNodeModulesBucket(ctx, dependencies, dirName, dirPath, logger.Fork("Building node_modules bucket "+dirName))
					task.result = result
					task.err = err
				})
			}
		}
		return nil, false
	})

	nodeModulesContainsDependency := func(nodeModulesDir tspath.Path, packageName string) bool {
		for _, task := range tasks {
			if task.entry.Key() == nodeModulesDir {
				return task.dependencyNames == nil || task.dependencyNames.Has(packageName)
			}
		}
		if bucket, ok := b.base.nodeModules[nodeModulesDir]; ok {
			return bucket.DependencyNames == nil || bucket.DependencyNames.Has(packageName)
		}
		return false
	}

	if project, ok := b.projects.Get(projectPath); ok {
		program := b.host.GetProgramForProject(projectPath)
		resolvedPackageNames := core.Memoize(func() *collections.Set[string] {
			return getResolvedPackageNames(ctx, program)
		})
		shouldRebuild := project.Value().state.hasDirtyFileBesides(change.RequestedFile)
		if !shouldRebuild && project.Value().state.newProgramStructure > 0 {
			// Exceptions from BucketState comment - check if new program's resolved package names include any
			// previously ignored, or if there are new non-node_modules files.
			// If not, we can skip rebuilding the project bucket.
			if project.Value().IgnoredPackageNames.Intersects(resolvedPackageNames()) || hasNewNonNodeModulesFiles(program, project.Value()) {
				shouldRebuild = true
			} else {
				project.Change(func(b *RegistryBucket) { b.state.newProgramStructure = newProgramStructureFalse })
			}
		}
		if shouldRebuild {
			task := &task{entry: project}
			tasks = append(tasks, task)
			wg.Go(func() {
				index, err := b.buildProjectBucket(
					ctx,
					projectPath,
					resolvedPackageNames(),
					nodeModulesContainsDependency,
					logger.Fork("Building project bucket "+string(projectPath)),
				)
				task.result = index
				task.err = err
			})
		}
	}

	start := time.Now()
	wg.Wait()

	for _, t := range tasks {
		if t.err != nil {
			continue
		}
		t.entry.Replace(t.result.bucket)
	}

	// If we failed to resolve any alias exports by ending up at a non-relative module specifier
	// that didn't resolve to another package, it's probably an ambient module declared in another package.
	// We recorded these failures, along with the name of every ambient module declared elsewhere, so we
	// can do a second pass on the failed files, this time including the ambient modules declarations that
	// were missing the first time. Example: node_modules/fs-extra/index.d.ts is simply `export * from "fs"`,
	// but when trying to resolve the `export *`, we don't know where "fs" is declared. The aliasResolver
	// tries to find packages named "fs" on the file system, but after failing, records "fs" as a failure
	// for fs-extra/index.d.ts. Meanwhile, if we also processed node_modules/@types/node/fs.d.ts, we
	// recorded that file as declaring the ambient module "fs". In the second pass, we combine those two
	// files and reprocess fs-extra/index.d.ts, this time finding "fs" declared in @types/node.
	secondPassStart := time.Now()
	var secondPassFileCount int
	for _, t := range tasks {
		if t.err != nil {
			continue
		}
		if t.result.possibleFailedAmbientModuleLookupTargets == nil {
			continue
		}
		rootFiles := make(map[string]*ast.SourceFile)
		for target := range t.result.possibleFailedAmbientModuleLookupTargets.Keys() {
			for _, fileName := range b.resolveAmbientModuleName(target, t.entry.Key()) {
				if _, exists := rootFiles[fileName]; exists {
					continue
				}
				rootFiles[fileName] = b.host.GetSourceFile(fileName, b.base.toPath(fileName))
				secondPassFileCount++
			}
		}
		if len(rootFiles) > 0 {
			aliasResolver := newAliasResolver(slices.Collect(maps.Values(rootFiles)), b.host, b.resolver, b.base.toPath, func(source ast.HasFileName, moduleName string) {
				// no-op
			})
			ch, _ := checker.NewChecker(aliasResolver)
			t.result.possibleFailedAmbientModuleLookupSources.Range(func(path tspath.Path, source *failedAmbientModuleLookupSource) bool {
				sourceFile := aliasResolver.GetSourceFile(source.fileName)
				extractor := b.newExportExtractor(t.entry.Key(), source.packageName, ch)
				fileExports := extractor.extractFromFile(sourceFile)
				t.result.bucket.Paths.Add(path)
				for _, exp := range fileExports {
					t.result.bucket.Index.insertAsWords(exp)
				}
				return true
			})
		}
	}

	if logger != nil && len(tasks) > 0 {
		if secondPassFileCount > 0 {
			logger.Logf("%d files required second pass, took %v", secondPassFileCount, time.Since(secondPassStart))
		}
		logger.Logf("Built %d indexes in %v", len(tasks), time.Since(start))
	}
}

func hasNewNonNodeModulesFiles(program *compiler.Program, bucket *RegistryBucket) bool {
	if bucket.state.newProgramStructure != newProgramStructureDifferentFileNames {
		return false
	}
	for _, file := range program.GetSourceFiles() {
		if strings.Contains(file.FileName(), "/node_modules/") || isIgnoredFile(program, file) {
			continue
		}
		if !bucket.Paths.Has(file.Path()) {
			return true
		}
	}
	return false
}

func isIgnoredFile(program *compiler.Program, file *ast.SourceFile) bool {
	return program.IsSourceFileDefaultLibrary(file.Path()) || program.IsGlobalTypingsFile(file.FileName())
}

type failedAmbientModuleLookupSource struct {
	mu          sync.Mutex
	fileName    string
	packageName string
}

type bucketBuildResult struct {
	bucket *RegistryBucket
	// File path to filename and package name
	possibleFailedAmbientModuleLookupSources *collections.SyncMap[tspath.Path, *failedAmbientModuleLookupSource]
	// Likely ambient module name
	possibleFailedAmbientModuleLookupTargets *collections.SyncSet[string]
}

func (b *registryBuilder) buildProjectBucket(
	ctx context.Context,
	projectPath tspath.Path,
	resolvedPackageNames *collections.Set[string],
	nodeModulesContainsDependency func(nodeModulesDir tspath.Path, packageName string) bool,
	logger *logging.LogTree,
) (*bucketBuildResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	start := time.Now()
	var mu sync.Mutex
	fileExcludePatterns := b.userPreferences.ParsedAutoImportFileExcludePatterns(b.host.FS().UseCaseSensitiveFileNames())
	result := &bucketBuildResult{bucket: &RegistryBucket{}}
	program := b.host.GetProgramForProject(projectPath)
	getChecker, closePool, checkerCount := createCheckerPool(program)
	defer closePool()
	exports := make(map[tspath.Path][]*Export)
	var wg sync.WaitGroup
	var ignoredPackageNames collections.Set[string]
	var skippedFileCount int
	var combinedStats extractorStats

outer:
	for _, file := range program.GetSourceFiles() {
		if isIgnoredFile(program, file) {
			continue
		}
		for _, excludePattern := range fileExcludePatterns {
			if matched, _ := excludePattern.MatchString(file.FileName()); matched {
				skippedFileCount++
				continue outer
			}
		}
		if packageName := modulespecifiers.GetPackageNameFromDirectory(file.FileName()); packageName != "" {
			// Only process this file if it is not going to be processed as part of a node_modules bucket
			// *and* if it was imported directly (not transitively) by a project file (i.e., this is part
			// of a package not listed in package.json, but imported anyway).
			pathComponents := tspath.GetPathComponents(string(file.Path()), "")
			nodeModulesDir := tspath.GetPathFromPathComponents(pathComponents[:slices.Index(pathComponents, "node_modules")])
			if nodeModulesContainsDependency(tspath.Path(nodeModulesDir), packageName) {
				continue
			}
			if !resolvedPackageNames.Has(packageName) {
				ignoredPackageNames.Add(packageName)
				continue
			}
		}
		wg.Go(func() {
			if ctx.Err() == nil {
				checker, done := getChecker()
				defer done()
				extractor := b.newExportExtractor("", "", checker)
				fileExports := extractor.extractFromFile(file)
				mu.Lock()
				exports[file.Path()] = fileExports
				mu.Unlock()
				stats := extractor.Stats()
				combinedStats.exports.Add(stats.exports.Load())
				combinedStats.usedChecker.Add(stats.usedChecker.Load())
			}
		})
	}

	wg.Wait()

	indexStart := time.Now()
	idx := &Index[*Export]{}
	for path, fileExports := range exports {
		result.bucket.Paths.Add(path)
		for _, exp := range fileExports {
			idx.insertAsWords(exp)
		}
	}

	result.bucket.Index = idx
	result.bucket.IgnoredPackageNames = &ignoredPackageNames
	result.bucket.state.fileExcludePatterns = b.userPreferences.AutoImportFileExcludePatterns

	if logger != nil {
		logger.Logf("Extracted exports: %v (%d exports, %d used checker, %d created checkers)", indexStart.Sub(start), combinedStats.exports.Load(), combinedStats.usedChecker.Load(), checkerCount())
		if skippedFileCount > 0 {
			logger.Logf("Skipped %d files due to exclude patterns", skippedFileCount)
		}
		logger.Logf("Built index: %v", time.Since(indexStart))
		logger.Logf("Bucket total: %v", time.Since(start))
	}
	return result, nil
}

func (b *registryBuilder) computeDependenciesForNodeModulesDirectory(change RegistryChange, dirName string, dirPath tspath.Path) *collections.Set[string] {
	// If any open files are in scope of this directory but not in scope of any package.json,
	// we need to add all packages in this node_modules directory.
	for path := range change.OpenFiles {
		if dirPath.ContainsPath(path) && b.getNearestAncestorDirectoryWithValidPackageJson(path) == nil {
			return nil
		}
	}

	// Get all package.jsons that have this node_modules directory in their spine
	dependencies := &collections.Set[string]{}
	b.directories.Range(func(entry *dirty.MapEntry[tspath.Path, *directory]) bool {
		if entry.Value().packageJson.Exists() && dirPath.ContainsPath(entry.Key()) {
			entry.Value().packageJson.Contents.RangeDependencies(func(name, _, field string) bool {
				if field == "dependencies" || field == "peerDendencies" {
					dependencies.Add(module.GetPackageNameFromTypesPackageName(name))
				}
				return true
			})
		}
		return true
	})
	return dependencies
}

func (b *registryBuilder) buildNodeModulesBucket(
	ctx context.Context,
	dependencies *collections.Set[string],
	dirName string,
	dirPath tspath.Path,
	logger *logging.LogTree,
) (*bucketBuildResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	start := time.Now()
	fileExcludePatterns := b.userPreferences.ParsedAutoImportFileExcludePatterns(b.host.FS().UseCaseSensitiveFileNames())
	directoryPackageNames, err := getPackageNamesInNodeModules(tspath.CombinePaths(dirName, "node_modules"), b.host.FS())
	if err != nil {
		return nil, err
	}

	extractorStart := time.Now()
	packageNames := core.Coalesce(dependencies, directoryPackageNames)

	var exportsMu sync.Mutex
	exports := make(map[tspath.Path][]*Export)
	ambientModuleNames := make(map[string][]string)

	var entrypointsMu sync.Mutex
	var entrypoints []*module.ResolvedEntrypoints
	var skippedEntrypointsCount int32
	var combinedStats extractorStats
	var possibleFailedAmbientModuleLookupTargets collections.SyncSet[string]
	var possibleFailedAmbientModuleLookupSources collections.SyncMap[tspath.Path, *failedAmbientModuleLookupSource]

	createAliasResolver := func(packageName string, entrypoints []*module.ResolvedEntrypoint) *aliasResolver {
		rootFiles := make([]*ast.SourceFile, len(entrypoints))
		var wg sync.WaitGroup
		for i, entrypoint := range entrypoints {
			wg.Go(func() {
				file := b.host.GetSourceFile(entrypoint.ResolvedFileName, b.base.toPath(entrypoint.ResolvedFileName))
				binder.BindSourceFile(file)
				rootFiles[i] = file
			})
		}
		wg.Wait()

		rootFiles = slices.DeleteFunc(rootFiles, func(f *ast.SourceFile) bool {
			return f == nil
		})

		return newAliasResolver(rootFiles, b.host, b.resolver, b.base.toPath, func(source ast.HasFileName, moduleName string) {
			possibleFailedAmbientModuleLookupTargets.Add(moduleName)
			possibleFailedAmbientModuleLookupSources.LoadOrStore(source.Path(), &failedAmbientModuleLookupSource{
				fileName: source.FileName(),
			})
		})
	}

	indexStart := time.Now()
	var wg sync.WaitGroup
	for packageName := range packageNames.Keys() {
		wg.Go(func() {
			if ctx.Err() != nil {
				return
			}

			typesPackageName := module.GetTypesPackageName(packageName)
			var packageJson *packagejson.InfoCacheEntry
			packageJson = b.host.GetPackageJson(tspath.CombinePaths(dirName, "node_modules", packageName, "package.json"))
			if !packageJson.DirectoryExists {
				packageJson = b.host.GetPackageJson(tspath.CombinePaths(dirName, "node_modules", typesPackageName, "package.json"))
			}
			packageEntrypoints := b.resolver.GetEntrypointsFromPackageJsonInfo(packageJson, packageName)
			if packageEntrypoints == nil {
				return
			}
			if len(fileExcludePatterns) > 0 {
				count := int32(len(packageEntrypoints.Entrypoints))
				packageEntrypoints.Entrypoints = slices.DeleteFunc(packageEntrypoints.Entrypoints, func(entrypoint *module.ResolvedEntrypoint) bool {
					for _, excludePattern := range fileExcludePatterns {
						if matched, _ := excludePattern.MatchString(entrypoint.ResolvedFileName); matched {
							return true
						}
					}
					return false
				})
				atomic.AddInt32(&skippedEntrypointsCount, count-int32(len(packageEntrypoints.Entrypoints)))
			}
			if len(packageEntrypoints.Entrypoints) == 0 {
				return
			}

			entrypointsMu.Lock()
			entrypoints = append(entrypoints, packageEntrypoints)
			entrypointsMu.Unlock()

			aliasResolver := createAliasResolver(packageName, packageEntrypoints.Entrypoints)
			checker, _ := checker.NewChecker(aliasResolver)
			extractor := b.newExportExtractor(dirPath, packageName, checker)
			seenFiles := collections.NewSetWithSizeHint[tspath.Path](len(packageEntrypoints.Entrypoints))
			for _, entrypoint := range aliasResolver.rootFiles {
				if !seenFiles.AddIfAbsent(entrypoint.Path()) {
					continue
				}

				if ctx.Err() != nil {
					return
				}

				fileExports := extractor.extractFromFile(entrypoint)
				exportsMu.Lock()
				for _, name := range entrypoint.AmbientModuleNames {
					ambientModuleNames[name] = append(ambientModuleNames[name], entrypoint.FileName())
				}
				if source, ok := possibleFailedAmbientModuleLookupSources.Load(entrypoint.Path()); !ok {
					// If we failed to resolve any ambient modules from this file, we'll try the
					// whole file again later, so don't add anything now.
					exports[entrypoint.Path()] = fileExports
				} else {
					// Record the package name so we can use it later during the second pass
					source.mu.Lock()
					source.packageName = packageName
					source.mu.Unlock()
				}
				exportsMu.Unlock()
			}
			if logger != nil {
				stats := extractor.Stats()
				combinedStats.exports.Add(stats.exports.Load())
				combinedStats.usedChecker.Add(stats.usedChecker.Load())
			}
		})
	}

	wg.Wait()

	result := &bucketBuildResult{
		bucket: &RegistryBucket{
			Index:              &Index[*Export]{},
			DependencyNames:    dependencies,
			PackageNames:       directoryPackageNames,
			AmbientModuleNames: ambientModuleNames,
			Paths:              *collections.NewSetWithSizeHint[tspath.Path](len(exports)),
			Entrypoints:        make(map[tspath.Path][]*module.ResolvedEntrypoint, len(exports)),
			state: BucketState{
				fileExcludePatterns: b.userPreferences.AutoImportFileExcludePatterns,
			},
		},
		possibleFailedAmbientModuleLookupSources: &possibleFailedAmbientModuleLookupSources,
		possibleFailedAmbientModuleLookupTargets: &possibleFailedAmbientModuleLookupTargets,
	}
	for path, fileExports := range exports {
		result.bucket.Paths.Add(path)
		for _, exp := range fileExports {
			result.bucket.Index.insertAsWords(exp)
		}
	}
	for _, entrypointSet := range entrypoints {
		for _, entrypoint := range entrypointSet.Entrypoints {
			path := b.base.toPath(entrypoint.ResolvedFileName)
			result.bucket.Entrypoints[path] = append(result.bucket.Entrypoints[path], entrypoint)
		}
	}

	if logger != nil {
		logger.Logf("Determined dependencies and package names: %v", extractorStart.Sub(start))
		logger.Logf("Extracted exports: %v (%d exports, %d used checker)", indexStart.Sub(extractorStart), combinedStats.exports.Load(), combinedStats.usedChecker.Load())
		if skippedEntrypointsCount > 0 {
			logger.Logf("Skipped %d entrypoints due to exclude patterns", skippedEntrypointsCount)
		}
		logger.Logf("Built index: %v", time.Since(indexStart))
		logger.Logf("Bucket total: %v", time.Since(start))
	}

	return result, ctx.Err()
}

func (b *registryBuilder) getNearestAncestorDirectoryWithValidPackageJson(filePath tspath.Path) *directory {
	return core.FirstResult(tspath.ForEachAncestorDirectoryPath(filePath.GetDirectoryPath(), func(dirPath tspath.Path) (result *directory, stop bool) {
		if dirEntry, ok := b.directories.Get(dirPath); ok && dirEntry.Value().packageJson.Exists() && dirEntry.Value().packageJson.Contents.Parseable {
			return dirEntry.Value(), true
		}
		return nil, false
	}))
}

func (b *registryBuilder) resolveAmbientModuleName(moduleName string, fromPath tspath.Path) []string {
	return core.FirstResult(tspath.ForEachAncestorDirectoryPath(fromPath, func(dirPath tspath.Path) (result []string, stop bool) {
		if bucket, ok := b.nodeModules.Get(dirPath); ok {
			if fileNames, ok := bucket.Value().AmbientModuleNames[moduleName]; ok {
				return fileNames, true
			}
		}
		return nil, false
	}))
}

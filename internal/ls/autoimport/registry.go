package autoimport

import (
	"cmp"
	"context"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/project/dirty"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type RegistryBucket struct {
	// !!! determine if dirty is only a package.json change, possible no-op if dependencies match
	dirty              bool
	Paths              map[tspath.Path]struct{}
	LookupLocations    map[tspath.Path]struct{}
	PackageNames       *collections.Set[string]
	AmbientModuleNames map[string][]string
	DependencyNames    *collections.Set[string]
	Entrypoints        map[tspath.Path][]*module.ResolvedEntrypoint
	Index              *Index[*Export]
}

func (b *RegistryBucket) Clone() *RegistryBucket {
	return &RegistryBucket{
		dirty:              b.dirty,
		Paths:              b.Paths,
		LookupLocations:    b.LookupLocations,
		AmbientModuleNames: b.AmbientModuleNames,
		DependencyNames:    b.DependencyNames,
		Entrypoints:        b.Entrypoints,
		Index:              b.Index,
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
	toPath func(fileName string) tspath.Path

	// exports      map[tspath.Path][]*RawExport
	directories map[tspath.Path]*directory

	nodeModules map[tspath.Path]*RegistryBucket
	projects    map[tspath.Path]*RegistryBucket

	// relativeSpecifierCache maps from importing file to target file to specifier
	relativeSpecifierCache map[tspath.Path]map[tspath.Path]string
}

func NewRegistry(toPath func(fileName string) tspath.Path) *Registry {
	return &Registry{
		toPath:      toPath,
		directories: make(map[tspath.Path]*directory),
	}
}

func (r *Registry) IsPreparedForImportingFile(fileName string, projectPath tspath.Path) bool {
	projectBucket, ok := r.projects[projectPath]
	if !ok {
		panic("project bucket missing")
	}
	if projectBucket.dirty {
		return false
	}
	path := r.toPath(fileName).GetDirectoryPath()
	for {
		if dirBucket, ok := r.nodeModules[path]; ok {
			if dirBucket.dirty {
				return false
			}
		}
		parent := path.GetDirectoryPath()
		if parent == path {
			break
		}
		path = parent
	}
	return true
}

func (r *Registry) Clone(ctx context.Context, change RegistryChange, host RegistryCloneHost, logger *logging.LogTree) (*Registry, error) {
	start := time.Now()
	if logger != nil {
		logger = logger.Fork("Building autoimport registry")
	}
	builder := newRegistryBuilder(r, host)
	builder.updateBucketAndDirectoryExistence(change, logger)
	builder.markBucketsDirty(change, logger)
	if change.RequestedFile != "" {
		builder.updateIndexes(ctx, change, logger)
	}
	// !!! deref removed source files
	if logger != nil {
		logger.Logf("Built autoimport registry in %v", time.Since(start))
	}
	return builder.Build(), nil
}

type BucketStats struct {
	Path        tspath.Path
	ExportCount int
	FileCount   int
	Dirty       bool
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
			Path:        path,
			ExportCount: exportCount,
			FileCount:   len(bucket.Paths),
			Dirty:       bucket.dirty,
		})
	}

	for path, bucket := range r.nodeModules {
		exportCount := 0
		if bucket.Index != nil {
			exportCount = len(bucket.Index.entries)
		}
		stats.NodeModulesBuckets = append(stats.NodeModulesBuckets, BucketStats{
			Path:        path,
			ExportCount: exportCount,
			FileCount:   len(bucket.Paths),
			Dirty:       bucket.dirty,
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
}

type RegistryCloneHost interface {
	module.ResolutionHost
	FS() vfs.FS
	GetDefaultProject(path tspath.Path) (tspath.Path, *compiler.Program)
	GetProgramForProject(projectPath tspath.Path) *compiler.Program
	GetPackageJson(fileName string) *packagejson.InfoCacheEntry
	GetSourceFile(fileName string, path tspath.Path) *ast.SourceFile
}

type registryBuilder struct {
	// exports     *dirty.MapBuilder[tspath.Path, []*RawExport, []*RawExport]
	host     RegistryCloneHost
	resolver *module.Resolver
	base     *Registry

	directories            *dirty.Map[tspath.Path, *directory]
	nodeModules            *dirty.Map[tspath.Path, *RegistryBucket]
	projects               *dirty.Map[tspath.Path, *RegistryBucket]
	relativeSpecifierCache *dirty.MapBuilder[tspath.Path, map[tspath.Path]string, map[tspath.Path]string]
}

func newRegistryBuilder(registry *Registry, host RegistryCloneHost) *registryBuilder {
	return &registryBuilder{
		host:     host,
		resolver: module.NewResolver(host, core.EmptyCompilerOptions, "", ""),
		base:     registry,

		directories:            dirty.NewMap(registry.directories),
		nodeModules:            dirty.NewMap(registry.nodeModules),
		projects:               dirty.NewMap(registry.projects),
		relativeSpecifierCache: dirty.NewMapBuilder(registry.relativeSpecifierCache, core.Identity, core.Identity),
	}
}

func (b *registryBuilder) Build() *Registry {
	return &Registry{
		toPath:                 b.base.toPath,
		directories:            core.FirstResult(b.directories.Finalize()),
		nodeModules:            core.FirstResult(b.nodeModules.Finalize()),
		projects:               core.FirstResult(b.projects.Finalize()),
		relativeSpecifierCache: core.FirstResult(b.relativeSpecifierCache.Build()),
	}
}

func (b *registryBuilder) updateBucketAndDirectoryExistence(change RegistryChange, logger *logging.LogTree) {
	start := time.Now()
	neededProjects := make(map[tspath.Path]struct{})
	neededDirectories := make(map[tspath.Path]string)
	for path, fileName := range change.OpenFiles {
		neededProjects[core.FirstResult(b.host.GetDefaultProject(path))] = struct{}{}
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

		if _, ok := b.base.relativeSpecifierCache[path]; !ok {
			b.relativeSpecifierCache.Set(path, make(map[tspath.Path]string))
		}
	}

	for path := range b.base.relativeSpecifierCache {
		if _, ok := change.OpenFiles[path]; !ok {
			b.relativeSpecifierCache.Delete(path)
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
			b.projects.Add(projectPath, &RegistryBucket{dirty: true})
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

	updateDirectory := func(dirPath tspath.Path, dirName string) {
		packageJsonFileName := tspath.CombinePaths(dirName, "package.json")
		packageJson := b.host.GetPackageJson(packageJsonFileName)
		hasNodeModules := b.host.FS().DirectoryExists(tspath.CombinePaths(dirName, "node_modules"))
		if entry, ok := b.directories.Get(dirPath); ok {
			entry.ChangeIf(func(dir *directory) bool {
				return dir.packageJson != packageJson || dir.hasNodeModules != hasNodeModules
			}, func(dir *directory) {
				dir.packageJson = packageJson
				dir.hasNodeModules = hasNodeModules
			})
		} else {
			b.directories.Add(dirPath, &directory{
				name:           dirName,
				packageJson:    packageJson,
				hasNodeModules: hasNodeModules,
			})
		}
		if hasNodeModules {
			if hasNodeModulesEntry, ok := b.nodeModules.Get(dirPath); ok {
				hasNodeModulesEntry.ChangeIf(func(bucket *RegistryBucket) bool {
					return !bucket.dirty
				}, func(bucket *RegistryBucket) {
					bucket.dirty = true
				})
			} else {
				b.nodeModules.Add(dirPath, &RegistryBucket{dirty: true})
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
			updateDirectory(dirPath, dirName)
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
			updateDirectory(dirPath, dirName)
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
	start := time.Now()
	cleanNodeModulesBuckets := make(map[tspath.Path]struct{})
	cleanProjectBuckets := make(map[tspath.Path]struct{})
	b.nodeModules.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
		if entry.Value().dirty {
			cleanNodeModulesBuckets[entry.Key()] = struct{}{}
		}
		return true
	})
	b.projects.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
		if entry.Value().dirty {
			cleanProjectBuckets[entry.Key()] = struct{}{}
		}
		return true
	})

	processURIs := func(uris map[lsproto.DocumentUri]struct{}) {
		if len(cleanNodeModulesBuckets) == 0 && len(cleanProjectBuckets) == 0 {
			return
		}
		for uri := range uris {
			// !!! handle package.json effect on node_modules (updateBucketAndDirectoryExistence already detected package.json change)
			path := b.base.toPath(uri.FileName())
			if len(cleanNodeModulesBuckets) > 0 {
				// For node_modules, mark the bucket dirty if anything changes in the directory
				if nodeModulesIndex := strings.Index(string(path), "/node_modules/"); nodeModulesIndex != -1 {
					dirPath := path[:nodeModulesIndex]
					if _, ok := cleanNodeModulesBuckets[dirPath]; ok {
						b.nodeModules.Change(dirPath, func(bucket *RegistryBucket) { bucket.dirty = true })
						delete(cleanNodeModulesBuckets, dirPath)
					}
				}
			}
			if len(cleanProjectBuckets) > 0 {
				// For projects, mark the bucket dirty if the bucket contains the file directly or as a lookup location
				for projectDirPath := range cleanProjectBuckets {
					entry, _ := b.projects.Get(projectDirPath)
					if _, ok := entry.Value().Paths[path]; ok {
						b.projects.Change(projectDirPath, func(bucket *RegistryBucket) { bucket.dirty = true })
						delete(cleanProjectBuckets, projectDirPath)
					} else if _, ok := entry.Value().LookupLocations[path]; ok {
						b.projects.Change(projectDirPath, func(bucket *RegistryBucket) { bucket.dirty = true })
						delete(cleanProjectBuckets, projectDirPath)
					}
				}
			}
		}
	}

	processURIs(change.Created.Keys())
	processURIs(change.Deleted.Keys())
	processURIs(change.Changed.Keys())

	if logger != nil {
		var dirtyNodeModulesPaths, dirtyProjectPaths []tspath.Path
		b.nodeModules.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
			if entry.Value().dirty {
				dirtyNodeModulesPaths = append(dirtyNodeModulesPaths, entry.Key())
			}
			return true
		})
		b.projects.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
			if entry.Value().dirty {
				dirtyProjectPaths = append(dirtyProjectPaths, entry.Key())
			}
			return true
		})
		for _, path := range dirtyNodeModulesPaths {
			logger.Logf("Dirty node_modules bucket: %s", path)
		}
		for _, path := range dirtyProjectPaths {
			logger.Logf("Dirty project bucket: %s", path)
		}
		logger.Logf("Marked buckets dirty in %v", time.Since(start))
	}
}

func (b *registryBuilder) updateIndexes(ctx context.Context, change RegistryChange, logger *logging.LogTree) {
	type task struct {
		entry  *dirty.MapEntry[tspath.Path, *RegistryBucket]
		result *bucketBuildResult
		err    error
	}

	var tasks []*task
	var projectTasks, nodeModulesTasks int
	var wg sync.WaitGroup
	projectPath, _ := b.host.GetDefaultProject(change.RequestedFile)
	if projectPath == "" {
		return
	}
	if project, ok := b.projects.Get(projectPath); ok {
		if project.Value().dirty {
			task := &task{entry: project}
			tasks = append(tasks, task)
			projectTasks++
			wg.Go(func() {
				index, err := b.buildProjectBucket(ctx, projectPath)
				task.result = index
				task.err = err
			})
		}
	}
	tspath.ForEachAncestorDirectoryPath(change.RequestedFile, func(dirPath tspath.Path) (any, bool) {
		if nodeModulesBucket, ok := b.nodeModules.Get(dirPath); ok {
			if nodeModulesBucket.Value().dirty {
				dirName := core.FirstResult(b.directories.Get(dirPath)).Value().name
				task := &task{entry: nodeModulesBucket}
				tasks = append(tasks, task)
				nodeModulesTasks++
				wg.Go(func() {
					result, err := b.buildNodeModulesBucket(ctx, change, dirName, dirPath)
					task.result = result
					task.err = err
				})
			}
		}
		return nil, false
	})

	if logger != nil && len(tasks) > 0 {
		logger.Logf("Building %d indexes (%d projects, %d node_modules)", len(tasks), projectTasks, nodeModulesTasks)
	}

	start := time.Now()
	wg.Wait()

	// !!! clean up this hot mess
	for _, t := range tasks {
		if t.err != nil {
			continue
		}
		t.entry.Replace(t.result.bucket)
	}

	secondPassStart := time.Now()
	var secondPassFileCount int
	for _, t := range tasks {
		if t.err != nil {
			continue
		}
		if t.result.possibleFailedAmbientModuleLookupTargets == nil {
			continue
		}
		rootFiles := make(map[string]struct{})
		for target := range t.result.possibleFailedAmbientModuleLookupTargets.Keys() {
			for _, fileName := range b.resolveAmbientModuleName(target, t.entry.Key()) {
				rootFiles[fileName] = struct{}{}
				secondPassFileCount++
			}
		}
		if len(rootFiles) > 0 {
			// !!! parallelize?
			aliasResolver := newAliasResolver(slices.Collect(maps.Keys(rootFiles)), b.host, b.resolver, b.base.toPath)
			ch := checker.NewChecker(aliasResolver)
			t.result.possibleFailedAmbientModuleLookupSources.Range(func(path tspath.Path, source *failedAmbientModuleLookupSource) bool {
				sourceFile := aliasResolver.GetSourceFile(source.fileName)
				extractor := b.newExportExtractor(t.entry.Key(), source.packageName, func() (*checker.Checker, func()) {
					return ch, func() {}
				})
				fileExports := extractor.extractFromFile(sourceFile)
				t.result.bucket.Paths[path] = struct{}{}
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

type bucketBuildResult struct {
	bucket *RegistryBucket
	// File path to filename and package name
	possibleFailedAmbientModuleLookupSources *collections.SyncMap[tspath.Path, *failedAmbientModuleLookupSource]
	// Likely ambient module name
	possibleFailedAmbientModuleLookupTargets *collections.SyncSet[string]
}

func (b *registryBuilder) buildProjectBucket(ctx context.Context, projectPath tspath.Path) (*bucketBuildResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var mu sync.Mutex
	result := &bucketBuildResult{bucket: &RegistryBucket{}}
	program := b.host.GetProgramForProject(projectPath)
	exports := make(map[tspath.Path][]*Export)
	var wg sync.WaitGroup
	getChecker, closePool := b.createCheckerPool(program)
	defer closePool()
	extractor := b.newExportExtractor("", "", getChecker)
	for _, file := range program.GetSourceFiles() {
		if strings.Contains(file.FileName(), "/node_modules/") || program.IsSourceFileDefaultLibrary(file.Path()) {
			continue
		}
		wg.Go(func() {
			if ctx.Err() == nil {
				// !!! we could consider doing ambient modules / augmentations more directly
				// from the program checker, instead of doing the syntax-based collection
				fileExports := extractor.extractFromFile(file)
				mu.Lock()
				exports[file.Path()] = fileExports
				mu.Unlock()
			}
		})
	}

	wg.Wait()
	idx := &Index[*Export]{}
	for path, fileExports := range exports {
		if result.bucket.Paths == nil {
			result.bucket.Paths = make(map[tspath.Path]struct{}, len(exports))
		}
		result.bucket.Paths[path] = struct{}{}
		for _, exp := range fileExports {
			idx.insertAsWords(exp)
		}
	}

	result.bucket.Index = idx
	return result, nil
}

func (b *registryBuilder) buildNodeModulesBucket(ctx context.Context, change RegistryChange, dirName string, dirPath tspath.Path) (*bucketBuildResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// If any open files are in scope of this directory but not in scope of any package.json,
	// we need to add all packages in this node_modules directory.
	// !!! ensure a different set of open files properly invalidates
	// buckets that are built but may be incomplete due to different package.json visibility
	// !!! should we really be preparing buckets for all open files? Could dirty tracking
	// be more granular? what are the actual inputs that determine whether a bucket is valid
	// for a given importing file?
	// For now, we'll always build for all open files. This `dependencies` computation
	// should be moved out and the result used to determine whether we need a rebuild.
	var dependencies *collections.Set[string]
	var packageNames *collections.Set[string]
	for path := range change.OpenFiles {
		if dirPath.ContainsPath(path) && b.getNearestAncestorDirectoryWithPackageJson(path) == nil {
			dependencies = nil
			break
		}
		dependencies = &collections.Set[string]{}
	}
	directoryPackageNames, err := getPackageNamesInNodeModules(tspath.CombinePaths(dirName, "node_modules"), b.host.FS())
	if err != nil {
		return nil, err
	}

	// Get all package.jsons that have this node_modules directory in their spine
	if dependencies != nil {
		b.directories.Range(func(entry *dirty.MapEntry[tspath.Path, *directory]) bool {
			if entry.Value().packageJson.Exists() && dirPath.ContainsPath(entry.Key()) {
				entry.Value().packageJson.Contents.RangeDependencies(func(name, _, _ string) bool {
					dependencies.Add(module.GetPackageNameFromTypesPackageName(name))
					return true
				})
			}
			return true
		})
		packageNames = dependencies
	} else {
		packageNames = directoryPackageNames
	}

	aliasResolver := newAliasResolver(nil, b.host, b.resolver, b.base.toPath)
	getChecker, closePool := b.createCheckerPool(aliasResolver)
	defer closePool()

	var exportsMu sync.Mutex
	exports := make(map[tspath.Path][]*Export)
	ambientModuleNames := make(map[string][]string)

	var entrypointsMu sync.Mutex
	var entrypoints []*module.ResolvedEntrypoints

	processFile := func(fileName string, path tspath.Path, packageName string) {
		sourceFile := b.host.GetSourceFile(fileName, path)
		binder.BindSourceFile(sourceFile)
		extractor := b.newExportExtractor(dirPath, packageName, getChecker)
		fileExports := extractor.extractFromFile(sourceFile)
		if source, ok := aliasResolver.possibleFailedAmbientModuleLookupSources.Load(sourceFile.Path()); !ok {
			// If we failed to resolve any ambient modules from this file, we'll try the
			// whole file again later, so don't add anything now.
			exportsMu.Lock()
			exports[path] = fileExports
			for _, name := range sourceFile.AmbientModuleNames {
				ambientModuleNames[name] = append(ambientModuleNames[name], fileName)
			}
			exportsMu.Unlock()
		} else {
			// Record the package name so we can use it later during the second pass
			// !!! perhaps we could store the whole set of partial exports and avoid
			//     repeating some work
			source.mu.Lock()
			source.packageName = packageName
			source.mu.Unlock()
		}
	}

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
			entrypointsMu.Lock()
			entrypoints = append(entrypoints, packageEntrypoints)
			entrypointsMu.Unlock()

			seenFiles := collections.NewSetWithSizeHint[tspath.Path](len(packageEntrypoints.Entrypoints))
			for _, entrypoint := range packageEntrypoints.Entrypoints {
				path := b.base.toPath(entrypoint.ResolvedFileName)
				if !seenFiles.AddIfAbsent(path) {
					continue
				}

				wg.Go(func() {
					if ctx.Err() != nil {
						return
					}
					processFile(entrypoint.ResolvedFileName, path, packageName)
				})
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
			Paths:              make(map[tspath.Path]struct{}, len(exports)),
			Entrypoints:        make(map[tspath.Path][]*module.ResolvedEntrypoint, len(exports)),
			LookupLocations:    make(map[tspath.Path]struct{}),
		},
		possibleFailedAmbientModuleLookupSources: &aliasResolver.possibleFailedAmbientModuleLookupSources,
		possibleFailedAmbientModuleLookupTargets: &aliasResolver.possibleFailedAmbientModuleLookupTargets,
	}
	for path, fileExports := range exports {
		result.bucket.Paths[path] = struct{}{}
		for _, exp := range fileExports {
			result.bucket.Index.insertAsWords(exp)
		}
	}
	for _, entrypointSet := range entrypoints {
		for _, entrypoint := range entrypointSet.Entrypoints {
			path := b.base.toPath(entrypoint.ResolvedFileName)
			result.bucket.Entrypoints[path] = append(result.bucket.Entrypoints[path], entrypoint)
		}
		for _, failedLocation := range entrypointSet.FailedLookupLocations {
			result.bucket.LookupLocations[b.base.toPath(failedLocation)] = struct{}{}
		}
	}

	return result, ctx.Err()
}

// !!! tune default size, create on demand
const checkerPoolSize = 16

func (b *registryBuilder) createCheckerPool(program checker.Program) (getChecker func() (*checker.Checker, func()), closePool func()) {
	pool := make(chan *checker.Checker, checkerPoolSize)
	for range checkerPoolSize {
		pool <- checker.NewChecker(program)
	}
	return func() (*checker.Checker, func()) {
			checker := <-pool
			return checker, func() {
				pool <- checker
			}
		}, func() {
			close(pool)
		}
}

func (b *registryBuilder) getNearestAncestorDirectoryWithPackageJson(filePath tspath.Path) *directory {
	return core.FirstResult(tspath.ForEachAncestorDirectoryPath(filePath.GetDirectoryPath(), func(dirPath tspath.Path) (result *directory, stop bool) {
		if dirEntry, ok := b.directories.Get(dirPath); ok && dirEntry.Value().packageJson.Exists() {
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

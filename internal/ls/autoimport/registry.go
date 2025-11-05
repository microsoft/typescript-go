package autoimport

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
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
	dirty           bool
	Paths           map[tspath.Path]struct{}
	LookupLocations map[tspath.Path]struct{}
	Dependencies    collections.Set[string]
	Entrypoints     map[tspath.Path][]*module.ResolvedEntrypoint
	Index           *Index[*RawExport]
}

func (b *RegistryBucket) Clone() *RegistryBucket {
	return &RegistryBucket{
		dirty:           b.dirty,
		Paths:           b.Paths,
		LookupLocations: b.LookupLocations,
		Dependencies:    b.Dependencies,
		Index:           b.Index,
	}
}

type directory struct {
	packageJson    *packagejson.InfoCacheEntry
	hasNodeModules bool
}

func (d *directory) Clone() *directory {
	return &directory{
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

	directories *dirty.Map[tspath.Path, *directory]
	nodeModules *dirty.Map[tspath.Path, *RegistryBucket]
	projects    *dirty.Map[tspath.Path, *RegistryBucket]
}

func newRegistryBuilder(registry *Registry, host RegistryCloneHost) *registryBuilder {
	return &registryBuilder{
		// exports:     dirty.NewMapBuilder(registry.exports, slices.Clone, core.Identity),
		host:     host,
		resolver: module.NewResolver(host, core.EmptyCompilerOptions, "", ""),
		base:     registry,

		directories: dirty.NewMap(registry.directories),
		nodeModules: dirty.NewMap(registry.nodeModules),
		projects:    dirty.NewMap(registry.projects),
	}
}

func (b *registryBuilder) Build() *Registry {
	return &Registry{
		toPath:      b.base.toPath,
		directories: core.FirstResult(b.directories.Finalize()),
		nodeModules: core.FirstResult(b.nodeModules.Finalize()),
		projects:    core.FirstResult(b.projects.Finalize()),
	}
}

func (b *registryBuilder) updateBucketAndDirectoryExistence(change RegistryChange, logger *logging.LogTree) {
	start := time.Now()
	neededProjects := make(map[tspath.Path]struct{})
	neededDirectories := make(map[tspath.Path]string)
	for path, fileName := range change.OpenFiles {
		neededProjects[core.FirstResult(b.host.GetDefaultProject(path))] = struct{}{}
		dir := fileName
		for {
			dir = tspath.GetDirectoryPath(dir)
			dirPath := path.GetDirectoryPath()
			if path == dirPath {
				break
			}
			if _, ok := neededDirectories[dirPath]; ok {
				break
			}
			neededDirectories[dirPath] = dir
			path = dirPath
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
		result *RegistryBucket
		err    error
	}

	var tasks []*task
	var projectTasks, nodeModulesTasks int
	wg := core.NewWorkGroup(false)
	projectPath, _ := b.host.GetDefaultProject(change.RequestedFile)
	if projectPath == "" {
		return
	}
	if project, ok := b.projects.Get(projectPath); ok {
		if project.Value().dirty {
			task := &task{entry: project}
			tasks = append(tasks, task)
			projectTasks++
			wg.Queue(func() {
				index, err := b.buildProjectIndex(ctx, projectPath)
				task.result = index
				task.err = err
			})
		}
	}
	tspath.ForEachAncestorDirectoryPath(change.RequestedFile, func(dirPath tspath.Path) (any, bool) {
		if nodeModulesBucket, ok := b.nodeModules.Get(dirPath); ok {
			if nodeModulesBucket.Value().dirty {
				task := &task{entry: nodeModulesBucket}
				tasks = append(tasks, task)
				nodeModulesTasks++
				wg.Queue(func() {
					index, err := b.buildNodeModulesIndex(ctx, dirPath)
					task.result = index
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
	wg.RunAndWait()
	if logger != nil && len(tasks) > 0 {
		logger.Logf("Built %d indexes in %v", len(tasks), time.Since(start))
	}
	for _, t := range tasks {
		if t.err != nil {
			continue
		}
		t.entry.Change(func(bucket *RegistryBucket) {
			bucket.dirty = false
			bucket.Index = t.result.Index
			bucket.Paths = t.result.Paths
			bucket.LookupLocations = t.result.LookupLocations
		})
	}
}

func (b *registryBuilder) buildProjectIndex(ctx context.Context, projectPath tspath.Path) (*RegistryBucket, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var mu sync.Mutex
	result := &RegistryBucket{}
	program := b.host.GetProgramForProject(projectPath)
	exports := make(map[tspath.Path][]*RawExport)
	wg := core.NewWorkGroup(false)
	for _, file := range program.GetSourceFiles() {
		if strings.Contains(file.FileName(), "/node_modules/") {
			continue
		}
		wg.Queue(func() {
			if ctx.Err() == nil {
				fileExports := Parse(file)
				mu.Lock()
				exports[file.Path()] = fileExports
				mu.Unlock()
			}
		})
	}

	wg.RunAndWait()
	idx := NewIndexBuilder[*RawExport](nil)
	for path, fileExports := range exports {
		if result.Paths == nil {
			result.Paths = make(map[tspath.Path]struct{}, len(exports))
		}
		result.Paths[path] = struct{}{}
		for _, exp := range fileExports {
			idx.InsertAsWords(exp)
		}
	}

	result.Index = idx.Index()
	return result, nil
}

func (b *registryBuilder) buildNodeModulesIndex(ctx context.Context, dirPath tspath.Path) (*RegistryBucket, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// get all package.jsons that have this node_modules directory in their spine
	// !!! distinguish between no dependencies and no package.jsons
	var dependencies collections.Set[string]
	b.directories.Range(func(entry *dirty.MapEntry[tspath.Path, *directory]) bool {
		if entry.Value().packageJson.Exists() && dirPath.ContainsPath(entry.Key()) {
			entry.Value().packageJson.Contents.RangeDependencies(func(name, _, _ string) bool {
				dependencies.Add(name)
				return true
			})
		}
		return true
	})

	var exportsMu sync.Mutex
	exports := make(map[tspath.Path][]*RawExport)
	var entrypointsMu sync.Mutex
	var entrypoints []*module.ResolvedEntrypoints
	wg := core.NewWorkGroup(false)

	for dep := range dependencies.Keys() {
		wg.Queue(func() {
			if ctx.Err() != nil {
				return
			}
			packageJson := b.host.GetPackageJson(tspath.CombinePaths(string(dirPath), "node_modules", dep, "package.json"))
			if !packageJson.Exists() {
				return
			}
			packageEntrypoints := b.resolver.GetEntrypointsFromPackageJsonInfo(packageJson)
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

				wg.Queue(func() {
					if ctx.Err() != nil {
						return
					}
					sourceFile := b.host.GetSourceFile(entrypoint.ResolvedFileName, path)
					binder.BindSourceFile(sourceFile)
					fileExports := Parse(sourceFile)
					exportsMu.Lock()
					exports[path] = fileExports
					exportsMu.Unlock()
				})
			}
		})
	}

	wg.RunAndWait()
	result := &RegistryBucket{
		Dependencies:    dependencies,
		Paths:           make(map[tspath.Path]struct{}, len(exports)),
		Entrypoints:     make(map[tspath.Path][]*module.ResolvedEntrypoint, len(exports)),
		LookupLocations: make(map[tspath.Path]struct{}),
	}
	idx := NewIndexBuilder[*RawExport](nil)
	for path, fileExports := range exports {
		result.Paths[path] = struct{}{}
		for _, exp := range fileExports {
			idx.InsertAsWords(exp)
		}
	}
	result.Index = idx.Index()
	for _, entrypointSet := range entrypoints {
		for _, entrypoint := range entrypointSet.Entrypoints {
			path := b.base.toPath(entrypoint.ResolvedFileName)
			result.Entrypoints[path] = append(result.Entrypoints[path], entrypoint)
		}
		for _, failedLocation := range entrypointSet.FailedLookupLocations {
			result.LookupLocations[b.base.toPath(failedLocation)] = struct{}{}
		}
	}

	return result, nil
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

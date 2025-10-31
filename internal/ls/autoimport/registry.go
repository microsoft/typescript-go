package autoimport

import (
	"context"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/project/dirty"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type RegistryBucket struct {
	// !!! determine if dirty is only a package.json change, possible no-op if dependencies match
	dirty           bool
	Paths           map[tspath.Path]struct{}
	LookupLocations map[tspath.Path]struct{}
	Dependencies    collections.Set[string]
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
	packageJson    *packagejson.DependencyFields
	hasNodeModules bool
}

func (d *directory) Clone() *directory {
	return &directory{
		packageJson:    d.packageJson,
		hasNodeModules: d.hasNodeModules,
	}
}

type Registry struct {
	toPath    func(fileName string) tspath.Path
	openFiles map[tspath.Path]string

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

type RegistryChange struct {
	RequestedFile tspath.Path
	OpenFiles     map[tspath.Path]string
	Changed       collections.Set[lsproto.DocumentUri]
	Created       collections.Set[lsproto.DocumentUri]
	Deleted       collections.Set[lsproto.DocumentUri]
}

type RegistryCloneHost interface {
	FS() vfs.FS
	GetDefaultProject(fileName string) (tspath.Path, *compiler.Program)
	GetProgramForProject(projectPath tspath.Path) *compiler.Program
	GetPackageJson(fileName string) *packagejson.DependencyFields
}

type registryBuilder struct {
	// exports     *dirty.MapBuilder[tspath.Path, []*RawExport, []*RawExport]
	host RegistryCloneHost
	base *Registry

	directories *dirty.Map[tspath.Path, *directory]
	nodeModules *dirty.Map[tspath.Path, *RegistryBucket]
	projects    *dirty.Map[tspath.Path, *RegistryBucket]
}

func newRegistryBuilder(registry *Registry, host RegistryCloneHost) *registryBuilder {
	return &registryBuilder{
		// exports:     dirty.NewMapBuilder(registry.exports, slices.Clone, core.Identity),
		host: host,
		base: registry,

		directories: dirty.NewMap(registry.directories),
		nodeModules: dirty.NewMap(registry.nodeModules),
		projects:    dirty.NewMap(registry.projects),
	}
}

func (b *registryBuilder) Build() *Registry {
	return &Registry{
		// exports:     b.exports.Build(),

	}
}

func (b *registryBuilder) updateBucketAndDirectoryExistence(change RegistryChange) {
	neededProjects := make(map[tspath.Path]struct{})
	neededDirectories := make(map[tspath.Path]string)
	for path, fileName := range change.OpenFiles {
		neededProjects[core.FirstResult(b.host.GetDefaultProject(fileName))] = struct{}{}
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

	core.DiffMapsFunc(
		b.base.projects,
		neededProjects,
		func(_ *RegistryBucket, _ struct{}) bool {
			panic("never called because onChanged is nil")
		},
		func(projectPath tspath.Path, _ struct{}) {
			// Need and don't have
			b.projects.Add(projectPath, &RegistryBucket{dirty: true})
		},
		func(projectPath tspath.Path, _ *RegistryBucket) {
			// Have and don't need
			b.projects.Delete(projectPath)
		},
		nil,
	)

	updateDirectory := func(dirPath tspath.Path, dirName string) {
		packageJsonFileName := tspath.CombinePaths(dirName, "package.json")
		packageJson := b.host.GetPackageJson(packageJsonFileName)
		hasNodeModules := b.host.FS().DirectoryExists(tspath.CombinePaths(dirName, "node_modules"))
		b.directories.Add(dirPath, &directory{
			packageJson:    packageJson,
			hasNodeModules: hasNodeModules,
		})
		if hasNodeModules {
			b.nodeModules.Add(dirPath, &RegistryBucket{dirty: true})
		} else {
			b.nodeModules.Delete(dirPath)
		}
	}

	core.DiffMapsFunc(
		b.base.directories,
		neededDirectories,
		func(dir *directory, dirName string) bool {
			packageJsonUri := lsconv.FileNameToDocumentURI(tspath.CombinePaths(dirName, "package.json"))
			return change.Changed.Has(packageJsonUri) || change.Deleted.Has(packageJsonUri) || change.Created.Has(packageJsonUri)
		},
		func(dirPath tspath.Path, dirName string) {
			// Need and don't have
			updateDirectory(dirPath, dirName)
		},
		func(dirPath tspath.Path, dir *directory) {
			// Have and don't need
			b.directories.Delete(dirPath)
			b.nodeModules.Delete(dirPath)
		},
		func(dirPath tspath.Path, dir *directory, dirName string) {
			// package.json may have changed
			updateDirectory(dirPath, dirName)
		},
	)
}

func (b *registryBuilder) markBucketsDirty(change RegistryChange) {
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
}

func (b *registryBuilder) updateIndexes(ctx context.Context, change RegistryChange) {
	type task struct {
		entry  *dirty.MapEntry[tspath.Path, *RegistryBucket]
		result *RegistryBucket
		err    error
	}

	var tasks []*task
	wg := core.NewWorkGroup(false)
	b.projects.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
		if entry.Value().dirty {
			task := &task{entry: entry}
			tasks = append(tasks, task)
			wg.Queue(func() {
				index, err := b.buildProjectIndex(ctx, entry.Key())
				task.result = index
				task.err = err
			})
		}
		return true
	})
	b.nodeModules.Range(func(entry *dirty.MapEntry[tspath.Path, *RegistryBucket]) bool {
		if entry.Value().dirty {
			task := &task{entry: entry}
			tasks = append(tasks, task)
			wg.Queue(func() {
				index, err := b.buildNodeModulesIndex(ctx, entry.Key())
				task.result = index
				task.err = err
			})
		}
		return true
	})

	wg.RunAndWait()
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
	var dependencies collections.Set[string]
	b.directories.Range(func(entry *dirty.MapEntry[tspath.Path, *directory]) bool {
		if entry.Value().packageJson != nil && dirPath.ContainsPath(entry.Key().GetDirectoryPath()) {
			entry.Value().packageJson.RangeDependencies(func(name, _, _ string) bool {
				dependencies.Add(name)
				return true
			})
		}
		return true
	})

	for dep := range dependencies.Keys() {
		packageJson := b.host.GetPackageJson(tspath.CombinePaths(string(dirPath), "node_modules", dep, "package.json"))
		if packageJson == nil {
			continue
		}

	}
}

func (r *Registry) Clone(ctx context.Context, change RegistryChange, host RegistryCloneHost) (*Registry, error) {
	builder := newRegistryBuilder(r, host)
	builder.updateBucketAndDirectoryExistence(change)
	builder.markBucketsDirty(change)
	builder.updateIndexes(ctx, change)
	return builder.Build(), nil
}

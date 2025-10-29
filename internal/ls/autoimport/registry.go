package autoimport

import (
	"context"
	"slices"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/project/dirty"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Registry struct {
	exports     map[tspath.Path][]*RawExport
	nodeModules map[tspath.Path]*Index[*RawExport]
	projects    map[tspath.Path]*Index[*RawExport]
}

type Project struct {
	Key     tspath.Path
	Program *compiler.Program
}

type RegistryChange struct {
	WithProject *Project
}

type registryBuilder struct {
	exports     *dirty.MapBuilder[tspath.Path, []*RawExport, []*RawExport]
	nodeModules *dirty.MapBuilder[tspath.Path, *Index[*RawExport], *IndexBuilder[*RawExport]]
	projects    *dirty.MapBuilder[tspath.Path, *Index[*RawExport], *IndexBuilder[*RawExport]]
}

func newRegistryBuilder(registry *Registry) *registryBuilder {
	if registry == nil {
		registry = &Registry{}
	}
	return &registryBuilder{
		exports:     dirty.NewMapBuilder(registry.exports, slices.Clone, core.Identity),
		nodeModules: dirty.NewMapBuilder(registry.nodeModules, NewIndexBuilder, (*IndexBuilder[*RawExport]).Index),
		projects:    dirty.NewMapBuilder(registry.projects, NewIndexBuilder, (*IndexBuilder[*RawExport]).Index),
	}
}

func (b *registryBuilder) Build() *Registry {
	return &Registry{
		exports:     b.exports.Build(),
		nodeModules: b.nodeModules.Build(),
		projects:    b.projects.Build(),
	}
}

// With what granularity will we perform updates? How do we remove stale entries?
// Will we always rebuild full tries, or update them? If rebuild, do we need TrieBuilder?

func (r *Registry) Clone(ctx context.Context, change RegistryChange) (*Registry, error) {
	builder := newRegistryBuilder(r)
	if change.WithProject != nil {
		var mu sync.Mutex
		exports := make(map[tspath.Path][]*RawExport)
		wg := core.NewWorkGroup(false)
		for _, file := range change.WithProject.Program.GetSourceFiles() {
			if strings.Contains(file.FileName(), "/node_modules/") {
				continue
			}
			wg.Queue(func() {
				if ctx.Err() == nil {
					// !!! check file hash
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
			builder.exports.Set(path, fileExports)
			for _, exp := range fileExports {
				idx.InsertAsWords(exp)
			}
		}
		builder.projects.Set(change.WithProject.Key, idx)
	}
	return builder.Build(), nil
}

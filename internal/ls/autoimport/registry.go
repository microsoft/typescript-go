package autoimport

import (
	"context"
	"slices"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/project/dirty"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Registry struct {
	exports     map[tspath.Path][]*RawExport
	nodeModules map[tspath.Path]*Trie[RawExport]
	projects    map[tspath.Path]*Trie[RawExport]
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
	nodeModules *dirty.MapBuilder[tspath.Path, *Trie[RawExport], *TrieBuilder[RawExport]]
	projects    *dirty.MapBuilder[tspath.Path, *Trie[RawExport], *TrieBuilder[RawExport]]
}

func newRegistryBuilder(corpus *Registry) *registryBuilder {
	return &registryBuilder{
		exports:     dirty.NewMapBuilder(corpus.exports, slices.Clone, core.Identity),
		nodeModules: dirty.NewMapBuilder(corpus.nodeModules, NewTrieBuilder, (*TrieBuilder[RawExport]).Trie),
		projects:    dirty.NewMapBuilder(corpus.projects, NewTrieBuilder, (*TrieBuilder[RawExport]).Trie),
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
		trie := NewTrieBuilder(nil)
		for path, fileExports := range exports {
			builder.exports.Set(path, fileExports)
			for _, exp := range fileExports {
				trie.InsertAsWords(exp.Name, exp)
			}
		}
	}
	return builder.Build(), nil
}

// Idea: one trie per package.json / node_modules directory?
// - *RawExport lives in shared structure by realpath
//   - could literally just live on SourceFile...
// - solves package shadowing, unreachable node_modules, dependency filtering
// - these tries should be shareable across different projects
// - non-node_modules files form tries per project?

func Collect(ctx context.Context, files []*ast.SourceFile) (*Registry, error) {
	var exports []*RawExport
	wg := core.NewWorkGroup(false)
	for _, file := range files {
		wg.Queue(func() {
			if ctx.Err() == nil {
				exports = append(exports, Parse(file)...)
			}
		})
	}
	wg.RunAndWait()

	var trie *Trie[RawExport]
	if ctx.Err() == nil {
		trie = &Trie[RawExport]{}
		for _, exp := range exports {
			trie.InsertAsWords(exp.Name, exp)
		}
	}
	return &Registry{Exports: exports, Trie: trie}, ctx.Err()
}

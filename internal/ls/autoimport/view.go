package autoimport

import (
	"context"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type View struct {
	registry      *Registry
	importingFile *ast.SourceFile
	program       *compiler.Program
	preferences   modulespecifiers.UserPreferences
	projectKey    tspath.Path

	existingImports *collections.MultiMap[ModuleID, existingImport]
}

func NewView(registry *Registry, importingFile *ast.SourceFile, projectKey tspath.Path, program *compiler.Program, preferences modulespecifiers.UserPreferences) *View {
	return &View{
		registry:      registry,
		importingFile: importingFile,
		program:       program,
		projectKey:    projectKey,
		preferences:   preferences,
	}
}

func (v *View) Search(prefix string) []*RawExport {
	// !!! deal with duplicates due to symlinks
	var results []*RawExport
	bucket, ok := v.registry.projects[v.projectKey]
	if ok {
		results = append(results, bucket.Index.Search(prefix, nil)...)
	}

	var excludePackages collections.Set[string]
	tspath.ForEachAncestorDirectoryPath(v.importingFile.Path().GetDirectoryPath(), func(dirPath tspath.Path) (result any, stop bool) {
		if nodeModulesBucket, ok := v.registry.nodeModules[dirPath]; ok {
			var filter func(e *RawExport) bool
			if excludePackages.Len() > 0 {
				filter = func(e *RawExport) bool {
					return !excludePackages.Has(e.PackageName)
				}
			}

			results = append(results, nodeModulesBucket.Index.Search(prefix, filter)...)
			excludePackages = *excludePackages.Union(nodeModulesBucket.PackageNames)
		}
		return nil, false
	})
	return results
}

type FixAndExport struct {
	Fix    *Fix
	Export *RawExport
}

func (v *View) GetCompletions(ctx context.Context, prefix string) []*FixAndExport {
	results := v.Search(prefix)

	type exportGroupKey struct {
		target            ExportID
		name              string
		ambientModuleName string
	}
	grouped := make(map[exportGroupKey][]*RawExport, len(results))
	for _, e := range results {
		if string(e.ModuleID) == string(v.importingFile.Path()) {
			// Don't auto-import from the importing file itself
			continue
		}
		target := e.ExportID
		if e.Target != (ExportID{}) {
			target = e.Target
		}
		key := exportGroupKey{
			target:            target,
			name:              e.Name(),
			ambientModuleName: e.AmbientModuleName(),
		}
		if existing, ok := grouped[key]; ok {
			for i, ex := range existing {
				if e.ExportID == ex.ExportID {
					grouped[key] = slices.Replace(existing, i, i+1, &RawExport{
						ExportID:                   e.ExportID,
						Syntax:                     e.Syntax,
						Flags:                      e.Flags | ex.Flags,
						ScriptElementKind:          min(e.ScriptElementKind, ex.ScriptElementKind),
						ScriptElementKindModifiers: *e.ScriptElementKindModifiers.Union(&ex.ScriptElementKindModifiers),
						localName:                  e.localName,
						Target:                     e.Target,
						FileName:                   e.FileName,
						Path:                       e.Path,
						NodeModulesDirectory:       e.NodeModulesDirectory,
					})
				}
			}
		}
		grouped[key] = append(grouped[key], e)
	}

	fixes := make([]*FixAndExport, 0, len(results))
	compareFixes := func(a, b *FixAndExport) int {
		return v.compareFixes(a.Fix, b.Fix)
	}

	for _, exps := range grouped {
		fixesForGroup := make([]*FixAndExport, 0, len(exps))
		for _, e := range exps {
			for _, fix := range v.GetFixes(ctx, e) {
				fixesForGroup = append(fixesForGroup, &FixAndExport{
					Fix:    fix,
					Export: e,
				})
			}
		}
		fixes = append(fixes, slices.MinFunc(fixesForGroup, compareFixes))
	}

	return fixes
}

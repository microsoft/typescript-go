package autoimport

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type View struct {
	registry      *Registry
	importingFile *ast.SourceFile
	program       *compiler.Program
	projectKey    tspath.Path

	existingImports *collections.MultiMap[ModuleID, existingImport]
}

func NewView(registry *Registry, importingFile *ast.SourceFile, projectKey tspath.Path, program *compiler.Program) *View {
	return &View{
		registry:      registry,
		importingFile: importingFile,
		program:       program,
		projectKey:    projectKey,
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

	groupedByTarget := make(map[ExportID][]*RawExport, len(results))
	for _, e := range results {
		if string(e.ModuleID) == string(v.importingFile.Path()) {
			// Don't auto-import from the importing file itself
			continue
		}
		key := e.ExportID
		if e.Target != (ExportID{}) {
			key = e.Target
		}
		if existing, ok := groupedByTarget[key]; ok {
			for i, ex := range existing {
				if e.ExportID == ex.ExportID {
					groupedByTarget[key] = slices.Replace(existing, i, i+1, &RawExport{
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
		groupedByTarget[key] = append(groupedByTarget[key], e)
	}

	mergedResults := make([]*RawExport, 0, len(results))
	for _, exps := range groupedByTarget {
		// !!! some kind of sort
		if len(exps) > 1 {
			var seenAmbientSpecifiers collections.Set[string]
			var seenNames collections.Set[string]
			for _, exp := range exps {
				if !tspath.IsExternalModuleNameRelative(string(exp.ModuleID)) && seenAmbientSpecifiers.AddIfAbsent(string(exp.ModuleID)) {
					mergedResults = append(mergedResults, exp)
				} else if !seenNames.Has(exp.Name()) {
					seenNames.Add(exp.Name())
					mergedResults = append(mergedResults, exp)
				}
			}
		} else {
			mergedResults = append(mergedResults, exps[0])
		}
	}

	return mergedResults
}

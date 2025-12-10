package autoimport

import (
	"context"
	"slices"
	"unicode"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type View struct {
	registry       *Registry
	importingFile  *ast.SourceFile
	program        *compiler.Program
	preferences    modulespecifiers.UserPreferences
	projectKey     tspath.Path
	allowedEndings []modulespecifiers.ModuleSpecifierEnding

	existingImports          *collections.MultiMap[ModuleID, existingImport]
	shouldUseRequireForFixes *bool
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

func (v *View) getAllowedEndings() []modulespecifiers.ModuleSpecifierEnding {
	if v.allowedEndings == nil {
		resolutionMode := v.program.GetDefaultResolutionModeForFile(v.importingFile)
		v.allowedEndings = modulespecifiers.GetAllowedEndingsInPreferredOrder(
			v.preferences,
			v.program,
			v.program.Options(),
			v.importingFile,
			"",
			resolutionMode,
		)
	}
	return v.allowedEndings
}

type QueryKind int

const (
	QueryKindWordPrefix QueryKind = iota
	QueryKindExactMatch
	QueryKindCaseInsensitiveMatch
)

func (v *View) Search(query string, kind QueryKind) []*Export {
	var results []*Export
	search := func(bucket *RegistryBucket) []*Export {
		switch kind {
		case QueryKindWordPrefix:
			return bucket.Index.SearchWordPrefix(query)
		case QueryKindExactMatch:
			return bucket.Index.Find(query, true)
		case QueryKindCaseInsensitiveMatch:
			return bucket.Index.Find(query, false)
		default:
			panic("unreachable")
		}
	}

	if bucket, ok := v.registry.projects[v.projectKey]; ok {
		exports := search(bucket)
		results = slices.Grow(results, len(exports))
		for _, e := range exports {
			if string(e.ModuleID) == string(v.importingFile.Path()) {
				// Don't auto-import from the importing file itself
				continue
			}
			results = append(results, e)
		}
	}

	var excludePackages *collections.Set[string]
	tspath.ForEachAncestorDirectoryPath(v.importingFile.Path().GetDirectoryPath(), func(dirPath tspath.Path) (result any, stop bool) {
		if nodeModulesBucket, ok := v.registry.nodeModules[dirPath]; ok {
			exports := search(nodeModulesBucket)
			if excludePackages.Len() > 0 {
				results = slices.Grow(results, len(exports))
				for _, e := range exports {
					if !excludePackages.Has(e.PackageName) {
						results = append(results, e)
					}
				}
			} else {
				results = append(results, exports...)
			}

			// As we go up the directory tree, exclude packages found in lower node_modules
			excludePackages = excludePackages.UnionedWith(nodeModulesBucket.PackageNames)
		}
		return nil, false
	})
	return results
}

type FixAndExport struct {
	Fix    *Fix
	Export *Export
}

func (v *View) GetCompletions(ctx context.Context, prefix string, forJSX bool, isTypeOnlyLocation bool) []*FixAndExport {
	results := v.Search(prefix, QueryKindWordPrefix)

	type exportGroupKey struct {
		target                     ExportID
		name                       string
		ambientModuleOrPackageName string
	}
	grouped := make(map[exportGroupKey][]*Export, len(results))
	for _, e := range results {
		name := e.Name()
		if forJSX && !(unicode.IsUpper(rune(name[0])) || e.IsRenameable()) {
			continue
		}
		target := e.ExportID
		if e.Target != (ExportID{}) {
			target = e.Target
		}
		key := exportGroupKey{
			target:                     target,
			name:                       name,
			ambientModuleOrPackageName: core.FirstNonZero(e.AmbientModuleName(), e.PackageName),
		}
		if existing, ok := grouped[key]; ok {
			for i, ex := range existing {
				if e.ExportID == ex.ExportID {
					grouped[key] = slices.Replace(existing, i, i+1, &Export{
						ExportID:                   e.ExportID,
						Syntax:                     e.Syntax,
						Flags:                      e.Flags | ex.Flags,
						ScriptElementKind:          min(e.ScriptElementKind, ex.ScriptElementKind),
						ScriptElementKindModifiers: *e.ScriptElementKindModifiers.UnionedWith(&ex.ScriptElementKindModifiers),
						localName:                  e.localName,
						Target:                     e.Target,
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
		return v.CompareFixes(a.Fix, b.Fix)
	}

	for _, exps := range grouped {
		fixesForGroup := make([]*FixAndExport, 0, len(exps))
		for _, e := range exps {
			for _, fix := range v.GetFixes(ctx, e, forJSX, isTypeOnlyLocation, nil) {
				fixesForGroup = append(fixesForGroup, &FixAndExport{
					Fix:    fix,
					Export: e,
				})
			}
		}
		fixes = append(fixes, core.MinAllFunc(fixesForGroup, compareFixes)...)
	}

	return fixes
}

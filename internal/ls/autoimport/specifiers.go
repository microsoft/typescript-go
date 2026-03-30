package autoimport

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/modulespecifiers"
)

func (v *View) GetModuleSpecifier(
	export *Export,
	userPreferences modulespecifiers.UserPreferences,
) (string, modulespecifiers.ResultKind) {
	// Ambient module
	if modulespecifiers.PathIsBareSpecifier(string(export.ModuleID)) {
		specifier := string(export.ModuleID)
		if modulespecifiers.IsExcludedByRegex(specifier, userPreferences.AutoImportSpecifierExcludeRegexes) {
			return "", modulespecifiers.ResultKindNone
		}
		return string(export.ModuleID), modulespecifiers.ResultKindAmbient
	}

	moduleFileName := export.ModuleFileName
	isSymlinkedPackageExport := false
	if export.PackageName != "" && moduleFileName != "" {
		realModuleFileName := v.program.Host().FS().Realpath(moduleFileName)
		isSymlinkedPackageExport = realModuleFileName != "" && realModuleFileName != moduleFileName
		if isSymlinkedPackageExport {
			moduleFileName = realModuleFileName
		}
	}

	if export.PackageName != "" && !isSymlinkedPackageExport {
		if entrypoints, ok := v.registry.entrypoints[export.Path]; ok {
			for _, entrypoint := range entrypoints {
				if entrypoint.IncludeConditions.IsSubsetOf(v.conditions) && !v.conditions.Intersects(entrypoint.ExcludeConditions) {
					specifier := modulespecifiers.ProcessEntrypointEnding(
						entrypoint,
						userPreferences,
						v.program,
						v.program.Options(),
						v.importingFile,
						v.getAllowedEndings(),
					)

					if !modulespecifiers.IsExcludedByRegex(specifier, userPreferences.AutoImportSpecifierExcludeRegexes) {
						return specifier, modulespecifiers.ResultKindNodeModules
					}
				}
			}
			return "", modulespecifiers.ResultKindNone
		}
	}

	cache := v.registry.specifierCache[v.importingFile.Path()]
	if export.PackageName == "" {
		if specifier, ok := cache.Load(export.Path); ok {
			if specifier == "" {
				return "", modulespecifiers.ResultKindNone
			}
			return specifier, modulespecifiers.ResultKindRelative
		}
	}

	specifiers, kind := modulespecifiers.GetModuleSpecifiersForFileWithInfo(
		v.importingFile,
		moduleFileName,
		v.program.Options(),
		v.program,
		userPreferences,
		modulespecifiers.ModuleSpecifierOptions{},
		true,
	)
	// !!! unsure when this could return multiple specifiers combined with the
	//     new node_modules code. Possibly with local symlinks, which should be
	//     very rare.
	for _, specifier := range specifiers {
		if strings.Contains(specifier, "/node_modules/") {
			continue
		}
		cache.Store(export.Path, specifier)
		return specifier, kind
	}
	cache.Store(export.Path, "")
	return "", modulespecifiers.ResultKindNone
}

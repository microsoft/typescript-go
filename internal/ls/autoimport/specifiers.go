package autoimport

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (v *View) GetModuleSpecifier(
	export *Export,
	userPreferences modulespecifiers.UserPreferences,
) (string, modulespecifiers.ResultKind) {
	// !!! try using existing import

	// Ambient module
	if modulespecifiers.PathIsBareSpecifier(string(export.ModuleID)) {
		specifier := string(export.ModuleID)
		if modulespecifiers.IsExcludedByRegex(specifier, userPreferences.AutoImportSpecifierExcludeRegexes) {
			return "", modulespecifiers.ResultKindNone
		}
		return string(export.ModuleID), modulespecifiers.ResultKindAmbient
	}

	if export.NodeModulesDirectory != "" {
		if entrypoints, ok := v.registry.nodeModules[export.NodeModulesDirectory].Entrypoints[export.Path]; ok {
			conditions := collections.NewSetFromItems(module.GetConditions(v.program.Options(), v.program.GetDefaultResolutionModeForFile(v.importingFile))...)
			for _, entrypoint := range entrypoints {
				if entrypoint.IncludeConditions.IsSubsetOf(conditions) && !conditions.Intersects(entrypoint.ExcludeConditions) {
					// !!! modulespecifiers.processEnding
					var specifier string
					switch entrypoint.Ending {
					case module.EndingFixed:
						specifier = entrypoint.ModuleSpecifier
					case module.EndingExtensionChangeable:
						dtsExtension := tspath.GetDeclarationFileExtension(entrypoint.ModuleSpecifier)
						if dtsExtension != "" {
							specifier = tspath.ChangeAnyExtension(entrypoint.ModuleSpecifier, modulespecifiers.GetJSExtensionForDeclarationFileExtension(dtsExtension), []string{dtsExtension}, false)
						} else {
							specifier = entrypoint.ModuleSpecifier
						}
					default:
						// !!! definitely wrong, lazy
						specifier = tspath.ChangeAnyExtension(entrypoint.ModuleSpecifier, "", []string{tspath.ExtensionDts, tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionJs, tspath.ExtensionJsx}, false)
					}

					if !modulespecifiers.IsExcludedByRegex(specifier, userPreferences.AutoImportSpecifierExcludeRegexes) {
						return specifier, modulespecifiers.ResultKindNodeModules
					}
				}
			}
			return "", modulespecifiers.ResultKindNone
		}
	}

	cache := v.registry.specifierCache[v.importingFile.Path()]
	if export.NodeModulesDirectory == "" {
		if specifier, ok := cache[export.Path]; ok {
			return specifier, modulespecifiers.ResultKindRelative
		}
	}

	specifiers, kind := modulespecifiers.GetModuleSpecifiersForFileWithInfo(
		v.importingFile,
		string(export.ExportID.ModuleID),
		v.program.Options(),
		v.program,
		userPreferences,
		modulespecifiers.ModuleSpecifierOptions{},
		true,
	)
	if len(specifiers) > 0 {
		// !!! sort/filter specifiers?
		specifier := specifiers[0]
		cache[export.Path] = specifier
		return specifier, kind
	}
	return "", modulespecifiers.ResultKindNone
}

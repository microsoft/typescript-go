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
		return string(export.ModuleID), modulespecifiers.ResultKindAmbient
	}

	if export.NodeModulesDirectory != "" {
		if entrypoints, ok := v.registry.nodeModules[export.NodeModulesDirectory].Entrypoints[export.Path]; ok {
			conditions := collections.NewSetFromItems(module.GetConditions(v.program.Options(), v.program.GetDefaultResolutionModeForFile(v.importingFile))...)
			for _, entrypoint := range entrypoints {
				if entrypoint.IncludeConditions.IsSubsetOf(conditions) && !conditions.Intersects(entrypoint.ExcludeConditions) {
					// !!! modulespecifiers.processEnding
					switch entrypoint.Ending {
					case module.EndingFixed:
						return entrypoint.ModuleSpecifier, modulespecifiers.ResultKindNodeModules
					case module.EndingExtensionChangeable:
						dtsExtension := tspath.GetDeclarationFileExtension(entrypoint.ModuleSpecifier)
						if dtsExtension != "" {
							return tspath.ChangeAnyExtension(entrypoint.ModuleSpecifier, modulespecifiers.GetJSExtensionForDeclarationFileExtension(dtsExtension), []string{dtsExtension}, false), modulespecifiers.ResultKindNodeModules
						}
						return entrypoint.ModuleSpecifier, modulespecifiers.ResultKindNodeModules
					default:
						// !!! definitely wrong, lazy
						return tspath.ChangeAnyExtension(entrypoint.ModuleSpecifier, "", []string{tspath.ExtensionDts, tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionJs, tspath.ExtensionJsx}, false), modulespecifiers.ResultKindNodeModules
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

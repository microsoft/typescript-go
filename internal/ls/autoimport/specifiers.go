package autoimport

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (v *View) GetModuleSpecifier(
	export *RawExport,
	userPreferences modulespecifiers.UserPreferences,
) string {
	// !!! try using existing import

	// Ambient module
	if modulespecifiers.PathIsBareSpecifier(string(export.ModuleID)) {
		return string(export.ModuleID)
	}

	if export.NodeModulesDirectory != "" {
		if entrypoints, ok := v.registry.nodeModules[export.NodeModulesDirectory].Entrypoints[export.Path]; ok {
			conditions := collections.NewSetFromItems(module.GetConditions(v.program.Options(), v.program.GetDefaultResolutionModeForFile(v.importingFile))...)
			for _, entrypoint := range entrypoints {
				if entrypoint.IncludeConditions.IsSubsetOf(conditions) && !conditions.Intersects(entrypoint.ExcludeConditions) {
					// !!! modulespecifiers.processEnding
					switch entrypoint.Ending {
					case module.EndingFixed:
						return entrypoint.ModuleSpecifier
					case module.EndingExtensionChangeable:
						dtsExtension := tspath.GetDeclarationFileExtension(entrypoint.ModuleSpecifier)
						if dtsExtension != "" {
							return tspath.ChangeAnyExtension(entrypoint.ModuleSpecifier, modulespecifiers.GetJSExtensionForDeclarationFileExtension(dtsExtension), []string{dtsExtension}, false)
						}
						return entrypoint.ModuleSpecifier
					default:
						// !!! definitely wrong, lazy
						return tspath.ChangeAnyExtension(entrypoint.ModuleSpecifier, "", []string{tspath.ExtensionDts, tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionJs, tspath.ExtensionJsx}, false)
					}
				}
			}
			return ""
		}
	}

	cache := v.registry.relativeSpecifierCache[v.importingFile.Path()]
	if export.NodeModulesDirectory == "" {
		if specifier, ok := cache[export.Path]; ok {
			return specifier
		}
	}

	specifiers, _ := modulespecifiers.GetModuleSpecifiersForFileWithInfo(
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
		return specifier
	}
	return ""
}

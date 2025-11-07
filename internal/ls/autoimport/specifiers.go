package autoimport

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
)

func (v *View) GetModuleSpecifier(
	export *RawExport,
	userPreferences modulespecifiers.UserPreferences,
) string {
	// !!! try using existing import

	if export.NodeModulesDirectory != "" {
		if entrypoints, ok := v.registry.nodeModules[export.NodeModulesDirectory].Entrypoints[export.Path]; ok {
			conditions := collections.NewSetFromItems(module.GetConditions(v.program.Options(), v.program.GetDefaultResolutionModeForFile(v.importingFile))...)
			for _, entrypoint := range entrypoints {
				if entrypoint.IncludeConditions.IsSubsetOf(conditions) && !conditions.Intersects(entrypoint.ExcludeConditions) {
					return entrypoint.ModuleSpecifier
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
		export.FileName,
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

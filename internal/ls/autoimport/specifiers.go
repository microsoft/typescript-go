package autoimport

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tspath"
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

	if export.NodeModulesDirectory != "" {
		bucket := v.registry.nodeModules[export.NodeModulesDirectory]
		if entrypoints, ok := bucket.Entrypoints[export.Path]; ok {
			// Get the package name from one of the entrypoints
			var packageNameForExport string
			if len(entrypoints) > 0 {
				packageNameForExport, _ = splitPackageSpecifier(entrypoints[0].ModuleSpecifier)
			}

			// Check if this package has exports using the precomputed set
			packageHasExports := bucket.PackagesWithExports != nil && bucket.PackagesWithExports.Has(packageNameForExport)

			for _, entrypoint := range entrypoints {
				// Skip EndingChangeable entrypoints only when the package has exports.
				// These files were discovered by reading the package directory directly
				// (not from exports), so they should not be suggested for auto-import
				// when the package has an exports field that doesn't expose them.
				// When the package doesn't have exports, EndingChangeable files are valid.
				if packageHasExports && entrypoint.Ending == module.EndingChangeable {
					continue
				}
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
		}
		// If the export is from a node_modules package but has no valid entrypoints,
		// it cannot be imported (e.g., internal files not exposed via package.json exports).
		return "", modulespecifiers.ResultKindNone
	}

	cache := v.registry.specifierCache[v.importingFile.Path()]
	if specifier, ok := cache.Load(export.Path); ok {
		if specifier == "" {
			return "", modulespecifiers.ResultKindNone
		}
		return specifier, modulespecifiers.ResultKindRelative
	}

	specifiers, kind := modulespecifiers.GetModuleSpecifiersForFileWithInfo(
		v.importingFile,
		export.ModuleFileName,
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
		// Check if this is a package specifier (starts with @ or doesn't start with ./ or ../)
		// that points to a non-entrypoint file in a symlinked package
		if v.isInvalidPackageSpecifier(specifier, export) {
			continue
		}
		cache.Store(export.Path, specifier)
		return specifier, kind
	}
	cache.Store(export.Path, "")
	return "", modulespecifiers.ResultKindNone
}

// isInvalidPackageSpecifier checks if a specifier is a bare package specifier
// (like @foo/bar/something or foo/something) that points to a file that is not
// a valid entrypoint in any node_modules bucket. This can happen when a symlinked
// monorepo package has internal files that are part of the program but should not
// be importable from outside the package via the package specifier.
func (v *View) isInvalidPackageSpecifier(specifier string, export *Export) bool {
	// Only check bare specifiers (not relative paths)
	if tspath.PathIsRelative(specifier) {
		return false
	}

	// Extract the package name from the specifier
	// For scoped packages like @foo/bar/path, package name is @foo/bar
	// For non-scoped packages like foo/path, package name is foo
	packageName, subpath := splitPackageSpecifier(specifier)
	if packageName == "" {
		return false
	}

	// If there's no subpath, this is a root import like "@foo/bar" or "foo"
	// These should be allowed as they go through normal resolution
	if subpath == "" {
		return false
	}

	// There's a subpath like "@foo/bar/something" - check if this is a valid entrypoint
	// Check if the export's file is a valid entrypoint in any node_modules bucket
	filePath := export.Path
	for _, bucket := range v.registry.nodeModules {
		if _, ok := bucket.Entrypoints[filePath]; ok {
			// File is a valid entrypoint
			return false
		}
	}

	// File is not a valid entrypoint, so this subpath specifier should not be generated
	// unless it can be resolved via the package's exports field (which it can't, or
	// it would have been handled by tryGetModuleNameFromExports in tryGetModuleNameAsNodeModule)
	return true
}

// splitPackageSpecifier splits a bare package specifier into package name and subpath.
// For "@scope/pkg/path", returns ("@scope/pkg", "path")
// For "pkg/path", returns ("pkg", "path")
// For "@scope/pkg" or "pkg", returns the package name and empty subpath
func splitPackageSpecifier(specifier string) (packageName string, subpath string) {
	if specifier == "" {
		return "", ""
	}

	// Handle scoped packages (@scope/package)
	if specifier[0] == '@' {
		// Find the second slash (after @scope/package)
		firstSlash := strings.Index(specifier, "/")
		if firstSlash == -1 {
			// Just "@scope" - not valid, but return as-is
			return specifier, ""
		}
		secondSlash := strings.Index(specifier[firstSlash+1:], "/")
		if secondSlash == -1 {
			// "@scope/package" - no subpath
			return specifier, ""
		}
		// "@scope/package/subpath"
		packageEnd := firstSlash + 1 + secondSlash
		return specifier[:packageEnd], specifier[packageEnd+1:]
	}

	// Non-scoped package
	firstSlash := strings.Index(specifier, "/")
	if firstSlash == -1 {
		// "package" - no subpath
		return specifier, ""
	}
	// "package/subpath"
	return specifier[:firstSlash], specifier[firstSlash+1:]
}

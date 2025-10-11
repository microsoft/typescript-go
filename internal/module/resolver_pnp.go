package module

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/module/pnp"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (r *resolutionState) loadPNPResolutionPath(moduleName string) (string, error) {
	resolution, err := pnp.ResolveToUnqualified(moduleName, r.containingDirectory, r.resolver.pnpResolutionConfig)
	if err != nil {
		return "", err
	}

	// trim trailing slash makes a bug in packageJsonInfoCache.Set
	// like @emotion/react/ -> @emotion/react after packageJsonInfoCache.Set
	// check why it's happening in packageJsonInfoCache and need to fix it
	return strings.TrimSuffix(resolution.Path, "/"), nil
}

func (r *resolutionState) loadModuleFromPNP(extensions extensions, typesScopeOnly bool) *resolved {
	packageName, rest := ParsePackageName(r.name)

	if !typesScopeOnly {
		pnpPath, err := r.loadPNPResolutionPath(packageName)

		if err == nil && r.resolver.host.FS().DirectoryExists(pnpPath) {
			candidate := tspath.NormalizePath(tspath.CombinePaths(pnpPath, rest))
			packageDirectory := pnpPath

			if result := r.loadModuleFromSpecificNodeModulesDirectory(extensions, candidate, packageDirectory, rest, true); !result.shouldContinueSearching() {
				return result
			}
		}
	}

	if extensions&extensionsDeclaration != 0 {
		typesPackageName := "@types/" + r.mangleScopedPackageName(packageName)
		pnpTypesPath, err := r.loadPNPResolutionPath(typesPackageName)
		if err == nil && r.resolver.host.FS().DirectoryExists(pnpTypesPath) {
			candidate := tspath.NormalizePath(tspath.CombinePaths(pnpTypesPath, rest))
			packageDirectory := pnpTypesPath

			if result := r.loadModuleFromSpecificNodeModulesDirectory(extensionsDeclaration, candidate, packageDirectory, rest, true); !result.shouldContinueSearching() {
				return result
			}
		}
	}

	return continueSearching()
}

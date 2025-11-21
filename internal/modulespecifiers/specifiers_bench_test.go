package modulespecifiers

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/symlinks"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func populateSymlinkCacheFromResolutions(importingFileName string, host *benchHost, compilerOptions *core.CompilerOptions, options ModuleSpecifierOptions, links *symlinks.KnownSymlinks) {
	packageJsonDir := host.GetNearestAncestorDirectoryWithPackageJson(tspath.GetDirectoryPath(importingFileName))
	if packageJsonDir == "" {
		return
	}

	packageJsonPath := tspath.CombinePaths(packageJsonDir, "package.json")

	pkgJsonInfo := host.GetPackageJsonInfo(packageJsonPath)
	if pkgJsonInfo == nil {
		return
	}

	pkgJson := pkgJsonInfo.GetContents()
	if pkgJson == nil {
		return
	}

	cwd := host.GetCurrentDirectory()
	caseSensitive := host.UseCaseSensitiveFileNames()

	// Helper to resolve dependencies without creating intermediate slices
	resolveDeps := func(deps map[string]string) {
		for depName := range deps {
			resolved := host.ResolveModuleName(depName, packageJsonPath, options.OverrideImportMode)
			if resolved != nil && resolved.OriginalPath != "" && resolved.ResolvedFileName != "" {
				processResolution(links, resolved.OriginalPath, resolved.ResolvedFileName, cwd, caseSensitive)
			}
		}
	}

	if deps, ok := pkgJson.Dependencies.GetValue(); ok {
		resolveDeps(deps)
	}
	if peerDeps, ok := pkgJson.PeerDependencies.GetValue(); ok {
		resolveDeps(peerDeps)
	}
	if optionalDeps, ok := pkgJson.OptionalDependencies.GetValue(); ok {
		resolveDeps(optionalDeps)
	}
}

func processResolution(links *symlinks.KnownSymlinks, originalPath string, resolvedFileName string, cwd string, caseSensitive bool) {
	originalPathKey := tspath.ToPath(originalPath, cwd, caseSensitive)
	links.SetFile(originalPathKey, resolvedFileName)

	commonResolved, commonOriginal := guessDirectorySymlink(originalPath, resolvedFileName, cwd, caseSensitive)
	if commonResolved != "" && commonOriginal != "" {
		symlinkPath := tspath.ToPath(commonOriginal, cwd, caseSensitive)
		if !tspath.ContainsIgnoredPath(string(symlinkPath)) {
			realPath := tspath.ToPath(commonResolved, cwd, caseSensitive)
			links.SetDirectory(
				commonOriginal,
				symlinkPath.EnsureTrailingDirectorySeparator(),
				&symlinks.KnownDirectoryLink{
					Real:     tspath.EnsureTrailingDirectorySeparator(commonResolved),
					RealPath: realPath.EnsureTrailingDirectorySeparator(),
				},
			)
		}
	}
}

func guessDirectorySymlink(originalPath string, resolvedFileName string, cwd string, caseSensitive bool) (string, string) {
	aParts := tspath.GetPathComponents(tspath.GetNormalizedAbsolutePath(resolvedFileName, cwd), "")
	bParts := tspath.GetPathComponents(tspath.GetNormalizedAbsolutePath(originalPath, cwd), "")
	isDirectory := false
	for len(aParts) >= 2 && len(bParts) >= 2 &&
		!isNodeModulesOrScopedPackageDirectory(aParts[len(aParts)-2], caseSensitive) &&
		!isNodeModulesOrScopedPackageDirectory(bParts[len(bParts)-2], caseSensitive) &&
		tspath.GetCanonicalFileName(aParts[len(aParts)-1], caseSensitive) == tspath.GetCanonicalFileName(bParts[len(bParts)-1], caseSensitive) {
		aParts = aParts[:len(aParts)-1]
		bParts = bParts[:len(bParts)-1]
		isDirectory = true
	}
	if isDirectory {
		return tspath.GetPathFromPathComponents(aParts), tspath.GetPathFromPathComponents(bParts)
	}
	return "", ""
}

func isNodeModulesOrScopedPackageDirectory(s string, caseSensitive bool) bool {
	return s != "" && (tspath.GetCanonicalFileName(s, caseSensitive) == "node_modules" || strings.HasPrefix(s, "@"))
}

type benchHost struct {
	mockModuleSpecifierGenerationHost
	resolveCount int
	packageJson  *packagejson.InfoCacheEntry
}

func (h *benchHost) ResolveModuleName(moduleName string, containingFile string, resolutionMode core.ResolutionMode) *module.ResolvedModule {
	h.resolveCount++
	return &module.ResolvedModule{
		ResolvedFileName: "/real/node_modules/" + moduleName + "/index.js",
		OriginalPath:     "/project/node_modules/" + moduleName + "/index.js",
	}
}

func (h *benchHost) GetPackageJsonInfo(pkgJsonPath string) *packagejson.InfoCacheEntry {
	return h.packageJson
}

func (h *benchHost) GetNearestAncestorDirectoryWithPackageJson(dirname string) string {
	return "/project"
}

type mockPackageJsonInfo struct {
	deps map[string]string
}

func (p *mockPackageJsonInfo) GetDirectory() string {
	return "/project"
}

func (p *mockPackageJsonInfo) GetContents() *packagejson.PackageJson {
	pkgJson := &packagejson.PackageJson{}
	pkgJson.Dependencies = packagejson.ExpectedOf(p.deps)
	return pkgJson
}

func BenchmarkPopulateSymlinkCacheFromResolutions(b *testing.B) {
	deps := make(map[string]string, 50)
	for i := range 50 {
		depName := "package-" + string(rune('a'+(i%26)))
		if i >= 26 {
			depName = depName + string(rune('a'+((i-26)%26)))
		}
		deps[depName] = "^1.0.0"
	}

	host := &benchHost{
		mockModuleSpecifierGenerationHost: mockModuleSpecifierGenerationHost{
			currentDir:                "/project",
			useCaseSensitiveFileNames: true,
			symlinkCache:              symlinks.NewKnownSymlink("/project", true),
		},
		packageJson: &packagejson.InfoCacheEntry{
			PackageDirectory: "/project",
			Contents: &packagejson.PackageJson{
				Fields: packagejson.Fields{
					DependencyFields: packagejson.DependencyFields{
						Dependencies: packagejson.ExpectedOf(deps),
					},
				},
			},
		},
	}

	compilerOptions := &core.CompilerOptions{}
	options := ModuleSpecifierOptions{
		OverrideImportMode: core.ResolutionModeNone,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		host.symlinkCache = symlinks.NewKnownSymlink("/project", true)
		host.resolveCount = 0

		for j := range 10 {
			importingFile := "/project/src/file" + string(rune('0'+j)) + ".ts"
			populateSymlinkCacheFromResolutions(importingFile, host, compilerOptions, options, host.symlinkCache)
		}
	}
}

func BenchmarkGetAllModulePaths(b *testing.B) {
	deps := make(map[string]string, 20)
	for i := range 20 {
		deps["package-"+string(rune('a'+i))] = "^1.0.0"
	}

	host := &benchHost{
		mockModuleSpecifierGenerationHost: mockModuleSpecifierGenerationHost{
			currentDir:                "/project",
			useCaseSensitiveFileNames: true,
			symlinkCache:              symlinks.NewKnownSymlink("/project", true),
		},
		packageJson: &packagejson.InfoCacheEntry{
			PackageDirectory: "/project",
			Contents: &packagejson.PackageJson{
				Fields: packagejson.Fields{
					DependencyFields: packagejson.DependencyFields{
						Dependencies: packagejson.ExpectedOf(deps),
					},
				},
			},
		},
	}

	info := getInfo(
		"/project/src/index.ts",
		host,
	)

	compilerOptions := &core.CompilerOptions{}
	options := ModuleSpecifierOptions{
		OverrideImportMode: core.ResolutionModeNone,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		getAllModulePathsWorker(
			info,
			"/real/node_modules/package-a/index.js",
			host,
			compilerOptions,
			options,
		)
	}
}

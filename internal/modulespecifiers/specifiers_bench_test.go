package modulespecifiers

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/symlinks"
)

type benchHost struct {
	mockModuleSpecifierGenerationHost
	resolveCount int
	packageJson  PackageJsonInfo
}

func (h *benchHost) ResolveModuleName(moduleName string, containingFile string, resolutionMode core.ResolutionMode) *module.ResolvedModule {
	h.resolveCount++
	return &module.ResolvedModule{
		ResolvedFileName: "/real/node_modules/" + moduleName + "/index.js",
		OriginalPath:     "/project/node_modules/" + moduleName + "/index.js",
	}
}

func (h *benchHost) GetPackageJsonInfo(pkgJsonPath string) PackageJsonInfo {
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
		packageJson: &mockPackageJsonInfo{deps: deps},
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
		packageJson: &mockPackageJsonInfo{deps: deps},
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

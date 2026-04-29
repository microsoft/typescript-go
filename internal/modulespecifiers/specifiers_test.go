package modulespecifiers

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/symlinks"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// Mock host for testing
type mockModuleSpecifierGenerationHost struct {
	currentDir                string
	useCaseSensitiveFileNames bool
	symlinkCache              *symlinks.KnownSymlinks
	packageJsonDirectory      string
	packageJsonInfo           *packagejson.InfoCacheEntry
	defaultResolutionMode     core.ResolutionMode
	resolveModuleName         func(moduleName string, containingFile string, resolutionMode core.ResolutionMode) *module.ResolvedModule
	resolveModuleNameCalls    []resolveModuleNameCall
}

type resolveModuleNameCall struct {
	moduleName     string
	containingFile string
	resolutionMode core.ResolutionMode
}

func (h *mockModuleSpecifierGenerationHost) GetCurrentDirectory() string {
	return h.currentDir
}

func (h *mockModuleSpecifierGenerationHost) UseCaseSensitiveFileNames() bool {
	return h.useCaseSensitiveFileNames
}

func (h *mockModuleSpecifierGenerationHost) GetSymlinkCache() *symlinks.KnownSymlinks {
	return h.symlinkCache
}

func (h *mockModuleSpecifierGenerationHost) ResolveModuleName(moduleName string, containingFile string, resolutionMode core.ResolutionMode) *module.ResolvedModule {
	h.resolveModuleNameCalls = append(h.resolveModuleNameCalls, resolveModuleNameCall{
		moduleName:     moduleName,
		containingFile: containingFile,
		resolutionMode: resolutionMode,
	})
	if h.resolveModuleName != nil {
		return h.resolveModuleName(moduleName, containingFile, resolutionMode)
	}
	return nil
}

func (h *mockModuleSpecifierGenerationHost) GetGlobalTypingsCacheLocation() string {
	return ""
}

func (h *mockModuleSpecifierGenerationHost) CommonSourceDirectory() string {
	return h.currentDir
}

func (h *mockModuleSpecifierGenerationHost) GetProjectReferenceFromSource(path tspath.Path) *tsoptions.SourceOutputAndProjectReference {
	return nil
}

func (h *mockModuleSpecifierGenerationHost) GetRedirectTargets(path tspath.Path) []string {
	return nil
}

func (h *mockModuleSpecifierGenerationHost) GetSourceOfProjectReferenceIfOutputIncluded(file ast.HasFileName) string {
	return file.FileName()
}

func (h *mockModuleSpecifierGenerationHost) FileExists(path string) bool {
	return true // Mock implementation
}

func (h *mockModuleSpecifierGenerationHost) GetNearestAncestorDirectoryWithPackageJson(dirname string) string {
	return h.packageJsonDirectory
}

func (h *mockModuleSpecifierGenerationHost) GetPackageJsonInfo(pkgJsonPath string) *packagejson.InfoCacheEntry {
	return h.packageJsonInfo
}

func (h *mockModuleSpecifierGenerationHost) GetDefaultResolutionModeForFile(file ast.HasFileName) core.ResolutionMode {
	return h.defaultResolutionMode
}

func (h *mockModuleSpecifierGenerationHost) GetResolvedModuleFromModuleSpecifier(file ast.HasFileName, moduleSpecifier *ast.StringLiteralLike) *module.ResolvedModule {
	return nil
}

func (h *mockModuleSpecifierGenerationHost) GetModeForUsageLocation(file ast.HasFileName, moduleSpecifier *ast.StringLiteralLike) core.ResolutionMode {
	return core.ResolutionModeNone
}

func TestGetEachFileNameOfModule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		importingFile  string
		importedFile   string
		preferSymlinks bool
		expectedCount  int
		expectedPaths  []string
	}{
		{
			name:           "basic file path",
			importingFile:  "/project/src/main.ts",
			importedFile:   "/project/lib/utils.ts",
			preferSymlinks: false,
			expectedCount:  1,
			expectedPaths:  []string{"/project/lib/utils.ts"},
		},
		{
			name:           "symlink preference false",
			importingFile:  "/project/src/main.ts",
			importedFile:   "/project/lib/utils.ts",
			preferSymlinks: false,
			expectedCount:  1,
		},
		{
			name:           "symlink preference true",
			importingFile:  "/project/src/main.ts",
			importedFile:   "/project/lib/utils.ts",
			preferSymlinks: true,
			expectedCount:  1,
		},
		{
			name:           "ignored path with no alternatives",
			importingFile:  "/project/src/main.ts",
			importedFile:   "/project/node_modules/.pnpm/file.ts",
			preferSymlinks: false,
			expectedCount:  1, // Should return 1 because there's no better option (all paths are ignored)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			host := &mockModuleSpecifierGenerationHost{
				currentDir:                "/project",
				useCaseSensitiveFileNames: true,
				symlinkCache:              symlinks.NewKnownSymlink("/project", true),
			}

			result := GetEachFileNameOfModule(tt.importingFile, tt.importedFile, host, tt.preferSymlinks)

			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d paths, got %d", tt.expectedCount, len(result))
			}

			if tt.expectedPaths != nil {
				for i, expectedPath := range tt.expectedPaths {
					if i >= len(result) {
						t.Errorf("Expected path %d: %s, but result has only %d paths", i, expectedPath, len(result))
						continue
					}
					if result[i].FileName != expectedPath {
						t.Errorf("Expected path %d to be %s, got %s", i, expectedPath, result[i].FileName)
					}
				}
			}

			for i, path := range result {
				if path.FileName == "" {
					t.Errorf("Path %d has empty FileName", i)
				}
			}
		})
	}
}

func TestGetEachFileNameOfModuleWithSymlinks(t *testing.T) {
	t.Parallel()
	host := &mockModuleSpecifierGenerationHost{
		currentDir:                "/project",
		useCaseSensitiveFileNames: true,
		symlinkCache:              symlinks.NewKnownSymlink("/project", true),
	}

	symlinkPath := tspath.ToPath("/project/symlink", "/project", true).EnsureTrailingDirectorySeparator()
	realDirectory := &symlinks.KnownDirectoryLink{
		Real:     "/real/path/",
		RealPath: tspath.ToPath("/real/path", "/project", true).EnsureTrailingDirectorySeparator(),
	}
	host.symlinkCache.SetDirectory("/project/symlink", symlinkPath, realDirectory)

	result := GetEachFileNameOfModule("/project/src/main.ts", "/real/path/file.ts", host, true)

	// Should find the symlink path
	found := false
	for _, path := range result {
		if path.FileName == "/project/symlink/file.ts" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find symlink path /project/symlink/file.ts")
	}
}

func TestGetAllModulePathsSeedsRuntimeDepsOncePerResolutionMode(t *testing.T) {
	t.Parallel()

	fields, err := packagejson.Parse([]byte(`{"dependencies":{"dep":"1.0.0"}}`))
	if err != nil {
		t.Fatal(err)
	}

	host := &mockModuleSpecifierGenerationHost{
		currentDir:                "/project",
		useCaseSensitiveFileNames: true,
		symlinkCache:              symlinks.NewKnownSymlink("/project", true),
		packageJsonDirectory:      "/project",
		packageJsonInfo: &packagejson.InfoCacheEntry{
			PackageDirectory: "/project",
			Contents: &packagejson.PackageJson{
				Fields:    fields,
				Parseable: true,
			},
		},
		defaultResolutionMode: core.ResolutionModeESM,
		resolveModuleName: func(moduleName string, containingFile string, resolutionMode core.ResolutionMode) *module.ResolvedModule {
			return &module.ResolvedModule{
				OriginalPath:     "/project/node_modules/dep/index.d.ts",
				ResolvedFileName: "/workspace/dep/index.d.ts",
			}
		},
	}
	source := ast.NewHasFileName("/project/src/main.mts", tspath.ToPath("/project/src/main.mts", "/project", true))
	info := getInfo(source, source.FileName(), host)

	getAllModulePathsWorker(info, "/workspace/dep/index.ts", host, &core.CompilerOptions{}, ModuleSpecifierOptions{})
	getAllModulePathsWorker(info, "/workspace/dep/index.ts", host, &core.CompilerOptions{}, ModuleSpecifierOptions{})

	if len(host.resolveModuleNameCalls) != 1 {
		t.Fatalf("expected one dependency resolution for repeated ESM calls, got %d", len(host.resolveModuleNameCalls))
	}
	if got := host.resolveModuleNameCalls[0].resolutionMode; got != core.ResolutionModeESM {
		t.Fatalf("expected default ESM resolution mode, got %v", got)
	}
	if got := host.resolveModuleNameCalls[0].containingFile; got != "/project/package.json" {
		t.Fatalf("expected package.json containing file, got %q", got)
	}

	getAllModulePathsWorker(info, "/workspace/dep/index.ts", host, &core.CompilerOptions{}, ModuleSpecifierOptions{OverrideImportMode: core.ResolutionModeCommonJS})

	if len(host.resolveModuleNameCalls) != 2 {
		t.Fatalf("expected a second dependency resolution for a different mode, got %d", len(host.resolveModuleNameCalls))
	}
	if got := host.resolveModuleNameCalls[1].resolutionMode; got != core.ResolutionModeCommonJS {
		t.Fatalf("expected override CommonJS resolution mode, got %v", got)
	}
}

func TestContainsNodeModules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "contains node_modules",
			path:     "/project/node_modules/lodash/index.js",
			expected: true,
		},
		{
			name:     "does not contain node_modules",
			path:     "/project/src/utils.ts",
			expected: false,
		},
		{
			name:     "node_modules in middle",
			path:     "/project/packages/node_modules/pkg/file.js",
			expected: true,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ContainsNodeModules(tt.path)
			if result != tt.expected {
				t.Errorf("ContainsNodeModules(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestContainsIgnoredPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "ignored path",
			path:     "/project/node_modules/.pnpm/file.ts",
			expected: true,
		},
		{
			name:     "not ignored path",
			path:     "/project/src/file.ts",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := containsIgnoredPath(tt.path)
			if result != tt.expected {
				t.Errorf("containsIgnoredPath(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestTryGetRealFileNameForNonJSDeclarationFileName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		fileName string
		expected string
	}{
		{
			name:     "json declaration file",
			fileName: "/project/foo.d.json.ts",
			expected: "/project/foo.json",
		},
		{
			name:     "multi-dot source extension declaration file",
			fileName: "/project/foo.module.d.css.ts",
			expected: "/project/foo.module.css",
		},
		{
			name:     "plain dts file ignored",
			fileName: "/project/foo.d.ts",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := TryGetRealFileNameForNonJSDeclarationFileName(tt.fileName); got != tt.expected {
				t.Errorf("TryGetRealFileNameForNonJSDeclarationFileName(%q) = %q, expected %q", tt.fileName, got, tt.expected)
			}
		})
	}
}

func TestTryGetModuleNameFromExportsOrImports(t *testing.T) {
	t.Parallel()
	t.Run("with exports pattern", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			targetFilePath string
			expected       string
		}{
			{
				name:           "match",
				targetFilePath: "/pkg/src/things/thing1/index.ts",
				expected:       "./src/things/thing1",
			},
			{
				name:           "mismatch with matching leading and trailing strings",
				targetFilePath: "/pkg/src/things/index.ts",
				expected:       "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result := tryGetModuleNameFromExportsOrImports(
					&core.CompilerOptions{},
					&mockModuleSpecifierGenerationHost{},
					tt.targetFilePath,
					"/pkg",
					"./src/things/*",
					packagejson.ExportsOrImports{
						JSONValue: packagejson.JSONValue{
							Type:  packagejson.JSONValueTypeString,
							Value: "./src/things/*/index.js",
						},
					},
					[]string{},
					MatchingModePattern,
					false,
					false,
				)
				if result != tt.expected {
					t.Errorf("tryGetModuleNameFromExportsOrImports(targetFilePath = %q) = %v, expected %v", tt.targetFilePath, result, tt.expected)
				}
			})
		}
	})
}

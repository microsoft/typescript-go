package autoimport_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/autoimporttestutil"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"gotest.tools/v3/assert"
)

func BenchmarkAutoImportRegistry_TypeScript(b *testing.B) {
	checkerURI := lsconv.FileNameToDocumentURI(tspath.CombinePaths(repo.TypeScriptSubmodulePath, "src/compiler/checker.ts"))
	checkerContent, ok := osvfs.FS().ReadFile(checkerURI.FileName())
	assert.Assert(b, ok, "failed to read checker.ts")

	for b.Loop() {
		b.StopTimer()
		session, _ := projecttestutil.SetupWithRealFS()
		session.DidOpenFile(context.Background(), checkerURI, 1, checkerContent, lsproto.LanguageKindTypeScript)
		b.StartTimer()

		_, err := session.GetLanguageServiceWithAutoImports(context.Background(), checkerURI)
		assert.NilError(b, err)
	}
}

func BenchmarkAutoImportRegistry_VSCode(b *testing.B) {
	mainURI := lsproto.DocumentUri("file:///Users/andrew/Developer/microsoft/vscode/src/main.ts")
	mainContent, ok := osvfs.FS().ReadFile(mainURI.FileName())
	assert.Assert(b, ok, "failed to read main.ts")

	for b.Loop() {
		b.StopTimer()
		session, _ := projecttestutil.SetupWithRealFS()
		session.DidOpenFile(context.Background(), mainURI, 1, mainContent, lsproto.LanguageKindTypeScript)
		b.StartTimer()

		_, err := session.GetLanguageServiceWithAutoImports(context.Background(), mainURI)
		assert.NilError(b, err)
	}
}

func TestRegistryLifecycle(t *testing.T) {
	t.Parallel()
	t.Run("preparesProjectAndNodeModulesBuckets", func(t *testing.T) {
		t.Parallel()
		fixture := autoimporttestutil.SetupLifecycleSession(t, lifecycleProjectRoot, 1)
		session := fixture.Session()
		project := fixture.SingleProject()
		mainFile := project.File(0)
		session.DidOpenFile(context.Background(), mainFile.URI(), 1, mainFile.Content(), lsproto.LanguageKindTypeScript)

		stats := autoImportStats(t, session)
		projectBucket := singleBucket(t, stats.ProjectBuckets)
		nodeModulesBucket := singleBucket(t, stats.NodeModulesBuckets)
		assert.Equal(t, true, projectBucket.State.Dirty())
		assert.Equal(t, 0, projectBucket.FileCount)
		assert.Equal(t, true, nodeModulesBucket.State.Dirty())
		assert.Equal(t, 0, nodeModulesBucket.FileCount)

		_, err := session.GetLanguageServiceWithAutoImports(context.Background(), mainFile.URI())
		assert.NilError(t, err)

		stats = autoImportStats(t, session)
		projectBucket = singleBucket(t, stats.ProjectBuckets)
		nodeModulesBucket = singleBucket(t, stats.NodeModulesBuckets)
		assert.Equal(t, false, projectBucket.State.Dirty())
		assert.Assert(t, projectBucket.ExportCount > 0)
		assert.Equal(t, false, nodeModulesBucket.State.Dirty())
		assert.Assert(t, nodeModulesBucket.ExportCount > 0)
	})

	t.Run("marksProjectBucketDirtyAfterEdit", func(t *testing.T) {
		t.Parallel()
		fixture := autoimporttestutil.SetupLifecycleSession(t, lifecycleProjectRoot, 2)
		session := fixture.Session()
		utils := fixture.Utils()
		project := fixture.SingleProject()
		mainFile := project.File(0)
		secondaryFile := project.File(1)
		session.DidOpenFile(context.Background(), mainFile.URI(), 1, mainFile.Content(), lsproto.LanguageKindTypeScript)
		session.DidOpenFile(context.Background(), secondaryFile.URI(), 1, secondaryFile.Content(), lsproto.LanguageKindTypeScript)
		_, err := session.GetLanguageServiceWithAutoImports(context.Background(), mainFile.URI())
		assert.NilError(t, err)

		updatedContent := mainFile.Content() + "// change\n"
		session.DidChangeFile(context.Background(), mainFile.URI(), 2, []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: updatedContent}},
		})

		_, err = session.GetLanguageService(context.Background(), mainFile.URI())
		assert.NilError(t, err)

		stats := autoImportStats(t, session)
		projectBucket := singleBucket(t, stats.ProjectBuckets)
		nodeModulesBucket := singleBucket(t, stats.NodeModulesBuckets)
		assert.Equal(t, projectBucket.State.Dirty(), true)
		assert.Equal(t, projectBucket.State.DirtyFile(), utils.ToPath(mainFile.FileName()))
		assert.Equal(t, nodeModulesBucket.State.Dirty(), false)
		assert.Equal(t, nodeModulesBucket.State.DirtyFile(), tspath.Path(""))

		// Bucket should not recompute when requesting same file changed
		_, err = session.GetLanguageServiceWithAutoImports(context.Background(), mainFile.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		projectBucket = singleBucket(t, stats.ProjectBuckets)
		assert.Equal(t, projectBucket.State.Dirty(), true)
		assert.Equal(t, projectBucket.State.DirtyFile(), utils.ToPath(mainFile.FileName()))

		// Bucket should recompute when other file has changed
		session.DidChangeFile(context.Background(), secondaryFile.URI(), 1, []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: "// new content"}},
		})
		_, err = session.GetLanguageServiceWithAutoImports(context.Background(), mainFile.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		projectBucket = singleBucket(t, stats.ProjectBuckets)
		assert.Equal(t, projectBucket.State.Dirty(), false)
	})

	t.Run("packageJsonDependencyChangesInvalidateNodeModulesBuckets", func(t *testing.T) {
		t.Parallel()
		fixture := autoimporttestutil.SetupLifecycleSession(t, lifecycleProjectRoot, 1)
		session := fixture.Session()
		sessionUtils := fixture.Utils()
		project := fixture.SingleProject()
		mainFile := project.File(0)
		nodePackage := project.NodeModules()[0]
		packageJSON := project.PackageJSONFile()
		ctx := context.Background()

		session.DidOpenFile(ctx, mainFile.URI(), 1, mainFile.Content(), lsproto.LanguageKindTypeScript)
		_, err := session.GetLanguageServiceWithAutoImports(ctx, mainFile.URI())
		assert.NilError(t, err)
		stats := autoImportStats(t, session)
		nodeModulesBucket := singleBucket(t, stats.NodeModulesBuckets)
		assert.Equal(t, nodeModulesBucket.State.Dirty(), false)

		fs := sessionUtils.FS()
		updatePackageJSON := func(content string) {
			assert.NilError(t, fs.WriteFile(packageJSON.FileName(), content, false))
			session.DidChangeWatchedFiles(ctx, []*lsproto.FileEvent{
				{Type: lsproto.FileChangeTypeChanged, Uri: packageJSON.URI()},
			})
		}

		sameDepsContent := fmt.Sprintf("{\n  \"name\": \"local-project-stable\",\n  \"dependencies\": {\n    \"%s\": \"*\"\n  }\n}\n", nodePackage.Name)
		updatePackageJSON(sameDepsContent)
		_, err = session.GetLanguageService(ctx, mainFile.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		nodeModulesBucket = singleBucket(t, stats.NodeModulesBuckets)
		assert.Equal(t, nodeModulesBucket.State.Dirty(), false)

		differentDepsContent := fmt.Sprintf("{\n  \"name\": \"local-project-stable\",\n  \"dependencies\": {\n    \"%s\": \"*\",\n    \"newpkg\": \"*\"\n  }\n}\n", nodePackage.Name)
		updatePackageJSON(differentDepsContent)
		_, err = session.GetLanguageServiceWithAutoImports(ctx, mainFile.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		assert.Check(t, singleBucket(t, stats.NodeModulesBuckets).DependencyNames.Has("newpkg"))
	})

	t.Run("nodeModulesBucketsDeletedWhenNoOpenFilesReferThem", func(t *testing.T) {
		t.Parallel()
		fixture := autoimporttestutil.SetupMonorepoLifecycleSession(t, autoimporttestutil.MonorepoSetupConfig{
			Root: monorepoProjectRoot,
			MonorepoPackageTemplate: autoimporttestutil.MonorepoPackageTemplate{
				Name:            "monorepo",
				NodeModuleNames: []string{"pkg-root"},
			},
			Packages: []autoimporttestutil.MonorepoPackageConfig{
				{FileCount: 1, MonorepoPackageTemplate: autoimporttestutil.MonorepoPackageTemplate{Name: "package-a", NodeModuleNames: []string{"pkg-a"}}},
				{FileCount: 1, MonorepoPackageTemplate: autoimporttestutil.MonorepoPackageTemplate{Name: "package-b", NodeModuleNames: []string{"pkg-b"}}},
			},
		})
		session := fixture.Session()
		monorepo := fixture.Monorepo()
		pkgA := monorepo.Package(0)
		pkgB := monorepo.Package(1)
		fileA := pkgA.File(0)
		fileB := pkgB.File(0)
		ctx := context.Background()

		// Open file in package-a, should create buckets for root and package-a node_modules
		session.DidOpenFile(ctx, fileA.URI(), 1, fileA.Content(), lsproto.LanguageKindTypeScript)
		_, err := session.GetLanguageServiceWithAutoImports(ctx, fileA.URI())
		assert.NilError(t, err)

		// Open file in package-b, should also create buckets for package-b
		session.DidOpenFile(ctx, fileB.URI(), 1, fileB.Content(), lsproto.LanguageKindTypeScript)
		_, err = session.GetLanguageServiceWithAutoImports(ctx, fileB.URI())
		assert.NilError(t, err)
		stats := autoImportStats(t, session)
		assert.Equal(t, len(stats.NodeModulesBuckets), 3)
		assert.Equal(t, len(stats.ProjectBuckets), 2)

		// Close file in package-a, package-a's node_modules bucket and project bucket should be removed
		session.DidCloseFile(ctx, fileA.URI())
		_, err = session.GetLanguageServiceWithAutoImports(ctx, fileB.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		assert.Equal(t, len(stats.NodeModulesBuckets), 2)
		assert.Equal(t, len(stats.ProjectBuckets), 1)
	})

	t.Run("dependencyAggregationChangesAsFilesOpenAndClose", func(t *testing.T) {
		t.Parallel()
		monorepoRoot := "/home/src/monorepo"
		packageADir := tspath.CombinePaths(monorepoRoot, "packages", "a")
		monorepoIndex := tspath.CombinePaths(monorepoRoot, "index.js")
		packageAIndex := tspath.CombinePaths(packageADir, "index.js")

		fixture := autoimporttestutil.SetupMonorepoLifecycleSession(t, autoimporttestutil.MonorepoSetupConfig{
			Root: monorepoRoot,
			MonorepoPackageTemplate: autoimporttestutil.MonorepoPackageTemplate{
				Name:            "monorepo",
				NodeModuleNames: []string{"pkg1", "pkg2", "pkg3"},
				DependencyNames: []string{"pkg1"},
			},
			Packages: []autoimporttestutil.MonorepoPackageConfig{
				{
					FileCount: 0,
					MonorepoPackageTemplate: autoimporttestutil.MonorepoPackageTemplate{
						Name:            "a",
						DependencyNames: []string{"pkg1", "pkg2"},
					},
				},
			},
			ExtraFiles: []autoimporttestutil.TextFileSpec{
				{Path: monorepoIndex, Content: "export const monorepoIndex = 1;\n"},
				{Path: packageAIndex, Content: "export const pkgA = 2;\n"},
			},
		})
		session := fixture.Session()
		monorepoHandle := fixture.ExtraFile(monorepoIndex)
		packageAHandle := fixture.ExtraFile(packageAIndex)

		ctx := context.Background()

		// Open monorepo root file: expect dependencies restricted to pkg1
		session.DidOpenFile(ctx, monorepoHandle.URI(), 1, monorepoHandle.Content(), lsproto.LanguageKindJavaScript)
		_, err := session.GetLanguageServiceWithAutoImports(ctx, monorepoHandle.URI())
		assert.NilError(t, err)
		stats := autoImportStats(t, session)
		assert.Assert(t, singleBucket(t, stats.NodeModulesBuckets).DependencyNames.Equals(collections.NewSetFromItems("pkg1")))

		// Open package-a file: pkg2 should be added to existing bucket
		session.DidOpenFile(ctx, packageAHandle.URI(), 1, packageAHandle.Content(), lsproto.LanguageKindJavaScript)
		_, err = session.GetLanguageServiceWithAutoImports(ctx, packageAHandle.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		assert.Assert(t, singleBucket(t, stats.NodeModulesBuckets).DependencyNames.Equals(collections.NewSetFromItems("pkg1", "pkg2")))

		// Close package-a file; only monorepo bucket should remain
		session.DidCloseFile(ctx, packageAHandle.URI())
		_, err = session.GetLanguageServiceWithAutoImports(ctx, monorepoHandle.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		assert.Assert(t, singleBucket(t, stats.NodeModulesBuckets).DependencyNames.Equals(collections.NewSetFromItems("pkg1")))

		// Close monorepo file; no node_modules buckets should remain
		session.DidCloseFile(ctx, monorepoHandle.URI())
		stats = autoImportStats(t, session)
		assert.Equal(t, len(stats.NodeModulesBuckets), 0)
	})
}

const (
	lifecycleProjectRoot = "/home/src/autoimport-lifecycle"
	monorepoProjectRoot  = "/home/src/autoimport-monorepo"
)

func autoImportStats(t *testing.T, session *project.Session) *autoimport.CacheStats {
	t.Helper()
	snapshot, release := session.Snapshot()
	defer release()
	registry := snapshot.AutoImportRegistry()
	if registry == nil {
		t.Fatal("auto import registry not initialized")
	}
	return registry.GetCacheStats()
}

func singleBucket(t *testing.T, buckets []autoimport.BucketStats) autoimport.BucketStats {
	t.Helper()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	return buckets[0]
}

package autoimport_test

import (
	"context"
	"fmt"
	"testing"

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
		assert.Equal(t, true, projectBucket.Dirty)
		assert.Equal(t, 0, projectBucket.FileCount)
		assert.Equal(t, true, nodeModulesBucket.Dirty)
		assert.Equal(t, 0, nodeModulesBucket.FileCount)

		_, err := session.GetLanguageServiceWithAutoImports(context.Background(), mainFile.URI())
		assert.NilError(t, err)

		stats = autoImportStats(t, session)
		projectBucket = singleBucket(t, stats.ProjectBuckets)
		nodeModulesBucket = singleBucket(t, stats.NodeModulesBuckets)
		assert.Equal(t, false, projectBucket.Dirty)
		assert.Assert(t, projectBucket.ExportCount > 0)
		assert.Equal(t, false, nodeModulesBucket.Dirty)
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
		assert.Equal(t, projectBucket.Dirty, true)
		assert.Equal(t, projectBucket.DirtyFile, utils.ToPath(mainFile.FileName()))
		assert.Equal(t, nodeModulesBucket.Dirty, false)
		assert.Equal(t, nodeModulesBucket.DirtyFile, tspath.Path(""))

		// Bucket should not recompute when requesting same file changed
		_, err = session.GetLanguageServiceWithAutoImports(context.Background(), mainFile.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		projectBucket = singleBucket(t, stats.ProjectBuckets)
		assert.Equal(t, projectBucket.Dirty, true)
		assert.Equal(t, projectBucket.DirtyFile, utils.ToPath(mainFile.FileName()))

		// Bucket should recompute when other file has changed
		session.DidChangeFile(context.Background(), secondaryFile.URI(), 1, []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: "// new content"}},
		})
		_, err = session.GetLanguageServiceWithAutoImports(context.Background(), mainFile.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		projectBucket = singleBucket(t, stats.ProjectBuckets)
		assert.Equal(t, projectBucket.Dirty, false)
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
		assert.Equal(t, nodeModulesBucket.Dirty, false)

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
		assert.Equal(t, nodeModulesBucket.Dirty, false)

		differentDepsContent := fmt.Sprintf("{\n  \"name\": \"local-project-stable\",\n  \"dependencies\": {\n    \"%s\": \"*\",\n    \"newpkg\": \"*\"\n  }\n}\n", nodePackage.Name)
		updatePackageJSON(differentDepsContent)
		_, err = session.GetLanguageServiceWithAutoImports(ctx, mainFile.URI())
		assert.NilError(t, err)
		stats = autoImportStats(t, session)
		assert.Check(t, singleBucket(t, stats.NodeModulesBuckets).DependencyNames.Has("newpkg"))
	})
}

const lifecycleProjectRoot = "/home/src/autoimport-lifecycle"

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

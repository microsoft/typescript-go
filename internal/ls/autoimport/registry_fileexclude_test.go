package autoimport_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/autoimporttestutil"
	"gotest.tools/v3/assert"
)

// TestAutoImportFileExcludePatternsRebuild tests that changing AutoImportFileExcludePatterns
// triggers a bucket rebuild, preventing the panic described in the issue where
// ErrNeedsAutoImports is returned even after calling GetLanguageServiceWithAutoImports.
//
// The bug: When autoImportFileExcludePatterns changes, possiblyNeedsRebuildForFile() returns true,
// but updateIndexes() doesn't check for changed exclude patterns, so the bucket isn't rebuilt.
// This causes IsPreparedForImportingFile to keep returning false, triggering ErrNeedsAutoImports
// even after GetLanguageServiceWithAutoImports is called, leading to a panic.
func TestAutoImportFileExcludePatternsRebuild(t *testing.T) {
	t.Parallel()

	// Setup a project with node_modules containing a package
	fixture := autoimporttestutil.SetupLifecycleSession(t, lifecycleProjectRoot, 1)
	session := fixture.Session()
	project := fixture.SingleProject()
	mainFile := project.File(0)

	ctx := context.Background()

	// Step 1: Open file and build auto-imports initially (no exclude patterns set)
	session.DidOpenFile(ctx, mainFile.URI(), 1, mainFile.Content(), lsproto.LanguageKindTypeScript)

	_, err := session.GetLanguageServiceWithAutoImports(ctx, mainFile.URI())
	assert.NilError(t, err)

	// Verify buckets are clean after initial build
	stats := autoImportStats(t, session)
	projectBucket := singleBucket(t, stats.ProjectBuckets)
	nodeModulesBucket := singleBucket(t, stats.NodeModulesBuckets)
	assert.Equal(t, false, projectBucket.State.Dirty(), "Project bucket should be clean after initial build")
	assert.Equal(t, false, nodeModulesBucket.State.Dirty(), "Node modules bucket should be clean after initial build")

	// Verify IsPreparedForImportingFile returns true with no exclude patterns
	snapshot, release := session.Snapshot()
	registry := snapshot.AutoImportRegistry()
	// Get the project from the snapshot to get the config path
	defaultProject := snapshot.GetDefaultProject(mainFile.URI())
	if defaultProject == nil {
		t.Fatal("No default project found for mainFile")
	}
	projectPath := defaultProject.ConfigFilePath()
	preferences := lsutil.NewDefaultUserPreferences()
	preferences.IncludeCompletionsForModuleExports = core.TSTrue
	preferences.IncludeCompletionsForImportStatements = core.TSTrue
	
	isPrepared := registry.IsPreparedForImportingFile(mainFile.FileName(), projectPath, preferences)
	release()
	assert.Assert(t, isPrepared, "IsPreparedForImportingFile should return true after initial build")

	// Step 2: Change the file exclude patterns preference
	// This simulates the user adding autoImportFileExcludePatterns in .vscode/settings.json
	newPreferences := lsutil.NewDefaultUserPreferences()
	newPreferences.IncludeCompletionsForModuleExports = core.TSTrue
	newPreferences.IncludeCompletionsForImportStatements = core.TSTrue
	newPreferences.AutoImportFileExcludePatterns = []string{"**/node_modules/**/*.d.ts"}

	// Update preferences via Configure
	newConfig := lsutil.NewUserConfig(newPreferences)
	session.Configure(newConfig)

	// Step 3: Check that IsPreparedForImportingFile now returns false
	// because the exclude patterns changed
	snapshot2, release2 := session.Snapshot()
	registry2 := snapshot2.AutoImportRegistry()
	isPrepared2 := registry2.IsPreparedForImportingFile(mainFile.FileName(), projectPath, newPreferences)
	release2()
	assert.Assert(t, !isPrepared2, "IsPreparedForImportingFile should return false after preferences change")

	// Step 4: Call GetLanguageServiceWithAutoImports - this SHOULD rebuild the buckets
	// BUG: Currently, updateIndexes() doesn't check for changed fileExcludePatterns,
	// so the buckets won't be rebuilt
	_, err = session.GetLanguageServiceWithAutoImports(ctx, mainFile.URI())
	assert.NilError(t, err)

	// Step 5: Verify that IsPreparedForImportingFile now returns true
	// This is where the bug manifests: if the bucket wasn't rebuilt,
	// IsPreparedForImportingFile will still return false, causing the panic
	snapshot3, release3 := session.Snapshot()
	registry3 := snapshot3.AutoImportRegistry()
	isPrepared3 := registry3.IsPreparedForImportingFile(mainFile.FileName(), projectPath, newPreferences)
	stats3 := registry3.GetCacheStats()
	release3()

	// Log bucket states for debugging
	if !isPrepared3 {
		t.Logf("BUG REPRODUCED: IsPreparedForImportingFile still returns false after GetLanguageServiceWithAutoImports")
		projectBucket3 := singleBucket(t, stats3.ProjectBuckets)
		nodeModulesBucket3 := singleBucket(t, stats3.NodeModulesBuckets)
		t.Logf("Project bucket dirty: %v", projectBucket3.State.Dirty())
		t.Logf("Node modules bucket dirty: %v", nodeModulesBucket3.State.Dirty())
		t.Logf("Project bucket fileExcludePatterns: %v", projectBucket3.State)
	}

	assert.Assert(t, isPrepared3, "IsPreparedForImportingFile should return true after GetLanguageServiceWithAutoImports rebuilds buckets")
}

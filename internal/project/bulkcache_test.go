package project_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestBulkCacheInvalidation(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Base file structure for testing
	baseFiles := map[string]any{
		"/project/tsconfig.json": `{
			"compilerOptions": {
				"strict": true,
				"target": "es2015"
			},
			"include": ["src/**/*"]
		}`,
		"/project/src/index.ts":     `import { helper } from "./helper"; console.log(helper);`,
		"/project/src/helper.ts":    `export const helper = "test";`,
		"/project/src/utils/lib.ts": `export function util() { return "util"; }`,
	}

	t.Run("large number of node_modules changes invalidates only node_modules cache", func(t *testing.T) {
		t.Parallel()
		session, utils := projecttestutil.Setup(baseFiles)

		// Open a file to create the project
		session.DidOpenFile(context.Background(), "file:///project/src/index.ts", 1, baseFiles["/project/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// Get initial snapshot and verify config
		ls, err := session.GetLanguageService(context.Background(), "file:///project/src/index.ts")
		assert.NilError(t, err)
		assert.Equal(t, ls.GetProgram().Options().Target, core.ScriptTargetES2015)

		snapshotBefore, release := session.Snapshot()
		defer release()
		configBefore := snapshotBefore.ConfigFileRegistry

		// Create excessive changes in node_modules (1001 changes to exceed threshold)
		fileEvents := generateFileEvents(1001, "file:///project/node_modules/generated/file%d.js", lsproto.FileChangeTypeCreated)

		// Update tsconfig.json on disk to test that configs don't get reloaded
		err = utils.FS().WriteFile("/project/tsconfig.json", `{
			"compilerOptions": {
				"strict": true,
				"target": "esnext"
			},
			"include": ["src/**/*"]
		}`, false)
		assert.NilError(t, err)

		// Process the excessive node_modules changes
		session.DidChangeWatchedFiles(context.Background(), fileEvents)

		// Get language service again to trigger snapshot update
		ls, err = session.GetLanguageService(context.Background(), "file:///project/src/index.ts")
		assert.NilError(t, err)

		snapshotAfter, release := session.Snapshot()
		defer release()
		configAfter := snapshotAfter.ConfigFileRegistry

		// Config should NOT have been reloaded (target should remain ES2015, not esnext)
		assert.Equal(t, ls.GetProgram().Options().Target, core.ScriptTargetES2015, "Config should not have been reloaded for node_modules-only changes")

		// Config registry should be the same instance (no configs reloaded)
		assert.Equal(t, configBefore, configAfter, "Config registry should not have changed for node_modules-only changes")
	})

	t.Run("large number of changes outside node_modules causes config reload", func(t *testing.T) {
		t.Parallel()
		session, utils := projecttestutil.Setup(baseFiles)

		// Open a file to create the project
		session.DidOpenFile(context.Background(), "file:///project/src/index.ts", 1, baseFiles["/project/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// Get initial state
		ls, err := session.GetLanguageService(context.Background(), "file:///project/src/index.ts")
		assert.NilError(t, err)
		assert.Equal(t, ls.GetProgram().Options().Target, core.ScriptTargetES2015)

		snapshotBefore, release := session.Snapshot()
		defer release()

		// Update tsconfig.json on disk
		err = utils.FS().WriteFile("/project/tsconfig.json", `{
			"compilerOptions": {
				"strict": true,
				"target": "esnext"
			},
			"include": ["src/**/*"]
		}`, false)
		assert.NilError(t, err)
		// Add root file
		err = utils.FS().WriteFile("/project/src/rootFile.ts", `console.log("root file")`, false)
		assert.NilError(t, err)

		// Create excessive changes outside node_modules (1001 changes to exceed threshold)
		fileEvents := generateFileEvents(1001, "file:///project/generated/file%d.ts", lsproto.FileChangeTypeCreated)

		// Process the excessive changes outside node_modules
		session.DidChangeWatchedFiles(context.Background(), fileEvents)

		// Get language service again to trigger snapshot update
		ls, err = session.GetLanguageService(context.Background(), "file:///project/src/index.ts")
		assert.NilError(t, err)

		snapshotAfter, release := session.Snapshot()
		defer release()

		// Config SHOULD have been reloaded (target should now be esnext and new root file present)
		assert.Equal(t, ls.GetProgram().Options().Target, core.ScriptTargetESNext, "Config should have been reloaded for changes outside node_modules")
		assert.Check(t, ls.GetProgram().GetSourceFile("/project/src/rootFile.ts") != nil, "New root file should be present")

		// Snapshots should be different
		assert.Assert(t, snapshotBefore != snapshotAfter, "Snapshot should have changed after bulk invalidation outside node_modules")
	})

	t.Run("large number of changes outside node_modules causes project reevaluation", func(t *testing.T) {
		t.Parallel()
		session, utils := projecttestutil.Setup(baseFiles)

		// Open a file that will initially use the root tsconfig
		session.DidOpenFile(context.Background(), "file:///project/src/utils/lib.ts", 1, baseFiles["/project/src/utils/lib.ts"].(string), lsproto.LanguageKindTypeScript)

		// Initially, the file should use the root project (strict mode)
		snapshot, release := session.Snapshot()
		defer release()
		initialProject := snapshot.GetDefaultProject("file:///project/src/utils/lib.ts")
		assert.Equal(t, initialProject.Name(), "/project/tsconfig.json", "Should initially use root tsconfig")

		// Get language service to verify initial strict mode
		ls, err := session.GetLanguageService(context.Background(), "file:///project/src/utils/lib.ts")
		assert.NilError(t, err)
		assert.Equal(t, ls.GetProgram().Options().Strict, core.TSTrue, "Should initially use strict mode from root config")

		// Now create the nested tsconfig (this would normally be detected, but we'll simulate a missed event)
		err = utils.FS().WriteFile("/project/src/utils/tsconfig.json", `{
			"compilerOptions": {
				"strict": false,
				"target": "esnext"
			}
		}`, false)
		assert.NilError(t, err)

		// Create excessive changes outside node_modules to trigger bulk invalidation
		// This simulates the scenario where the nested tsconfig creation was missed in the flood of changes
		fileEvents := generateFileEvents(1001, "file:///project/src/generated/file%d.ts", lsproto.FileChangeTypeCreated)

		// Process the excessive changes - this should trigger project reevaluation
		session.DidChangeWatchedFiles(context.Background(), fileEvents)

		// Get language service - this should now find the nested config and switch projects
		ls, err = session.GetLanguageService(context.Background(), "file:///project/src/utils/lib.ts")
		assert.NilError(t, err)

		snapshot, release = session.Snapshot()
		defer release()
		newProject := snapshot.GetDefaultProject("file:///project/src/utils/lib.ts")

		// The file should now use the nested tsconfig
		assert.Equal(t, newProject.Name(), "/project/src/utils/tsconfig.json", "Should now use nested tsconfig after bulk invalidation")
		assert.Equal(t, ls.GetProgram().Options().Strict, core.TSFalse, "Should now use non-strict mode from nested config")
		assert.Equal(t, ls.GetProgram().Options().Target, core.ScriptTargetESNext, "Should use esnext target from nested config")
	})

	t.Run("excessive changes only in node_modules does not affect config file names cache", func(t *testing.T) {
		t.Parallel()
		testConfigFileNamesCacheBehavior(t, "file:///project/node_modules/generated/file%d.js", false, "node_modules changes should not clear config cache")
	})

	t.Run("excessive changes outside node_modules clears config file names cache", func(t *testing.T) {
		t.Parallel()
		testConfigFileNamesCacheBehavior(t, "file:///project/generated/file%d.ts", true, "non-node_modules changes should clear config cache")
	})
}

// Helper function to generate excessive file change events
func generateFileEvents(count int, pathTemplate string, changeType lsproto.FileChangeType) []*lsproto.FileEvent {
	var events []*lsproto.FileEvent
	for i := range count {
		events = append(events, &lsproto.FileEvent{
			Uri:  lsproto.DocumentUri(fmt.Sprintf(pathTemplate, i)),
			Type: changeType,
		})
	}
	return events
}

// Helper function to test config file names cache behavior
func testConfigFileNamesCacheBehavior(t *testing.T, eventPathTemplate string, expectConfigDiscovery bool, testName string) {
	files := map[string]any{
		"/project/src/index.ts": `console.log("test");`, // No tsconfig initially
	}
	session, utils := projecttestutil.Setup(files)

	// Open file without tsconfig - should create inferred project
	session.DidOpenFile(context.Background(), "file:///project/src/index.ts", 1, files["/project/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

	snapshot, release := session.Snapshot()
	defer release()
	assert.Assert(t, snapshot.ProjectCollection.InferredProject() != nil, "Should have inferred project")
	assert.Equal(t, snapshot.GetDefaultProject("file:///project/src/index.ts").Kind, project.KindInferred)

	// Create a tsconfig that would affect this file (simulating a missed creation event)
	err := utils.FS().WriteFile("/project/tsconfig.json", `{
		"compilerOptions": {
			"strict": true
		},
		"include": ["src/**/*"]
	}`, false)
	assert.NilError(t, err)

	// Create excessive changes to trigger bulk invalidation
	fileEvents := generateFileEvents(1001, eventPathTemplate, lsproto.FileChangeTypeCreated)

	// Process the changes
	session.DidChangeWatchedFiles(context.Background(), fileEvents)

	// Get language service to trigger config discovery
	_, err = session.GetLanguageService(context.Background(), "file:///project/src/index.ts")
	assert.NilError(t, err)

	snapshot, release = session.Snapshot()
	defer release()
	newProject := snapshot.GetDefaultProject("file:///project/src/index.ts")

	// Check expected behavior
	if expectConfigDiscovery {
		// Should now use configured project instead of inferred
		assert.Equal(t, newProject.Kind, project.KindConfigured, "Should now use configured project after cache invalidation")
		assert.Equal(t, newProject.Name(), "/project/tsconfig.json", "Should use the newly discovered tsconfig")
	} else {
		// Should still use inferred project (config file names cache not cleared)
		assert.Assert(t, newProject == snapshot.ProjectCollection.InferredProject(), "Should still use inferred project after node_modules-only changes")
	}
}

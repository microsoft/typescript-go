package project_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestProjectLifetime(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	t.Run("configured project", func(t *testing.T) {
		t.Parallel()
		files := map[string]any{
			"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true,
					"module": "nodenext",
					"strict": true
				},
				"include": ["src"]
			}`,
			"/home/projects/TS/p1/src/index.ts": `import { x } from "./x";`,
			"/home/projects/TS/p1/src/x.ts":     `export const x = 1;`,
			"/home/projects/TS/p1/config.ts":    `let x = 1, y = 2;`,
			"/home/projects/TS/p2/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true,
					"module": "nodenext",
					"strict": true
				},
				"include": ["src"]
			}`,
			"/home/projects/TS/p2/src/index.ts": `import { x } from "./x";`,
			"/home/projects/TS/p2/src/x.ts":     `export const x = 1;`,
			"/home/projects/TS/p2/config.ts":    `let x = 1, y = 2;`,
			"/home/projects/TS/p3/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true,
					"module": "nodenext",
					"strict": true
				},
				"include": ["src"]
			}`,
			"/home/projects/TS/p3/src/index.ts": `import { x } from "./x";`,
			"/home/projects/TS/p3/src/x.ts":     `export const x = 1;`,
			"/home/projects/TS/p3/config.ts":    `let x = 1, y = 2;`,
		}
		session, utils := projecttestutil.Setup(files)
		snapshot := session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 0)

		// Open files in two projects
		uri1 := lsproto.DocumentUri("file:///home/projects/TS/p1/src/index.ts")
		uri2 := lsproto.DocumentUri("file:///home/projects/TS/p2/src/index.ts")
		session.DidOpenFile(context.Background(), uri1, 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		session.DidOpenFile(context.Background(), uri2, 1, files["/home/projects/TS/p2/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		session.WaitForBackgroundTasks()
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 2)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p2/tsconfig.json")) != nil)
		assert.Equal(t, len(utils.Client().WatchFilesCalls()), 1)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p2/tsconfig.json")) != nil)

		// Close p1 file and open p3 file
		session.DidCloseFile(context.Background(), uri1)
		uri3 := lsproto.DocumentUri("file:///home/projects/TS/p3/src/index.ts")
		session.DidOpenFile(context.Background(), uri3, 1, files["/home/projects/TS/p3/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		session.WaitForBackgroundTasks()
		// Should still have two projects, but p1 replaced by p3
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 2)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) == nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p2/tsconfig.json")) != nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p3/tsconfig.json")) != nil)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p1/tsconfig.json")) == nil)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p2/tsconfig.json")) != nil)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p3/tsconfig.json")) != nil)
		assert.Equal(t, len(utils.Client().WatchFilesCalls()), 1)
		assert.Equal(t, len(utils.Client().UnwatchFilesCalls()), 0)

		// Close p2 and p3 files, open p1 file again
		session.DidCloseFile(context.Background(), uri2)
		session.DidCloseFile(context.Background(), uri3)
		session.DidOpenFile(context.Background(), uri1, 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		session.WaitForBackgroundTasks()
		// Should have one project (p1)
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p2/tsconfig.json")) == nil)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p3/tsconfig.json")) == nil)
		assert.Equal(t, len(utils.Client().WatchFilesCalls()), 1)
		assert.Equal(t, len(utils.Client().UnwatchFilesCalls()), 0)
	})

	t.Run("unrooted inferred projects", func(t *testing.T) {
		t.Parallel()
		files := map[string]any{
			"/home/projects/TS/p1/src/index.ts": `import { x } from "./x";`,
			"/home/projects/TS/p1/src/x.ts":     `export const x = 1;`,
			"/home/projects/TS/p1/config.ts":    `let x = 1, y = 2;`,
			"/home/projects/TS/p2/src/index.ts": `import { x } from "./x";`,
			"/home/projects/TS/p2/src/x.ts":     `export const x = 1;`,
			"/home/projects/TS/p2/config.ts":    `let x = 1, y = 2;`,
			"/home/projects/TS/p3/src/index.ts": `import { x } from "./x";`,
			"/home/projects/TS/p3/src/x.ts":     `export const x = 1;`,
			"/home/projects/TS/p3/config.ts":    `let x = 1, y = 2;`,
		}
		session, _ := projecttestutil.Setup(files)
		snapshot := session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 0)

		// Open files without workspace roots (empty string) - should create single inferred project
		uri1 := lsproto.DocumentUri("file:///home/projects/TS/p1/src/index.ts")
		uri2 := lsproto.DocumentUri("file:///home/projects/TS/p2/src/index.ts")
		session.DidOpenFile(context.Background(), uri1, 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		session.DidOpenFile(context.Background(), uri2, 1, files["/home/projects/TS/p2/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// Should have one inferred project
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() != nil)

		// Close p1 file and open p3 file
		session.DidCloseFile(context.Background(), uri1)
		uri3 := lsproto.DocumentUri("file:///home/projects/TS/p3/src/index.ts")
		session.DidOpenFile(context.Background(), uri3, 1, files["/home/projects/TS/p3/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// Should still have one inferred project
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() != nil)

		// Close p2 and p3 files, open p1 file again
		session.DidCloseFile(context.Background(), uri2)
		session.DidCloseFile(context.Background(), uri3)
		session.DidOpenFile(context.Background(), uri1, 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// Should still have one inferred project
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() != nil)
	})

	t.Run("file moves from inferred to configured project", func(t *testing.T) {
		t.Parallel()
		files := map[string]any{
			"/home/projects/ts/foo.ts": `export const foo = 1;`,
			"/home/projects/ts/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true,
					"module": "nodenext",
					"strict": true
				},
				"include": ["main.ts"]
			}`,
			"/home/projects/ts/p1/main.ts": `import { foo } from "../foo"; console.log(foo);`,
		}
		session, _ := projecttestutil.Setup(files)

		// Open foo.ts first - should create inferred project since no tsconfig found initially
		fooUri := lsproto.DocumentUri("file:///home/projects/ts/foo.ts")
		session.DidOpenFile(context.Background(), fooUri, 1, files["/home/projects/ts/foo.ts"].(string), lsproto.LanguageKindTypeScript)

		// Should have one inferred project
		snapshot := session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() != nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) == nil)

		// Now open main.ts - should trigger discovery of tsconfig.json and move foo.ts to configured project
		mainUri := lsproto.DocumentUri("file:///home/projects/ts/p1/main.ts")
		session.DidOpenFile(context.Background(), mainUri, 1, files["/home/projects/ts/p1/main.ts"].(string), lsproto.LanguageKindTypeScript)

		// Should now have one configured project and no inferred project
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() == nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)

		// Config file should be present
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)

		// Close main.ts - configured project should remain because foo.ts is still open
		session.DidCloseFile(context.Background(), mainUri)
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)

		// Close foo.ts - configured project should be retained until next file open
		session.DidCloseFile(context.Background(), fooUri)
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ConfigFileRegistry.GetConfig(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)
	})

	t.Run("file move from inferred to configured via didOpen/didClose sequence", func(t *testing.T) {
		t.Parallel()
		// Start with tsconfig.json that includes "src" but file is at root level
		files := map[string]any{
			"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"include": ["src"]
			}`,
			"/home/projects/TS/p1/index.ts": `export const x = 1;`,
		}
		session, utils := projecttestutil.Setup(files)

		// Open index.ts at root level - should create inferred project since it's not under src/
		// Creates config file registry entry, but has no files
		indexUri := lsproto.DocumentUri("file:///home/projects/TS/p1/index.ts")
		session.DidOpenFile(context.Background(), indexUri, 1, files["/home/projects/TS/p1/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// Should have one inferred project only (file is not included by tsconfig)
		snapshot := session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() != nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) == nil)

		// Simulate file move: create src/index.ts on disk
		err := utils.FS().WriteFile("/home/projects/TS/p1/src/index.ts", files["/home/projects/TS/p1/index.ts"].(string))
		assert.NilError(t, err)
		err = utils.FS().Remove("/home/projects/TS/p1/index.ts")
		assert.NilError(t, err)

		// Simulate file move sequence as it would happen in an editor:
		// 1. didOpen src/index.ts (new location)
		// Open comes in before file create event, so the config file is not marked as needing a file name reload,
		// so it's not turned into a configured project yet. This is probably not ideal, but it should sort itself
		// out momentarily after the file watcher events are processed. When we try the config file, we mark it
		// as "retained by src/index.ts" so the config entry doesn't get deleted before src/index.ts is closed.
		// Even though we currently think src/index.ts doesn't belong to the config, the config is in its directory
		// path, so we'll always see it as a candidate for containing src/index.ts.
		srcIndexUri := lsproto.DocumentUri("file:///home/projects/TS/p1/src/index.ts")
		session.DidOpenFile(context.Background(), srcIndexUri, 1, files["/home/projects/TS/p1/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// 2. didClose index.ts (old location)
		session.DidCloseFile(context.Background(), indexUri)

		// 3. didChangeWatchedFiles: create src/index.ts and delete index.ts
		// The creation event for src/index.ts now hits the config file registry, and we should notice we
		// got a creation event for a file that retained the config, triggering a filename reload.
		session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
			{
				Uri:  srcIndexUri,
				Type: lsproto.FileChangeTypeCreated,
			},
			{
				Uri:  indexUri,
				Type: lsproto.FileChangeTypeDeleted,
			},
		})

		// Should now have one configured project only (file is now under src/)
		_, err = session.GetLanguageService(context.Background(), srcIndexUri)
		assert.NilError(t, err)
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() == nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)
	})

	t.Run("tsconfig move from subdirectory to parent via didChangeWatchedFiles", func(t *testing.T) {
		t.Parallel()
		// Start with tsconfig.json in src/ that includes "src" - file won't be included initially
		files := map[string]any{
			"/home/projects/TS/p1/src/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"include": ["src"]
			}`,
			"/home/projects/TS/p1/src/index.ts": `export const x = 1;`,
		}
		session, utils := projecttestutil.Setup(files)

		// Open src/index.ts - should create inferred project since tsconfig.json includes "src"
		// relative to its location (src/src/ which doesn't exist)
		indexUri := lsproto.DocumentUri("file:///home/projects/TS/p1/src/index.ts")
		session.DidOpenFile(context.Background(), indexUri, 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		// Should have one inferred project only (file is not included by tsconfig at src/tsconfig.json)
		snapshot := session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() != nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/src/tsconfig.json")) == nil)

		// Simulate tsconfig.json move: create tsconfig.json at parent level, delete from src/
		tsconfigContent := files["/home/projects/TS/p1/src/tsconfig.json"].(string)
		err := utils.FS().WriteFile("/home/projects/TS/p1/tsconfig.json", tsconfigContent)
		assert.NilError(t, err)
		err = utils.FS().Remove("/home/projects/TS/p1/src/tsconfig.json")
		assert.NilError(t, err)

		// Simulate file move via didChangeWatchedFiles
		newTsconfigUri := lsproto.DocumentUri("file:///home/projects/TS/p1/tsconfig.json")
		oldTsconfigUri := lsproto.DocumentUri("file:///home/projects/TS/p1/src/tsconfig.json")
		session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
			{
				Uri:  newTsconfigUri,
				Type: lsproto.FileChangeTypeCreated,
			},
			{
				Uri:  oldTsconfigUri,
				Type: lsproto.FileChangeTypeDeleted,
			},
		})

		// Should now have one configured project only (tsconfig.json now includes src/index.ts)
		_, err = session.GetLanguageService(context.Background(), indexUri)
		assert.NilError(t, err)
		snapshot = session.Snapshot()
		assert.Equal(t, len(snapshot.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapshot.ProjectCollection.InferredProject() == nil)
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/ts/p1/tsconfig.json")) != nil)
	})

	t.Run("workspace is always watched", func(t *testing.T) {
		t.Parallel()
		// Demonstrate what happens to LSP file watchers when we:
		//   1. Open a file in project A (watchers created)
		//   2. Close project A's file
		//   3. Open a file in project B (project A removed, project B added)
		//
		// When the workspace directory is "/" all project globs collapse to "/**/*",
		// so the ref-counted watcher is never removed. To exercise real watcher churn,
		// we set CurrentDirectory to a narrow directory so that projects outside it
		// produce distinct outside-workspace watcher globs.
		files := map[string]any{
			// Project A - inside the workspace
			"/home/projects/workspace/projA/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				}
			}`,
			"/home/projects/workspace/projA/a.ts": `export const a = 1;`,
			// Project B - outside the workspace
			"/external/projB/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				}
			}`,
			"/external/projB/b.ts": `export const b = 2;`,
		}
		session, utils := projecttestutil.SetupWithOptions(files, &project.SessionOptions{
			CurrentDirectory:       "/home/projects/workspace",
			DefaultLibraryPath:     bundled.LibPath(),
			PositionEncoding:       lsproto.PositionEncodingKindUTF8,
			WatchEnabled:           true,
			LoggingEnabled:         true,
			PushDiagnosticsEnabled: true,
		})

		// --- Step 1: Open project A ---
		uriA := lsproto.DocumentUri("file:///home/projects/workspace/projA/a.ts")
		session.DidOpenFile(context.Background(), uriA, 1, files["/home/projects/workspace/projA/a.ts"].(string), lsproto.LanguageKindTypeScript)
		_, err := session.GetLanguageService(context.Background(), uriA)
		assert.NilError(t, err)
		session.WaitForBackgroundTasks()

		watchCallsAfterOpenA := len(utils.Client().WatchFilesCalls())
		unwatchCallsAfterOpenA := len(utils.Client().UnwatchFilesCalls())
		assert.Assert(t, watchCallsAfterOpenA > 0, "expected at least one WatchFiles call after opening project A")
		assert.Equal(t, unwatchCallsAfterOpenA, 0, "expected no UnwatchFiles calls after opening project A")

		snapshot := session.Snapshot()
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/workspace/proja/tsconfig.json")) != nil,
			"project A should exist")

		// Record the watcher IDs registered for project A
		var watcherIDsA []project.WatcherID
		for _, call := range utils.Client().WatchFilesCalls() {
			watcherIDsA = append(watcherIDsA, call.ID)
		}

		// --- Step 2: Close project A, open project B ---
		session.DidCloseFile(context.Background(), uriA)
		uriB := lsproto.DocumentUri("file:///external/projB/b.ts")
		session.DidOpenFile(context.Background(), uriB, 1, files["/external/projB/b.ts"].(string), lsproto.LanguageKindTypeScript)
		_, err = session.GetLanguageService(context.Background(), uriB)
		assert.NilError(t, err)
		session.WaitForBackgroundTasks()

		snapshot = session.Snapshot()
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/workspace/proja/tsconfig.json")) == nil,
			"project A should be removed after closing its only open file and opening another project")
		assert.Assert(t, snapshot.ProjectCollection.ConfiguredProject(tspath.Path("/external/projb/tsconfig.json")) != nil,
			"project B should exist")

		watchCallsAfterSwitch := len(utils.Client().WatchFilesCalls())
		unwatchCallsAfterSwitch := len(utils.Client().UnwatchFilesCalls())

		// Record all watcher IDs that were unwatched
		var unwatchedIDs []project.WatcherID
		for _, call := range utils.Client().UnwatchFilesCalls() {
			unwatchedIDs = append(unwatchedIDs, call.ID)
		}

		// New watchers should have been created for project B
		assert.Assert(t, watchCallsAfterSwitch > watchCallsAfterOpenA,
			"expected new WatchFiles calls for project B, got %d total (was %d)", watchCallsAfterSwitch, watchCallsAfterOpenA)

		// Project A's watchers should have been cleaned up (globs no longer needed)
		assert.Assert(t, unwatchCallsAfterSwitch > 0,
			"expected UnwatchFiles calls when project A is removed, got %d", unwatchCallsAfterSwitch)

		// --- Step 3: Verify project B's files are watched ---
		assert.Assert(t, utils.WatchesFile("/external/projb/b.ts"),
			"project B's files should be watched")
		// Workspace should still be watched via the workspace watcher,
		// even though project A was removed.
		assert.Assert(t, utils.WatchesFile("/home/projects/workspace/proja/a.ts"),
			"workspace should still be watched via workspace watcher")
	})

	t.Run("deleted open file remains in project until closed", func(t *testing.T) {
		t.Parallel()
		// Scenario:
		// 1. Start with two files included by a tsconfig, both open
		// 2. In a single batch change, delete one of the files but leave it open, and create a new file included by the tsconfig
		// 3. Request a LS for the deleted but open file
		// 4. Project should include both the new file and the deleted open file
		// 5. Close the deleted file
		// 6. On next LS request, the project should exclude the deleted file

		files := map[string]any{
			"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"include": ["src"]
			}`,
			"/home/projects/TS/p1/src/index.ts": ``,
			"/home/projects/TS/p1/src/x.ts":     `export const x = 1;`,
		}
		session, utils := projecttestutil.Setup(files)

		// Step 1: Open both files
		indexUri := lsproto.DocumentUri("file:///home/projects/TS/p1/src/index.ts")
		xUri := lsproto.DocumentUri("file:///home/projects/TS/p1/src/x.ts")
		session.DidOpenFile(context.Background(), indexUri, 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		session.DidOpenFile(context.Background(), xUri, 1, files["/home/projects/TS/p1/src/x.ts"].(string), lsproto.LanguageKindTypeScript)

		// Verify initial state - both files should be in the project
		ls, err := session.GetLanguageService(context.Background(), indexUri)
		assert.NilError(t, err)
		program := ls.GetProgram()
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/index.ts") != nil, "index.ts should be in project")
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/x.ts") != nil, "x.ts should be in project")

		// Step 2: In a single batch change:
		// - Delete x.ts from disk (but leave it open)
		// - Create a new file y.ts on disk
		err = utils.FS().Remove("/home/projects/TS/p1/src/x.ts")
		assert.NilError(t, err)
		err = utils.FS().WriteFile("/home/projects/TS/p1/src/y.ts", `export const y = 2;`)
		assert.NilError(t, err)

		// Send both events in a single batch
		session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
			{
				Uri:  xUri,
				Type: lsproto.FileChangeTypeDeleted,
			},
			{
				Uri:  lsproto.DocumentUri("file:///home/projects/TS/p1/src/y.ts"),
				Type: lsproto.FileChangeTypeCreated,
			},
		})

		// Step 3 & 4: Request LS for the deleted but still open file
		// Project should include: index.ts, x.ts (open overlay), y.ts (new disk file)
		ls, err = session.GetLanguageService(context.Background(), xUri)
		assert.NilError(t, err)
		program = ls.GetProgram()
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/index.ts") != nil, "index.ts should still be in project")
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/x.ts") != nil, "x.ts should still be in project (open overlay)")
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/y.ts") != nil, "y.ts should be in project (new file)")

		// Step 5: Close the deleted file
		session.DidCloseFile(context.Background(), xUri)

		// Step 6: On next LS request, x.ts should be excluded
		ls, err = session.GetLanguageService(context.Background(), indexUri)
		assert.NilError(t, err)
		program = ls.GetProgram()
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/index.ts") != nil, "index.ts should still be in project")
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/x.ts") == nil, "x.ts should no longer be in project (closed and deleted)")
		assert.Assert(t, program.GetSourceFile("/home/projects/TS/p1/src/y.ts") != nil, "y.ts should still be in project")
	})
}

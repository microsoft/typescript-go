package project_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestCustomConfigFileName(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/src/tsconfig.json":     `{"compilerOptions": {"strict": false}}`,
		"/src/tsconfig.all.json": `{"compilerOptions": {"strict": true}}`,
		"/src/index.ts":          `export const x = 1;`,
	}
	uri := lsproto.DocumentUri("file:///src/index.ts")

	t.Run("picks up custom config and switches on preference change", func(t *testing.T) {
		t.Parallel()
		session, _ := projecttestutil.Setup(files)

		session.DidOpenFile(context.Background(), uri, 1, files["/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		ls, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)

		snapshot := session.Snapshot()
		assert.Equal(t, snapshot.GetDefaultProject(uri).Name(), "/src/tsconfig.json")
		assert.Equal(t, ls.GetProgram().Options().Strict, core.TSFalse)

		prefs := lsutil.NewDefaultUserPreferences()
		prefs.CustomConfigFileName = "tsconfig.all.json"
		session.Configure(lsutil.NewUserConfig(prefs))

		ls, err = session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)

		snapshot = session.Snapshot()
		assert.Equal(t, snapshot.GetDefaultProject(uri).Name(), "/src/tsconfig.all.json")
		assert.Equal(t, ls.GetProgram().Options().Strict, core.TSTrue)
	})

	t.Run("uses tsconfig.json when customConfigFileName is empty", func(t *testing.T) {
		t.Parallel()
		session, _ := projecttestutil.Setup(files)

		prefs := lsutil.NewDefaultUserPreferences()
		// default for CustomConfigFileName is "".
		assert.Equal(t, prefs.CustomConfigFileName, "")
		session.Configure(lsutil.NewUserConfig(prefs))

		session.DidOpenFile(context.Background(), uri, 1, files["/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		_, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)

		snapshot := session.Snapshot()
		assert.Equal(t, snapshot.GetDefaultProject(uri).Name(), "/src/tsconfig.json")
	})

	t.Run("falls back to tsconfig.json when custom config missing", func(t *testing.T) {
		t.Parallel()
		session, _ := projecttestutil.Setup(files)

		prefs := lsutil.NewDefaultUserPreferences()
		prefs.CustomConfigFileName = "tsconfig.nonexistent.json"
		session.Configure(lsutil.NewUserConfig(prefs))

		session.DidOpenFile(context.Background(), uri, 1, files["/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		_, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)

		snapshot := session.Snapshot()
		assert.Equal(t, snapshot.GetDefaultProject(uri).Name(), "/src/tsconfig.json")
	})

	t.Run("reverts to tsconfig.json when custom config preference is cleared", func(t *testing.T) {
		t.Parallel()
		session, _ := projecttestutil.Setup(files)

		// Step 1: Open file, verify it uses tsconfig.json (strict: false)
		session.DidOpenFile(context.Background(), uri, 1, files["/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		ls, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)

		snapshot := session.Snapshot()
		assert.Equal(t, snapshot.GetDefaultProject(uri).Name(), "/src/tsconfig.json")
		assert.Equal(t, ls.GetProgram().Options().Strict, core.TSFalse)

		// Step 2: Switch to custom config (strict: true)
		prefs := lsutil.NewDefaultUserPreferences()
		prefs.CustomConfigFileName = "tsconfig.all.json"
		session.Configure(lsutil.NewUserConfig(prefs))

		ls, err = session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)

		snapshot = session.Snapshot()
		assert.Equal(t, snapshot.GetDefaultProject(uri).Name(), "/src/tsconfig.all.json")
		assert.Equal(t, ls.GetProgram().Options().Strict, core.TSTrue)

		// Step 3: Clear custom config preference, should revert to tsconfig.json (strict: false)
		prefs = lsutil.NewDefaultUserPreferences()
		prefs.CustomConfigFileName = ""
		session.Configure(lsutil.NewUserConfig(prefs))

		ls, err = session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)

		snapshot = session.Snapshot()
		assert.Equal(t, snapshot.GetDefaultProject(uri).Name(), "/src/tsconfig.json")
		assert.Equal(t, ls.GetProgram().Options().Strict, core.TSFalse)
	})

	// This test demonstrates the bug reported in #2020: after changing
	// customConfigFileName, the server does not schedule a diagnostics refresh,
	// so the VS Code client never knows to re-pull diagnostics and shows stale results.
	t.Run("schedules diagnostics refresh when custom config preference changes", func(t *testing.T) {
		t.Parallel()
		session, utils := projecttestutil.Setup(files)

		session.DidOpenFile(context.Background(), uri, 1, files["/src/index.ts"].(string), lsproto.LanguageKindTypeScript)
		_, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)
		session.WaitForBackgroundTasks()

		// Record baseline refresh call count
		baselineRefreshCount := len(utils.Client().RefreshDiagnosticsCalls())

		// Change the custom config preference
		prefs := lsutil.NewDefaultUserPreferences()
		prefs.CustomConfigFileName = "tsconfig.all.json"
		session.Configure(lsutil.NewUserConfig(prefs))

		// GetLanguageService triggers the snapshot update with the new config
		_, err = session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)
		session.WaitForBackgroundTasks()

		// The server should have scheduled a diagnostics refresh to tell the client
		// to re-pull diagnostics with the new project configuration.
		refreshCount := len(utils.Client().RefreshDiagnosticsCalls())
		assert.Assert(t, refreshCount > baselineRefreshCount,
			"expected RefreshDiagnostics to be called after customConfigFileName change, got %d calls (baseline %d)",
			refreshCount, baselineRefreshCount)
	})
}

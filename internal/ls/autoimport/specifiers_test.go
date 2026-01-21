package autoimport_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestGetModuleSpecifierCrashWithUntitledFile(t *testing.T) {
	// Regression test for https://github.com/microsoft/typescript-go/issues/2550
	// Crash when requesting completions for an untitled file when there's also
	// a saved file with exports but no tsconfig.
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Set up a file with exports but no tsconfig
	files := map[string]any{
		"/home/src/project/utils.ts": `export function helper() {}
`,
	}

	session, _ := projecttestutil.Setup(files)
	t.Cleanup(session.Close)

	ctx := context.Background()

	// Open the file with exports
	session.DidOpenFile(ctx, "file:///home/src/project/utils.ts", 1, `export function helper() {}
`, lsproto.LanguageKindTypeScript)

	// Open an untitled file
	session.DidOpenFile(ctx, "untitled:Untitled-1", 1, "", lsproto.LanguageKindTypeScript)

	// Get language service with auto imports for the untitled file
	// This should not crash
	languageService, err := session.GetLanguageServiceWithAutoImports(ctx, "untitled:Untitled-1")
	assert.NilError(t, err)

	// Request completions - this previously caused a nil pointer dereference
	// in GetModuleSpecifier when accessing the specifierCache
	_, err = languageService.ProvideCompletion(
		ctx,
		"untitled:Untitled-1",
		lsproto.Position{Line: 0, Character: 0},
		nil,
	)
	assert.NilError(t, err)
}

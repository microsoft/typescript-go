package project_test

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestContentMapperInProject(t *testing.T) {
	t.Parallel()
	files := map[string]any{
		"/home/project/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true },
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`,
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.ExecName),
		"/home/project/app.box":                          "export const version = #{target};\n",
		"/home/project/main.ts":                          "import { version } from \"./app.box\";\nexport const twice: number = version * 2;\n",
	}

	newSession := func(trusted bool) (*project.Session, *projecttestutil.SessionUtils) {
		init, utils := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
			CurrentDirectory:               "/home/project",
			DefaultLibraryPath:             bundled.LibPath(),
			TypingsLocation:                projecttestutil.TestTypingsLocation,
			PositionEncoding:               lsproto.PositionEncodingKindUTF8,
			LoggingEnabled:                 true,
			DangerouslyLoadExternalPlugins: trusted,
		}, nil)
		init.Spawner = contentmappertest.NewSpawner()
		return project.NewSession(init), utils
	}

	t.Run("trusted workspace transforms the content-mapped file", func(t *testing.T) {
		t.Parallel()
		session, utils := newSession(true)
		defer session.Close()

		session.DidOpenFile(context.Background(), "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		ls, err := session.GetLanguageService(context.Background(), "file:///home/project/main.ts")
		assert.NilError(t, err)

		boxFile := ls.GetProgram().GetSourceFile("/home/project/app.box")
		assert.Assert(t, boxFile != nil, "expected app.box to be loaded into the program")
		// The #{target} token was substituted with the es2020 target value (7) by the content mapper.
		assert.Assert(t, strings.Contains(boxFile.Text(), "export const version = 7;"), "app.box was not transformed: %q", boxFile.Text())

		// The config's .box mapper should have been registered for text document synchronization.
		session.WaitForBackgroundTasks()
		calls := utils.Client().RegisterContentMapperExtensionsCalls()
		assert.Assert(t, len(calls) > 0, "expected RegisterContentMapperExtensions to be called")
		assert.DeepEqual(t, calls[len(calls)-1].Extensions, []string{".box"})
	})

	t.Run("untrusted workspace does not run the content mapper", func(t *testing.T) {
		t.Parallel()
		session, utils := newSession(false)
		defer session.Close()

		session.DidOpenFile(context.Background(), "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		ls, err := session.GetLanguageService(context.Background(), "file:///home/project/main.ts")
		assert.NilError(t, err)

		// Without workspace trust, the content mapper gate drops the mappers, so .box is not a recognized
		// extension and app.box never enters the program.
		boxFile := ls.GetProgram().GetSourceFile("/home/project/app.box")
		assert.Assert(t, boxFile == nil, "app.box should not be loaded without trust")

		// No content mapper extensions should be registered without trust.
		session.WaitForBackgroundTasks()
		for _, call := range utils.Client().RegisterContentMapperExtensionsCalls() {
			assert.Equal(t, len(call.Extensions), 0, "expected no content mapper extensions to be registered without trust")
		}
	})

	t.Run("editing an open content-mapped file reparses it through the mapper", func(t *testing.T) {
		t.Parallel()
		session, _ := newSession(true)
		defer session.Close()

		ctx := context.Background()
		session.DidOpenFile(ctx, "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		// Open the .box with a foreign language id so its overlay script kind is Unknown, matching how an
		// editor opens a content-mapped file. This is what made the incremental reparse panic.
		session.DidOpenFile(ctx, "file:///home/project/app.box", 1, files["/home/project/app.box"].(string), lsproto.LanguageKind("box"))
		_, err := session.GetLanguageService(ctx, "file:///home/project/main.ts")
		assert.NilError(t, err)

		// Editing the open .box file drives the single-file incremental reparse path
		// (Program.UpdateProgram), which must re-run the content mapper transform rather than parse the
		// raw foreign text.
		session.DidChangeFile(ctx, "file:///home/project/app.box", 2, []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: "export const version = #{target};\nexport const extra = 1;\n"}},
		})
		ls, err := session.GetLanguageService(ctx, "file:///home/project/main.ts")
		assert.NilError(t, err)

		boxFile := ls.GetProgram().GetSourceFile("/home/project/app.box")
		assert.Assert(t, boxFile != nil, "expected app.box to be loaded")
		assert.Assert(t, strings.Contains(boxFile.Text(), "export const version = 7;"), "reparsed app.box was not transformed: %q", boxFile.Text())
		assert.Assert(t, strings.Contains(boxFile.Text(), "export const extra = 1;"), "reparsed app.box missing the edit: %q", boxFile.Text())
	})

	t.Run("unchanged content-mapped file is reused from the cache across a full rebuild", func(t *testing.T) {
		t.Parallel()
		session, utils := newSession(true)
		defer session.Close()

		ctx := context.Background()
		session.DidOpenFile(ctx, "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		ls, err := session.GetLanguageService(ctx, "file:///home/project/main.ts")
		assert.NilError(t, err)
		boxFile := ls.GetProgram().GetSourceFile("/home/project/app.box")
		assert.Assert(t, boxFile != nil, "expected app.box to be loaded")
		assert.Assert(t, strings.Contains(boxFile.Text(), "export const version = 7;"), "app.box was not transformed: %q", boxFile.Text())

		// Changing a compiler option the mapper does not depend on (strict) forces a full program
		// rebuild while leaving app.box's content and the mapper's transform identity unchanged, so the
		// transformed file must be served from the parse cache rather than re-transformed.
		err = utils.FS().WriteFile("/home/project/tsconfig.json", `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": false },
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`)
		assert.NilError(t, err)
		session.DidChangeWatchedFiles(ctx, []*lsproto.FileEvent{{Uri: "file:///home/project/tsconfig.json", Type: lsproto.FileChangeTypeChanged}})

		ls, err = session.GetLanguageService(ctx, "file:///home/project/main.ts")
		assert.NilError(t, err)
		rebuiltBox := ls.GetProgram().GetSourceFile("/home/project/app.box")
		assert.Assert(t, rebuiltBox == boxFile, "expected the unchanged content-mapped file to be reused from the parse cache, not re-transformed")
	})
}

// TestContentMapperDiagnostics verifies that compiler diagnostics on a content-mapped file are reported
// against the original (untransformed) text: the mapper prepends a synthesized preamble line, so a
// diagnostic in the body would land one line too low if its transformed range were used verbatim.
func TestContentMapperDiagnostics(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}
	files := map[string]any{
		"/home/project/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true, "skipLibCheck": true },
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`,
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.ExecName),
		"/home/project/app.box":                          "export const version = #{target};\nexport const bad: string = version;\n",
		"/home/project/main.ts":                          "import \"./app.box\";\n",
	}

	init, _ := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
		CurrentDirectory:               "/home/project",
		DefaultLibraryPath:             bundled.LibPath(),
		TypingsLocation:                projecttestutil.TestTypingsLocation,
		PositionEncoding:               lsproto.PositionEncodingKindUTF8,
		DangerouslyLoadExternalPlugins: true,
	}, nil)
	init.Spawner = contentmappertest.NewSpawner()
	session := project.NewSession(init)
	defer session.Close()

	ctx := context.Background()
	session.DidOpenFile(ctx, "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
	// Open the .box with a foreign language id, matching how an editor opens a content-mapped file.
	session.DidOpenFile(ctx, "file:///home/project/app.box", 1, files["/home/project/app.box"].(string), lsproto.LanguageKind("box"))

	ls, err := session.GetLanguageService(ctx, "file:///home/project/app.box")
	assert.NilError(t, err)

	resp, err := ls.ProvideDiagnostics(ctx, "file:///home/project/app.box")
	assert.NilError(t, err)
	report := resp.FullDocumentDiagnosticReport
	assert.Assert(t, report != nil, "expected a full diagnostic report")

	// The assignability error (TS2322) is on `version` in `export const bad: string = version;`, which is
	// original line 1. The transformed text prepends a synthesized preamble line, so an unmapped range
	// would report it on line 2; the span map must place it back on line 1.
	var found bool
	for _, d := range report.Items {
		if d.Code != nil && d.Code.Integer != nil && *d.Code.Integer == 2322 {
			found = true
			assert.Equal(t, d.Range.Start.Line, uint32(1), "type error should map back to the original .box line")
		}
	}
	assert.Assert(t, found, "expected a type-assignability diagnostic on app.box")
}

// TestContentMapperSynthesizedDiagnostics verifies that compiler errors in a mapper's fully synthesized
// output (which have no location in the original file) are not dropped or scattered at position 0, but
// collected into a single aggregate diagnostic at the top of the file whose related information carries
// the real messages.
func TestContentMapperSynthesizedDiagnostics(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}
	files := map[string]any{
		"/home/project/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true, "skipLibCheck": true },
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`,
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.SynthesizingExec),
		"/home/project/app.box":                          "anything\n",
		"/home/project/main.ts":                          "import \"./app.box\";\n",
	}

	init, _ := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
		CurrentDirectory:               "/home/project",
		DefaultLibraryPath:             bundled.LibPath(),
		TypingsLocation:                projecttestutil.TestTypingsLocation,
		PositionEncoding:               lsproto.PositionEncodingKindUTF8,
		DangerouslyLoadExternalPlugins: true,
	}, nil)
	init.Spawner = contentmappertest.NewSpawner()
	session := project.NewSession(init)
	defer session.Close()

	ctx := context.Background()
	session.DidOpenFile(ctx, "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
	session.DidOpenFile(ctx, "file:///home/project/app.box", 1, files["/home/project/app.box"].(string), lsproto.LanguageKind("box"))

	ls, err := session.GetLanguageService(ctx, "file:///home/project/app.box")
	assert.NilError(t, err)

	// The aggregate carries the underlying messages as related information, which the converter only emits
	// when the client advertises support for it.
	caps := &lsproto.ResolvedClientCapabilities{}
	caps.TextDocument.Diagnostic.RelatedInformation = true
	diagCtx := lsproto.WithClientCapabilities(ctx, caps)
	resp, err := ls.ProvideDiagnostics(diagCtx, "file:///home/project/app.box")
	assert.NilError(t, err)
	report := resp.FullDocumentDiagnosticReport
	assert.Assert(t, report != nil, "expected a full diagnostic report")

	// The synthesizing mapper emits `export const el = jsxRuntime(Widget);` with a fully synthesized span
	// map, so its "Cannot find name" errors have no location in app.box. They must be folded into a single
	// aggregate at the top of the file, with the real messages preserved as related information.
	assert.Equal(t, len(report.Items), 1, "synthesized errors should collapse to a single aggregate diagnostic")
	aggregate := report.Items[0]
	assert.Equal(t, aggregate.Range.Start.Line, uint32(0), "aggregate should sit at the top of the file")
	assert.Assert(t, aggregate.Code != nil && aggregate.Code.Integer != nil && *aggregate.Code.Integer == 100037, "aggregate should carry the content-mapper diagnostic code")
	assert.Assert(t, aggregate.RelatedInformation != nil, "aggregate should carry related information")
	related := *aggregate.RelatedInformation
	assert.Assert(t, len(related) >= 1, "expected the real messages as related information")
	var sawCannotFindName bool
	for _, info := range related {
		if strings.Contains(info.Message, "Cannot find name") {
			sawCannotFindName = true
		}
	}
	assert.Assert(t, sawCannotFindName, "expected the underlying compiler messages to be preserved")
}

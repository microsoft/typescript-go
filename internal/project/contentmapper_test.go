package project_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls"
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

// TestContentMapperOutputMapping verifies that a language-service result whose location lands in a
// content-mapped file's transformed text is reported in the file's original coordinates. A go-to-definition
// from a .ts file into a .box declaration must point at the original position, not the transformed one
// (which is shifted by the mapper's synthesized preamble).
func TestContentMapperOutputMapping(t *testing.T) {
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
		"/home/project/app.box":                          "export const version = #{target};\n",
		"/home/project/main.ts":                          "import { version } from \"./app.box\";\nexport const twice: number = version * 2;\n",
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

	ls, err := session.GetLanguageService(ctx, "file:///home/project/main.ts")
	assert.NilError(t, err)

	// `version` appears at character 29 on line 1 of main.ts ("...number = version * 2;").
	definition, err := ls.ProvideDefinition(ctx, "file:///home/project/main.ts", lsproto.Position{Line: 1, Character: 32})
	assert.NilError(t, err)
	assert.Assert(t, definition.Locations != nil, "expected a definition result")
	locs := *definition.Locations
	assert.Equal(t, len(locs), 1, "expected exactly one definition location")

	// In app.box the declaration `export const version` places `version` at original line 0, characters
	// 13-20. The mapper's synthesized preamble shifts it to a later line in the transformed text, so an
	// unmapped range would report the wrong position.
	def := locs[0]
	assert.Equal(t, string(def.Uri), "file:///home/project/app.box")
	assert.Equal(t, def.Range.Start.Line, uint32(0), "definition should map back to the original .box line")
	assert.Equal(t, def.Range.Start.Character, uint32(13), "definition should map back to the original .box character")
	assert.Equal(t, def.Range.End.Character, uint32(20))
}

func TestContentMapperSynthesizedDefinitionFallsBackToFile(t *testing.T) {
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
		"/home/project/app.box":                          "component source with no direct TypeScript span\n",
		"/home/project/main.ts":                          "import { el } from \"./app.box\";\nexport const value = el;\n",
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
	ls, err := session.GetLanguageService(ctx, "file:///home/project/main.ts")
	assert.NilError(t, err)

	definition, err := ls.ProvideDefinition(ctx, "file:///home/project/main.ts", lsproto.Position{Line: 1, Character: 22})
	assert.NilError(t, err)
	assert.Assert(t, definition.Locations != nil, "expected a file-level definition result")
	locations := *definition.Locations
	assert.Equal(t, len(locations), 1)
	assert.Equal(t, string(locations[0].Uri), "file:///home/project/app.box")
	assert.Equal(t, locations[0].Range, lsproto.Range{}, "synthesized definition should navigate to the file without claiming a source span")

	symbols, err := ls.ProvideDocumentSymbols(ctx, "file:///home/project/app.box")
	assert.NilError(t, err)
	assert.Assert(t, symbols.SymbolInformations != nil)
	assert.Equal(t, len(*symbols.SymbolInformations), 0, "synthesized declarations should not appear as symbols at unrelated source ranges")
}

// TestContentMapperRequestOrigin verifies that a language-service request whose position originates inside
// a content-mapped file is answered by mapping the original position forward into the transformed text. A
// hover over a real identifier in the .box file must resolve, while a hover in a fully synthesized mapper's
// output (which has no original counterpart) must return nothing rather than a misplaced result.
func TestContentMapperRequestOrigin(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	newSession := func(t *testing.T, execName, boxContent string) *project.Session {
		files := map[string]any{
			"/home/project/tsconfig.json": `{
				"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true, "skipLibCheck": true },
				"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
			}`,
			"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(execName),
			"/home/project/app.box":                          boxContent,
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
		ctx := context.Background()
		session.DidOpenFile(ctx, "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		session.DidOpenFile(ctx, "file:///home/project/app.box", 1, boxContent, lsproto.LanguageKind("box"))
		return session
	}

	t.Run("hover over a verbatim identifier resolves through the mapped position", func(t *testing.T) {
		t.Parallel()
		session := newSession(t, contentmappertest.ExecName, "export const version = #{target};\n")
		defer session.Close()

		ls, err := session.GetLanguageService(context.Background(), "file:///home/project/app.box")
		assert.NilError(t, err)

		// `version` sits at original line 0, characters 13-20. Its original position must be mapped forward
		// into the transformed text for the checker to resolve it.
		hover, err := ls.ProvideHover(context.Background(), &lsproto.HoverParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: "file:///home/project/app.box"},
			Position:     lsproto.Position{Line: 0, Character: 15},
		})
		assert.NilError(t, err)
		assert.Assert(t, hover.Hover != nil, "expected a hover result for the .box identifier")
		assert.Assert(t, strings.Contains(hover.Hover.Contents.MarkupContent.Value, "version"), "hover should describe `version`: %+v", hover.Hover.Contents)
	})

	t.Run("hover in fully synthesized output returns nothing", func(t *testing.T) {
		t.Parallel()
		session := newSession(t, contentmappertest.SynthesizingExec, "anything\n")
		defer session.Close()

		ls, err := session.GetLanguageService(context.Background(), "file:///home/project/app.box")
		assert.NilError(t, err)

		// The synthesizing mapper's output has no original counterpart, so the position cannot be mapped
		// forward and the request must yield no result.
		hover, err := ls.ProvideHover(context.Background(), &lsproto.HoverParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: "file:///home/project/app.box"},
			Position:     lsproto.Position{Line: 0, Character: 0},
		})
		assert.NilError(t, err)
		assert.Assert(t, hover.Hover == nil, "expected no hover result in fully synthesized output")
	})
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

// TestContentMapperCompletions verifies that completions are offered inside verbatim spans of a
// content-mapped file (mapped to the original position) but declined outside them, and that auto-import
// completions whose import edit would land in generated code are filtered out of the list.
func TestContentMapperCompletions(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	newSession := func(t *testing.T, files map[string]any) *project.Session {
		init, _ := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
			CurrentDirectory:               "/home/project",
			DefaultLibraryPath:             bundled.LibPath(),
			TypingsLocation:                projecttestutil.TestTypingsLocation,
			PositionEncoding:               lsproto.PositionEncodingKindUTF8,
			DangerouslyLoadExternalPlugins: true,
		}, nil)
		init.Spawner = contentmappertest.NewSpawner()
		session := project.NewSession(init)
		ctx := context.Background()
		session.DidOpenFile(ctx, "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		session.DidOpenFile(ctx, "file:///home/project/app.box", 1, files["/home/project/app.box"].(string), lsproto.LanguageKind("box"))
		return session
	}

	baseFiles := func(boxContent string) map[string]any {
		return map[string]any{
			"/home/project/tsconfig.json": `{
				"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true, "skipLibCheck": true },
				"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
			}`,
			"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.ExecName),
			"/home/project/app.box":                          boxContent,
			"/home/project/main.ts":                          "import \"./app.box\";\n",
		}
	}

	t.Run("offers completions in a verbatim span", func(t *testing.T) {
		t.Parallel()
		session := newSession(t, baseFiles("export const foo = 1;\nfoo\n"))
		defer session.Close()

		// Complete at the end of `foo` on original line 1 (a verbatim span).
		resp := completeContentMapped(t, session, "file:///home/project/app.box", lsproto.Position{Line: 1, Character: 3})
		assert.Assert(t, resp.List != nil, "expected a completion list in a verbatim span")
		var sawFoo bool
		for _, item := range resp.List.Items {
			if item.Label == "foo" {
				sawFoo = true
			}
		}
		assert.Assert(t, sawFoo, "expected `foo` to be offered as a completion")

		// The edit range must land on the original cursor line (1). If it were left in transformed
		// coordinates (shifted by the synthesized preamble), it would not contain the cursor and VS Code
		// would silently filter out every item.
		assert.Assert(t, resp.List.ItemDefaults != nil && resp.List.ItemDefaults.EditRange != nil, "expected an edit range default")
		insert := resp.List.ItemDefaults.EditRange.EditRangeWithInsertReplace.Insert
		assert.Equal(t, insert.End.Line, uint32(1), "the edit range must map back to the original cursor line")
		assert.Equal(t, insert.End.Character, uint32(3), "the edit range must end at the original cursor character")
	})

	t.Run("declines completions outside a verbatim span", func(t *testing.T) {
		t.Parallel()
		session := newSession(t, baseFiles("export const version = #{target};\n"))
		defer session.Close()

		// The `#{target}` token maps to an atom span (a mapper substitution), which is not editable.
		resp := completeContentMapped(t, session, "file:///home/project/app.box", lsproto.Position{Line: 0, Character: 25})
		assert.Assert(t, resp.List == nil && resp.Items == nil, "expected no completions inside an atom span")
	})

	t.Run("filters auto-imports that would edit generated code", func(t *testing.T) {
		t.Parallel()
		files := baseFiles("export const version = #{target};\nconst x = help;\n")
		files["/home/project/dep.ts"] = "export const helper = 1;\n"
		session := newSession(t, files)
		defer session.Close()

		// `help` on original line 1 would suggest auto-importing `helper` from ./dep, but the new import
		// statement would be inserted into the synthesized preamble, so the item must be dropped.
		resp := completeContentMapped(t, session, "file:///home/project/app.box", lsproto.Position{Line: 1, Character: 14})
		assert.Assert(t, resp.List != nil, "expected a completion list")
		for _, item := range resp.List.Items {
			if item.Data != nil && item.Data.AutoImport != nil {
				t.Fatalf("auto-import completion %q should have been filtered out (edit lands in generated code)", item.Label)
			}
		}
	})

	t.Run("keeps auto-imports whose edit lands in a verbatim span", func(t *testing.T) {
		t.Parallel()
		files := baseFiles("import { a } from \"./dep\";\nconst x = help;\n")
		files["/home/project/dep.ts"] = "export const a = 1;\nexport const helper = 2;\n"
		session := newSession(t, files)
		defer session.Close()

		// `./dep` is already imported in the verbatim body, so completing `help` adds `helper` to the
		// existing import on original line 0 — an edit fully inside a verbatim span, which must be kept.
		resp := completeContentMapped(t, session, "file:///home/project/app.box", lsproto.Position{Line: 1, Character: 14})
		assert.Assert(t, resp.List != nil, "expected a completion list")
		var helper *lsproto.CompletionItem
		for _, item := range resp.List.Items {
			if item.Label == "helper" && item.Data != nil && item.Data.AutoImport != nil {
				helper = item
			}
		}
		assert.Assert(t, helper != nil, "expected the `helper` auto-import to be offered")
		assert.Assert(t, helper.AdditionalTextEdits != nil, "expected the import edit to be attached eagerly")
		for _, edit := range *helper.AdditionalTextEdits {
			assert.Equal(t, edit.Range.Start.Line, uint32(0), "the import edit should map to original line 0")
		}
	})
}

// completeContentMapped issues a completion request, mirroring the server's two-phase auto-import protocol
// (retrying with an auto-import-prepared language service when the first attempt reports ErrNeedsAutoImports).
// It advertises editRange/insertReplace support so the response carries edit ranges, as VS Code does.
func completeContentMapped(t *testing.T, session *project.Session, uri lsproto.DocumentUri, pos lsproto.Position) lsproto.CompletionResponse {
	t.Helper()
	caps := &lsproto.ResolvedClientCapabilities{}
	caps.TextDocument.Completion.CompletionList.ItemDefaults = []string{"editRange", "commitCharacters"}
	caps.TextDocument.Completion.CompletionItem.InsertReplaceSupport = true
	ctx := lsproto.WithClientCapabilities(context.Background(), caps)
	var resp lsproto.CompletionResponse
	_, err := session.WithLanguageServiceAndSnapshot(ctx, uri, func(languageService *ls.LanguageService, snapshot *project.Snapshot) (func() error, error) {
		r, e := languageService.ProvideCompletion(ctx, uri, pos, &lsproto.CompletionContext{})
		if errors.Is(e, ls.ErrNeedsAutoImports) {
			languageService, e = session.GetLanguageServiceWithAutoImports(ctx, snapshot, uri)
			if e != nil {
				return nil, e
			}
			r, e = languageService.ProvideCompletion(ctx, uri, pos, &lsproto.CompletionContext{})
		}
		resp = r
		return nil, e
	})
	assert.NilError(t, err)
	return resp
}

// TestContentMapperRename verifies that a rename originating in a verbatim span rewrites the verbatim
// occurrences in original coordinates, and that prepareRename declines a position outside a verbatim span.
func TestContentMapperRename(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	newSession := func(t *testing.T, boxContent string) *project.Session {
		files := map[string]any{
			"/home/project/tsconfig.json": `{
				"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true, "skipLibCheck": true },
				"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
			}`,
			"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.ExecName),
			"/home/project/app.box":                          boxContent,
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
		ctx := context.Background()
		session.DidOpenFile(ctx, "file:///home/project/main.ts", 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		session.DidOpenFile(ctx, "file:///home/project/app.box", 1, boxContent, lsproto.LanguageKind("box"))
		return session
	}

	t.Run("renames verbatim occurrences in original coordinates", func(t *testing.T) {
		t.Parallel()
		session := newSession(t, "export const foo = 1;\nexport const bar = foo;\n")
		defer session.Close()

		ctx := context.Background()
		ls, err := session.GetLanguageService(ctx, "file:///home/project/app.box")
		assert.NilError(t, err)

		// `foo` is declared at original line 0, char 13 (a verbatim span), so rename is allowed.
		info := ls.GetRenameInfo(ctx, "baz", "file:///home/project/app.box", lsproto.Position{Line: 0, Character: 14})
		assert.Assert(t, info.CanRename, "expected rename to be allowed at a verbatim identifier")

		result, err := ls.ProvideRename(ctx, &lsproto.RenameParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: "file:///home/project/app.box"},
			Position:     lsproto.Position{Line: 0, Character: 14},
			NewName:      "baz",
		}, nil)
		assert.NilError(t, err)
		assert.Assert(t, result.WorkspaceEdit != nil && result.WorkspaceEdit.Changes != nil, "expected a workspace edit")
		edits := (*result.WorkspaceEdit.Changes)["file:///home/project/app.box"]
		assert.Equal(t, len(edits), 2, "expected the declaration and use to be renamed")
		lines := []uint32{}
		for _, edit := range edits {
			lines = append(lines, edit.Range.Start.Line)
			assert.Equal(t, edit.NewText, "baz")
		}
		assert.Assert(t, slices.Contains(lines, uint32(0)) && slices.Contains(lines, uint32(1)), "edits should map to original lines 0 and 1, got %v", lines)
	})

	t.Run("declines rename outside a verbatim span", func(t *testing.T) {
		t.Parallel()
		session := newSession(t, "export const version = #{target};\n")
		defer session.Close()

		ctx := context.Background()
		ls, err := session.GetLanguageService(ctx, "file:///home/project/app.box")
		assert.NilError(t, err)

		// The `#{target}` token maps to an atom span, which cannot be edited in the original text.
		info := ls.GetRenameInfo(ctx, "whatever", "file:///home/project/app.box", lsproto.Position{Line: 0, Character: 25})
		assert.Assert(t, !info.CanRename, "expected rename to be declined inside an atom span")
	})
}

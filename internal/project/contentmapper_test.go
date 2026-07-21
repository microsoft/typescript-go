package project_test

import (
	"context"
	"errors"
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
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.TransformingMapper),
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

	t.Run("watch change to a content-mapped file updates the program", func(t *testing.T) {
		t.Parallel()
		session, utils := newSession(true)
		defer session.Close()

		ctx := context.Background()
		mainURI := lsproto.DocumentUri("file:///home/project/main.ts")
		session.DidOpenFile(ctx, mainURI, 1, files["/home/project/main.ts"].(string), lsproto.LanguageKindTypeScript)
		languageService, err := session.GetLanguageService(ctx, mainURI)
		assert.NilError(t, err)
		original := languageService.GetProgram().GetSourceFile("/home/project/app.box")
		assert.Assert(t, original != nil, "expected app.box to be loaded")

		// Wait until the configured extension set has been published; watch filtering uses the set captured
		// when the snapshot change is created.
		session.WaitForBackgroundTasks()
		updatedContent := "export const version = #{target};\nexport const watched = true;\n"
		assert.NilError(t, utils.FS().WriteFile("/home/project/app.box", updatedContent))
		session.DidChangeWatchedFiles(ctx, []*lsproto.FileEvent{{
			Uri:  "file:///home/project/app.box",
			Type: lsproto.FileChangeTypeChanged,
		}})

		languageService, err = session.GetLanguageService(ctx, mainURI)
		assert.NilError(t, err)
		updatedSnapshot := session.Snapshot()
		configuredProject := updatedSnapshot.GetDefaultProject(mainURI)
		assert.Assert(t, configuredProject != nil, "expected configured project")
		assert.Equal(t, configuredProject.ProgramUpdateKind, project.ProgramUpdateKindCloned)
		assert.Equal(t, configuredProject.ProgramLastUpdate, updatedSnapshot.ID())
		updated := languageService.GetProgram().GetSourceFile("/home/project/app.box")
		assert.Assert(t, updated != nil, "expected app.box to remain loaded")
		assert.Assert(t, updated != original, "expected the watched content-mapped file to be reparsed")
		assert.Assert(t, strings.Contains(updated.Text(), "export const version = 7;"), "updated app.box was not transformed: %q", updated.Text())
		assert.Assert(t, strings.Contains(updated.Text(), "export const watched = true;"), "updated app.box missing watched change: %q", updated.Text())
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

func TestContentMapperOpenFileExcludedByConfigChange(t *testing.T) {
	t.Parallel()
	files := map[string]any{
		"/home/project/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true },
			"include": ["src"],
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`,
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.TransformingMapper),
		"/home/project/src/app.box":                      "export const version = #{target};\n",
		"/home/project/src/main.ts":                      "export const main = true;\n",
	}
	init, utils := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
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
	boxURI := lsproto.DocumentUri("file:///home/project/src/app.box")
	session.DidOpenFile(ctx, boxURI, 1, files["/home/project/src/app.box"].(string), lsproto.LanguageKind("box"))
	languageService, err := session.GetLanguageService(ctx, boxURI)
	assert.NilError(t, err)
	assert.Assert(t, languageService.GetProgram().GetSourceFile("/home/project/src/app.box") != nil)

	assert.NilError(t, utils.FS().WriteFile("/home/project/tsconfig.json", `{
		"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true },
		"include": ["src/**/*.ts"],
		"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
	}`))
	session.DidChangeWatchedFiles(ctx, []*lsproto.FileEvent{{
		Uri:  "file:///home/project/tsconfig.json",
		Type: lsproto.FileChangeTypeChanged,
	}})

	languageService, err = session.GetLanguageService(ctx, boxURI)
	assert.NilError(t, err)
	defaultProject := session.Snapshot().GetDefaultProject(boxURI)
	assert.Assert(t, defaultProject != nil, "expected a default project for the open app.box")
	assert.Equal(t, defaultProject.Kind, project.KindInferred)
	boxFile := languageService.GetProgram().GetSourceFile("/home/project/src/app.box")
	assert.Assert(t, boxFile != nil, "expected the open app.box in the inferred project")
	assert.Assert(t, boxFile.ContentMapper() != "", "expected app.box to retain its content mapper")
	assert.Assert(t, !strings.Contains(boxFile.Text(), "#{target}"), "expected app.box to be transformed: %q", boxFile.Text())
}

func TestContentMapperInferredProjectUsesSessionMappers(t *testing.T) {
	t.Parallel()
	files := map[string]any{
		"/home/configured/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler" },
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`,
		"/home/configured/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.TransformingMapper),
		"/home/configured/main.ts":                          "export const main = true;\n",
		"/home/loose/app.box":                               "export const version = #{target};\n",
	}
	init, _ := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
		CurrentDirectory:               "/home",
		DefaultLibraryPath:             bundled.LibPath(),
		TypingsLocation:                projecttestutil.TestTypingsLocation,
		PositionEncoding:               lsproto.PositionEncodingKindUTF8,
		DangerouslyLoadExternalPlugins: true,
	}, nil)
	init.Spawner = contentmappertest.NewSpawner()
	session := project.NewSession(init)
	defer session.Close()

	ctx := context.Background()
	configuredURI := lsproto.DocumentUri("file:///home/configured/main.ts")
	session.DidOpenFile(ctx, configuredURI, 1, files["/home/configured/main.ts"].(string), lsproto.LanguageKindTypeScript)
	_, err := session.GetLanguageService(ctx, configuredURI)
	assert.NilError(t, err)

	boxURI := lsproto.DocumentUri("file:///home/loose/app.box")
	session.DidOpenFile(ctx, boxURI, 1, files["/home/loose/app.box"].(string), lsproto.LanguageKind("box"))
	languageService, err := session.GetLanguageService(ctx, boxURI)
	assert.NilError(t, err)
	defaultProject := session.Snapshot().GetDefaultProject(boxURI)
	assert.Assert(t, defaultProject != nil, "expected a default project for the loose app.box")
	assert.Equal(t, defaultProject.Kind, project.KindInferred)
	boxFile := languageService.GetProgram().GetSourceFile("/home/loose/app.box")
	assert.Assert(t, boxFile != nil, "expected loose app.box in the inferred project")
	assert.Assert(t, boxFile.ContentMapper() != "", "expected loose app.box to use the session mapper union")
	assert.Assert(t, !strings.Contains(boxFile.Text(), "#{target}"), "expected loose app.box to be transformed: %q", boxFile.Text())
}

func TestDiscoverContentMapperExtensions(t *testing.T) {
	t.Parallel()
	files := map[string]any{
		"/home/project/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true },
			"contentMappers": [ { "package": "mapper", "extensions": [".vue"] } ]
		}`,
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.ComponentMapper),
		"/home/project/ProfileCard.vue":                  `<template><h1>Profile</h1></template>`,
		"/home/other/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler" },
			"contentMappers": [ { "package": "mapper", "extensions": [".svelte"] } ]
		}`,
		"/home/other/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.ComponentMapper),
		"/home/other/Widget.svelte":                    `<template><h1>Widget</h1></template>`,
	}
	init, utils := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
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
	assert.DeepEqual(t, session.DiscoverContentMapperExtensions(ctx, []lsproto.DocumentUri{
		"file:///home/other/Widget.svelte",
	}, []string{".svelte"}), []string{".svelte"})

	matched := session.DiscoverContentMapperExtensions(ctx, []lsproto.DocumentUri{
		"file:///home/project/ProfileCard.vue",
		"file:///home/project/ignored.svelte",
	}, []string{".vue", ".svelte", ".vue", "vue", "../bad"})
	assert.DeepEqual(t, matched, []string{".vue"})

	calls := utils.Client().RegisterContentMapperExtensionsCalls()
	assert.Assert(t, len(calls) > 0, "expected content mapper extensions to be registered")
	assert.DeepEqual(t, calls[len(calls)-1].Extensions, []string{".svelte", ".vue"})

	snapshot := session.Snapshot()
	assert.Assert(t, snapshot.GetDefaultProject("file:///home/project/ProfileCard.vue") != nil, "expected the configured project to be discovered")
	assert.Assert(t, snapshot.GetDefaultProject("file:///home/project/ignored.svelte") == nil, "unmatched extensions should not load a project")

	untrustedInit, untrustedUtils := projecttestutil.GetSessionInitOptions(files, &project.SessionOptions{
		CurrentDirectory:               "/home/project",
		DefaultLibraryPath:             bundled.LibPath(),
		TypingsLocation:                projecttestutil.TestTypingsLocation,
		PositionEncoding:               lsproto.PositionEncodingKindUTF8,
		DangerouslyLoadExternalPlugins: false,
	}, nil)
	untrustedInit.Spawner = contentmappertest.NewSpawner()
	untrusted := project.NewSession(untrustedInit)
	defer untrusted.Close()
	assert.Equal(t, len(untrusted.DiscoverContentMapperExtensions(ctx, []lsproto.DocumentUri{"file:///home/project/ProfileCard.vue"}, []string{".vue"})), 0)
	assert.Equal(t, len(untrustedUtils.Client().RegisterContentMapperExtensionsCalls()), 0)
}

func TestContentMapperSynthesizedDocumentSymbols(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}
	files := map[string]any{
		"/home/project/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true, "skipLibCheck": true },
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`,
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.SynthesizingMapper),
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

	symbols, err := ls.ProvideDocumentSymbols(ctx, "file:///home/project/app.box")
	assert.NilError(t, err)
	assert.Assert(t, symbols.SymbolInformations != nil)
	assert.Equal(t, len(*symbols.SymbolInformations), 0, "synthesized declarations should not appear as symbols at unrelated source ranges")
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
			"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.TransformingMapper),
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

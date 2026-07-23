package project_test

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func TestWorkspaceRootSymlinkProjectReferenceRepro(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/home": vfstest.Symlink("/real/home"),

		"/real/home/workspace/tsconfig.json": `{
			"files": [],
			"references": [
				{ "path": "./app" },
				{ "path": "./pkg" }
			]
		}`,
		"/real/home/workspace/app/tsconfig.json": `{
			"compilerOptions": {
				"composite": true,
				"module": "nodenext",
				"moduleResolution": "nodenext",
				"rootDir": "./src",
				"outDir": "./dist"
			},
			"files": ["src/index.ts"],
			"references": [
				{ "path": "../pkg" }
			]
		}`,
		"/real/home/workspace/app/src/index.ts": `import { value } from "mylib";
export const answer = value;`,
		"/real/home/workspace/app/node_modules/mylib": vfstest.Symlink("/real/home/workspace/pkg"),

		"/real/home/workspace/pkg/tsconfig.json": `{
			"compilerOptions": {
				"composite": true,
				"module": "nodenext",
				"moduleResolution": "nodenext",
				"rootDir": "./src",
				"outDir": "./dist",
				"declaration": true
			},
			"files": ["src/index.ts"]
		}`,
		"/real/home/workspace/pkg/package.json": `{
			"name": "mylib",
			"main": "./dist/index.js",
			"types": "./dist/index.d.ts"
		}`,
		"/real/home/workspace/pkg/src/index.ts": `export const value = 42;`,
	}

	t.Run("realpath open works", func(t *testing.T) {
		t.Parallel()

		session, _ := projecttestutil.SetupWithOptions(files, &project.SessionOptions{
			CurrentDirectory:   "/real/home/workspace/app",
			DefaultLibraryPath: bundled.LibPath(),
			TypingsLocation:    projecttestutil.TestTypingsLocation,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       true,
			LoggingEnabled:     true,
		})

		uri := lsproto.DocumentUri("file:///real/home/workspace/app/src/index.ts")
		session.DidOpenFile(context.Background(), uri, 1, files["/real/home/workspace/app/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		defaultProject, _, allProjects, err := session.GetLanguageServiceAndProjectsForFile(context.Background(), uri)
		assert.NilError(t, err)
		assert.Equal(t, defaultProject.Name(), "/real/home/workspace/app/tsconfig.json")
		assert.Equal(t, len(allProjects), 1)

		ls, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)
		program := ls.GetProgram()
		sourceFile := program.GetSourceFile("/real/home/workspace/app/src/index.ts")
		assert.Assert(t, sourceFile != nil)
		resolved := program.ResolveModuleName("mylib", sourceFile.FileName(), core.ResolutionModeCommonJS)
		assert.Assert(t, resolved != nil)
		assert.Equal(t, resolved.OriginalPath, "/real/home/workspace/app/node_modules/mylib/dist/index.d.ts")
		assert.Equal(t, resolved.ResolvedFileName, "/real/home/workspace/pkg/dist/index.d.ts")
		diags := program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), sourceFile)
		assert.Equal(t, len(diags), 0)
	})

	t.Run("symlink open works", func(t *testing.T) {
		t.Parallel()

		session, _ := projecttestutil.SetupWithOptions(files, &project.SessionOptions{
			CurrentDirectory:   "/home/workspace/app",
			DefaultLibraryPath: bundled.LibPath(),
			TypingsLocation:    projecttestutil.TestTypingsLocation,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       true,
			LoggingEnabled:     true,
		})

		uri := lsproto.DocumentUri("file:///home/workspace/app/src/index.ts")
		session.DidOpenFile(context.Background(), uri, 1, files["/real/home/workspace/app/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

		defaultProject, _, allProjects, err := session.GetLanguageServiceAndProjectsForFile(context.Background(), uri)
		assert.NilError(t, err)
		assert.Equal(t, defaultProject.Name(), "/home/workspace/app/tsconfig.json")
		assert.Equal(t, len(allProjects), 1)

		ls, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)
		program := ls.GetProgram()
		sourceFile := program.GetSourceFile("/home/workspace/app/src/index.ts")
		assert.Assert(t, sourceFile != nil)
		resolved := program.ResolveModuleName("mylib", sourceFile.FileName(), core.ResolutionModeCommonJS)
		assert.Assert(t, resolved != nil)
		assert.Equal(t, resolved.OriginalPath, "/home/workspace/app/node_modules/mylib/dist/index.d.ts")
		assert.Equal(t, resolved.ResolvedFileName, "/real/home/workspace/pkg/dist/index.d.ts")
		diags := program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), sourceFile)
		assert.Equal(t, len(diags), 0)
	})
}

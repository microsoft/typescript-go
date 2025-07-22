package projectv2_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/testutil/projectv2testutil"
	"gotest.tools/v3/assert"
)

func TestWatch(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	t.Run("handling changes", func(t *testing.T) {
		defaultFiles := map[string]any{
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
		}

		t.Run("change open file", func(t *testing.T) {
			t.Parallel()
			session, utils := projectv2testutil.Setup(defaultFiles)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/x.ts", 1, defaultFiles["/home/projects/TS/p1/src/x.ts"].(string), lsproto.LanguageKindTypeScript)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, defaultFiles["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			lsBefore, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			programBefore := lsBefore.GetProgram()

			err = utils.FS().WriteFile("/home/projects/TS/p1/src/x.ts", `export const x = 2;`, false)
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeChanged,
					Uri:  "file:///home/projects/TS/p1/src/x.ts",
				},
			})

			lsAfter, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			// Program should remain the same since the file is open and changes are handled through DidChangeTextDocument
			assert.Equal(t, programBefore, lsAfter.GetProgram())
		})

		t.Run("change closed program file", func(t *testing.T) {
			t.Parallel()
			session, utils := projectv2testutil.Setup(defaultFiles)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, defaultFiles["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			lsBefore, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			programBefore := lsBefore.GetProgram()

			err = utils.FS().WriteFile("/home/projects/TS/p1/src/x.ts", `export const x = 2;`, false)
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeChanged,
					Uri:  "file:///home/projects/TS/p1/src/x.ts",
				},
			})

			lsAfter, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			assert.Check(t, lsAfter.GetProgram() != programBefore)
		})

		t.Run("change config file", func(t *testing.T) {
			t.Parallel()
			files := map[string]any{
				"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true,
					"strict": false
				}
			}`,
				"/home/projects/TS/p1/src/x.ts": `export declare const x: number | undefined;`,
				"/home/projects/TS/p1/src/index.ts": `
				import { x } from "./x";
				let y: number = x;`,
			}

			session, utils := projectv2testutil.Setup(files)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program := ls.GetProgram()
			// Should have 0 errors with strict: false
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 0)

			err = utils.FS().WriteFile("/home/projects/TS/p1/tsconfig.json", `{
			"compilerOptions": {
				"noLib": false,
				"strict": true
			}
		}`, false)
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeChanged,
					Uri:  "file:///home/projects/TS/p1/tsconfig.json",
				},
			})

			ls, err = session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program = ls.GetProgram()
			// Should have 1 error with strict: true
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 1)
		})

		t.Run("delete explicitly included file", func(t *testing.T) {
			t.Parallel()
			files := map[string]any{
				"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"files": ["src/index.ts", "src/x.ts"]
			}`,
				"/home/projects/TS/p1/src/x.ts":     `export declare const x: number | undefined;`,
				"/home/projects/TS/p1/src/index.ts": `import { x } from "./x";`,
			}
			session, utils := projectv2testutil.Setup(files)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program := ls.GetProgram()
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 0)

			err = utils.FS().Remove("/home/projects/TS/p1/src/x.ts")
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeDeleted,
					Uri:  "file:///home/projects/TS/p1/src/x.ts",
				},
			})

			ls, err = session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program = ls.GetProgram()
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 1)
			assert.Check(t, program.GetSourceFile("/home/projects/TS/p1/src/x.ts") == nil)
		})

		t.Run("delete wildcard included file", func(t *testing.T) {
			t.Parallel()
			files := map[string]any{
				"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"include": ["src"]
			}`,
				"/home/projects/TS/p1/src/index.ts": `let x = 2;`,
				"/home/projects/TS/p1/src/x.ts":     `let y = x;`,
			}
			session, utils := projectv2testutil.Setup(files)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/x.ts", 1, files["/home/projects/TS/p1/src/x.ts"].(string), lsproto.LanguageKindTypeScript)

			ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/x.ts")
			assert.NilError(t, err)
			program := ls.GetProgram()
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/x.ts"))), 0)

			err = utils.FS().Remove("/home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeDeleted,
					Uri:  "file:///home/projects/TS/p1/src/index.ts",
				},
			})

			ls, err = session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/x.ts")
			assert.NilError(t, err)
			program = ls.GetProgram()
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/x.ts"))), 1)
		})

		t.Run("create explicitly included file", func(t *testing.T) {
			t.Parallel()
			files := map[string]any{
				"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"files": ["src/index.ts", "src/y.ts"]
			}`,
				"/home/projects/TS/p1/src/index.ts": `import { y } from "./y";`,
			}
			session, utils := projectv2testutil.Setup(files)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program := ls.GetProgram()

			// Initially should have an error because y.ts is missing
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 1)

			// Add the missing file
			err = utils.FS().WriteFile("/home/projects/TS/p1/src/y.ts", `export const y = 1;`, false)
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeCreated,
					Uri:  "file:///home/projects/TS/p1/src/y.ts",
				},
			})

			// Error should be resolved
			ls, err = session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program = ls.GetProgram()
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 0)
			assert.Check(t, program.GetSourceFile("/home/projects/TS/p1/src/y.ts") != nil)
		})

		t.Run("create failed lookup location", func(t *testing.T) {
			t.Parallel()
			files := map[string]any{
				"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"files": ["src/index.ts"]
			}`,
				"/home/projects/TS/p1/src/index.ts": `import { z } from "./z";`,
			}
			session, utils := projectv2testutil.Setup(files)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program := ls.GetProgram()

			// Initially should have an error because z.ts is missing
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 1)

			// Add a new file through failed lookup watch
			err = utils.FS().WriteFile("/home/projects/TS/p1/src/z.ts", `export const z = 1;`, false)
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeCreated,
					Uri:  "file:///home/projects/TS/p1/src/z.ts",
				},
			})

			// Error should be resolved and the new file should be included in the program
			ls, err = session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program = ls.GetProgram()
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 0)
			assert.Check(t, program.GetSourceFile("/home/projects/TS/p1/src/z.ts") != nil)
		})

		t.Run("create wildcard included file", func(t *testing.T) {
			t.Parallel()
			files := map[string]any{
				"/home/projects/TS/p1/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true
				},
				"include": ["src"]
			}`,
				"/home/projects/TS/p1/src/index.ts": `a;`,
			}
			session, utils := projectv2testutil.Setup(files)
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program := ls.GetProgram()

			// Initially should have an error because declaration for 'a' is missing
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 1)

			// Add a new file through wildcard watch
			err = utils.FS().WriteFile("/home/projects/TS/p1/src/a.ts", `const a = 1;`, false)
			assert.NilError(t, err)

			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Type: lsproto.FileChangeTypeCreated,
					Uri:  "file:///home/projects/TS/p1/src/a.ts",
				},
			})

			// Error should be resolved and the new file should be included in the program
			ls, err = session.GetLanguageService(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			assert.NilError(t, err)
			program = ls.GetProgram()
			assert.Equal(t, len(program.GetSemanticDiagnostics(projecttestutil.WithRequestID(t.Context()), program.GetSourceFile("/home/projects/TS/p1/src/index.ts"))), 0)
			assert.Check(t, program.GetSourceFile("/home/projects/TS/p1/src/a.ts") != nil)
		})
	})
}

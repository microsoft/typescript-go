package lspservertests

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func TestRename(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	testCases := []*lspServerTest{
		{
			subScenario: "ancestor and project ref management",
			files: func() map[string]any {
				return map[string]any{
					"/user/username/projects/temp/temp.ts":       "let x = 10",
					"/user/username/projects/temp/tsconfig.json": "{}",
					"/user/username/projects/container/lib/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"composite": true,
								},
								references: [],
								files: [
									"index.ts",
								],
							}`),
					"/user/username/projects/container/lib/index.ts": stringtestutil.Dedent(`
							export const myConst = 30;`),
					"/user/username/projects/container/exec/tsconfig.json": stringtestutil.Dedent(`
							{
								"files": ["./index.ts"],
								"references": [
									{ "path": "../lib" },
								],
							}`),
					"/user/username/projects/container/exec/index.ts": stringtestutil.Dedent(`
							import { myConst } from "../lib";
							export function getMyConst() {
								return myConst;
							}`),
					"/user/username/projects/container/compositeExec/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"composite": true,
								},
								"files": ["./index.ts"],
								"references": [
									{ "path": "../lib" },
								],
							}`),
					"/user/username/projects/container/compositeExec/index.ts": stringtestutil.Dedent(`
							import { myConst } from "../lib";
							export function getMyConst() {
								return myConst;
							}`),
					"/user/username/projects/container/tsconfig.json": stringtestutil.Dedent(`
							{
								"files": [],
								"include": [],
								"references": [
									{ "path": "./exec" },
									{ "path": "./compositeExec" },
								],
							}`),
				}
			},
			test: func(server *testServer) {
				file := "/user/username/projects/container/compositeExec/index.ts"
				temp := "/user/username/projects/temp/temp.ts"
				server.openFile(file, lsproto.LanguageKindTypeScript)
				// Open temp file and verify all projects alive
				server.openFile(temp, lsproto.LanguageKindTypeScript)

				// Ref projects are loaded after as part of this command
				server.baselineRename(file, lsptestutil.PositionToLineAndCharacter(file, server.content(file), "myConst", 0))

				// Open temp file and verify all projects alive
				server.closeFile(temp)
				server.openFile(temp, lsproto.LanguageKindTypeScript)

				// Close all files and open temp file, only inferred project should be alive
				server.closeFile(file)
				server.closeFile(temp)
				server.openFile(temp, lsproto.LanguageKindTypeScript)
			},
		},
		{
			subScenario: "rename in common file renames all project",
			files: func() map[string]any {
				return map[string]any{
					"/users/username/projects/a/a.ts": stringtestutil.Dedent(`
						import {C} from "./c/fc";
						console.log(C)
					`),
					"/users/username/projects/a/tsconfig.json": "{}",
					"/users/username/projects/a/c":             vfstest.Symlink("/users/username/projects/c"),
					"/users/username/projects/b/b.ts": stringtestutil.Dedent(`
						import {C} from "../c/fc";
						console.log(C)
					`),
					"/users/username/projects/b/tsconfig.json": "{}",
					"/users/username/projects/b/c":             vfstest.Symlink("/users/username/projects/c"),
					"/users/username/projects/c/fc.ts": stringtestutil.Dedent(`
						export const C = 42;
					`),
				}
			},
			test: func(server *testServer) {
				aTs := "/users/username/projects/a/a.ts"
				bTs := "/users/username/projects/b/b.ts"
				aFc := "/users/username/projects/a/c/fc.ts"
				bFc := "/users/username/projects/b/c/fc.ts"
				server.openFile(aTs, lsproto.LanguageKindTypeScript)
				server.openFile(bTs, lsproto.LanguageKindTypeScript)
				cContent := server.content("/users/username/projects/c/fc.ts")
				server.openFileWithContent(aFc, cContent, lsproto.LanguageKindTypeScript)
				server.openFileWithContent(bFc, cContent, lsproto.LanguageKindTypeScript)
				server.baselineRename(aFc, lsptestutil.PositionToLineAndCharacter(aFc, cContent, "C", 0))
			},
		},
	}

	for _, test := range testCases {
		test.run(t, "rename")
	}
}

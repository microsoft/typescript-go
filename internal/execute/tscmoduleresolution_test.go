package execute_test

import (
	"fmt"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func TestTscModuleResolution(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		getBuildModuleResolutionInProjectRefTestCase(false),
		getBuildModuleResolutionInProjectRefTestCase(true),
		{
			subScenario: `type reference resolution uses correct options for different resolution options referenced project`,
			files: FileMap{
				"/home/src/workspaces/project/packages/pkg1_index.ts": `export const theNum: TheNum = "type1";`,
				"/home/src/workspaces/project/packages/pkg1.tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": {
                            "composite": true,
                            "typeRoots": ["./typeroot1"]
                        },
                        "files": ["./pkg1_index.ts"],
                    }
                `),
				"/home/src/workspaces/project/packages/typeroot1/sometype/index.d.ts": `declare type TheNum = "type1";`,
				"/home/src/workspaces/project/packages/pkg2_index.ts":                 `export const theNum: TheNum2 = "type2";`,
				"/home/src/workspaces/project/packages/pkg2.tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": {
                            "composite": true,
                            "typeRoots": ["./typeroot2"]
                        },
                        "files": ["./pkg2_index.ts"],
                    }
                `),
				"/home/src/workspaces/project/packages/typeroot2/sometype/index.d.ts": `declare type TheNum2 = "type2";`,
			},
			commandLineArgs: []string{"-b", "packages/pkg1.tsconfig.json", "packages/pkg2.tsconfig.json", "--verbose", "--traceResolution"},
		},
		{
			subScenario: "impliedNodeFormat differs between projects for shared file",
			files: FileMap{
				"/home/src/workspaces/project/a/src/index.ts": "",
				"/home/src/workspaces/project/a/tsconfig.json": stringtestutil.Dedent(`
				{
                    "compilerOptions": {
						"strict": true
					}
				}
                `),
				"/home/src/workspaces/project/b/src/index.ts": stringtestutil.Dedent(`
                    import pg from "pg";
                    pg.foo();
                `),
				"/home/src/workspaces/project/b/tsconfig.json": stringtestutil.Dedent(`
				{
                    "compilerOptions": { 
						"strict": true,
						"module": "node16"
					},
                }`),
				"/home/src/workspaces/project/b/package.json": stringtestutil.Dedent(`
				{
                    "name": "b",
                    "type": "module"
                }`),
				"/home/src/workspaces/project/node_modules/@types/pg/index.d.ts": "export function foo(): void;",
				"/home/src/workspaces/project/node_modules/@types/pg/package.json": stringtestutil.Dedent(`
				{
                    "name": "@types/pg",
                    "types": "index.d.ts"
                }`),
			},
			commandLineArgs: []string{"-b", "a", "b", "--verbose", "--traceResolution", "--explainFiles"},
			edits:           noChangeOnlyEdit,
		},
		{
			subScenario:     "shared resolution should not report error",
			files:           getTscModuleResolutionSharingFileMap(),
			commandLineArgs: []string{"-b", "packages/b", "--verbose", "--traceResolution", "--explainFiles"},
		},
		{
			subScenario:     "when resolution is not shared",
			files:           getTscModuleResolutionSharingFileMap(),
			commandLineArgs: []string{"-b", "packages/a", "--verbose", "--traceResolution", "--explainFiles"},
			edits: []*tscEdit{
				{
					caption:         "build b",
					commandLineArgs: []string{"-b", "packages/b", "--verbose", "--traceResolution", "--explainFiles"},
				},
			},
		},
		{
			subScenario: "pnpm style layout",
			files: FileMap{
				// button@0.0.1
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.1/node_modules/@component-type-checker/button/src/index.ts": stringtestutil.Dedent(`
                    export interface Button {
                        a: number;
                        b: number;
                    }
                    export function createButton(): Button {
                        return {
                            a: 0,
                            b: 1,
                        };
                    }
                `),
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.1/node_modules/@component-type-checker/button/package.json": stringtestutil.Dedent(`
					{
						"name": "@component-type-checker/button",
						"version": "0.0.1",
						"main": "./src/index.ts"
					}`),

				// button@0.0.2
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.2/node_modules/@component-type-checker/button/src/index.ts": stringtestutil.Dedent(`
                    export interface Button {
                        a: number;
                        c: number;
                    }
                    export function createButton(): Button {
                        return {
                            a: 0,
                            c: 2,
                        };
                    }
                `),
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.2/node_modules/@component-type-checker/button/package.json": stringtestutil.Dedent(`
                    {
                        "name": "@component-type-checker/button",
                        "version": "0.0.2",
                        "main": "./src/index.ts"
                    }`),

				// @component-type-checker+components@0.0.1_@component-type-checker+button@0.0.1
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.1/node_modules/@component-type-checker/button": vfstest.Symlink(
					"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.1/node_modules/@component-type-checker/button",
				),
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.1/node_modules/@component-type-checker/components/src/index.ts": stringtestutil.Dedent(`
                    export { createButton, Button } from "@component-type-checker/button";
                `),
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.1/node_modules/@component-type-checker/components/package.json": stringtestutil.Dedent(`
					{
						"name": "@component-type-checker/components",
						"version": "0.0.1",
						"main": "./src/index.ts",
						"peerDependencies": {
							"@component-type-checker/button": "*"
						},
						"devDependencies": {
							"@component-type-checker/button": "0.0.2"
						}
					}`),

				// @component-type-checker+components@0.0.1_@component-type-checker+button@0.0.2
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.2/node_modules/@component-type-checker/button": vfstest.Symlink(
					"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.2/node_modules/@component-type-checker/button",
				),
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.2/node_modules/@component-type-checker/components/src/index.ts": stringtestutil.Dedent(`
                    export { createButton, Button } from "@component-type-checker/button";
                `),
				"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.2/node_modules/@component-type-checker/components/package.json": stringtestutil.Dedent(`
					{
						"name": "@component-type-checker/components",
						"version": "0.0.1",
						"main": "./src/index.ts",
						"peerDependencies": {
							"@component-type-checker/button": "*"
						},
						"devDependencies": {
							"@component-type-checker/button": "0.0.2"
						}
					}`),

				// sdk => @component-type-checker+components@0.0.1_@component-type-checker+button@0.0.1
				"/home/src/projects/component-type-checker/packages/sdk/src/index.ts": stringtestutil.Dedent(`
                    export { Button, createButton } from "@component-type-checker/components";
                    export const VERSION = "0.0.2";
                `),
				"/home/src/projects/component-type-checker/packages/sdk/package.json": stringtestutil.Dedent(`
                    {
                        "name": "@component-type-checker/sdk1",
                        "version": "0.0.2",
                        "main": "./src/index.ts",
                        "dependencies": {
                            "@component-type-checker/components": "0.0.1",
                            "@component-type-checker/button": "0.0.1"
                        }
                    }`),
				"/home/src/projects/component-type-checker/packages/sdk/node_modules/@component-type-checker/button": vfstest.Symlink(
					"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.1/node_modules/@component-type-checker/button",
				),
				"/home/src/projects/component-type-checker/packages/sdk/node_modules/@component-type-checker/components": vfstest.Symlink(
					"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.1/node_modules/@component-type-checker/components",
				),

				// app => @component-type-checker+components@0.0.1_@component-type-checker+button@0.0.2
				"/home/src/projects/component-type-checker/packages/app/src/app.tsx": stringtestutil.Dedent(`
                    import { VERSION } from "@component-type-checker/sdk";
                    import { Button } from "@component-type-checker/components";
                    import { createButton } from "@component-type-checker/button";
                    const button: Button = createButton();
                `),
				"/home/src/projects/component-type-checker/packages/app/package.json": stringtestutil.Dedent(`
					{
						"name": "app",
						"version": "1.0.0",
						"dependencies": {
							"@component-type-checker/button": "0.0.2",
							"@component-type-checker/components": "0.0.1",
							"@component-type-checker/sdk": "0.0.2"
						}
					}`),
				"/home/src/projects/component-type-checker/packages/app/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"target": "es5",
							"module": "esnext",
							"lib": ["ES5"],
							"moduleResolution": "node",
							"outDir": "dist",
						},
						"include": ["src"],
					}`),
				"/home/src/projects/component-type-checker/packages/app/node_modules/@component-type-checker/button": vfstest.Symlink(
					"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+button@0.0.2/node_modules/@component-type-checker/button",
				),
				"/home/src/projects/component-type-checker/packages/app/node_modules/@component-type-checker/components": vfstest.Symlink(
					"/home/src/projects/component-type-checker/node_modules/.pnpm/@component-type-checker+components@0.0.1_@component-type-checker+button@0.0.2/node_modules/@component-type-checker/components",
				),
				"/home/src/projects/component-type-checker/packages/app/node_modules/@component-type-checker/sdk": vfstest.Symlink(
					"/home/src/projects/component-type-checker/packages/sdk",
				),
			},
			cwd:             "/home/src/projects/component-type-checker/packages/app",
			commandLineArgs: []string{"--traceResolution", "--explainFiles"},
		},
		{
			subScenario: "package json scope",
			files: FileMap{
				"/home/src/workspaces/project/src/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"target": "ES2016",
							"composite": true,
							"module": "Node16",
							"traceResolution": true,
						},
						"files": [
							"main.ts",
							"fileA.ts",
							"fileB.mts",
						],
					}`),
				"/home/src/workspaces/project/src/main.ts": "export const x = 10;",
				"/home/src/workspaces/project/src/fileA.ts": stringtestutil.Dedent(`
                    import { foo } from "./fileB.mjs";
                    foo();
                `),
				"/home/src/workspaces/project/src/fileB.mts": "export function foo() {}",
				"/home/src/workspaces/project/package.json": stringtestutil.Dedent(`
                    {
                        "name": "app",
                        "version": "1.0.0"
                    }
                `),
			},
			commandLineArgs: []string{"-p", "src", "--explainFiles", "--extendedDiagnostics"},
			edits: []*tscEdit{
				{
					caption: "Delete package.json",
					edit: func(sys *testSys) {
						sys.removeNoError("/home/src/workspaces/project/package.json")
					},
					// !!! repopulateInfo on diagnostics not yet implemented
					expectedDiff: "Currently we arent repopulating error chain so errors will be different",
				},
			},
		},
		{
			subScenario: "alternateResult",
			files: FileMap{
				"/home/src/projects/project/node_modules/@types/bar/package.json":  getTscModuleResolutionAlternateResultAtTypesPackageJson("bar" /*addTypesCondition*/, false),
				"/home/src/projects/project/node_modules/@types/bar/index.d.ts":    getTscModuleResolutionAlternateResultDts("bar"),
				"/home/src/projects/project/node_modules/bar/package.json":         getTscModuleResolutionAlternateResultPackageJson("bar" /*addTypes*/, false /*addTypesCondition*/, false),
				"/home/src/projects/project/node_modules/bar/index.js":             getTscModuleResolutionAlternateResultJs("bar"),
				"/home/src/projects/project/node_modules/bar/index.mjs":            getTscModuleResolutionAlternateResultMjs("bar"),
				"/home/src/projects/project/node_modules/foo/package.json":         getTscModuleResolutionAlternateResultPackageJson("foo" /*addTypes*/, true /*addTypesCondition*/, false),
				"/home/src/projects/project/node_modules/foo/index.js":             getTscModuleResolutionAlternateResultJs("foo"),
				"/home/src/projects/project/node_modules/foo/index.mjs":            getTscModuleResolutionAlternateResultMjs("foo"),
				"/home/src/projects/project/node_modules/foo/index.d.ts":           getTscModuleResolutionAlternateResultDts("foo"),
				"/home/src/projects/project/node_modules/@types/bar2/package.json": getTscModuleResolutionAlternateResultAtTypesPackageJson("bar2" /*addTypesCondition*/, true),
				"/home/src/projects/project/node_modules/@types/bar2/index.d.ts":   getTscModuleResolutionAlternateResultDts("bar2"),
				"/home/src/projects/project/node_modules/bar2/package.json":        getTscModuleResolutionAlternateResultPackageJson("bar2" /*addTypes*/, false /*addTypesCondition*/, false),
				"/home/src/projects/project/node_modules/bar2/index.js":            getTscModuleResolutionAlternateResultJs("bar2"),
				"/home/src/projects/project/node_modules/bar2/index.mjs":           getTscModuleResolutionAlternateResultMjs("bar2"),
				"/home/src/projects/project/node_modules/foo2/package.json":        getTscModuleResolutionAlternateResultPackageJson("foo2" /*addTypes*/, true /*addTypesCondition*/, true),
				"/home/src/projects/project/node_modules/foo2/index.js":            getTscModuleResolutionAlternateResultJs("foo2"),
				"/home/src/projects/project/node_modules/foo2/index.mjs":           getTscModuleResolutionAlternateResultMjs("foo2"),
				"/home/src/projects/project/node_modules/foo2/index.d.ts":          getTscModuleResolutionAlternateResultDts("foo2"),
				"/home/src/projects/project/index.mts": stringtestutil.Dedent(`
					import { foo } from "foo";
					import { bar } from "bar";
					import { foo2 } from "foo2";
					import { bar2 } from "bar2";
				`),
				"/home/src/projects/project/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"module": "node16",
							"moduleResolution": "node16",
							"traceResolution": true,
							"incremental": true,
							"strict": true,
							"types": [],
						},
						"files": ["index.mts"],
					}`),
			},
			cwd: "/home/src/projects/project",
			edits: []*tscEdit{
				{
					caption: "delete the alternateResult in @types",
					edit: func(sys *testSys) {
						sys.removeNoError("/home/src/projects/project/node_modules/@types/bar/index.d.ts")
					},
					// !!! repopulateInfo on diagnostics not yet implemented
					expectedDiff: "Currently we arent repopulating error chain so errors will be different",
				},
				{
					caption: "delete the node10Result in package/types",
					edit: func(sys *testSys) {
						sys.removeNoError("/home/src/projects/project/node_modules/foo/index.d.ts")
					},
					// !!! repopulateInfo on diagnostics not yet implemented
					expectedDiff: "Currently we arent repopulating error chain so errors will be different",
				},
				{
					caption: "add the alternateResult in @types",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/@types/bar/index.d.ts", getTscModuleResolutionAlternateResultDts("bar"), false)
					},
					// !!! repopulateInfo on diagnostics not yet implemented
					expectedDiff: "Currently we arent repopulating error chain so errors will be different",
				},
				{
					caption: "add the alternateResult in package/types",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/foo/index.d.ts", getTscModuleResolutionAlternateResultDts("foo"), false)
					},
				},
				{
					caption: "update package.json from @types so error is fixed",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/@types/bar/package.json", getTscModuleResolutionAlternateResultAtTypesPackageJson("bar" /*addTypesCondition*/, true), false)
					},
				},
				{
					caption: "update package.json so error is fixed",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/foo/package.json", getTscModuleResolutionAlternateResultPackageJson("foo" /*addTypes*/, true /*addTypesCondition*/, true), false)
					},
				},
				{
					caption: "update package.json from @types so error is introduced",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/@types/bar2/package.json", getTscModuleResolutionAlternateResultAtTypesPackageJson("bar2" /*addTypesCondition*/, false), false)
					},
				},
				{
					caption: "update package.json so error is introduced",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/foo2/package.json", getTscModuleResolutionAlternateResultPackageJson("foo2" /*addTypes*/, true /*addTypesCondition*/, false), false)
					},
				},
				{
					caption: "delete the alternateResult in @types",
					edit: func(sys *testSys) {
						sys.removeNoError("/home/src/projects/project/node_modules/@types/bar2/index.d.ts")
					},
					// !!! repopulateInfo on diagnostics not yet implemented
					expectedDiff: "Currently we arent repopulating error chain so errors will be different",
				},
				{
					caption: "delete the node10Result in package/types",
					edit: func(sys *testSys) {
						sys.removeNoError("/home/src/projects/project/node_modules/foo2/index.d.ts")
					},
					// !!! repopulateInfo on diagnostics not yet implemented
					expectedDiff: "Currently we arent repopulating error chain so errors will be different",
				},
				{
					caption: "add the alternateResult in @types",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/@types/bar2/index.d.ts", getTscModuleResolutionAlternateResultDts("bar2"), false)
					},
					// !!! repopulateInfo on diagnostics not yet implemented
					expectedDiff: "Currently we arent repopulating error chain so errors will be different",
				},
				{
					caption: "add the ndoe10Result in package/types",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/projects/project/node_modules/foo2/index.d.ts", getTscModuleResolutionAlternateResultDts("foo2"), false)
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test.run(t, "moduleResolution")
	}
}

func getBuildModuleResolutionInProjectRefTestCase(preserveSymlinks bool) *tscInput {
	return &tscInput{
		subScenario: `resolves specifier in output declaration file from referenced project correctly` + core.IfElse(preserveSymlinks, " with preserveSymlinks", ""),
		files: FileMap{
			`/user/username/projects/myproject/packages/pkg1/index.ts`: stringtestutil.Dedent(`
				import type { TheNum } from 'pkg2'
				export const theNum: TheNum = 42;`),
			`/user/username/projects/myproject/packages/pkg1/tsconfig.json`: stringtestutil.Dedent(fmt.Sprintf(`
                {
                    "compilerOptions": { 
						"outDir": "build",
						"preserveSymlinks": %t
					},
                    "references": [{ "path": "../pkg2" }]
                }
            `, preserveSymlinks)),
			`/user/username/projects/myproject/packages/pkg2/const.ts`: stringtestutil.Dedent(`
                export type TheNum = 42;
            `),
			`/user/username/projects/myproject/packages/pkg2/index.ts`: stringtestutil.Dedent(`
                export type { TheNum } from 'const';
            `),
			`/user/username/projects/myproject/packages/pkg2/tsconfig.json`: stringtestutil.Dedent(fmt.Sprintf(`
                {
                    "compilerOptions": {
                        "composite": true,
                        "outDir": "build",
                        "paths": {
                            "const": ["./const"]
                        },
                        "preserveSymlinks": %t,
                    },
                }
            `, preserveSymlinks)),
			`/user/username/projects/myproject/packages/pkg2/package.json`: stringtestutil.Dedent(`
                {
                    "name": "pkg2",
                    "version": "1.0.0",
                    "main": "build/index.js"
                }
            `),
			`/user/username/projects/myproject/node_modules/pkg2`: vfstest.Symlink(`/user/username/projects/myproject/packages/pkg2`),
		},
		cwd:             "/user/username/projects/myproject",
		commandLineArgs: []string{"-b", "packages/pkg1", "--verbose", "--traceResolution"},
	}
}

func getTscModuleResolutionSharingFileMap() FileMap {
	return FileMap{
		"/home/src/workspaces/project/packages/a/index.js":      `export const a = 'a';`,
		"/home/src/workspaces/project/packages/a/test/index.js": `import 'a';`,
		"/home/src/workspaces/project/packages/a/tsconfig.json": stringtestutil.Dedent(`
			{
                "compilerOptions": {
                    "checkJs": true,
                    "composite": true,
                    "declaration": true,
                    "emitDeclarationOnly": true,
                    "module": "nodenext",
                    "outDir": "types",
                },
            }`),
		"/home/src/workspaces/project/packages/a/package.json": stringtestutil.Dedent(`
			{
                "name": "a",
                "version": "0.0.0",
                "type": "module",
                "exports": {
                    ".": {
                        "types": "./types/index.d.ts",
                        "default": "./index.js"
                    }
                }
            }`),
		"/home/src/workspaces/project/packages/b/index.js": `export { a } from 'a';`,
		"/home/src/workspaces/project/packages/b/tsconfig.json": stringtestutil.Dedent(`
			{
               "references": [{ "path": "../a" }],
                "compilerOptions": {
                    "checkJs": true,
                    "module": "nodenext",
                    "noEmit": true,
                    "noImplicitAny": true,
                },
            }`),
		"/home/src/workspaces/project/packages/b/package.json": stringtestutil.Dedent(`
			{
                "name": "b",
                "version": "0.0.0",
                "type": "module"
            }`),
		"/home/src/workspaces/project/node_modules/a": vfstest.Symlink("/home/src/workspaces/project/packages/a"),
	}
}

func getTscModuleResolutionAlternateResultAtTypesPackageJson(packageName string, addTypesCondition bool) string {
	var typesString string
	if addTypesCondition {
		typesString = `"types": "./index.d.ts",`
	}
	return stringtestutil.Dedent(fmt.Sprintf(`
		{
			"name": "@types/%s",
			"version": "1.0.0",
			"types": "index.d.ts",
			"exports": {
				".": {
					%s
					"require": "./index.d.ts"
				}
			}
		}`, packageName, typesString))
}

func getTscModuleResolutionAlternateResultPackageJson(packageName string, addTypes bool, addTypesCondition bool) string {
	var types string
	if addTypes {
		types = `"types": "index.d.ts",`
	}
	var typesString string
	if addTypesCondition {
		typesString = `"types": "./index.d.ts",`
	}
	return stringtestutil.Dedent(fmt.Sprintf(`
	{
        "name": "%s",
        "version": "1.0.0",
        "main": "index.js",
        %s
        "exports": {
            ".": {
                %s
                "import": "./index.mjs",
                "require": "./index.js"
            }
        }
    }`, packageName, types, typesString))
}

func getTscModuleResolutionAlternateResultDts(packageName string) string {
	return fmt.Sprintf(`export declare const %s: number;`, packageName)
}

func getTscModuleResolutionAlternateResultJs(packageName string) string {
	return fmt.Sprintf(`module.exports = { %s: 1 };`, packageName)
}

func getTscModuleResolutionAlternateResultMjs(packageName string) string {
	return fmt.Sprintf(`export const %s = 1;`, packageName)
}

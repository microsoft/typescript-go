package execute_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestBuildCommandLine(t *testing.T) {
	t.Parallel()
	testCases := slices.Concat(
		[]*tscInput{
			{
				subScenario:     "help",
				files:           FileMap{},
				commandLineArgs: []string{"--build", "--help"},
			},
			{
				subScenario:     "different options",
				files:           getBuildCommandLineDifferentOptionsMap("composite"),
				commandLineArgs: []string{"--build", "--verbose"},
				edits: []*tscEdit{
					{
						caption:         "with sourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--sourceMap"},
					},
					{
						caption: "should re-emit only js so they dont contain sourcemap",
					},
					{
						caption:         "with declaration should not emit anything",
						commandLineArgs: []string{"--build", "--verbose", "--declaration"},
					},
					noChange,
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
					},
					{
						caption: "should re-emit only dts so they dont contain sourcemap",
					},
					{
						caption:         "with emitDeclarationOnly should not emit anything",
						commandLineArgs: []string{"--build", "--verbose", "--emitDeclarationOnly"},
					},
					noChange,
					{
						caption: "local change",
						edit: func(sys *testSys) {
							sys.replaceFileText("/home/src/workspaces/project/a.ts", "Local = 1", "Local = 10")
						},
					},
					{
						caption:         "with declaration should not emit anything",
						commandLineArgs: []string{"--build", "--verbose", "--declaration"},
					},
					{
						caption:         "with inlineSourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--inlineSourceMap"},
					},
					{
						caption:         "with sourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--sourceMap"},
					},
				},
			},
			{
				subScenario:     "different options with incremental",
				files:           getBuildCommandLineDifferentOptionsMap("incremental"),
				commandLineArgs: []string{"--build", "--verbose"},
				edits: []*tscEdit{
					{
						caption:         "with sourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--sourceMap"},
					},
					{
						caption: "should re-emit only js so they dont contain sourcemap",
					},
					{
						caption:         "with declaration, emit Dts and should not emit js",
						commandLineArgs: []string{"--build", "--verbose", "--declaration"},
					},
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
					},
					noChange,
					{
						caption: "local change",
						edit: func(sys *testSys) {
							sys.replaceFileText("/home/src/workspaces/project/a.ts", "Local = 1", "Local = 10")
						},
					},
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
					},
					noChange,
					{
						caption:         "with inlineSourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--inlineSourceMap"},
					},
					{
						caption:         "with sourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--sourceMap"},
					},
					{
						caption: "emit js files",
					},
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
					},
					{
						caption:         "with declaration and declarationMap, should not re-emit",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
					},
				},
			},
		},
		getBuildCommandLineEmitDeclarationOnlyTestCases([]string{"composite"}, ""),
		getBuildCommandLineEmitDeclarationOnlyTestCases([]string{"incremental", "declaration"}, " with declaration and incremental"),
		getBuildCommandLineEmitDeclarationOnlyTestCases([]string{"declaration"}, " with declaration"),
	)

	for _, test := range testCases {
		test.run(t, "commandLine")
	}
}

func TestBuildClean(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "file name and output name clashing",
			files: FileMap{
				"/home/src/workspaces/solution/index.js": "",
				"/home/src/workspaces/solution/bar.ts":   "",
				"/home/src/workspaces/solution/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": { "allowJs": true }
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--clean"},
		},
		{
			subScenario: "tsx with dts emit",
			files: FileMap{
				"/home/src/workspaces/solution/project/src/main.tsx": "export const x = 10;",
				"/home/src/workspaces/solution/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": { "declaration": true },
					"include": ["src/**/*.tsx", "src/**/*.ts"]
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "project", "-v", "--explainFiles"},
			edits: []*tscEdit{
				noChange,
				{
					caption:         "clean build",
					commandLineArgs: []string{"-b", "project", "--clean"},
				},
			},
		},
	}

	for _, test := range testCases {
		test.run(t, "clean")
	}
}

func TestBuildConfigFileErrors(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "when tsconfig extends the missing file",
			files: FileMap{
				"/home/src/workspaces/project/tsconfig.first.json": stringtestutil.Dedent(`
					{
						"extends": "./foobar.json",
						"compilerOptions": {
							"composite": true
						}
					}`),
				"/home/src/workspaces/project/tsconfig.second.json": stringtestutil.Dedent(`
					{
						"extends": "./foobar.json",
						"compilerOptions": {
							"composite": true
						}
					}`),
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"composite": true
						},
						"references": [
							{ "path": "./tsconfig.first.json" },
							{ "path": "./tsconfig.second.json" }
						]
					}`),
			},
			commandLineArgs: []string{"--b"},
		},
		{
			subScenario: "reports syntax errors in config file",
			files: FileMap{
				"/home/src/workspaces/project/a.ts": "export function foo() { }",
				"/home/src/workspaces/project/b.ts": "export function bar() { }",
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"composite": true,
						},
						"files": [
							"a.ts"
							"b.ts"
						]
					}`),
			},
			commandLineArgs: []string{"--b"},
			edits: []*tscEdit{
				{
					caption: "reports syntax errors after change to config file",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/tsconfig.json", ",", `, "declaration": true`)
					},
				},
				{
					caption: "reports syntax errors after change to ts file",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/project/a.ts", "export function fooBar() { }")
					},
				},
				noChange,
				{
					caption: "builds after fixing config file errors",
					edit: func(sys *testSys) {
						sys.writeFileNoError("/home/src/workspaces/project/tsconfig.json", stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"composite": true, "declaration": true
								},
								"files": [
									"a.ts",
									"b.ts"
								]
							}`), false)
					},
				},
			},
		},
		{
			subScenario:     "missing config file",
			files:           FileMap{},
			commandLineArgs: []string{"--b", "bogus.json"},
		},
	}

	for _, test := range testCases {
		test.run(t, "configFileErrors")
	}
}

func TestBuildEmitDeclarationOnly(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		getBuildEmitDeclarationOnlyTestCase(false),
		getBuildEmitDeclarationOnlyTestCase(true),
		{
			subScenario:     `only dts output in non circular imports project with emitDeclarationOnly`,
			files:           getBuildEmitDeclarationOnlyImportFileMap(true, false),
			commandLineArgs: []string{"--b", "--verbose"},
			edits: []*tscEdit{
				{
					caption: "incremental-declaration-doesnt-change",
					edit: func(sys *testSys) {
						sys.replaceFileText(
							"/home/src/workspaces/project/src/a.ts",
							"export interface A {",
							stringtestutil.Dedent(`
								class C { }
								export interface A {`),
						)
					},
				},
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/src/a.ts", "b: B;", "b: B; foo: any;")
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test.run(t, "emitDeclarationOnly")
	}
}

func TestBuildFileDelete(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "detects deleted file",
			files: FileMap{
				"/home/src/workspaces/solution/child/child.ts": stringtestutil.Dedent(`
					import { child2 } from "../child/child2";
					export function child() {
						child2();
					}
				`),
				"/home/src/workspaces/solution/child/child2.ts": stringtestutil.Dedent(`
					export function child2() {
					}
				`),
				"/home/src/workspaces/solution/child/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": { "composite": true }
					}
				`),
				"/home/src/workspaces/solution/main/main.ts": stringtestutil.Dedent(`
                    import { child } from "../child/child";
                    export function main() {
                        child();
                    }
                `),
				"/home/src/workspaces/solution/main/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": { "composite": true },
						"references": [{ "path": "../child" }],
					}
				`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "main/tsconfig.json", "-v", "--traceResolution", "--explainFiles"},
			edits: []*tscEdit{
				{
					caption: "delete child2 file",
					edit: func(sys *testSys) {
						sys.removeNoError("/home/src/workspaces/solution/child/child2.ts")
						sys.removeNoError("/home/src/workspaces/solution/child/child2.js")
						sys.removeNoError("/home/src/workspaces/solution/child/child2.d.ts")
					},
				},
			},
		},
		{
			subScenario: "deleted file without composite",
			files: FileMap{
				"/home/src/workspaces/solution/child/child.ts": stringtestutil.Dedent(`
					import { child2 } from "../child/child2";
					export function child() {
						child2();
					}
				`),
				"/home/src/workspaces/solution/child/child2.ts": stringtestutil.Dedent(`
					export function child2() {
					}
				`),
				"/home/src/workspaces/solution/child/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": { }
					}
				`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "child/tsconfig.json", "-v", "--traceResolution", "--explainFiles"},
			edits: []*tscEdit{
				{
					caption: "delete child2 file",
					edit: func(sys *testSys) {
						sys.removeNoError("/home/src/workspaces/solution/child/child2.ts")
						sys.removeNoError("/home/src/workspaces/solution/child/child2.js")
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test.run(t, "fileDelete")
	}
}

func TestBuildInferredTypeFromTransitiveModule(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "inferred type from transitive module",
			files:           getBuildInferredTypeFromTransitiveModuleMap(false, ""),
			commandLineArgs: []string{"--b", "--verbose"},
			edits: []*tscEdit{
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/bar.ts", "param: string", "")
					},
				},
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/bar.ts", "foobar()", "foobar(param: string)")
					},
				},
			},
		},
		{
			subScenario:     "inferred type from transitive module with isolatedModules",
			files:           getBuildInferredTypeFromTransitiveModuleMap(true, ""),
			commandLineArgs: []string{"--b", "--verbose"},
			edits: []*tscEdit{
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/bar.ts", "param: string", "")
					},
				},
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/bar.ts", "foobar()", "foobar(param: string)")
					},
				},
			},
		},
		{
			subScenario: "reports errors in files affected by change in signature with isolatedModules",
			files: getBuildInferredTypeFromTransitiveModuleMap(true, stringtestutil.Dedent(`
				import { default as bar } from './bar';
				bar("hello");
			`)),
			commandLineArgs: []string{"--b", "--verbose"},
			edits: []*tscEdit{
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/bar.ts", "param: string", "")
					},
				},
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/bar.ts", "foobar()", "foobar(param: string)")
					},
				},
				{
					caption: "incremental-declaration-changes",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/bar.ts", "param: string", "")
					},
				},
				{
					caption: "Fix Error",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/lazyIndex.ts", `bar("hello")`, "bar()")
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test.run(t, "inferredTypeFromTransitiveModule")
	}
}

func TestBuildJavascriptProjectEmit(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			// !!! sheetal errors seem different
			subScenario: "loads js-based projects and emits them correctly",
			files: FileMap{
				"/home/src/workspaces/solution/common/nominal.js": stringtestutil.Dedent(`
                    /**
                     * @template T, Name
                     * @typedef {T & {[Symbol.species]: Name}} Nominal
                     */
                    module.exports = {};
				`),
				"/home/src/workspaces/solution/common/tsconfig.json": stringtestutil.Dedent(`
					{
						"extends": "../tsconfig.base.json",
						"compilerOptions": {
							"composite": true,
						},
						"include": ["nominal.js"],
					}
				`),
				"/home/src/workspaces/solution/sub-project/index.js": stringtestutil.Dedent(`
                    import { Nominal } from '../common/nominal';

                    /**
                     * @typedef {Nominal<string, 'MyNominal'>} MyNominal
                     */
				`),
				"/home/src/workspaces/solution/sub-project/tsconfig.json": stringtestutil.Dedent(`
				{
					"extends": "../tsconfig.base.json",
					"compilerOptions": {
						"composite": true,
					},
					"references": [
						{ "path": "../common" },
					],
					"include": ["./index.js"],
				}`),
				"/home/src/workspaces/solution/sub-project-2/index.js": stringtestutil.Dedent(`
                    import { MyNominal } from '../sub-project/index';

                    const variable = {
                        key: /** @type {MyNominal} */('value'),
                    };

                    /**
                     * @return {keyof typeof variable}
                     */
                    export function getVar() {
                        return 'key';
                    }
				`),
				"/home/src/workspaces/solution/sub-project-2/tsconfig.json": stringtestutil.Dedent(`
				{
                    "extends": "../tsconfig.base.json",
                    "compilerOptions": {
                        "composite": true,
                    },
                    "references": [
                        { "path": "../sub-project" },
                    ],
                    "include": ["./index.js"],
                }`),
				"/home/src/workspaces/solution/tsconfig.json": stringtestutil.Dedent(`
				{
                    "compilerOptions": {
                        "composite": true,
                    },
                    "references": [
                        { "path": "./sub-project" },
                        { "path": "./sub-project-2" },
                    ],
                    "include": [],
                }`),
				"/home/src/workspaces/solution/tsconfig.base.json": stringtestutil.Dedent(`
				{
                    "compilerOptions": {
                        "skipLibCheck": true,
                        "rootDir": "./",
                        "outDir": "../lib",
                        "allowJs": true,
                        "checkJs": true,
                        "declaration": true,
                    },
                }`),
				tscLibPath + "/lib.d.ts": strings.Replace(tscDefaultLibContent, "interface SymbolConstructor {", "interface SymbolConstructor {\n    readonly species: symbol;", 1),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b"},
		},
		{
			subScenario: `loads js-based projects with non-moved json files and emits them correctly`,
			files: FileMap{
				"/home/src/workspaces/solution/common/obj.json": stringtestutil.Dedent(`
				{
                    "val": 42,
                }`),
				"/home/src/workspaces/solution/common/index.ts": stringtestutil.Dedent(`
                    import x = require("./obj.json");
                    export = x;
                `),
				"/home/src/workspaces/solution/common/tsconfig.json": stringtestutil.Dedent(`
				{
                    "extends": "../tsconfig.base.json",
                    "compilerOptions": {
                        "outDir": null,
                        "composite": true,
                    },
                    "include": ["index.ts", "obj.json"],
                }`),
				"/home/src/workspaces/solution/sub-project/index.js": stringtestutil.Dedent(`
                    import mod from '../common';

                    export const m = mod;
				`),
				"/home/src/workspaces/solution/sub-project/tsconfig.json": stringtestutil.Dedent(`
				{
                    "extends": "../tsconfig.base.json",
                    "compilerOptions": {
                        "composite": true,
                    },
                    "references": [
                        { "path": "../common" },
                    ],
                    "include": ["./index.js"],
                }`),
				"/home/src/workspaces/solution/sub-project-2/index.js": stringtestutil.Dedent(`
                    import { m } from '../sub-project/index';

                    const variable = {
                        key: m,
                    };

                    export function getVar() {
                        return variable;
                    }
				`),
				"/home/src/workspaces/solution/sub-project-2/tsconfig.json": stringtestutil.Dedent(`
				{
					"extends": "../tsconfig.base.json",
					"compilerOptions": {
						"composite": true,
					},
                    "references": [
                        { "path": "../sub-project" },
                    ],
                    "include": ["./index.js"],
                }`),
				"/home/src/workspaces/solution/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
					},
					"references": [
						{ "path": "./sub-project" },
						{ "path": "./sub-project-2" },
                    ],
                    "include": [],
                }`),
				"/home/src/workspaces/solution/tsconfig.base.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"skipLibCheck": true,
						"rootDir": "./",
						"outDir": "../out",
						"allowJs": true,
						"checkJs": true,
						"resolveJsonModule": true,
						"esModuleInterop": true,
						"declaration": true,
					},
                }`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"-b"},
		},
	}

	for _, test := range testCases {
		test.run(t, "javascriptProjectEmit")
	}
}

func TestBuildLateBoundSymbol(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "interface is merged and contains late bound member",
			files: FileMap{
				"/home/src/workspaces/project/src/globals.d.ts": stringtestutil.Dedent(`
                    interface SymbolConstructor {
                        (description?: string | number): symbol;
                    }
                    declare var Symbol: SymbolConstructor;
                `),
				"/home/src/workspaces/project/src/hkt.ts": `export interface HKT<T> { }`,
				"/home/src/workspaces/project/src/main.ts": stringtestutil.Dedent(`
                    import { HKT } from "./hkt";

                    const sym = Symbol();

                    declare module "./hkt" {
                        interface HKT<T> {
                            [sym]: { a: T }
                        }
                    }
                    const x = 10;
                    type A = HKT<number>[typeof sym];
                `),
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
                    "compilerOptions": {
                        "rootDir": "src",
                        "incremental": true,
                    },
                }`),
			},
			commandLineArgs: []string{"--b", "--verbose"},
			edits: []*tscEdit{
				{
					caption: "incremental-declaration-doesnt-change",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/src/main.ts", "const x = 10;", "")
					},
				},
				{
					caption: "incremental-declaration-doesnt-change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/project/src/main.ts", "const x = 10;")
					},
				},
			},
		},
	}

	for _, test := range testCases {
		test.run(t, "lateBoundSymbol")
	}
}

func TestBuildSolutionProject(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "verify that subsequent builds after initial build doesnt build anything",
			files: FileMap{
				"/home/src/workspaces/solution/src/folder/index.ts": `export const x = 10;`,
				"/home/src/workspaces/solution/src/folder/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "files": ["index.ts"],
                        "compilerOptions": {
                            "composite": true
                        }
                    }
                `),
				"/home/src/workspaces/solution/src/folder2/index.ts": `export const x = 10;`,
				"/home/src/workspaces/solution/src/folder2/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "files": ["index.ts"],
                        "compilerOptions": {
                            "composite": true
                        }
                    }
                `),
				"/home/src/workspaces/solution/src/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "files": [],
                        "compilerOptions": {
                            "composite": true
                        },
						"references": [
							{ "path": "./folder" },
							{ "path": "./folder2" },
						]
                }`),
				"/home/src/workspaces/solution/tests/index.ts": `export const x = 10;`,
				"/home/src/workspaces/solution/tests/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "files": ["index.ts"],
                        "compilerOptions": {
                            "composite": true
                        },
                        "references": [
                            { "path": "../src" }
                        ]
                    }
                `),
				"/home/src/workspaces/solution/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "files": [],
                        "compilerOptions": {
                            "composite": true
                        },
                        "references": [
                            { "path": "./src" },
                            { "path": "./tests" }
                        ]
                    }
                `),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--v"},
			edits:           noChangeOnlyEdit,
		},
		{
			subScenario: "when solution is referenced indirectly",
			files: FileMap{
				"/home/src/workspaces/solution/project1/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": { "composite": true },
                        "references": []
                    }
                `),
				"/home/src/workspaces/solution/project2/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": { "composite": true },
                        "references": []
                    }
                `),
				"/home/src/workspaces/solution/project2/src/b.ts": "export const b = 10;",
				"/home/src/workspaces/solution/project3/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": { "composite": true },
                        "references": [
							{ "path": "../project1" },
							{ "path": "../project2" }
						]
                    }
                `),
				"/home/src/workspaces/solution/project3/src/c.ts": "export const c = 10;",
				"/home/src/workspaces/solution/project4/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": { "composite": true },
                        "references": [{ "path": "../project3" }]
                    }
                `),
				"/home/src/workspaces/solution/project4/src/d.ts": "export const d = 10;",
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "project4", "--verbose", "--explainFiles"},
			edits: []*tscEdit{
				{
					caption: "modify project3 file",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/solution/project3/src/c.ts", "c = ", "cc = ")
					},
				},
			},
		},
		{
			subScenario: "has empty files diagnostic when files is empty and no references are provided",
			files: FileMap{
				"/home/src/workspaces/solution/no-references/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "references": [],
                        "files": [],
                        "compilerOptions": {
                            "composite": true,
                            "declaration": true,
                            "forceConsistentCasingInFileNames": true,
                            "skipDefaultLibCheck": true,
                        },
                    }`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "no-references"},
		},
		{
			subScenario: "does not have empty files diagnostic when files is empty and references are provided",
			files: FileMap{
				"/home/src/workspaces/solution/core/index.ts": "export function multiply(a: number, b: number) { return a * b; }",
				"/home/src/workspaces/solution/core/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": {
                            "composite": true,
                            "declaration": true,
                            "declarationMap": true,
                            "skipDefaultLibCheck": true,
                        },
                    }`),
				"/home/src/workspaces/solution/with-references/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "references": [
                            { "path": "../core" },
                        ],
                        "files": [],
                        "compilerOptions": {
                            "composite": true,
                            "declaration": true,
                            "forceConsistentCasingInFileNames": true,
                            "skipDefaultLibCheck": true,
                        },
                    }`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "with-references"},
		},
	}

	for _, test := range testCases {
		test.run(t, "solution")
	}
}

func getBuildCommandLineDifferentOptionsMap(optionName string) FileMap {
	return FileMap{
		"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"%s": true
				}
			}`, optionName)),
		"/home/src/workspaces/project/a.ts": `export const a = 10;const aLocal = 10;`,
		"/home/src/workspaces/project/b.ts": `export const b = 10;const bLocal = 10;`,
		"/home/src/workspaces/project/c.ts": `import { a } from "./a";export const c = a;`,
		"/home/src/workspaces/project/d.ts": `import { b } from "./b";export const d = b;`,
	}
}

func getBuildCommandLineEmitDeclarationOnlyMap(options []string) FileMap {
	compilerOptionsStr := strings.Join(core.Map(options, func(opt string) string {
		return fmt.Sprintf(`"%s": true`, opt)
	}), ", ")
	return FileMap{
		"/home/src/workspaces/solution/project1/src/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": { %s }
			}`, compilerOptionsStr)),
		"/home/src/workspaces/solution/project1/src/a.ts": `export const a = 10;const aLocal = 10;`,
		"/home/src/workspaces/solution/project1/src/b.ts": `export const b = 10;const bLocal = 10;`,
		"/home/src/workspaces/solution/project1/src/c.ts": `import { a } from "./a";export const c = a;`,
		"/home/src/workspaces/solution/project1/src/d.ts": `import { b } from "./b";export const d = b;`,
		"/home/src/workspaces/solution/project2/src/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": { %s },
				"references": [{ "path": "../../project1/src" }]
			}`, compilerOptionsStr)),
		"/home/src/workspaces/solution/project2/src/e.ts": `export const e = 10;`,
		"/home/src/workspaces/solution/project2/src/f.ts": `import { a } from "../../project1/src/a"; export const f = a;`,
		"/home/src/workspaces/solution/project2/src/g.ts": `import { b } from "../../project1/src/b"; export const g = b;`,
	}
}

func getBuildCommandLineEmitDeclarationOnlyTestCases(options []string, suffix string) []*tscInput {
	return []*tscInput{
		{
			subScenario:     "emitDeclarationOnly on commandline" + suffix,
			files:           getBuildCommandLineEmitDeclarationOnlyMap(options),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly"},
			edits: []*tscEdit{
				noChange,
				{
					caption: "local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/a.ts", "const aa = 10;")
					},
				},
				{
					caption: "non local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/a.ts", "export const aaa = 10;")
					},
				},
				{
					caption:         "emit js files",
					commandLineArgs: []string{"--b", "project2/src", "--verbose"},
				},
				noChange,
				{
					caption: "js emit with change without emitDeclarationOnly",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "const alocal = 10;")
					},
					commandLineArgs: []string{"--b", "project2/src", "--verbose"},
				},
				{
					caption: "local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "const aaaa = 10;")
					},
				},
				{
					caption: "non local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "export const aaaaa = 10;")
					},
				},
				{
					caption: "js emit with change without emitDeclarationOnly",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "export const a2 = 10;")
					},
					commandLineArgs: []string{"--b", "project2/src", "--verbose"},
				},
			},
		},
		{
			subScenario:     "emitDeclarationOnly false on commandline" + suffix,
			files:           getBuildCommandLineEmitDeclarationOnlyMap(slices.Concat(options, []string{"emitDeclarationOnly"})),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "project2/src", "--verbose"},
			edits: []*tscEdit{
				noChange,
				{
					caption: "change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/a.ts", "const aa = 10;")
					},
				},
				{
					caption:         "emit js files",
					commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly", "false"},
				},
				noChange,
				{
					caption:         "no change run with js emit",
					commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly", "false"},
				},
				{
					caption: "js emit with change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "const blocal = 10;")
					},
					commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly", "false"},
				},
			},
		},
	}
}

func getBuildEmitDeclarationOnlyImportFileMap(declarationMap bool, circularRef bool) FileMap {
	files := FileMap{
		"/home/src/workspaces/project/src/a.ts": stringtestutil.Dedent(`
			import { B } from "./b";

			export interface A {
				b: B;
			}
		`),
		"/home/src/workspaces/project/src/b.ts": stringtestutil.Dedent(`
			import { C } from "./c";

			export interface B {
				b: C;
			}
		`),
		"/home/src/workspaces/project/src/c.ts": stringtestutil.Dedent(`
			import { A } from "./a";

			export interface C {
				a: A;
			}
		`),
		"/home/src/workspaces/project/src/index.ts": stringtestutil.Dedent(`
			export { A } from "./a";
			export { B } from "./b";
			export { C } from "./c";
		`),
		"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"incremental": true,
					"target": "es5",
					"module": "commonjs",
					"declaration": true,
					"declarationMap": %t,
					"sourceMap": true,
					"outDir": "./lib",
					"composite": true,
					"strict": true,
					"esModuleInterop": true,
					"alwaysStrict": true,
					"rootDir": "src",
					"emitDeclarationOnly": true,
				},
			}`, declarationMap)),
	}
	if !circularRef {
		delete(files, "/home/src/workspaces/project/src/index.ts")
		files["/home/src/workspaces/project/src/a.ts"] = stringtestutil.Dedent(`
			export class B { prop = "hello"; }

			export interface A {
				b: B;
			}
		`)
	}
	return files
}

func getBuildEmitDeclarationOnlyTestCase(declarationMap bool) *tscInput {
	return &tscInput{
		subScenario:     `only dts output in circular import project with emitDeclarationOnly` + core.IfElse(declarationMap, " and declarationMap", ""),
		files:           getBuildEmitDeclarationOnlyImportFileMap(declarationMap, true),
		commandLineArgs: []string{"--b", "--verbose"},
		edits: []*tscEdit{
			{
				caption: "incremental-declaration-changes",
				edit: func(sys *testSys) {
					sys.replaceFileText("/home/src/workspaces/project/src/a.ts", "b: B;", "b: B; foo: any;")
				},
			},
		},
	}
}

func getBuildInferredTypeFromTransitiveModuleMap(isolatedModules bool, lazyExtraContents string) FileMap {
	return FileMap{
		"/home/src/workspaces/project/bar.ts": stringtestutil.Dedent(`
			interface RawAction {
				(...args: any[]): Promise<any> | void;
			}
			interface ActionFactory {
				<T extends RawAction>(target: T): T;
			}
			declare function foo<U extends any[] = any[]>(): ActionFactory;
			export default foo()(function foobar(param: string): void {
			});
		`),
		"/home/src/workspaces/project/bundling.ts": stringtestutil.Dedent(`
			export class LazyModule<TModule> {
				constructor(private importCallback: () => Promise<TModule>) {}
			}

			export class LazyAction<
				TAction extends (...args: any[]) => any,
				TModule
			>  {
				constructor(_lazyModule: LazyModule<TModule>, _getter: (module: TModule) => TAction) {
				}
			}
		`),
		"/home/src/workspaces/project/global.d.ts": stringtestutil.Dedent(`
			interface PromiseConstructor {
				new <T>(): Promise<T>;
			}
			declare var Promise: PromiseConstructor;
			interface Promise<T> {
			}
		`),
		"/home/src/workspaces/project/index.ts": stringtestutil.Dedent(`
			import { LazyAction, LazyModule } from './bundling';
			const lazyModule = new LazyModule(() =>
				import('./lazyIndex')
			);
			export const lazyBar = new LazyAction(lazyModule, m => m.bar);
		`),
		"/home/src/workspaces/project/lazyIndex.ts": stringtestutil.Dedent(`
			export { default as bar } from './bar';
		`) + lazyExtraContents,
		"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
                "compilerOptions": {
                    "target": "es5",
                    "declaration": true,
                    "outDir": "obj",
                    "incremental": true,
					"isolatedModules": %t,
                },
            }`, isolatedModules)),
	}
}

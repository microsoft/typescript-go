package execute_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
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
				files:           getCommandLineDifferentOptionsMap("composite"),
				commandLineArgs: []string{"--build", "--verbose"},
				edits: []*testTscEdit{
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
				files:           getCommandLineDifferentOptionsMap("incremental"),
				commandLineArgs: []string{"--build", "--verbose"},
				edits: []*testTscEdit{
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
		getCommandLineEmitDeclarationOnlyTestCases([]string{"composite"}, ""),
		getCommandLineEmitDeclarationOnlyTestCases([]string{"incremental", "declaration"}, " with declaration and incremental"),
		getCommandLineEmitDeclarationOnlyTestCases([]string{"declaration"}, " with declaration"),
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
			edits: []*testTscEdit{
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
			edits: []*testTscEdit{
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

func TestConfigFileExtends(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "when building solution with projects extends config with include",
			files:           getConfigFileExtendsFileMap(),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--v", "--listFiles"},
		},
		{
			subScenario:     "when building project uses reference and both extend config with include",
			files:           getConfigFileExtendsFileMap(),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "webpack/tsconfig.json", "--v", "--listFiles"},
		},
		{
			subScenario: "resolves the symlink path",
			files: FileMap{
				"/users/user/projects/myconfigs/node_modules/@something/tsconfig-node/tsconfig.json": stringtestutil.Dedent(`
					{
						"extends": "@something/tsconfig-base/tsconfig.json",
						"compilerOptions": {
							"removeComments": true
						}
					}
				`),
				"/users/user/projects/myconfigs/node_modules/@something/tsconfig-base/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": { "composite": true }
					}
				`),
				"/users/user/projects/myproject/src/index.ts": stringtestutil.Dedent(`
					// some comment
					export const x = 10;
				`),
				"/users/user/projects/myproject/src/tsconfig.json": stringtestutil.Dedent(`
					{
						"extends": "@something/tsconfig-node/tsconfig.json"
					}`),
				"/users/user/projects/myproject/node_modules/@something/tsconfig-node": vfstest.Symlink("/users/user/projects/myconfigs/node_modules/@something/tsconfig-node"),
			},
			cwd:             "/users/user/projects/myproject",
			commandLineArgs: []string{"--b", "src", "--extendedDiagnostics"},
		},
		{
			subScenario: "configDir template",
			files: FileMap{
				"/home/src/projects/configs/first/tsconfig.json": stringtestutil.Dedent(`
					{
						"extends": "../second/tsconfig.json",
						"include": ["${configDir}/src"],
						"compilerOptions": {
							"typeRoots": ["root1", "${configDir}/root2", "root3"],
							"types": [],
						},
					}`),
				"/home/src/projects/configs/second/tsconfig.json": stringtestutil.Dedent(`
					{
						"files": ["${configDir}/main.ts"],
						"compilerOptions": {
							"declarationDir": "${configDir}/decls",
							"paths": {
								"@myscope/*": ["${configDir}/types/*"],
							},
						},
						"watchOptions": {
							"excludeFiles": ["${configDir}/main.ts"],
						},
					}`),
				"/home/src/projects/myproject/tsconfig.json": stringtestutil.Dedent(`
					{
						"extends": "../configs/first/tsconfig.json",
						"compilerOptions": {
							"declaration": true,
							"outDir": "outDir",
							"traceResolution": true,
						},
					}`),
				"/home/src/projects/myproject/main.ts": stringtestutil.Dedent(`
					// some comment
					export const y = 10;
					import { x } from "@myscope/sometype";
				`),
				"/home/src/projects/myproject/types/sometype.ts": stringtestutil.Dedent(`
					export const x = 10;
				`),
			},
			cwd:             "/home/src/projects/myproject",
			commandLineArgs: []string{"--b", "--explainFiles", "--v"},
		},
	}

	for _, test := range testCases {
		test.run(t, "configFileExtends")
	}
}

func TestDeclarationEmit(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "when declaration file is referenced through triple slash",
			files:           getDeclarationEmitDtsReferenceAsTrippleSlashMap(),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--verbose"},
		},
		{
			subScenario:     "when declaration file is referenced through triple slash but uses no references",
			files:           getDeclarationEmitDtsReferenceAsTrippleSlashMapNoRef(),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--verbose"},
		},
		{
			subScenario: "when declaration file used inferred type from referenced project",
			files: FileMap{
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"composite": true,
							"paths": { "@fluentui/*": ["./packages/*/src"] },
						},
					}`),
				"/home/src/workspaces/project/packages/pkg1/src/index.ts": stringtestutil.Dedent(`
					export interface IThing {
						a: string;
					}
					export interface IThings {
						thing1: IThing;
					}
				`),
				"/home/src/workspaces/project/packages/pkg1/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "extends": "../../tsconfig",
                        "compilerOptions": { "outDir": "lib" },
                        "include": ["src"],
                    }
                `),
				"/home/src/workspaces/project/packages/pkg2/src/index.ts": stringtestutil.Dedent(`
					import { IThings } from '@fluentui/pkg1';
					export function fn4() {
						const a: IThings = { thing1: { a: 'b' } };
						return a.thing1;
					}
				`),
				"/home/src/workspaces/project/packages/pkg2/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "extends": "../../tsconfig",
                        "compilerOptions": { "outDir": "lib" },
                        "include": ["src"],
                        "references": [{ "path": "../pkg1" }],
                    }
                `),
			},
			commandLineArgs: []string{"--b", "packages/pkg2/tsconfig.json", "--verbose"},
		},
		{
			subScenario:     "reports dts generation errors",
			files:           getDeclarationEmitDtsErrorsFileMap(false),
			commandLineArgs: []string{"-b", "--explainFiles", "--listEmittedFiles", "--v"},
			edits:           noChangeOnlyEdit,
		},
		{
			subScenario:     "reports dts generation errors with incremental",
			files:           getDeclarationEmitDtsErrorsFileMap(true),
			commandLineArgs: []string{"-b", "--explainFiles", "--listEmittedFiles", "--v"},
			edits:           noChangeOnlyEdit,
		},
	}

	for _, test := range testCases {
		test.run(t, "declarationEmit")
	}
}

func TestBuildDemoProject(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "in master branch with everything setup correctly and reports no error",
			files:           getDemoFileMap(),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
			edits:           noChangeOnlyEdit,
		},
		{
			subScenario:     "in circular branch reports the error about it by stopping build",
			files:           getDemoCircularFileMap(),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
		},
		{
			subScenario:     "in bad-ref branch reports the error about files not in rootDir at the import location",
			files:           getDemoBadRefFileMap(),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
		},
	}

	for _, test := range testCases {
		test.run(t, "demo")
	}
}

func TestEmitDeclarationOnly(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		getEmitDeclarationOnlyTestCase(false),
		getEmitDeclarationOnlyTestCase(true),
		{
			subScenario:     `only dts output in non circular imports project with emitDeclarationOnly`,
			files:           getEmitDeclarationOnlyNonCircularImportFileMap(),
			commandLineArgs: []string{"--b", "--verbose"},
			edits: []*testTscEdit{
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

func TestSolutionProject(t *testing.T) {
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
			edits: []*testTscEdit{
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

func getCommandLineDifferentOptionsMap(optionName string) FileMap {
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

func getCommandLineEmitDeclarationOnlyMap(options []string) FileMap {
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

func getCommandLineEmitDeclarationOnlyTestCases(options []string, suffix string) []*tscInput {
	return []*tscInput{
		{
			subScenario:     "emitDeclarationOnly on commandline" + suffix,
			files:           getCommandLineEmitDeclarationOnlyMap(options),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly"},
			edits: []*testTscEdit{
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
			files:           getCommandLineEmitDeclarationOnlyMap(slices.Concat(options, []string{"emitDeclarationOnly"})),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "project2/src", "--verbose"},
			edits: []*testTscEdit{
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

func getConfigFileExtendsFileMap() FileMap {
	return FileMap{
		"/home/src/workspaces/solution/tsconfig.json": stringtestutil.Dedent(`
			{
                "references": [
                    { "path": "./shared/tsconfig.json" },
                    { "path": "./webpack/tsconfig.json" },
                ],
                "files": [],
            }`),
		"/home/src/workspaces/solution/shared/tsconfig-base.json": stringtestutil.Dedent(`
			{
                "include": ["./typings-base/"],
            }`),
		"/home/src/workspaces/solution/shared/typings-base/globals.d.ts": `type Unrestricted = any;`,
		"/home/src/workspaces/solution/shared/tsconfig.json": stringtestutil.Dedent(`
			{
                "extends": "./tsconfig-base.json",
                "compilerOptions": {
                    "composite": true,
                    "outDir": "../target-tsc-build/",
                    "rootDir": "..",
                },
                "files": ["./index.ts"],
            }`),
		"/home/src/workspaces/solution/shared/index.ts": `export const a: Unrestricted = 1;`,
		"/home/src/workspaces/solution/webpack/tsconfig.json": stringtestutil.Dedent(`
			{
                "extends": "../shared/tsconfig-base.json",
                "compilerOptions": {
                    "composite": true,
                    "outDir": "../target-tsc-build/",
                    "rootDir": "..",
                },
                "files": ["./index.ts"],
                "references": [{ "path": "../shared/tsconfig.json" }],
            }`),
		"/home/src/workspaces/solution/webpack/index.ts": `export const b: Unrestricted = 1;`,
	}
}

func getDeclarationEmitDtsReferenceAsTrippleSlashMap() FileMap {
	return FileMap{
		"/home/src/workspaces/solution/tsconfig.base.json": stringtestutil.Dedent(`
			{
                "compilerOptions": {
                    "rootDir": "./",
                    "outDir": "lib",
                },
            }`),
		"/home/src/workspaces/solution/tsconfig.json": stringtestutil.Dedent(`
			{
                "compilerOptions": { "composite": true },
                "references": [{ "path": "./src" }],
                "include": [],
            }`),
		"/home/src/workspaces/solution/src/tsconfig.json": stringtestutil.Dedent(`
			{
                "compilerOptions": { "composite": true },
                "references": [{ "path": "./subProject" }, { "path": "./subProject2" }],
                "include": [],
            }`),
		"/home/src/workspaces/solution/src/subProject/tsconfig.json": stringtestutil.Dedent(`
			{
                "extends": "../../tsconfig.base.json",
                "compilerOptions": { "composite": true },
                "references": [{ "path": "../common" }],
                "include": ["./index.ts"],
            }`),
		"/home/src/workspaces/solution/src/subProject/index.ts": stringtestutil.Dedent(`
			import { Nominal } from '../common/nominal';
			export type MyNominal = Nominal<string, 'MyNominal'>;`),
		"/home/src/workspaces/solution/src/subProject2/tsconfig.json": stringtestutil.Dedent(`
			{
                "extends": "../../tsconfig.base.json",
                "compilerOptions": { "composite": true },
                "references": [{ "path": "../subProject" }],
                "include": ["./index.ts"],
            }`),
		"/home/src/workspaces/solution/src/subProject2/index.ts": stringtestutil.Dedent(`
			import { MyNominal } from '../subProject/index';
			const variable = {
				key: 'value' as MyNominal,
			};
			export function getVar(): keyof typeof variable {
				return 'key';
			}`),
		"/home/src/workspaces/solution/src/common/tsconfig.json": stringtestutil.Dedent(`
			{
				"extends": "../../tsconfig.base.json",
				"compilerOptions": { "composite": true },
				"include": ["./nominal.ts"],
			}`),
		"/home/src/workspaces/solution/src/common/nominal.ts": stringtestutil.Dedent(`
			/// <reference path="./types.d.ts" preserve="true" />
			export declare type Nominal<T, Name extends string> = MyNominal<T, Name>;`),
		"/home/src/workspaces/solution/src/common/types.d.ts": stringtestutil.Dedent(`
			declare type MyNominal<T, Name extends string> = T & {
				specialKey: Name;
			};`),
	}
}

func getDeclarationEmitDtsReferenceAsTrippleSlashMapNoRef() FileMap {
	files := getDeclarationEmitDtsReferenceAsTrippleSlashMap()
	files["/home/src/workspaces/solution/tsconfig.json"] = stringtestutil.Dedent(`
		{
			"extends": "./tsconfig.base.json",
			"compilerOptions": { "composite": true },
			"include": ["./src/**/*.ts"],
		}`)
	return files
}

func getDeclarationEmitDtsErrorsFileMap(incremental bool) FileMap {
	return FileMap{
		"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"module": "NodeNext",
					"moduleResolution": "NodeNext",
					"incremental": %t,
					"declaration": true,
					"skipLibCheck": true,
					"skipDefaultLibCheck": true,
				},
			}`, incremental)),
		"/home/src/workspaces/project/index.ts": stringtestutil.Dedent(`
            import ky from 'ky';
            export const api = ky.extend({});
        `),
		"/home/src/workspaces/project/package.json": stringtestutil.Dedent(`
			{
				"type": "module",
			}`),
		"/home/src/workspaces/project/node_modules/ky/distribution/index.d.ts": stringtestutil.Dedent(`
            type KyInstance = {
                extend(options: Record<string,unknown>): KyInstance;
            }
            declare const ky: KyInstance;
            export default ky;
        `),
		"/home/src/workspaces/project/node_modules/ky/package.json": stringtestutil.Dedent(`
            {
                "name": "ky",
                "type": "module",
                "main": "./distribution/index.js"
            }
        `),
	}
}

func getDemoFileMap() FileMap {
	return FileMap{
		"/user/username/projects/demo/animals/animal.ts": stringtestutil.Dedent(`
            export type Size = "small" | "medium" | "large";
            export default interface Animal {
                size: Size;
            }
        `),
		"/user/username/projects/demo/animals/dog.ts": stringtestutil.Dedent(`
            import Animal from '.';
            import { makeRandomName } from '../core/utilities';

            export interface Dog extends Animal {
                woof(): void;
                name: string;
            }

            export function createDog(): Dog {
                return ({
                    size: "medium",
                    woof: function(this: Dog) {
                        console.log(` + "`" + `${ this.name } says "Woof"!` + "`" + `);
                    },
                    name: makeRandomName()
                });
            }
        `),
		"/user/username/projects/demo/animals/index.ts": stringtestutil.Dedent(`
            import Animal from './animal';

            export default Animal;
            import { createDog, Dog } from './dog';
            export { createDog, Dog };
        `),
		"/user/username/projects/demo/animals/tsconfig.json": stringtestutil.Dedent(`
            {
                "extends": "../tsconfig-base.json",
                "compilerOptions": {
                    "outDir": "../lib/animals",
                    "rootDir": "."
                },
                "references": [
                    { "path": "../core" }
                ]
            }
        `),
		"/user/username/projects/demo/core/utilities.ts": stringtestutil.Dedent(`

            export function makeRandomName() {
                return "Bob!?! ";
            }

            export function lastElementOf<T>(arr: T[]): T | undefined {
                if (arr.length === 0) return undefined;
                return arr[arr.length - 1];
            }
        `),
		"/user/username/projects/demo/core/tsconfig.json": stringtestutil.Dedent(`
			{
				"extends": "../tsconfig-base.json",
				"compilerOptions": {
					"outDir": "../lib/core",
					"rootDir": "."
				},
			}
		`),
		"/user/username/projects/demo/zoo/zoo.ts": stringtestutil.Dedent(`
            import { Dog, createDog } from '../animals/index';

            export function createZoo(): Array<Dog> {
                return [
                    createDog()
                ];
            }
        `),
		"/user/username/projects/demo/zoo/tsconfig.json": stringtestutil.Dedent(`
            {
                "extends": "../tsconfig-base.json",
                "compilerOptions": {
                    "outDir": "../lib/zoo",
                    "rootDir": "."
                },
				"references": [
					{
						"path": "../animals"
					}
				]
        	}
        `),
		"/user/username/projects/demo/tsconfig-base.json": stringtestutil.Dedent(`
			{
				"compilerOptions": {
					"declaration": true,
					"target": "es5",
					"module": "commonjs",
					"strict": true,
					"noUnusedLocals": true,
					"noUnusedParameters": true,
					"noImplicitReturns": true,
					"noFallthroughCasesInSwitch": true,
					"composite": true,
				},
			}
		`),
		"/user/username/projects/demo/tsconfig.json": stringtestutil.Dedent(`
            {
                "files": [],
                "references": [
					{
						"path": "./core"
					},
					{
						"path": "./animals",
					},
					{
						"path": "./zoo",
					},
				],
        	}
		`),
	}
}

func getDemoCircularFileMap() FileMap {
	files := getDemoFileMap()
	files["/user/username/projects/demo/core/tsconfig.json"] = stringtestutil.Dedent(`
		{
			"extends": "../tsconfig-base.json",
			"compilerOptions": {
				"outDir": "../lib/core",
				"rootDir": "."
			},
			"references": [
				{
					"path": "../zoo",
				}
			]
		}
	`)
	return files
}

func getDemoBadRefFileMap() FileMap {
	files := getDemoFileMap()
	files["/user/username/projects/demo/core/utilities.ts"] = `import * as A from '../animals'
` + files["/user/username/projects/demo/core/utilities.ts"].(string)
	return files
}

func getEmitDeclarationOnlyCircularImportFileMap(declarationMap bool) FileMap {
	return FileMap{
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
}

func getEmitDeclarationOnlyTestCase(declarationMap bool) *tscInput {
	return &tscInput{
		subScenario:     `only dts output in circular import project with emitDeclarationOnly` + core.IfElse(declarationMap, " and declarationMap", ""),
		files:           getEmitDeclarationOnlyCircularImportFileMap(declarationMap),
		commandLineArgs: []string{"--b", "--verbose"},
		edits: []*testTscEdit{
			{
				caption: "incremental-declaration-changes",
				edit: func(sys *testSys) {
					sys.replaceFileText("/home/src/workspaces/project/src/a.ts", "b: B;", "b: B; foo: any;")
				},
			},
		},
	}
}

func getEmitDeclarationOnlyNonCircularImportFileMap() FileMap {
	files := getEmitDeclarationOnlyCircularImportFileMap(true)
	delete(files, "/home/src/workspaces/project/src/index.ts")
	files["/home/src/workspaces/project/src/a.ts"] = stringtestutil.Dedent(`
		export class B { prop = "hello"; }

		export interface A {
			b: B;
		}
	`)
	return files
}

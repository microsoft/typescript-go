package execute_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

var (
	expectedVerboseDiff     = "Verbose output status will be different because of up-to-date-ness checks"
	noChangeWithVerboseDiff = &testTscEdit{
		caption:      "no change",
		expectedDiff: expectedVerboseDiff,
	}
	noChangeWithVerboseDiffOOnlyEdit = []*testTscEdit{noChangeWithVerboseDiff}
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
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:      "should re-emit only js so they dont contain sourcemap",
						expectedDiff: expectedVerboseDiff,
					},
					{
						caption:         "with declaration should not emit anything",
						commandLineArgs: []string{"--build", "--verbose", "--declaration"},
						expectedDiff:    expectedVerboseDiff,
					},
					noChangeWithVerboseDiff,
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:      "should re-emit only dts so they dont contain sourcemap",
						expectedDiff: expectedVerboseDiff,
					},
					{
						caption:         "with emitDeclarationOnly should not emit anything",
						commandLineArgs: []string{"--build", "--verbose", "--emitDeclarationOnly"},
						expectedDiff:    expectedVerboseDiff,
					},
					noChangeWithVerboseDiff,
					{
						caption: "local change",
						edit: func(sys *testSys) {
							sys.replaceFileText("/home/src/workspaces/project/a.ts", "Local = 1", "Local = 10")
						},
						expectedDiff: expectedVerboseDiff,
					},
					{
						caption:         "with declaration should not emit anything",
						commandLineArgs: []string{"--build", "--verbose", "--declaration"},
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:         "with inlineSourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--inlineSourceMap"},
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:         "with sourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--sourceMap"},
						expectedDiff:    expectedVerboseDiff,
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
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:      "should re-emit only js so they dont contain sourcemap",
						expectedDiff: expectedVerboseDiff,
					},
					{
						caption:         "with declaration, emit Dts and should not emit js",
						commandLineArgs: []string{"--build", "--verbose", "--declaration"},
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
						expectedDiff:    expectedVerboseDiff,
					},
					noChangeWithVerboseDiff,
					{
						caption: "local change",
						edit: func(sys *testSys) {
							sys.replaceFileText("/home/src/workspaces/project/a.ts", "Local = 1", "Local = 10")
						},
						expectedDiff: expectedVerboseDiff,
					},
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
						expectedDiff:    expectedVerboseDiff,
					},
					noChangeWithVerboseDiff,
					{
						caption:         "with inlineSourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--inlineSourceMap"},
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:         "with sourceMap",
						commandLineArgs: []string{"--build", "--verbose", "--sourceMap"},
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:      "emit js files",
						expectedDiff: expectedVerboseDiff,
					},
					{
						caption:         "with declaration and declarationMap",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
						expectedDiff:    expectedVerboseDiff,
					},
					{
						caption:         "with declaration and declarationMap, should not re-emit",
						commandLineArgs: []string{"--build", "--verbose", "--declaration", "--declarationMap"},
						expectedDiff:    expectedVerboseDiff,
					},
				},
			},
		},
		// !!! sheetal currently errors for project reference needing to have composite is not reported
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
				noChangeWithVerboseDiff,
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
			// !!! sheetal error is not reported
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
	}

	for _, test := range testCases {
		test.run(t, "configFileExtends")
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
			edits:           noChangeWithVerboseDiffOOnlyEdit,
		},
		{
			subScenario:     "in circular branch reports the error about it by stopping build",
			files:           getCircularDemoFileMap(),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
		},
		{
			subScenario:     "in bad-ref branch reports the error about files not in rootDir at the import location",
			files:           getBadRefDemoFileMap(),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
		},
	}

	for _, test := range testCases {
		test.run(t, "demo")
	}
}

// !!! sheetal working on this
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
			edits:           noChangeWithVerboseDiffOOnlyEdit,
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
				noChangeWithVerboseDiff,
				{
					caption: "local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/a.ts", "const aa = 10;")
					},
					expectedDiff: expectedVerboseDiff,
				},
				{
					caption: "non local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/a.ts", "export const aaa = 10;")
					},
					expectedDiff: expectedVerboseDiff,
				},
				{
					caption:         "emit js files",
					commandLineArgs: []string{"--b", "project2/src", "--verbose"},
					expectedDiff:    expectedVerboseDiff,
				},
				noChangeWithVerboseDiff,
				{
					caption: "js emit with change without emitDeclarationOnly",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "const alocal = 10;")
					},
					commandLineArgs: []string{"--b", "project2/src", "--verbose"},
					expectedDiff:    expectedVerboseDiff,
				},
				{
					caption: "local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "const aaaa = 10;")
					},
					expectedDiff: expectedVerboseDiff,
				},
				{
					caption: "non local change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "export const aaaaa = 10;")
					},
					expectedDiff: expectedVerboseDiff,
				},
				{
					caption: "js emit with change without emitDeclarationOnly",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "export const a2 = 10;")
					},
					commandLineArgs: []string{"--b", "project2/src", "--verbose"},
					expectedDiff:    expectedVerboseDiff,
				},
			},
		},
		{
			subScenario:     "emitDeclarationOnly false on commandline" + suffix,
			files:           getCommandLineEmitDeclarationOnlyMap(slices.Concat(options, []string{"emitDeclarationOnly"})),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "project2/src", "--verbose"},
			edits: []*testTscEdit{
				noChangeWithVerboseDiff,
				{
					caption: "change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/a.ts", "const aa = 10;")
					},
					expectedDiff: expectedVerboseDiff,
				},
				{
					caption:         "emit js files",
					commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly", "false"},
					expectedDiff:    expectedVerboseDiff,
				},
				noChangeWithVerboseDiff,
				{
					caption:         "no change run with js emit",
					commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly", "false"},
					expectedDiff:    expectedVerboseDiff,
				},
				{
					caption: "js emit with change",
					edit: func(sys *testSys) {
						sys.appendFile("/home/src/workspaces/solution/project1/src/b.ts", "const blocal = 10;")
					},
					commandLineArgs: []string{"--b", "project2/src", "--verbose", "--emitDeclarationOnly", "false"},
					expectedDiff:    expectedVerboseDiff,
				},
			},
		},
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

func getCircularDemoFileMap() FileMap {
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

func getBadRefDemoFileMap() FileMap {
	files := getDemoFileMap()
	files["/user/username/projects/demo/core/utilities.ts"] = `import * as A from '../animals'
` + files["/user/username/projects/demo/core/utilities.ts"].(string)
	return files
}

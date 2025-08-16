package execute_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestTscLibraryResolution(t *testing.T) {
	t.Parallel()
	testCases := slices.Concat(
		getTscLibResolutionTestCases([]string{"-b", "project1", "project2", "project3", "project4", "--verbose", "--explainFiles"}),
		getTscLibResolutionTestCases([]string{"-p", "project1", "--explainFiles"}),
		[]*tscInput{
			{
				subScenario:     "unknown lib",
				files:           getTscLibraryResolutionUnknown(),
				cwd:             "/home/src/workspace/projects",
				commandLineArgs: []string{"-p", "project1", "--explainFiles"},
			},
			{
				subScenario: "when noLib toggles",
				files: FileMap{
					"/home/src/workspaces/project/a.d.ts": `declare const a = "hello";`,
					"/home/src/workspaces/project/b.ts":   `const b = 10;`,
					"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
                    {
                        "compilerOptions": {
                            "declaration": true,
                            "incremental": true,
                            "lib": ["es6"],
                        },
                    }
                `),
				},
				edits: []*tscEdit{
					{
						caption:         "with --noLib",
						commandLineArgs: []string{"--noLib"},
					},
				},
			},
		},
	)

	for _, test := range testCases {
		test.run(t, "libraryResolution")
	}
}

func getTscLibraryResolutionFileMap(libReplacement bool) FileMap {
	files := FileMap{
		"/home/src/workspace/projects/project1/utils.d.ts": `export const y = 10;`,
		"/home/src/workspace/projects/project1/file.ts":    `export const file = 10;`,
		"/home/src/workspace/projects/project1/core.d.ts":  `export const core = 10;`,
		"/home/src/workspace/projects/project1/index.ts":   `export const x = "type1";`,
		"/home/src/workspace/projects/project1/file2.ts": stringtestutil.Dedent(`
            /// <reference lib="webworker"/>
            /// <reference lib="scripthost"/>
            /// <reference lib="es5"/>
        `),
		"/home/src/workspace/projects/project1/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"composite": true,
					"typeRoots": ["./typeroot1"],
					"lib": ["es5", "dom"],
					"traceResolution": true,
					"libReplacement": %t
				}
			}
		`, libReplacement)),
		"/home/src/workspace/projects/project1/typeroot1/sometype/index.d.ts": `export type TheNum = "type1";`,
		"/home/src/workspace/projects/project2/utils.d.ts":                    `export const y = 10;`,
		"/home/src/workspace/projects/project2/index.ts":                      `export const y = 10`,
		"/home/src/workspace/projects/project2/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
            {
                "compilerOptions": {
                    "composite": true,
                    "lib": ["es5", "dom"],
                    "traceResolution": true,
                    "libReplacement": %t
                }
            }
        `, libReplacement)),
		"/home/src/workspace/projects/project3/utils.d.ts": `export const y = 10;`,
		"/home/src/workspace/projects/project3/index.ts":   `export const z = 10`,
		"/home/src/workspace/projects/project3/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
            {
                "compilerOptions": {
                    "composite": true,
                    "lib": ["es5", "dom"],
                    "traceResolution": true,
                    "libReplacement": %t
                }
            }
        `, libReplacement)),
		"/home/src/workspace/projects/project4/utils.d.ts": `export const y = 10;`,
		"/home/src/workspace/projects/project4/index.ts":   `export const z = 10`,
		"/home/src/workspace/projects/project4/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
            {
                "compilerOptions": {
                    "composite": true,
                    "lib": ["esnext", "dom", "webworker"],
                    "traceResolution": true,
                    "libReplacement": %t
                }
            }
        `, libReplacement)),
		getTestLibPathFor("dom"):        "interface DOMInterface { }",
		getTestLibPathFor("webworker"):  "interface WebWorkerInterface { }",
		getTestLibPathFor("scripthost"): "interface ScriptHostInterface { }",
		"/home/src/workspace/projects/node_modules/@typescript/unlreated/index.d.ts": "export const unrelated = 10;",
	}
	if libReplacement {
		files["/home/src/workspace/projects/node_modules/@typescript/lib-es5/index.d.ts"] = tscDefaultLibContent
		files["/home/src/workspace/projects/node_modules/@typescript/lib-esnext/index.d.ts"] = tscDefaultLibContent
		files["/home/src/workspace/projects/node_modules/@typescript/lib-dom/index.d.ts"] = "interface DOMInterface { }"
		files["/home/src/workspace/projects/node_modules/@typescript/lib-webworker/index.d.ts"] = "interface WebWorkerInterface { }"
		files["/home/src/workspace/projects/node_modules/@typescript/lib-scripthost/index.d.ts"] = "interface ScriptHostInterface { }"
	}
	return files
}

func getTscLibResolutionTestCases(commandLineArgs []string) []*tscInput {
	return []*tscInput{
		{
			subScenario:     "with config",
			files:           getTscLibraryResolutionFileMap(false),
			cwd:             "/home/src/workspace/projects",
			commandLineArgs: commandLineArgs,
		},
		{
			subScenario:     "with config with libReplacement",
			files:           getTscLibraryResolutionFileMap(true),
			cwd:             "/home/src/workspace/projects",
			commandLineArgs: commandLineArgs,
		},
	}
}

func getTscLibraryResolutionUnknown() FileMap {
	return FileMap{
		"/home/src/workspace/projects/project1/utils.d.ts": `export const y = 10;`,
		"/home/src/workspace/projects/project1/file.ts":    `export const file = 10;`,
		"/home/src/workspace/projects/project1/core.d.ts":  `export const core = 10;`,
		"/home/src/workspace/projects/project1/index.ts":   `export const x = "type1";`,
		"/home/src/workspace/projects/project1/file2.ts": stringtestutil.Dedent(`
            /// <reference lib="webworker2"/>
            /// <reference lib="unknownlib"/>
            /// <reference lib="scripthost"/>
        `),
		"/home/src/workspace/projects/project1/tsconfig.json": stringtestutil.Dedent(`
		{
            "compilerOptions": {
                "composite": true,
                "traceResolution": true,
                "libReplacement": true
            }
        }`),
		getTestLibPathFor("webworker"):  "interface WebWorkerInterface { }",
		getTestLibPathFor("scripthost"): "interface ScriptHostInterface { }",
	}
}

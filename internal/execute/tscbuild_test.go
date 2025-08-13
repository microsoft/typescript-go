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

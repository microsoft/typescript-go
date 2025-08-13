package execute_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func TestTscExtends(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "when building solution with projects extends config with include",
			files:           getBuildConfigFileExtendsFileMap(),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--v", "--listFiles"},
		},
		{
			subScenario:     "when building project uses reference and both extend config with include",
			files:           getBuildConfigFileExtendsFileMap(),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "webpack/tsconfig.json", "--v", "--listFiles"},
		},
		getTscExtendsWithSymlinkTestCase("-p"),
		getTscExtendsWithSymlinkTestCase("-b"),
		getTscExtendsConfigDirTestCase("", []string{"--explainFiles"}),
		getTscExtendsConfigDirTestCase(" showConfig", []string{"--showConfig"}),
		getTscExtendsConfigDirTestCase(" with commandline", []string{"--explainFiles", "--outDir", "${configDir}/outDir"}),
		getTscExtendsConfigDirTestCase("", []string{"--b", "--explainFiles", "--v"}),
	}

	for _, test := range testCases {
		test.run(t, "extends")
	}
}

func getBuildConfigFileExtendsFileMap() FileMap {
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

func getTscExtendsWithSymlinkTestCase(builtType string) *tscInput {
	return &tscInput{
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
		commandLineArgs: []string{builtType, "src", "--extendedDiagnostics"},
	}
}

func getTscExtendsConfigDirTestCase(subScenarioSufix string, commandLineArgs []string) *tscInput {
	return &tscInput{
		subScenario: "configDir template" + subScenarioSufix,
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
		commandLineArgs: commandLineArgs,
	}
}

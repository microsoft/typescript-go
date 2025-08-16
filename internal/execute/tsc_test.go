package execute_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestTscCommandline(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "show help with ExitStatus.DiagnosticsPresent_OutputsSkipped",
			env: map[string]string{
				"TS_TEST_TERMINAL_WIDTH": "120",
			},
			commandLineArgs: nil,
		},
		{
			subScenario:     "show help with ExitStatus.DiagnosticsPresent_OutputsSkipped when host cannot provide terminal width",
			commandLineArgs: nil,
		},
		{
			subScenario: "does not add color when NO_COLOR is set",
			env: map[string]string{
				"NO_COLOR": "true",
			},
			commandLineArgs: nil,
		},
		{
			subScenario:     "when build not first argument",
			commandLineArgs: []string{"--verbose", "--build"},
		},
		{
			subScenario:     "help",
			commandLineArgs: []string{"--help"},
		},
		{
			subScenario:     "help all",
			commandLineArgs: []string{"--help", "--all"},
		},
		{
			subScenario:     "Parse --lib option with file name",
			files:           FileMap{"/home/src/workspaces/project/first.ts": `export const Key = Symbol()`},
			commandLineArgs: []string{"--lib", "es6 ", "first.ts"},
		},
		{
			subScenario: "Project is empty string",
			files: FileMap{
				"/home/src/workspaces/project/first.ts": `export const a = 1`,
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"strict": true,
						"noEmit": true
					}
				}`),
			},
			commandLineArgs: []string{},
		},
		{
			subScenario: "Parse -p",
			files: FileMap{
				"/home/src/workspaces/project/first.ts": `export const a = 1`,
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"strict": true,
						"noEmit": true
					}
				}`),
			},
			commandLineArgs: []string{"-p", "."},
		},
		{
			subScenario: "Parse -p with path to tsconfig file",
			files: FileMap{
				"/home/src/workspaces/project/first.ts": `export const a = 1`,
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"strict": true,
						"noEmit": true
					}
				}`),
			},
			commandLineArgs: []string{"-p", "/home/src/workspaces/project/tsconfig.json"},
		},
		{
			subScenario: "Parse -p with path to tsconfig folder",
			files: FileMap{
				"/home/src/workspaces/project/first.ts": `export const a = 1`,
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"strict": true,
						"noEmit": true
					}
				}`),
			},
			commandLineArgs: []string{"-p", "/home/src/workspaces/project"},
		},
		{
			subScenario:     "Parse enum type options",
			commandLineArgs: []string{"--moduleResolution", "nodenext ", "first.ts", "--module", "nodenext", "--target", "esnext", "--moduleDetection", "auto", "--jsx", "react", "--newLine", "crlf"},
		},
		{
			subScenario: "Parse watch interval option",
			files: FileMap{
				"/home/src/workspaces/project/first.ts": `export const a = 1`,
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"strict": true,
						"noEmit": true
					}
				}`),
			},
			commandLineArgs: []string{"-w", "--watchInterval", "1000"},
		},
		{
			subScenario:     "Parse watch interval option without tsconfig.json",
			commandLineArgs: []string{"-w", "--watchInterval", "1000"},
		},
	}

	for _, testCase := range testCases {
		testCase.run(t, "commandLine")
	}
}

func TestTscComposite(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "when setting composite false on command line",
			files: FileMap{
				"/home/src/workspaces/project/src/main.ts": "export const x = 10;",
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"target": "es5",
						"module": "commonjs",
						"composite": true,
					},
					"include": [
						"src/**/*.ts",
					],
				}`),
			},
			commandLineArgs: []string{"--composite", "false"},
		},
		{
			// !!! sheetal null is not reflected in final options
			subScenario: "when setting composite null on command line",
			files: FileMap{
				"/home/src/workspaces/project/src/main.ts": "export const x = 10;",
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"target": "es5",
						"module": "commonjs",
						"composite": true,
					},
					"include": [
						"src/**/*.ts",
					],
				}`),
			},
			commandLineArgs: []string{"--composite", "null"},
		},
		{
			subScenario: "when setting composite false on command line but has tsbuild info in config",
			files: FileMap{
				"/home/src/workspaces/project/src/main.ts": "export const x = 10;",
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"target": "es5",
						"module": "commonjs",
						"composite": true,
						"tsBuildInfoFile": "tsconfig.json.tsbuildinfo",
					},
					"include": [
						"src/**/*.ts",
					],
				}`),
			},
			commandLineArgs: []string{"--composite", "false"},
		},
		{
			subScenario: "when setting composite false and tsbuildinfo as null on command line but has tsbuild info in config",
			files: FileMap{
				"/home/src/workspaces/project/src/main.ts": "export const x = 10;",
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"target": "es5",
						"module": "commonjs",
						"composite": true,
						"tsBuildInfoFile": "tsconfig.json.tsbuildinfo",
					},
					"include": [
						"src/**/*.ts",
					],
				}`),
			},
			commandLineArgs: []string{"--composite", "false", "--tsBuildInfoFile", "null"},
		},
		{
			subScenario: "converting to modules",
			files: FileMap{
				"/home/src/workspaces/project/src/main.ts": "const x = 10;",
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"module": "none",
						"composite": true,
					},
				}`),
			},
			edits: []*tscEdit{
				{
					caption: "convert to modules",
					edit: func(sys *testSys) {
						sys.replaceFileText("/home/src/workspaces/project/tsconfig.json", "none", "es2015")
					},
				},
			},
		},
		{
			subScenario: "synthetic jsx import of ESM module from CJS module no crash no jsx element",
			files: FileMap{
				"/home/src/projects/project/src/main.ts": "export default 42;",
				"/home/src/projects/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"module": "Node16",
						"jsx": "react-jsx",
						"jsxImportSource": "solid-js",
					},
				}`),
				"/home/src/projects/project/node_modules/solid-js/package.json": stringtestutil.Dedent(`
					{
						"name": "solid-js",
						"type": "module"
					}
				`),
				"/home/src/projects/project/node_modules/solid-js/jsx-runtime.d.ts": stringtestutil.Dedent(`
					export namespace JSX {
						type IntrinsicElements = { div: {}; };
					}
				`),
			},
			cwd: "/home/src/projects/project",
		},
		{
			subScenario: "synthetic jsx import of ESM module from CJS module error on jsx element",
			files: FileMap{
				"/home/src/projects/project/src/main.tsx": "export default <div/>;",
				"/home/src/projects/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"module": "Node16",
						"jsx": "react-jsx",
						"jsxImportSource": "solid-js",
					},
				}`),
				"/home/src/projects/project/node_modules/solid-js/package.json": stringtestutil.Dedent(`
					{
						"name": "solid-js",
						"type": "module"
					}
				`),
				"/home/src/projects/project/node_modules/solid-js/jsx-runtime.d.ts": stringtestutil.Dedent(`
					export namespace JSX {
						type IntrinsicElements = { div: {}; };
					}
				`),
			},
			cwd: "/home/src/projects/project",
		},
	}

	for _, testCase := range testCases {
		testCase.run(t, "composite")
	}
}

func TestTscListFilesOnly(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario: "loose file",
			files: FileMap{
				"/home/src/workspaces/project/test.ts": "export const x = 1;",
			},
			commandLineArgs: []string{"test.ts", "--listFilesOnly"},
		},
		{
			subScenario: "combined with incremental",
			files: FileMap{
				"/home/src/workspaces/project/test.ts":       "export const x = 1;",
				"/home/src/workspaces/project/tsconfig.json": "{}",
			},
			commandLineArgs: []string{"--incremental", "--listFilesOnly"},
			edits: []*tscEdit{
				{
					caption:         "incremental actual build",
					commandLineArgs: []string{"--incremental"},
				},
				noChange,
				{
					caption:         "incremental should not build",
					commandLineArgs: []string{"--incremental"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase.run(t, "listFilesOnly")
	}
}

func TestTypeAcquisition(t *testing.T) {
	t.Parallel()
	(&tscInput{
		subScenario: "parse tsconfig with typeAcquisition",
		files: FileMap{
			"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
			{
				"compilerOptions": {
					"composite": true,
					"noEmit": true,
				},
				"typeAcquisition": {
					"enable": true,
					"include": ["0.d.ts", "1.d.ts"],
					"exclude": ["0.js", "1.js"],
					"disableFilenameBasedTypeAcquisition": true,
				},
			}`),
		},
		commandLineArgs: []string{},
	}).run(t, "typeAcquisition")
}

package execute_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestTscNoEmit(t *testing.T) {
	t.Parallel()
	noEmitScenarios := []*tscNoEmitScenario{
		{
			subScenario: "syntax errors",
			aText:       `const a = "hello`,
		},
		{
			subScenario: "semantic errors",
			aText:       `const a: number = "hello"`,
		},
		{
			subScenario: "dts errors",
			aText:       `const a = class { private p = 10; };`,
			dtsEnabled:  true,
		},
		{
			subScenario: "dts errors without dts enabled",
			aText:       `const a = class { private p = 10; };`,
		},
	}
	noEmitChangesScenarios := []*tscNoEmitChangesScenario{
		{
			// !!! sheetal missing initial reporting of Duplicate_identifier_arguments_Compiler_uses_arguments_to_initialize_rest_parameters is absent
			subScenario:   "composite",
			optionsString: `"composite": true`,
		},
		{
			subScenario:   "incremental declaration",
			optionsString: `"incremental": true, "declaration": true`,
		},
		{
			subScenario:   "incremental",
			optionsString: `"incremental": true`,
		},
	}
	testCases := slices.Concat(
		[]*tscInput{
			{
				subScenario: "when project has strict true",
				files: FileMap{
					"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
						{
							"compilerOptions": {
								"incremental": true,
								"strict": true
							}
						}`),
					"/home/src/workspaces/project/class1.ts": `export class class1 {}`,
				},
				commandLineArgs: []string{"--noEmit"},
				edits:           noChangeOnlyEdit,
			},
		},
		getTscNoEmitAndErrorsTestCases(noEmitScenarios, []string{}),
		getTscNoEmitAndErrorsTestCases(noEmitScenarios, []string{"-b", "-v"}),
		getTscNoEmitChangesTestCases(noEmitChangesScenarios, []string{}),
		getTscNoEmitChangesTestCases(noEmitChangesScenarios, []string{"-b", "-v"}),
		getTscNoEmitDtsChangesTestCases(),
		getTscNoEmitDtsChangesMultiFileErrorsTestCases([]string{}),
		getTscNoEmitDtsChangesMultiFileErrorsTestCases([]string{"-b", "-v"}),
	)

	for _, test := range testCases {
		test.run(t, "noEmit")
	}
}

func getTscNoEmitAndErrorsFileMap(scenario *tscNoEmitScenario, incremental bool, asModules bool) FileMap {
	files := FileMap{
		"/home/src/projects/project/a.ts": core.IfElse(asModules, `export `, "") + scenario.aText,
		"/home/src/projects/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"incremental": %t,
					"declaration": %t
				}
			}
		`, incremental, scenario.dtsEnabled)),
	}
	if asModules {
		files["/home/src/projects/project/b.ts"] = `export const b = 10;`
	}
	return files
}

func getTscNoEmitAndErrorsEdits(scenario *tscNoEmitScenario, commandLineArgs []string, asModules bool) []*tscEdit {
	fixedATsContent := core.IfElse(asModules, "export ", "") + `const a = "hello";`
	return []*tscEdit{
		noChange,
		{
			caption: "Fix error",
			edit: func(sys *testSys) {
				sys.writeFileNoError("/home/src/projects/project/a.ts", fixedATsContent, false)
			},
		},
		noChange,
		{
			caption:         "Emit after fixing error",
			commandLineArgs: commandLineArgs,
		},
		noChange,
		{
			caption: "Introduce error",
			edit: func(sys *testSys) {
				sys.writeFileNoError("/home/src/projects/project/a.ts", scenario.aText, false)
			},
		},
		{
			caption:         "Emit when error",
			commandLineArgs: commandLineArgs,
		},
		noChange,
	}
}

type tscNoEmitScenario struct {
	subScenario string
	aText       string
	dtsEnabled  bool
}

func getTscNoEmitAndErrorsTestCases(scenarios []*tscNoEmitScenario, commandLineArgs []string) []*tscInput {
	testingCases := make([]*tscInput, 0, len(scenarios)*3)
	for _, scenario := range scenarios {
		testingCases = append(
			testingCases,
			&tscInput{
				subScenario:     scenario.subScenario,
				commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit"}),
				files:           getTscNoEmitAndErrorsFileMap(scenario, false, false),
				cwd:             "/home/src/projects/project",
				edits:           getTscNoEmitAndErrorsEdits(scenario, commandLineArgs, false),
			},
			&tscInput{
				subScenario:     scenario.subScenario + " with incremental",
				commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit"}),
				files:           getTscNoEmitAndErrorsFileMap(scenario, true, false),
				cwd:             "/home/src/projects/project",
				edits:           getTscNoEmitAndErrorsEdits(scenario, commandLineArgs, false),
			},
			&tscInput{
				subScenario:     scenario.subScenario + " with incremental as modules",
				commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit"}),
				files:           getTscNoEmitAndErrorsFileMap(scenario, true, true),
				cwd:             "/home/src/projects/project",
				edits:           getTscNoEmitAndErrorsEdits(scenario, commandLineArgs, true),
			},
		)
	}
	return testingCases
}

func getTscNoEmitChangesFileMap(optionsStr string) FileMap {
	return FileMap{
		"/home/src/workspaces/project/src/class.ts": stringtestutil.Dedent(`
			export class classC {
				prop = 1;
			}`),
		"/home/src/workspaces/project/src/indirectClass.ts": stringtestutil.Dedent(`
			import { classC } from './class';
			export class indirectClass {
				classC = new classC();
			}`),
		"/home/src/workspaces/project/src/directUse.ts": stringtestutil.Dedent(`
			import { indirectClass } from './indirectClass';
			new indirectClass().classC.prop;`),
		"/home/src/workspaces/project/src/indirectUse.ts": stringtestutil.Dedent(`
			import { indirectClass } from './indirectClass';
			new indirectClass().classC.prop;`),
		"/home/src/workspaces/project/src/noChangeFile.ts": stringtestutil.Dedent(`
			export function writeLog(s: string) {
			}`),
		"/home/src/workspaces/project/src/noChangeFileWithEmitSpecificError.ts": stringtestutil.Dedent(`
			function someFunc(arguments: boolean, ...rest: any[]) {
			}`),
		"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions":  { %s }
			}`, optionsStr)),
	}
}

type tscNoEmitChangesScenario struct {
	subScenario   string
	optionsString string
}

func getTscNoEmitChangesTestCases(scenarios []*tscNoEmitChangesScenario, commandLineArgs []string) []*tscInput {
	noChangeWithNoEmit := &tscEdit{
		caption:         "No Change run with noEmit",
		commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit"}),
	}
	noChangeWithEmit := &tscEdit{
		caption:         "No Change run with emit",
		commandLineArgs: commandLineArgs,
	}
	introduceError := func(sys *testSys) {
		sys.replaceFileText("/home/src/workspaces/project/src/class.ts", "prop", "prop1")
	}
	fixError := func(sys *testSys) {
		sys.replaceFileText("/home/src/workspaces/project/src/class.ts", "prop1", "prop")
	}
	testCases := make([]*tscInput, 0, len(scenarios))
	for _, scenario := range scenarios {
		testCases = append(
			testCases,
			&tscInput{
				subScenario:     "changes " + scenario.subScenario,
				commandLineArgs: commandLineArgs,
				files:           getTscNoEmitChangesFileMap(scenario.optionsString),
				edits: []*tscEdit{
					noChangeWithNoEmit,
					noChangeWithNoEmit,
					{
						caption:         "Introduce error but still noEmit",
						commandLineArgs: noChangeWithNoEmit.commandLineArgs,
						edit:            introduceError,
					},
					{
						caption: "Fix error and emit",
						edit:    fixError,
					},
					noChangeWithEmit,
					noChangeWithNoEmit,
					noChangeWithNoEmit,
					noChangeWithEmit,
					{
						caption: "Introduce error and emit",
						edit:    introduceError,
					},
					noChangeWithEmit,
					noChangeWithNoEmit,
					noChangeWithNoEmit,
					noChangeWithEmit,
					{
						caption:         "Fix error and no emit",
						commandLineArgs: noChangeWithNoEmit.commandLineArgs,
						edit:            fixError,
					},
					noChangeWithEmit,
					noChangeWithNoEmit,
					noChangeWithNoEmit,
					noChangeWithEmit,
				},
			},
			&tscInput{
				subScenario:     "changes with initial noEmit " + scenario.subScenario,
				commandLineArgs: noChangeWithNoEmit.commandLineArgs,
				files:           getTscNoEmitChangesFileMap(scenario.optionsString),
				edits: []*tscEdit{
					noChangeWithEmit,
					{
						caption:         "Introduce error with emit",
						commandLineArgs: commandLineArgs,
						edit:            introduceError,
					},
					{
						caption: "Fix error and no emit",
						edit:    fixError,
					},
					noChangeWithEmit,
				},
			},
		)
	}
	return testCases
}

func getTscNoEmitDtsChangesFileMap(incremental bool, asModules bool) FileMap {
	files := FileMap{
		"/home/src/projects/project/a.ts": core.IfElse(asModules, `export const a = class { private p = 10; };`, `const a = class { private p = 10; };`),
		"/home/src/projects/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"incremental": %t,
				}
			}
		`, incremental)),
	}
	if asModules {
		files["/home/src/projects/project/b.ts"] = `export const b = 10;`
	}
	return files
}

func getTscNoEmitDtsChangesEdits(commandLineArgs []string) []*tscEdit {
	return []*tscEdit{
		noChange,
		{
			caption:         "With declaration enabled noEmit - Should report errors",
			commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit", "--declaration"}),
		},
		{
			caption:         "With declaration and declarationMap noEmit - Should report errors",
			commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit", "--declaration", "--declarationMap"}),
		},
		noChange,
		{
			caption:         "Dts Emit with error",
			commandLineArgs: slices.Concat(commandLineArgs, []string{"--declaration"}),
		},
		{
			caption: "Fix the error",
			edit: func(sys *testSys) {
				sys.replaceFileText("/home/src/projects/project/a.ts", "private", "public")
			},
		},
		{
			caption:         "With declaration enabled noEmit",
			commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit", "--declaration"}),
		},
		{
			caption:         "With declaration and declarationMap noEmit",
			commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit", "--declaration", "--declarationMap"}),
		},
	}
}

func getTscNoEmitDtsChangesTestCases() []*tscInput {
	return []*tscInput{
		{
			subScenario:     "dts errors with declaration enable changes",
			commandLineArgs: []string{"-b", "-v", "--noEmit"},
			files:           getTscNoEmitDtsChangesFileMap(false, false),
			cwd:             "/home/src/projects/project",
			edits:           getTscNoEmitDtsChangesEdits([]string{"-b", "-v"}),
		},
		{
			subScenario:     "dts errors with declaration enable changes with incremental",
			commandLineArgs: []string{"-b", "-v", "--noEmit"},
			files:           getTscNoEmitDtsChangesFileMap(true, false),
			cwd:             "/home/src/projects/project",
			edits:           getTscNoEmitDtsChangesEdits([]string{"-b", "-v"}),
		},
		{
			subScenario:     "dts errors with declaration enable changes with incremental as modules",
			commandLineArgs: []string{"-b", "-v", "--noEmit"},
			files:           getTscNoEmitDtsChangesFileMap(true, true),
			cwd:             "/home/src/projects/project",
			edits:           getTscNoEmitDtsChangesEdits([]string{"-b", "-v"}),
		},
	}
}

func getTscNoEmitDtsChangesMultiFileErrorsTestCases(commandLineArgs []string) []*tscInput {
	aContent := `export const a = class { private p = 10; };`
	return []*tscInput{
		{
			subScenario:     "dts errors with declaration enable changes with multiple files",
			commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit"}),
			files: FileMap{
				"/home/src/projects/project/a.ts": aContent,
				"/home/src/projects/project/b.ts": `export const b = 10;`,
				"/home/src/projects/project/c.ts": strings.Replace(aContent, "a", "c", 1),
				"/home/src/projects/project/d.ts": strings.Replace(aContent, "a", "d", 1),
				"/home/src/projects/project/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"incremental": true,
						}
					}
				`),
			},
			cwd: "/home/src/projects/project",
			edits: slices.Concat(
				getTscNoEmitDtsChangesEdits(commandLineArgs),
				[]*tscEdit{
					{
						caption: "Fix the another ",
						edit: func(sys *testSys) {
							sys.replaceFileText("/home/src/projects/project/c.ts", "private", "public")
						},
						commandLineArgs: slices.Concat(commandLineArgs, []string{"--noEmit", "--declaration", "--declarationMap"}),
					},
				},
			),
		},
	}
}

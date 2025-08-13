package execute_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

type noCheckScenario struct {
	subScenario string
	aText       string
}

func TestTscNoCheck(t *testing.T) {
	t.Parallel()
	cases := []noCheckScenario{
		{"syntax errors", `export const a = "hello`},
		{"semantic errors", `export const a: number = "hello";`},
		{"dts errors", `export const a = class { private p = 10; };`},
	}
	testCases := core.FlatMap(cases, func(c noCheckScenario) []*tscInput {
		return []*tscInput{
			getTscNoCheckTestCase(&c, false, []string{}),
			getTscNoCheckTestCase(&c, true, []string{}),
			getTscNoCheckTestCase(&c, false, []string{"-b", "-v"}),
			getTscNoCheckTestCase(&c, true, []string{"-b", "-v"}),
		}
	})
	for _, test := range testCases {
		test.run(t, "noCheck")
	}
}

func getTscNoCheckTestCase(scenario *noCheckScenario, incremental bool, commandLineArgs []string) *tscInput {
	noChangeWithCheck := &tscEdit{
		caption:         "No Change run with checking",
		commandLineArgs: commandLineArgs,
	}
	fixErrorNoCheck := &tscEdit{
		caption: "Fix `a` error with noCheck",
		edit: func(sys *testSys) {
			sys.writeFileNoError("/home/src/workspaces/project/a.ts", `export const a = "hello";`, false)
		},
	}
	addErrorNoCheck := &tscEdit{
		caption: "Introduce error with noCheck",
		edit: func(sys *testSys) {
			sys.writeFileNoError("/home/src/workspaces/project/a.ts", scenario.aText, false)
		},
	}
	return &tscInput{
		subScenario: scenario.subScenario + core.IfElse(incremental, " with incremental", ""),
		files: FileMap{
			"/home/src/workspaces/project/a.ts": scenario.aText,
			"/home/src/workspaces/project/b.ts": `export const b = 10;`,
			"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"declaration": true,
					"incremental": %t
				}
			}`, incremental)),
		},
		commandLineArgs: slices.Concat(commandLineArgs, []string{"--noCheck"}),
		edits: []*tscEdit{
			noChange,
			fixErrorNoCheck,   // Fix error with noCheck
			noChange,          // Should be no op
			noChangeWithCheck, // Check errors - should not report any errors - update buildInfo
			noChangeWithCheck, // Should be no op
			noChange,          // Should be no op
			addErrorNoCheck,
			noChange,          // Should be no op
			noChangeWithCheck, // Should check errors and update buildInfo
			fixErrorNoCheck,   // Fix error with noCheck
			noChangeWithCheck, // Should check errors and update buildInfo
			{
				caption: "Add file with error",
				edit: func(sys *testSys) {
					sys.writeFileNoError("/home/src/workspaces/project/c.ts", `export const c: number = "hello";`, false)
				},
				commandLineArgs: commandLineArgs,
			},
			addErrorNoCheck,
			fixErrorNoCheck,
			noChangeWithCheck,
			noChange,          // Should be no op
			noChangeWithCheck, // Should be no op
		},
	}
}

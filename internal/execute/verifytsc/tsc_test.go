package verifytsc

import (
	"testing"
)

func TestTsc(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			scenario:    "commandLine",
			subScenario: "show help with ExitStatus.DiagnosticsPresent_OutputsSkipped",
			sys:         NewTestSys(nil, ""),
			// , {
			// 	environmentVariables: new Map([["TS_TEST_TERMINAL_WIDTH", "120"]]),
			// }),
			commandLineArgs: nil,
		},
		{
			scenario:        "commandLine",
			subScenario:     "show help with ExitStatus.DiagnosticsPresent_OutputsSkipped when host can't provide terminal width",
			sys:             NewTestSys(nil, ""),
			commandLineArgs: nil,
		},
		{
			scenario:    "commandLine",
			subScenario: "does not add color when NO_COLOR is set",
			sys:         NewTestSys(nil, ""),
			// , {
			// 		environmentVariables: new Map([["NO_COLOR", "true"]]),
			// 	}),
			commandLineArgs: nil,
		},
		{
			scenario:    "commandLine",
			subScenario: "does not add color when NO_COLOR is set",
			sys:         NewTestSys(nil, ""),
			// , {
			// 	environmentVariables: new Map([["NO_COLOR", "true"]]),
			// }
			// ),
			commandLineArgs: nil,
		},
		{
			scenario:        "commandLine",
			subScenario:     "when build not first argument",
			sys:             NewTestSys(nil, ""),
			commandLineArgs: []string{"--verbose", "--build"},
		},
		{
			scenario:        "commandLine",
			subScenario:     "help",
			sys:             NewTestSys(nil, ""),
			commandLineArgs: []string{"--help"},
		},
		{
			scenario:        "commandLine",
			subScenario:     "help all",
			sys:             NewTestSys(nil, ""),
			commandLineArgs: []string{"--help", "--all"},
		},
	}

	for _, testCase := range testCases {
		testCase.verify(t)
	}

	// todo: temp test, checking that the initial implementention of tsc in tsgo will parse correctly
	(&tscInput{
		scenario:        "commandLine",
		subScenario:     "noEmit and Strict",
		sys:             NewTestSys(nil, ""),
		commandLineArgs: []string{"--noEmit", "--strict"},
	}).verify(t)
}

func TestNoEmit(t *testing.T) {
	t.Parallel()

	(&tscInput{
		scenario:    "noEmit",
		subScenario: "when project has strict true",
		sys: NewTestSys(FileMap{
			"/home/src/workspaces/project/tsconfig.json": `{
	compilerOptions: {
		incremental: true,
		strict: true,
	},
}`,
			"/home/src/workspaces/project/class1.ts": `export class class1 {}`,
		}, ""),
		commandLineArgs: []string{"--noEmit"},
	}).verify(t)
}

func TestProjectReferences(t *testing.T) {
	t.Parallel()
	(&tscInput{
		scenario:    "projectReferences",
		subScenario: "when project references composite project with noEmit",
		sys: NewTestSys(FileMap{
			"/home/src/workspaces/solution/src/utils/index.ts": "export const x = 10;",
			"/home/src/workspaces/solution/src/utils/tsconfig.json": `{
	compilerOptions: {
		composite: true,
		noEmit: true,
	},
})`,
			"/home/src/workspaces/solution/project/index.ts": `import { x } from "../utils";`,
			"/home/src/workspaces/solution/project/tsconfig.json": `{
		references: [
			{ path: "../utils" },
		],
	}),
},`,
		},
			"/home/src/workspaces/solution",
		),
		commandLineArgs: []string{"--p", "project"},
	}).verify(t)
}

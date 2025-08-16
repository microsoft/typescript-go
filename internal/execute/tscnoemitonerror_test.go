package execute_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

type tscNoEmitOnErrorScenario struct {
	subScenario       string
	mainErrorContent  string
	fixedErrorContent string
}

func TestTscNoEmitOnError(t *testing.T) {
	t.Parallel()
	scenarios := []*tscNoEmitOnErrorScenario{
		{
			subScenario: "syntax errors",
			mainErrorContent: stringtestutil.Dedent(`
                import { A } from "../shared/types/db";
                const a = {
                    lastName: 'sdsd'
                ;
            `),
			fixedErrorContent: stringtestutil.Dedent(`
                import { A } from "../shared/types/db";
                const a = {
                    lastName: 'sdsd'
                };`),
		},
		{
			subScenario: "semantic errors",
			mainErrorContent: stringtestutil.Dedent(`
                import { A } from "../shared/types/db";
                const a: string = 10;`),
			fixedErrorContent: stringtestutil.Dedent(`
                import { A } from "../shared/types/db";
                const a: string = "hello";`),
		},
		{
			subScenario: "dts errors",
			mainErrorContent: stringtestutil.Dedent(`
                import { A } from "../shared/types/db";
                export const a = class { private p = 10; };
            `),
			fixedErrorContent: stringtestutil.Dedent(`
                import { A } from "../shared/types/db";
                export const a = class { p = 10; };
            `),
		},
	}
	testCases := slices.Concat(
		getTscNoEmitOnErrorTestCases(scenarios, []string{}),
		getTscNoEmitOnErrorTestCases(scenarios, []string{"-b", "-v"}),
		[]*tscInput{
			{
				subScenario: `when declarationMap changes`,
				files: FileMap{
					"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
						{
							"compilerOptions": {
								"noEmitOnError": true,
								"declaration": true,
								"composite": true,
							},
						}`),
					"/home/src/workspaces/project/a.ts": "const x = 10;",
					"/home/src/workspaces/project/b.ts": "const y = 10;",
				},
				edits: []*tscEdit{
					{
						caption: "error and enable declarationMap",
						edit: func(sys *testSys) {
							sys.replaceFileText("/home/src/workspaces/project/a.ts", "x", "x: 20")
						},
						commandLineArgs: []string{"--declarationMap"},
					},
					{
						caption: "fix error declarationMap",
						edit: func(sys *testSys) {
							sys.replaceFileText("/home/src/workspaces/project/a.ts", "x: 20", "x")
						},
						commandLineArgs: []string{"--declarationMap"},
					},
				},
			},
			{
				subScenario: "file deleted before fixing error with noEmitOnError",
				files: FileMap{
					"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
						{
							"compilerOptions": {
								"outDir": "outDir",
								"noEmitOnError": true,
							},
						}`),
					"/home/src/workspaces/project/file1.ts": `export const x: 30 = "hello";`,
					"/home/src/workspaces/project/file2.ts": `export class D { }`,
				},
				commandLineArgs: []string{"-i"},
				edits: []*tscEdit{
					{
						caption: "delete file without error",
						edit: func(sys *testSys) {
							sys.removeNoError("/home/src/workspaces/project/file2.ts")
						},
					},
				},
			},
		},
	)

	for _, test := range testCases {
		test.run(t, "noEmitOnError")
	}
}

func getTscNoEmitOnErrorFileMap(scenario *tscNoEmitOnErrorScenario, declaration bool, incremental bool) FileMap {
	return FileMap{
		"/user/username/projects/noEmitOnError/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
		 {
            "compilerOptions": {
				"outDir": "./dev-build",
                "declaration": %t,
                "incremental": %t,
                "noEmitOnError": true,
            },
        }`, declaration, incremental)),
		"/user/username/projects/noEmitOnError/shared/types/db.ts": stringtestutil.Dedent(`
            export interface A {
                name: string;
            }
        `),
		"/user/username/projects/noEmitOnError/src/main.ts": scenario.mainErrorContent,
		"/user/username/projects/noEmitOnError/src/other.ts": stringtestutil.Dedent(`
            console.log("hi");
            export { }
        `),
	}
}

func getTscNoEmitOnErrorTestCases(scenarios []*tscNoEmitOnErrorScenario, commandLineArgs []string) []*tscInput {
	testCases := make([]*tscInput, 0, len(scenarios)*4)
	for _, scenario := range scenarios {
		edits := []*tscEdit{
			noChange,
			{
				caption: "Fix error",
				edit: func(sys *testSys) {
					sys.writeFileNoError("/user/username/projects/noEmitOnError/src/main.ts", scenario.fixedErrorContent, false)
				},
			},
			noChange,
		}
		testCases = append(
			testCases,
			&tscInput{
				subScenario:     scenario.subScenario,
				files:           getTscNoEmitOnErrorFileMap(scenario, false, false),
				cwd:             "/user/username/projects/noEmitOnError",
				commandLineArgs: commandLineArgs,
				edits:           edits,
			},
			&tscInput{
				subScenario:     scenario.subScenario + " with declaration",
				files:           getTscNoEmitOnErrorFileMap(scenario, true, false),
				cwd:             "/user/username/projects/noEmitOnError",
				commandLineArgs: commandLineArgs,
				edits:           edits,
			},
			&tscInput{
				subScenario:     scenario.subScenario + " with incremental",
				files:           getTscNoEmitOnErrorFileMap(scenario, false, true),
				cwd:             "/user/username/projects/noEmitOnError",
				commandLineArgs: commandLineArgs,
				edits:           edits,
			},
			&tscInput{
				subScenario:     scenario.subScenario + " with declaration with incremental",
				files:           getTscNoEmitOnErrorFileMap(scenario, true, true),
				cwd:             "/user/username/projects/noEmitOnError",
				commandLineArgs: commandLineArgs,
				edits:           edits,
			},
		)
	}
	return testCases
}

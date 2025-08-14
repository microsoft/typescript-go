package execute_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestBuildDemoProject(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "in master branch with everything setup correctly and reports no error",
			files:           getBuildDemoFileMap(demoBranchMain),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
			edits:           noChangeOnlyEdit,
		},
		{
			subScenario:     "in circular branch reports the error about it by stopping build",
			files:           getBuildDemoFileMap(demoBranchCircularRef),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
		},
		{
			// !!! sheetal - this has missing errors from strada about files not in rootDir (3) and value is declared but not used (1)
			subScenario:     "in bad-ref branch reports the error about files not in rootDir at the import location",
			files:           getBuildDemoFileMap(demoBranchBadRef),
			cwd:             "/user/username/projects/demo",
			commandLineArgs: []string{"--b", "--verbose"},
		},
	}

	for _, test := range testCases {
		test.run(t, "demo")
	}
}

type demoBranch uint

const (
	demoBranchMain = iota
	demoBranchCircularRef
	demoBranchBadRef
)

func getBuildDemoFileMap(demoType demoBranch) FileMap {
	files := FileMap{
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
	switch demoType {
	case demoBranchCircularRef:
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
	case demoBranchBadRef:
		files["/user/username/projects/demo/core/utilities.ts"] = `import * as A from '../animals'
` + files["/user/username/projects/demo/core/utilities.ts"].(string)
	}
	return files
}

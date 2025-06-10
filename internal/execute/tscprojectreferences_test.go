package execute_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
)

func TestProjectReferences(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}

	cases := []tscInput{
		{
			subScenario: "when project references composite project with noEmit",
			sys: newTestSys(FileMap{
				"/home/src/workspaces/solution/utils/index.ts": "export const x = 10;",
				"/home/src/workspaces/solution/utils/tsconfig.json": `{
			"compilerOptions": {
				"composite": true,
				"noEmit": true,
			},
		}`,
				"/home/src/workspaces/solution/project/index.ts": `import { x } from "../utils";`,
				"/home/src/workspaces/solution/project/tsconfig.json": `{
			"references": [
				{ "path": "../utils" },
			],
		}`,
			},
				"/home/src/workspaces/solution",
			),
			commandLineArgs: []string{"--p", "project"},
		},
		{
			subScenario: "when project references composite",
			sys: newTestSys(FileMap{
				"/home/src/workspaces/solution/utils/index.ts":   "export const x = 10;",
				"/home/src/workspaces/solution/utils/index.d.ts": "export declare const x = 10;",
				"/home/src/workspaces/solution/utils/tsconfig.json": `{
	"compilerOptions": {
		"composite": true,
	},
}`,
				"/home/src/workspaces/solution/project/index.ts": `import { x } from "../utils";`,
				"/home/src/workspaces/solution/project/tsconfig.json": `{
	"references": [
		{ "path": "../utils" },
	],
}`,
			}, "/home/src/workspaces/solution"),
			commandLineArgs: []string{"--p", "project"},
		},
	}

	for _, c := range cases {
		c.verify(t, "projectReferences")
	}
}

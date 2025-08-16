package execute_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestProjectReferences(t *testing.T) {
	t.Parallel()
	cases := []tscInput{
		{
			subScenario: "when project references composite project with noEmit",
			files: FileMap{
				"/home/src/workspaces/solution/utils/index.ts": "export const x = 10;",
				"/home/src/workspaces/solution/utils/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"noEmit": true
					}
				}`),
				"/home/src/workspaces/solution/project/index.ts": `import { x } from "../utils";`,
				"/home/src/workspaces/solution/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"references": [
						{ "path": "../utils" },
					],
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "project"},
		},
		{
			subScenario: "when project references composite",
			files: FileMap{
				"/home/src/workspaces/solution/utils/index.ts":   "export const x = 10;",
				"/home/src/workspaces/solution/utils/index.d.ts": "export declare const x = 10;",
				"/home/src/workspaces/solution/utils/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true
					}
				}`),
				"/home/src/workspaces/solution/project/index.ts": `import { x } from "../utils";`,
				"/home/src/workspaces/solution/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"references": [
						{ "path": "../utils" },
					],
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "project"},
		},
		{
			subScenario: "when project reference is not built",
			files: FileMap{
				"/home/src/workspaces/solution/utils/index.ts": "export const x = 10;",
				"/home/src/workspaces/solution/utils/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true
					}
				}`),
				"/home/src/workspaces/solution/project/index.ts": `import { x } from "../utils";`,
				"/home/src/workspaces/solution/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"references": [
						{ "path": "../utils" },
					],
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "project"},
		},
		{
			subScenario: "when project contains invalid project reference",
			files: FileMap{
				"/home/src/workspaces/solution/project/index.ts": `export const x = 10;`,
				"/home/src/workspaces/solution/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"references": [
						{ "path": "../utils" },
					],
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "project"},
		},
		{
			subScenario: "default import interop uses referenced project settings",
			files: FileMap{
				"/home/src/workspaces/project/node_modules/ambiguous-package/package.json": stringtestutil.Dedent(`
				{
					"name": "ambiguous-package"
				}`),
				"/home/src/workspaces/project/node_modules/ambiguous-package/index.d.ts": "export declare const ambiguous: number;",
				"/home/src/workspaces/project/node_modules/esm-package/package.json": stringtestutil.Dedent(`
				{
					"name": "esm-package",
					"type": "module"
				}`),
				"/home/src/workspaces/project/node_modules/esm-package/index.d.ts": "export declare const esm: number;",
				"/home/src/workspaces/project/lib/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"declaration": true,
						"rootDir": "src",
						"outDir": "dist",
						"module": "esnext",
						"moduleResolution": "bundler",
					},
					"include": ["src"],
				}`),
				"/home/src/workspaces/project/lib/src/a.ts":    "export const a = 0;",
				"/home/src/workspaces/project/lib/dist/a.d.ts": "export declare const a = 0;",
				"/home/src/workspaces/project/app/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"module": "esnext",
						"moduleResolution": "bundler",
						"rootDir": "src",
						"outDir": "dist",
					},
					"include": ["src"],
					"references": [
						{ "path": "../lib" },
					],
				}`),
				"/home/src/workspaces/project/app/src/local.ts": "export const local = 0;",
				"/home/src/workspaces/project/app/src/index.ts": stringtestutil.Dedent(`
					import local from "./local"; // Error
					import esm from "esm-package"; // Error
					import referencedSource from "../../lib/src/a"; // Error
					import referencedDeclaration from "../../lib/dist/a"; // Error
					import ambiguous from "ambiguous-package"; // Ok`),
			},
			commandLineArgs: []string{"--p", "app", "--pretty", "false"},
		},
		{
			subScenario: "referencing ambient const enum from referenced project with preserveConstEnums",
			files: FileMap{
				"/home/src/workspaces/solution/utils/index.ts":   "export const enum E { A = 1 }",
				"/home/src/workspaces/solution/utils/index.d.ts": "export declare const enum E { A = 1 }",
				"/home/src/workspaces/solution/utils/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"declaration": true,
						"preserveConstEnums": true,
					},
				}`),
				"/home/src/workspaces/solution/project/index.ts": `import { E } from "../utils"; E.A;`,
				"/home/src/workspaces/solution/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"isolatedModules": true,
					},
					"references": [
						{ "path": "../utils" },
					],
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "project"},
		},
		{
			subScenario: "importing const enum from referenced project with preserveConstEnums and verbatimModuleSyntax",
			files: FileMap{
				"/home/src/workspaces/solution/preserve/index.ts":   "export const enum E { A = 1 }",
				"/home/src/workspaces/solution/preserve/index.d.ts": "export declare const enum E { A = 1 }",
				"/home/src/workspaces/solution/preserve/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"declaration": true,
						"preserveConstEnums": true,
					},
				}`),
				"/home/src/workspaces/solution/no-preserve/index.ts":   "export const enum E { A = 1 }",
				"/home/src/workspaces/solution/no-preserve/index.d.ts": "export declare const enum F { A = 1 }",
				"/home/src/workspaces/solution/no-preserve/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"declaration": true,
						"preserveConstEnums": false,
					},
				}`),
				"/home/src/workspaces/solution/project/index.ts": stringtestutil.Dedent(`
					import { E } from "../preserve";
					import { F } from "../no-preserve";
					E.A;
					F.A;`),
				"/home/src/workspaces/solution/project/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"module": "preserve",
						"verbatimModuleSyntax": true,
					},
					"references": [
						{ "path": "../preserve" },
						{ "path": "../no-preserve" },
					],
				}`),
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "project", "--pretty", "false"},
		},
		{
			subScenario: "rewriteRelativeImportExtensionsProjectReferences1",
			files: FileMap{
				"/home/src/workspaces/packages/common/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"rootDir": "src",
						"outDir": "dist", 
						"module": "nodenext"
					}
				}`),
				"/home/src/workspaces/packages/common/package.json": stringtestutil.Dedent(`
				{
						"name": "common",
						"version": "1.0.0",
						"type": "module",
						"exports": {
							".": {
								"source": "./src/index.ts",
								"default": "./dist/index.js"
							}
						}
				}`),
				"/home/src/workspaces/packages/common/src/index.ts":    "export {};",
				"/home/src/workspaces/packages/common/dist/index.d.ts": "export {};",
				"/home/src/workspaces/packages/main/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"module": "nodenext",
						"rewriteRelativeImportExtensions": true,
						"rootDir": "src",
						"outDir": "dist"
					},
					"references": [
						{ "path": "../common" }
					]
				}`),
				"/home/src/workspaces/packages/main/package.json": stringtestutil.Dedent(`
				{
					"type": "module"
				}`),
				"/home/src/workspaces/packages/main/src/index.ts": `import {} from "../../common/src/index.ts";`,
			},
			cwd:             "/home/src/workspaces",
			commandLineArgs: []string{"-p", "packages/main", "--pretty", "false"},
		},
		{
			subScenario: "rewriteRelativeImportExtensionsProjectReferences2",
			files: FileMap{
				"/home/src/workspaces/solution/src/tsconfig-base.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"module": "nodenext",
						"composite": true,
						"rootDir": ".",
						"outDir": "../dist",
						"rewriteRelativeImportExtensions": true
					}
				}`),
				"/home/src/workspaces/solution/src/compiler/tsconfig.json": stringtestutil.Dedent(`
				{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {}
				}`),
				"/home/src/workspaces/solution/src/compiler/parser.ts":    "export {};",
				"/home/src/workspaces/solution/dist/compiler/parser.d.ts": "export {};",
				"/home/src/workspaces/solution/src/services/tsconfig.json": stringtestutil.Dedent(`
				{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {},
					"references": [
						{ "path": "../compiler" }
					]
				}`),
				"/home/src/workspaces/solution/src/services/services.ts": `import {} from "../compiler/parser.ts";`,
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "src/services", "--pretty", "false"},
		},
		{
			subScenario: "rewriteRelativeImportExtensionsProjectReferences3",
			files: FileMap{
				"/home/src/workspaces/solution/src/tsconfig-base.json": stringtestutil.Dedent(`
				{
					"compilerOptions": { 
						"module": "nodenext",
						"composite": true,
						"rewriteRelativeImportExtensions": true
					}
				}`),
				"/home/src/workspaces/solution/src/compiler/tsconfig.json": stringtestutil.Dedent(`
				{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {
						"rootDir": ".",
						"outDir": "../../dist/compiler"
					}
				}`),
				"/home/src/workspaces/solution/src/compiler/parser.ts":    "export {};",
				"/home/src/workspaces/solution/dist/compiler/parser.d.ts": "export {};",
				"/home/src/workspaces/solution/src/services/tsconfig.json": stringtestutil.Dedent(`
				{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {
						"rootDir": ".", 
						"outDir": "../../dist/services"
					},
					"references": [
						{ "path": "../compiler" }
					]
				}`),
				"/home/src/workspaces/solution/src/services/services.ts": `import {} from "../compiler/parser.ts";`,
			},
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--p", "src/services", "--pretty", "false"},
		},
		{
			subScenario: "default setup was created correctly",
			files: FileMap{
				"/home/src/workspaces/project/primary/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					}
				}`),
				"/home/src/workspaces/project/primary/a.ts": "export { };",
				"/home/src/workspaces/project/secondary/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"references": [{
						"path": "../primary"
					}]
				}`),
				"/home/src/workspaces/project/secondary/b.ts": `import * as mod_1 from "../primary/a";`,
			},
			commandLineArgs: []string{"--p", "primary/tsconfig.json"},
		},
		{
			subScenario: "errors when declaration = false",
			files: FileMap{
				"/home/src/workspaces/project/primary/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
						"declaration": false
					}
				}`),
				"/home/src/workspaces/project/primary/a.ts": "export { };",
			},
			commandLineArgs: []string{"--p", "primary/tsconfig.json"},
		},
		{
			subScenario: "errors when the referenced project doesnt have composite",
			files: FileMap{
				"/home/src/workspaces/project/primary/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": false,
						"outDir": "bin",
					}
				}`),
				"/home/src/workspaces/project/primary/a.ts": "export { };",
				"/home/src/workspaces/project/reference/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"files": [ "b.ts" ],
					"references": [ { "path": "../primary" } ]
				}`),
				"/home/src/workspaces/project/reference/b.ts": `import * as mod_1 from "../primary/a";`,
			},
			commandLineArgs: []string{"--p", "reference/tsconfig.json"},
		},
		{
			subScenario: "does not error when the referenced project doesnt have composite if its a container project",
			files: FileMap{
				"/home/src/workspaces/project/primary/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": false,
						"outDir": "bin",
					}
				}`),
				"/home/src/workspaces/project/primary/a.ts": "export { };",
				"/home/src/workspaces/project/reference/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"files": [ ],
					"references": [{
						"path": "../primary"
					}]
				}`),
				"/home/src/workspaces/project/reference/b.ts": `import * as mod_1 from "../primary/a";`,
			},
			commandLineArgs: []string{"--p", "reference/tsconfig.json"},
		},
		{
			subScenario: "errors when the file list is not exhaustive",
			files: FileMap{
				"/home/src/workspaces/project/primary/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"files": [ "a.ts" ]
				}`),
				"/home/src/workspaces/project/primary/a.ts": "import * as b from './b'",
				"/home/src/workspaces/project/primary/b.ts": "export {}",
			},
			commandLineArgs: []string{"--p", "primary/tsconfig.json"},
		},
		{
			subScenario: "errors when the referenced project doesnt exist",
			files: FileMap{
				"/home/src/workspaces/project/primary/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"references": [{
						"path": "../foo"
					}]
				}`),
				"/home/src/workspaces/project/primary/a.ts": "export { };",
			},
			commandLineArgs: []string{"--p", "primary/tsconfig.json"},
		},
		{
			subScenario: "redirects to the output dts file",
			files: FileMap{
				"/home/src/workspaces/project/alpha/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					}
				}`),
				"/home/src/workspaces/project/alpha/a.ts":       "export const m: number = 3;",
				"/home/src/workspaces/project/alpha/bin/a.d.ts": "export { };",
				"/home/src/workspaces/project/beta/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"references": [ { "path": "../alpha" } ]
				}`),
				"/home/src/workspaces/project/beta/b.ts": "import { m } from '../alpha/a'",
			},
			commandLineArgs: []string{"--p", "beta/tsconfig.json", "--explainFiles"},
		},
		{
			subScenario: "issues a nice error when the input file is missing",
			files: FileMap{
				"/home/src/workspaces/project/alpha/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"references": []
				}`),
				"/home/src/workspaces/project/alpha/a.ts": "export const m: number = 3;",
				"/home/src/workspaces/project/beta/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"references": [ { "path": "../alpha" } ]
				}`),
				"/home/src/workspaces/project/beta/b.ts": "import { m } from '../alpha/a'",
			},
			commandLineArgs: []string{"--p", "beta/tsconfig.json"},
		},
		{
			subScenario: "issues a nice error when the input file is missing when module reference is not relative",
			files: FileMap{
				"/home/src/workspaces/project/alpha/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					}
				}`),
				"/home/src/workspaces/project/alpha/a.ts": "export const m: number = 3;",
				"/home/src/workspaces/project/beta/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
						"paths": {
                            "@alpha/*": ["../alpha/*"],
                        },
					},
					"references": [ { "path": "../alpha" } ]
				}`),
				"/home/src/workspaces/project/beta/b.ts": "import { m } from '@alpha/a'",
			},
			commandLineArgs: []string{"--p", "beta/tsconfig.json"},
		},
		{
			subScenario: "doesnt infer the rootDir from source paths",
			files: FileMap{
				"/home/src/workspaces/project/alpha/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"references": []
				}`),
				"/home/src/workspaces/project/alpha/src/a.ts": "export const m: number = 3;",
			},
			commandLineArgs: []string{"--p", "alpha/tsconfig.json"},
		},
		{
			// !!! sheetal rootDir error not reported
			subScenario: "errors when a file is outside the rootdir",
			files: FileMap{
				"/home/src/workspaces/project/alpha/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "bin",
					},
					"references": []
				}`),
				"/home/src/workspaces/project/alpha/src/a.ts": "import * as b from '../../beta/b'",
				"/home/src/workspaces/project/beta/b.ts":      "export { }",
			},
			commandLineArgs: []string{"--p", "alpha/tsconfig.json"},
		},
	}

	for _, c := range cases {
		c.run(t, "projectReferences")
	}
}

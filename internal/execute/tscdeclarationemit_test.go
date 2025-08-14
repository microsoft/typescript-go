package execute_test

import (
	"fmt"
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func TestTscDeclarationEmit(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "when declaration file is referenced through triple slash",
			files:           getBuildDeclarationEmitDtsReferenceAsTrippleSlashMap(false),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--verbose"},
		},
		{
			subScenario:     "when declaration file is referenced through triple slash but uses no references",
			files:           getBuildDeclarationEmitDtsReferenceAsTrippleSlashMap(true),
			cwd:             "/home/src/workspaces/solution",
			commandLineArgs: []string{"--b", "--verbose"},
		},
		{
			subScenario: "when declaration file used inferred type from referenced project",
			files: FileMap{
				"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"composite": true,
							"paths": { "@fluentui/*": ["./packages/*/src"] },
						},
					}`),
				"/home/src/workspaces/project/packages/pkg1/src/index.ts": stringtestutil.Dedent(`
					export interface IThing {
						a: string;
					}
					export interface IThings {
						thing1: IThing;
					}
				`),
				"/home/src/workspaces/project/packages/pkg1/tsconfig.json": stringtestutil.Dedent(`
					{
						"extends": "../../tsconfig",
						"compilerOptions": { "outDir": "lib" },
						"include": ["src"],
					}
				`),
				"/home/src/workspaces/project/packages/pkg2/src/index.ts": stringtestutil.Dedent(`
					import { IThings } from '@fluentui/pkg1';
					export function fn4() {
						const a: IThings = { thing1: { a: 'b' } };
						return a.thing1;
					}
				`),
				"/home/src/workspaces/project/packages/pkg2/tsconfig.json": stringtestutil.Dedent(`
					{
						"extends": "../../tsconfig",
						"compilerOptions": { "outDir": "lib" },
						"include": ["src"],
						"references": [{ "path": "../pkg1" }],
					}
				`),
			},
			commandLineArgs: []string{"--b", "packages/pkg2/tsconfig.json", "--verbose"},
		},
		{
			subScenario:     "reports dts generation errors",
			files:           getTscDeclarationEmitDtsErrorsFileMap(false, false),
			commandLineArgs: []string{"-b", "--explainFiles", "--listEmittedFiles", "--v"},
			edits:           noChangeOnlyEdit,
		},
		{
			subScenario:     "reports dts generation errors with incremental",
			files:           getTscDeclarationEmitDtsErrorsFileMap(false, true),
			commandLineArgs: []string{"-b", "--explainFiles", "--listEmittedFiles", "--v"},
			edits:           noChangeOnlyEdit,
		},
		{
			subScenario:     "reports dts generation errors",
			files:           getTscDeclarationEmitDtsErrorsFileMap(false, false),
			commandLineArgs: []string{"--explainFiles", "--listEmittedFiles"},
			edits: []*tscEdit{
				noChange,
				{
					caption:         "build -b",
					commandLineArgs: []string{"-b", "--explainFiles", "--listEmittedFiles", "--v"},
				},
			},
		},
		{
			subScenario:     "reports dts generation errors with incremental",
			files:           getTscDeclarationEmitDtsErrorsFileMap(true, true),
			commandLineArgs: []string{"--explainFiles", "--listEmittedFiles"},
			edits: []*tscEdit{
				noChange,
				{
					caption:         "build -b",
					commandLineArgs: []string{"-b", "--explainFiles", "--listEmittedFiles", "--v"},
				},
			},
		},
		{
			subScenario: "when using Windows paths and uppercase letters",
			files: FileMap{
				"D:/Work/pkg1/package.json": stringtestutil.Dedent(`
				{
					"name": "ts-specifier-bug",
					"version": "1.0.0",
					"main": "index.js"
				}`),
				"D:/Work/pkg1/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"declaration": true,
						"target": "es2017",
						"outDir": "./dist",
					},
					"include": ["src"],
				}`),
				"D:/Work/pkg1/src/main.ts": stringtestutil.Dedent(`
					import { PartialType } from './utils';

					class Common {}
					
					export class Sub extends PartialType(Common) {
						id: string;
					}
				`),
				"D:/Work/pkg1/src/utils/index.ts": stringtestutil.Dedent(`
					import { MyType, MyReturnType } from './type-helpers';

					export function PartialType<T>(classRef: MyType<T>) {
						abstract class PartialClassType {
							constructor() {}
						}
					
						return PartialClassType as MyReturnType;
					}
				`),
				"D:/Work/pkg1/src/utils/type-helpers.ts": stringtestutil.Dedent(`
					export type MyReturnType = {	
						new (...args: any[]): any;
					};
				
					export interface MyType<T = any> extends Function {
						new (...args: any[]): T;
					}
				`),
			},
			cwd:              "D:/Work/pkg1",
			windowsStyleRoot: "D:/",
			ignoreCase:       true,
			commandLineArgs:  []string{"-p", "D:\\Work\\pkg1", "--explainFiles"},
		},
		{
			// !!! sheetal redirected files not yet implemented
			subScenario: "when same version is referenced through source and another symlinked package",
			files: FileMap{
				`/user/username/projects/myproject/plugin-two/index.d.ts`:                               pluginTwoDts(),
				`/user/username/projects/myproject/plugin-two/node_modules/typescript-fsa/package.json`: fsaPackageJson(),
				`/user/username/projects/myproject/plugin-two/node_modules/typescript-fsa/index.d.ts`:   fsaIndex(),
				`/user/username/projects/myproject/plugin-one/tsconfig.json`:                            pluginOneConfig(),
				`/user/username/projects/myproject/plugin-one/index.ts`:                                 pluginOneIndex(),
				`/user/username/projects/myproject/plugin-one/action.ts`:                                pluginOneAction(),
				`/user/username/projects/myproject/plugin-one/node_modules/typescript-fsa/package.json`: fsaPackageJson(),
				`/user/username/projects/myproject/plugin-one/node_modules/typescript-fsa/index.d.ts`:   fsaIndex(),
				`/user/username/projects/myproject/plugin-one/node_modules/plugin-two`:                  vfstest.Symlink(`/user/username/projects/myproject/plugin-two`),
			},
			cwd:             "/user/username/projects/myproject",
			commandLineArgs: []string{"-p", "plugin-one", "--explainFiles"},
		},
		{
			// !!! sheetal redirected files not yet implemented
			subScenario: "when same version is referenced through source and another symlinked package with indirect link",
			files: FileMap{
				`/user/username/projects/myproject/plugin-two/package.json`: stringtestutil.Dedent(`
				{
					"name": "plugin-two",
					"version": "0.1.3",
					"main": "dist/commonjs/index.js"
				}`),
				`/user/username/projects/myproject/plugin-two/dist/commonjs/index.d.ts`:                 pluginTwoDts(),
				`/user/username/projects/myproject/plugin-two/node_modules/typescript-fsa/package.json`: fsaPackageJson(),
				`/user/username/projects/myproject/plugin-two/node_modules/typescript-fsa/index.d.ts`:   fsaIndex(),
				`/user/username/projects/myproject/plugin-one/tsconfig.json`:                            pluginOneConfig(),
				`/user/username/projects/myproject/plugin-one/index.ts`:                                 pluginOneIndex() + "\n" + pluginOneAction(),
				`/user/username/projects/myproject/plugin-one/node_modules/typescript-fsa/package.json`: fsaPackageJson(),
				`/user/username/projects/myproject/plugin-one/node_modules/typescript-fsa/index.d.ts`:   fsaIndex(),
				`/temp/yarn/data/link/plugin-two`:                                                       vfstest.Symlink(`/user/username/projects/myproject/plugin-two`),
				`/user/username/projects/myproject/plugin-one/node_modules/plugin-two`:                  vfstest.Symlink(`/temp/yarn/data/link/plugin-two`),
			},
			cwd:             "/user/username/projects/myproject",
			commandLineArgs: []string{"-p", "plugin-one", "--explainFiles"},
		},
		{
			// !!! sheetal strada has error for d.ts generation in pkg3/src/keys.ts but corsa doesnt have that
			subScenario: "when pkg references sibling package through indirect symlink",
			files: FileMap{
				`/user/username/projects/myproject/pkg1/dist/index.d.ts`: `export * from './types';`,
				`/user/username/projects/myproject/pkg1/dist/types.d.ts`: stringtestutil.Dedent(`
					export declare type A = {
						id: string;
					};
					export declare type B = {
						id: number;
					};
					export declare type IdType = A | B;
					export declare class MetadataAccessor<T, D extends IdType = IdType> {
						readonly key: string;
						private constructor();
						toString(): string;
						static create<T, D extends IdType = IdType>(key: string): MetadataAccessor<T, D>;
					}`),
				`/user/username/projects/myproject/pkg1/package.json`: stringtestutil.Dedent(`
					{
						"name": "@raymondfeng/pkg1",
						"version": "1.0.0",
						"main": "dist/index.js",
						"typings": "dist/index.d.ts"
					}`),
				`/user/username/projects/myproject/pkg2/dist/index.d.ts`: `export * from './types';`,
				`/user/username/projects/myproject/pkg2/dist/types.d.ts`: `export {MetadataAccessor} from '@raymondfeng/pkg1';`,
				`/user/username/projects/myproject/pkg2/package.json`: stringtestutil.Dedent(`
					{
						"name": "@raymondfeng/pkg2",
						"version": "1.0.0",
						"main": "dist/index.js",
						"typings": "dist/index.d.ts"
					}`),
				`/user/username/projects/myproject/pkg3/src/index.ts`: `export * from './keys';`,
				`/user/username/projects/myproject/pkg3/src/keys.ts`: stringtestutil.Dedent(`
					import {MetadataAccessor} from "@raymondfeng/pkg2";
					export const ADMIN = MetadataAccessor.create<boolean>('1');`),
				`/user/username/projects/myproject/pkg3/tsconfig.json`: stringtestutil.Dedent(`
                    {
                        "compilerOptions": {
                            "outDir": "dist",
                            "rootDir": "src",
                            "target": "es5",
                            "module": "commonjs",
                            "strict": true,
                            "esModuleInterop": true,
                            "declaration": true,
                        },
                    }`),
				`/user/username/projects/myproject/pkg2/node_modules/@raymondfeng/pkg1`: vfstest.Symlink(`/user/username/projects/myproject/pkg1`),
				`/user/username/projects/myproject/pkg3/node_modules/@raymondfeng/pkg2`: vfstest.Symlink(`/user/username/projects/myproject/pkg2`),
			},
			cwd:             "/user/username/projects/myproject",
			commandLineArgs: []string{"-p", "pkg3", "--explainFiles"},
		},
	}

	for _, test := range testCases {
		test.run(t, "declarationEmit")
	}
}

func getBuildDeclarationEmitDtsReferenceAsTrippleSlashMap(useNoRef bool) FileMap {
	files := FileMap{
		"/home/src/workspaces/solution/tsconfig.base.json": stringtestutil.Dedent(`
			{
                "compilerOptions": {
                    "rootDir": "./",
                    "outDir": "lib",
                },
            }`),
		"/home/src/workspaces/solution/tsconfig.json": stringtestutil.Dedent(`
			{
                "compilerOptions": { "composite": true },
                "references": [{ "path": "./src" }],
                "include": [],
            }`),
		"/home/src/workspaces/solution/src/tsconfig.json": stringtestutil.Dedent(`
			{
                "compilerOptions": { "composite": true },
                "references": [{ "path": "./subProject" }, { "path": "./subProject2" }],
                "include": [],
            }`),
		"/home/src/workspaces/solution/src/subProject/tsconfig.json": stringtestutil.Dedent(`
			{
                "extends": "../../tsconfig.base.json",
                "compilerOptions": { "composite": true },
                "references": [{ "path": "../common" }],
                "include": ["./index.ts"],
            }`),
		"/home/src/workspaces/solution/src/subProject/index.ts": stringtestutil.Dedent(`
			import { Nominal } from '../common/nominal';
			export type MyNominal = Nominal<string, 'MyNominal'>;`),
		"/home/src/workspaces/solution/src/subProject2/tsconfig.json": stringtestutil.Dedent(`
			{
                "extends": "../../tsconfig.base.json",
                "compilerOptions": { "composite": true },
                "references": [{ "path": "../subProject" }],
                "include": ["./index.ts"],
            }`),
		"/home/src/workspaces/solution/src/subProject2/index.ts": stringtestutil.Dedent(`
			import { MyNominal } from '../subProject/index';
			const variable = {
				key: 'value' as MyNominal,
			};
			export function getVar(): keyof typeof variable {
				return 'key';
			}`),
		"/home/src/workspaces/solution/src/common/tsconfig.json": stringtestutil.Dedent(`
			{
				"extends": "../../tsconfig.base.json",
				"compilerOptions": { "composite": true },
				"include": ["./nominal.ts"],
			}`),
		"/home/src/workspaces/solution/src/common/nominal.ts": stringtestutil.Dedent(`
			/// <reference path="./types.d.ts" preserve="true" />
			export declare type Nominal<T, Name extends string> = MyNominal<T, Name>;`),
		"/home/src/workspaces/solution/src/common/types.d.ts": stringtestutil.Dedent(`
			declare type MyNominal<T, Name extends string> = T & {
				specialKey: Name;
			};`),
	}
	if useNoRef {
		files["/home/src/workspaces/solution/tsconfig.json"] = stringtestutil.Dedent(`
		{
			"extends": "./tsconfig.base.json",
			"compilerOptions": { "composite": true },
			"include": ["./src/**/*.ts"],
		}`)
	}
	return files
}

func getTscDeclarationEmitDtsErrorsFileMap(composite bool, incremental bool) FileMap {
	return FileMap{
		"/home/src/workspaces/project/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"module": "NodeNext",
					"moduleResolution": "NodeNext",
					"composite": %t,
					"incremental": %t,
					"declaration": true,
					"skipLibCheck": true,
					"skipDefaultLibCheck": true,
				},
			}`, composite, incremental)),
		"/home/src/workspaces/project/index.ts": stringtestutil.Dedent(`
            import ky from 'ky';
            export const api = ky.extend({});
        `),
		"/home/src/workspaces/project/package.json": stringtestutil.Dedent(`
			{
				"type": "module"
			}`),
		"/home/src/workspaces/project/node_modules/ky/distribution/index.d.ts": stringtestutil.Dedent(`
            type KyInstance = {
                extend(options: Record<string,unknown>): KyInstance;
            }
            declare const ky: KyInstance;
            export default ky;
        `),
		"/home/src/workspaces/project/node_modules/ky/package.json": stringtestutil.Dedent(`
            {
                "name": "ky",
                "type": "module",
                "main": "./distribution/index.js"
            }
        `),
	}
}

func pluginOneConfig() string {
	return stringtestutil.Dedent(`
	{
		"compilerOptions": {
			"target": "es5",
			"declaration": true,
			"traceResolution": true,
		},
	}`)
}

func pluginOneIndex() string {
	return `import pluginTwo from "plugin-two"; // include this to add reference to symlink`
}

func pluginOneAction() string {
	return stringtestutil.Dedent(`
		import { actionCreatorFactory } from "typescript-fsa"; // Include version of shared lib
		const action = actionCreatorFactory("somekey");
		const featureOne = action<{ route: string }>("feature-one");
		export const actions = { featureOne };`)
}

func pluginTwoDts() string {
	return stringtestutil.Dedent(`
		declare const _default: {
			features: {
				featureOne: {
					actions: {
						featureOne: {
							(payload: {
								name: string;
								order: number;
							}, meta?: {
								[key: string]: any;
							}): import("typescript-fsa").Action<{
								name: string;
								order: number;
							}>;
						};
					};
					path: string;
				};
			};
		};
		export default _default;`)
}

func fsaPackageJson() string {
	return stringtestutil.Dedent(`
		{
			"name": "typescript-fsa",
			"version": "3.0.0-beta-2"
		}`)
}

func fsaIndex() string {
	return stringtestutil.Dedent(`
		export interface Action<Payload> {
			type: string;
			payload: Payload;
		}
		export declare type ActionCreator<Payload> = {
			type: string;
			(payload: Payload): Action<Payload>;
		}
		export interface ActionCreatorFactory {
			<Payload = void>(type: string): ActionCreator<Payload>;
		}
		export declare function actionCreatorFactory(prefix?: string | null): ActionCreatorFactory;
		export default actionCreatorFactory;`)
}

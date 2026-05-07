currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/packages/app/node_modules/@lab/feature-gating] -> /home/src/workspaces/project/packages/feature-gating *new*
//// [/home/src/workspaces/project/packages/app/package.json] *new* 
{
    "name": "@lab/app",
    "version": "0.0.0",
    "private": true,
    "type": "module",
    "main": "./dist/index.js",
    "types": "./dist/index.d.ts",
    "dependencies": {
        "@lab/feature-gating": "workspace:*"
    }
}
//// [/home/src/workspaces/project/packages/app/src/index.ts] *new* 
import { createFeatureGateSelector } from "@lab/feature-gating";

export const isFooEnabled = createFeatureGateSelector("foo");
//// [/home/src/workspaces/project/packages/app/tsconfig.json] *new* 
{
    "extends": "../../tsconfig.base.json",
    "compilerOptions": {
        "rootDir": "src",
        "outDir": "dist",
        "tsBuildInfoFile": "dist/.tsbuildinfo"
    },
    "references": [
        { "path": "../feature-gating" }
    ],
    "include": ["src/**/*.ts"]
}
//// [/home/src/workspaces/project/packages/feature-gating/package.json] *new* 
{
    "name": "@lab/feature-gating",
    "version": "0.0.0",
    "private": true,
    "type": "module",
    "main": "./dist/index.js",
    "types": "./dist/index.d.ts",
    "exports": {
        ".": {
            "types": "./dist/index.d.ts",
            "default": "./dist/index.js"
        }
    }
}
//// [/home/src/workspaces/project/packages/feature-gating/src/index.ts] *new* 
import type { State } from "./types.js";

export type Selector<TState, TResult> = (state: TState) => TResult;

export const createFeatureGateSelector =
    (featureGate: string) =>
    (state: State): boolean =>
        state.featureGates[0] === featureGate;
//// [/home/src/workspaces/project/packages/feature-gating/src/types.ts] *new* 
export interface State {
    featureGates: string[];
}
//// [/home/src/workspaces/project/packages/feature-gating/tsconfig.json] *new* 
{
    "extends": "../../tsconfig.base.json",
    "compilerOptions": {
        "rootDir": "src",
        "outDir": "dist",
        "tsBuildInfoFile": "dist/.tsbuildinfo"
    },
    "include": ["src/**/*.ts"]
}
//// [/home/src/workspaces/project/tsconfig.base.json] *new* 
{
    "compilerOptions": {
        "target": "es2022",
        "module": "node16",
        "moduleResolution": "node16",
        "strict": true,
        "declaration": true,
        "composite": true,
        "skipLibCheck": true
    }
}

tsgo -b packages/app/tsconfig.json
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96mpackages/app/src/index.ts[0m:[93m3[0m:[93m14[0m - [91merror[0m[90m TS2883: [0mThe inferred type of 'isFooEnabled' cannot be named without a reference to 'State' from '../node_modules/@lab/feature-gating/src/types.js'. This is likely not portable. A type annotation is necessary.

[7m3[0m export const isFooEnabled = createFeatureGateSelector("foo");
[7m [0m [91m             ~~~~~~~~~~~~[0m


Found 1 error in packages/app/src/index.ts[90m:3[0m

//// [/home/src/tslibs/TS/Lib/lib.es2022.full.d.ts] *Lib*
/// <reference no-default-lib="true"/>
interface Boolean {}
interface Function {}
interface CallableFunction {}
interface NewableFunction {}
interface IArguments {}
interface Number { toExponential: any; }
interface Object {}
interface RegExp {}
interface String { charAt: any; }
interface Array<T> { length: number; [n: number]: T; }
interface ReadonlyArray<T> {}
interface SymbolConstructor {
    (desc?: string | number): symbol;
    for(name: string): symbol;
    readonly toStringTag: symbol;
}
declare var Symbol: SymbolConstructor;
interface Symbol {
    readonly [Symbol.toStringTag]: string;
}
declare const console: { log(msg: any): void; };
//// [/home/src/workspaces/project/packages/app/dist/.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[4],"fileNames":["lib.es2022.full.d.ts","../../feature-gating/dist/types.d.ts","../../feature-gating/dist/index.d.ts","../src/index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n","impliedNodeFormat":99},{"version":"e56efd2ac2313f4027ecef90d5909f7d-import type { State } from \"./types.js\";\nexport type Selector<TState, TResult> = (state: TState) => TResult;\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n","impliedNodeFormat":99},{"version":"298b85581cfc3c43cb5f799407adfdde-import { createFeatureGateSelector } from \"@lab/feature-gating\";\n\nexport const isFooEnabled = createFeatureGateSelector(\"foo\");","signature":"7d83c37faa2e3c49d11d44ba64eb0ff5-export declare const isFooEnabled: any;\n\n(79,12): error2883: The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883\nisFooEnabled\n../node_modules/@lab/feature-gating/src/types.js\nState\n\n(79,12): error2883: The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883\nisFooEnabled\n../node_modules/@lab/feature-gating/src/types.js\nState\n","impliedNodeFormat":99}],"fileIdsList":[[3],[2]],"options":{"composite":true,"declaration":true,"module":100,"outDir":"./","rootDir":"../src","skipLibCheck":true,"strict":true,"target":9,"tsBuildInfoFile":"./.tsbuildinfo"},"referencedMap":[[4,1],[3,2]],"emitDiagnosticsPerFile":[[4,[{"pos":79,"end":91,"code":2883,"category":1,"messageKey":"The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883","messageArgs":["isFooEnabled","../node_modules/@lab/feature-gating/src/types.js","State"]},{"pos":79,"end":91,"code":2883,"category":1,"messageKey":"The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883","messageArgs":["isFooEnabled","../node_modules/@lab/feature-gating/src/types.js","State"]}]]],"latestChangedDtsFile":"./index.d.ts","emitSignatures":[[4,"8ffe87a81ada97360480c6e889089fff-export declare const isFooEnabled: any;\n"]]}
//// [/home/src/workspaces/project/packages/app/dist/.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../src/index.ts"
      ],
      "original": 4
    }
  ],
  "fileNames": [
    "lib.es2022.full.d.ts",
    "../../feature-gating/dist/types.d.ts",
    "../../feature-gating/dist/index.d.ts",
    "../src/index.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2022.full.d.ts",
      "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../feature-gating/dist/types.d.ts",
      "version": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
      "signature": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
        "impliedNodeFormat": 99
      }
    },
    {
      "fileName": "../../feature-gating/dist/index.d.ts",
      "version": "e56efd2ac2313f4027ecef90d5909f7d-import type { State } from \"./types.js\";\nexport type Selector<TState, TResult> = (state: TState) => TResult;\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
      "signature": "e56efd2ac2313f4027ecef90d5909f7d-import type { State } from \"./types.js\";\nexport type Selector<TState, TResult> = (state: TState) => TResult;\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "e56efd2ac2313f4027ecef90d5909f7d-import type { State } from \"./types.js\";\nexport type Selector<TState, TResult> = (state: TState) => TResult;\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
        "impliedNodeFormat": 99
      }
    },
    {
      "fileName": "../src/index.ts",
      "version": "298b85581cfc3c43cb5f799407adfdde-import { createFeatureGateSelector } from \"@lab/feature-gating\";\n\nexport const isFooEnabled = createFeatureGateSelector(\"foo\");",
      "signature": "7d83c37faa2e3c49d11d44ba64eb0ff5-export declare const isFooEnabled: any;\n\n(79,12): error2883: The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883\nisFooEnabled\n../node_modules/@lab/feature-gating/src/types.js\nState\n\n(79,12): error2883: The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883\nisFooEnabled\n../node_modules/@lab/feature-gating/src/types.js\nState\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "298b85581cfc3c43cb5f799407adfdde-import { createFeatureGateSelector } from \"@lab/feature-gating\";\n\nexport const isFooEnabled = createFeatureGateSelector(\"foo\");",
        "signature": "7d83c37faa2e3c49d11d44ba64eb0ff5-export declare const isFooEnabled: any;\n\n(79,12): error2883: The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883\nisFooEnabled\n../node_modules/@lab/feature-gating/src/types.js\nState\n\n(79,12): error2883: The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883\nisFooEnabled\n../node_modules/@lab/feature-gating/src/types.js\nState\n",
        "impliedNodeFormat": 99
      }
    }
  ],
  "fileIdsList": [
    [
      "../../feature-gating/dist/index.d.ts"
    ],
    [
      "../../feature-gating/dist/types.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declaration": true,
    "module": 100,
    "outDir": "./",
    "rootDir": "../src",
    "skipLibCheck": true,
    "strict": true,
    "target": 9,
    "tsBuildInfoFile": "./.tsbuildinfo"
  },
  "referencedMap": {
    "../src/index.ts": [
      "../../feature-gating/dist/index.d.ts"
    ],
    "../../feature-gating/dist/index.d.ts": [
      "../../feature-gating/dist/types.d.ts"
    ]
  },
  "emitDiagnosticsPerFile": [
    [
      "../src/index.ts",
      [
        {
          "pos": 79,
          "end": 91,
          "code": 2883,
          "category": 1,
          "messageKey": "The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883",
          "messageArgs": [
            "isFooEnabled",
            "../node_modules/@lab/feature-gating/src/types.js",
            "State"
          ]
        },
        {
          "pos": 79,
          "end": 91,
          "code": 2883,
          "category": 1,
          "messageKey": "The_inferred_type_of_0_cannot_be_named_without_a_reference_to_2_from_1_This_is_likely_not_portable_A_2883",
          "messageArgs": [
            "isFooEnabled",
            "../node_modules/@lab/feature-gating/src/types.js",
            "State"
          ]
        }
      ]
    ]
  ],
  "latestChangedDtsFile": "./index.d.ts",
  "emitSignatures": [
    {
      "file": "../src/index.ts",
      "signature": "8ffe87a81ada97360480c6e889089fff-export declare const isFooEnabled: any;\n",
      "original": [
        4,
        "8ffe87a81ada97360480c6e889089fff-export declare const isFooEnabled: any;\n"
      ]
    }
  ],
  "size": 2973
}
//// [/home/src/workspaces/project/packages/app/dist/index.d.ts] *new* 
export declare const isFooEnabled: any;

//// [/home/src/workspaces/project/packages/app/dist/index.js] *new* 
import { createFeatureGateSelector } from "@lab/feature-gating";
export const isFooEnabled = createFeatureGateSelector("foo");

//// [/home/src/workspaces/project/packages/feature-gating/dist/.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,3]],"fileNames":["lib.es2022.full.d.ts","../src/types.ts","../src/index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"f78db82625b6990d70f7fe3dce325aef-export interface State {\n    featureGates: string[];\n}","signature":"7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n","impliedNodeFormat":99},{"version":"fb89ef3c627c238871cca890a4b8c2c6-import type { State } from \"./types.js\";\n\nexport type Selector<TState, TResult> = (state: TState) => TResult;\n\nexport const createFeatureGateSelector =\n    (featureGate: string) =>\n    (state: State): boolean =>\n        state.featureGates[0] === featureGate;","signature":"e56efd2ac2313f4027ecef90d5909f7d-import type { State } from \"./types.js\";\nexport type Selector<TState, TResult> = (state: TState) => TResult;\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n","impliedNodeFormat":99}],"fileIdsList":[[2]],"options":{"composite":true,"declaration":true,"module":100,"outDir":"./","rootDir":"../src","skipLibCheck":true,"strict":true,"target":9,"tsBuildInfoFile":"./.tsbuildinfo"},"referencedMap":[[3,1]],"latestChangedDtsFile":"./index.d.ts"}
//// [/home/src/workspaces/project/packages/feature-gating/dist/.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../src/types.ts",
        "../src/index.ts"
      ],
      "original": [
        2,
        3
      ]
    }
  ],
  "fileNames": [
    "lib.es2022.full.d.ts",
    "../src/types.ts",
    "../src/index.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2022.full.d.ts",
      "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/types.ts",
      "version": "f78db82625b6990d70f7fe3dce325aef-export interface State {\n    featureGates: string[];\n}",
      "signature": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "f78db82625b6990d70f7fe3dce325aef-export interface State {\n    featureGates: string[];\n}",
        "signature": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
        "impliedNodeFormat": 99
      }
    },
    {
      "fileName": "../src/index.ts",
      "version": "fb89ef3c627c238871cca890a4b8c2c6-import type { State } from \"./types.js\";\n\nexport type Selector<TState, TResult> = (state: TState) => TResult;\n\nexport const createFeatureGateSelector =\n    (featureGate: string) =>\n    (state: State): boolean =>\n        state.featureGates[0] === featureGate;",
      "signature": "e56efd2ac2313f4027ecef90d5909f7d-import type { State } from \"./types.js\";\nexport type Selector<TState, TResult> = (state: TState) => TResult;\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "fb89ef3c627c238871cca890a4b8c2c6-import type { State } from \"./types.js\";\n\nexport type Selector<TState, TResult> = (state: TState) => TResult;\n\nexport const createFeatureGateSelector =\n    (featureGate: string) =>\n    (state: State): boolean =>\n        state.featureGates[0] === featureGate;",
        "signature": "e56efd2ac2313f4027ecef90d5909f7d-import type { State } from \"./types.js\";\nexport type Selector<TState, TResult> = (state: TState) => TResult;\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
        "impliedNodeFormat": 99
      }
    }
  ],
  "fileIdsList": [
    [
      "../src/types.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declaration": true,
    "module": 100,
    "outDir": "./",
    "rootDir": "../src",
    "skipLibCheck": true,
    "strict": true,
    "target": 9,
    "tsBuildInfoFile": "./.tsbuildinfo"
  },
  "referencedMap": {
    "../src/index.ts": [
      "../src/types.ts"
    ]
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1988
}
//// [/home/src/workspaces/project/packages/feature-gating/dist/index.d.ts] *new* 
import type { State } from "./types.js";
export type Selector<TState, TResult> = (state: TState) => TResult;
export declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;

//// [/home/src/workspaces/project/packages/feature-gating/dist/index.js] *new* 
export const createFeatureGateSelector = (featureGate) => (state) => state.featureGates[0] === featureGate;

//// [/home/src/workspaces/project/packages/feature-gating/dist/types.d.ts] *new* 
export interface State {
    featureGates: string[];
}

//// [/home/src/workspaces/project/packages/feature-gating/dist/types.js] *new* 
export {};


packages/feature-gating/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2022.full.d.ts
*refresh*    /home/src/workspaces/project/packages/feature-gating/src/types.ts
*refresh*    /home/src/workspaces/project/packages/feature-gating/src/index.ts
Signatures::
(stored at emit) /home/src/workspaces/project/packages/feature-gating/src/types.ts
(stored at emit) /home/src/workspaces/project/packages/feature-gating/src/index.ts

packages/app/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2022.full.d.ts
*refresh*    /home/src/workspaces/project/packages/feature-gating/dist/types.d.ts
*refresh*    /home/src/workspaces/project/packages/feature-gating/dist/index.d.ts
*refresh*    /home/src/workspaces/project/packages/app/src/index.ts
Signatures::
(stored at emit) /home/src/workspaces/project/packages/app/src/index.ts

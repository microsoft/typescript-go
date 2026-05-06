currentDirectory::/home/src/workspaces/solution
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/solution/package.json] *new* 
{
    "name": "monorepo-root",
    "private": true,
    "workspaces": ["packages/*"]
}
//// [/home/src/workspaces/solution/packages/app/node_modules/@lab/feature-gating] -> /home/src/workspaces/solution/packages/feature-gating *new*
//// [/home/src/workspaces/solution/packages/app/package.json] *new* 
{
    "name": "@lab/app",
    "version": "1.0.0",
    "private": true,
    "type": "module",
    "main": "./src/index.ts",
    "types": "./src/index.ts",
    "exports": {
        ".": "./src/index.ts"
    },
    "dependencies": {
        "@lab/feature-gating": "workspace:*"
    }
}
//// [/home/src/workspaces/solution/packages/app/src/index.ts] *new* 
import { createFeatureGateSelector } from "@lab/feature-gating";

export const isFooEnabled = createFeatureGateSelector("foo");
//// [/home/src/workspaces/solution/packages/app/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declaration": true,
        "emitDeclarationOnly": true,
        "module": "ESNext",
        "moduleResolution": "Bundler",
        "target": "ES2022",
        "outDir": "./out",
        "rootDir": "./src"
    },
    "include": ["src/**/*"],
    "references": [{ "path": "../feature-gating" }]
}
//// [/home/src/workspaces/solution/packages/feature-gating/package.json] *new* 
{
    "name": "@lab/feature-gating",
    "version": "1.0.0",
    "private": true,
    "type": "module",
    "main": "./src/index.ts",
    "types": "./src/index.ts",
    "exports": {
        ".": "./src/index.ts"
    }
}
//// [/home/src/workspaces/solution/packages/feature-gating/src/index.ts] *new* 
import type { State } from "./types.js";

export const createFeatureGateSelector =
    (featureGate: string) =>
    (state: State): boolean =>
        state.featureGates.includes(featureGate);
//// [/home/src/workspaces/solution/packages/feature-gating/src/types.ts] *new* 
export interface State {
    featureGates: string[];
}
//// [/home/src/workspaces/solution/packages/feature-gating/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declaration": true,
        "emitDeclarationOnly": true,
        "module": "ESNext",
        "moduleResolution": "Bundler",
        "target": "ES2022",
        "outDir": "./out",
        "rootDir": "./src"
    },
    "include": ["src/**/*"]
}
//// [/home/src/workspaces/solution/tsconfig.json] *new* 
{
    "files": [],
    "include": [],
    "references": [
        { "path": "packages/feature-gating" },
        { "path": "packages/app" }
    ]
}

tsgo --b --verbose
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * packages/feature-gating/tsconfig.json
    * packages/app/tsconfig.json
    * tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'packages/feature-gating/tsconfig.json' is out of date because output file 'packages/feature-gating/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'packages/feature-gating/tsconfig.json'...

[96mpackages/feature-gating/src/index.ts[0m:[93m6[0m:[93m28[0m - [91merror[0m[90m TS2550: [0mProperty 'includes' does not exist on type 'string[]'. Do you need to change your target library? Try changing the 'lib' compiler option to 'es2016' or later.

[7m6[0m         state.featureGates.includes(featureGate);
[7m [0m [91m                           ~~~~~~~~[0m

[[90mHH:MM:SS AM[0m] Project 'packages/app/tsconfig.json' is out of date because output file 'packages/app/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'packages/app/tsconfig.json'...


Found 1 error in packages/feature-gating/src/index.ts[90m:6[0m

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
//// [/home/src/workspaces/solution/packages/app/out/index.d.ts] *new* 
export declare const isFooEnabled: (state: import("../../feature-gating/out/types").State) => boolean;

//// [/home/src/workspaces/solution/packages/app/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[4],"fileNames":["lib.es2022.full.d.ts","../feature-gating/out/types.d.ts","../feature-gating/out/index.d.ts","./src/index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n","4c5093097e0b2bf49691652acef51be6-import type { State } from \"./types.js\";\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",{"version":"298b85581cfc3c43cb5f799407adfdde-import { createFeatureGateSelector } from \"@lab/feature-gating\";\n\nexport const isFooEnabled = createFeatureGateSelector(\"foo\");","signature":"a4ed3aed68fab3d71b32e1ca0fc7e0ef-export declare const isFooEnabled: (state: import(\"../../feature-gating/out/types\").State) => boolean;\n","impliedNodeFormat":1}],"fileIdsList":[[3],[2]],"options":{"composite":true,"emitDeclarationOnly":true,"declaration":true,"module":99,"outDir":"./out","rootDir":"./src","target":9},"referencedMap":[[4,1],[3,2]],"latestChangedDtsFile":"./out/index.d.ts"}
//// [/home/src/workspaces/solution/packages/app/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./src/index.ts"
      ],
      "original": 4
    }
  ],
  "fileNames": [
    "lib.es2022.full.d.ts",
    "../feature-gating/out/types.d.ts",
    "../feature-gating/out/index.d.ts",
    "./src/index.ts"
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
      "fileName": "../feature-gating/out/types.d.ts",
      "version": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
      "signature": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../feature-gating/out/index.d.ts",
      "version": "4c5093097e0b2bf49691652acef51be6-import type { State } from \"./types.js\";\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
      "signature": "4c5093097e0b2bf49691652acef51be6-import type { State } from \"./types.js\";\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/index.ts",
      "version": "298b85581cfc3c43cb5f799407adfdde-import { createFeatureGateSelector } from \"@lab/feature-gating\";\n\nexport const isFooEnabled = createFeatureGateSelector(\"foo\");",
      "signature": "a4ed3aed68fab3d71b32e1ca0fc7e0ef-export declare const isFooEnabled: (state: import(\"../../feature-gating/out/types\").State) => boolean;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "298b85581cfc3c43cb5f799407adfdde-import { createFeatureGateSelector } from \"@lab/feature-gating\";\n\nexport const isFooEnabled = createFeatureGateSelector(\"foo\");",
        "signature": "a4ed3aed68fab3d71b32e1ca0fc7e0ef-export declare const isFooEnabled: (state: import(\"../../feature-gating/out/types\").State) => boolean;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../feature-gating/out/index.d.ts"
    ],
    [
      "../feature-gating/out/types.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "emitDeclarationOnly": true,
    "declaration": true,
    "module": 99,
    "outDir": "./out",
    "rootDir": "./src",
    "target": 9
  },
  "referencedMap": {
    "./src/index.ts": [
      "../feature-gating/out/index.d.ts"
    ],
    "../feature-gating/out/index.d.ts": [
      "../feature-gating/out/types.d.ts"
    ]
  },
  "latestChangedDtsFile": "./out/index.d.ts",
  "size": 1807
}
//// [/home/src/workspaces/solution/packages/feature-gating/out/index.d.ts] *new* 
import type { State } from "./types.js";
export declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;

//// [/home/src/workspaces/solution/packages/feature-gating/out/types.d.ts] *new* 
export interface State {
    featureGates: string[];
}

//// [/home/src/workspaces/solution/packages/feature-gating/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,3]],"fileNames":["lib.es2022.full.d.ts","./src/types.ts","./src/index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"f78db82625b6990d70f7fe3dce325aef-export interface State {\n    featureGates: string[];\n}","signature":"7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n","impliedNodeFormat":1},{"version":"0d68b3494c1f711ae216bf424b83044f-import type { State } from \"./types.js\";\n\nexport const createFeatureGateSelector =\n    (featureGate: string) =>\n    (state: State): boolean =>\n        state.featureGates.includes(featureGate);","signature":"4c5093097e0b2bf49691652acef51be6-import type { State } from \"./types.js\";\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"composite":true,"emitDeclarationOnly":true,"declaration":true,"module":99,"outDir":"./out","rootDir":"./src","target":9},"referencedMap":[[3,1]],"semanticDiagnosticsPerFile":[[3,[{"pos":170,"end":178,"code":2550,"category":1,"messageKey":"Property_0_does_not_exist_on_type_1_Do_you_need_to_change_your_target_library_Try_changing_the_lib_c_2550","messageArgs":["includes","string[]","es2016"]}]]],"latestChangedDtsFile":"./out/index.d.ts"}
//// [/home/src/workspaces/solution/packages/feature-gating/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./src/types.ts",
        "./src/index.ts"
      ],
      "original": [
        2,
        3
      ]
    }
  ],
  "fileNames": [
    "lib.es2022.full.d.ts",
    "./src/types.ts",
    "./src/index.ts"
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
      "fileName": "./src/types.ts",
      "version": "f78db82625b6990d70f7fe3dce325aef-export interface State {\n    featureGates: string[];\n}",
      "signature": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f78db82625b6990d70f7fe3dce325aef-export interface State {\n    featureGates: string[];\n}",
        "signature": "7d532181b5d800e5494f045792170875-export interface State {\n    featureGates: string[];\n}\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/index.ts",
      "version": "0d68b3494c1f711ae216bf424b83044f-import type { State } from \"./types.js\";\n\nexport const createFeatureGateSelector =\n    (featureGate: string) =>\n    (state: State): boolean =>\n        state.featureGates.includes(featureGate);",
      "signature": "4c5093097e0b2bf49691652acef51be6-import type { State } from \"./types.js\";\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0d68b3494c1f711ae216bf424b83044f-import type { State } from \"./types.js\";\n\nexport const createFeatureGateSelector =\n    (featureGate: string) =>\n    (state: State): boolean =>\n        state.featureGates.includes(featureGate);",
        "signature": "4c5093097e0b2bf49691652acef51be6-import type { State } from \"./types.js\";\nexport declare const createFeatureGateSelector: (featureGate: string) => (state: State) => boolean;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/types.ts"
    ]
  ],
  "options": {
    "composite": true,
    "emitDeclarationOnly": true,
    "declaration": true,
    "module": 99,
    "outDir": "./out",
    "rootDir": "./src",
    "target": 9
  },
  "referencedMap": {
    "./src/index.ts": [
      "./src/types.ts"
    ]
  },
  "semanticDiagnosticsPerFile": [
    [
      "./src/index.ts",
      [
        {
          "pos": 170,
          "end": 178,
          "code": 2550,
          "category": 1,
          "messageKey": "Property_0_does_not_exist_on_type_1_Do_you_need_to_change_your_target_library_Try_changing_the_lib_c_2550",
          "messageArgs": [
            "includes",
            "string[]",
            "es2016"
          ]
        }
      ]
    ]
  ],
  "latestChangedDtsFile": "./out/index.d.ts",
  "size": 2062
}

packages/feature-gating/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2022.full.d.ts
*refresh*    /home/src/workspaces/solution/packages/feature-gating/src/types.ts
*refresh*    /home/src/workspaces/solution/packages/feature-gating/src/index.ts
Signatures::
(stored at emit) /home/src/workspaces/solution/packages/feature-gating/src/types.ts
(stored at emit) /home/src/workspaces/solution/packages/feature-gating/src/index.ts

packages/app/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2022.full.d.ts
*refresh*    /home/src/workspaces/solution/packages/feature-gating/out/types.d.ts
*refresh*    /home/src/workspaces/solution/packages/feature-gating/out/index.d.ts
*refresh*    /home/src/workspaces/solution/packages/app/src/index.ts
Signatures::
(stored at emit) /home/src/workspaces/solution/packages/app/src/index.ts

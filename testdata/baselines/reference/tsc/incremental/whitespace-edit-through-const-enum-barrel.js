currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/src/core-barrel.ts] *new* 
export * from './core';

//// [/home/src/workspaces/project/src/core.ts] *new* 
export interface CoreThing { a: number }
export const CORE = 1;
//// [/home/src/workspaces/project/src/index.ts] *new* 
export * from './mid';

//// [/home/src/workspaces/project/src/leaf0.ts] *new* 
import { MID } from './index';
export const L0: number = MID + 0;
//// [/home/src/workspaces/project/src/leaf1.ts] *new* 
import { MID } from './index';
export const L1: number = MID + 1;
//// [/home/src/workspaces/project/src/leaf2.ts] *new* 
import { MID } from './index';
export const L2: number = MID + 2;
//// [/home/src/workspaces/project/src/mid.ts] *new* 
import { CORE, CoreThing } from './core-barrel';
export const MID: number = CORE + 1;
export function useCore(x: CoreThing): number { return x.a + MID; }
export const enum MidMode { A = 1, B = 2 }
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "incremental": true,
        "module": "esnext",
        "moduleResolution": "bundler",
        "strict": true,
    },
    "include": ["src/**/*.ts"],
}

tsgo 
ExitStatus:: Success
Output::
//// [/home/src/tslibs/TS/Lib/lib.es2025.full.d.ts] *Lib*
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
//// [/home/src/workspaces/project/src/core-barrel.js] *new* 
export * from './core';

//// [/home/src/workspaces/project/src/core.js] *new* 
export const CORE = 1;

//// [/home/src/workspaces/project/src/index.js] *new* 
export * from './mid';

//// [/home/src/workspaces/project/src/leaf0.js] *new* 
import { MID } from './index';
export const L0 = MID + 0;

//// [/home/src/workspaces/project/src/leaf1.js] *new* 
import { MID } from './index';
export const L1 = MID + 1;

//// [/home/src/workspaces/project/src/leaf2.js] *new* 
import { MID } from './index';
export const L2 = MID + 2;

//// [/home/src/workspaces/project/src/mid.js] *new* 
import { CORE } from './core-barrel';
export const MID = CORE + 1;
export function useCore(x) { return x.a + MID; }

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,8]],"fileNames":["lib.es2025.full.d.ts","./src/core.ts","./src/core-barrel.ts","./src/mid.ts","./src/index.ts","./src/leaf0.ts","./src/leaf1.ts","./src/leaf2.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;","signature":"ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n","impliedNodeFormat":1},"9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",{"version":"fa7b0a93637d28d11404a7083d543bb7-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }","signature":"3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n","impliedNodeFormat":1},"d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",{"version":"08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;","signature":"a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n","impliedNodeFormat":1},{"version":"19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;","signature":"4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n","impliedNodeFormat":1},{"version":"fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;","signature":"e0004efaf732955da575bd44946afa47-export declare const L2: number;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[4],[5],[3]],"options":{"module":99,"strict":true},"referencedMap":[[3,1],[5,2],[6,3],[7,3],[8,3],[4,4]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./src/core.ts",
        "./src/core-barrel.ts",
        "./src/mid.ts",
        "./src/index.ts",
        "./src/leaf0.ts",
        "./src/leaf1.ts",
        "./src/leaf2.ts"
      ],
      "original": [
        2,
        8
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./src/core.ts",
    "./src/core-barrel.ts",
    "./src/mid.ts",
    "./src/index.ts",
    "./src/leaf0.ts",
    "./src/leaf1.ts",
    "./src/leaf2.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2025.full.d.ts",
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
      "fileName": "./src/core.ts",
      "version": "e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;",
      "signature": "ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;",
        "signature": "ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/core-barrel.ts",
      "version": "9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",
      "signature": "9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/mid.ts",
      "version": "fa7b0a93637d28d11404a7083d543bb7-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }",
      "signature": "3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "fa7b0a93637d28d11404a7083d543bb7-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }",
        "signature": "3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/index.ts",
      "version": "d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",
      "signature": "d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/leaf0.ts",
      "version": "08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;",
      "signature": "a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;",
        "signature": "a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/leaf1.ts",
      "version": "19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;",
      "signature": "4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;",
        "signature": "4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/leaf2.ts",
      "version": "fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;",
      "signature": "e0004efaf732955da575bd44946afa47-export declare const L2: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;",
        "signature": "e0004efaf732955da575bd44946afa47-export declare const L2: number;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/core.ts"
    ],
    [
      "./src/mid.ts"
    ],
    [
      "./src/index.ts"
    ],
    [
      "./src/core-barrel.ts"
    ]
  ],
  "options": {
    "module": 99,
    "strict": true
  },
  "referencedMap": {
    "./src/core-barrel.ts": [
      "./src/core.ts"
    ],
    "./src/index.ts": [
      "./src/mid.ts"
    ],
    "./src/leaf0.ts": [
      "./src/index.ts"
    ],
    "./src/leaf1.ts": [
      "./src/index.ts"
    ],
    "./src/leaf2.ts": [
      "./src/index.ts"
    ],
    "./src/mid.ts": [
      "./src/core-barrel.ts"
    ]
  },
  "size": 2662
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/src/core.ts
*refresh*    /home/src/workspaces/project/src/core-barrel.ts
*refresh*    /home/src/workspaces/project/src/mid.ts
*refresh*    /home/src/workspaces/project/src/index.ts
*refresh*    /home/src/workspaces/project/src/leaf0.ts
*refresh*    /home/src/workspaces/project/src/leaf1.ts
*refresh*    /home/src/workspaces/project/src/leaf2.ts
Signatures::


Edit [0]:: whitespace only edit of mid.ts
//// [/home/src/workspaces/project/src/mid.ts] *modified* 
import { CORE, CoreThing } from './core-barrel';
export const MID: number = CORE + 1;
export function useCore(x: CoreThing): number { return x.a + MID; }
export const enum MidMode { A = 1, B = 2 }


tsgo 
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/src/mid.js] *rewrite with same content*
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,8]],"fileNames":["lib.es2025.full.d.ts","./src/core.ts","./src/core-barrel.ts","./src/mid.ts","./src/index.ts","./src/leaf0.ts","./src/leaf1.ts","./src/leaf2.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;","signature":"ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n","impliedNodeFormat":1},"9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",{"version":"8fe87069c86f2dc1185a7c3fcbbafdff-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }\n","signature":"3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n","impliedNodeFormat":1},"d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",{"version":"08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;","signature":"a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n","impliedNodeFormat":1},{"version":"19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;","signature":"4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n","impliedNodeFormat":1},{"version":"fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;","signature":"e0004efaf732955da575bd44946afa47-export declare const L2: number;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[4],[5],[3]],"options":{"module":99,"strict":true},"referencedMap":[[3,1],[5,2],[6,3],[7,3],[8,3],[4,4]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./src/core.ts",
        "./src/core-barrel.ts",
        "./src/mid.ts",
        "./src/index.ts",
        "./src/leaf0.ts",
        "./src/leaf1.ts",
        "./src/leaf2.ts"
      ],
      "original": [
        2,
        8
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./src/core.ts",
    "./src/core-barrel.ts",
    "./src/mid.ts",
    "./src/index.ts",
    "./src/leaf0.ts",
    "./src/leaf1.ts",
    "./src/leaf2.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2025.full.d.ts",
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
      "fileName": "./src/core.ts",
      "version": "e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;",
      "signature": "ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;",
        "signature": "ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/core-barrel.ts",
      "version": "9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",
      "signature": "9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/mid.ts",
      "version": "8fe87069c86f2dc1185a7c3fcbbafdff-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }\n",
      "signature": "3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8fe87069c86f2dc1185a7c3fcbbafdff-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }\n",
        "signature": "3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/index.ts",
      "version": "d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",
      "signature": "d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/leaf0.ts",
      "version": "08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;",
      "signature": "a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;",
        "signature": "a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/leaf1.ts",
      "version": "19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;",
      "signature": "4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;",
        "signature": "4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/leaf2.ts",
      "version": "fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;",
      "signature": "e0004efaf732955da575bd44946afa47-export declare const L2: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;",
        "signature": "e0004efaf732955da575bd44946afa47-export declare const L2: number;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/core.ts"
    ],
    [
      "./src/mid.ts"
    ],
    [
      "./src/index.ts"
    ],
    [
      "./src/core-barrel.ts"
    ]
  ],
  "options": {
    "module": 99,
    "strict": true
  },
  "referencedMap": {
    "./src/core-barrel.ts": [
      "./src/core.ts"
    ],
    "./src/index.ts": [
      "./src/mid.ts"
    ],
    "./src/leaf0.ts": [
      "./src/index.ts"
    ],
    "./src/leaf1.ts": [
      "./src/index.ts"
    ],
    "./src/leaf2.ts": [
      "./src/index.ts"
    ],
    "./src/mid.ts": [
      "./src/core-barrel.ts"
    ]
  },
  "size": 2664
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/mid.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/src/mid.ts


Edit [1]:: same whitespace only edit of mid.ts again
//// [/home/src/workspaces/project/src/mid.ts] *modified* 
import { CORE, CoreThing } from './core-barrel';
export const MID: number = CORE + 1;
export function useCore(x: CoreThing): number { return x.a + MID; }
export const enum MidMode { A = 1, B = 2 }



tsgo 
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/src/mid.js] *rewrite with same content*
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,8]],"fileNames":["lib.es2025.full.d.ts","./src/core.ts","./src/core-barrel.ts","./src/mid.ts","./src/index.ts","./src/leaf0.ts","./src/leaf1.ts","./src/leaf2.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;","signature":"ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n","impliedNodeFormat":1},"9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",{"version":"6c1cd5f0aafc87eb884063d62bb4f3eb-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }\n\n","signature":"3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n","impliedNodeFormat":1},"d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",{"version":"08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;","signature":"a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n","impliedNodeFormat":1},{"version":"19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;","signature":"4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n","impliedNodeFormat":1},{"version":"fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;","signature":"e0004efaf732955da575bd44946afa47-export declare const L2: number;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[4],[5],[3]],"options":{"module":99,"strict":true},"referencedMap":[[3,1],[5,2],[6,3],[7,3],[8,3],[4,4]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./src/core.ts",
        "./src/core-barrel.ts",
        "./src/mid.ts",
        "./src/index.ts",
        "./src/leaf0.ts",
        "./src/leaf1.ts",
        "./src/leaf2.ts"
      ],
      "original": [
        2,
        8
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./src/core.ts",
    "./src/core-barrel.ts",
    "./src/mid.ts",
    "./src/index.ts",
    "./src/leaf0.ts",
    "./src/leaf1.ts",
    "./src/leaf2.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2025.full.d.ts",
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
      "fileName": "./src/core.ts",
      "version": "e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;",
      "signature": "ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e5491f3dd6b53b47da4e4fb95c9a7b29-export interface CoreThing { a: number }\nexport const CORE = 1;",
        "signature": "ba4df4ecef1094d4eac01d5521f11cf1-export interface CoreThing {\n    a: number;\n}\nexport declare const CORE = 1;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/core-barrel.ts",
      "version": "9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",
      "signature": "9052442094c7d6b8a0042b8cf78b0d7c-export * from './core';\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/mid.ts",
      "version": "6c1cd5f0aafc87eb884063d62bb4f3eb-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }\n\n",
      "signature": "3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "6c1cd5f0aafc87eb884063d62bb4f3eb-import { CORE, CoreThing } from './core-barrel';\nexport const MID: number = CORE + 1;\nexport function useCore(x: CoreThing): number { return x.a + MID; }\nexport const enum MidMode { A = 1, B = 2 }\n\n",
        "signature": "3a937d8c1f5eac7c2ccdc796215b98e2-import { CoreThing } from './core-barrel';\nexport declare const MID: number;\nexport declare function useCore(x: CoreThing): number;\nexport declare const enum MidMode {\n    A = 1,\n    B = 2\n}\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/index.ts",
      "version": "d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",
      "signature": "d29114db824119c790ab0abbfa71ed26-export * from './mid';\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/leaf0.ts",
      "version": "08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;",
      "signature": "a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "08ba7e0251a0c703fdbcaeb58aff98a4-import { MID } from './index';\nexport const L0: number = MID + 0;",
        "signature": "a75cb491cae2980bea107af43ccdb683-export declare const L0: number;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/leaf1.ts",
      "version": "19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;",
      "signature": "4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "19a6135ac807112d3f56bffe4c701da0-import { MID } from './index';\nexport const L1: number = MID + 1;",
        "signature": "4436bfc413b2bf1223720335abbfdc92-export declare const L1: number;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/leaf2.ts",
      "version": "fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;",
      "signature": "e0004efaf732955da575bd44946afa47-export declare const L2: number;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "fc85df162774370c1ab0da8fda8dcb32-import { MID } from './index';\nexport const L2: number = MID + 2;",
        "signature": "e0004efaf732955da575bd44946afa47-export declare const L2: number;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/core.ts"
    ],
    [
      "./src/mid.ts"
    ],
    [
      "./src/index.ts"
    ],
    [
      "./src/core-barrel.ts"
    ]
  ],
  "options": {
    "module": 99,
    "strict": true
  },
  "referencedMap": {
    "./src/core-barrel.ts": [
      "./src/core.ts"
    ],
    "./src/index.ts": [
      "./src/mid.ts"
    ],
    "./src/leaf0.ts": [
      "./src/index.ts"
    ],
    "./src/leaf1.ts": [
      "./src/index.ts"
    ],
    "./src/leaf2.ts": [
      "./src/index.ts"
    ],
    "./src/mid.ts": [
      "./src/core-barrel.ts"
    ]
  },
  "size": 2666
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/mid.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/src/mid.ts

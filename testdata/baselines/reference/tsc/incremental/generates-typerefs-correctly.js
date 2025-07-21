currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/src/box.ts] *new* 
export interface Box<T> {
    unbox(): T
}
//// [/home/src/workspaces/project/src/bug.js] *new* 
import * as B from "./box.js"
import * as W from "./wrap.js"

/**
 * @template {object} C
 * @param {C} source
 * @returns {W.Wrap<C>}
 */
const wrap = source => {
throw source
}

/**
 * @returns {B.Box<number>}
 */
const box = (n = 0) => ({ unbox: () => n })

export const bug = wrap({ n: box(1) });
//// [/home/src/workspaces/project/src/wrap.ts] *new* 
export type Wrap<C> = {
    [K in keyof C]: { wrapped: C[K] }
}
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "outDir": "outDir",
        "checkJs": true
    },
    "include": ["src"],
}

tsgo 
ExitStatus:: Success
Output::
//// [/home/src/tslibs/TS/Lib/lib.d.ts] *Lib*
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
//// [/home/src/workspaces/project/outDir/src/box.d.ts] *new* 
export interface Box<T> {
    unbox(): T;
}

//// [/home/src/workspaces/project/outDir/src/box.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });

//// [/home/src/workspaces/project/outDir/src/bug.d.ts] *new* 
import * as B from "./box.js";
import * as W from "./wrap.js";
export declare const bug: W.Wrap<{
    n: B.Box<number>;
}>;

//// [/home/src/workspaces/project/outDir/src/bug.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.bug = void 0;
const B = require("./box.js");
const W = require("./wrap.js");
/**
 * @template {object} C
 * @param {C} source
 * @returns {W.Wrap<C>}
 */
const wrap = source => {
    throw source;
};
/**
 * @returns {B.Box<number>}
 */
const box = (n = 0) => ({ unbox: () => n });
exports.bug = wrap({ n: box(1) });

//// [/home/src/workspaces/project/outDir/src/wrap.d.ts] *new* 
export type Wrap<C> = {
    [K in keyof C]: {
        wrapped: C[K];
    };
};

//// [/home/src/workspaces/project/outDir/src/wrap.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });

//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../../tslibs/TS/Lib/lib.d.ts","../src/box.ts","../src/wrap.ts","../src/bug.js"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e-export interface Box\u003cT\u003e {\n    unbox(): T\n}","signature":"5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448-export interface Box\u003cT\u003e {\n    unbox(): T;\n}\n","impliedNodeFormat":1},{"version":"a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: { wrapped: C[K] }\n}","signature":"4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: {\n        wrapped: C[K];\n    };\n};\n","impliedNodeFormat":1},{"version":"f4375ea4700b8e5e8d2c1cc7552de717928f9b5333a1be759bc15f196daea88e-import * as B from \"./box.js\"\nimport * as W from \"./wrap.js\"\n\n/**\n * @template {object} C\n * @param {C} source\n * @returns {W.Wrap\u003cC\u003e}\n */\nconst wrap = source =\u003e {\nthrow source\n}\n\n/**\n * @returns {B.Box\u003cnumber\u003e}\n */\nconst box = (n = 0) =\u003e ({ unbox: () =\u003e n })\n\nexport const bug = wrap({ n: box(1) });","signature":"d50883e16164d399f0e9534640b44ee6c05a7c690534d0201a645262e7caab22-import * as B from \"./box.js\";\nimport * as W from \"./wrap.js\";\nexport declare const bug: W.Wrap\u003c{\n    n: B.Box\u003cnumber\u003e;\n}\u003e;\n","impliedNodeFormat":1}],"fileIdsList":[[2,3]],"options":{"checkJs":true,"composite":true,"outDir":"./"},"referencedMap":[[4,1]],"latestChangedDtsFile":"./src/wrap.d.ts"}
//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../../tslibs/TS/Lib/lib.d.ts",
    "../src/box.ts",
    "../src/wrap.ts",
    "../src/bug.js"
  ],
  "fileInfos": [
    {
      "fileName": "../../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/box.ts",
      "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e-export interface Box\u003cT\u003e {\n    unbox(): T\n}",
      "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448-export interface Box\u003cT\u003e {\n    unbox(): T;\n}\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e-export interface Box\u003cT\u003e {\n    unbox(): T\n}",
        "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448-export interface Box\u003cT\u003e {\n    unbox(): T;\n}\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/wrap.ts",
      "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: { wrapped: C[K] }\n}",
      "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: {\n        wrapped: C[K];\n    };\n};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: { wrapped: C[K] }\n}",
        "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: {\n        wrapped: C[K];\n    };\n};\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/bug.js",
      "version": "f4375ea4700b8e5e8d2c1cc7552de717928f9b5333a1be759bc15f196daea88e-import * as B from \"./box.js\"\nimport * as W from \"./wrap.js\"\n\n/**\n * @template {object} C\n * @param {C} source\n * @returns {W.Wrap\u003cC\u003e}\n */\nconst wrap = source =\u003e {\nthrow source\n}\n\n/**\n * @returns {B.Box\u003cnumber\u003e}\n */\nconst box = (n = 0) =\u003e ({ unbox: () =\u003e n })\n\nexport const bug = wrap({ n: box(1) });",
      "signature": "d50883e16164d399f0e9534640b44ee6c05a7c690534d0201a645262e7caab22-import * as B from \"./box.js\";\nimport * as W from \"./wrap.js\";\nexport declare const bug: W.Wrap\u003c{\n    n: B.Box\u003cnumber\u003e;\n}\u003e;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f4375ea4700b8e5e8d2c1cc7552de717928f9b5333a1be759bc15f196daea88e-import * as B from \"./box.js\"\nimport * as W from \"./wrap.js\"\n\n/**\n * @template {object} C\n * @param {C} source\n * @returns {W.Wrap\u003cC\u003e}\n */\nconst wrap = source =\u003e {\nthrow source\n}\n\n/**\n * @returns {B.Box\u003cnumber\u003e}\n */\nconst box = (n = 0) =\u003e ({ unbox: () =\u003e n })\n\nexport const bug = wrap({ n: box(1) });",
        "signature": "d50883e16164d399f0e9534640b44ee6c05a7c690534d0201a645262e7caab22-import * as B from \"./box.js\";\nimport * as W from \"./wrap.js\";\nexport declare const bug: W.Wrap\u003c{\n    n: B.Box\u003cnumber\u003e;\n}\u003e;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../src/box.ts",
      "../src/wrap.ts"
    ]
  ],
  "options": {
    "checkJs": true,
    "composite": true,
    "outDir": "./"
  },
  "referencedMap": {
    "../src/bug.js": [
      "../src/box.ts",
      "../src/wrap.ts"
    ]
  },
  "latestChangedDtsFile": "./src/wrap.d.ts",
  "size": 2450
}

SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/src/box.ts
*refresh*    /home/src/workspaces/project/src/wrap.ts
*refresh*    /home/src/workspaces/project/src/bug.js
Signatures::
(stored at emit) /home/src/workspaces/project/src/box.ts
(stored at emit) /home/src/workspaces/project/src/wrap.ts
(stored at emit) /home/src/workspaces/project/src/bug.js


Edit [0]:: modify js file
//// [/home/src/workspaces/project/src/bug.js] *modified* 
import * as B from "./box.js"
import * as W from "./wrap.js"

/**
 * @template {object} C
 * @param {C} source
 * @returns {W.Wrap<C>}
 */
const wrap = source => {
throw source
}

/**
 * @returns {B.Box<number>}
 */
const box = (n = 0) => ({ unbox: () => n })

export const bug = wrap({ n: box(1) });export const something = 1;

tsgo 
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/outDir/src/bug.d.ts] *modified* 
import * as B from "./box.js";
import * as W from "./wrap.js";
export declare const bug: W.Wrap<{
    n: B.Box<number>;
}>;
export declare const something = 1;

//// [/home/src/workspaces/project/outDir/src/bug.js] *modified* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.something = exports.bug = void 0;
const B = require("./box.js");
const W = require("./wrap.js");
/**
 * @template {object} C
 * @param {C} source
 * @returns {W.Wrap<C>}
 */
const wrap = source => {
    throw source;
};
/**
 * @returns {B.Box<number>}
 */
const box = (n = 0) => ({ unbox: () => n });
exports.bug = wrap({ n: box(1) });
exports.something = 1;

//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../../tslibs/TS/Lib/lib.d.ts","../src/box.ts","../src/wrap.ts","../src/bug.js"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e-export interface Box\u003cT\u003e {\n    unbox(): T\n}","signature":"5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448-export interface Box\u003cT\u003e {\n    unbox(): T;\n}\n","impliedNodeFormat":1},{"version":"a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: { wrapped: C[K] }\n}","signature":"4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: {\n        wrapped: C[K];\n    };\n};\n","impliedNodeFormat":1},{"version":"676e1613c74cee0aa880dc1108ce9e0a07e184e35a7d78282bf4a480bf6124db-import * as B from \"./box.js\"\nimport * as W from \"./wrap.js\"\n\n/**\n * @template {object} C\n * @param {C} source\n * @returns {W.Wrap\u003cC\u003e}\n */\nconst wrap = source =\u003e {\nthrow source\n}\n\n/**\n * @returns {B.Box\u003cnumber\u003e}\n */\nconst box = (n = 0) =\u003e ({ unbox: () =\u003e n })\n\nexport const bug = wrap({ n: box(1) });export const something = 1;","signature":"6a80e6a4e41ee496beaac682ed90a523d5f8b87987da12ca01326156d552d025-import * as B from \"./box.js\";\nimport * as W from \"./wrap.js\";\nexport declare const bug: W.Wrap\u003c{\n    n: B.Box\u003cnumber\u003e;\n}\u003e;\nexport declare const something = 1;\n","impliedNodeFormat":1}],"fileIdsList":[[2,3]],"options":{"checkJs":true,"composite":true,"outDir":"./"},"referencedMap":[[4,1]],"latestChangedDtsFile":"./src/bug.d.ts"}
//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../../tslibs/TS/Lib/lib.d.ts",
    "../src/box.ts",
    "../src/wrap.ts",
    "../src/bug.js"
  ],
  "fileInfos": [
    {
      "fileName": "../../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/box.ts",
      "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e-export interface Box\u003cT\u003e {\n    unbox(): T\n}",
      "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448-export interface Box\u003cT\u003e {\n    unbox(): T;\n}\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e-export interface Box\u003cT\u003e {\n    unbox(): T\n}",
        "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448-export interface Box\u003cT\u003e {\n    unbox(): T;\n}\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/wrap.ts",
      "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: { wrapped: C[K] }\n}",
      "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: {\n        wrapped: C[K];\n    };\n};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: { wrapped: C[K] }\n}",
        "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b-export type Wrap\u003cC\u003e = {\n    [K in keyof C]: {\n        wrapped: C[K];\n    };\n};\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/bug.js",
      "version": "676e1613c74cee0aa880dc1108ce9e0a07e184e35a7d78282bf4a480bf6124db-import * as B from \"./box.js\"\nimport * as W from \"./wrap.js\"\n\n/**\n * @template {object} C\n * @param {C} source\n * @returns {W.Wrap\u003cC\u003e}\n */\nconst wrap = source =\u003e {\nthrow source\n}\n\n/**\n * @returns {B.Box\u003cnumber\u003e}\n */\nconst box = (n = 0) =\u003e ({ unbox: () =\u003e n })\n\nexport const bug = wrap({ n: box(1) });export const something = 1;",
      "signature": "6a80e6a4e41ee496beaac682ed90a523d5f8b87987da12ca01326156d552d025-import * as B from \"./box.js\";\nimport * as W from \"./wrap.js\";\nexport declare const bug: W.Wrap\u003c{\n    n: B.Box\u003cnumber\u003e;\n}\u003e;\nexport declare const something = 1;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "676e1613c74cee0aa880dc1108ce9e0a07e184e35a7d78282bf4a480bf6124db-import * as B from \"./box.js\"\nimport * as W from \"./wrap.js\"\n\n/**\n * @template {object} C\n * @param {C} source\n * @returns {W.Wrap\u003cC\u003e}\n */\nconst wrap = source =\u003e {\nthrow source\n}\n\n/**\n * @returns {B.Box\u003cnumber\u003e}\n */\nconst box = (n = 0) =\u003e ({ unbox: () =\u003e n })\n\nexport const bug = wrap({ n: box(1) });export const something = 1;",
        "signature": "6a80e6a4e41ee496beaac682ed90a523d5f8b87987da12ca01326156d552d025-import * as B from \"./box.js\";\nimport * as W from \"./wrap.js\";\nexport declare const bug: W.Wrap\u003c{\n    n: B.Box\u003cnumber\u003e;\n}\u003e;\nexport declare const something = 1;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../src/box.ts",
      "../src/wrap.ts"
    ]
  ],
  "options": {
    "checkJs": true,
    "composite": true,
    "outDir": "./"
  },
  "referencedMap": {
    "../src/bug.js": [
      "../src/box.ts",
      "../src/wrap.ts"
    ]
  },
  "latestChangedDtsFile": "./src/bug.d.ts",
  "size": 2513
}

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/bug.js
Signatures::
(computed .d.ts) /home/src/workspaces/project/src/bug.js

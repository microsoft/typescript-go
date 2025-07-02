
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

ExitStatus:: 0

CompilerOptions::{}
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
{"version":"FakeTSVersion","fileNames":["../../../tslibs/TS/Lib/lib.d.ts","../src/box.ts","../src/wrap.ts","../src/bug.js"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e","signature":"5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448","impliedNodeFormat":1},{"version":"a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7","signature":"4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b","impliedNodeFormat":1},{"version":"f4375ea4700b8e5e8d2c1cc7552de717928f9b5333a1be759bc15f196daea88e","signature":"d50883e16164d399f0e9534640b44ee6c05a7c690534d0201a645262e7caab22","impliedNodeFormat":1}],"fileIdsList":[[2,3]],"options":{"checkJs":true,"composite":true,"outDir":"./"},"referencedMap":[[4,1]],"latestChangedDtsFile":"./src/wrap.d.ts"}
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
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/box.ts",
      "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e",
      "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e",
        "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/wrap.ts",
      "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7",
      "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7",
        "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/bug.js",
      "version": "f4375ea4700b8e5e8d2c1cc7552de717928f9b5333a1be759bc15f196daea88e",
      "signature": "d50883e16164d399f0e9534640b44ee6c05a7c690534d0201a645262e7caab22",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f4375ea4700b8e5e8d2c1cc7552de717928f9b5333a1be759bc15f196daea88e",
        "signature": "d50883e16164d399f0e9534640b44ee6c05a7c690534d0201a645262e7caab22",
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
  "size": 950
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


Edit:: modify js file
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

ExitStatus:: 0
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
{"version":"FakeTSVersion","fileNames":["../../../tslibs/TS/Lib/lib.d.ts","../src/box.ts","../src/wrap.ts","../src/bug.js"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e","signature":"5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448","impliedNodeFormat":1},{"version":"a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7","signature":"4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b","impliedNodeFormat":1},{"version":"676e1613c74cee0aa880dc1108ce9e0a07e184e35a7d78282bf4a480bf6124db","signature":"6a80e6a4e41ee496beaac682ed90a523d5f8b87987da12ca01326156d552d025","impliedNodeFormat":1}],"fileIdsList":[[2,3]],"options":{"checkJs":true,"composite":true,"outDir":"./"},"referencedMap":[[4,1]],"latestChangedDtsFile":"./src/bug.d.ts"}
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
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/box.ts",
      "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e",
      "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0368ea83346f71eb3a3ca19fa727980a0456898b66b88a4c998681504e53158e",
        "signature": "5869e012eabb24957b3b2a27561b2b2717204cde82b935ea3195bf53b2217448",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/wrap.ts",
      "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7",
      "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a3491cb9265a4ed19398a6a100ae6e6789c7272cfe887c80b4c439101afd7dc7",
        "signature": "4c8e88aa97caafb769b257fcde8ff8a50ef3e962fc7f3c51fba1c7c7e905dc3b",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/bug.js",
      "version": "676e1613c74cee0aa880dc1108ce9e0a07e184e35a7d78282bf4a480bf6124db",
      "signature": "6a80e6a4e41ee496beaac682ed90a523d5f8b87987da12ca01326156d552d025",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "676e1613c74cee0aa880dc1108ce9e0a07e184e35a7d78282bf4a480bf6124db",
        "signature": "6a80e6a4e41ee496beaac682ed90a523d5f8b87987da12ca01326156d552d025",
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
  "size": 949
}


SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/bug.js

Signatures::
(computed .d.ts) /home/src/workspaces/project/src/bug.js

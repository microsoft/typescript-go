currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/node_modules/.pnpm/dep@1.0.0/node_modules/dep/index.d.ts] *new* 
export declare function value(): number;
export declare function other(): string;

//// [/home/src/workspaces/project/node_modules/.pnpm/dep@1.0.0/node_modules/dep/package.json] *new* 
{ "name": "dep", "version": "1.0.0", "types": "index.d.ts" }
//// [/home/src/workspaces/project/node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts] *new* 
export declare function value(): number;
export declare function other(): string;

//// [/home/src/workspaces/project/node_modules/.pnpm/dep@2.0.0/node_modules/dep/package.json] *new* 
{ "name": "dep", "version": "2.0.0", "types": "index.d.ts" }
//// [/home/src/workspaces/project/node_modules/dep] -> /home/src/workspaces/project/node_modules/.pnpm/dep@1.0.0/node_modules/dep *new*
//// [/home/src/workspaces/project/src/route0.ts] *new* 
import { value } from "./shared";
export const r0: number = value();

//// [/home/src/workspaces/project/src/route1.ts] *new* 
import { value } from "./shared";
export const r1: number = value();

//// [/home/src/workspaces/project/src/route10.ts] *new* 
import { value } from "./shared";
export const r10: number = value();

//// [/home/src/workspaces/project/src/route11.ts] *new* 
import { value } from "./shared";
export const r11: number = value();

//// [/home/src/workspaces/project/src/route2.ts] *new* 
import { value } from "./shared";
export const r2: number = value();

//// [/home/src/workspaces/project/src/route3.ts] *new* 
import { value } from "./shared";
export const r3: number = value();

//// [/home/src/workspaces/project/src/route4.ts] *new* 
import { value } from "./shared";
export const r4: number = value();

//// [/home/src/workspaces/project/src/route5.ts] *new* 
import { value } from "./shared";
export const r5: number = value();

//// [/home/src/workspaces/project/src/route6.ts] *new* 
import { value } from "./shared";
export const r6: number = value();

//// [/home/src/workspaces/project/src/route7.ts] *new* 
import { value } from "./shared";
export const r7: number = value();

//// [/home/src/workspaces/project/src/route8.ts] *new* 
import { value } from "./shared";
export const r8: number = value();

//// [/home/src/workspaces/project/src/route9.ts] *new* 
import { value } from "./shared";
export const r9: number = value();

//// [/home/src/workspaces/project/src/shared.ts] *new* 
export { value } from "dep";

//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "incremental": true,
        "module": "esnext",
        "moduleResolution": "bundler",
        "strict": true,
        "noEmit": true,
    },
    "include": ["src/**/*.ts"],
}

tsgo --noEmit
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
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[3,15]],"packageJsons":["./node_modules/.pnpm/dep@1.0.0/node_modules/dep/package.json"],"fileNames":["lib.es2025.full.d.ts","./node_modules/.pnpm/dep@1.0.0/node_modules/dep/index.d.ts","./src/shared.ts","./src/route0.ts","./src/route1.ts","./src/route10.ts","./src/route11.ts","./src/route2.ts","./src/route3.ts","./src/route4.ts","./src/route5.ts","./src/route6.ts","./src/route7.ts","./src/route8.ts","./src/route9.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"ad2894a2ed059494de0d2a936fc203af-export declare function value(): number;\nexport declare function other(): string;\n","72efd51a0841b29d44444b0c071c9dfd-export { value } from \"dep\";\n","643e2357762975a68e62450d2a141ed7-import { value } from \"./shared\";\nexport const r0: number = value();\n","e70e0f57d1482b0aa84294501d5fc901-import { value } from \"./shared\";\nexport const r1: number = value();\n","8b70baa8a8e3f7fe4a640b6d0e75bbf5-import { value } from \"./shared\";\nexport const r10: number = value();\n","0e2481ead9651bcd9f25f1820a3846cb-import { value } from \"./shared\";\nexport const r11: number = value();\n","45ec354de757bb084c982f43e57d9364-import { value } from \"./shared\";\nexport const r2: number = value();\n","ea27de969e3a4e3331fff8b5cdd62b37-import { value } from \"./shared\";\nexport const r3: number = value();\n","02bd2af513f08ee5c874cbb56ca9ac52-import { value } from \"./shared\";\nexport const r4: number = value();\n","04a9b25767a947555005619663c77e2f-import { value } from \"./shared\";\nexport const r5: number = value();\n","4984282c908ecbe067f619d668810b58-import { value } from \"./shared\";\nexport const r6: number = value();\n","19c07f9bda4fea693b022ded99656b8a-import { value } from \"./shared\";\nexport const r7: number = value();\n","e3b98729cdd5987e8a8b76c24214bb96-import { value } from \"./shared\";\nexport const r8: number = value();\n","3c249c1ddc90a1a70a1e9dd2d1d0f004-import { value } from \"./shared\";\nexport const r9: number = value();\n"],"fileIdsList":[[3],[2]],"options":{"module":99,"strict":true},"referencedMap":[[4,1],[5,1],[6,1],[7,1],[8,1],[9,1],[10,1],[11,1],[12,1],[13,1],[14,1],[15,1],[3,2]],"affectedFilesPendingEmit":[4,5,6,7,8,9,10,11,12,13,14,15,3]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./src/shared.ts",
        "./src/route0.ts",
        "./src/route1.ts",
        "./src/route10.ts",
        "./src/route11.ts",
        "./src/route2.ts",
        "./src/route3.ts",
        "./src/route4.ts",
        "./src/route5.ts",
        "./src/route6.ts",
        "./src/route7.ts",
        "./src/route8.ts",
        "./src/route9.ts"
      ],
      "original": [
        3,
        15
      ]
    }
  ],
  "packageJsons": [
    "./node_modules/.pnpm/dep@1.0.0/node_modules/dep/package.json"
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./node_modules/.pnpm/dep@1.0.0/node_modules/dep/index.d.ts",
    "./src/shared.ts",
    "./src/route0.ts",
    "./src/route1.ts",
    "./src/route10.ts",
    "./src/route11.ts",
    "./src/route2.ts",
    "./src/route3.ts",
    "./src/route4.ts",
    "./src/route5.ts",
    "./src/route6.ts",
    "./src/route7.ts",
    "./src/route8.ts",
    "./src/route9.ts"
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
      "fileName": "./node_modules/.pnpm/dep@1.0.0/node_modules/dep/index.d.ts",
      "version": "ad2894a2ed059494de0d2a936fc203af-export declare function value(): number;\nexport declare function other(): string;\n",
      "signature": "ad2894a2ed059494de0d2a936fc203af-export declare function value(): number;\nexport declare function other(): string;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/shared.ts",
      "version": "72efd51a0841b29d44444b0c071c9dfd-export { value } from \"dep\";\n",
      "signature": "72efd51a0841b29d44444b0c071c9dfd-export { value } from \"dep\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route0.ts",
      "version": "643e2357762975a68e62450d2a141ed7-import { value } from \"./shared\";\nexport const r0: number = value();\n",
      "signature": "643e2357762975a68e62450d2a141ed7-import { value } from \"./shared\";\nexport const r0: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route1.ts",
      "version": "e70e0f57d1482b0aa84294501d5fc901-import { value } from \"./shared\";\nexport const r1: number = value();\n",
      "signature": "e70e0f57d1482b0aa84294501d5fc901-import { value } from \"./shared\";\nexport const r1: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route10.ts",
      "version": "8b70baa8a8e3f7fe4a640b6d0e75bbf5-import { value } from \"./shared\";\nexport const r10: number = value();\n",
      "signature": "8b70baa8a8e3f7fe4a640b6d0e75bbf5-import { value } from \"./shared\";\nexport const r10: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route11.ts",
      "version": "0e2481ead9651bcd9f25f1820a3846cb-import { value } from \"./shared\";\nexport const r11: number = value();\n",
      "signature": "0e2481ead9651bcd9f25f1820a3846cb-import { value } from \"./shared\";\nexport const r11: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route2.ts",
      "version": "45ec354de757bb084c982f43e57d9364-import { value } from \"./shared\";\nexport const r2: number = value();\n",
      "signature": "45ec354de757bb084c982f43e57d9364-import { value } from \"./shared\";\nexport const r2: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route3.ts",
      "version": "ea27de969e3a4e3331fff8b5cdd62b37-import { value } from \"./shared\";\nexport const r3: number = value();\n",
      "signature": "ea27de969e3a4e3331fff8b5cdd62b37-import { value } from \"./shared\";\nexport const r3: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route4.ts",
      "version": "02bd2af513f08ee5c874cbb56ca9ac52-import { value } from \"./shared\";\nexport const r4: number = value();\n",
      "signature": "02bd2af513f08ee5c874cbb56ca9ac52-import { value } from \"./shared\";\nexport const r4: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route5.ts",
      "version": "04a9b25767a947555005619663c77e2f-import { value } from \"./shared\";\nexport const r5: number = value();\n",
      "signature": "04a9b25767a947555005619663c77e2f-import { value } from \"./shared\";\nexport const r5: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route6.ts",
      "version": "4984282c908ecbe067f619d668810b58-import { value } from \"./shared\";\nexport const r6: number = value();\n",
      "signature": "4984282c908ecbe067f619d668810b58-import { value } from \"./shared\";\nexport const r6: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route7.ts",
      "version": "19c07f9bda4fea693b022ded99656b8a-import { value } from \"./shared\";\nexport const r7: number = value();\n",
      "signature": "19c07f9bda4fea693b022ded99656b8a-import { value } from \"./shared\";\nexport const r7: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route8.ts",
      "version": "e3b98729cdd5987e8a8b76c24214bb96-import { value } from \"./shared\";\nexport const r8: number = value();\n",
      "signature": "e3b98729cdd5987e8a8b76c24214bb96-import { value } from \"./shared\";\nexport const r8: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route9.ts",
      "version": "3c249c1ddc90a1a70a1e9dd2d1d0f004-import { value } from \"./shared\";\nexport const r9: number = value();\n",
      "signature": "3c249c1ddc90a1a70a1e9dd2d1d0f004-import { value } from \"./shared\";\nexport const r9: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./src/shared.ts"
    ],
    [
      "./node_modules/.pnpm/dep@1.0.0/node_modules/dep/index.d.ts"
    ]
  ],
  "options": {
    "module": 99,
    "strict": true
  },
  "referencedMap": {
    "./src/route0.ts": [
      "./src/shared.ts"
    ],
    "./src/route1.ts": [
      "./src/shared.ts"
    ],
    "./src/route10.ts": [
      "./src/shared.ts"
    ],
    "./src/route11.ts": [
      "./src/shared.ts"
    ],
    "./src/route2.ts": [
      "./src/shared.ts"
    ],
    "./src/route3.ts": [
      "./src/shared.ts"
    ],
    "./src/route4.ts": [
      "./src/shared.ts"
    ],
    "./src/route5.ts": [
      "./src/shared.ts"
    ],
    "./src/route6.ts": [
      "./src/shared.ts"
    ],
    "./src/route7.ts": [
      "./src/shared.ts"
    ],
    "./src/route8.ts": [
      "./src/shared.ts"
    ],
    "./src/route9.ts": [
      "./src/shared.ts"
    ],
    "./src/shared.ts": [
      "./node_modules/.pnpm/dep@1.0.0/node_modules/dep/index.d.ts"
    ]
  },
  "affectedFilesPendingEmit": [
    [
      "./src/route0.ts",
      "Js",
      4
    ],
    [
      "./src/route1.ts",
      "Js",
      5
    ],
    [
      "./src/route10.ts",
      "Js",
      6
    ],
    [
      "./src/route11.ts",
      "Js",
      7
    ],
    [
      "./src/route2.ts",
      "Js",
      8
    ],
    [
      "./src/route3.ts",
      "Js",
      9
    ],
    [
      "./src/route4.ts",
      "Js",
      10
    ],
    [
      "./src/route5.ts",
      "Js",
      11
    ],
    [
      "./src/route6.ts",
      "Js",
      12
    ],
    [
      "./src/route7.ts",
      "Js",
      13
    ],
    [
      "./src/route8.ts",
      "Js",
      14
    ],
    [
      "./src/route9.ts",
      "Js",
      15
    ],
    [
      "./src/shared.ts",
      "Js",
      3
    ]
  ],
  "size": 2964
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/node_modules/.pnpm/dep@1.0.0/node_modules/dep/index.d.ts
*refresh*    /home/src/workspaces/project/src/shared.ts
*refresh*    /home/src/workspaces/project/src/route0.ts
*refresh*    /home/src/workspaces/project/src/route1.ts
*refresh*    /home/src/workspaces/project/src/route10.ts
*refresh*    /home/src/workspaces/project/src/route11.ts
*refresh*    /home/src/workspaces/project/src/route2.ts
*refresh*    /home/src/workspaces/project/src/route3.ts
*refresh*    /home/src/workspaces/project/src/route4.ts
*refresh*    /home/src/workspaces/project/src/route5.ts
*refresh*    /home/src/workspaces/project/src/route6.ts
*refresh*    /home/src/workspaces/project/src/route7.ts
*refresh*    /home/src/workspaces/project/src/route8.ts
*refresh*    /home/src/workspaces/project/src/route9.ts
Signatures::


Edit [0]:: pnpm dependency version bump repoints node_modules/dep symlink

tsgo --noEmit
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[3,15]],"packageJsons":["./node_modules/.pnpm/dep@2.0.0/node_modules/dep/package.json"],"fileNames":["lib.es2025.full.d.ts","./node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts","./src/shared.ts","./src/route0.ts","./src/route1.ts","./src/route10.ts","./src/route11.ts","./src/route2.ts","./src/route3.ts","./src/route4.ts","./src/route5.ts","./src/route6.ts","./src/route7.ts","./src/route8.ts","./src/route9.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"ad2894a2ed059494de0d2a936fc203af-export declare function value(): number;\nexport declare function other(): string;\n","72efd51a0841b29d44444b0c071c9dfd-export { value } from \"dep\";\n","643e2357762975a68e62450d2a141ed7-import { value } from \"./shared\";\nexport const r0: number = value();\n","e70e0f57d1482b0aa84294501d5fc901-import { value } from \"./shared\";\nexport const r1: number = value();\n","8b70baa8a8e3f7fe4a640b6d0e75bbf5-import { value } from \"./shared\";\nexport const r10: number = value();\n","0e2481ead9651bcd9f25f1820a3846cb-import { value } from \"./shared\";\nexport const r11: number = value();\n","45ec354de757bb084c982f43e57d9364-import { value } from \"./shared\";\nexport const r2: number = value();\n","ea27de969e3a4e3331fff8b5cdd62b37-import { value } from \"./shared\";\nexport const r3: number = value();\n","02bd2af513f08ee5c874cbb56ca9ac52-import { value } from \"./shared\";\nexport const r4: number = value();\n","04a9b25767a947555005619663c77e2f-import { value } from \"./shared\";\nexport const r5: number = value();\n","4984282c908ecbe067f619d668810b58-import { value } from \"./shared\";\nexport const r6: number = value();\n","19c07f9bda4fea693b022ded99656b8a-import { value } from \"./shared\";\nexport const r7: number = value();\n","e3b98729cdd5987e8a8b76c24214bb96-import { value } from \"./shared\";\nexport const r8: number = value();\n","3c249c1ddc90a1a70a1e9dd2d1d0f004-import { value } from \"./shared\";\nexport const r9: number = value();\n"],"fileIdsList":[[3],[2]],"options":{"module":99,"strict":true},"referencedMap":[[4,1],[5,1],[6,1],[7,1],[8,1],[9,1],[10,1],[11,1],[12,1],[13,1],[14,1],[15,1],[3,2]],"affectedFilesPendingEmit":[4,5,6,7,8,9,10,11,12,13,14,15,3]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./src/shared.ts",
        "./src/route0.ts",
        "./src/route1.ts",
        "./src/route10.ts",
        "./src/route11.ts",
        "./src/route2.ts",
        "./src/route3.ts",
        "./src/route4.ts",
        "./src/route5.ts",
        "./src/route6.ts",
        "./src/route7.ts",
        "./src/route8.ts",
        "./src/route9.ts"
      ],
      "original": [
        3,
        15
      ]
    }
  ],
  "packageJsons": [
    "./node_modules/.pnpm/dep@2.0.0/node_modules/dep/package.json"
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts",
    "./src/shared.ts",
    "./src/route0.ts",
    "./src/route1.ts",
    "./src/route10.ts",
    "./src/route11.ts",
    "./src/route2.ts",
    "./src/route3.ts",
    "./src/route4.ts",
    "./src/route5.ts",
    "./src/route6.ts",
    "./src/route7.ts",
    "./src/route8.ts",
    "./src/route9.ts"
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
      "fileName": "./node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts",
      "version": "ad2894a2ed059494de0d2a936fc203af-export declare function value(): number;\nexport declare function other(): string;\n",
      "signature": "ad2894a2ed059494de0d2a936fc203af-export declare function value(): number;\nexport declare function other(): string;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/shared.ts",
      "version": "72efd51a0841b29d44444b0c071c9dfd-export { value } from \"dep\";\n",
      "signature": "72efd51a0841b29d44444b0c071c9dfd-export { value } from \"dep\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route0.ts",
      "version": "643e2357762975a68e62450d2a141ed7-import { value } from \"./shared\";\nexport const r0: number = value();\n",
      "signature": "643e2357762975a68e62450d2a141ed7-import { value } from \"./shared\";\nexport const r0: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route1.ts",
      "version": "e70e0f57d1482b0aa84294501d5fc901-import { value } from \"./shared\";\nexport const r1: number = value();\n",
      "signature": "e70e0f57d1482b0aa84294501d5fc901-import { value } from \"./shared\";\nexport const r1: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route10.ts",
      "version": "8b70baa8a8e3f7fe4a640b6d0e75bbf5-import { value } from \"./shared\";\nexport const r10: number = value();\n",
      "signature": "8b70baa8a8e3f7fe4a640b6d0e75bbf5-import { value } from \"./shared\";\nexport const r10: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route11.ts",
      "version": "0e2481ead9651bcd9f25f1820a3846cb-import { value } from \"./shared\";\nexport const r11: number = value();\n",
      "signature": "0e2481ead9651bcd9f25f1820a3846cb-import { value } from \"./shared\";\nexport const r11: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route2.ts",
      "version": "45ec354de757bb084c982f43e57d9364-import { value } from \"./shared\";\nexport const r2: number = value();\n",
      "signature": "45ec354de757bb084c982f43e57d9364-import { value } from \"./shared\";\nexport const r2: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route3.ts",
      "version": "ea27de969e3a4e3331fff8b5cdd62b37-import { value } from \"./shared\";\nexport const r3: number = value();\n",
      "signature": "ea27de969e3a4e3331fff8b5cdd62b37-import { value } from \"./shared\";\nexport const r3: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route4.ts",
      "version": "02bd2af513f08ee5c874cbb56ca9ac52-import { value } from \"./shared\";\nexport const r4: number = value();\n",
      "signature": "02bd2af513f08ee5c874cbb56ca9ac52-import { value } from \"./shared\";\nexport const r4: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route5.ts",
      "version": "04a9b25767a947555005619663c77e2f-import { value } from \"./shared\";\nexport const r5: number = value();\n",
      "signature": "04a9b25767a947555005619663c77e2f-import { value } from \"./shared\";\nexport const r5: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route6.ts",
      "version": "4984282c908ecbe067f619d668810b58-import { value } from \"./shared\";\nexport const r6: number = value();\n",
      "signature": "4984282c908ecbe067f619d668810b58-import { value } from \"./shared\";\nexport const r6: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route7.ts",
      "version": "19c07f9bda4fea693b022ded99656b8a-import { value } from \"./shared\";\nexport const r7: number = value();\n",
      "signature": "19c07f9bda4fea693b022ded99656b8a-import { value } from \"./shared\";\nexport const r7: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route8.ts",
      "version": "e3b98729cdd5987e8a8b76c24214bb96-import { value } from \"./shared\";\nexport const r8: number = value();\n",
      "signature": "e3b98729cdd5987e8a8b76c24214bb96-import { value } from \"./shared\";\nexport const r8: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/route9.ts",
      "version": "3c249c1ddc90a1a70a1e9dd2d1d0f004-import { value } from \"./shared\";\nexport const r9: number = value();\n",
      "signature": "3c249c1ddc90a1a70a1e9dd2d1d0f004-import { value } from \"./shared\";\nexport const r9: number = value();\n",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./src/shared.ts"
    ],
    [
      "./node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts"
    ]
  ],
  "options": {
    "module": 99,
    "strict": true
  },
  "referencedMap": {
    "./src/route0.ts": [
      "./src/shared.ts"
    ],
    "./src/route1.ts": [
      "./src/shared.ts"
    ],
    "./src/route10.ts": [
      "./src/shared.ts"
    ],
    "./src/route11.ts": [
      "./src/shared.ts"
    ],
    "./src/route2.ts": [
      "./src/shared.ts"
    ],
    "./src/route3.ts": [
      "./src/shared.ts"
    ],
    "./src/route4.ts": [
      "./src/shared.ts"
    ],
    "./src/route5.ts": [
      "./src/shared.ts"
    ],
    "./src/route6.ts": [
      "./src/shared.ts"
    ],
    "./src/route7.ts": [
      "./src/shared.ts"
    ],
    "./src/route8.ts": [
      "./src/shared.ts"
    ],
    "./src/route9.ts": [
      "./src/shared.ts"
    ],
    "./src/shared.ts": [
      "./node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts"
    ]
  },
  "affectedFilesPendingEmit": [
    [
      "./src/route0.ts",
      "Js",
      4
    ],
    [
      "./src/route1.ts",
      "Js",
      5
    ],
    [
      "./src/route10.ts",
      "Js",
      6
    ],
    [
      "./src/route11.ts",
      "Js",
      7
    ],
    [
      "./src/route2.ts",
      "Js",
      8
    ],
    [
      "./src/route3.ts",
      "Js",
      9
    ],
    [
      "./src/route4.ts",
      "Js",
      10
    ],
    [
      "./src/route5.ts",
      "Js",
      11
    ],
    [
      "./src/route6.ts",
      "Js",
      12
    ],
    [
      "./src/route7.ts",
      "Js",
      13
    ],
    [
      "./src/route8.ts",
      "Js",
      14
    ],
    [
      "./src/route9.ts",
      "Js",
      15
    ],
    [
      "./src/shared.ts",
      "Js",
      3
    ]
  ],
  "size": 2964
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts
*refresh*    /home/src/workspaces/project/src/shared.ts
*refresh*    /home/src/workspaces/project/src/route0.ts
*refresh*    /home/src/workspaces/project/src/route1.ts
*refresh*    /home/src/workspaces/project/src/route10.ts
*refresh*    /home/src/workspaces/project/src/route11.ts
*refresh*    /home/src/workspaces/project/src/route2.ts
*refresh*    /home/src/workspaces/project/src/route3.ts
*refresh*    /home/src/workspaces/project/src/route4.ts
*refresh*    /home/src/workspaces/project/src/route5.ts
*refresh*    /home/src/workspaces/project/src/route6.ts
*refresh*    /home/src/workspaces/project/src/route7.ts
*refresh*    /home/src/workspaces/project/src/route8.ts
*refresh*    /home/src/workspaces/project/src/route9.ts
Signatures::
(used version)   /home/src/workspaces/project/node_modules/.pnpm/dep@2.0.0/node_modules/dep/index.d.ts
(computed .d.ts) /home/src/workspaces/project/src/shared.ts
(used version)   /home/src/workspaces/project/src/route0.ts
(used version)   /home/src/workspaces/project/src/route1.ts
(used version)   /home/src/workspaces/project/src/route10.ts
(used version)   /home/src/workspaces/project/src/route11.ts
(used version)   /home/src/workspaces/project/src/route2.ts
(used version)   /home/src/workspaces/project/src/route3.ts
(used version)   /home/src/workspaces/project/src/route4.ts
(used version)   /home/src/workspaces/project/src/route5.ts
(used version)   /home/src/workspaces/project/src/route6.ts
(used version)   /home/src/workspaces/project/src/route7.ts
(used version)   /home/src/workspaces/project/src/route8.ts
(used version)   /home/src/workspaces/project/src/route9.ts


Edit [1]:: no change

tsgo --noEmit
ExitStatus:: Success
Output::

tsconfig.json::
SemanticDiagnostics::
Signatures::

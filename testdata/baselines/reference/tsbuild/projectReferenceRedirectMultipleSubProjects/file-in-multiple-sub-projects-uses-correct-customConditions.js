currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/app/app.ts] *new* 
import { platform } from "../pkg/index.js";
const p: "native" = platform;
//// [/home/src/workspaces/project/app/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declaration": true,
        "module": "nodenext",
        "moduleResolution": "nodenext",
        "customConditions": ["react-native"],
        "rootDir": ".",
        "outDir": "../dist/app",
        "skipDefaultLibCheck": true
    },
    "references": [
        { "path": "../pkg" }
    ]
}
//// [/home/src/workspaces/project/pkg/index.native.ts] *new* 
export { platform } from "./src/util.js";
//// [/home/src/workspaces/project/pkg/index.ts] *new* 
export { platform } from "./src/util.js";
//// [/home/src/workspaces/project/pkg/package.json] *new* 
{
    "name": "pkg",
    "type": "module",
    "exports": {
        ".": {
            "react-native": {
                "types": "./index.native.ts",
                "default": "./index.native.js"
            },
            "types": "./index.ts",
            "default": "./index.js"
        }
    }
}
//// [/home/src/workspaces/project/pkg/src/util.native.ts] *new* 
export const platform = "native" as const;
//// [/home/src/workspaces/project/pkg/src/util.ts] *new* 
export const platform = "web" as const;
//// [/home/src/workspaces/project/pkg/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declaration": true,
        "module": "nodenext",
        "moduleResolution": "nodenext",
        "rootDir": ".",
        "outDir": "../dist/pkg",
        "skipDefaultLibCheck": true
    },
    "references": [
        { "path": "./tsconfig.native.json" }
    ]
}
//// [/home/src/workspaces/project/pkg/tsconfig.native.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declaration": true,
        "module": "nodenext",
        "moduleResolution": "nodenext",
        "customConditions": ["react-native"],
        "moduleSuffixes": [".native", ""],
        "rootDir": ".",
        "outDir": "../dist/pkg-native",
        "skipDefaultLibCheck": true
    }
}

tsgo --b app --verbose
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * pkg/tsconfig.native.json
    * pkg/tsconfig.json
    * app/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'pkg/tsconfig.native.json' is out of date because output file 'dist/pkg-native/tsconfig.native.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg/tsconfig.native.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg/tsconfig.json' is out of date because output file 'dist/pkg/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'app/tsconfig.json' is out of date because output file 'dist/app/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'app/tsconfig.json'...

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
//// [/home/src/workspaces/project/dist/app/app.d.ts] *new* 
export {};

//// [/home/src/workspaces/project/dist/app/app.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index_js_1 = require("../pkg/index.js");
const p = index_js_1.platform;

//// [/home/src/workspaces/project/dist/app/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[4],"fileNames":["lib.es2025.full.d.ts","../pkg-native/src/util.native.d.ts","../pkg-native/index.d.ts","../../app/app.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n","c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",{"version":"783d622927808e0e8c2767aac18bff30-import { platform } from \"../pkg/index.js\";\nconst p: \"native\" = platform;","signature":"abe7d9981d6018efb6b2b794f40a1607-export {};\n","impliedNodeFormat":1}],"fileIdsList":[[3],[2]],"options":{"composite":true,"declaration":true,"module":199,"outDir":"./","rootDir":"../../app","skipDefaultLibCheck":true},"referencedMap":[[4,1],[3,2]],"latestChangedDtsFile":"./app.d.ts"}
//// [/home/src/workspaces/project/dist/app/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../../app/app.ts"
      ],
      "original": 4
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../pkg-native/src/util.native.d.ts",
    "../pkg-native/index.d.ts",
    "../../app/app.ts"
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
      "fileName": "../pkg-native/src/util.native.d.ts",
      "version": "1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n",
      "signature": "1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../pkg-native/index.d.ts",
      "version": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "signature": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../../app/app.ts",
      "version": "783d622927808e0e8c2767aac18bff30-import { platform } from \"../pkg/index.js\";\nconst p: \"native\" = platform;",
      "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "783d622927808e0e8c2767aac18bff30-import { platform } from \"../pkg/index.js\";\nconst p: \"native\" = platform;",
        "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../pkg-native/index.d.ts"
    ],
    [
      "../pkg-native/src/util.native.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declaration": true,
    "module": 199,
    "outDir": "./",
    "rootDir": "../../app",
    "skipDefaultLibCheck": true
  },
  "referencedMap": {
    "../../app/app.ts": [
      "../pkg-native/index.d.ts"
    ],
    "../pkg-native/index.d.ts": [
      "../pkg-native/src/util.native.d.ts"
    ]
  },
  "latestChangedDtsFile": "./app.d.ts",
  "size": 1525
}
//// [/home/src/workspaces/project/dist/pkg-native/index.d.ts] *new* 
export { platform } from "./src/util.js";

//// [/home/src/workspaces/project/dist/pkg-native/index.js] *new* 
export { platform } from "./src/util.js";

//// [/home/src/workspaces/project/dist/pkg-native/index.native.d.ts] *new* 
export { platform } from "./src/util.js";

//// [/home/src/workspaces/project/dist/pkg-native/index.native.js] *new* 
export { platform } from "./src/util.js";

//// [/home/src/workspaces/project/dist/pkg-native/src/util.d.ts] *new* 
export declare const platform: "web";

//// [/home/src/workspaces/project/dist/pkg-native/src/util.js] *new* 
export const platform = "web";

//// [/home/src/workspaces/project/dist/pkg-native/src/util.native.d.ts] *new* 
export declare const platform: "native";

//// [/home/src/workspaces/project/dist/pkg-native/src/util.native.js] *new* 
export const platform = "native";

//// [/home/src/workspaces/project/dist/pkg-native/tsconfig.native.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["lib.es2025.full.d.ts","../../pkg/src/util.native.ts","../../pkg/index.native.ts","../../pkg/index.ts","../../pkg/src/util.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"c1f8db351239c175fbb3960262e28684-export const platform = \"native\" as const;","signature":"1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n","impliedNodeFormat":99},{"version":"922e0f658a8807c3f0f559560501905c-export { platform } from \"./src/util.js\";","signature":"c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n","impliedNodeFormat":99},{"version":"922e0f658a8807c3f0f559560501905c-export { platform } from \"./src/util.js\";","signature":"c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n","impliedNodeFormat":99},{"version":"7941de8fb997b556c0afef2b586d7205-export const platform = \"web\" as const;","signature":"5082e4a38cc5cc308625a8754198c0e3-export declare const platform: \"web\";\n","impliedNodeFormat":99}],"fileIdsList":[[2]],"options":{"composite":true,"declaration":true,"module":199,"outDir":"./","rootDir":"../../pkg","skipDefaultLibCheck":true},"referencedMap":[[3,1],[4,1]],"latestChangedDtsFile":"./src/util.d.ts"}
//// [/home/src/workspaces/project/dist/pkg-native/tsconfig.native.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../../pkg/src/util.native.ts",
        "../../pkg/index.native.ts",
        "../../pkg/index.ts",
        "../../pkg/src/util.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../../pkg/src/util.native.ts",
    "../../pkg/index.native.ts",
    "../../pkg/index.ts",
    "../../pkg/src/util.ts"
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
      "fileName": "../../pkg/src/util.native.ts",
      "version": "c1f8db351239c175fbb3960262e28684-export const platform = \"native\" as const;",
      "signature": "1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "c1f8db351239c175fbb3960262e28684-export const platform = \"native\" as const;",
        "signature": "1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n",
        "impliedNodeFormat": 99
      }
    },
    {
      "fileName": "../../pkg/index.native.ts",
      "version": "922e0f658a8807c3f0f559560501905c-export { platform } from \"./src/util.js\";",
      "signature": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "922e0f658a8807c3f0f559560501905c-export { platform } from \"./src/util.js\";",
        "signature": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
        "impliedNodeFormat": 99
      }
    },
    {
      "fileName": "../../pkg/index.ts",
      "version": "922e0f658a8807c3f0f559560501905c-export { platform } from \"./src/util.js\";",
      "signature": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "922e0f658a8807c3f0f559560501905c-export { platform } from \"./src/util.js\";",
        "signature": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
        "impliedNodeFormat": 99
      }
    },
    {
      "fileName": "../../pkg/src/util.ts",
      "version": "7941de8fb997b556c0afef2b586d7205-export const platform = \"web\" as const;",
      "signature": "5082e4a38cc5cc308625a8754198c0e3-export declare const platform: \"web\";\n",
      "impliedNodeFormat": "ESNext",
      "original": {
        "version": "7941de8fb997b556c0afef2b586d7205-export const platform = \"web\" as const;",
        "signature": "5082e4a38cc5cc308625a8754198c0e3-export declare const platform: \"web\";\n",
        "impliedNodeFormat": 99
      }
    }
  ],
  "fileIdsList": [
    [
      "../../pkg/src/util.native.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declaration": true,
    "module": 199,
    "outDir": "./",
    "rootDir": "../../pkg",
    "skipDefaultLibCheck": true
  },
  "referencedMap": {
    "../../pkg/index.native.ts": [
      "../../pkg/src/util.native.ts"
    ],
    "../../pkg/index.ts": [
      "../../pkg/src/util.native.ts"
    ]
  },
  "latestChangedDtsFile": "./src/util.d.ts",
  "size": 2004
}
//// [/home/src/workspaces/project/dist/pkg/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["lib.es2025.full.d.ts","../pkg-native/src/util.native.d.ts","../pkg-native/index.native.d.ts","../pkg-native/index.d.ts","../pkg-native/src/util.d.ts","../../pkg/src/util.native.ts","../../pkg/index.native.ts","../../pkg/index.ts","../../pkg/src/util.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n","c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n","c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n","5082e4a38cc5cc308625a8754198c0e3-export declare const platform: \"web\";\n"],"fileIdsList":[[2]],"options":{"composite":true,"declaration":true,"module":199,"outDir":"./","rootDir":"../../pkg","skipDefaultLibCheck":true},"referencedMap":[[4,1],[3,1]],"resolvedRoot":[[2,6],[3,7],[4,8],[5,9]]}
//// [/home/src/workspaces/project/dist/pkg/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../pkg-native/src/util.native.d.ts",
        "../pkg-native/index.native.d.ts",
        "../pkg-native/index.d.ts",
        "../pkg-native/src/util.d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../pkg-native/src/util.native.d.ts",
    "../pkg-native/index.native.d.ts",
    "../pkg-native/index.d.ts",
    "../pkg-native/src/util.d.ts",
    "../../pkg/src/util.native.ts",
    "../../pkg/index.native.ts",
    "../../pkg/index.ts",
    "../../pkg/src/util.ts"
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
      "fileName": "../pkg-native/src/util.native.d.ts",
      "version": "1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n",
      "signature": "1f876b2eee633f65aa2e7817bfee737a-export declare const platform: \"native\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../pkg-native/index.native.d.ts",
      "version": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "signature": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../pkg-native/index.d.ts",
      "version": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "signature": "c28ca27b5e491b0e7ccd4741b9f9aba4-export { platform } from \"./src/util.js\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../pkg-native/src/util.d.ts",
      "version": "5082e4a38cc5cc308625a8754198c0e3-export declare const platform: \"web\";\n",
      "signature": "5082e4a38cc5cc308625a8754198c0e3-export declare const platform: \"web\";\n",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "../pkg-native/src/util.native.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declaration": true,
    "module": 199,
    "outDir": "./",
    "rootDir": "../../pkg",
    "skipDefaultLibCheck": true
  },
  "referencedMap": {
    "../pkg-native/index.d.ts": [
      "../pkg-native/src/util.native.d.ts"
    ],
    "../pkg-native/index.native.d.ts": [
      "../pkg-native/src/util.native.d.ts"
    ]
  },
  "resolvedRoot": [
    [
      "../pkg-native/src/util.native.d.ts",
      "../../pkg/src/util.native.ts"
    ],
    [
      "../pkg-native/index.native.d.ts",
      "../../pkg/index.native.ts"
    ],
    [
      "../pkg-native/index.d.ts",
      "../../pkg/index.ts"
    ],
    [
      "../pkg-native/src/util.d.ts",
      "../../pkg/src/util.ts"
    ]
  ],
  "size": 1629
}

pkg/tsconfig.native.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/pkg/src/util.native.ts
*refresh*    /home/src/workspaces/project/pkg/index.native.ts
*refresh*    /home/src/workspaces/project/pkg/index.ts
*refresh*    /home/src/workspaces/project/pkg/src/util.ts
Signatures::
(stored at emit) /home/src/workspaces/project/pkg/src/util.native.ts
(stored at emit) /home/src/workspaces/project/pkg/index.native.ts
(stored at emit) /home/src/workspaces/project/pkg/index.ts
(stored at emit) /home/src/workspaces/project/pkg/src/util.ts

pkg/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/dist/pkg-native/src/util.native.d.ts
*refresh*    /home/src/workspaces/project/dist/pkg-native/index.native.d.ts
*refresh*    /home/src/workspaces/project/dist/pkg-native/index.d.ts
*refresh*    /home/src/workspaces/project/dist/pkg-native/src/util.d.ts
Signatures::

app/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/dist/pkg-native/src/util.native.d.ts
*refresh*    /home/src/workspaces/project/dist/pkg-native/index.d.ts
*refresh*    /home/src/workspaces/project/app/app.ts
Signatures::
(stored at emit) /home/src/workspaces/project/app/app.ts

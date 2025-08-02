currentDirectory::/home/src/workspaces/solution
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/solution/project1/src/a.ts] *new* 
export const a = 10;const aLocal = 10;
//// [/home/src/workspaces/solution/project1/src/b.ts] *new* 
export const b = 10;const bLocal = 10;
//// [/home/src/workspaces/solution/project1/src/c.ts] *new* 
import { a } from "./a";export const c = a;
//// [/home/src/workspaces/solution/project1/src/d.ts] *new* 
import { b } from "./b";export const d = b;
//// [/home/src/workspaces/solution/project1/src/tsconfig.json] *new* 
{
    "compilerOptions": { "incremental": true, "declaration": true }
}
//// [/home/src/workspaces/solution/project2/src/e.ts] *new* 
export const e = 10;
//// [/home/src/workspaces/solution/project2/src/f.ts] *new* 
import { a } from "../../project1/src/a"; export const f = a;
//// [/home/src/workspaces/solution/project2/src/g.ts] *new* 
import { b } from "../../project1/src/b"; export const g = b;
//// [/home/src/workspaces/solution/project2/src/tsconfig.json] *new* 
{
    "compilerOptions": { "incremental": true, "declaration": true },
    "references": [{ "path": "../../project1/src" }]
}

tsgo --b project2/src --verbose --emitDeclarationOnly
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

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
//// [/home/src/workspaces/solution/project1/src/a.d.ts] *new* 
export declare const a = 10;

//// [/home/src/workspaces/solution/project1/src/b.d.ts] *new* 
export declare const b = 10;

//// [/home/src/workspaces/solution/project1/src/c.d.ts] *new* 
export declare const c = 10;

//// [/home/src/workspaces/solution/project1/src/d.d.ts] *new* 
export declare const d = 10;

//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"7dea2bd009b9cd0fd54ca48f1d10ffe0-export const a = 10;const aLocal = 10;","signature":"589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n","impliedNodeFormat":1},{"version":"016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;","signature":"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "7dea2bd009b9cd0fd54ca48f1d10ffe0-export const a = 10;const aLocal = 10;",
      "signature": "589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dea2bd009b9cd0fd54ca48f1d10ffe0-export const a = 10;const aLocal = 10;",
        "signature": "589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
        "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1783
}
//// [/home/src/workspaces/solution/project2/src/e.d.ts] *new* 
export declare const e = 10;

//// [/home/src/workspaces/solution/project2/src/f.d.ts] *new* 
export declare const f = 10;

//// [/home/src/workspaces/solution/project2/src/g.d.ts] *new* 
export declare const g = 10;

//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2,4,6],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./e.ts","../../project1/src/a.d.ts","./f.ts","../../project1/src/b.d.ts","./g.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"f0e6085d8cf835a263334b8d65c348d6-export const e = 10;","signature":"03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n","impliedNodeFormat":1},"589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n",{"version":"7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;","signature":"b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n","impliedNodeFormat":1},"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",{"version":"8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;","signature":"0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[3],[5]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[6,2]]}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./e.ts"
      ],
      "original": 2
    },
    {
      "files": [
        "./f.ts"
      ],
      "original": 4
    },
    {
      "files": [
        "./g.ts"
      ],
      "original": 6
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./e.ts",
    "../../project1/src/a.d.ts",
    "./f.ts",
    "../../project1/src/b.d.ts",
    "./g.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./e.ts",
      "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
      "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
        "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/a.d.ts",
      "version": "589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n",
      "signature": "589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./f.ts",
      "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
      "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
        "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/b.d.ts",
      "version": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./g.ts",
      "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
      "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
        "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../../project1/src/a.d.ts"
    ],
    [
      "../../project1/src/b.d.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./f.ts": [
      "../../project1/src/a.d.ts"
    ],
    "./g.ts": [
      "../../project1/src/b.d.ts"
    ]
  },
  "size": 1794
}

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/solution/project1/src/a.ts
*refresh*    /home/src/workspaces/solution/project1/src/b.ts
*refresh*    /home/src/workspaces/solution/project1/src/c.ts
*refresh*    /home/src/workspaces/solution/project1/src/d.ts
Signatures::
(stored at emit) /home/src/workspaces/solution/project1/src/a.ts
(stored at emit) /home/src/workspaces/solution/project1/src/b.ts
(stored at emit) /home/src/workspaces/solution/project1/src/c.ts
(stored at emit) /home/src/workspaces/solution/project1/src/d.ts

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/solution/project2/src/e.ts
*refresh*    /home/src/workspaces/solution/project1/src/a.d.ts
*refresh*    /home/src/workspaces/solution/project2/src/f.ts
*refresh*    /home/src/workspaces/solution/project1/src/b.d.ts
*refresh*    /home/src/workspaces/solution/project2/src/g.ts
Signatures::
(stored at emit) /home/src/workspaces/solution/project2/src/e.ts
(stored at emit) /home/src/workspaces/solution/project2/src/f.ts
(stored at emit) /home/src/workspaces/solution/project2/src/g.ts


Edit [0]:: no change

tsgo --b project2/src --verbose --emitDeclarationOnly
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is up to date because newest input 'project1/src/a.ts' is older than output 'project1/src/tsconfig.tsbuildinfo'

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is up to date because newest input 'project2/src/e.ts' is older than output 'project2/src/tsconfig.tsbuildinfo'




Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,7 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
-
-[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...
-
-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
-
-[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is up to date because newest input 'project1/src/a.ts' is older than output 'project1/src/tsconfig.tsbuildinfo'
+
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is up to date because newest input 'project2/src/e.ts' is older than output 'project2/src/tsconfig.tsbuildinfo'


Edit [1]:: local change
//// [/home/src/workspaces/solution/project1/src/a.ts] *modified* 
export const a = 10;const aLocal = 10;const aa = 10;

tsgo --b project2/src --verbose --emitDeclarationOnly
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/a.ts'

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'project2/src/tsconfig.json'...

//// [/home/src/workspaces/solution/project1/src/a.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"cee20a32c75d7b26d548b27e24acb62a-export const a = 10;const aLocal = 10;const aa = 10;","signature":"589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n","impliedNodeFormat":1},{"version":"016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;","signature":"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "cee20a32c75d7b26d548b27e24acb62a-export const a = 10;const aLocal = 10;const aa = 10;",
      "signature": "589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "cee20a32c75d7b26d548b27e24acb62a-export const a = 10;const aLocal = 10;const aa = 10;",
        "signature": "589173ef1057b7817772d5a947bd33ba-export declare const a = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
        "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1797
}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *mTime changed*

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/a.ts
Signatures::
(computed .d.ts) /home/src/workspaces/solution/project1/src/a.ts

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
Signatures::


Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,13 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/a.ts'

 [[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

 [[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...
+
+[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'project2/src/tsconfig.json'...


Edit [2]:: non local change
//// [/home/src/workspaces/solution/project1/src/a.ts] *modified* 
export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;

tsgo --b project2/src --verbose --emitDeclarationOnly
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/a.ts'

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

//// [/home/src/workspaces/solution/project1/src/a.d.ts] *modified* 
export declare const a = 10;
export declare const aaa = 10;

//// [/home/src/workspaces/solution/project1/src/c.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;","signature":"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n","impliedNodeFormat":1},{"version":"016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;","signature":"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
        "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
        "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1851
}
//// [/home/src/workspaces/solution/project2/src/f.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2,4,6],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./e.ts","../../project1/src/a.d.ts","./f.ts","../../project1/src/b.d.ts","./g.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"f0e6085d8cf835a263334b8d65c348d6-export const e = 10;","signature":"03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n","impliedNodeFormat":1},"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",{"version":"7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;","signature":"b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n","impliedNodeFormat":1},"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",{"version":"8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;","signature":"0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[3],[5]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[6,2]]}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./e.ts"
      ],
      "original": 2
    },
    {
      "files": [
        "./f.ts"
      ],
      "original": 4
    },
    {
      "files": [
        "./g.ts"
      ],
      "original": 6
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./e.ts",
    "../../project1/src/a.d.ts",
    "./f.ts",
    "../../project1/src/b.d.ts",
    "./g.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./e.ts",
      "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
      "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
        "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/a.d.ts",
      "version": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./f.ts",
      "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
      "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
        "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/b.d.ts",
      "version": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./g.ts",
      "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
      "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
        "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../../project1/src/a.d.ts"
    ],
    [
      "../../project1/src/b.d.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./f.ts": [
      "../../project1/src/a.d.ts"
    ],
    "./g.ts": [
      "../../project1/src/b.d.ts"
    ]
  },
  "size": 1826
}

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/a.ts
*refresh*    /home/src/workspaces/solution/project1/src/c.ts
Signatures::
(computed .d.ts) /home/src/workspaces/solution/project1/src/a.ts
(computed .d.ts) /home/src/workspaces/solution/project1/src/c.ts

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/a.d.ts
*refresh*    /home/src/workspaces/solution/project2/src/f.ts
Signatures::
(used version)   /home/src/workspaces/solution/project1/src/a.d.ts
(computed .d.ts) /home/src/workspaces/solution/project2/src/f.ts


Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,11 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/a.ts'

 [[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

 [[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...


Edit [3]:: emit js files

tsgo --b project2/src --verbose
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because buildinfo file 'project1/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because buildinfo file 'project2/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

//// [/home/src/workspaces/solution/project1/src/a.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.aaa = exports.a = void 0;
exports.a = 10;
const aLocal = 10;
const aa = 10;
exports.aaa = 10;

//// [/home/src/workspaces/solution/project1/src/b.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.b = void 0;
exports.b = 10;
const bLocal = 10;

//// [/home/src/workspaces/solution/project1/src/c.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.c = void 0;
const a_1 = require("./a");
exports.c = a_1.a;

//// [/home/src/workspaces/solution/project1/src/d.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.d = void 0;
const b_1 = require("./b");
exports.d = b_1.b;

//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;","signature":"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n","impliedNodeFormat":1},{"version":"016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;","signature":"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
        "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "016155b7db85dc5a88a3933712286fb6-export const b = 10;const bLocal = 10;",
        "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1824
}
//// [/home/src/workspaces/solution/project2/src/e.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.e = void 0;
exports.e = 10;

//// [/home/src/workspaces/solution/project2/src/f.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.f = void 0;
const a_1 = require("../../project1/src/a");
exports.f = a_1.a;

//// [/home/src/workspaces/solution/project2/src/g.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.g = void 0;
const b_1 = require("../../project1/src/b");
exports.g = b_1.b;

//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2,4,6],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./e.ts","../../project1/src/a.d.ts","./f.ts","../../project1/src/b.d.ts","./g.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"f0e6085d8cf835a263334b8d65c348d6-export const e = 10;","signature":"03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n","impliedNodeFormat":1},"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",{"version":"7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;","signature":"b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n","impliedNodeFormat":1},"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",{"version":"8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;","signature":"0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[3],[5]],"options":{"declaration":true},"referencedMap":[[4,1],[6,2]]}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./e.ts"
      ],
      "original": 2
    },
    {
      "files": [
        "./f.ts"
      ],
      "original": 4
    },
    {
      "files": [
        "./g.ts"
      ],
      "original": 6
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./e.ts",
    "../../project1/src/a.d.ts",
    "./f.ts",
    "../../project1/src/b.d.ts",
    "./g.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./e.ts",
      "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
      "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
        "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/a.d.ts",
      "version": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./f.ts",
      "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
      "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
        "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/b.d.ts",
      "version": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./g.ts",
      "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
      "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
        "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../../project1/src/a.d.ts"
    ],
    [
      "../../project1/src/b.d.ts"
    ]
  ],
  "options": {
    "declaration": true
  },
  "referencedMap": {
    "./f.ts": [
      "../../project1/src/a.d.ts"
    ],
    "./g.ts": [
      "../../project1/src/b.d.ts"
    ]
  },
  "size": 1799
}

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
Signatures::

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
Signatures::


Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,11 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because buildinfo file 'project1/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

 [[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because buildinfo file 'project2/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

 [[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...


Edit [4]:: no change

tsgo --b project2/src --verbose --emitDeclarationOnly
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is up to date because newest input 'project1/src/a.ts' is older than output 'project1/src/tsconfig.tsbuildinfo'

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is up to date because newest input 'project2/src/e.ts' is older than output 'project2/src/tsconfig.tsbuildinfo'




Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,7 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
-
-[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...
-
-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
-
-[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is up to date because newest input 'project1/src/a.ts' is older than output 'project1/src/tsconfig.tsbuildinfo'
+
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is up to date because newest input 'project2/src/e.ts' is older than output 'project2/src/tsconfig.tsbuildinfo'


Edit [5]:: js emit with change without emitDeclarationOnly
//// [/home/src/workspaces/solution/project1/src/b.ts] *modified* 
export const b = 10;const bLocal = 10;const alocal = 10;

tsgo --b project2/src --verbose
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/b.ts'

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'project2/src/tsconfig.json'...

//// [/home/src/workspaces/solution/project1/src/b.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/b.js] *modified* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.b = void 0;
exports.b = 10;
const bLocal = 10;
const alocal = 10;

//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;","signature":"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n","impliedNodeFormat":1},{"version":"92dce007d82268ecfc5505c7cff769dc-export const b = 10;const bLocal = 10;const alocal = 10;","signature":"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
        "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "92dce007d82268ecfc5505c7cff769dc-export const b = 10;const bLocal = 10;const alocal = 10;",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "92dce007d82268ecfc5505c7cff769dc-export const b = 10;const bLocal = 10;const alocal = 10;",
        "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1842
}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *mTime changed*

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/b.ts
Signatures::
(computed .d.ts) /home/src/workspaces/solution/project1/src/b.ts

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
Signatures::


Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,13 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/b.ts'

 [[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

 [[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...
+
+[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'project2/src/tsconfig.json'...


Edit [6]:: local change
//// [/home/src/workspaces/solution/project1/src/b.ts] *modified* 
export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;

tsgo --b project2/src --verbose --emitDeclarationOnly
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/b.ts'

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'project2/src/tsconfig.json'...

//// [/home/src/workspaces/solution/project1/src/b.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;","signature":"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n","impliedNodeFormat":1},{"version":"a8b6ed7f958eee637900d942e6462f1c-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;","signature":"7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
        "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "a8b6ed7f958eee637900d942e6462f1c-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;",
      "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a8b6ed7f958eee637900d942e6462f1c-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;",
        "signature": "7c03652c2857b770a29bade9a57608cd-export declare const b = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1885
}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *mTime changed*

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/b.ts
Signatures::
(computed .d.ts) /home/src/workspaces/solution/project1/src/b.ts

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
Signatures::


Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,13 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/b.ts'

 [[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

 [[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...
+
+[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'project2/src/tsconfig.json'...


Edit [7]:: non local change
//// [/home/src/workspaces/solution/project1/src/b.ts] *modified* 
export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;

tsgo --b project2/src --verbose --emitDeclarationOnly
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/b.ts'

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

//// [/home/src/workspaces/solution/project1/src/b.d.ts] *modified* 
export declare const b = 10;
export declare const aaaaa = 10;

//// [/home/src/workspaces/solution/project1/src/d.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;","signature":"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n","impliedNodeFormat":1},{"version":"77770f43a5a5ad2db2ae0de1f18c0477-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;","signature":"7e77209a132838642d57821a72236b9e-export declare const b = 10;\nexport declare const aaaaa = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
        "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "77770f43a5a5ad2db2ae0de1f18c0477-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;",
      "signature": "7e77209a132838642d57821a72236b9e-export declare const b = 10;\nexport declare const aaaaa = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "77770f43a5a5ad2db2ae0de1f18c0477-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;",
        "signature": "7e77209a132838642d57821a72236b9e-export declare const b = 10;\nexport declare const aaaaa = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1943
}
//// [/home/src/workspaces/solution/project2/src/g.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2,4,6],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./e.ts","../../project1/src/a.d.ts","./f.ts","../../project1/src/b.d.ts","./g.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"f0e6085d8cf835a263334b8d65c348d6-export const e = 10;","signature":"03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n","impliedNodeFormat":1},"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",{"version":"7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;","signature":"b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n","impliedNodeFormat":1},"7e77209a132838642d57821a72236b9e-export declare const b = 10;\nexport declare const aaaaa = 10;\n",{"version":"8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;","signature":"0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[3],[5]],"options":{"emitDeclarationOnly":true,"declaration":true},"referencedMap":[[4,1],[6,2]]}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./e.ts"
      ],
      "original": 2
    },
    {
      "files": [
        "./f.ts"
      ],
      "original": 4
    },
    {
      "files": [
        "./g.ts"
      ],
      "original": 6
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./e.ts",
    "../../project1/src/a.d.ts",
    "./f.ts",
    "../../project1/src/b.d.ts",
    "./g.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./e.ts",
      "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
      "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
        "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/a.d.ts",
      "version": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./f.ts",
      "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
      "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
        "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/b.d.ts",
      "version": "7e77209a132838642d57821a72236b9e-export declare const b = 10;\nexport declare const aaaaa = 10;\n",
      "signature": "7e77209a132838642d57821a72236b9e-export declare const b = 10;\nexport declare const aaaaa = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./g.ts",
      "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
      "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
        "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../../project1/src/a.d.ts"
    ],
    [
      "../../project1/src/b.d.ts"
    ]
  ],
  "options": {
    "emitDeclarationOnly": true,
    "declaration": true
  },
  "referencedMap": {
    "./f.ts": [
      "../../project1/src/a.d.ts"
    ],
    "./g.ts": [
      "../../project1/src/b.d.ts"
    ]
  },
  "size": 1860
}

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/b.ts
*refresh*    /home/src/workspaces/solution/project1/src/d.ts
Signatures::
(computed .d.ts) /home/src/workspaces/solution/project1/src/b.ts
(computed .d.ts) /home/src/workspaces/solution/project1/src/d.ts

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/b.d.ts
*refresh*    /home/src/workspaces/solution/project2/src/g.ts
Signatures::
(used version)   /home/src/workspaces/solution/project1/src/b.d.ts
(computed .d.ts) /home/src/workspaces/solution/project2/src/g.ts


Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,11 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output 'project1/src/tsconfig.tsbuildinfo' is older than input 'project1/src/b.ts'

 [[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output 'project2/src/tsconfig.tsbuildinfo' is older than input 'project1/src'

 [[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...


Edit [8]:: js emit with change without emitDeclarationOnly
//// [/home/src/workspaces/solution/project1/src/b.ts] *modified* 
export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;export const a2 = 10;

tsgo --b project2/src --verbose
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * project1/src/tsconfig.json
    * project2/src/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because buildinfo file 'project1/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

[[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because buildinfo file 'project2/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

[[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

//// [/home/src/workspaces/solution/project1/src/a.js] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/b.d.ts] *modified* 
export declare const b = 10;
export declare const aaaaa = 10;
export declare const a2 = 10;

//// [/home/src/workspaces/solution/project1/src/b.js] *modified* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.a2 = exports.aaaaa = exports.b = void 0;
exports.b = 10;
const bLocal = 10;
const alocal = 10;
const aaaa = 10;
exports.aaaaa = 10;
exports.a2 = 10;

//// [/home/src/workspaces/solution/project1/src/c.js] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/d.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/d.js] *rewrite with same content*
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,5]],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./a.ts","./b.ts","./c.ts","./d.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;","signature":"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n","impliedNodeFormat":1},{"version":"afaee26f7edaee0dd3ce5d62d9d12490-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;export const a2 = 10;","signature":"195d8237d70e02ed04ff98ea3f31f6a5-export declare const b = 10;\nexport declare const aaaaa = 10;\nexport declare const a2 = 10;\n","impliedNodeFormat":1},{"version":"54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;","signature":"8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n","impliedNodeFormat":1},{"version":"48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;","signature":"cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[2],[3]],"options":{"declaration":true},"referencedMap":[[4,1],[5,2]]}
//// [/home/src/workspaces/solution/project1/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./a.ts",
        "./b.ts",
        "./c.ts",
        "./d.ts"
      ],
      "original": [
        2,
        5
      ]
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./a.ts",
    "./b.ts",
    "./c.ts",
    "./d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e719bd7527b04c1b42d578fca4e36b9f-export const a = 10;const aLocal = 10;const aa = 10;export const aaa = 10;",
        "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./b.ts",
      "version": "afaee26f7edaee0dd3ce5d62d9d12490-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;export const a2 = 10;",
      "signature": "195d8237d70e02ed04ff98ea3f31f6a5-export declare const b = 10;\nexport declare const aaaaa = 10;\nexport declare const a2 = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "afaee26f7edaee0dd3ce5d62d9d12490-export const b = 10;const bLocal = 10;const alocal = 10;const aaaa = 10;export const aaaaa = 10;export const a2 = 10;",
        "signature": "195d8237d70e02ed04ff98ea3f31f6a5-export declare const b = 10;\nexport declare const aaaaa = 10;\nexport declare const a2 = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./c.ts",
      "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
      "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "54735153f9b8943813e5b419aa52078f-import { a } from \"./a\";export const c = a;",
        "signature": "8e8d7e212457b775e32f78ac43bac800-export declare const c = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./d.ts",
      "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
      "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "48be508424a37d4ba69a717688ca1847-import { b } from \"./b\";export const d = b;",
        "signature": "cbe8d1524c57b7913aea74050cf9fabb-export declare const d = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./a.ts"
    ],
    [
      "./b.ts"
    ]
  ],
  "options": {
    "declaration": true
  },
  "referencedMap": {
    "./c.ts": [
      "./a.ts"
    ],
    "./d.ts": [
      "./b.ts"
    ]
  },
  "size": 1968
}
//// [/home/src/workspaces/solution/project2/src/e.js] *rewrite with same content*
//// [/home/src/workspaces/solution/project2/src/f.js] *rewrite with same content*
//// [/home/src/workspaces/solution/project2/src/g.d.ts] *rewrite with same content*
//// [/home/src/workspaces/solution/project2/src/g.js] *rewrite with same content*
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2,4,6],"fileNames":["../../../../tslibs/TS/Lib/lib.d.ts","./e.ts","../../project1/src/a.d.ts","./f.ts","../../project1/src/b.d.ts","./g.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"f0e6085d8cf835a263334b8d65c348d6-export const e = 10;","signature":"03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n","impliedNodeFormat":1},"a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",{"version":"7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;","signature":"b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n","impliedNodeFormat":1},"195d8237d70e02ed04ff98ea3f31f6a5-export declare const b = 10;\nexport declare const aaaaa = 10;\nexport declare const a2 = 10;\n",{"version":"8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;","signature":"0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n","impliedNodeFormat":1}],"fileIdsList":[[3],[5]],"options":{"declaration":true},"referencedMap":[[4,1],[6,2]]}
//// [/home/src/workspaces/solution/project2/src/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./e.ts"
      ],
      "original": 2
    },
    {
      "files": [
        "./f.ts"
      ],
      "original": 4
    },
    {
      "files": [
        "./g.ts"
      ],
      "original": 6
    }
  ],
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.d.ts",
    "./e.ts",
    "../../project1/src/a.d.ts",
    "./f.ts",
    "../../project1/src/b.d.ts",
    "./g.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./e.ts",
      "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
      "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f0e6085d8cf835a263334b8d65c348d6-export const e = 10;",
        "signature": "03228703b057b7967d3b34e3e293afe6-export declare const e = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/a.d.ts",
      "version": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "signature": "a63a12b1c0db30a4f4db8fab30590511-export declare const a = 10;\nexport declare const aaa = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./f.ts",
      "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
      "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7642b795d1c8af619b85a871556ad795-import { a } from \"../../project1/src/a\"; export const f = a;",
        "signature": "b32d51b52857b786b7f8bc9da8254719-export declare const f = 10;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../../project1/src/b.d.ts",
      "version": "195d8237d70e02ed04ff98ea3f31f6a5-export declare const b = 10;\nexport declare const aaaaa = 10;\nexport declare const a2 = 10;\n",
      "signature": "195d8237d70e02ed04ff98ea3f31f6a5-export declare const b = 10;\nexport declare const aaaaa = 10;\nexport declare const a2 = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./g.ts",
      "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
      "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8ac6555c89cc836d9f2e654843fd5d85-import { b } from \"../../project1/src/b\"; export const g = b;",
        "signature": "0cd91ad43457b78bfa48c8f690fb287c-export declare const g = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../../project1/src/a.d.ts"
    ],
    [
      "../../project1/src/b.d.ts"
    ]
  ],
  "options": {
    "declaration": true
  },
  "referencedMap": {
    "./f.ts": [
      "../../project1/src/a.d.ts"
    ],
    "./g.ts": [
      "../../project1/src/b.d.ts"
    ]
  },
  "size": 1864
}

/home/src/workspaces/solution/project1/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/b.ts
*refresh*    /home/src/workspaces/solution/project1/src/d.ts
Signatures::
(computed .d.ts) /home/src/workspaces/solution/project1/src/b.ts
(computed .d.ts) /home/src/workspaces/solution/project1/src/d.ts

/home/src/workspaces/solution/project2/src/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/solution/project1/src/b.d.ts
*refresh*    /home/src/workspaces/solution/project2/src/g.ts
Signatures::
(used version)   /home/src/workspaces/solution/project1/src/b.d.ts
(computed .d.ts) /home/src/workspaces/solution/project2/src/g.ts


Diff:: Verbose output status will be different because of up-to-date-ness checks
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -2,11 +2,11 @@
     * project1/src/tsconfig.json
     * project2/src/tsconfig.json

-[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because output file 'project1/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project1/src/tsconfig.json' is out of date because buildinfo file 'project1/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

 [[90mHH:MM:SS AM[0m] Building project 'project1/src/tsconfig.json'...

-[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because output file 'project2/src/tsconfig.tsbuildinfo' does not exist
+[[90mHH:MM:SS AM[0m] Project 'project2/src/tsconfig.json' is out of date because buildinfo file 'project2/src/tsconfig.tsbuildinfo' indicates there is change in compilerOptions

 [[90mHH:MM:SS AM[0m] Building project 'project2/src/tsconfig.json'...

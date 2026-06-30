currentDirectory::/user/username/projects/facade
useCaseSensitiveFileNames::true
Input::
//// [/user/username/projects/facade/consumer/index.ts] *new* 
import type { Status } from "../facade/index";

export const status: "old" = null as unknown as Status;
//// [/user/username/projects/facade/consumer/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../facade" }],
}
//// [/user/username/projects/facade/facade/index.ts] *new* 
export * from "../source/index";
//// [/user/username/projects/facade/facade/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../source" }],
}
//// [/user/username/projects/facade/source/index.ts] *new* 
export type Status = "old";
//// [/user/username/projects/facade/source/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
}

tsgo -b consumer --verbose
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * source/tsconfig.json
    * facade/tsconfig.json
    * consumer/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'source/tsconfig.json' is out of date because output file 'source/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'source/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'facade/tsconfig.json' is out of date because output file 'facade/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'facade/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'consumer/tsconfig.json' is out of date because output file 'consumer/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'consumer/tsconfig.json'...

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
//// [/user/username/projects/facade/consumer/index.d.ts] *new* 
export declare const status: "old";

//// [/user/username/projects/facade/consumer/index.js] *new* 
export const status = null;

//// [/user/username/projects/facade/consumer/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[4],"fileNames":["lib.es2025.full.d.ts","../source/index.d.ts","../facade/index.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n","7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",{"version":"f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;","signature":"69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n","impliedNodeFormat":1}],"fileIdsList":[[3],[2]],"options":{"composite":true},"referencedMap":[[4,1],[3,2]],"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/facade/consumer/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 4
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../source/index.d.ts",
    "../facade/index.d.ts",
    "./index.ts"
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
      "fileName": "../source/index.d.ts",
      "version": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../facade/index.d.ts",
      "version": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./index.ts",
      "version": "f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;",
      "signature": "69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;",
        "signature": "69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../facade/index.d.ts"
    ],
    [
      "../source/index.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./index.ts": [
      "../facade/index.d.ts"
    ],
    "../facade/index.d.ts": [
      "../source/index.d.ts"
    ]
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1444
}
//// [/user/username/projects/facade/facade/index.d.ts] *new* 
export * from "../source/index";

//// [/user/username/projects/facade/facade/index.js] *new* 
export * from "../source/index";

//// [/user/username/projects/facade/facade/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[3],"fileNames":["lib.es2025.full.d.ts","../source/index.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",{"version":"e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";","signature":"7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"composite":true},"referencedMap":[[3,1]],"latestChangedDtsFile":"./index.d.ts","emitSignatures":[[3,"158d569ea17ecd4dcfd73f9a3b1a115f-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n"]]}
//// [/user/username/projects/facade/facade/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 3
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../source/index.d.ts",
    "./index.ts"
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
      "fileName": "../source/index.d.ts",
      "version": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./index.ts",
      "version": "e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";",
      "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";",
        "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../source/index.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./index.ts": [
      "../source/index.d.ts"
    ]
  },
  "latestChangedDtsFile": "./index.d.ts",
  "emitSignatures": [
    {
      "file": "./index.ts",
      "signature": "158d569ea17ecd4dcfd73f9a3b1a115f-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "original": [
        3,
        "158d569ea17ecd4dcfd73f9a3b1a115f-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n"
      ]
    }
  ],
  "size": 1471
}
//// [/user/username/projects/facade/source/index.d.ts] *new* 
export type Status = "old";

//// [/user/username/projects/facade/source/index.js] *new* 
export {};

//// [/user/username/projects/facade/source/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"805cb8841c400106fada7ec033d625cc-export type Status = \"old\";","signature":"5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/facade/source/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 2
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./index.ts"
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
      "fileName": "./index.ts",
      "version": "805cb8841c400106fada7ec033d625cc-export type Status = \"old\";",
      "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "805cb8841c400106fada7ec033d625cc-export type Status = \"old\";",
        "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1117
}

source/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/facade/source/index.ts
Signatures::
(stored at emit) /user/username/projects/facade/source/index.ts

facade/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/facade/source/index.d.ts
*refresh*    /user/username/projects/facade/facade/index.ts
Signatures::
(stored at emit) /user/username/projects/facade/facade/index.ts

consumer/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/facade/source/index.d.ts
*refresh*    /user/username/projects/facade/facade/index.d.ts
*refresh*    /user/username/projects/facade/consumer/index.ts
Signatures::
(stored at emit) /user/username/projects/facade/consumer/index.ts


Edit [0]:: Change re-exported type in source project
//// [/user/username/projects/facade/source/index.ts] *modified* 
export type Status = "new";

tsgo -b consumer --verbose
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * source/tsconfig.json
    * facade/tsconfig.json
    * consumer/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'source/tsconfig.json' is out of date because output 'source/tsconfig.tsbuildinfo' is older than input 'source/index.ts'

[[90mHH:MM:SS AM[0m] Building project 'source/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'facade/tsconfig.json' is out of date because output 'facade/tsconfig.tsbuildinfo' is older than input 'source/index.d.ts'

[[90mHH:MM:SS AM[0m] Building project 'facade/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'consumer/tsconfig.json' is out of date because output 'consumer/tsconfig.tsbuildinfo' is older than input 'source/index.d.ts'

[[90mHH:MM:SS AM[0m] Building project 'consumer/tsconfig.json'...

[96mconsumer/index.ts[0m:[93m3[0m:[93m14[0m - [91merror[0m[90m TS2322: [0mType '"new"' is not assignable to type '"old"'.

[7m3[0m export const status: "old" = null as unknown as Status;
[7m [0m [91m             ~~~~~~[0m


Found 1 error in consumer/index.ts[90m:3[0m

//// [/user/username/projects/facade/consumer/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[4],"fileNames":["lib.es2025.full.d.ts","../source/index.d.ts","../facade/index.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n","7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",{"version":"f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;","signature":"69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n","impliedNodeFormat":1}],"fileIdsList":[[3],[2]],"options":{"composite":true},"referencedMap":[[4,1],[3,2]],"semanticDiagnosticsPerFile":[[4,[{"pos":61,"end":67,"code":2322,"category":1,"messageKey":"Type_0_is_not_assignable_to_type_1_2322","messageArgs":["\"new\"","\"old\""]}]]],"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/facade/consumer/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 4
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../source/index.d.ts",
    "../facade/index.d.ts",
    "./index.ts"
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
      "fileName": "../source/index.d.ts",
      "version": "9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",
      "signature": "9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../facade/index.d.ts",
      "version": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./index.ts",
      "version": "f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;",
      "signature": "69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;",
        "signature": "69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../facade/index.d.ts"
    ],
    [
      "../source/index.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./index.ts": [
      "../facade/index.d.ts"
    ],
    "../facade/index.d.ts": [
      "../source/index.d.ts"
    ]
  },
  "semanticDiagnosticsPerFile": [
    [
      "./index.ts",
      [
        {
          "pos": 61,
          "end": 67,
          "code": 2322,
          "category": 1,
          "messageKey": "Type_0_is_not_assignable_to_type_1_2322",
          "messageArgs": [
            "\"new\"",
            "\"old\""
          ]
        }
      ]
    ]
  ],
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1617
}
//// [/user/username/projects/facade/facade/index.d.ts] *rewrite with same content*
//// [/user/username/projects/facade/facade/index.js] *rewrite with same content*
//// [/user/username/projects/facade/facade/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[3],"fileNames":["lib.es2025.full.d.ts","../source/index.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",{"version":"e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";","signature":"7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"composite":true},"referencedMap":[[3,1]],"latestChangedDtsFile":"./index.d.ts","emitSignatures":[[3,"aa26532ed10dd2bca3829448f2ca8eb3-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n"]]}
//// [/user/username/projects/facade/facade/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 3
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../source/index.d.ts",
    "./index.ts"
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
      "fileName": "../source/index.d.ts",
      "version": "9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",
      "signature": "9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./index.ts",
      "version": "e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";",
      "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";",
        "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../source/index.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./index.ts": [
      "../source/index.d.ts"
    ]
  },
  "latestChangedDtsFile": "./index.d.ts",
  "emitSignatures": [
    {
      "file": "./index.ts",
      "signature": "aa26532ed10dd2bca3829448f2ca8eb3-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",
      "original": [
        3,
        "aa26532ed10dd2bca3829448f2ca8eb3-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n"
      ]
    }
  ],
  "size": 1471
}
//// [/user/username/projects/facade/source/index.d.ts] *modified* 
export type Status = "new";

//// [/user/username/projects/facade/source/index.js] *rewrite with same content*
//// [/user/username/projects/facade/source/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"bbd3073d7cd39697a145b49e4798b14a-export type Status = \"new\";","signature":"9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/facade/source/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 2
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./index.ts"
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
      "fileName": "./index.ts",
      "version": "bbd3073d7cd39697a145b49e4798b14a-export type Status = \"new\";",
      "signature": "9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "bbd3073d7cd39697a145b49e4798b14a-export type Status = \"new\";",
        "signature": "9b0886d42f77b756f7e196fff7faf468-export type Status = \"new\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1117
}

source/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/facade/source/index.ts
Signatures::
(computed .d.ts) /user/username/projects/facade/source/index.ts

facade/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/facade/source/index.d.ts
*refresh*    /user/username/projects/facade/facade/index.ts
Signatures::
(used version)   /user/username/projects/facade/source/index.d.ts
(computed .d.ts) /user/username/projects/facade/facade/index.ts

consumer/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/facade/source/index.d.ts
*refresh*    /user/username/projects/facade/facade/index.d.ts
*refresh*    /user/username/projects/facade/consumer/index.ts
Signatures::
(used version)   /user/username/projects/facade/source/index.d.ts
(used version)   /user/username/projects/facade/facade/index.d.ts
(stored at emit) /user/username/projects/facade/consumer/index.ts


Edit [1]:: Revert re-exported type in source project
//// [/user/username/projects/facade/source/index.ts] *modified* 
export type Status = "old";

tsgo -b consumer --verbose
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * source/tsconfig.json
    * facade/tsconfig.json
    * consumer/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'source/tsconfig.json' is out of date because output 'source/tsconfig.tsbuildinfo' is older than input 'source/index.ts'

[[90mHH:MM:SS AM[0m] Building project 'source/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'facade/tsconfig.json' is out of date because output 'facade/tsconfig.tsbuildinfo' is older than input 'source/index.d.ts'

[[90mHH:MM:SS AM[0m] Building project 'facade/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'consumer/tsconfig.json' is out of date because buildinfo file 'consumer/tsconfig.tsbuildinfo' indicates that program needs to report errors.

[[90mHH:MM:SS AM[0m] Building project 'consumer/tsconfig.json'...

//// [/user/username/projects/facade/consumer/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[4],"fileNames":["lib.es2025.full.d.ts","../source/index.d.ts","../facade/index.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n","7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",{"version":"f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;","signature":"69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n","impliedNodeFormat":1}],"fileIdsList":[[3],[2]],"options":{"composite":true},"referencedMap":[[4,1],[3,2]],"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/facade/consumer/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 4
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../source/index.d.ts",
    "../facade/index.d.ts",
    "./index.ts"
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
      "fileName": "../source/index.d.ts",
      "version": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../facade/index.d.ts",
      "version": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./index.ts",
      "version": "f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;",
      "signature": "69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f4b394d2e4781ef1ae46991d46f2bf7e-import type { Status } from \"../facade/index\";\n\nexport const status: \"old\" = null as unknown as Status;",
        "signature": "69bc43828c83e73e7ba3bd205e41d031-export declare const status: \"old\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../facade/index.d.ts"
    ],
    [
      "../source/index.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./index.ts": [
      "../facade/index.d.ts"
    ],
    "../facade/index.d.ts": [
      "../source/index.d.ts"
    ]
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1444
}
//// [/user/username/projects/facade/facade/index.d.ts] *rewrite with same content*
//// [/user/username/projects/facade/facade/index.js] *rewrite with same content*
//// [/user/username/projects/facade/facade/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[3],"fileNames":["lib.es2025.full.d.ts","../source/index.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",{"version":"e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";","signature":"7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"composite":true},"referencedMap":[[3,1]],"latestChangedDtsFile":"./index.d.ts","emitSignatures":[[3,"158d569ea17ecd4dcfd73f9a3b1a115f-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n"]]}
//// [/user/username/projects/facade/facade/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 3
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "../source/index.d.ts",
    "./index.ts"
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
      "fileName": "../source/index.d.ts",
      "version": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./index.ts",
      "version": "e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";",
      "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e6a451c08997b100385569b7eaa8927c-export * from \"../source/index\";",
        "signature": "7c2372fe894981acfb34a06a40fa7323-export * from \"../source/index\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../source/index.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./index.ts": [
      "../source/index.d.ts"
    ]
  },
  "latestChangedDtsFile": "./index.d.ts",
  "emitSignatures": [
    {
      "file": "./index.ts",
      "signature": "158d569ea17ecd4dcfd73f9a3b1a115f-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "original": [
        3,
        "158d569ea17ecd4dcfd73f9a3b1a115f-export * from \"../source/index\";\n\n/user/username/projects/facade/source/index.d.ts=5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n"
      ]
    }
  ],
  "size": 1471
}
//// [/user/username/projects/facade/source/index.d.ts] *modified* 
export type Status = "old";

//// [/user/username/projects/facade/source/index.js] *rewrite with same content*
//// [/user/username/projects/facade/source/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"805cb8841c400106fada7ec033d625cc-export type Status = \"old\";","signature":"5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/facade/source/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.ts"
      ],
      "original": 2
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./index.ts"
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
      "fileName": "./index.ts",
      "version": "805cb8841c400106fada7ec033d625cc-export type Status = \"old\";",
      "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "805cb8841c400106fada7ec033d625cc-export type Status = \"old\";",
        "signature": "5a09cbbef2c0e6f8324b2a0c72b092db-export type Status = \"old\";\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1117
}

source/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/facade/source/index.ts
Signatures::
(computed .d.ts) /user/username/projects/facade/source/index.ts

facade/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/facade/source/index.d.ts
*refresh*    /user/username/projects/facade/facade/index.ts
Signatures::
(used version)   /user/username/projects/facade/source/index.d.ts
(computed .d.ts) /user/username/projects/facade/facade/index.ts

consumer/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/facade/source/index.d.ts
*refresh*    /user/username/projects/facade/facade/index.d.ts
*refresh*    /user/username/projects/facade/consumer/index.ts
Signatures::
(used version)   /user/username/projects/facade/source/index.d.ts
(used version)   /user/username/projects/facade/facade/index.d.ts
(stored at emit) /user/username/projects/facade/consumer/index.ts

currentDirectory::/user/username/projects/myproject
useCaseSensitiveFileNames::true
Input::
//// [/user/username/projects/myproject/pkg0/index.ts] *new* 
export const pkg0 = 0;
//// [/user/username/projects/myproject/pkg0/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },

}
//// [/user/username/projects/myproject/pkg1/index.ts] *new* 
export const pkg1 = 1;
//// [/user/username/projects/myproject/pkg1/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../pkg0" }],
}
//// [/user/username/projects/myproject/pkg2/index.ts] *new* 
export const pkg2 = 2;
//// [/user/username/projects/myproject/pkg2/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../pkg0" }],
}
//// [/user/username/projects/myproject/pkg3/index.ts] *new* 
export const pkg3 = 3;
//// [/user/username/projects/myproject/pkg3/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../pkg0" }],
}
//// [/user/username/projects/myproject/pkg4/index.ts] *new* 
export const pkg4 = 4;
//// [/user/username/projects/myproject/pkg4/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../pkg0" }],
}
//// [/user/username/projects/myproject/pkg5/index.ts] *new* 
export const pkg5 = 5;
//// [/user/username/projects/myproject/pkg5/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../pkg0" }],
}
//// [/user/username/projects/myproject/pkg6/index.ts] *new* 
export const pkg6 = 6;
//// [/user/username/projects/myproject/pkg6/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../pkg0" }],
}
//// [/user/username/projects/myproject/pkg7/index.ts] *new* 
export const pkg7 = 7;
//// [/user/username/projects/myproject/pkg7/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [{ "path": "../pkg0" }],
}
//// [/user/username/projects/myproject/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "references": [
        { "path": "./pkg0" },
        { "path": "./pkg1" },
        { "path": "./pkg2" },
        { "path": "./pkg3" },
        { "path": "./pkg4" },
        { "path": "./pkg5" },
        { "path": "./pkg6" },
        { "path": "./pkg7" }
    ]
}

tsgo -b -w -v --builders 3
ExitStatus:: Success
Output::
[2J[3J[H[[90mHH:MM:SS AM[0m] Starting compilation in watch mode...

[[90mHH:MM:SS AM[0m] Projects in this build: 
    * pkg0/tsconfig.json
    * pkg1/tsconfig.json
    * pkg2/tsconfig.json
    * pkg3/tsconfig.json
    * pkg4/tsconfig.json
    * pkg5/tsconfig.json
    * pkg6/tsconfig.json
    * pkg7/tsconfig.json
    * tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'pkg0/tsconfig.json' is out of date because output file 'pkg0/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg0/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg1/tsconfig.json' is out of date because output file 'pkg1/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg1/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg2/tsconfig.json' is out of date because output file 'pkg2/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg2/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg3/tsconfig.json' is out of date because output file 'pkg3/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg3/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg4/tsconfig.json' is out of date because output file 'pkg4/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg4/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg5/tsconfig.json' is out of date because output file 'pkg5/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg5/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg6/tsconfig.json' is out of date because output file 'pkg6/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg6/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg7/tsconfig.json' is out of date because output file 'pkg7/tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'pkg7/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'tsconfig.json' is out of date because output file 'tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'tsconfig.json'...

[[90mHH:MM:SS AM[0m] Found 0 errors. Watching for file changes.

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
//// [/user/username/projects/myproject/pkg0/index.d.ts] *new* 
export const pkg0 = 0;

//// [/user/username/projects/myproject/pkg0/index.js] *new* 
export const pkg0 = 0;

//// [/user/username/projects/myproject/pkg0/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"44b96319d3533fefbaaf61cad5d90c48-export const pkg0 = 0;","signature":"6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg0/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "44b96319d3533fefbaaf61cad5d90c48-export const pkg0 = 0;",
      "signature": "6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "44b96319d3533fefbaaf61cad5d90c48-export const pkg0 = 0;",
        "signature": "6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/pkg1/index.d.ts] *new* 
export const pkg1 = 1;

//// [/user/username/projects/myproject/pkg1/index.js] *new* 
export const pkg1 = 1;

//// [/user/username/projects/myproject/pkg1/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"dd791c6b74e4cf9af9283579215cad88-export const pkg1 = 1;","signature":"50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg1/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "dd791c6b74e4cf9af9283579215cad88-export const pkg1 = 1;",
      "signature": "50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "dd791c6b74e4cf9af9283579215cad88-export const pkg1 = 1;",
        "signature": "50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/pkg2/index.d.ts] *new* 
export const pkg2 = 2;

//// [/user/username/projects/myproject/pkg2/index.js] *new* 
export const pkg2 = 2;

//// [/user/username/projects/myproject/pkg2/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0cf39c01e85273ecae99a05645f3b18b-export const pkg2 = 2;","signature":"bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg2/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "0cf39c01e85273ecae99a05645f3b18b-export const pkg2 = 2;",
      "signature": "bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0cf39c01e85273ecae99a05645f3b18b-export const pkg2 = 2;",
        "signature": "bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/pkg3/index.d.ts] *new* 
export const pkg3 = 3;

//// [/user/username/projects/myproject/pkg3/index.js] *new* 
export const pkg3 = 3;

//// [/user/username/projects/myproject/pkg3/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"3e7c6d2a2c682228b50c62e065dd86eb-export const pkg3 = 3;","signature":"fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg3/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "3e7c6d2a2c682228b50c62e065dd86eb-export const pkg3 = 3;",
      "signature": "fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "3e7c6d2a2c682228b50c62e065dd86eb-export const pkg3 = 3;",
        "signature": "fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/pkg4/index.d.ts] *new* 
export const pkg4 = 4;

//// [/user/username/projects/myproject/pkg4/index.js] *new* 
export const pkg4 = 4;

//// [/user/username/projects/myproject/pkg4/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"83c6d6f47b9957960cfb2d779ba0fcf4-export const pkg4 = 4;","signature":"d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg4/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "83c6d6f47b9957960cfb2d779ba0fcf4-export const pkg4 = 4;",
      "signature": "d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "83c6d6f47b9957960cfb2d779ba0fcf4-export const pkg4 = 4;",
        "signature": "d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/pkg5/index.d.ts] *new* 
export const pkg5 = 5;

//// [/user/username/projects/myproject/pkg5/index.js] *new* 
export const pkg5 = 5;

//// [/user/username/projects/myproject/pkg5/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"c2daf00ad21fc855a340358c8c204b48-export const pkg5 = 5;","signature":"ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg5/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "c2daf00ad21fc855a340358c8c204b48-export const pkg5 = 5;",
      "signature": "ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "c2daf00ad21fc855a340358c8c204b48-export const pkg5 = 5;",
        "signature": "ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/pkg6/index.d.ts] *new* 
export const pkg6 = 6;

//// [/user/username/projects/myproject/pkg6/index.js] *new* 
export const pkg6 = 6;

//// [/user/username/projects/myproject/pkg6/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"ca770928977c270ca6068d9ec3f4ed53-export const pkg6 = 6;","signature":"0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg6/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "ca770928977c270ca6068d9ec3f4ed53-export const pkg6 = 6;",
      "signature": "0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "ca770928977c270ca6068d9ec3f4ed53-export const pkg6 = 6;",
        "signature": "0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/pkg7/index.d.ts] *new* 
export const pkg7 = 7;

//// [/user/username/projects/myproject/pkg7/index.js] *new* 
export const pkg7 = 7;

//// [/user/username/projects/myproject/pkg7/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"d445d242c6358b7d6b7a2989cd41e84a-export const pkg7 = 7;","signature":"d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg7/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
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
      "version": "d445d242c6358b7d6b7a2989cd41e84a-export const pkg7 = 7;",
      "signature": "d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d445d242c6358b7d6b7a2989cd41e84a-export const pkg7 = 7;",
        "signature": "d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1103
}
//// [/user/username/projects/myproject/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,9]],"fileNames":["lib.es2025.full.d.ts","./pkg0/index.d.ts","./pkg1/index.d.ts","./pkg2/index.d.ts","./pkg3/index.d.ts","./pkg4/index.d.ts","./pkg5/index.d.ts","./pkg6/index.d.ts","./pkg7/index.d.ts","./pkg0/index.ts","./pkg1/index.ts","./pkg2/index.ts","./pkg3/index.ts","./pkg4/index.ts","./pkg5/index.ts","./pkg6/index.ts","./pkg7/index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n","50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n","bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n","fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n","d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n","ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n","0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n","d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n"],"options":{"composite":true},"resolvedRoot":[[2,10],[3,11],[4,12],[5,13],[6,14],[7,15],[8,16],[9,17]]}
//// [/user/username/projects/myproject/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./pkg0/index.d.ts",
        "./pkg1/index.d.ts",
        "./pkg2/index.d.ts",
        "./pkg3/index.d.ts",
        "./pkg4/index.d.ts",
        "./pkg5/index.d.ts",
        "./pkg6/index.d.ts",
        "./pkg7/index.d.ts"
      ],
      "original": [
        2,
        9
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./pkg0/index.d.ts",
    "./pkg1/index.d.ts",
    "./pkg2/index.d.ts",
    "./pkg3/index.d.ts",
    "./pkg4/index.d.ts",
    "./pkg5/index.d.ts",
    "./pkg6/index.d.ts",
    "./pkg7/index.d.ts",
    "./pkg0/index.ts",
    "./pkg1/index.ts",
    "./pkg2/index.ts",
    "./pkg3/index.ts",
    "./pkg4/index.ts",
    "./pkg5/index.ts",
    "./pkg6/index.ts",
    "./pkg7/index.ts"
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
      "fileName": "./pkg0/index.d.ts",
      "version": "6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n",
      "signature": "6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg1/index.d.ts",
      "version": "50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n",
      "signature": "50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg2/index.d.ts",
      "version": "bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n",
      "signature": "bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg3/index.d.ts",
      "version": "fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n",
      "signature": "fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg4/index.d.ts",
      "version": "d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n",
      "signature": "d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg5/index.d.ts",
      "version": "ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n",
      "signature": "ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg6/index.d.ts",
      "version": "0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n",
      "signature": "0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg7/index.d.ts",
      "version": "d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n",
      "signature": "d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "options": {
    "composite": true
  },
  "resolvedRoot": [
    [
      "./pkg0/index.d.ts",
      "./pkg0/index.ts"
    ],
    [
      "./pkg1/index.d.ts",
      "./pkg1/index.ts"
    ],
    [
      "./pkg2/index.d.ts",
      "./pkg2/index.ts"
    ],
    [
      "./pkg3/index.d.ts",
      "./pkg3/index.ts"
    ],
    [
      "./pkg4/index.d.ts",
      "./pkg4/index.ts"
    ],
    [
      "./pkg5/index.d.ts",
      "./pkg5/index.ts"
    ],
    [
      "./pkg6/index.d.ts",
      "./pkg6/index.ts"
    ],
    [
      "./pkg7/index.d.ts",
      "./pkg7/index.ts"
    ]
  ],
  "size": 1749
}

pkg0/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg0/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg0/index.ts

pkg1/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg1/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg1/index.ts

pkg2/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg2/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg2/index.ts

pkg3/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg3/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg3/index.ts

pkg4/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg4/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg4/index.ts

pkg5/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg5/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg5/index.ts

pkg6/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg6/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg6/index.ts

pkg7/tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg7/index.ts
Signatures::
(stored at emit) /user/username/projects/myproject/pkg7/index.ts

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /user/username/projects/myproject/pkg0/index.d.ts
*refresh*    /user/username/projects/myproject/pkg1/index.d.ts
*refresh*    /user/username/projects/myproject/pkg2/index.d.ts
*refresh*    /user/username/projects/myproject/pkg3/index.d.ts
*refresh*    /user/username/projects/myproject/pkg4/index.d.ts
*refresh*    /user/username/projects/myproject/pkg5/index.d.ts
*refresh*    /user/username/projects/myproject/pkg6/index.d.ts
*refresh*    /user/username/projects/myproject/pkg7/index.d.ts
Signatures::


Edit [0]:: dts doesn't change
//// [/user/username/projects/myproject/pkg0/index.ts] *modified* 
export const pkg0 = 0;const someConst2 = 10;


Output::
[2J[3J[H[[90mHH:MM:SS AM[0m] File change detected. Starting incremental compilation...

[[90mHH:MM:SS AM[0m] Projects in this build: 
    * pkg0/tsconfig.json
    * pkg1/tsconfig.json
    * pkg2/tsconfig.json
    * pkg3/tsconfig.json
    * pkg4/tsconfig.json
    * pkg5/tsconfig.json
    * pkg6/tsconfig.json
    * pkg7/tsconfig.json
    * tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'pkg0/tsconfig.json' is out of date because output 'pkg0/tsconfig.tsbuildinfo' is older than input 'pkg0/index.ts'

[[90mHH:MM:SS AM[0m] Building project 'pkg0/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg1/tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'pkg1/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg2/tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'pkg2/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg3/tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'pkg3/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg4/tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'pkg4/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg5/tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'pkg5/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg6/tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'pkg6/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg7/tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'pkg7/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'tsconfig.json' is up to date with .d.ts files from its dependencies

[[90mHH:MM:SS AM[0m] Updating output timestamps of project 'tsconfig.json'...

[[90mHH:MM:SS AM[0m] Found 0 errors. Watching for file changes.

//// [/user/username/projects/myproject/pkg0/index.js] *modified* 
export const pkg0 = 0;
const someConst2 = 10;

//// [/user/username/projects/myproject/pkg0/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"87e698a84292a7346d74cb79ef9ee973-export const pkg0 = 0;const someConst2 = 10;","signature":"6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg0/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
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
      "version": "87e698a84292a7346d74cb79ef9ee973-export const pkg0 = 0;const someConst2 = 10;",
      "signature": "6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "87e698a84292a7346d74cb79ef9ee973-export const pkg0 = 0;const someConst2 = 10;",
        "signature": "6acd84c9889b3f2cdeec5b209e9e475b-export const pkg0 = 0;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1125
}
//// [/user/username/projects/myproject/pkg1/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg2/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg3/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg4/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg5/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg6/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg7/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/tsconfig.tsbuildinfo] *mTime changed*

pkg0/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/myproject/pkg0/index.ts
Signatures::
(computed .d.ts) /user/username/projects/myproject/pkg0/index.ts


Edit [1]:: no change


Output::



Edit [2]:: dts change
//// [/user/username/projects/myproject/pkg0/index.ts] *modified* 
export const pkg0 = 0;const someConst2 = 10;export const someConst = 10;


Output::
[2J[3J[H[[90mHH:MM:SS AM[0m] File change detected. Starting incremental compilation...

[[90mHH:MM:SS AM[0m] Projects in this build: 
    * pkg0/tsconfig.json
    * pkg1/tsconfig.json
    * pkg2/tsconfig.json
    * pkg3/tsconfig.json
    * pkg4/tsconfig.json
    * pkg5/tsconfig.json
    * pkg6/tsconfig.json
    * pkg7/tsconfig.json
    * tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'pkg0/tsconfig.json' is out of date because output 'pkg0/tsconfig.tsbuildinfo' is older than input 'pkg0/index.ts'

[[90mHH:MM:SS AM[0m] Building project 'pkg0/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg1/tsconfig.json' is out of date because output 'pkg1/index.js' is older than input 'pkg0/tsconfig.json'

[[90mHH:MM:SS AM[0m] Building project 'pkg1/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'pkg1/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg2/tsconfig.json' is out of date because output 'pkg2/index.js' is older than input 'pkg0/tsconfig.json'

[[90mHH:MM:SS AM[0m] Building project 'pkg2/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'pkg2/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg3/tsconfig.json' is out of date because output 'pkg3/index.js' is older than input 'pkg0/tsconfig.json'

[[90mHH:MM:SS AM[0m] Building project 'pkg3/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'pkg3/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg4/tsconfig.json' is out of date because output 'pkg4/index.js' is older than input 'pkg0/tsconfig.json'

[[90mHH:MM:SS AM[0m] Building project 'pkg4/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'pkg4/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg5/tsconfig.json' is out of date because output 'pkg5/index.js' is older than input 'pkg0/tsconfig.json'

[[90mHH:MM:SS AM[0m] Building project 'pkg5/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'pkg5/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg6/tsconfig.json' is out of date because output 'pkg6/index.js' is older than input 'pkg0/tsconfig.json'

[[90mHH:MM:SS AM[0m] Building project 'pkg6/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'pkg6/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'pkg7/tsconfig.json' is out of date because output 'pkg7/index.js' is older than input 'pkg0/tsconfig.json'

[[90mHH:MM:SS AM[0m] Building project 'pkg7/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Updating unchanged output timestamps of project 'pkg7/tsconfig.json'...

[[90mHH:MM:SS AM[0m] Project 'tsconfig.json' is out of date because output 'tsconfig.tsbuildinfo' is older than input 'pkg0/index.ts'

[[90mHH:MM:SS AM[0m] Building project 'tsconfig.json'...

[[90mHH:MM:SS AM[0m] Found 0 errors. Watching for file changes.

//// [/user/username/projects/myproject/pkg0/index.d.ts] *modified* 
export const pkg0 = 0;
export const someConst = 10;

//// [/user/username/projects/myproject/pkg0/index.js] *modified* 
export const pkg0 = 0;
const someConst2 = 10;
export const someConst = 10;

//// [/user/username/projects/myproject/pkg0/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"78d1f0ba95fa9081a8a366af27252e24-export const pkg0 = 0;const someConst2 = 10;export const someConst = 10;","signature":"d1d261972f7dc844736c08011a310c14-export const pkg0 = 0;\nexport const someConst = 10;\n","impliedNodeFormat":1}],"options":{"composite":true},"latestChangedDtsFile":"./index.d.ts"}
//// [/user/username/projects/myproject/pkg0/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
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
      "version": "78d1f0ba95fa9081a8a366af27252e24-export const pkg0 = 0;const someConst2 = 10;export const someConst = 10;",
      "signature": "d1d261972f7dc844736c08011a310c14-export const pkg0 = 0;\nexport const someConst = 10;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "78d1f0ba95fa9081a8a366af27252e24-export const pkg0 = 0;const someConst2 = 10;export const someConst = 10;",
        "signature": "d1d261972f7dc844736c08011a310c14-export const pkg0 = 0;\nexport const someConst = 10;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true
  },
  "latestChangedDtsFile": "./index.d.ts",
  "size": 1183
}
//// [/user/username/projects/myproject/pkg1/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg2/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg3/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg4/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg5/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg6/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/pkg7/tsconfig.tsbuildinfo] *mTime changed*
//// [/user/username/projects/myproject/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,9]],"fileNames":["lib.es2025.full.d.ts","./pkg0/index.d.ts","./pkg1/index.d.ts","./pkg2/index.d.ts","./pkg3/index.d.ts","./pkg4/index.d.ts","./pkg5/index.d.ts","./pkg6/index.d.ts","./pkg7/index.d.ts","./pkg0/index.ts","./pkg1/index.ts","./pkg2/index.ts","./pkg3/index.ts","./pkg4/index.ts","./pkg5/index.ts","./pkg6/index.ts","./pkg7/index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"d1d261972f7dc844736c08011a310c14-export const pkg0 = 0;\nexport const someConst = 10;\n","50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n","bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n","fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n","d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n","ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n","0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n","d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n"],"options":{"composite":true},"resolvedRoot":[[2,10],[3,11],[4,12],[5,13],[6,14],[7,15],[8,16],[9,17]]}
//// [/user/username/projects/myproject/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./pkg0/index.d.ts",
        "./pkg1/index.d.ts",
        "./pkg2/index.d.ts",
        "./pkg3/index.d.ts",
        "./pkg4/index.d.ts",
        "./pkg5/index.d.ts",
        "./pkg6/index.d.ts",
        "./pkg7/index.d.ts"
      ],
      "original": [
        2,
        9
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./pkg0/index.d.ts",
    "./pkg1/index.d.ts",
    "./pkg2/index.d.ts",
    "./pkg3/index.d.ts",
    "./pkg4/index.d.ts",
    "./pkg5/index.d.ts",
    "./pkg6/index.d.ts",
    "./pkg7/index.d.ts",
    "./pkg0/index.ts",
    "./pkg1/index.ts",
    "./pkg2/index.ts",
    "./pkg3/index.ts",
    "./pkg4/index.ts",
    "./pkg5/index.ts",
    "./pkg6/index.ts",
    "./pkg7/index.ts"
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
      "fileName": "./pkg0/index.d.ts",
      "version": "d1d261972f7dc844736c08011a310c14-export const pkg0 = 0;\nexport const someConst = 10;\n",
      "signature": "d1d261972f7dc844736c08011a310c14-export const pkg0 = 0;\nexport const someConst = 10;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg1/index.d.ts",
      "version": "50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n",
      "signature": "50fcd00ffeb770539c9a6dd979533700-export const pkg1 = 1;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg2/index.d.ts",
      "version": "bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n",
      "signature": "bd7cabc9c8189344fd1dfff12cef1d8d-export const pkg2 = 2;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg3/index.d.ts",
      "version": "fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n",
      "signature": "fb83be37ecc4355915ba3cd93d24a0c0-export const pkg3 = 3;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg4/index.d.ts",
      "version": "d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n",
      "signature": "d5ce0c39ab0f66ef460bc592c4aeff0d-export const pkg4 = 4;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg5/index.d.ts",
      "version": "ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n",
      "signature": "ddf6787cbff3b373c211028a3d91846a-export const pkg5 = 5;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg6/index.d.ts",
      "version": "0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n",
      "signature": "0bbfd9db1d98fbbe9043868218bd015b-export const pkg6 = 6;\n",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./pkg7/index.d.ts",
      "version": "d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n",
      "signature": "d031b9984f48c07cc2ae01df7854074c-export const pkg7 = 7;\n",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "options": {
    "composite": true
  },
  "resolvedRoot": [
    [
      "./pkg0/index.d.ts",
      "./pkg0/index.ts"
    ],
    [
      "./pkg1/index.d.ts",
      "./pkg1/index.ts"
    ],
    [
      "./pkg2/index.d.ts",
      "./pkg2/index.ts"
    ],
    [
      "./pkg3/index.d.ts",
      "./pkg3/index.ts"
    ],
    [
      "./pkg4/index.d.ts",
      "./pkg4/index.ts"
    ],
    [
      "./pkg5/index.d.ts",
      "./pkg5/index.ts"
    ],
    [
      "./pkg6/index.d.ts",
      "./pkg6/index.ts"
    ],
    [
      "./pkg7/index.d.ts",
      "./pkg7/index.ts"
    ]
  ],
  "size": 1779
}

pkg0/tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/myproject/pkg0/index.ts
Signatures::
(computed .d.ts) /user/username/projects/myproject/pkg0/index.ts

pkg1/tsconfig.json::
SemanticDiagnostics::
Signatures::

pkg2/tsconfig.json::
SemanticDiagnostics::
Signatures::

pkg3/tsconfig.json::
SemanticDiagnostics::
Signatures::

pkg4/tsconfig.json::
SemanticDiagnostics::
Signatures::

pkg5/tsconfig.json::
SemanticDiagnostics::
Signatures::

pkg6/tsconfig.json::
SemanticDiagnostics::
Signatures::

pkg7/tsconfig.json::
SemanticDiagnostics::
Signatures::

tsconfig.json::
SemanticDiagnostics::
*refresh*    /user/username/projects/myproject/pkg0/index.d.ts
Signatures::
(used version)   /user/username/projects/myproject/pkg0/index.d.ts


Edit [3]:: no change


Output::


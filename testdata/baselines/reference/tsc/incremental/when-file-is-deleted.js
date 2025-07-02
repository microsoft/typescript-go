
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/file1.ts] *new* 
export class  C { }
//// [/home/src/workspaces/project/file2.ts] *new* 
export class D { }
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "outDir": "outDir"
    }
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
//// [/home/src/workspaces/project/outDir/file1.d.ts] *new* 
export declare class C {
}

//// [/home/src/workspaces/project/outDir/file1.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = void 0;
class C {
}
exports.C = C;

//// [/home/src/workspaces/project/outDir/file2.d.ts] *new* 
export declare class D {
}

//// [/home/src/workspaces/project/outDir/file2.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.D = void 0;
class D {
}
exports.D = D;

//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../../tslibs/TS/Lib/lib.d.ts","../file1.ts","../file2.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e8197812b523db314f9f43881cab045172bec1a6893c27b62a52b128afbb19da","signature":"33031a47f740dde897da491c7c6ac0ef2224f9c807448ba058aadba8abd00433","impliedNodeFormat":1},{"version":"2d42470676839be6ca4923b34e799e3a318398eb2ff7c6273c676358d80093e6","signature":"f7f62800d2d53e363dcd48e24d95af396e4f0bbafc1713aca098f7644aeb0331","impliedNodeFormat":1}],"options":{"composite":true,"outDir":"./"},"latestChangedDtsFile":"./file2.d.ts"}
//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../../tslibs/TS/Lib/lib.d.ts",
    "../file1.ts",
    "../file2.ts"
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
      "fileName": "../file1.ts",
      "version": "e8197812b523db314f9f43881cab045172bec1a6893c27b62a52b128afbb19da",
      "signature": "33031a47f740dde897da491c7c6ac0ef2224f9c807448ba058aadba8abd00433",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e8197812b523db314f9f43881cab045172bec1a6893c27b62a52b128afbb19da",
        "signature": "33031a47f740dde897da491c7c6ac0ef2224f9c807448ba058aadba8abd00433",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../file2.ts",
      "version": "2d42470676839be6ca4923b34e799e3a318398eb2ff7c6273c676358d80093e6",
      "signature": "f7f62800d2d53e363dcd48e24d95af396e4f0bbafc1713aca098f7644aeb0331",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "2d42470676839be6ca4923b34e799e3a318398eb2ff7c6273c676358d80093e6",
        "signature": "f7f62800d2d53e363dcd48e24d95af396e4f0bbafc1713aca098f7644aeb0331",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true,
    "outDir": "./"
  },
  "latestChangedDtsFile": "./file2.d.ts",
  "size": 685
}


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/file1.ts
*refresh*    /home/src/workspaces/project/file2.ts

Signatures::
(stored at emit) /home/src/workspaces/project/file1.ts
(stored at emit) /home/src/workspaces/project/file2.ts


Edit:: delete file with imports
//// [/home/src/workspaces/project/file2.ts] *deleted*

ExitStatus:: 0
Output::
//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../../tslibs/TS/Lib/lib.d.ts","../file1.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e8197812b523db314f9f43881cab045172bec1a6893c27b62a52b128afbb19da","signature":"33031a47f740dde897da491c7c6ac0ef2224f9c807448ba058aadba8abd00433","impliedNodeFormat":1}],"options":{"composite":true,"outDir":"./"},"latestChangedDtsFile":"./file2.d.ts"}
//// [/home/src/workspaces/project/outDir/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../../tslibs/TS/Lib/lib.d.ts",
    "../file1.ts"
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
      "fileName": "../file1.ts",
      "version": "e8197812b523db314f9f43881cab045172bec1a6893c27b62a52b128afbb19da",
      "signature": "33031a47f740dde897da491c7c6ac0ef2224f9c807448ba058aadba8abd00433",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e8197812b523db314f9f43881cab045172bec1a6893c27b62a52b128afbb19da",
        "signature": "33031a47f740dde897da491c7c6ac0ef2224f9c807448ba058aadba8abd00433",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true,
    "outDir": "./"
  },
  "latestChangedDtsFile": "./file2.d.ts",
  "size": 491
}


SemanticDiagnostics::

Signatures::

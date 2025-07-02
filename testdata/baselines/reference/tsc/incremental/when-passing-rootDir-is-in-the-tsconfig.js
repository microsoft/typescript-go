
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/src/main.ts] *new* 
export const x = 10;
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "incremental": true,
        "outDir": "dist",
        "rootDir": "./"
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
//// [/home/src/workspaces/project/dist/src/main.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
exports.x = 10;

//// [/home/src/workspaces/project/dist/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../../tslibs/TS/Lib/lib.d.ts","../src/main.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"03da4d6a46cc7950ba861120c64b47c14bc80b3c64f47ef17b61cb454358afd6"],"options":{"outDir":"./","rootDir":".."}}
//// [/home/src/workspaces/project/dist/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../../tslibs/TS/Lib/lib.d.ts",
    "../src/main.ts"
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
      "fileName": "../src/main.ts",
      "version": "03da4d6a46cc7950ba861120c64b47c14bc80b3c64f47ef17b61cb454358afd6",
      "signature": "03da4d6a46cc7950ba861120c64b47c14bc80b3c64f47ef17b61cb454358afd6",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "options": {
    "outDir": "./",
    "rootDir": ".."
  },
  "size": 341
}


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/src/main.ts

Signatures::


Edit:: no change

ExitStatus:: 0
Output::


SemanticDiagnostics::

Signatures::

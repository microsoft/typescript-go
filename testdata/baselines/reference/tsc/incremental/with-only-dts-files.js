currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/src/another.d.ts] *new* 
export const y = 10;
//// [/home/src/workspaces/project/src/main.d.ts] *new* 
export const x = 10;
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{}

tsgo --incremental
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
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/another.d.ts","./src/main.d.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"4aa16e9a67a4820d1dc51507221b4c73b5626b3a759d79d7147ad4eabe37ef49","03da4d6a46cc7950ba861120c64b47c14bc80b3c64f47ef17b61cb454358afd6"]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/another.d.ts",
    "./src/main.d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
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
      "fileName": "./src/another.d.ts",
      "version": "4aa16e9a67a4820d1dc51507221b4c73b5626b3a759d79d7147ad4eabe37ef49",
      "signature": "4aa16e9a67a4820d1dc51507221b4c73b5626b3a759d79d7147ad4eabe37ef49",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/main.d.ts",
      "version": "03da4d6a46cc7950ba861120c64b47c14bc80b3c64f47ef17b61cb454358afd6",
      "signature": "03da4d6a46cc7950ba861120c64b47c14bc80b3c64f47ef17b61cb454358afd6",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "size": 386
}

SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/src/another.d.ts
*refresh*    /home/src/workspaces/project/src/main.d.ts
Signatures::


Edit [0]:: no change

tsgo --incremental
ExitStatus:: Success
Output::

SemanticDiagnostics::
Signatures::


Edit [1]:: modify d.ts file
//// [/home/src/workspaces/project/src/main.d.ts] *modified* 
export const x = 10;export const xy = 100;

tsgo --incremental
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/another.d.ts","./src/main.d.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"4aa16e9a67a4820d1dc51507221b4c73b5626b3a759d79d7147ad4eabe37ef49","a701af2196ad5afd5fe2fb1f6c1c9ad538b85b9e4c37d747738ecc7f5d609540"]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/another.d.ts",
    "./src/main.d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
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
      "fileName": "./src/another.d.ts",
      "version": "4aa16e9a67a4820d1dc51507221b4c73b5626b3a759d79d7147ad4eabe37ef49",
      "signature": "4aa16e9a67a4820d1dc51507221b4c73b5626b3a759d79d7147ad4eabe37ef49",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/main.d.ts",
      "version": "a701af2196ad5afd5fe2fb1f6c1c9ad538b85b9e4c37d747738ecc7f5d609540",
      "signature": "a701af2196ad5afd5fe2fb1f6c1c9ad538b85b9e4c37d747738ecc7f5d609540",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "size": 386
}

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/main.d.ts
Signatures::
(used version)   /home/src/workspaces/project/src/main.d.ts

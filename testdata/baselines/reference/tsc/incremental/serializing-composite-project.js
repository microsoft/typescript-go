currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/index.tsx] *new* 
export const a = 1;
//// [/home/src/workspaces/project/other.ts] *new* 
export const b = 2;
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "strict": true,
        "module": "esnext",
    },
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
//// [/home/src/workspaces/project/index.d.ts] *new* 
export declare const a = 1;

//// [/home/src/workspaces/project/index.js] *new* 
export const a = 1;

//// [/home/src/workspaces/project/other.d.ts] *new* 
export declare const b = 2;

//// [/home/src/workspaces/project/other.js] *new* 
export const b = 2;

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./index.tsx","./other.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"683314ed22112e8dea8095c8c6173afa2c61279f5fe07968ebe0e21fff16871d-export const a = 1;","signature":"f0f1286d442f3c09fa07d37db7d31755cb3761daed3c8008fbfce412770425c6-export declare const a = 1;\n","impliedNodeFormat":1},{"version":"34f0f66ce649a0df0d3d5bad537c3b867b11b2fbeb5eee37e1c75a795544a4ed-export const b = 2;","signature":"0fb7bb75ad82d403bd7ba1f151c0297ef6a9167d0039e90a9067289387a719a7-export declare const b = 2;\n","impliedNodeFormat":1}],"options":{"composite":true,"module":99,"strict":true},"latestChangedDtsFile":"./other.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./index.tsx",
    "./other.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
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
      "fileName": "./index.tsx",
      "version": "683314ed22112e8dea8095c8c6173afa2c61279f5fe07968ebe0e21fff16871d-export const a = 1;",
      "signature": "f0f1286d442f3c09fa07d37db7d31755cb3761daed3c8008fbfce412770425c6-export declare const a = 1;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "683314ed22112e8dea8095c8c6173afa2c61279f5fe07968ebe0e21fff16871d-export const a = 1;",
        "signature": "f0f1286d442f3c09fa07d37db7d31755cb3761daed3c8008fbfce412770425c6-export declare const a = 1;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./other.ts",
      "version": "34f0f66ce649a0df0d3d5bad537c3b867b11b2fbeb5eee37e1c75a795544a4ed-export const b = 2;",
      "signature": "0fb7bb75ad82d403bd7ba1f151c0297ef6a9167d0039e90a9067289387a719a7-export declare const b = 2;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "34f0f66ce649a0df0d3d5bad537c3b867b11b2fbeb5eee37e1c75a795544a4ed-export const b = 2;",
        "signature": "0fb7bb75ad82d403bd7ba1f151c0297ef6a9167d0039e90a9067289387a719a7-export declare const b = 2;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true,
    "module": 99,
    "strict": true
  },
  "latestChangedDtsFile": "./other.d.ts",
  "size": 1498
}

SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/index.tsx
*refresh*    /home/src/workspaces/project/other.ts
Signatures::
(stored at emit) /home/src/workspaces/project/index.tsx
(stored at emit) /home/src/workspaces/project/other.ts

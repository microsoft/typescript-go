currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/index.ts] *new* 
import { util } from "./lib/util";
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{}

tsgo --watch
ExitStatus:: Success
Output::
build starting at HH:MM:SS AM
[96mindex.ts[0m:[93m1[0m:[93m22[0m - [91merror[0m[90m TS2307: [0mCannot find module './lib/util' or its corresponding type declarations.

[7m1[0m import { util } from "./lib/util";
[7m [0m [91m                     ~~~~~~~~~~~~[0m


Found 1 error in index.ts[90m:1[0m

build finished in d.ddds
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
//// [/home/src/workspaces/project/index.js] *new* 
export {};


tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/index.ts
Signatures::


Edit [0]:: create directory and imported file
//// [/home/src/workspaces/project/lib/util.ts] *new* 
export const util = "hello";


Output::
build starting at HH:MM:SS AM
build finished in d.ddds
//// [/home/src/workspaces/project/index.js] *rewrite with same content*
//// [/home/src/workspaces/project/lib/util.js] *new* 
export const util = "hello";


tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/lib/util.ts
*refresh*    /home/src/workspaces/project/index.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/lib/util.ts
(computed .d.ts) /home/src/workspaces/project/index.ts

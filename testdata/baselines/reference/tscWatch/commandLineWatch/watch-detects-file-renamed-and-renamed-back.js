currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/helper.ts] *new* 
export const helper = 1;
//// [/home/src/workspaces/project/index.ts] *new* 
import { helper } from "./helper";
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{}

tsgo --watch
ExitStatus:: Success
Output::
build starting at HH:MM:SS AM
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
//// [/home/src/workspaces/project/helper.js] *new* 
export const helper = 1;

//// [/home/src/workspaces/project/index.js] *new* 
export {};


tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/helper.ts
*refresh*    /home/src/workspaces/project/index.ts
Signatures::


Edit [0]:: rename helper to helper2
//// [/home/src/workspaces/project/helper.ts] *deleted*
//// [/home/src/workspaces/project/helper2.ts] *new* 
export const helper = 1;


Output::
build starting at HH:MM:SS AM
[96mindex.ts[0m:[93m1[0m:[93m24[0m - [91merror[0m[90m TS7016: [0mCould not find a declaration file for module './helper'. '/home/src/workspaces/project/helper.js' implicitly has an 'any' type.

[7m1[0m import { helper } from "./helper";
[7m [0m [91m                       ~~~~~~~~~~[0m


Found 1 error in index.ts[90m:1[0m

build finished in d.ddds
//// [/home/src/workspaces/project/index.js] *rewrite with same content*

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/index.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/index.ts


Diff:: incremental resolves to .js output from prior build while clean build cannot find module
--- nonIncremental /home/src/workspaces/project/helper2.js
+++ incremental /home/src/workspaces/project/helper2.js
@@ -1,1 +0,0 @@
-export const helper = 1;
--- nonIncremental.output.txt
+++ incremental.output.txt
@@ -1,4 +1,4 @@
-[96mindex.ts[0m:[93m1[0m:[93m24[0m - [91merror[0m[90m TS2307: [0mCannot find module './helper' or its corresponding type declarations.
+[96mindex.ts[0m:[93m1[0m:[93m24[0m - [91merror[0m[90m TS7016: [0mCould not find a declaration file for module './helper'. '/home/src/workspaces/project/helper.js' implicitly has an 'any' type.

 [7m1[0m import { helper } from "./helper";
 [7m [0m [91m                       ~~~~~~~~~~[0m

Edit [1]:: rename back to helper
//// [/home/src/workspaces/project/helper.ts] *new* 
export const helper = 1;
//// [/home/src/workspaces/project/helper2.ts] *deleted*


Output::
build starting at HH:MM:SS AM
build finished in d.ddds
//// [/home/src/workspaces/project/helper.js] *rewrite with same content*
//// [/home/src/workspaces/project/index.js] *rewrite with same content*

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/helper.ts
*refresh*    /home/src/workspaces/project/index.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/helper.ts
(computed .d.ts) /home/src/workspaces/project/index.ts

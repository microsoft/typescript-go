currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/src/a.ts] *new* 
export const a = 1;
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
	"compilerOptions": {},
	"include": ["src/**/*.ts"]
}

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
//// [/home/src/workspaces/project/src/a.js] *new* 
export const a = 1;


tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/src/a.ts
Signatures::


Edit [0]:: add new file to existing src directory
//// [/home/src/workspaces/project/src/b.ts] *new* 
export const b = 2;


Output::

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/src/a.ts
Signatures::


Diff:: incremental skips emit for new unreferenced file
--- nonIncremental /home/src/workspaces/project/src/b.js
+++ incremental /home/src/workspaces/project/src/b.js
@@ -1,1 +0,0 @@
-export const b = 2;

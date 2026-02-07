currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/a.ts] *new* 
export const a = 1;
//// [/home/src/workspaces/project/b.ts] *new* 
export const b = 2;
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{"compilerOptions": {"outDir": "dist"}}

tsgo --watch
ExitStatus:: Success
Output::
build starting at HH:MM:SS AM
build finished in d.ddds
//// [/home/src/tslibs/TS/Lib/lib.es2024.full.d.ts] *Lib*
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
//// [/home/src/workspaces/project/dist/a.js] *new* 
export const a = 1;

//// [/home/src/workspaces/project/dist/b.js] *new* 
export const b = 2;


tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2024.full.d.ts
*refresh*    /home/src/workspaces/project/a.ts
*refresh*    /home/src/workspaces/project/b.ts
Signatures::


Edit [0]:: delete b.ts
//// [/home/src/workspaces/project/b.ts] *deleted*


Output::
build starting at HH:MM:SS AM
build finished in d.ddds
//// [/home/src/workspaces/project/dist/b.js] *deleted*

tsconfig.json::
SemanticDiagnostics::
Signatures::

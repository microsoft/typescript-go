
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::-w
//// [/home/src/workspaces/project/a.ts] *new* 
const a = class { private p = 10; };
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
	"compilerOptions": {
            "noEmit": true
	}
}

ExitStatus:: 0

CompilerOptions::{
    "watch": true
}
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


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/a.ts

Signatures::


Edit:: fix syntax error
//// [/home/src/workspaces/project/a.ts] *modified* 
const a = "hello";


Output::


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/a.ts

Signatures::
(used version)   /home/src/tslibs/TS/Lib/lib.d.ts


Edit:: emit after fixing error
//// [/home/src/workspaces/project/tsconfig.json] *modified* 
{
	"compilerOptions": {
            
	}
}


Output::
//// [/home/src/workspaces/project/a.js] *new* 
const a = "hello";



SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/a.ts

Signatures::
(used version)   /home/src/tslibs/TS/Lib/lib.d.ts
(computed .d.ts) /home/src/workspaces/project/a.ts


Edit:: no emit run after fixing error
//// [/home/src/workspaces/project/tsconfig.json] *modified* 
{
	"compilerOptions": {
            "noEmit": true,
            
	}
}


Output::


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/a.ts

Signatures::
(used version)   /home/src/tslibs/TS/Lib/lib.d.ts
(computed .d.ts) /home/src/workspaces/project/a.ts


Edit:: introduce error
//// [/home/src/workspaces/project/a.ts] *modified* 
const a = class { private p = 10; };


Output::


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/a.ts

Signatures::
(used version)   /home/src/tslibs/TS/Lib/lib.d.ts


Edit:: emit when error
//// [/home/src/workspaces/project/tsconfig.json] *modified* 
{
	"compilerOptions": {
            
	}
}


Output::
//// [/home/src/workspaces/project/a.js] *modified* 
const a = class {
    p = 10;
};



SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/a.ts

Signatures::
(used version)   /home/src/tslibs/TS/Lib/lib.d.ts
(computed .d.ts) /home/src/workspaces/project/a.ts


Edit:: no emit run when error
//// [/home/src/workspaces/project/tsconfig.json] *modified* 
{
	"compilerOptions": {
            "noEmit": true,
            
	}
}


Output::


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/a.ts

Signatures::
(used version)   /home/src/tslibs/TS/Lib/lib.d.ts
(computed .d.ts) /home/src/workspaces/project/a.ts

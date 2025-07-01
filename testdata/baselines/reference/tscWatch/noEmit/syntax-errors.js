
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::-w
//// [/home/src/workspaces/project/a.ts] *new* 
const a = "hello
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
	"compilerOptions": {
            "noEmit": true
	}
}



CompilerOptions::{
    "watch": true
}


Output::
[96ma.ts[0m:[93m1[0m:[93m17[0m - [91merror[0m[90m TS1002: [0mUnterminated string literal.

[7m1[0m const a = "hello
[7m [0m [91m                ~[0m


Found 1 error in a.ts[90m:1[0m

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



Edit:: fix syntax error

Output::
//// [/home/src/workspaces/project/a.ts] *modified* 
&{const a = "hello"; 0xc00095d5c0}



Edit:: emit after fixing error

Output::
//// [/home/src/workspaces/project/a.js] *new* 
const a = "hello";

//// [/home/src/workspaces/project/tsconfig.json] *modified* 
&{{
	"compilerOptions": {
            
	}
} 0xc000cc3ec0}



Edit:: no emit run after fixing error

Output::
//// [/home/src/workspaces/project/tsconfig.json] *modified* 
&{{
	"compilerOptions": {
            "noEmit": true,
            
	}
} 0xc0011af2c0}



Edit:: introduce error

Output::
[96ma.ts[0m:[93m1[0m:[93m17[0m - [91merror[0m[90m TS1002: [0mUnterminated string literal.

[7m1[0m const a = "hello
[7m [0m [91m                ~[0m


Found 1 error in a.ts[90m:1[0m

//// [/home/src/workspaces/project/a.ts] *modified* 
&{const a = "hello 0xc001508990}



Edit:: emit when error

Output::
[96ma.ts[0m:[93m1[0m:[93m17[0m - [91merror[0m[90m TS1002: [0mUnterminated string literal.

[7m1[0m const a = "hello
[7m [0m [91m                ~[0m


Found 1 error in a.ts[90m:1[0m

//// [/home/src/workspaces/project/a.js] *modified* 
&{const a = "hello;
 0xc000fb26f0}
//// [/home/src/workspaces/project/tsconfig.json] *modified* 
&{{
	"compilerOptions": {
            
	}
} 0xc000fb2750}



Edit:: no emit run when error

Output::
[96ma.ts[0m:[93m1[0m:[93m17[0m - [91merror[0m[90m TS1002: [0mUnterminated string literal.

[7m1[0m const a = "hello
[7m [0m [91m                ~[0m


Found 1 error in a.ts[90m:1[0m

//// [/home/src/workspaces/project/tsconfig.json] *modified* 
&{{
	"compilerOptions": {
            "noEmit": true,
            
	}
} 0xc00062bce0}


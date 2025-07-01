
currentDirectory::/home/src/projects/myproject
useCaseSensitiveFileNames::true
Input::--explainFiles --outDir ${configDir}/outDir
//// [/home/src/projects/configs/first/tsconfig.json] *new* 
{
	"extends": "../second/tsconfig.json",
	"include": ["${configDir}/src"],
	"compilerOptions": {
		"typeRoots": ["root1", "${configDir}/root2", "root3"],
		"types": [],
	},
}
//// [/home/src/projects/configs/second/tsconfig.json] *new* 
{
	"files": ["${configDir}/main.ts"],
	"compilerOptions": {
		"declarationDir": "${configDir}/decls",
		"paths": {
			"@myscope/*": ["${configDir}/types/*"],
			"other/*": ["other/*"],
		},
		"baseUrl": "${configDir}",
	},
	"watchOptions": {
		"excludeFiles": ["${configDir}/main.ts"],
	},
}
//// [/home/src/projects/myproject/main.ts] *new* 

	// some comment
	export const y = 10;
	import { x } from "@myscope/sometype";

//// [/home/src/projects/myproject/root2/other/sometype2/index.d.ts] *new* 

	export const k = 10;

//// [/home/src/projects/myproject/src/secondary.ts] *new* 

	// some comment
	export const z = 10;
	import { k } from "other/sometype2";

//// [/home/src/projects/myproject/tsconfig.json] *new* 
{
	"extends": "../configs/first/tsconfig.json",
	"compilerOptions": {
		"declaration": true,
		"outDir": "outDir",
		"traceResolution": true,
	},
}
//// [/home/src/projects/myproject/types/sometype.ts] *new* 

	export const x = 10;


ExitStatus:: 2

CompilerOptions::{
    "outDir": "/home/src/projects/myproject/${configDir}/outDir",
    "explainFiles": true
}
Output::
[96msrc/secondary.ts[0m:[93m4[0m:[93m20[0m - [91merror[0m[90m TS2307: [0mCannot find module 'other/sometype2' or its corresponding type declarations.

[7m4[0m  import { k } from "other/sometype2";
[7m [0m [91m                   ~~~~~~~~~~~~~~~~~[0m


Found 1 error in src/secondary.ts[90m:4[0m

//// [/home/src/projects/myproject/${configDir}/outDir/main.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.y = void 0;
// some comment
exports.y = 10;

//// [/home/src/projects/myproject/${configDir}/outDir/src/secondary.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.z = void 0;
// some comment
exports.z = 10;

//// [/home/src/projects/myproject/${configDir}/outDir/types/sometype.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
exports.x = 10;

//// [/home/src/projects/myproject/decls/main.d.ts] *new* 
// some comment
export declare const y = 10;

//// [/home/src/projects/myproject/decls/src/secondary.d.ts] *new* 
// some comment
export declare const z = 10;

//// [/home/src/projects/myproject/decls/types/sometype.d.ts] *new* 
export declare const x = 10;

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


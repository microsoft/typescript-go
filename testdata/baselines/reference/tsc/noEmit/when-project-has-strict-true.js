
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::--noEmit
//// [/home/src/workspaces/project/class1.ts] *new* 
export class class1 {}
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
	"compilerOptions": {
		"incremental": true,
		"strict": true,
	},
}

ExitStatus:: 0

CompilerOptions::{
    "noEmit": true
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
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./class1.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"a7765a20d4489ae259632d5fe609919af401c278b7a90516894ef2774ce3bc97"],"options":{"strict":true},"affectedFilesPendingEmit":[2]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./class1.ts"
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
      "fileName": "./class1.ts",
      "version": "a7765a20d4489ae259632d5fe609919af401c278b7a90516894ef2774ce3bc97",
      "signature": "a7765a20d4489ae259632d5fe609919af401c278b7a90516894ef2774ce3bc97",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "options": {
    "strict": true
  },
  "affectedFilesPendingEmit": [
    [
      "./class1.ts",
      "Js",
      2
    ]
  ],
  "size": 351
}


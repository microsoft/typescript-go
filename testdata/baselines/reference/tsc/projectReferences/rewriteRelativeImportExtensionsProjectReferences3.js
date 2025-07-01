
currentDirectory::/home/src/workspaces/solution
useCaseSensitiveFileNames::true
Input::--p src/services --pretty false
//// [/home/src/workspaces/solution/dist/compiler/parser.d.ts] *new* 
export {};
//// [/home/src/workspaces/solution/src/compiler/parser.ts] *new* 
export {};
//// [/home/src/workspaces/solution/src/compiler/tsconfig.json] *new* 
{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {
						"rootDir": ".",
						"outDir": "../../dist/compiler"
					}
				}
//// [/home/src/workspaces/solution/src/services/services.ts] *new* 
import {} from "../compiler/parser.ts";
//// [/home/src/workspaces/solution/src/services/tsconfig.json] *new* 
{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {
						"rootDir": ".", 
						"outDir": "../../dist/services"
					},
					"references": [
						{ "path": "../compiler" }
					]
				}
//// [/home/src/workspaces/solution/src/tsconfig-base.json] *new* 
{
					"compilerOptions": { 
						"module": "nodenext",
						"composite": true,
						"rewriteRelativeImportExtensions": true
					}
				}

ExitStatus:: 0

CompilerOptions::{
    "project": "/home/src/workspaces/solution/src/services",
    "pretty": false
}
Output::
No output
//// [/home/src/tslibs/TS/Lib/lib.esnext.full.d.ts] *Lib*
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
//// [/home/src/workspaces/solution/dist/services/services.d.ts] *new* 
export {};

//// [/home/src/workspaces/solution/dist/services/services.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });

//// [/home/src/workspaces/solution/dist/services/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../../../tslibs/TS/Lib/lib.esnext.full.d.ts","../compiler/parser.d.ts","../../src/services/services.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"2e29cd9a98755c46896f7a2d56524db2d6d96b248e36db46de14c30bf47c8d05",{"version":"407537635fda1a543a422ecdd456c1402aaa2083cde5acfb4eb424ab02fc0612","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"composite":true,"module":199,"outDir":"./","rewriteRelativeImportExtensions":true,"rootDir":"../../src/services"},"referencedMap":[[3,1]],"latestChangedDtsFile":"./services.d.ts"}
//// [/home/src/workspaces/solution/dist/services/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../../../tslibs/TS/Lib/lib.esnext.full.d.ts",
    "../compiler/parser.d.ts",
    "../../src/services/services.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../../../tslibs/TS/Lib/lib.esnext.full.d.ts",
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
      "fileName": "../compiler/parser.d.ts",
      "version": "2e29cd9a98755c46896f7a2d56524db2d6d96b248e36db46de14c30bf47c8d05",
      "signature": "2e29cd9a98755c46896f7a2d56524db2d6d96b248e36db46de14c30bf47c8d05",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../../src/services/services.ts",
      "version": "407537635fda1a543a422ecdd456c1402aaa2083cde5acfb4eb424ab02fc0612",
      "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "407537635fda1a543a422ecdd456c1402aaa2083cde5acfb4eb424ab02fc0612",
        "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../compiler/parser.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "module": 199,
    "outDir": "./",
    "rewriteRelativeImportExtensions": true,
    "rootDir": "../../src/services"
  },
  "referencedMap": {
    "../../src/services/services.ts": [
      "../compiler/parser.d.ts"
    ]
  },
  "latestChangedDtsFile": "./services.d.ts",
  "size": 748
}


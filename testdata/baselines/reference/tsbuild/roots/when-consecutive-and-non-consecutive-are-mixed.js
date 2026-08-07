currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/anotherNonConsecutive.ts] *new* 
import { random } from "./random2";
export const nonConsecutive = "hello";
//// [/home/src/workspaces/project/asArray1.ts] *new* 
import { random } from "./random1";
export const x = "hello";
//// [/home/src/workspaces/project/asArray2.ts] *new* 
export const x = "hello";
//// [/home/src/workspaces/project/asArray3.ts] *new* 
export const x = "hello";
//// [/home/src/workspaces/project/file1.ts] *new* 
export const x = "hello";
//// [/home/src/workspaces/project/file2.ts] *new* 
export const y = "world";
//// [/home/src/workspaces/project/nonconsecutive.ts] *new* 
import { random } from "./random";
    export const nonConsecutive = "hello";
//// [/home/src/workspaces/project/random.d.ts] *new* 
export const random = "hello";
//// [/home/src/workspaces/project/random1.d.ts] *new* 
export const random = "hello";
//// [/home/src/workspaces/project/random2.d.ts] *new* 
export const random = "hello";
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "include": ["file*.ts", "nonconsecutive*.ts", "asArray*.ts", "anotherNonConsecutive.ts"],
}

tsgo --b -v
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'tsconfig.json' is out of date because output file 'tsconfig.tsbuildinfo' does not exist

[[90mHH:MM:SS AM[0m] Building project 'tsconfig.json'...

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
//// [/home/src/workspaces/project/anotherNonConsecutive.d.ts] *new* 
export declare const nonConsecutive = "hello";

//// [/home/src/workspaces/project/anotherNonConsecutive.js] *new* 
export const nonConsecutive = "hello";

//// [/home/src/workspaces/project/asArray1.d.ts] *new* 
export declare const x = "hello";

//// [/home/src/workspaces/project/asArray1.js] *new* 
export const x = "hello";

//// [/home/src/workspaces/project/asArray2.d.ts] *new* 
export declare const x = "hello";

//// [/home/src/workspaces/project/asArray2.js] *new* 
export const x = "hello";

//// [/home/src/workspaces/project/asArray3.d.ts] *new* 
export declare const x = "hello";

//// [/home/src/workspaces/project/asArray3.js] *new* 
export const x = "hello";

//// [/home/src/workspaces/project/file1.d.ts] *new* 
export declare const x = "hello";

//// [/home/src/workspaces/project/file1.js] *new* 
export const x = "hello";

//// [/home/src/workspaces/project/file2.d.ts] *new* 
export declare const y = "world";

//// [/home/src/workspaces/project/file2.js] *new* 
export const y = "world";

//// [/home/src/workspaces/project/nonconsecutive.d.ts] *new* 
export declare const nonConsecutive = "hello";

//// [/home/src/workspaces/project/nonconsecutive.js] *new* 
export const nonConsecutive = "hello";

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,3],5,[7,9],11],"fileNames":["lib.es2025.full.d.ts","./file1.ts","./file2.ts","./random.d.ts","./nonconsecutive.ts","./random1.d.ts","./asArray1.ts","./asArray2.ts","./asArray3.ts","./random2.d.ts","./anotherNonConsecutive.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";","42f7437ec4aafe7a26fc38ec3ba035fe-export const y = \"world\";","cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";","805bd15692edf65c8f946bb322e0a876-import { random } from \"./random\";\n    export const nonConsecutive = \"hello\";","cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";","f1ff3ce44c97592c5ef5a510c8ed269d-import { random } from \"./random1\";\nexport const x = \"hello\";","cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";","cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";","cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";","f5ecedece6d42fea3feb8b5832aacd89-import { random } from \"./random2\";\nexport const nonConsecutive = \"hello\";"],"fileIdsList":[[10],[6],[4]],"options":{"composite":true},"referencedMap":[[11,1],[7,2],[5,3]],"latestChangedDtsFile":"./anotherNonConsecutive.d.ts","emitSignatures":[[2,"0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"],[3,"64ca81919be0c5adb4964999189ddb2c-export declare const y = \"world\";\n"],[5,"ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"],[7,"0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"],[8,"0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"],[9,"0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"],[11,"ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./file1.ts",
        "./file2.ts"
      ],
      "original": [
        2,
        3
      ]
    },
    {
      "files": [
        "./nonconsecutive.ts"
      ],
      "original": 5
    },
    {
      "files": [
        "./asArray1.ts",
        "./asArray2.ts",
        "./asArray3.ts"
      ],
      "original": [
        7,
        9
      ]
    },
    {
      "files": [
        "./anotherNonConsecutive.ts"
      ],
      "original": 11
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./file1.ts",
    "./file2.ts",
    "./random.d.ts",
    "./nonconsecutive.ts",
    "./random1.d.ts",
    "./asArray1.ts",
    "./asArray2.ts",
    "./asArray3.ts",
    "./random2.d.ts",
    "./anotherNonConsecutive.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2025.full.d.ts",
      "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./file1.ts",
      "version": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "signature": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./file2.ts",
      "version": "42f7437ec4aafe7a26fc38ec3ba035fe-export const y = \"world\";",
      "signature": "42f7437ec4aafe7a26fc38ec3ba035fe-export const y = \"world\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./random.d.ts",
      "version": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "signature": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./nonconsecutive.ts",
      "version": "805bd15692edf65c8f946bb322e0a876-import { random } from \"./random\";\n    export const nonConsecutive = \"hello\";",
      "signature": "805bd15692edf65c8f946bb322e0a876-import { random } from \"./random\";\n    export const nonConsecutive = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./random1.d.ts",
      "version": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "signature": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./asArray1.ts",
      "version": "f1ff3ce44c97592c5ef5a510c8ed269d-import { random } from \"./random1\";\nexport const x = \"hello\";",
      "signature": "f1ff3ce44c97592c5ef5a510c8ed269d-import { random } from \"./random1\";\nexport const x = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./asArray2.ts",
      "version": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "signature": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./asArray3.ts",
      "version": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "signature": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./random2.d.ts",
      "version": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "signature": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./anotherNonConsecutive.ts",
      "version": "f5ecedece6d42fea3feb8b5832aacd89-import { random } from \"./random2\";\nexport const nonConsecutive = \"hello\";",
      "signature": "f5ecedece6d42fea3feb8b5832aacd89-import { random } from \"./random2\";\nexport const nonConsecutive = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./random2.d.ts"
    ],
    [
      "./random1.d.ts"
    ],
    [
      "./random.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./anotherNonConsecutive.ts": [
      "./random2.d.ts"
    ],
    "./asArray1.ts": [
      "./random1.d.ts"
    ],
    "./nonconsecutive.ts": [
      "./random.d.ts"
    ]
  },
  "latestChangedDtsFile": "./anotherNonConsecutive.d.ts",
  "emitSignatures": [
    {
      "file": "./file1.ts",
      "signature": "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n",
      "original": [
        2,
        "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"
      ]
    },
    {
      "file": "./file2.ts",
      "signature": "64ca81919be0c5adb4964999189ddb2c-export declare const y = \"world\";\n",
      "original": [
        3,
        "64ca81919be0c5adb4964999189ddb2c-export declare const y = \"world\";\n"
      ]
    },
    {
      "file": "./nonconsecutive.ts",
      "signature": "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n",
      "original": [
        5,
        "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"
      ]
    },
    {
      "file": "./asArray1.ts",
      "signature": "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n",
      "original": [
        7,
        "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"
      ]
    },
    {
      "file": "./asArray2.ts",
      "signature": "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n",
      "original": [
        8,
        "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"
      ]
    },
    {
      "file": "./asArray3.ts",
      "signature": "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n",
      "original": [
        9,
        "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"
      ]
    },
    {
      "file": "./anotherNonConsecutive.ts",
      "signature": "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n",
      "original": [
        11,
        "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"
      ]
    }
  ],
  "size": 2574
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/file1.ts
*refresh*    /home/src/workspaces/project/file2.ts
*refresh*    /home/src/workspaces/project/random.d.ts
*refresh*    /home/src/workspaces/project/nonconsecutive.ts
*refresh*    /home/src/workspaces/project/random1.d.ts
*refresh*    /home/src/workspaces/project/asArray1.ts
*refresh*    /home/src/workspaces/project/asArray2.ts
*refresh*    /home/src/workspaces/project/asArray3.ts
*refresh*    /home/src/workspaces/project/random2.d.ts
*refresh*    /home/src/workspaces/project/anotherNonConsecutive.ts
Signatures::
(stored at emit) /home/src/workspaces/project/file1.ts
(stored at emit) /home/src/workspaces/project/file2.ts
(stored at emit) /home/src/workspaces/project/nonconsecutive.ts
(stored at emit) /home/src/workspaces/project/asArray1.ts
(stored at emit) /home/src/workspaces/project/asArray2.ts
(stored at emit) /home/src/workspaces/project/asArray3.ts
(stored at emit) /home/src/workspaces/project/anotherNonConsecutive.ts


Edit [0]:: delete file1
//// [/home/src/workspaces/project/file1.d.ts] *deleted*
//// [/home/src/workspaces/project/file1.js] *deleted*
//// [/home/src/workspaces/project/file1.ts] *deleted*

tsgo --b -v
ExitStatus:: Success
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'tsconfig.json' is out of date because buildinfo file 'tsconfig.tsbuildinfo' indicates that file 'file1.ts' was root file of compilation but not any more.

[[90mHH:MM:SS AM[0m] Building project 'tsconfig.json'...

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2,4,[6,8],10],"fileNames":["lib.es2025.full.d.ts","./file2.ts","./random.d.ts","./nonconsecutive.ts","./random1.d.ts","./asArray1.ts","./asArray2.ts","./asArray3.ts","./random2.d.ts","./anotherNonConsecutive.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"42f7437ec4aafe7a26fc38ec3ba035fe-export const y = \"world\";","cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";","805bd15692edf65c8f946bb322e0a876-import { random } from \"./random\";\n    export const nonConsecutive = \"hello\";","cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";","f1ff3ce44c97592c5ef5a510c8ed269d-import { random } from \"./random1\";\nexport const x = \"hello\";","cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";","cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";","cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";","f5ecedece6d42fea3feb8b5832aacd89-import { random } from \"./random2\";\nexport const nonConsecutive = \"hello\";"],"fileIdsList":[[9],[5],[3]],"options":{"composite":true},"referencedMap":[[10,1],[6,2],[4,3]],"latestChangedDtsFile":"./anotherNonConsecutive.d.ts","emitSignatures":[[2,"64ca81919be0c5adb4964999189ddb2c-export declare const y = \"world\";\n"],[4,"ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"],[6,"0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"],[7,"0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"],[8,"0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"],[10,"ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./file2.ts"
      ],
      "original": 2
    },
    {
      "files": [
        "./nonconsecutive.ts"
      ],
      "original": 4
    },
    {
      "files": [
        "./asArray1.ts",
        "./asArray2.ts",
        "./asArray3.ts"
      ],
      "original": [
        6,
        8
      ]
    },
    {
      "files": [
        "./anotherNonConsecutive.ts"
      ],
      "original": 10
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./file2.ts",
    "./random.d.ts",
    "./nonconsecutive.ts",
    "./random1.d.ts",
    "./asArray1.ts",
    "./asArray2.ts",
    "./asArray3.ts",
    "./random2.d.ts",
    "./anotherNonConsecutive.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2025.full.d.ts",
      "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./file2.ts",
      "version": "42f7437ec4aafe7a26fc38ec3ba035fe-export const y = \"world\";",
      "signature": "42f7437ec4aafe7a26fc38ec3ba035fe-export const y = \"world\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./random.d.ts",
      "version": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "signature": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./nonconsecutive.ts",
      "version": "805bd15692edf65c8f946bb322e0a876-import { random } from \"./random\";\n    export const nonConsecutive = \"hello\";",
      "signature": "805bd15692edf65c8f946bb322e0a876-import { random } from \"./random\";\n    export const nonConsecutive = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./random1.d.ts",
      "version": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "signature": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./asArray1.ts",
      "version": "f1ff3ce44c97592c5ef5a510c8ed269d-import { random } from \"./random1\";\nexport const x = \"hello\";",
      "signature": "f1ff3ce44c97592c5ef5a510c8ed269d-import { random } from \"./random1\";\nexport const x = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./asArray2.ts",
      "version": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "signature": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./asArray3.ts",
      "version": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "signature": "cc7052ed344567798ec87f1c0f8f276c-export const x = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./random2.d.ts",
      "version": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "signature": "cc6fd496b46fb211a7b922d3e31cdc4e-export const random = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./anotherNonConsecutive.ts",
      "version": "f5ecedece6d42fea3feb8b5832aacd89-import { random } from \"./random2\";\nexport const nonConsecutive = \"hello\";",
      "signature": "f5ecedece6d42fea3feb8b5832aacd89-import { random } from \"./random2\";\nexport const nonConsecutive = \"hello\";",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./random2.d.ts"
    ],
    [
      "./random1.d.ts"
    ],
    [
      "./random.d.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./anotherNonConsecutive.ts": [
      "./random2.d.ts"
    ],
    "./asArray1.ts": [
      "./random1.d.ts"
    ],
    "./nonconsecutive.ts": [
      "./random.d.ts"
    ]
  },
  "latestChangedDtsFile": "./anotherNonConsecutive.d.ts",
  "emitSignatures": [
    {
      "file": "./file2.ts",
      "signature": "64ca81919be0c5adb4964999189ddb2c-export declare const y = \"world\";\n",
      "original": [
        2,
        "64ca81919be0c5adb4964999189ddb2c-export declare const y = \"world\";\n"
      ]
    },
    {
      "file": "./nonconsecutive.ts",
      "signature": "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n",
      "original": [
        4,
        "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"
      ]
    },
    {
      "file": "./asArray1.ts",
      "signature": "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n",
      "original": [
        6,
        "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"
      ]
    },
    {
      "file": "./asArray2.ts",
      "signature": "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n",
      "original": [
        7,
        "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"
      ]
    },
    {
      "file": "./asArray3.ts",
      "signature": "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n",
      "original": [
        8,
        "0c71c4d05f424f4dc52c978a9207cdf6-export declare const x = \"hello\";\n"
      ]
    },
    {
      "file": "./anotherNonConsecutive.ts",
      "signature": "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n",
      "original": [
        10,
        "ed1a51017d86c01f126fd5c93a77072d-export declare const nonConsecutive = \"hello\";\n"
      ]
    }
  ],
  "size": 2416
}

tsconfig.json::
SemanticDiagnostics::
Signatures::

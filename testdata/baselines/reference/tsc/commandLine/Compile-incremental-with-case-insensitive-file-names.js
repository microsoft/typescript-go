currentDirectory::/home/project
useCaseSensitiveFileNames::false
Input::
//// [/home/node_modules/lib1/index.d.ts] *new* 
import type { Foo } from 'someLib';
export type { Foo as Foo1 };
//// [/home/node_modules/lib1/package.json] *new* 
{
    "name": "lib1"
}
//// [/home/node_modules/lib2/index.d.ts] *new* 
import type { Foo } from 'somelib';
export type { Foo as Foo2 };
export declare const foo2: Foo;
//// [/home/node_modules/lib2/package.json] *new* 
{
    "name": "lib2"
}
//// [/home/node_modules/otherLib/index.d.ts] *new* 
export type Str = string;
//// [/home/node_modules/otherLib/package.json] *new* 
{
    "name": "otherlib"
}
//// [/home/node_modules/someLib/index.d.ts] *new* 
import type { Str } from 'otherLib';
export type Foo = { foo: Str; };
//// [/home/node_modules/someLib/package.json] *new* 
        {
"name": "somelib"
        }
//// [/home/project/src/index.ts] *new* 
import type { Foo1 } from 'lib1';
import type { Foo2 } from 'lib2';
export const foo1: Foo1 = { foo: "a" };
export const foo2: Foo2 = { foo: "b" };
//// [/home/project/tsconfig.json] *new* 
        {
"compilerOptions": {
				"incremental": true
},
        }

tsgo -p .
ExitStatus:: Success
Output::
//// [/home/project/src/index.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.foo2 = exports.foo1 = void 0;
exports.foo1 = { foo: "a" };
exports.foo2 = { foo: "b" };

//// [/home/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../src/tslibs/ts/lib/lib.d.ts","../node_modules/otherlib/index.d.ts","../node_modules/somelib/index.d.ts","../node_modules/lib1/index.d.ts","../node_modules/lib2/index.d.ts","./src/index.ts"],"fileInfos":[{"version":"eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"a40368814f37d4357a34c3f941f3e677-export type Str = string;","c83d508c75590e2f7051911d4a827ceb-import type { Str } from 'otherLib';\nexport type Foo = { foo: Str; };","a508f6d5ab12f0c2168b5488ebbd34be-import type { Foo } from 'someLib';\nexport type { Foo as Foo1 };","528165de17b5f9f1dbbf984a062d9b57-import type { Foo } from 'somelib';\nexport type { Foo as Foo2 };\nexport declare const foo2: Foo;","c367835e98822462afe5cacbac4d94c2-import type { Foo1 } from 'lib1';\nimport type { Foo2 } from 'lib2';\nexport const foo1: Foo1 = { foo: \"a\" };\nexport const foo2: Foo2 = { foo: \"b\" };"],"fileIdsList":[[3],[2],[4,5]],"referencedMap":[[4,1],[5,1],[3,2],[6,3]]}
//// [/home/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../src/tslibs/ts/lib/lib.d.ts",
    "../node_modules/otherlib/index.d.ts",
    "../node_modules/somelib/index.d.ts",
    "../node_modules/lib1/index.d.ts",
    "../node_modules/lib2/index.d.ts",
    "./src/index.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../src/tslibs/ts/lib/lib.d.ts",
      "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eae9e83ef0f77eeb2e35dc9b91facce1-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../node_modules/otherlib/index.d.ts",
      "version": "a40368814f37d4357a34c3f941f3e677-export type Str = string;",
      "signature": "a40368814f37d4357a34c3f941f3e677-export type Str = string;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../node_modules/somelib/index.d.ts",
      "version": "c83d508c75590e2f7051911d4a827ceb-import type { Str } from 'otherLib';\nexport type Foo = { foo: Str; };",
      "signature": "c83d508c75590e2f7051911d4a827ceb-import type { Str } from 'otherLib';\nexport type Foo = { foo: Str; };",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../node_modules/lib1/index.d.ts",
      "version": "a508f6d5ab12f0c2168b5488ebbd34be-import type { Foo } from 'someLib';\nexport type { Foo as Foo1 };",
      "signature": "a508f6d5ab12f0c2168b5488ebbd34be-import type { Foo } from 'someLib';\nexport type { Foo as Foo1 };",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../node_modules/lib2/index.d.ts",
      "version": "528165de17b5f9f1dbbf984a062d9b57-import type { Foo } from 'somelib';\nexport type { Foo as Foo2 };\nexport declare const foo2: Foo;",
      "signature": "528165de17b5f9f1dbbf984a062d9b57-import type { Foo } from 'somelib';\nexport type { Foo as Foo2 };\nexport declare const foo2: Foo;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/index.ts",
      "version": "c367835e98822462afe5cacbac4d94c2-import type { Foo1 } from 'lib1';\nimport type { Foo2 } from 'lib2';\nexport const foo1: Foo1 = { foo: \"a\" };\nexport const foo2: Foo2 = { foo: \"b\" };",
      "signature": "c367835e98822462afe5cacbac4d94c2-import type { Foo1 } from 'lib1';\nimport type { Foo2 } from 'lib2';\nexport const foo1: Foo1 = { foo: \"a\" };\nexport const foo2: Foo2 = { foo: \"b\" };",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "../node_modules/somelib/index.d.ts"
    ],
    [
      "../node_modules/otherlib/index.d.ts"
    ],
    [
      "../node_modules/lib1/index.d.ts",
      "../node_modules/lib2/index.d.ts"
    ]
  ],
  "referencedMap": {
    "../node_modules/lib1/index.d.ts": [
      "../node_modules/somelib/index.d.ts"
    ],
    "../node_modules/lib2/index.d.ts": [
      "../node_modules/somelib/index.d.ts"
    ],
    "../node_modules/somelib/index.d.ts": [
      "../node_modules/otherlib/index.d.ts"
    ],
    "./src/index.ts": [
      "../node_modules/lib1/index.d.ts",
      "../node_modules/lib2/index.d.ts"
    ]
  },
  "size": 1681
}
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
*refresh*    /home/node_modules/otherLib/index.d.ts
*refresh*    /home/node_modules/somelib/index.d.ts
*refresh*    /home/node_modules/lib1/index.d.ts
*refresh*    /home/node_modules/lib2/index.d.ts
*refresh*    /home/project/src/index.ts
Signatures::

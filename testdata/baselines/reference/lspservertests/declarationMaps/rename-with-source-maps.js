UseCaseSensitiveFileNames: false
//// [/user/username/projects/myproject/decls/FnS.d.ts] *new* 
export declare function fn1(): void;
export declare function fn2(): void;
export declare function fn3(): void;
export declare function fn4(): void;
export declare function fn5(): void;
//# sourceMappingURL=FnS.d.ts.map
//// [/user/username/projects/myproject/decls/FnS.d.ts.map] *new* 
{"version":3,"file":"FnS.d.ts","sourceRoot":"","sources":["../dependency/FnS.ts"],"names":[],"mappings":"AAAA,wBAAgB,GAAG,SAAM;AACzB,wBAAgB,GAAG,SAAM;AACzB,wBAAgB,GAAG,SAAM;AACzB,wBAAgB,GAAG,SAAM;AACzB,wBAAgB,GAAG,SAAM"}
//// [/user/username/projects/myproject/dependency/FnS.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.fn1 = fn1;
exports.fn2 = fn2;
exports.fn3 = fn3;
exports.fn4 = fn4;
exports.fn5 = fn5;
function fn1() { }
function fn2() { }
function fn3() { }
function fn4() { }
function fn5() { }

//// [/user/username/projects/myproject/dependency/FnS.ts] *new* 
export function fn1() { }
export function fn2() { }
export function fn3() { }
export function fn4() { }
export function fn5() { }

//// [/user/username/projects/myproject/dependency/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
        "declarationDir": "../decls"
    }
}
//// [/user/username/projects/myproject/dependency/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.d.ts","./FnS.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"dabc23167a661a3e71744cf63ff5adcc-export function fn1() { }\nexport function fn2() { }\nexport function fn3() { }\nexport function fn4() { }\nexport function fn5() { }\n","signature":"bf56b5172c2fd500384b437b7700d594-export declare function fn1(): void;\nexport declare function fn2(): void;\nexport declare function fn3(): void;\nexport declare function fn4(): void;\nexport declare function fn5(): void;\n","impliedNodeFormat":1}],"options":{"composite":true,"declarationDir":"../decls","declarationMap":true},"latestChangedDtsFile":"../decls/FnS.d.ts"}
//// [/user/username/projects/myproject/dependency/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./FnS.ts"
      ],
      "original": 2
    }
  ],
  "fileNames": [
    "lib.d.ts",
    "./FnS.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.d.ts",
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
      "fileName": "./FnS.ts",
      "version": "dabc23167a661a3e71744cf63ff5adcc-export function fn1() { }\nexport function fn2() { }\nexport function fn3() { }\nexport function fn4() { }\nexport function fn5() { }\n",
      "signature": "bf56b5172c2fd500384b437b7700d594-export declare function fn1(): void;\nexport declare function fn2(): void;\nexport declare function fn3(): void;\nexport declare function fn4(): void;\nexport declare function fn5(): void;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "dabc23167a661a3e71744cf63ff5adcc-export function fn1() { }\nexport function fn2() { }\nexport function fn3() { }\nexport function fn4() { }\nexport function fn5() { }\n",
        "signature": "bf56b5172c2fd500384b437b7700d594-export declare function fn1(): void;\nexport declare function fn2(): void;\nexport declare function fn3(): void;\nexport declare function fn4(): void;\nexport declare function fn5(): void;\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "composite": true,
    "declarationDir": "../decls",
    "declarationMap": true
  },
  "latestChangedDtsFile": "../decls/FnS.d.ts",
  "size": 1423
}
//// [/user/username/projects/myproject/main/main.d.ts] *new* 
export {};
//# sourceMappingURL=main.d.ts.map
//// [/user/username/projects/myproject/main/main.d.ts.map] *new* 
{"version":3,"file":"main.d.ts","sourceRoot":"","sources":["main.ts"],"names":[],"mappings":""}
//// [/user/username/projects/myproject/main/main.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const FnS_1 = require("../decls/FnS");
(0, FnS_1.fn1)();
(0, FnS_1.fn2)();
(0, FnS_1.fn3)();
(0, FnS_1.fn4)();
(0, FnS_1.fn5)();

//// [/user/username/projects/myproject/main/main.ts] *new* 
import {
    fn1,
    fn2,
    fn3,
    fn4,
    fn5
} from "../decls/FnS";

fn1();
fn2();
fn3();
fn4();
fn5();
//// [/user/username/projects/myproject/main/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
    },
}
//// [/user/username/projects/myproject/main/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[3],"fileNames":["lib.d.ts","../decls/FnS.d.ts","./main.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"ccbe7f3de9ff1ded77b994bef4cf9dce-export declare function fn1(): void;\nexport declare function fn2(): void;\nexport declare function fn3(): void;\nexport declare function fn4(): void;\nexport declare function fn5(): void;\n//# sourceMappingURL=FnS.d.ts.map",{"version":"4fd5832b607095e96718cc3f14637677-import {\n    fn1,\n    fn2,\n    fn3,\n    fn4,\n    fn5\n} from \"../decls/FnS\";\n\nfn1();\nfn2();\nfn3();\nfn4();\nfn5();","signature":"abe7d9981d6018efb6b2b794f40a1607-export {};\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"composite":true,"declarationMap":true},"referencedMap":[[3,1]],"latestChangedDtsFile":"./main.d.ts"}
//// [/user/username/projects/myproject/main/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./main.ts"
      ],
      "original": 3
    }
  ],
  "fileNames": [
    "lib.d.ts",
    "../decls/FnS.d.ts",
    "./main.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.d.ts",
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
      "fileName": "../decls/FnS.d.ts",
      "version": "ccbe7f3de9ff1ded77b994bef4cf9dce-export declare function fn1(): void;\nexport declare function fn2(): void;\nexport declare function fn3(): void;\nexport declare function fn4(): void;\nexport declare function fn5(): void;\n//# sourceMappingURL=FnS.d.ts.map",
      "signature": "ccbe7f3de9ff1ded77b994bef4cf9dce-export declare function fn1(): void;\nexport declare function fn2(): void;\nexport declare function fn3(): void;\nexport declare function fn4(): void;\nexport declare function fn5(): void;\n//# sourceMappingURL=FnS.d.ts.map",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./main.ts",
      "version": "4fd5832b607095e96718cc3f14637677-import {\n    fn1,\n    fn2,\n    fn3,\n    fn4,\n    fn5\n} from \"../decls/FnS\";\n\nfn1();\nfn2();\nfn3();\nfn4();\nfn5();",
      "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4fd5832b607095e96718cc3f14637677-import {\n    fn1,\n    fn2,\n    fn3,\n    fn4,\n    fn5\n} from \"../decls/FnS\";\n\nfn1();\nfn2();\nfn3();\nfn4();\nfn5();",
        "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../decls/FnS.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declarationMap": true
  },
  "referencedMap": {
    "./main.ts": [
      "../decls/FnS.d.ts"
    ]
  },
  "latestChangedDtsFile": "./main.d.ts",
  "size": 1525
}
//// [/user/username/projects/myproject/tsconfig.json] *new* 
{
    "references": [
        { "path": "main" }
    ]
}
//// [/user/username/projects/myproject/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":["./dependency/FnS.ts","./main/main.ts"]}
//// [/user/username/projects/myproject/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./dependency/FnS.ts"
      ],
      "original": "./dependency/FnS.ts"
    },
    {
      "files": [
        "./main/main.ts"
      ],
      "original": "./main/main.ts"
    }
  ],
  "size": 75
}
//// [/user/username/projects/random/random.ts] *new* 
export const a = 10;
//// [/user/username/projects/random/tsconfig.json] *new* 
{}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/dependency/FnS.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export function fn1() { }\nexport function fn2() { }\nexport function fn3() { }\nexport function fn4() { }\nexport function fn5() { }\n"
    }
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] *new*
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/tsconfig.json] *new*
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] *new*
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/dependency/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/dependency/fns.ts  
Config File Names::
  [/user/username/projects/myproject/dependency/fns.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/dependency/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/dependency/tsconfig.json  /user/username/projects/myproject/tsconfig.json 
      /user/username/projects/myproject/tsconfig.json              
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const a = 10;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] 
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/tsconfig.json] 
  [/user/username/projects/random/tsconfig.json] *new*
    /user/username/projects/random/random.ts  
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] 
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
  [/user/username/projects/random/random.ts] *new*
    /user/username/projects/random/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/myproject/dependency/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/dependency/fns.ts  
  [/user/username/projects/random/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/random/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/random/random.ts  
Config File Names::
  [/user/username/projects/myproject/dependency/fns.ts] 
    NearestConfigFileName: /user/username/projects/myproject/dependency/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/dependency/tsconfig.json  /user/username/projects/myproject/tsconfig.json 
      /user/username/projects/myproject/tsconfig.json              
  [/user/username/projects/random/random.ts] *new*
    NearestConfigFileName: /user/username/projects/random/tsconfig.json
{
  "method": "textDocument/rename",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/dependency/FnS.ts"
    },
    "position": {
      "line": 2,
      "character": 16
    },
    "newName": "?"
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] 
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/tsconfig.json] *modified*
    /user/username/projects/myproject/decls/FnS.d.ts     *new*
    /user/username/projects/myproject/dependency/FnS.ts  *new*
    /user/username/projects/myproject/main/main.ts       *new*
  [/user/username/projects/random/tsconfig.json] 
    /user/username/projects/random/random.ts  
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] *modified*
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
    /user/username/projects/myproject/tsconfig.json             *new*
  [/user/username/projects/random/random.ts] 
    /user/username/projects/random/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/myproject/dependency/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/dependency/fns.ts  
  [/user/username/projects/myproject/main/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/tsconfig.json  
  [/user/username/projects/myproject/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/tsconfig.json  
  [/user/username/projects/random/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/random/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/random/random.ts  
// === /user/username/projects/myproject/dependency/FnS.ts ===
// export function fn1() { }
// export function fn2() { }
// export function /*RENAME*/[|fn3RENAME|]() { }
// export function fn4() { }
// export function fn5() { }
// 
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] 
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
    /user/username/projects/myproject/tsconfig.json             
  [/user/username/projects/random/random.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const a = 10;"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] 
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
    /user/username/projects/myproject/tsconfig.json             
  [/user/username/projects/random/random.ts] *new*
    /user/username/projects/random/tsconfig.json  (default) 
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/dependency/FnS.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] *closed*
  [/user/username/projects/random/random.ts] 
    /user/username/projects/random/tsconfig.json  (default) 
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts"
    }
  }
}
Open Files::
  [/user/username/projects/random/random.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const a = 10;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] *deleted*
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/tsconfig.json] *deleted*
    /user/username/projects/myproject/decls/FnS.d.ts     
    /user/username/projects/myproject/dependency/FnS.ts  
    /user/username/projects/myproject/main/main.ts       
  [/user/username/projects/random/tsconfig.json] 
    /user/username/projects/random/random.ts  
Open Files::
  [/user/username/projects/random/random.ts] *new*
    /user/username/projects/random/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] *deleted*
  [/user/username/projects/myproject/main/tsconfig.json] *deleted*
  [/user/username/projects/myproject/tsconfig.json] *deleted*
  [/user/username/projects/random/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/random/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/random/random.ts  
Config File Names::
  [/user/username/projects/myproject/dependency/fns.ts] *deleted*
  [/user/username/projects/random/random.ts] 
    NearestConfigFileName: /user/username/projects/random/tsconfig.json

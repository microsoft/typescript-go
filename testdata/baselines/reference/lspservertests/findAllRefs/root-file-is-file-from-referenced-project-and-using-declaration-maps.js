UseCaseSensitiveFileNames: false
//// [/user/username/projects/project/out/input/keyboard.d.ts] *new* 
export declare function evaluateKeyboardEvent(): void;
//# sourceMappingURL=keyboard.d.ts.map
//// [/user/username/projects/project/out/input/keyboard.d.ts.map] *new* 
{"version":3,"file":"keyboard.d.ts","sourceRoot":"","sources":["../../src/common/input/keyboard.ts"],"names":[],"mappings":"AACA,wBAAgB,qBAAqB,SAAM"}
//// [/user/username/projects/project/out/input/keyboard.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.evaluateKeyboardEvent = evaluateKeyboardEvent;
function bar() { return "just a random function so .d.ts location doesnt match"; }
function evaluateKeyboardEvent() { }

//// [/user/username/projects/project/out/input/keyboard.test.d.ts] *new* 
export {};
//# sourceMappingURL=keyboard.test.d.ts.map
//// [/user/username/projects/project/out/input/keyboard.test.d.ts.map] *new* 
{"version":3,"file":"keyboard.test.d.ts","sourceRoot":"","sources":["../../src/common/input/keyboard.test.ts"],"names":[],"mappings":""}
//// [/user/username/projects/project/out/input/keyboard.test.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const keyboard_1 = require("common/input/keyboard");
function testEvaluateKeyboardEvent() {
    return (0, keyboard_1.evaluateKeyboardEvent)();
}

//// [/user/username/projects/project/out/src.tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,4]],"fileNames":["lib.d.ts","./input/keyboard.d.ts","../src/terminal.ts","./input/keyboard.test.d.ts","../src/common/input/keyboard.ts","../src/common/input/keyboard.test.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"362bf1943c7d2a927f25978d1f24d61f-export declare function evaluateKeyboardEvent(): void;\n//# sourceMappingURL=keyboard.d.ts.map",{"version":"03ece765cef2ff28231304b4d7007649-import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction foo() {\n    return evaluateKeyboardEvent();\n}","signature":"abe7d9981d6018efb6b2b794f40a1607-export {};\n","impliedNodeFormat":1},"aaff3960e003a07ef50db935cb91a696-export {};\n//# sourceMappingURL=keyboard.test.d.ts.map"],"fileIdsList":[[2]],"options":{"composite":true,"declarationMap":true,"outDir":"./","tsBuildInfoFile":"./src.tsconfig.tsbuildinfo"},"referencedMap":[[3,1]],"latestChangedDtsFile":"./terminal.d.ts","resolvedRoot":[[2,5],[4,6]]}
//// [/user/username/projects/project/out/src.tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./input/keyboard.d.ts",
        "../src/terminal.ts",
        "./input/keyboard.test.d.ts"
      ],
      "original": [
        2,
        4
      ]
    }
  ],
  "fileNames": [
    "lib.d.ts",
    "./input/keyboard.d.ts",
    "../src/terminal.ts",
    "./input/keyboard.test.d.ts",
    "../src/common/input/keyboard.ts",
    "../src/common/input/keyboard.test.ts"
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
      "fileName": "./input/keyboard.d.ts",
      "version": "362bf1943c7d2a927f25978d1f24d61f-export declare function evaluateKeyboardEvent(): void;\n//# sourceMappingURL=keyboard.d.ts.map",
      "signature": "362bf1943c7d2a927f25978d1f24d61f-export declare function evaluateKeyboardEvent(): void;\n//# sourceMappingURL=keyboard.d.ts.map",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../src/terminal.ts",
      "version": "03ece765cef2ff28231304b4d7007649-import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction foo() {\n    return evaluateKeyboardEvent();\n}",
      "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "03ece765cef2ff28231304b4d7007649-import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction foo() {\n    return evaluateKeyboardEvent();\n}",
        "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./input/keyboard.test.d.ts",
      "version": "aaff3960e003a07ef50db935cb91a696-export {};\n//# sourceMappingURL=keyboard.test.d.ts.map",
      "signature": "aaff3960e003a07ef50db935cb91a696-export {};\n//# sourceMappingURL=keyboard.test.d.ts.map",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./input/keyboard.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declarationMap": true,
    "outDir": "./",
    "tsBuildInfoFile": "./src.tsconfig.tsbuildinfo"
  },
  "referencedMap": {
    "../src/terminal.ts": [
      "./input/keyboard.d.ts"
    ]
  },
  "latestChangedDtsFile": "./terminal.d.ts",
  "resolvedRoot": [
    [
      "./input/keyboard.d.ts",
      "../src/common/input/keyboard.ts"
    ],
    [
      "./input/keyboard.test.d.ts",
      "../src/common/input/keyboard.test.ts"
    ]
  ],
  "size": 1695
}
//// [/user/username/projects/project/out/terminal.d.ts] *new* 
export {};
//# sourceMappingURL=terminal.d.ts.map
//// [/user/username/projects/project/out/terminal.d.ts.map] *new* 
{"version":3,"file":"terminal.d.ts","sourceRoot":"","sources":["../src/terminal.ts"],"names":[],"mappings":""}
//// [/user/username/projects/project/out/terminal.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const keyboard_1 = require("common/input/keyboard");
function foo() {
    return (0, keyboard_1.evaluateKeyboardEvent)();
}

//// [/user/username/projects/project/out/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,3]],"fileNames":["lib.d.ts","../src/common/input/keyboard.ts","../src/common/input/keyboard.test.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"6e37abb18685aa8e2642ab0218eef699-function bar() { return \"just a random function so .d.ts location doesnt match\"; }\nexport function evaluateKeyboardEvent() { }","signature":"569df1f274cf52f322c6e7a2f4e891fe-export declare function evaluateKeyboardEvent(): void;\n","impliedNodeFormat":1},{"version":"80dddc2f6e9d5c8a92919e5a78fab3c0-import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction testEvaluateKeyboardEvent() {\n    return evaluateKeyboardEvent();\n}","signature":"abe7d9981d6018efb6b2b794f40a1607-export {};\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"composite":true,"declarationMap":true,"outDir":"./"},"referencedMap":[[3,1]],"latestChangedDtsFile":"./input/keyboard.test.d.ts"}
//// [/user/username/projects/project/out/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../src/common/input/keyboard.ts",
        "../src/common/input/keyboard.test.ts"
      ],
      "original": [
        2,
        3
      ]
    }
  ],
  "fileNames": [
    "lib.d.ts",
    "../src/common/input/keyboard.ts",
    "../src/common/input/keyboard.test.ts"
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
      "fileName": "../src/common/input/keyboard.ts",
      "version": "6e37abb18685aa8e2642ab0218eef699-function bar() { return \"just a random function so .d.ts location doesnt match\"; }\nexport function evaluateKeyboardEvent() { }",
      "signature": "569df1f274cf52f322c6e7a2f4e891fe-export declare function evaluateKeyboardEvent(): void;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "6e37abb18685aa8e2642ab0218eef699-function bar() { return \"just a random function so .d.ts location doesnt match\"; }\nexport function evaluateKeyboardEvent() { }",
        "signature": "569df1f274cf52f322c6e7a2f4e891fe-export declare function evaluateKeyboardEvent(): void;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../src/common/input/keyboard.test.ts",
      "version": "80dddc2f6e9d5c8a92919e5a78fab3c0-import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction testEvaluateKeyboardEvent() {\n    return evaluateKeyboardEvent();\n}",
      "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "80dddc2f6e9d5c8a92919e5a78fab3c0-import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction testEvaluateKeyboardEvent() {\n    return evaluateKeyboardEvent();\n}",
        "signature": "abe7d9981d6018efb6b2b794f40a1607-export {};\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../src/common/input/keyboard.ts"
    ]
  ],
  "options": {
    "composite": true,
    "declarationMap": true,
    "outDir": "./"
  },
  "referencedMap": {
    "../src/common/input/keyboard.test.ts": [
      "../src/common/input/keyboard.ts"
    ]
  },
  "latestChangedDtsFile": "./input/keyboard.test.d.ts",
  "size": 1660
}
//// [/user/username/projects/project/src/common/input/keyboard.test.ts] *new* 
import { evaluateKeyboardEvent } from 'common/input/keyboard';
function testEvaluateKeyboardEvent() {
    return evaluateKeyboardEvent();
}
//// [/user/username/projects/project/src/common/input/keyboard.ts] *new* 
function bar() { return "just a random function so .d.ts location doesnt match"; }
export function evaluateKeyboardEvent() { }
//// [/user/username/projects/project/src/common/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
        "outDir": "../../out",
        "disableSourceOfProjectReferenceRedirect": true,
        "paths": {
            "*": ["../*"],
        },
    },
    "include": ["./**/*"]
}
//// [/user/username/projects/project/src/terminal.ts] *new* 
import { evaluateKeyboardEvent } from 'common/input/keyboard';
function foo() {
    return evaluateKeyboardEvent();
}
//// [/user/username/projects/project/src/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
        "outDir": "../out",
        "disableSourceOfProjectReferenceRedirect": true,
        "paths": {
            "common/*": ["./common/*"],
        },
        "tsBuildInfoFile": "../out/src.tsconfig.tsbuildinfo"
    },
    "include": ["./**/*"],
    "references": [
        { "path": "./common" },
    ],
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/project/src/common/input/keyboard.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "function bar() { return \"just a random function so .d.ts location doesnt match\"; }\nexport function evaluateKeyboardEvent() { }"
    }
  }
}
Projects::
  [/user/username/projects/project/src/common/tsconfig.json] *new*
    /user/username/projects/project/src/common/input/keyboard.ts       
    /user/username/projects/project/src/common/input/keyboard.test.ts  
  [/user/username/projects/project/src/tsconfig.json] *new*
Open Files::
  [/user/username/projects/project/src/common/input/keyboard.ts] *new*
    /user/username/projects/project/src/common/tsconfig.json  (default) 
Config::
  [/user/username/projects/project/src/common/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/project/src/common/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/project/src/common/input/keyboard.ts  
Config File Names::
  [/user/username/projects/project/src/common/input/keyboard.ts] *new*
    NearestConfigFileName: /user/username/projects/project/src/common/tsconfig.json
    Ancestors:
      /user/username/projects/project/src/common/tsconfig.json  /user/username/projects/project/src/tsconfig.json 
      /user/username/projects/project/src/tsconfig.json          
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/project/src/terminal.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction foo() {\n    return evaluateKeyboardEvent();\n}"
    }
  }
}
Projects::
  [/user/username/projects/project/src/common/tsconfig.json] 
    /user/username/projects/project/src/common/input/keyboard.ts       
    /user/username/projects/project/src/common/input/keyboard.test.ts  
  [/user/username/projects/project/src/tsconfig.json] *modified*
    /user/username/projects/project/out/input/keyboard.d.ts       *new*
    /user/username/projects/project/src/terminal.ts               *new*
    /user/username/projects/project/out/input/keyboard.test.d.ts  *new*
Open Files::
  [/user/username/projects/project/src/common/input/keyboard.ts] 
    /user/username/projects/project/src/common/tsconfig.json  (default) 
  [/user/username/projects/project/src/terminal.ts] *new*
    /user/username/projects/project/src/tsconfig.json  (default) 
Config::
  [/user/username/projects/project/src/common/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/project/src/common/tsconfig.json  
      /user/username/projects/project/src/tsconfig.json         *new*
    RetainingOpenFiles:
      /user/username/projects/project/src/common/input/keyboard.ts  
  [/user/username/projects/project/src/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/project/src/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/project/src/terminal.ts  
Config File Names::
  [/user/username/projects/project/src/common/input/keyboard.ts] 
    NearestConfigFileName: /user/username/projects/project/src/common/tsconfig.json
    Ancestors:
      /user/username/projects/project/src/common/tsconfig.json  /user/username/projects/project/src/tsconfig.json 
      /user/username/projects/project/src/tsconfig.json          
  [/user/username/projects/project/src/terminal.ts] *new*
    NearestConfigFileName: /user/username/projects/project/src/tsconfig.json
    Ancestors:
      /user/username/projects/project/src/tsconfig.json   
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/project/src/common/input/keyboard.ts"
    },
    "position": {
      "line": 1,
      "character": 16
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /user/username/projects/project/src/common/input/keyboard.test.ts ===
// import { [|evaluateKeyboardEvent|] } from 'common/input/keyboard';
// function testEvaluateKeyboardEvent() {
//     return [|evaluateKeyboardEvent|]();
// }

// === /user/username/projects/project/src/common/input/keyboard.ts ===
// function bar() { return "just a random function so .d.ts location doesnt match"; }
// export function /*FIND ALL REFS*/[|evaluateKeyboardEvent|]() { }

// === /user/username/projects/project/src/terminal.ts ===
// import { [|evaluateKeyboardEvent|] } from 'common/input/keyboard';
// function foo() {
//     return [|evaluateKeyboardEvent|]();
// }

currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/hub.ts] *new* 
export const value = 1;
//// [/home/src/workspaces/project/importer.ts] *new* 
import { value } from "./reexport";
console.log(value);
//// [/home/src/workspaces/project/reexport.ts] *new* 
export { value } from "./hub";
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "incremental": true,
        "noEmit": true,
        "module": "esnext",
    },
}

tsgo 
ExitStatus:: Success
Output::
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
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[[2,4]],"fileNames":["lib.es2025.full.d.ts","./hub.ts","./reexport.ts","./importer.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"9ae38e3a9bd5acd9f384aed0787571ff-export const value = 1;","d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";","23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);"],"fileIdsList":[[3],[2]],"options":{"module":99},"referencedMap":[[4,1],[3,2]],"affectedFilesPendingEmit":[2,4,3]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./hub.ts",
        "./reexport.ts",
        "./importer.ts"
      ],
      "original": [
        2,
        4
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./hub.ts",
    "./reexport.ts",
    "./importer.ts"
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
      "fileName": "./hub.ts",
      "version": "9ae38e3a9bd5acd9f384aed0787571ff-export const value = 1;",
      "signature": "9ae38e3a9bd5acd9f384aed0787571ff-export const value = 1;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./reexport.ts",
      "version": "d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";",
      "signature": "d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./importer.ts",
      "version": "23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);",
      "signature": "23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./reexport.ts"
    ],
    [
      "./hub.ts"
    ]
  ],
  "options": {
    "module": 99
  },
  "referencedMap": {
    "./importer.ts": [
      "./reexport.ts"
    ],
    "./reexport.ts": [
      "./hub.ts"
    ]
  },
  "affectedFilesPendingEmit": [
    [
      "./hub.ts",
      "Js",
      2
    ],
    [
      "./importer.ts",
      "Js",
      4
    ],
    [
      "./reexport.ts",
      "Js",
      3
    ]
  ],
  "size": 1240
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/hub.ts
*refresh*    /home/src/workspaces/project/reexport.ts
*refresh*    /home/src/workspaces/project/importer.ts
Signatures::


Edit [0]:: append ordinary comment
//// [/home/src/workspaces/project/hub.ts] *modified* 
export const value = 1;
// comment

tsgo 
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,4]],"fileNames":["lib.es2025.full.d.ts","./hub.ts","./reexport.ts","./importer.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"d30072c6a70ae5d26a3b9f1a0c63af03-export const value = 1;\n// comment","signature":"ba5c8ef2b0978fc944b312728d63c3f1-export declare const value = 1;\n","impliedNodeFormat":1},"d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";","23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);"],"fileIdsList":[[3],[2]],"options":{"module":99},"referencedMap":[[4,1],[3,2]],"affectedFilesPendingEmit":[2,4,3]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./hub.ts",
        "./reexport.ts",
        "./importer.ts"
      ],
      "original": [
        2,
        4
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./hub.ts",
    "./reexport.ts",
    "./importer.ts"
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
      "fileName": "./hub.ts",
      "version": "d30072c6a70ae5d26a3b9f1a0c63af03-export const value = 1;\n// comment",
      "signature": "ba5c8ef2b0978fc944b312728d63c3f1-export declare const value = 1;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d30072c6a70ae5d26a3b9f1a0c63af03-export const value = 1;\n// comment",
        "signature": "ba5c8ef2b0978fc944b312728d63c3f1-export declare const value = 1;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./reexport.ts",
      "version": "d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";",
      "signature": "d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./importer.ts",
      "version": "23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);",
      "signature": "23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./reexport.ts"
    ],
    [
      "./hub.ts"
    ]
  ],
  "options": {
    "module": 99
  },
  "referencedMap": {
    "./importer.ts": [
      "./reexport.ts"
    ],
    "./reexport.ts": [
      "./hub.ts"
    ]
  },
  "affectedFilesPendingEmit": [
    [
      "./hub.ts",
      "Js",
      2
    ],
    [
      "./importer.ts",
      "Js",
      4
    ],
    [
      "./reexport.ts",
      "Js",
      3
    ]
  ],
  "size": 1367
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/hub.ts
*refresh*    /home/src/workspaces/project/reexport.ts
*refresh*    /home/src/workspaces/project/importer.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/hub.ts
(used version)   /home/src/workspaces/project/importer.ts


Edit [1]:: change exported value
//// [/home/src/workspaces/project/hub.ts] *modified* 
export const value = 'one';
// comment

tsgo 
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[[2,4]],"fileNames":["lib.es2025.full.d.ts","./hub.ts","./reexport.ts","./importer.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"7406586c3ea25aed2dc08660dfb2cc13-export const value = 'one';\n// comment","signature":"e4163a3e4c4379367c549366a939940a-export declare const value = \"one\";\n","impliedNodeFormat":1},"d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";","23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);"],"fileIdsList":[[3],[2]],"options":{"module":99},"referencedMap":[[4,1],[3,2]],"affectedFilesPendingEmit":[2,4,3]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./hub.ts",
        "./reexport.ts",
        "./importer.ts"
      ],
      "original": [
        2,
        4
      ]
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./hub.ts",
    "./reexport.ts",
    "./importer.ts"
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
      "fileName": "./hub.ts",
      "version": "7406586c3ea25aed2dc08660dfb2cc13-export const value = 'one';\n// comment",
      "signature": "e4163a3e4c4379367c549366a939940a-export declare const value = \"one\";\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7406586c3ea25aed2dc08660dfb2cc13-export const value = 'one';\n// comment",
        "signature": "e4163a3e4c4379367c549366a939940a-export declare const value = \"one\";\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./reexport.ts",
      "version": "d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";",
      "signature": "d5041c1ab88dbbd3cd99d36379352deb-export { value } from \"./hub\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./importer.ts",
      "version": "23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);",
      "signature": "23fbccfab572480bb3fa6f41a1ea8483-import { value } from \"./reexport\";\nconsole.log(value);",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./reexport.ts"
    ],
    [
      "./hub.ts"
    ]
  ],
  "options": {
    "module": 99
  },
  "referencedMap": {
    "./importer.ts": [
      "./reexport.ts"
    ],
    "./reexport.ts": [
      "./hub.ts"
    ]
  },
  "affectedFilesPendingEmit": [
    [
      "./hub.ts",
      "Js",
      2
    ],
    [
      "./importer.ts",
      "Js",
      4
    ],
    [
      "./reexport.ts",
      "Js",
      3
    ]
  ],
  "size": 1377
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/hub.ts
*refresh*    /home/src/workspaces/project/reexport.ts
*refresh*    /home/src/workspaces/project/importer.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/hub.ts
(used version)   /home/src/workspaces/project/importer.ts

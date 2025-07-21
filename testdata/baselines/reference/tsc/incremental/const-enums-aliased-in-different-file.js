currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/a.ts] *new* 
import {A} from "./c"
let a = A.ONE
//// [/home/src/workspaces/project/b.d.ts] *new* 
export { AWorker as A } from "./worker";
//// [/home/src/workspaces/project/c.ts] *new* 
import {A} from "./b"
let b = A.ONE
export {A}
//// [/home/src/workspaces/project/worker.d.ts] *new* 
export const enum AWorker {
    ONE = 1
}

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
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
//// [/home/src/workspaces/project/a.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
let a = c_1.A.ONE;

//// [/home/src/workspaces/project/a.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./worker.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"6bbf56b58e48a1823a06aa76603127859d31869809dab154ab8c5903d6465602-export const enum AWorker {\n    ONE = 1\n}","d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";","f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}","f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE"],"fileIdsList":[[4],[2],[3]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[5,1],[3,2],[4,3]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./worker.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./worker.d.ts",
      "version": "6bbf56b58e48a1823a06aa76603127859d31869809dab154ab8c5903d6465602-export const enum AWorker {\n    ONE = 1\n}",
      "signature": "6bbf56b58e48a1823a06aa76603127859d31869809dab154ab8c5903d6465602-export const enum AWorker {\n    ONE = 1\n}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./b.d.ts",
      "version": "d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";",
      "signature": "d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "signature": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
    ],
    [
      "./worker.d.ts"
    ],
    [
      "./b.d.ts"
    ]
  ],
  "options": {
    "tsBuildInfoFile": "./a.tsbuildinfo"
  },
  "referencedMap": {
    "./a.ts": [
      "./c.ts"
    ],
    "./b.d.ts": [
      "./worker.d.ts"
    ],
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 1520
}
//// [/home/src/workspaces/project/c.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
let b = b_1.A.ONE;


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/worker.d.ts
*refresh*    /home/src/workspaces/project/b.d.ts
*refresh*    /home/src/workspaces/project/c.ts
*refresh*    /home/src/workspaces/project/a.ts
Signatures::


Edit [0]:: change enum value
//// [/home/src/workspaces/project/worker.d.ts] *modified* 
export const enum AWorker {
    ONE = 2
}

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.js] *modified time*
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./worker.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"4749e1d328fb4b3d350a6fd1629b976d85cb877d3fb71070eef63605206f3a60-export const enum AWorker {\n    ONE = 2\n}","d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";","f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}","f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE"],"fileIdsList":[[4],[2],[3]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[5,1],[3,2],[4,3]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./worker.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./worker.d.ts",
      "version": "4749e1d328fb4b3d350a6fd1629b976d85cb877d3fb71070eef63605206f3a60-export const enum AWorker {\n    ONE = 2\n}",
      "signature": "4749e1d328fb4b3d350a6fd1629b976d85cb877d3fb71070eef63605206f3a60-export const enum AWorker {\n    ONE = 2\n}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./b.d.ts",
      "version": "d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";",
      "signature": "d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "signature": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
    ],
    [
      "./worker.d.ts"
    ],
    [
      "./b.d.ts"
    ]
  ],
  "options": {
    "tsBuildInfoFile": "./a.tsbuildinfo"
  },
  "referencedMap": {
    "./a.ts": [
      "./c.ts"
    ],
    "./b.d.ts": [
      "./worker.d.ts"
    ],
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 1520
}
//// [/home/src/workspaces/project/c.js] *modified time*

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/worker.d.ts
*refresh*    /home/src/workspaces/project/b.d.ts
*refresh*    /home/src/workspaces/project/c.ts
*refresh*    /home/src/workspaces/project/a.ts
Signatures::
(used version)   /home/src/workspaces/project/worker.d.ts
(used version)   /home/src/workspaces/project/b.d.ts
(used version)   /home/src/workspaces/project/c.ts
(used version)   /home/src/workspaces/project/a.ts


Edit [1]:: change enum value again
//// [/home/src/workspaces/project/worker.d.ts] *modified* 
export const enum AWorker {
    ONE = 3
}

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.js] *modified time*
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./worker.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}","d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";","f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}","f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE"],"fileIdsList":[[4],[2],[3]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[5,1],[3,2],[4,3]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./worker.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./worker.d.ts",
      "version": "9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}",
      "signature": "9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./b.d.ts",
      "version": "d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";",
      "signature": "d1cae7e0884009696a49cf892aeb5da4fc8bc6050db4bce9ee30b657ccba9636-export { AWorker as A } from \"./worker\";",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "signature": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
    ],
    [
      "./worker.d.ts"
    ],
    [
      "./b.d.ts"
    ]
  ],
  "options": {
    "tsBuildInfoFile": "./a.tsbuildinfo"
  },
  "referencedMap": {
    "./a.ts": [
      "./c.ts"
    ],
    "./b.d.ts": [
      "./worker.d.ts"
    ],
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 1520
}
//// [/home/src/workspaces/project/c.js] *modified time*

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/worker.d.ts
*refresh*    /home/src/workspaces/project/b.d.ts
*refresh*    /home/src/workspaces/project/c.ts
*refresh*    /home/src/workspaces/project/a.ts
Signatures::
(used version)   /home/src/workspaces/project/worker.d.ts
(used version)   /home/src/workspaces/project/b.d.ts
(used version)   /home/src/workspaces/project/c.ts
(used version)   /home/src/workspaces/project/a.ts


Edit [2]:: something else changes in b.d.ts
//// [/home/src/workspaces/project/b.d.ts] *modified* 
export { AWorker as A } from "./worker";export const randomThing = 10;

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.js] *modified time*
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./worker.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}","eececfa518b405636ca59764618daf54087efc29f152e951a20f543569bc9da6-export { AWorker as A } from \"./worker\";export const randomThing = 10;",{"version":"f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}","signature":"1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732-import { A } from \"./b\";\nexport { A };\n","impliedNodeFormat":1},{"version":"f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n","impliedNodeFormat":1}],"fileIdsList":[[4],[2],[3]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[5,1],[3,2],[4,3]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./worker.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./worker.d.ts",
      "version": "9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}",
      "signature": "9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./b.d.ts",
      "version": "eececfa518b405636ca59764618daf54087efc29f152e951a20f543569bc9da6-export { AWorker as A } from \"./worker\";export const randomThing = 10;",
      "signature": "eececfa518b405636ca59764618daf54087efc29f152e951a20f543569bc9da6-export { AWorker as A } from \"./worker\";export const randomThing = 10;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732-import { A } from \"./b\";\nexport { A };\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
        "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732-import { A } from \"./b\";\nexport { A };\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
        "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
    ],
    [
      "./worker.d.ts"
    ],
    [
      "./b.d.ts"
    ]
  ],
  "options": {
    "tsBuildInfoFile": "./a.tsbuildinfo"
  },
  "referencedMap": {
    "./a.ts": [
      "./c.ts"
    ],
    "./b.d.ts": [
      "./worker.d.ts"
    ],
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 1833
}
//// [/home/src/workspaces/project/c.js] *modified time*

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/b.d.ts
*refresh*    /home/src/workspaces/project/c.ts
*refresh*    /home/src/workspaces/project/a.ts
Signatures::
(used version)   /home/src/workspaces/project/b.d.ts
(computed .d.ts) /home/src/workspaces/project/c.ts
(computed .d.ts) /home/src/workspaces/project/a.ts


Edit [3]:: something else changes in b.d.ts again
//// [/home/src/workspaces/project/b.d.ts] *modified* 
export { AWorker as A } from "./worker";export const randomThing = 10;export const randomThing2 = 10;

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./worker.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}","fc3df7c9dd67a05ea3329c1a9ef9c5f81ddb2109cc13c2089aef6c5dbef34340-export { AWorker as A } from \"./worker\";export const randomThing = 10;export const randomThing2 = 10;",{"version":"f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}","signature":"1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732-import { A } from \"./b\";\nexport { A };\n","impliedNodeFormat":1},"f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE"],"fileIdsList":[[4],[2],[3]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[5,1],[3,2],[4,3]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./worker.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./worker.d.ts",
      "version": "9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}",
      "signature": "9429bb2451db350f9bceba9041b64c9b80446ebfd5c423802bdee3ac2cd0fccc-export const enum AWorker {\n    ONE = 3\n}",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./b.d.ts",
      "version": "fc3df7c9dd67a05ea3329c1a9ef9c5f81ddb2109cc13c2089aef6c5dbef34340-export { AWorker as A } from \"./worker\";export const randomThing = 10;export const randomThing2 = 10;",
      "signature": "fc3df7c9dd67a05ea3329c1a9ef9c5f81ddb2109cc13c2089aef6c5dbef34340-export { AWorker as A } from \"./worker\";export const randomThing = 10;export const randomThing2 = 10;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
      "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732-import { A } from \"./b\";\nexport { A };\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a-import {A} from \"./b\"\nlet b = A.ONE\nexport {A}",
        "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732-import { A } from \"./b\";\nexport { A };\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09-import {A} from \"./c\"\nlet a = A.ONE",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
    ],
    [
      "./worker.d.ts"
    ],
    [
      "./b.d.ts"
    ]
  ],
  "options": {
    "tsBuildInfoFile": "./a.tsbuildinfo"
  },
  "referencedMap": {
    "./a.ts": [
      "./c.ts"
    ],
    "./b.d.ts": [
      "./worker.d.ts"
    ],
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 1738
}
//// [/home/src/workspaces/project/c.js] *modified time*

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/b.d.ts
*refresh*    /home/src/workspaces/project/c.ts
*refresh*    /home/src/workspaces/project/a.ts
Signatures::
(used version)   /home/src/workspaces/project/b.d.ts
(computed .d.ts) /home/src/workspaces/project/c.ts
(used version)   /home/src/workspaces/project/a.ts

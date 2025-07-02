currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/class1.ts] *new* 
const a: MagicNumber = 1;
console.log(a);
//// [/home/src/workspaces/project/constants.ts] *new* 
export default 1;
//// [/home/src/workspaces/project/reexport.ts] *new* 
export { default as ConstantNumber } from "./constants"
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true
    }
}
//// [/home/src/workspaces/project/types.d.ts] *new* 
type MagicNumber = typeof import('./reexport').ConstantNumber

tsgo 
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
//// [/home/src/workspaces/project/class1.d.ts] *new* 
declare const a = 1;

//// [/home/src/workspaces/project/class1.js] *new* 
const a = 1;
console.log(a);

//// [/home/src/workspaces/project/constants.d.ts] *new* 
declare const _default: number;
export default _default;

//// [/home/src/workspaces/project/constants.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = 1;

//// [/home/src/workspaces/project/reexport.d.ts] *new* 
export { default as ConstantNumber } from "./constants";

//// [/home/src/workspaces/project/reexport.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConstantNumber = void 0;
const constants_1 = require("./constants");
Object.defineProperty(exports, "ConstantNumber", { enumerable: true, get: function () { return constants_1.default; } });

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./class1.ts","./constants.ts","./reexport.ts","./types.d.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e23285158ce57995e96701e24b3979513455090a1acd8091dfcc03759c20f3a7","signature":"48becb2a5e6a58c58a56481b0edb338fd157be16188619433b47b766a8a71b60","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"56332e0a55734bc2b73df56a2df8635ed5c5b24b6d7a456b41de7cab9a2f3814","signature":"36b18171f705a4860606c23c78134f2e86f6e695d994dc8c849df608075d2e5a","impliedNodeFormat":1},{"version":"d358f9090c6427d5f6dd68671cba18871b08d2fcb193da8e979e324d68000cb3","signature":"84ff27208288e3f7379462f8745935fe8205e24ac2e2ef7c9ddef2dc2fbfeb39","impliedNodeFormat":1},{"version":"4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[3],[4]],"options":{"composite":true},"referencedMap":[[4,1],[5,2]],"latestChangedDtsFile":"./reexport.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./class1.ts",
    "./constants.ts",
    "./reexport.ts",
    "./types.d.ts"
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
      "version": "e23285158ce57995e96701e24b3979513455090a1acd8091dfcc03759c20f3a7",
      "signature": "48becb2a5e6a58c58a56481b0edb338fd157be16188619433b47b766a8a71b60",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e23285158ce57995e96701e24b3979513455090a1acd8091dfcc03759c20f3a7",
        "signature": "48becb2a5e6a58c58a56481b0edb338fd157be16188619433b47b766a8a71b60",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./constants.ts",
      "version": "56332e0a55734bc2b73df56a2df8635ed5c5b24b6d7a456b41de7cab9a2f3814",
      "signature": "36b18171f705a4860606c23c78134f2e86f6e695d994dc8c849df608075d2e5a",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "56332e0a55734bc2b73df56a2df8635ed5c5b24b6d7a456b41de7cab9a2f3814",
        "signature": "36b18171f705a4860606c23c78134f2e86f6e695d994dc8c849df608075d2e5a",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./reexport.ts",
      "version": "d358f9090c6427d5f6dd68671cba18871b08d2fcb193da8e979e324d68000cb3",
      "signature": "84ff27208288e3f7379462f8745935fe8205e24ac2e2ef7c9ddef2dc2fbfeb39",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d358f9090c6427d5f6dd68671cba18871b08d2fcb193da8e979e324d68000cb3",
        "signature": "84ff27208288e3f7379462f8745935fe8205e24ac2e2ef7c9ddef2dc2fbfeb39",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./types.d.ts",
      "version": "4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2",
      "signature": "4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./constants.ts"
    ],
    [
      "./reexport.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./reexport.ts": [
      "./constants.ts"
    ],
    "./types.d.ts": [
      "./reexport.ts"
    ]
  },
  "latestChangedDtsFile": "./reexport.d.ts",
  "size": 1092
}

SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/class1.ts
*refresh*    /home/src/workspaces/project/constants.ts
*refresh*    /home/src/workspaces/project/reexport.ts
*refresh*    /home/src/workspaces/project/types.d.ts
Signatures::
(stored at emit) /home/src/workspaces/project/class1.ts
(stored at emit) /home/src/workspaces/project/constants.ts
(stored at emit) /home/src/workspaces/project/reexport.ts


Edit [0]:: Modify imports used in global file
//// [/home/src/workspaces/project/constants.ts] *modified* 
export default 2;

tsgo 
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/constants.js] *modified* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = 2;

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./class1.ts","./constants.ts","./reexport.ts","./types.d.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e23285158ce57995e96701e24b3979513455090a1acd8091dfcc03759c20f3a7","signature":"48becb2a5e6a58c58a56481b0edb338fd157be16188619433b47b766a8a71b60","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"a903ba200e1d43efd9a49f4a8f57c622efb2ca72b9a00222dacd16d6ba6f3ba0","signature":"36b18171f705a4860606c23c78134f2e86f6e695d994dc8c849df608075d2e5a","impliedNodeFormat":1},{"version":"d358f9090c6427d5f6dd68671cba18871b08d2fcb193da8e979e324d68000cb3","signature":"84ff27208288e3f7379462f8745935fe8205e24ac2e2ef7c9ddef2dc2fbfeb39","impliedNodeFormat":1},{"version":"4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[3],[4]],"options":{"composite":true},"referencedMap":[[4,1],[5,2]],"latestChangedDtsFile":"./reexport.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./class1.ts",
    "./constants.ts",
    "./reexport.ts",
    "./types.d.ts"
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
      "version": "e23285158ce57995e96701e24b3979513455090a1acd8091dfcc03759c20f3a7",
      "signature": "48becb2a5e6a58c58a56481b0edb338fd157be16188619433b47b766a8a71b60",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e23285158ce57995e96701e24b3979513455090a1acd8091dfcc03759c20f3a7",
        "signature": "48becb2a5e6a58c58a56481b0edb338fd157be16188619433b47b766a8a71b60",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./constants.ts",
      "version": "a903ba200e1d43efd9a49f4a8f57c622efb2ca72b9a00222dacd16d6ba6f3ba0",
      "signature": "36b18171f705a4860606c23c78134f2e86f6e695d994dc8c849df608075d2e5a",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a903ba200e1d43efd9a49f4a8f57c622efb2ca72b9a00222dacd16d6ba6f3ba0",
        "signature": "36b18171f705a4860606c23c78134f2e86f6e695d994dc8c849df608075d2e5a",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./reexport.ts",
      "version": "d358f9090c6427d5f6dd68671cba18871b08d2fcb193da8e979e324d68000cb3",
      "signature": "84ff27208288e3f7379462f8745935fe8205e24ac2e2ef7c9ddef2dc2fbfeb39",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d358f9090c6427d5f6dd68671cba18871b08d2fcb193da8e979e324d68000cb3",
        "signature": "84ff27208288e3f7379462f8745935fe8205e24ac2e2ef7c9ddef2dc2fbfeb39",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./types.d.ts",
      "version": "4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2",
      "signature": "4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4364705e72bd2d53c70faabaa6d59be927f81efc6cb62cdfb1b7b18c72ec95f2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./constants.ts"
    ],
    [
      "./reexport.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./reexport.ts": [
      "./constants.ts"
    ],
    "./types.d.ts": [
      "./reexport.ts"
    ]
  },
  "latestChangedDtsFile": "./reexport.d.ts",
  "size": 1092
}

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/constants.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/constants.ts

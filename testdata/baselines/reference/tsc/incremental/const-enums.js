currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/a.ts] *new* 
import {A} from "./c"
let a = A.ONE
//// [/home/src/workspaces/project/b.d.ts] *new* 
export const enum A {
    ONE = 1
}
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
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"c00e44762baea954363536f3ccceae15fd84bf2b246f8cb077164d106cd5c9ae","f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a","f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09"],"fileIdsList":[[3],[2]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[4,1],[3,2]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
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
      "fileName": "./b.d.ts",
      "version": "c00e44762baea954363536f3ccceae15fd84bf2b246f8cb077164d106cd5c9ae",
      "signature": "c00e44762baea954363536f3ccceae15fd84bf2b246f8cb077164d106cd5c9ae",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
      "signature": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
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
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 545
}
//// [/home/src/workspaces/project/c.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
let b = b_1.A.ONE;


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/b.d.ts
*refresh*    /home/src/workspaces/project/c.ts
*refresh*    /home/src/workspaces/project/a.ts
Signatures::


Edit [0]:: change enum value
//// [/home/src/workspaces/project/b.d.ts] *modified* 
export const enum A {
    ONE = 2
}

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.js] *modified time*
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"6ab9b377a00d6157a880e7be8f490a45f91a751d3aab7a46a21fd0169189e21c",{"version":"f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a","signature":"1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732","impliedNodeFormat":1},{"version":"f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881","impliedNodeFormat":1}],"fileIdsList":[[3],[2]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[4,1],[3,2]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
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
      "fileName": "./b.d.ts",
      "version": "6ab9b377a00d6157a880e7be8f490a45f91a751d3aab7a46a21fd0169189e21c",
      "signature": "6ab9b377a00d6157a880e7be8f490a45f91a751d3aab7a46a21fd0169189e21c",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
      "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
        "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
        "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
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
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 771
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


Edit [1]:: change enum value again
//// [/home/src/workspaces/project/b.d.ts] *modified* 
export const enum A {
    ONE = 3
}

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.js] *modified time*
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"9b7fcb1de7abaa3bee27ad197ab712966678462da9b016fd18a98202a022e062",{"version":"f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a","signature":"1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732","impliedNodeFormat":1},"f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09"],"fileIdsList":[[3],[2]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[4,1],[3,2]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
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
      "fileName": "./b.d.ts",
      "version": "9b7fcb1de7abaa3bee27ad197ab712966678462da9b016fd18a98202a022e062",
      "signature": "9b7fcb1de7abaa3bee27ad197ab712966678462da9b016fd18a98202a022e062",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
      "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
        "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
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
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 658
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


Edit [2]:: something else changes in b.d.ts
//// [/home/src/workspaces/project/b.d.ts] *modified* 
export const enum A {
    ONE = 3
}export const randomThing = 10;

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.js] *modified time*
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"8eb5f9388406068b4f938250bc5ba6c8f62cb17db7bd81dad8d2b8e411d344ed",{"version":"f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a","signature":"1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732","impliedNodeFormat":1},"f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09"],"fileIdsList":[[3],[2]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[4,1],[3,2]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
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
      "fileName": "./b.d.ts",
      "version": "8eb5f9388406068b4f938250bc5ba6c8f62cb17db7bd81dad8d2b8e411d344ed",
      "signature": "8eb5f9388406068b4f938250bc5ba6c8f62cb17db7bd81dad8d2b8e411d344ed",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
      "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
        "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
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
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 658
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


Edit [3]:: something else changes in b.d.ts again
//// [/home/src/workspaces/project/b.d.ts] *modified* 
export const enum A {
    ONE = 3
}export const randomThing = 10;export const randomThing2 = 10;

tsgo -i a.ts --tsbuildinfofile a.tsbuildinfo
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/a.js] *modified time*
//// [/home/src/workspaces/project/a.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./b.d.ts","./c.ts","./a.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"0e2a47b6ac5d1eb00a443643049e47e98f3bcacb8852a925f600902222323516",{"version":"f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a","signature":"1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732","impliedNodeFormat":1},"f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09"],"fileIdsList":[[3],[2]],"options":{"tsBuildInfoFile":"./a.tsbuildinfo"},"referencedMap":[[4,1],[3,2]]}
//// [/home/src/workspaces/project/a.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./b.d.ts",
    "./c.ts",
    "./a.ts"
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
      "fileName": "./b.d.ts",
      "version": "0e2a47b6ac5d1eb00a443643049e47e98f3bcacb8852a925f600902222323516",
      "signature": "0e2a47b6ac5d1eb00a443643049e47e98f3bcacb8852a925f600902222323516",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./c.ts",
      "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
      "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "f96307fd3b2531524ba186a0b7862e9605da65828664c51363d0b608f8141c8a",
        "signature": "1affdd1113604735d4499c03d6271d13972094ddab6991610e72d53c00d14732",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./a.ts",
      "version": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "signature": "f5a433d8f46180a7988f1820c3e70520cbbc864a870e4f6cdd4857edf3688e09",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./c.ts"
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
    "./c.ts": [
      "./b.d.ts"
    ]
  },
  "size": 658
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

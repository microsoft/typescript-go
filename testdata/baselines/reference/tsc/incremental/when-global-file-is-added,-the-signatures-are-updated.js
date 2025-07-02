
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/src/anotherFileWithSameReferenes.ts] *new* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function anotherFileWithSameReferenes() { }
//// [/home/src/workspaces/project/src/filePresent.ts] *new* 
function something() { return 10; }
//// [/home/src/workspaces/project/src/main.ts] *new* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": { "composite": true },
    "include": ["src/**/*.ts"],
}

ExitStatus:: 0

CompilerOptions::{}
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
//// [/home/src/workspaces/project/src/anotherFileWithSameReferenes.d.ts] *new* 
declare function anotherFileWithSameReferenes(): void;

//// [/home/src/workspaces/project/src/anotherFileWithSameReferenes.js] *new* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function anotherFileWithSameReferenes() { }

//// [/home/src/workspaces/project/src/filePresent.d.ts] *new* 
declare function something(): number;

//// [/home/src/workspaces/project/src/filePresent.js] *new* 
function something() { return 10; }

//// [/home/src/workspaces/project/src/main.d.ts] *new* 
declare function main(): void;

//// [/home/src/workspaces/project/src/main.js] *new* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/filePresent.ts","./src/anotherFileWithSameReferenes.ts","./src/main.ts","./src/fileNotFound.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5","signature":"4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d","signature":"d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"9f4be7f0076e4a80708374aafbed16d8c40107ca5556a6155e1c3b9d2a4bfd1c","signature":"ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[2,5]],"options":{"composite":true},"referencedMap":[[3,1],[4,1]],"latestChangedDtsFile":"./src/main.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/filePresent.ts",
    "./src/anotherFileWithSameReferenes.ts",
    "./src/main.ts",
    "./src/fileNotFound.ts"
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
      "fileName": "./src/filePresent.ts",
      "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
      "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
        "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/anotherFileWithSameReferenes.ts",
      "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
      "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
        "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/main.ts",
      "version": "9f4be7f0076e4a80708374aafbed16d8c40107ca5556a6155e1c3b9d2a4bfd1c",
      "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "9f4be7f0076e4a80708374aafbed16d8c40107ca5556a6155e1c3b9d2a4bfd1c",
        "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./src/anotherFileWithSameReferenes.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    "./src/main.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ]
  },
  "latestChangedDtsFile": "./src/main.d.ts",
  "size": 1056
}


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/src/filePresent.ts
*refresh*    /home/src/workspaces/project/src/anotherFileWithSameReferenes.ts
*refresh*    /home/src/workspaces/project/src/main.ts

Signatures::
(stored at emit) /home/src/workspaces/project/src/filePresent.ts
(stored at emit) /home/src/workspaces/project/src/anotherFileWithSameReferenes.ts
(stored at emit) /home/src/workspaces/project/src/main.ts


Edit:: no change

ExitStatus:: 0
Output::


SemanticDiagnostics::

Signatures::


Edit:: Modify main file
//// [/home/src/workspaces/project/src/main.ts] *modified* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }something();

ExitStatus:: 0
Output::
//// [/home/src/workspaces/project/src/main.js] *modified* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }
something();

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/filePresent.ts","./src/anotherFileWithSameReferenes.ts","./src/main.ts","./src/fileNotFound.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5","signature":"4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d","signature":"d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"958bd60a28f4c68b48b6efabb4498a30aae1c5f7207bbc2ef3e6b639a76075e3","signature":"ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[2,5]],"options":{"composite":true},"referencedMap":[[3,1],[4,1]],"latestChangedDtsFile":"./src/main.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/filePresent.ts",
    "./src/anotherFileWithSameReferenes.ts",
    "./src/main.ts",
    "./src/fileNotFound.ts"
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
      "fileName": "./src/filePresent.ts",
      "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
      "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
        "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/anotherFileWithSameReferenes.ts",
      "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
      "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
        "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/main.ts",
      "version": "958bd60a28f4c68b48b6efabb4498a30aae1c5f7207bbc2ef3e6b639a76075e3",
      "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "958bd60a28f4c68b48b6efabb4498a30aae1c5f7207bbc2ef3e6b639a76075e3",
        "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./src/anotherFileWithSameReferenes.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    "./src/main.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ]
  },
  "latestChangedDtsFile": "./src/main.d.ts",
  "size": 1056
}


SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/main.ts

Signatures::
(computed .d.ts) /home/src/workspaces/project/src/main.ts


Edit:: Modify main file again
//// [/home/src/workspaces/project/src/main.ts] *modified* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }something();something();

ExitStatus:: 0
Output::
//// [/home/src/workspaces/project/src/main.js] *modified* 
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }
something();
something();

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/filePresent.ts","./src/anotherFileWithSameReferenes.ts","./src/main.ts","./src/fileNotFound.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5","signature":"4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d","signature":"d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"b6d6648dd2368dd226aaed824ed05d767f851bc8c7346baf6478588fc8c53f6d","signature":"ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[2,5]],"options":{"composite":true},"referencedMap":[[3,1],[4,1]],"latestChangedDtsFile":"./src/main.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/filePresent.ts",
    "./src/anotherFileWithSameReferenes.ts",
    "./src/main.ts",
    "./src/fileNotFound.ts"
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
      "fileName": "./src/filePresent.ts",
      "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
      "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
        "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/anotherFileWithSameReferenes.ts",
      "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
      "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
        "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/main.ts",
      "version": "b6d6648dd2368dd226aaed824ed05d767f851bc8c7346baf6478588fc8c53f6d",
      "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "b6d6648dd2368dd226aaed824ed05d767f851bc8c7346baf6478588fc8c53f6d",
        "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./src/anotherFileWithSameReferenes.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    "./src/main.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ]
  },
  "latestChangedDtsFile": "./src/main.d.ts",
  "size": 1056
}


SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/main.ts

Signatures::
(computed .d.ts) /home/src/workspaces/project/src/main.ts


Edit:: Add new file and update main file
//// [/home/src/workspaces/project/src/main.ts] *modified* 
/// <reference path="./newFile.ts"/>
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }something();something();foo();
//// [/home/src/workspaces/project/src/newFile.ts] *new* 
function foo() { return 20; }

ExitStatus:: 0
Output::
//// [/home/src/workspaces/project/src/main.js] *modified* 
/// <reference path="./newFile.ts"/>
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }
something();
something();
foo();

//// [/home/src/workspaces/project/src/newFile.d.ts] *new* 
declare function foo(): number;

//// [/home/src/workspaces/project/src/newFile.js] *new* 
function foo() { return 20; }

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/filePresent.ts","./src/anotherFileWithSameReferenes.ts","./src/newFile.ts","./src/main.ts","./src/fileNotFound.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5","signature":"4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d","signature":"d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1","signature":"f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"930f3c7023d87506648db3c352e86f7a2baf38e91f743773fddf4a520aff393a","signature":"ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[2,6],[2,4,6]],"options":{"composite":true},"referencedMap":[[3,1],[5,2]],"latestChangedDtsFile":"./src/newFile.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/filePresent.ts",
    "./src/anotherFileWithSameReferenes.ts",
    "./src/newFile.ts",
    "./src/main.ts",
    "./src/fileNotFound.ts"
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
      "fileName": "./src/filePresent.ts",
      "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
      "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
        "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/anotherFileWithSameReferenes.ts",
      "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
      "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
        "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/newFile.ts",
      "version": "a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1",
      "signature": "f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1",
        "signature": "f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/main.ts",
      "version": "930f3c7023d87506648db3c352e86f7a2baf38e91f743773fddf4a520aff393a",
      "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "930f3c7023d87506648db3c352e86f7a2baf38e91f743773fddf4a520aff393a",
        "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    [
      "./src/filePresent.ts",
      "./src/newFile.ts",
      "./src/fileNotFound.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./src/anotherFileWithSameReferenes.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    "./src/main.ts": [
      "./src/filePresent.ts",
      "./src/newFile.ts",
      "./src/fileNotFound.ts"
    ]
  },
  "latestChangedDtsFile": "./src/newFile.d.ts",
  "size": 1292
}


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/src/newFile.ts
*refresh*    /home/src/workspaces/project/src/main.ts

Signatures::
(computed .d.ts) /home/src/workspaces/project/src/newFile.ts
(computed .d.ts) /home/src/workspaces/project/src/main.ts


Edit:: Write file that could not be resolved
//// [/home/src/workspaces/project/src/fileNotFound.ts] *new* 
function something2() { return 20; }

ExitStatus:: 0
Output::
//// [/home/src/workspaces/project/src/anotherFileWithSameReferenes.js] *modified time*
//// [/home/src/workspaces/project/src/fileNotFound.d.ts] *new* 
declare function something2(): number;

//// [/home/src/workspaces/project/src/fileNotFound.js] *new* 
function something2() { return 20; }

//// [/home/src/workspaces/project/src/main.js] *modified time*
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/filePresent.ts","./src/fileNotFound.ts","./src/anotherFileWithSameReferenes.ts","./src/newFile.ts","./src/main.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5","signature":"4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"6c5b229cbf53b2c6867ab14b139eeac37ed3ec0c1564ba561f7faa869aaba32c","signature":"14ba6a62cd6d3e47b343358b2c3e1b7e34b488b489f0d9b915c796cd2e61bbad","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d","signature":"d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1","signature":"f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"930f3c7023d87506648db3c352e86f7a2baf38e91f743773fddf4a520aff393a","signature":"ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[2,3],[2,3,5]],"options":{"composite":true},"referencedMap":[[4,1],[6,2]],"latestChangedDtsFile":"./src/fileNotFound.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/filePresent.ts",
    "./src/fileNotFound.ts",
    "./src/anotherFileWithSameReferenes.ts",
    "./src/newFile.ts",
    "./src/main.ts"
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
      "fileName": "./src/filePresent.ts",
      "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
      "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
        "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/fileNotFound.ts",
      "version": "6c5b229cbf53b2c6867ab14b139eeac37ed3ec0c1564ba561f7faa869aaba32c",
      "signature": "14ba6a62cd6d3e47b343358b2c3e1b7e34b488b489f0d9b915c796cd2e61bbad",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "6c5b229cbf53b2c6867ab14b139eeac37ed3ec0c1564ba561f7faa869aaba32c",
        "signature": "14ba6a62cd6d3e47b343358b2c3e1b7e34b488b489f0d9b915c796cd2e61bbad",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/anotherFileWithSameReferenes.ts",
      "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
      "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
        "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/newFile.ts",
      "version": "a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1",
      "signature": "f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1",
        "signature": "f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/main.ts",
      "version": "930f3c7023d87506648db3c352e86f7a2baf38e91f743773fddf4a520aff393a",
      "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "930f3c7023d87506648db3c352e86f7a2baf38e91f743773fddf4a520aff393a",
        "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts",
      "./src/newFile.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./src/anotherFileWithSameReferenes.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    "./src/main.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts",
      "./src/newFile.ts"
    ]
  },
  "latestChangedDtsFile": "./src/fileNotFound.d.ts",
  "size": 1503
}


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/src/fileNotFound.ts
*refresh*    /home/src/workspaces/project/src/anotherFileWithSameReferenes.ts
*refresh*    /home/src/workspaces/project/src/main.ts

Signatures::
(computed .d.ts) /home/src/workspaces/project/src/fileNotFound.ts
(computed .d.ts) /home/src/workspaces/project/src/anotherFileWithSameReferenes.ts
(computed .d.ts) /home/src/workspaces/project/src/main.ts


Edit:: Modify main file
//// [/home/src/workspaces/project/src/main.ts] *modified* 
/// <reference path="./newFile.ts"/>
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }something();something();foo();something();

ExitStatus:: 0
Output::
//// [/home/src/workspaces/project/src/main.js] *modified* 
/// <reference path="./newFile.ts"/>
/// <reference path="./filePresent.ts"/>
/// <reference path="./fileNotFound.ts"/>
function main() { }
something();
something();
foo();
something();

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/filePresent.ts","./src/fileNotFound.ts","./src/anotherFileWithSameReferenes.ts","./src/newFile.ts","./src/main.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5","signature":"4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"6c5b229cbf53b2c6867ab14b139eeac37ed3ec0c1564ba561f7faa869aaba32c","signature":"14ba6a62cd6d3e47b343358b2c3e1b7e34b488b489f0d9b915c796cd2e61bbad","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d","signature":"d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1","signature":"f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"a8ae65557452d18998812e950c84bacd81f817ef323ccfe777c3b5450519c167","signature":"ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896","affectsGlobalScope":true,"impliedNodeFormat":1}],"fileIdsList":[[2,3],[2,3,5]],"options":{"composite":true},"referencedMap":[[4,1],[6,2]],"latestChangedDtsFile":"./src/fileNotFound.d.ts"}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/filePresent.ts",
    "./src/fileNotFound.ts",
    "./src/anotherFileWithSameReferenes.ts",
    "./src/newFile.ts",
    "./src/main.ts"
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
      "fileName": "./src/filePresent.ts",
      "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
      "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8f597c315d3bba69a79551601042fdcfe05d35f763762db4908b255c7f17c7d5",
        "signature": "4f3eeb0c183707d474ecb20d55e49f78d8a3fa3ac388d3b7a318d603ad8478c2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/fileNotFound.ts",
      "version": "6c5b229cbf53b2c6867ab14b139eeac37ed3ec0c1564ba561f7faa869aaba32c",
      "signature": "14ba6a62cd6d3e47b343358b2c3e1b7e34b488b489f0d9b915c796cd2e61bbad",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "6c5b229cbf53b2c6867ab14b139eeac37ed3ec0c1564ba561f7faa869aaba32c",
        "signature": "14ba6a62cd6d3e47b343358b2c3e1b7e34b488b489f0d9b915c796cd2e61bbad",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/anotherFileWithSameReferenes.ts",
      "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
      "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "552b902790cfbb88a7ed6a7b04b24bb18c58a7f52bcfe7808912e126359e258d",
        "signature": "d4dbe375786b0b36d5425a70f140bbb3d377883027d2fa29aa022a2bd446fbda",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/newFile.ts",
      "version": "a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1",
      "signature": "f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a4c88e3619994da0f5e4da2dc210f6038e710b9bb831003767da68c882137fb1",
        "signature": "f0d67d5e01f8fff5f52028627fc0fb5a78b24df03e482ddac513fa1f873934ee",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/main.ts",
      "version": "a8ae65557452d18998812e950c84bacd81f817ef323ccfe777c3b5450519c167",
      "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "a8ae65557452d18998812e950c84bacd81f817ef323ccfe777c3b5450519c167",
        "signature": "ed4b087ea2a2e4a58647864cf512c7534210bfc2f9d236a2f9ed5245cf7a0896",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts",
      "./src/newFile.ts"
    ]
  ],
  "options": {
    "composite": true
  },
  "referencedMap": {
    "./src/anotherFileWithSameReferenes.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts"
    ],
    "./src/main.ts": [
      "./src/filePresent.ts",
      "./src/fileNotFound.ts",
      "./src/newFile.ts"
    ]
  },
  "latestChangedDtsFile": "./src/fileNotFound.d.ts",
  "size": 1503
}


SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/src/main.ts

Signatures::
(computed .d.ts) /home/src/workspaces/project/src/main.ts

currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/a.ts] *new* 
const items = [1, 2, 3];
// The spread ...items overwrites the [Symbol.iterator] property specified above it
const obj = {
    length: items.length,
    [Symbol.iterator]: function* () {
        for (const item of items) yield item;
    },
    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'
};
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "incremental": true,
        "strict": true,
        "target": "esnext",
        "declaration": true,
        "outDir": "./out",
        "tsBuildInfoFile": "./out/test.tsbuildinfo",
        "lib": ["esnext"]
    }
}

tsgo 
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96ma.ts[0m:[93m4[0m:[93m5[0m - [91merror[0m[90m TS2783: [0m'length' is specified more than once, so this usage will be overwritten.

[7m4[0m     length: items.length,
[7m [0m [91m    ~~~~~~~~~~~~~~~~~~~~[0m

  [96ma.ts[0m:[93m8[0m:[93m5[0m - This spread always overwrites this property.
    [7m8[0m     ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'
    [7m [0m [96m    ~~~~~~~~[0m

[96ma.ts[0m:[93m5[0m:[93m13[0m - [91merror[0m[90m TS2339: [0mProperty 'iterator' does not exist on type 'SymbolConstructor'.

[7m5[0m     [Symbol.iterator]: function* () {
[7m [0m [91m            ~~~~~~~~[0m

[96ma.ts[0m:[93m6[0m:[93m28[0m - [91merror[0m[90m TS2488: [0mType 'number[]' must have a '[Symbol.iterator]()' method that returns an iterator.

[7m6[0m         for (const item of items) yield item;
[7m [0m [91m                           ~~~~~[0m


Found 3 errors in the same file, starting at: a.ts[90m:4[0m

//// [/home/src/tslibs/TS/Lib/lib.esnext.d.ts] *Lib*
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
//// [/home/src/workspaces/project/out/a.d.ts] *new* 
declare const items: number[];
declare const obj: {
    [x: number]: number | (() => {});
    length: number;
};

//// [/home/src/workspaces/project/out/a.js] *new* 
"use strict";
const items = [1, 2, 3];
// The spread ...items overwrites the [Symbol.iterator] property specified above it
const obj = {
    length: items.length,
    [Symbol.iterator]: function* () {
        for (const item of items)
            yield item;
    },
    ...items, // This spread overwrites 'length' and '[Symbol.iterator]'
};

//// [/home/src/workspaces/project/out/test.tsbuildinfo] *new* 
{"version":"FakeTSVersion","errors":true,"root":[2],"fileNames":["lib.esnext.d.ts","../a.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"5e4e96b6ea24b37212fe59cffecadaf5-const items = [1, 2, 3];\n// The spread ...items overwrites the [Symbol.iterator] property specified above it\nconst obj = {\n    length: items.length,\n    [Symbol.iterator]: function* () {\n        for (const item of items) yield item;\n    },\n    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'\n};","signature":"ac720d6e8894b43c25e86df6766c3e7f-declare const items: number[];\ndeclare const obj: {\n    [x: number]: number | (() => {});\n    length: number;\n};\n","affectsGlobalScope":true,"impliedNodeFormat":1}],"options":{"declaration":true,"outDir":"./","strict":true,"target":99,"tsBuildInfoFile":"./test.tsbuildinfo"},"semanticDiagnosticsPerFile":[[2,[{"pos":127,"end":147,"code":2783,"category":1,"messageKey":"_0_is_specified_more_than_once_so_this_usage_will_be_overwritten_2783","messageArgs":["length"],"relatedInformation":[{"pos":244,"end":252,"code":2785,"category":1,"messageKey":"This_spread_always_overwrites_this_property_2785"}]},{"pos":161,"end":169,"code":2339,"category":1,"messageKey":"Property_0_does_not_exist_on_type_1_2339","messageArgs":["iterator","SymbolConstructor"]},{"pos":214,"end":219,"code":2488,"category":1,"messageKey":"Type_0_must_have_a_Symbol_iterator_method_that_returns_an_iterator_2488","messageArgs":["number[]"]}]]]}
//// [/home/src/workspaces/project/out/test.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "errors": true,
  "root": [
    {
      "files": [
        "../a.ts"
      ],
      "original": 2
    }
  ],
  "fileNames": [
    "lib.esnext.d.ts",
    "../a.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.esnext.d.ts",
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
      "fileName": "../a.ts",
      "version": "5e4e96b6ea24b37212fe59cffecadaf5-const items = [1, 2, 3];\n// The spread ...items overwrites the [Symbol.iterator] property specified above it\nconst obj = {\n    length: items.length,\n    [Symbol.iterator]: function* () {\n        for (const item of items) yield item;\n    },\n    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'\n};",
      "signature": "ac720d6e8894b43c25e86df6766c3e7f-declare const items: number[];\ndeclare const obj: {\n    [x: number]: number | (() => {});\n    length: number;\n};\n",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "5e4e96b6ea24b37212fe59cffecadaf5-const items = [1, 2, 3];\n// The spread ...items overwrites the [Symbol.iterator] property specified above it\nconst obj = {\n    length: items.length,\n    [Symbol.iterator]: function* () {\n        for (const item of items) yield item;\n    },\n    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'\n};",
        "signature": "ac720d6e8894b43c25e86df6766c3e7f-declare const items: number[];\ndeclare const obj: {\n    [x: number]: number | (() => {});\n    length: number;\n};\n",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "declaration": true,
    "outDir": "./",
    "strict": true,
    "target": 99,
    "tsBuildInfoFile": "./test.tsbuildinfo"
  },
  "semanticDiagnosticsPerFile": [
    [
      "../a.ts",
      [
        {
          "pos": 127,
          "end": 147,
          "code": 2783,
          "category": 1,
          "messageKey": "_0_is_specified_more_than_once_so_this_usage_will_be_overwritten_2783",
          "messageArgs": [
            "length"
          ],
          "relatedInformation": [
            {
              "pos": 244,
              "end": 252,
              "code": 2785,
              "category": 1,
              "messageKey": "This_spread_always_overwrites_this_property_2785"
            }
          ]
        },
        {
          "pos": 161,
          "end": 169,
          "code": 2339,
          "category": 1,
          "messageKey": "Property_0_does_not_exist_on_type_1_2339",
          "messageArgs": [
            "iterator",
            "SymbolConstructor"
          ]
        },
        {
          "pos": 214,
          "end": 219,
          "code": 2488,
          "category": 1,
          "messageKey": "Type_0_must_have_a_Symbol_iterator_method_that_returns_an_iterator_2488",
          "messageArgs": [
            "number[]"
          ]
        }
      ]
    ]
  ],
  "size": 2213
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.esnext.d.ts
*refresh*    /home/src/workspaces/project/a.ts
Signatures::
(stored at emit) /home/src/workspaces/project/a.ts


Edit [0]:: no change

tsgo 
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96ma.ts[0m:[93m4[0m:[93m5[0m - [91merror[0m[90m TS2783: [0m'length' is specified more than once, so this usage will be overwritten.

[7m4[0m     length: items.length,
[7m [0m [91m    ~~~~~~~~~~~~~~~~~~~~[0m

  [96ma.ts[0m:[93m8[0m:[93m5[0m - This spread always overwrites this property.
    [7m8[0m     ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'
    [7m [0m [96m    ~~~~~~~~[0m

[96ma.ts[0m:[93m5[0m:[93m13[0m - [91merror[0m[90m TS2339: [0mProperty 'iterator' does not exist on type 'SymbolConstructor'.

[7m5[0m     [Symbol.iterator]: function* () {
[7m [0m [91m            ~~~~~~~~[0m

[96ma.ts[0m:[93m6[0m:[93m28[0m - [91merror[0m[90m TS2488: [0mType 'number[]' must have a '[Symbol.iterator]()' method that returns an iterator.

[7m6[0m         for (const item of items) yield item;
[7m [0m [91m                           ~~~~~[0m


Found 3 errors in the same file, starting at: a.ts[90m:4[0m

//// [/home/src/workspaces/project/out/test.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.esnext.d.ts","../a.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"5e4e96b6ea24b37212fe59cffecadaf5-const items = [1, 2, 3];\n// The spread ...items overwrites the [Symbol.iterator] property specified above it\nconst obj = {\n    length: items.length,\n    [Symbol.iterator]: function* () {\n        for (const item of items) yield item;\n    },\n    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'\n};","signature":"ac720d6e8894b43c25e86df6766c3e7f-declare const items: number[];\ndeclare const obj: {\n    [x: number]: number | (() => {});\n    length: number;\n};\n","affectsGlobalScope":true,"impliedNodeFormat":1}],"options":{"declaration":true,"outDir":"./","strict":true,"target":99,"tsBuildInfoFile":"./test.tsbuildinfo"},"semanticDiagnosticsPerFile":[[2,[{"pos":127,"end":147,"code":2783,"category":1,"messageKey":"_0_is_specified_more_than_once_so_this_usage_will_be_overwritten_2783","messageArgs":["length"],"relatedInformation":[{"pos":244,"end":252,"code":2785,"category":1,"messageKey":"This_spread_always_overwrites_this_property_2785"}]},{"pos":161,"end":169,"code":2339,"category":1,"messageKey":"Property_0_does_not_exist_on_type_1_2339","messageArgs":["iterator","SymbolConstructor"]},{"pos":214,"end":219,"code":2488,"category":1,"messageKey":"Type_0_must_have_a_Symbol_iterator_method_that_returns_an_iterator_2488","messageArgs":["number[]"]}]]]}
//// [/home/src/workspaces/project/out/test.tsbuildinfo.readable.baseline.txt] *modified* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../a.ts"
      ],
      "original": 2
    }
  ],
  "fileNames": [
    "lib.esnext.d.ts",
    "../a.ts"
  ],
  "fileInfos": [
    {
      "fileName": "lib.esnext.d.ts",
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
      "fileName": "../a.ts",
      "version": "5e4e96b6ea24b37212fe59cffecadaf5-const items = [1, 2, 3];\n// The spread ...items overwrites the [Symbol.iterator] property specified above it\nconst obj = {\n    length: items.length,\n    [Symbol.iterator]: function* () {\n        for (const item of items) yield item;\n    },\n    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'\n};",
      "signature": "ac720d6e8894b43c25e86df6766c3e7f-declare const items: number[];\ndeclare const obj: {\n    [x: number]: number | (() => {});\n    length: number;\n};\n",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "5e4e96b6ea24b37212fe59cffecadaf5-const items = [1, 2, 3];\n// The spread ...items overwrites the [Symbol.iterator] property specified above it\nconst obj = {\n    length: items.length,\n    [Symbol.iterator]: function* () {\n        for (const item of items) yield item;\n    },\n    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'\n};",
        "signature": "ac720d6e8894b43c25e86df6766c3e7f-declare const items: number[];\ndeclare const obj: {\n    [x: number]: number | (() => {});\n    length: number;\n};\n",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "declaration": true,
    "outDir": "./",
    "strict": true,
    "target": 99,
    "tsBuildInfoFile": "./test.tsbuildinfo"
  },
  "semanticDiagnosticsPerFile": [
    [
      "../a.ts",
      [
        {
          "pos": 127,
          "end": 147,
          "code": 2783,
          "category": 1,
          "messageKey": "_0_is_specified_more_than_once_so_this_usage_will_be_overwritten_2783",
          "messageArgs": [
            "length"
          ],
          "relatedInformation": [
            {
              "pos": 244,
              "end": 252,
              "code": 2785,
              "category": 1,
              "messageKey": "This_spread_always_overwrites_this_property_2785"
            }
          ]
        },
        {
          "pos": 161,
          "end": 169,
          "code": 2339,
          "category": 1,
          "messageKey": "Property_0_does_not_exist_on_type_1_2339",
          "messageArgs": [
            "iterator",
            "SymbolConstructor"
          ]
        },
        {
          "pos": 214,
          "end": 219,
          "code": 2488,
          "category": 1,
          "messageKey": "Type_0_must_have_a_Symbol_iterator_method_that_returns_an_iterator_2488",
          "messageArgs": [
            "number[]"
          ]
        }
      ]
    ]
  ],
  "size": 2199
}

tsconfig.json::
SemanticDiagnostics::
Signatures::


Edit [1]:: no change and tsc -b

tsgo -b -v
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[[90mHH:MM:SS AM[0m] Projects in this build: 
    * tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'tsconfig.json' is out of date because buildinfo file 'out/test.tsbuildinfo' indicates that program needs to report errors.

[[90mHH:MM:SS AM[0m] Building project 'tsconfig.json'...

[96ma.ts[0m:[93m4[0m:[93m5[0m - [91merror[0m[90m TS2783: [0m'length' is specified more than once, so this usage will be overwritten.

[7m4[0m     length: items.length,
[7m [0m [91m    ~~~~~~~~~~~~~~~~~~~~[0m

  [96ma.ts[0m:[93m8[0m:[93m5[0m - This spread always overwrites this property.
    [7m8[0m     ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'
    [7m [0m [96m    ~~~~~~~~[0m

[96ma.ts[0m:[93m5[0m:[93m13[0m - [91merror[0m[90m TS2339: [0mProperty 'iterator' does not exist on type 'SymbolConstructor'.

[7m5[0m     [Symbol.iterator]: function* () {
[7m [0m [91m            ~~~~~~~~[0m

[96ma.ts[0m:[93m6[0m:[93m28[0m - [91merror[0m[90m TS2488: [0mType 'number[]' must have a '[Symbol.iterator]()' method that returns an iterator.

[7m6[0m         for (const item of items) yield item;
[7m [0m [91m                           ~~~~~[0m


Found 3 errors in the same file, starting at: a.ts[90m:4[0m


tsconfig.json::
SemanticDiagnostics::
Signatures::

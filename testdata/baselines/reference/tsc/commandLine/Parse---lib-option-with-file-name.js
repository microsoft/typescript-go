
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::--lib es6  first.ts
//// [/home/src/workspaces/project/first.ts] *new* 
export const Key = Symbol()

ExitStatus:: 0

ParsedCommandLine::{
    "parsedConfig": {
        "compilerOptions": {
            "lib": [
                "lib.es2015.d.ts"
            ]
        },
        "watchOptions": {
            "watchInterval": null,
            "watchFile": 0,
            "watchDirectory": 0,
            "fallbackPolling": 0,
            "synchronousWatchDirectory": null,
            "excludeDirectories": null,
            "excludeFiles": null
        },
        "typeAcquisition": null,
        "fileNames": [
            "first.ts"
        ],
        "projectReferences": null
    },
    "configFile": null,
    "errors": [],
    "raw": {
        "lib": [
            "lib.es2015.d.ts"
        ]
    },
    "compileOnSave": null
}
Output::
//// [/home/src/tslibs/TS/Lib/lib.es2015.d.ts] *Lib*
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
//// [/home/src/workspaces/project/first.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Key = void 0;
exports.Key = Symbol();



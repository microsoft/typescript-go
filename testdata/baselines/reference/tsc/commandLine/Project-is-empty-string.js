
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/first.ts] *new* 
export const a = 1
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{ "compilerOptions": { "strict": true, "noEmit": true  } }

ExitStatus:: 0

ParsedCommandLine::{
    "parsedConfig": {
        "compilerOptions": {},
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
        "fileNames": [],
        "projectReferences": null
    },
    "configFile": null,
    "errors": [],
    "raw": {},
    "compileOnSave": null
}
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


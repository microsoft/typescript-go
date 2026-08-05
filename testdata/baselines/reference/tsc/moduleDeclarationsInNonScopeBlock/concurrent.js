currentDirectory::/home/src/projects/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/projects/project/a.ts] *new* 
{
    export { a } from "exportNamed";
    export * from "exportStar";
    import { b } from "importNamed";
    import c = require("importEquals");
    import "sideEffect";
}

tsgo --noEmit a.ts
ExitStatus:: DiagnosticsPresent_OutputsSkipped
Output::
[96ma.ts[0m:[93m2[0m:[93m5[0m - [91merror[0m[90m TS1233: [0mAn export declaration can only be used at the top level of a namespace or module.

[7m2[0m     export { a } from "exportNamed";
[7m [0m [91m    ~~~~~~[0m

[96ma.ts[0m:[93m2[0m:[93m23[0m - [91merror[0m[90m TS2307: [0mCannot find module 'exportNamed' or its corresponding type declarations.

[7m2[0m     export { a } from "exportNamed";
[7m [0m [91m                      ~~~~~~~~~~~~~[0m

[96ma.ts[0m:[93m3[0m:[93m5[0m - [91merror[0m[90m TS1233: [0mAn export declaration can only be used at the top level of a namespace or module.

[7m3[0m     export * from "exportStar";
[7m [0m [91m    ~~~~~~[0m

[96ma.ts[0m:[93m3[0m:[93m19[0m - [91merror[0m[90m TS2307: [0mCannot find module 'exportStar' or its corresponding type declarations.

[7m3[0m     export * from "exportStar";
[7m [0m [91m                  ~~~~~~~~~~~~[0m

[96ma.ts[0m:[93m4[0m:[93m5[0m - [91merror[0m[90m TS1232: [0mAn import declaration can only be used at the top level of a namespace or module.

[7m4[0m     import { b } from "importNamed";
[7m [0m [91m    ~~~~~~[0m

[96ma.ts[0m:[93m4[0m:[93m23[0m - [91merror[0m[90m TS2307: [0mCannot find module 'importNamed' or its corresponding type declarations.

[7m4[0m     import { b } from "importNamed";
[7m [0m [91m                      ~~~~~~~~~~~~~[0m

[96ma.ts[0m:[93m5[0m:[93m5[0m - [91merror[0m[90m TS1232: [0mAn import declaration can only be used at the top level of a namespace or module.

[7m5[0m     import c = require("importEquals");
[7m [0m [91m    ~~~~~~[0m

[96ma.ts[0m:[93m5[0m:[93m24[0m - [91merror[0m[90m TS2307: [0mCannot find module 'importEquals' or its corresponding type declarations.

[7m5[0m     import c = require("importEquals");
[7m [0m [91m                       ~~~~~~~~~~~~~~[0m

[96ma.ts[0m:[93m6[0m:[93m5[0m - [91merror[0m[90m TS1232: [0mAn import declaration can only be used at the top level of a namespace or module.

[7m6[0m     import "sideEffect";
[7m [0m [91m    ~~~~~~[0m


Found 9 errors in the same file, starting at: a.ts[90m:2[0m

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


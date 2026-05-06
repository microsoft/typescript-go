//// [tests/cases/compiler/declarationEmitModuleWithScopeMarker.ts] ////

//// [declarationEmitModuleWithScopeMarker.ts]
declare module "bar" {
    var before: typeof func;

    export function normal(): void;

    export default function func(): typeof func;

    var after: typeof func;

    export {}
}


//// [declarationEmitModuleWithScopeMarker.js]
"use strict";


//// [declarationEmitModuleWithScopeMarker.d.ts]
module "bar" {
    var before: typeof func;
    export function normal(): void;
    export default function func(): typeof func;
    var after: typeof func;
    export {};
}


//// [DtsFileErrors]


declarationEmitModuleWithScopeMarker.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitModuleWithScopeMarker.d.ts (1 errors) ====
    module "bar" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var before: typeof func;
        export function normal(): void;
        export default function func(): typeof func;
        var after: typeof func;
        export {};
    }
    
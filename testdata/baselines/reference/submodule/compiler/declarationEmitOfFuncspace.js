//// [tests/cases/compiler/declarationEmitOfFuncspace.ts] ////

//// [expando.ts]
// #27032
function ExpandoMerge(n: number) {
    return n;
}
namespace ExpandoMerge {
    export interface I { }
}


//// [expando.js]
"use strict";
// #27032
function ExpandoMerge(n) {
    return n;
}


//// [expando.d.ts]
function ExpandoMerge(n: number): number;
namespace ExpandoMerge {
    interface I {
    }
}


//// [DtsFileErrors]


expando.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== expando.d.ts (1 errors) ====
    function ExpandoMerge(n: number): number;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    namespace ExpandoMerge {
        interface I {
        }
    }
    
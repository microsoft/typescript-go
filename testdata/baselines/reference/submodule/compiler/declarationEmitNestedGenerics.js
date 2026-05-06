//// [tests/cases/compiler/declarationEmitNestedGenerics.ts] ////

//// [declarationEmitNestedGenerics.ts]
function f<T>(p: T) {
    let g: <T>(x: T) => typeof p = null as any;
    return g;
}

function g<T>(x: T) {
    let y: typeof x extends (infer T)[] ? T : typeof x = null as any;
    return y;
}

//// [declarationEmitNestedGenerics.js]
"use strict";
function f(p) {
    let g = null;
    return g;
}
function g(x) {
    let y = null;
    return y;
}


//// [declarationEmitNestedGenerics.d.ts]
function f<T>(p: T): <T_1>(x: T_1) => typeof p;
function g<T>(x: T): T extends (infer T_1)[] ? T_1 : T;


//// [DtsFileErrors]


declarationEmitNestedGenerics.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitNestedGenerics.d.ts (1 errors) ====
    function f<T>(p: T): <T_1>(x: T_1) => typeof p;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function g<T>(x: T): T extends (infer T_1)[] ? T_1 : T;
    
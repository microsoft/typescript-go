//// [tests/cases/compiler/declarationEmitFunctionDuplicateNamespace.ts] ////

//// [declarationEmitFunctionDuplicateNamespace.ts]
function f(a: 0): 0;
function f(a: 1): 1;
function f(a: 0 | 1) {
    return a;
}

f.x = 2;


//// [declarationEmitFunctionDuplicateNamespace.js]
"use strict";
function f(a) {
    return a;
}
f.x = 2;


//// [declarationEmitFunctionDuplicateNamespace.d.ts]
function f(a: 0): 0;
function f(a: 1): 1;
declare namespace f {
    var x: number;
}


//// [DtsFileErrors]


declarationEmitFunctionDuplicateNamespace.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitFunctionDuplicateNamespace.d.ts (1 errors) ====
    function f(a: 0): 0;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f(a: 1): 1;
    declare namespace f {
        var x: number;
    }
    
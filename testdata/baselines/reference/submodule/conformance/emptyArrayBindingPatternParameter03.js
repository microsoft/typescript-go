//// [tests/cases/conformance/es6/destructuring/emptyArrayBindingPatternParameter03.ts] ////

//// [emptyArrayBindingPatternParameter03.ts]
function f(a, []) {
    var x, y, z;
}

//// [emptyArrayBindingPatternParameter03.js]
"use strict";
function f(a, []) {
    var x, y, z;
}


//// [emptyArrayBindingPatternParameter03.d.ts]
function f(a: any, []: Iterable<any, void, undefined>): void;


//// [DtsFileErrors]


emptyArrayBindingPatternParameter03.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyArrayBindingPatternParameter03.d.ts (1 errors) ====
    function f(a: any, []: Iterable<any, void, undefined>): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
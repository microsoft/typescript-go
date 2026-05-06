//// [tests/cases/conformance/es6/destructuring/emptyObjectBindingPatternParameter02.ts] ////

//// [emptyObjectBindingPatternParameter02.ts]
function f(a, {}) {
    var x, y, z;
}

//// [emptyObjectBindingPatternParameter02.js]
"use strict";
function f(a, {}) {
    var x, y, z;
}


//// [emptyObjectBindingPatternParameter02.d.ts]
function f(a: any, {}: {}): void;


//// [DtsFileErrors]


emptyObjectBindingPatternParameter02.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyObjectBindingPatternParameter02.d.ts (1 errors) ====
    function f(a: any, {}: {}): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
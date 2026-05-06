//// [tests/cases/conformance/es6/destructuring/emptyObjectBindingPatternParameter01.ts] ////

//// [emptyObjectBindingPatternParameter01.ts]
function f({}) {
    var x, y, z;
}

//// [emptyObjectBindingPatternParameter01.js]
"use strict";
function f({}) {
    var x, y, z;
}


//// [emptyObjectBindingPatternParameter01.d.ts]
function f({}: {}): void;


//// [DtsFileErrors]


emptyObjectBindingPatternParameter01.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyObjectBindingPatternParameter01.d.ts (1 errors) ====
    function f({}: {}): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
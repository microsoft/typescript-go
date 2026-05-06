//// [tests/cases/conformance/es6/destructuring/emptyObjectBindingPatternParameter04.ts] ////

//// [emptyObjectBindingPatternParameter04.ts]
function f({} = {a: 1, b: "2", c: true}) {
    var x, y, z;
}

//// [emptyObjectBindingPatternParameter04.js]
"use strict";
function f({} = { a: 1, b: "2", c: true }) {
    var x, y, z;
}


//// [emptyObjectBindingPatternParameter04.d.ts]
function f({}?: {
    a: number;
    b: string;
    c: boolean;
}): void;


//// [DtsFileErrors]


emptyObjectBindingPatternParameter04.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyObjectBindingPatternParameter04.d.ts (1 errors) ====
    function f({}?: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a: number;
        b: string;
        c: boolean;
    }): void;
    
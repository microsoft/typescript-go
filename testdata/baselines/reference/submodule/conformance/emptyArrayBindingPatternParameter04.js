//// [tests/cases/conformance/es6/destructuring/emptyArrayBindingPatternParameter04.ts] ////

//// [emptyArrayBindingPatternParameter04.ts]
function f([] = [1,2,3,4]) {
    var x, y, z;
}

//// [emptyArrayBindingPatternParameter04.js]
"use strict";
function f([] = [1, 2, 3, 4]) {
    var x, y, z;
}


//// [emptyArrayBindingPatternParameter04.d.ts]
function f([]?: number[]): void;


//// [DtsFileErrors]


emptyArrayBindingPatternParameter04.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyArrayBindingPatternParameter04.d.ts (1 errors) ====
    function f([]?: number[]): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
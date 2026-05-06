//// [tests/cases/conformance/es6/destructuring/emptyArrayBindingPatternParameter01.ts] ////

//// [emptyArrayBindingPatternParameter01.ts]
function f([]) {
    var x, y, z;
}

//// [emptyArrayBindingPatternParameter01.js]
"use strict";
function f([]) {
    var x, y, z;
}


//// [emptyArrayBindingPatternParameter01.d.ts]
function f([]: Iterable<any, void, undefined>): void;


//// [DtsFileErrors]


emptyArrayBindingPatternParameter01.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyArrayBindingPatternParameter01.d.ts (1 errors) ====
    function f([]: Iterable<any, void, undefined>): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
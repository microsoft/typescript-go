//// [tests/cases/compiler/functionReturningItself.ts] ////

//// [functionReturningItself.ts]
function somefn() {
    return somefn;
}

//// [functionReturningItself.js]
"use strict";
function somefn() {
    return somefn;
}


//// [functionReturningItself.d.ts]
function somefn(): typeof somefn;


//// [DtsFileErrors]


functionReturningItself.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== functionReturningItself.d.ts (1 errors) ====
    function somefn(): typeof somefn;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
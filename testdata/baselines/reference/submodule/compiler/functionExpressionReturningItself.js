//// [tests/cases/compiler/functionExpressionReturningItself.ts] ////

//// [functionExpressionReturningItself.ts]
var x = function somefn() { return somefn; };

//// [functionExpressionReturningItself.js]
"use strict";
var x = function somefn() { return somefn; };


//// [functionExpressionReturningItself.d.ts]
var x: () => () => /*elided*/ any;


//// [DtsFileErrors]


functionExpressionReturningItself.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== functionExpressionReturningItself.d.ts (1 errors) ====
    var x: () => () => /*elided*/ any;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
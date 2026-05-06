//// [tests/cases/compiler/genericFunctions1.ts] ////

//// [genericFunctions1.ts]
function foo<T > (x: T) { return x; }

var x = foo(5); // 'x' should be number

//// [genericFunctions1.js]
"use strict";
function foo(x) { return x; }
var x = foo(5); // 'x' should be number


//// [genericFunctions1.d.ts]
function foo<T>(x: T): T;
var x: number;


//// [DtsFileErrors]


genericFunctions1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== genericFunctions1.d.ts (1 errors) ====
    function foo<T>(x: T): T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var x: number;
    
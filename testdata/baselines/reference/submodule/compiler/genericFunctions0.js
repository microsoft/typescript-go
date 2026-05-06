//// [tests/cases/compiler/genericFunctions0.ts] ////

//// [genericFunctions0.ts]
function foo<T > (x: T) { return x; }

var x = foo<number>(5); // 'x' should be number

//// [genericFunctions0.js]
"use strict";
function foo(x) { return x; }
var x = foo(5); // 'x' should be number


//// [genericFunctions0.d.ts]
function foo<T>(x: T): T;
var x: number;


//// [DtsFileErrors]


genericFunctions0.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== genericFunctions0.d.ts (1 errors) ====
    function foo<T>(x: T): T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var x: number;
    
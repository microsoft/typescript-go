//// [tests/cases/compiler/selfReferentialFunctionType.ts] ////

//// [selfReferentialFunctionType.ts]
declare function f<T>(args: typeof f<T>): T;
declare function g<T = typeof g>(args: T): T;
declare function h<T>(): typeof h<T>;


//// [selfReferentialFunctionType.js]
"use strict";


//// [selfReferentialFunctionType.d.ts]
function f<T>(args: typeof f<T>): T;
function g<T = typeof g>(args: T): T;
function h<T>(): typeof h<T>;


//// [DtsFileErrors]


selfReferentialFunctionType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== selfReferentialFunctionType.d.ts (1 errors) ====
    function f<T>(args: typeof f<T>): T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function g<T = typeof g>(args: T): T;
    function h<T>(): typeof h<T>;
    
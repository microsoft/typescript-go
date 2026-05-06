//// [tests/cases/compiler/declarationEmitDestructuringArrayPattern5.ts] ////

//// [declarationEmitDestructuringArrayPattern5.ts]
var [, , z] = [1, 2, 4];
var [, a, , ] = [3, 4, 5];
var [, , [, b, ]] = [3,5,[0, 1]];

//// [declarationEmitDestructuringArrayPattern5.js]
"use strict";
var [, , z] = [1, 2, 4];
var [, a, ,] = [3, 4, 5];
var [, , [, b,]] = [3, 5, [0, 1]];


//// [declarationEmitDestructuringArrayPattern5.d.ts]
var z: number;
var a: number;
var b: number;


//// [DtsFileErrors]


declarationEmitDestructuringArrayPattern5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuringArrayPattern5.d.ts (1 errors) ====
    var z: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var a: number;
    var b: number;
    
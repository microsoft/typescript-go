//// [tests/cases/compiler/declFileForVarList.ts] ////

//// [declFileForVarList.ts]
var x, y, z = 1;
var x1 = 1, y2 = 2, z2 = 3;

//// [declFileForVarList.js]
"use strict";
var x, y, z = 1;
var x1 = 1, y2 = 2, z2 = 3;


//// [declFileForVarList.d.ts]
var x: any, y: any, z: number;
var x1: number, y2: number, z2: number;


//// [DtsFileErrors]


declFileForVarList.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileForVarList.d.ts (1 errors) ====
    var x: any, y: any, z: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var x1: number, y2: number, z2: number;
    
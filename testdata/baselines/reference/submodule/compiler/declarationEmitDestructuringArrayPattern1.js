//// [tests/cases/compiler/declarationEmitDestructuringArrayPattern1.ts] ////

//// [declarationEmitDestructuringArrayPattern1.ts]
var [] = [1, "hello"]; // Dont emit anything
var [x] = [1, "hello"]; // emit x: number
var [x1, y1] = [1, "hello"]; // emit x1: number, y1: string
var [, , z1] = [0, 1, 2]; // emit z1: number

var a = [1, "hello"];
var [x2] = a;          // emit x2: number | string
var [x3, y3, z3] = a;  // emit x3, y3, z3 

//// [declarationEmitDestructuringArrayPattern1.js]
"use strict";
var [] = [1, "hello"]; // Dont emit anything
var [x] = [1, "hello"]; // emit x: number
var [x1, y1] = [1, "hello"]; // emit x1: number, y1: string
var [, , z1] = [0, 1, 2]; // emit z1: number
var a = [1, "hello"];
var [x2] = a; // emit x2: number | string
var [x3, y3, z3] = a; // emit x3, y3, z3 


//// [declarationEmitDestructuringArrayPattern1.d.ts]
var x: number;
var x1: number, y1: string;
var z1: number;
var a: (string | number)[];
var x2: string | number;
var x3: string | number, y3: string | number, z3: string | number;


//// [DtsFileErrors]


declarationEmitDestructuringArrayPattern1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuringArrayPattern1.d.ts (1 errors) ====
    var x: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var x1: number, y1: string;
    var z1: number;
    var a: (string | number)[];
    var x2: string | number;
    var x3: string | number, y3: string | number, z3: string | number;
    
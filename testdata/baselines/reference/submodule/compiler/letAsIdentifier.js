//// [tests/cases/compiler/letAsIdentifier.ts] ////

//// [letAsIdentifier.ts]
var let = 10;
var a = 10;
let = 30;
let
a;

//// [letAsIdentifier.js]
"use strict";
var let = 10;
var a = 10;
let = 30;
let a;


//// [letAsIdentifier.d.ts]
var let: number;
var a: number;
let a: any;

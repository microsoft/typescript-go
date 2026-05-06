//// [tests/cases/compiler/genericArray0.ts] ////

//// [genericArray0.ts]
var x:number[];


var y = x; 

function map<U>() {
    var ys: U[] = [];
}


//// [genericArray0.js]
"use strict";
var x;
var y = x;
function map() {
    var ys = [];
}


//// [genericArray0.d.ts]
var x: number[];
var y: number[];
function map<U>(): void;

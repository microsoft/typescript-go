//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesWithVariousOperators02.ts] ////

//// [stringLiteralTypesWithVariousOperators02.ts]
declare let abc: "ABC";
declare let xyz: "XYZ";
declare let abcOrXyz: "ABC" | "XYZ";
declare let abcOrXyzOrNumber: "ABC" | "XYZ" | number;

let a = abcOrXyzOrNumber + 100;
let b = 100 + abcOrXyzOrNumber;
let c = abcOrXyzOrNumber + abcOrXyzOrNumber;
let d = abcOrXyzOrNumber + true;
let e = false + abcOrXyzOrNumber;
let f = abcOrXyzOrNumber++;
let g = --abcOrXyzOrNumber;
let h = abcOrXyzOrNumber ^ 10;
let i = abcOrXyzOrNumber | 10;
let j = abc < xyz;
let k = abc === xyz;
let l = abc != xyz;

//// [stringLiteralTypesWithVariousOperators02.js]
"use strict";
let a = abcOrXyzOrNumber + 100;
let b = 100 + abcOrXyzOrNumber;
let c = abcOrXyzOrNumber + abcOrXyzOrNumber;
let d = abcOrXyzOrNumber + true;
let e = false + abcOrXyzOrNumber;
let f = abcOrXyzOrNumber++;
let g = --abcOrXyzOrNumber;
let h = abcOrXyzOrNumber ^ 10;
let i = abcOrXyzOrNumber | 10;
let j = abc < xyz;
let k = abc === xyz;
let l = abc != xyz;


//// [stringLiteralTypesWithVariousOperators02.d.ts]
let abc: "ABC";
let xyz: "XYZ";
let abcOrXyz: "ABC" | "XYZ";
let abcOrXyzOrNumber: "ABC" | "XYZ" | number;
let a: any;
let b: any;
let c: any;
let d: any;
let e: any;
let f: number;
let g: number;
let h: number;
let i: number;
let j: boolean;
let k: boolean;
let l: boolean;

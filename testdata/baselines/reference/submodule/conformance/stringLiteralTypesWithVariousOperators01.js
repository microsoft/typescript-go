//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesWithVariousOperators01.ts] ////

//// [stringLiteralTypesWithVariousOperators01.ts]
declare let abc: "ABC";
declare let xyz: "XYZ";
declare let abcOrXyz: "ABC" | "XYZ";
declare let abcOrXyzOrNumber: "ABC" | "XYZ" | number;

let a = "" + abc;
let b = abc + "";
let c = 10 + abc;
let d = abc + 10;
let e = xyz + abc;
let f = abc + xyz;
let g = true + abc;
let h = abc + true;
let i = abc + abcOrXyz + xyz;
let j = abcOrXyz + abcOrXyz;
let k = +abcOrXyz;
let l = -abcOrXyz;
let m = abcOrXyzOrNumber + "";
let n = "" + abcOrXyzOrNumber;
let o = abcOrXyzOrNumber + abcOrXyz;
let p = abcOrXyz + abcOrXyzOrNumber;
let q = !abcOrXyzOrNumber;
let r = ~abcOrXyzOrNumber;
let s = abcOrXyzOrNumber < abcOrXyzOrNumber;
let t = abcOrXyzOrNumber >= abcOrXyz;
let u = abc === abcOrXyz;
let v = abcOrXyz === abcOrXyzOrNumber;

//// [stringLiteralTypesWithVariousOperators01.js]
"use strict";
let a = "" + abc;
let b = abc + "";
let c = 10 + abc;
let d = abc + 10;
let e = xyz + abc;
let f = abc + xyz;
let g = true + abc;
let h = abc + true;
let i = abc + abcOrXyz + xyz;
let j = abcOrXyz + abcOrXyz;
let k = +abcOrXyz;
let l = -abcOrXyz;
let m = abcOrXyzOrNumber + "";
let n = "" + abcOrXyzOrNumber;
let o = abcOrXyzOrNumber + abcOrXyz;
let p = abcOrXyz + abcOrXyzOrNumber;
let q = !abcOrXyzOrNumber;
let r = ~abcOrXyzOrNumber;
let s = abcOrXyzOrNumber < abcOrXyzOrNumber;
let t = abcOrXyzOrNumber >= abcOrXyz;
let u = abc === abcOrXyz;
let v = abcOrXyz === abcOrXyzOrNumber;


//// [stringLiteralTypesWithVariousOperators01.d.ts]
let abc: "ABC";
let xyz: "XYZ";
let abcOrXyz: "ABC" | "XYZ";
let abcOrXyzOrNumber: "ABC" | "XYZ" | number;
let a: string;
let b: string;
let c: string;
let d: string;
let e: string;
let f: string;
let g: string;
let h: string;
let i: string;
let j: string;
let k: number;
let l: number;
let m: string;
let n: string;
let o: string;
let p: string;
let q: boolean;
let r: number;
let s: boolean;
let t: boolean;
let u: boolean;
let v: boolean;


//// [DtsFileErrors]


stringLiteralTypesWithVariousOperators01.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesWithVariousOperators01.d.ts (1 errors) ====
    let abc: "ABC";
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    let xyz: "XYZ";
    let abcOrXyz: "ABC" | "XYZ";
    let abcOrXyzOrNumber: "ABC" | "XYZ" | number;
    let a: string;
    let b: string;
    let c: string;
    let d: string;
    let e: string;
    let f: string;
    let g: string;
    let h: string;
    let i: string;
    let j: string;
    let k: number;
    let l: number;
    let m: string;
    let n: string;
    let o: string;
    let p: string;
    let q: boolean;
    let r: number;
    let s: boolean;
    let t: boolean;
    let u: boolean;
    let v: boolean;
    
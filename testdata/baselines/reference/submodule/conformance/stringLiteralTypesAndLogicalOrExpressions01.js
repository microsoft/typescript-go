//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesAndLogicalOrExpressions01.ts] ////

//// [stringLiteralTypesAndLogicalOrExpressions01.ts]
declare function myRandBool(): boolean;

let a: "foo" = "foo";
let b = a || "foo";
let c: "foo" = b;
let d = b || "bar";
let e: "foo" | "bar" = d;


//// [stringLiteralTypesAndLogicalOrExpressions01.js]
"use strict";
let a = "foo";
let b = a || "foo";
let c = b;
let d = b || "bar";
let e = d;


//// [stringLiteralTypesAndLogicalOrExpressions01.d.ts]
function myRandBool(): boolean;
let a: "foo";
let b: "foo";
let c: "foo";
let d: "foo";
let e: "foo" | "bar";


//// [DtsFileErrors]


stringLiteralTypesAndLogicalOrExpressions01.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesAndLogicalOrExpressions01.d.ts (1 errors) ====
    function myRandBool(): boolean;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    let a: "foo";
    let b: "foo";
    let c: "foo";
    let d: "foo";
    let e: "foo" | "bar";
    
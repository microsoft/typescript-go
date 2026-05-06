//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesAndParenthesizedExpressions01.ts] ////

//// [stringLiteralTypesAndParenthesizedExpressions01.ts]
declare function myRandBool(): boolean;

let a: "foo" = ("foo");
let b: "foo" | "bar" = ("foo");
let c: "foo" = (myRandBool ? "foo" : ("foo"));
let d: "foo" | "bar" = (myRandBool ? "foo" : ("bar"));


//// [stringLiteralTypesAndParenthesizedExpressions01.js]
"use strict";
let a = ("foo");
let b = ("foo");
let c = (myRandBool ? "foo" : ("foo"));
let d = (myRandBool ? "foo" : ("bar"));


//// [stringLiteralTypesAndParenthesizedExpressions01.d.ts]
function myRandBool(): boolean;
let a: "foo";
let b: "foo" | "bar";
let c: "foo";
let d: "foo" | "bar";

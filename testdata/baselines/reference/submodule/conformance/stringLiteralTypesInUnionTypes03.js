//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesInUnionTypes03.ts] ////

//// [stringLiteralTypesInUnionTypes03.ts]
type T = number | "foo" | "bar";

var x: "foo" | "bar" | number;
var y: T = undefined;

if (x === "foo") {
    let a = x;
}
else if (x !== "bar") {
    let b = x || y;
}
else {
    let c = x;
    let d = y;
    let e: (typeof x) | (typeof y) = c || d;
}

x = y;
y = x;

//// [stringLiteralTypesInUnionTypes03.js]
"use strict";
var x;
var y = undefined;
if (x === "foo") {
    let a = x;
}
else if (x !== "bar") {
    let b = x || y;
}
else {
    let c = x;
    let d = y;
    let e = c || d;
}
x = y;
y = x;


//// [stringLiteralTypesInUnionTypes03.d.ts]
type T = number | "foo" | "bar";
var x: "foo" | "bar" | number;
var y: T;


//// [DtsFileErrors]


stringLiteralTypesInUnionTypes03.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesInUnionTypes03.d.ts (1 errors) ====
    type T = number | "foo" | "bar";
    var x: "foo" | "bar" | number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var y: T;
    
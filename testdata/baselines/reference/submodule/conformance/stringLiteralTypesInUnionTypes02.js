//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesInUnionTypes02.ts] ////

//// [stringLiteralTypesInUnionTypes02.ts]
type T = string | "foo" | "bar" | "baz";

var x: "foo" | "bar" | "baz" | string = undefined;
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

//// [stringLiteralTypesInUnionTypes02.js]
"use strict";
var x = undefined;
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


//// [stringLiteralTypesInUnionTypes02.d.ts]
type T = string | "foo" | "bar" | "baz";
var x: "foo" | "bar" | "baz" | string;
var y: T;


//// [DtsFileErrors]


stringLiteralTypesInUnionTypes02.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesInUnionTypes02.d.ts (1 errors) ====
    type T = string | "foo" | "bar" | "baz";
    var x: "foo" | "bar" | "baz" | string;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var y: T;
    
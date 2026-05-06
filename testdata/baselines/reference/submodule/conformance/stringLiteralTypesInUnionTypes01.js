//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesInUnionTypes01.ts] ////

//// [stringLiteralTypesInUnionTypes01.ts]
type T = "foo" | "bar" | "baz";

var x: "foo" | "bar" | "baz" = undefined;
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

//// [stringLiteralTypesInUnionTypes01.js]
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


//// [stringLiteralTypesInUnionTypes01.d.ts]
type T = "foo" | "bar" | "baz";
var x: "foo" | "bar" | "baz";
var y: T;


//// [DtsFileErrors]


stringLiteralTypesInUnionTypes01.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesInUnionTypes01.d.ts (1 errors) ====
    type T = "foo" | "bar" | "baz";
    var x: "foo" | "bar" | "baz";
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var y: T;
    
//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesAsTypeParameterConstraint02.ts] ////

//// [stringLiteralTypesAsTypeParameterConstraint02.ts]
function foo<T extends "foo">(f: (x: T) => T) {
    return f;
}

let f = foo((y: "foo" | "bar") => y === "foo" ? y : "foo");
let fResult = f("foo");

//// [stringLiteralTypesAsTypeParameterConstraint02.js]
"use strict";
function foo(f) {
    return f;
}
let f = foo((y) => y === "foo" ? y : "foo");
let fResult = f("foo");


//// [stringLiteralTypesAsTypeParameterConstraint02.d.ts]
function foo<T extends "foo">(f: (x: T) => T): (x: T) => T;
let f: (x: "foo") => "foo";
let fResult: "foo";


//// [DtsFileErrors]


stringLiteralTypesAsTypeParameterConstraint02.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesAsTypeParameterConstraint02.d.ts (1 errors) ====
    function foo<T extends "foo">(f: (x: T) => T): (x: T) => T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    let f: (x: "foo") => "foo";
    let fResult: "foo";
    
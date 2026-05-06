//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesAsTypeParameterConstraint01.ts] ////

//// [stringLiteralTypesAsTypeParameterConstraint01.ts]
function foo<T extends "foo">(f: (x: T) => T) {
    return f;
}

function bar<T extends "foo" | "bar">(f: (x: T) => T) {
    return f;
}

let f = foo(x => x);
let fResult = f("foo");

let g = foo((x => x));
let gResult = g("foo");

let h = bar(x => x);
let hResult = h("foo");
hResult = h("bar");

//// [stringLiteralTypesAsTypeParameterConstraint01.js]
"use strict";
function foo(f) {
    return f;
}
function bar(f) {
    return f;
}
let f = foo(x => x);
let fResult = f("foo");
let g = foo((x => x));
let gResult = g("foo");
let h = bar(x => x);
let hResult = h("foo");
hResult = h("bar");


//// [stringLiteralTypesAsTypeParameterConstraint01.d.ts]
function foo<T extends "foo">(f: (x: T) => T): (x: T) => T;
function bar<T extends "foo" | "bar">(f: (x: T) => T): (x: T) => T;
let f: (x: "foo") => "foo";
let fResult: "foo";
let g: (x: "foo") => "foo";
let gResult: "foo";
let h: (x: "bar" | "foo") => "bar" | "foo";
let hResult: "bar" | "foo";


//// [DtsFileErrors]


stringLiteralTypesAsTypeParameterConstraint01.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesAsTypeParameterConstraint01.d.ts (1 errors) ====
    function foo<T extends "foo">(f: (x: T) => T): (x: T) => T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function bar<T extends "foo" | "bar">(f: (x: T) => T): (x: T) => T;
    let f: (x: "foo") => "foo";
    let fResult: "foo";
    let g: (x: "foo") => "foo";
    let gResult: "foo";
    let h: (x: "bar" | "foo") => "bar" | "foo";
    let hResult: "bar" | "foo";
    
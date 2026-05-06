//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesOverloadAssignability03.ts] ////

//// [stringLiteralTypesOverloadAssignability03.ts]
function f(x: "foo"): number;
function f(x: string): number {
    return 0;
}

function g(x: "foo"): number;
function g(x: string): number {
    return 0;
}

let a = f;
let b = g;

a = b;
b = a;

//// [stringLiteralTypesOverloadAssignability03.js]
"use strict";
function f(x) {
    return 0;
}
function g(x) {
    return 0;
}
let a = f;
let b = g;
a = b;
b = a;


//// [stringLiteralTypesOverloadAssignability03.d.ts]
function f(x: "foo"): number;
function g(x: "foo"): number;
let a: typeof f;
let b: typeof g;


//// [DtsFileErrors]


stringLiteralTypesOverloadAssignability03.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesOverloadAssignability03.d.ts (1 errors) ====
    function f(x: "foo"): number;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function g(x: "foo"): number;
    let a: typeof f;
    let b: typeof g;
    
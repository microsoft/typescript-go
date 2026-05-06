//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesOverloads04.ts] ////

//// [stringLiteralTypesOverloads04.ts]
declare function f(x: (p: "foo" | "bar") => "foo");

f(y => {
    const z = y = "foo";
    return z;
})

//// [stringLiteralTypesOverloads04.js]
"use strict";
f(y => {
    const z = y = "foo";
    return z;
});


//// [stringLiteralTypesOverloads04.d.ts]
function f(x: (p: "foo" | "bar") => "foo"): any;


//// [DtsFileErrors]


stringLiteralTypesOverloads04.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesOverloads04.d.ts (1 errors) ====
    function f(x: (p: "foo" | "bar") => "foo"): any;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
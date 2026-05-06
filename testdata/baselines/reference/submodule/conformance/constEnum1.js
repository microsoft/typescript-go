//// [tests/cases/conformance/constEnums/constEnum1.ts] ////

//// [constEnum1.ts]
// An enum declaration that specifies a const modifier is a constant enum declaration.
// In a constant enum declaration, all members must have constant values and
// it is an error for a member declaration to specify an expression that isn't classified as a constant enum expression.

const enum E {
    a = 10,
    b = a,
    c = (a+1),
    e,
    d = ~e,
    f = a << 2 >> 1,
    g = a << 2 >>> 1,
    h = a | b
}

//// [constEnum1.js]
"use strict";
// An enum declaration that specifies a const modifier is a constant enum declaration.
// In a constant enum declaration, all members must have constant values and
// it is an error for a member declaration to specify an expression that isn't classified as a constant enum expression.


//// [constEnum1.d.ts]
const enum E {
    a = 10,
    b = 10,
    c = 11,
    e = 12,
    d = -13,
    f = 20,
    g = 20,
    h = 10
}


//// [DtsFileErrors]


constEnum1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== constEnum1.d.ts (1 errors) ====
    const enum E {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a = 10,
        b = 10,
        c = 11,
        e = 12,
        d = -13,
        f = 20,
        g = 20,
        h = 10
    }
    
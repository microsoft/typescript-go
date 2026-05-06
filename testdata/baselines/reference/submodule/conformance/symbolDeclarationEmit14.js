//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit14.ts] ////

//// [symbolDeclarationEmit14.ts]
class C {
    get [Symbol.toPrimitive]() { return ""; }
    get [Symbol.toStringTag]() { return ""; }
}

//// [symbolDeclarationEmit14.js]
"use strict";
class C {
    get [Symbol.toPrimitive]() { return ""; }
    get [Symbol.toStringTag]() { return ""; }
}


//// [symbolDeclarationEmit14.d.ts]
class C {
    get [Symbol.toPrimitive](): string;
    get [Symbol.toStringTag](): string;
}


//// [DtsFileErrors]


symbolDeclarationEmit14.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit14.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        get [Symbol.toPrimitive](): string;
        get [Symbol.toStringTag](): string;
    }
    
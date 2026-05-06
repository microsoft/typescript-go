//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit4.ts] ////

//// [symbolDeclarationEmit4.ts]
class C {
    get [Symbol.toPrimitive]() { return ""; }
    set [Symbol.toPrimitive](x) { }
}

//// [symbolDeclarationEmit4.js]
"use strict";
class C {
    get [Symbol.toPrimitive]() { return ""; }
    set [Symbol.toPrimitive](x) { }
}


//// [symbolDeclarationEmit4.d.ts]
class C {
    get [Symbol.toPrimitive](): string;
    set [Symbol.toPrimitive](x: string);
}


//// [DtsFileErrors]


symbolDeclarationEmit4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit4.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        get [Symbol.toPrimitive](): string;
        set [Symbol.toPrimitive](x: string);
    }
    
//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit13.ts] ////

//// [symbolDeclarationEmit13.ts]
class C {
    get [Symbol.toPrimitive]() { return ""; }
    set [Symbol.toStringTag](x) { }
}

//// [symbolDeclarationEmit13.js]
"use strict";
class C {
    get [Symbol.toPrimitive]() { return ""; }
    set [Symbol.toStringTag](x) { }
}


//// [symbolDeclarationEmit13.d.ts]
class C {
    get [Symbol.toPrimitive](): string;
    set [Symbol.toStringTag](x: any);
}


//// [DtsFileErrors]


symbolDeclarationEmit13.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit13.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        get [Symbol.toPrimitive](): string;
        set [Symbol.toStringTag](x: any);
    }
    
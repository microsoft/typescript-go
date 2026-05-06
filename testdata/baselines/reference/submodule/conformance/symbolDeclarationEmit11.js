//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit11.ts] ////

//// [symbolDeclarationEmit11.ts]
class C {
    static [Symbol.iterator] = 0;
    static [Symbol.isConcatSpreadable]() { }
    static get [Symbol.toPrimitive]() { return ""; }
    static set [Symbol.toPrimitive](x) { }
}

//// [symbolDeclarationEmit11.js]
"use strict";
var _a;
class C {
    static [(_a = Symbol.iterator, Symbol.isConcatSpreadable)]() { }
    static get [Symbol.toPrimitive]() { return ""; }
    static set [Symbol.toPrimitive](x) { }
}
C[_a] = 0;


//// [symbolDeclarationEmit11.d.ts]
class C {
    static [Symbol.iterator]: number;
    static [Symbol.isConcatSpreadable](): void;
    static get [Symbol.toPrimitive](): string;
    static set [Symbol.toPrimitive](x: string);
}


//// [DtsFileErrors]


symbolDeclarationEmit11.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit11.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        static [Symbol.iterator]: number;
        static [Symbol.isConcatSpreadable](): void;
        static get [Symbol.toPrimitive](): string;
        static set [Symbol.toPrimitive](x: string);
    }
    
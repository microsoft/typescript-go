//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit3.ts] ////

//// [symbolDeclarationEmit3.ts]
class C {
    [Symbol.toPrimitive](x: number);
    [Symbol.toPrimitive](x: string);
    [Symbol.toPrimitive](x: any) { }
}

//// [symbolDeclarationEmit3.js]
"use strict";
class C {
    [Symbol.toPrimitive](x) { }
}


//// [symbolDeclarationEmit3.d.ts]
class C {
    [Symbol.toPrimitive](x: number): any;
    [Symbol.toPrimitive](x: string): any;
}


//// [DtsFileErrors]


symbolDeclarationEmit3.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit3.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.toPrimitive](x: number): any;
        [Symbol.toPrimitive](x: string): any;
    }
    
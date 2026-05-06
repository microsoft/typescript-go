//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit2.ts] ////

//// [symbolDeclarationEmit2.ts]
class C {
    [Symbol.toPrimitive] = "";
}

//// [symbolDeclarationEmit2.js]
"use strict";
var _a;
class C {
    constructor() {
        this[_a] = "";
    }
}
_a = Symbol.toPrimitive;


//// [symbolDeclarationEmit2.d.ts]
class C {
    [Symbol.toPrimitive]: string;
}


//// [DtsFileErrors]


symbolDeclarationEmit2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit2.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.toPrimitive]: string;
    }
    
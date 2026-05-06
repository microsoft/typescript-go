//// [tests/cases/conformance/classes/constructorDeclarations/constructorParameters/declarationEmitReadonly.ts] ////

//// [declarationEmitReadonly.ts]
class C {
    constructor(readonly x: number) {}
}

//// [declarationEmitReadonly.js]
"use strict";
class C {
    constructor(x) {
        this.x = x;
    }
}


//// [declarationEmitReadonly.d.ts]
class C {
    readonly x: number;
    constructor(x: number);
}


//// [DtsFileErrors]


declarationEmitReadonly.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitReadonly.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        readonly x: number;
        constructor(x: number);
    }
    
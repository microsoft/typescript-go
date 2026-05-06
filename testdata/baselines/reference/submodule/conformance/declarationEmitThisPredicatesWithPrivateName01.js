//// [tests/cases/conformance/declarationEmit/typePredicates/declarationEmitThisPredicatesWithPrivateName01.ts] ////

//// [declarationEmitThisPredicatesWithPrivateName01.ts]
export class C {
    m(): this is D {
        return this instanceof D;
    }
}

class D extends C {
}

//// [declarationEmitThisPredicatesWithPrivateName01.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = void 0;
class C {
    m() {
        return this instanceof D;
    }
}
exports.C = C;
class D extends C {
}


//// [declarationEmitThisPredicatesWithPrivateName01.d.ts]
export class C {
    m(): this is D;
}
class D extends C {
}
export {};


//// [DtsFileErrors]


declarationEmitThisPredicatesWithPrivateName01.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitThisPredicatesWithPrivateName01.d.ts (1 errors) ====
    export class C {
        m(): this is D;
    }
    class D extends C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    export {};
    
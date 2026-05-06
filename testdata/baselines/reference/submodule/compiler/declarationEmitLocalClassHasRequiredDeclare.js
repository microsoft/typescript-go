//// [tests/cases/compiler/declarationEmitLocalClassHasRequiredDeclare.ts] ////

//// [declarationEmitLocalClassHasRequiredDeclare.ts]
export declare namespace A {
    namespace X { }
}

class X { }

export class A {
    static X = X;
}

export declare namespace Y {

}

export class Y { }

//// [declarationEmitLocalClassHasRequiredDeclare.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Y = exports.A = void 0;
class X {
}
class A {
}
exports.A = A;
A.X = X;
class Y {
}
exports.Y = Y;


//// [declarationEmitLocalClassHasRequiredDeclare.d.ts]
export namespace A {
    namespace X { }
}
class X {
}
export class A {
    static X: typeof X;
}
export namespace Y {
}
export class Y {
}
export {};


//// [DtsFileErrors]


declarationEmitLocalClassHasRequiredDeclare.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitLocalClassHasRequiredDeclare.d.ts (1 errors) ====
    export namespace A {
        namespace X { }
    }
    class X {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    export class A {
        static X: typeof X;
    }
    export namespace Y {
    }
    export class Y {
    }
    export {};
    
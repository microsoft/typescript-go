//// [tests/cases/compiler/declarationEmitPrivateReadonlyLiterals.ts] ////

//// [declarationEmitPrivateReadonlyLiterals.ts]
class Foo {
    private static readonly A = "a";
    private readonly B = "b";
    private static readonly C = 42;
    private readonly D = 42;
}


//// [declarationEmitPrivateReadonlyLiterals.js]
"use strict";
class Foo {
    constructor() {
        this.B = "b";
        this.D = 42;
    }
}
Foo.A = "a";
Foo.C = 42;


//// [declarationEmitPrivateReadonlyLiterals.d.ts]
class Foo {
    private static readonly A;
    private readonly B;
    private static readonly C;
    private readonly D;
}


//// [DtsFileErrors]


declarationEmitPrivateReadonlyLiterals.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitPrivateReadonlyLiterals.d.ts (1 errors) ====
    class Foo {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        private static readonly A;
        private readonly B;
        private static readonly C;
        private readonly D;
    }
    
//// [tests/cases/compiler/classExpressionInClassStaticDeclarations.ts] ////

//// [classExpressionInClassStaticDeclarations.ts]
class C {
    static D = class extends C {};
}

//// [classExpressionInClassStaticDeclarations.js]
"use strict";
class C {
}
C.D = class extends C {
};


//// [classExpressionInClassStaticDeclarations.d.ts]
class C {
    static D: {
        new (): {};
        D: /*elided*/ any;
    };
}


//// [DtsFileErrors]


classExpressionInClassStaticDeclarations.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== classExpressionInClassStaticDeclarations.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        static D: {
            new (): {};
            D: /*elided*/ any;
        };
    }
    
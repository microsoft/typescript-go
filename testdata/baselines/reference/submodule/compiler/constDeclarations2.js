//// [tests/cases/compiler/constDeclarations2.ts] ////

//// [constDeclarations2.ts]
// No error
namespace M {
    export const c1 = false;
    export const c2: number = 23;
    export const c3 = 0, c4 :string = "", c5 = null;
}


//// [constDeclarations2.js]
"use strict";
// No error
var M;
(function (M) {
    M.c1 = false;
    M.c2 = 23;
    M.c3 = 0, M.c4 = "", M.c5 = null;
})(M || (M = {}));


//// [constDeclarations2.d.ts]
namespace M {
    const c1 = false;
    const c2: number;
    const c3 = 0, c4: string, c5: any;
}


//// [DtsFileErrors]


constDeclarations2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== constDeclarations2.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        const c1 = false;
        const c2: number;
        const c3 = 0, c4: string, c5: any;
    }
    
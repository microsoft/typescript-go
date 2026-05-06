//// [tests/cases/compiler/constEnumDeclarations.ts] ////

//// [constEnumDeclarations.ts]
const enum E {
    A = 1,
    B = 2,
    C = A | B
}

const enum E2 {
    A = 1,
    B,
    C
}

//// [constEnumDeclarations.js]
"use strict";


//// [constEnumDeclarations.d.ts]
const enum E {
    A = 1,
    B = 2,
    C = 3
}
const enum E2 {
    A = 1,
    B = 2,
    C = 3
}


//// [DtsFileErrors]


constEnumDeclarations.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== constEnumDeclarations.d.ts (1 errors) ====
    const enum E {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        A = 1,
        B = 2,
        C = 3
    }
    const enum E2 {
        A = 1,
        B = 2,
        C = 3
    }
    
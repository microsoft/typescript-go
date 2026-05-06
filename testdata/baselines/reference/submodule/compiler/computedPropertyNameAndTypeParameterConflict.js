//// [tests/cases/compiler/computedPropertyNameAndTypeParameterConflict.ts] ////

//// [computedPropertyNameAndTypeParameterConflict.ts]
declare const O: unique symbol;
declare class Bar<O> {
    [O]: number;
}



//// [computedPropertyNameAndTypeParameterConflict.js]
"use strict";


//// [computedPropertyNameAndTypeParameterConflict.d.ts]
const O: unique symbol;
class Bar<O> {
    [O]: number;
}


//// [DtsFileErrors]


computedPropertyNameAndTypeParameterConflict.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== computedPropertyNameAndTypeParameterConflict.d.ts (1 errors) ====
    const O: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    class Bar<O> {
        [O]: number;
    }
    
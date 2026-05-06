//// [tests/cases/compiler/requiredInitializedParameter4.ts] ////

//// [requiredInitializedParameter4.ts]
class C1 {
    method(a = 0, b) { }
}

//// [requiredInitializedParameter4.js]
"use strict";
class C1 {
    method(a = 0, b) { }
}


//// [requiredInitializedParameter4.d.ts]
class C1 {
    method(a: number, b: any): void;
}


//// [DtsFileErrors]


requiredInitializedParameter4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== requiredInitializedParameter4.d.ts (1 errors) ====
    class C1 {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        method(a: number, b: any): void;
    }
    
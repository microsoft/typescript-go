//// [tests/cases/compiler/requiredInitializedParameter3.ts] ////

//// [requiredInitializedParameter3.ts]
interface I1 {
    method();
}

class C1 implements I1 {
    method(a = 0, b?) { }
}

//// [requiredInitializedParameter3.js]
"use strict";
class C1 {
    method(a = 0, b) { }
}


//// [requiredInitializedParameter3.d.ts]
interface I1 {
    method(): any;
}
class C1 implements I1 {
    method(a?: number, b?: any): void;
}


//// [DtsFileErrors]


requiredInitializedParameter3.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== requiredInitializedParameter3.d.ts (1 errors) ====
    interface I1 {
        method(): any;
    }
    class C1 implements I1 {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        method(a?: number, b?: any): void;
    }
    
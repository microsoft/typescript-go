//// [tests/cases/compiler/declFileForClassWithPrivateOverloadedFunction.ts] ////

//// [declFileForClassWithPrivateOverloadedFunction.ts]
class C {
    private foo(x: number);
    private foo(x: string);
    private foo(x: any) { }
}

//// [declFileForClassWithPrivateOverloadedFunction.js]
"use strict";
class C {
    foo(x) { }
}


//// [declFileForClassWithPrivateOverloadedFunction.d.ts]
class C {
    private foo;
}


//// [DtsFileErrors]


declFileForClassWithPrivateOverloadedFunction.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileForClassWithPrivateOverloadedFunction.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        private foo;
    }
    
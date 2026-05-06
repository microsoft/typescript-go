//// [tests/cases/compiler/declFileForFunctionTypeAsTypeParameter.ts] ////

//// [declFileForFunctionTypeAsTypeParameter.ts]
class X<T> {
}
class C extends X<() => number> {
}
interface I extends X<() => number> {
}



//// [declFileForFunctionTypeAsTypeParameter.js]
"use strict";
class X {
}
class C extends X {
}


//// [declFileForFunctionTypeAsTypeParameter.d.ts]
class X<T> {
}
class C extends X<() => number> {
}
interface I extends X<() => number> {
}


//// [DtsFileErrors]


declFileForFunctionTypeAsTypeParameter.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileForFunctionTypeAsTypeParameter.d.ts (1 errors) ====
    class X<T> {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    class C extends X<() => number> {
    }
    interface I extends X<() => number> {
    }
    
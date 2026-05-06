//// [tests/cases/compiler/declFileClassExtendsNull.ts] ////

//// [declFileClassExtendsNull.ts]
class ExtendsNull extends null {
}

//// [declFileClassExtendsNull.js]
"use strict";
class ExtendsNull extends null {
}


//// [declFileClassExtendsNull.d.ts]
class ExtendsNull extends null {
}


//// [DtsFileErrors]


declFileClassExtendsNull.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileClassExtendsNull.d.ts (1 errors) ====
    class ExtendsNull extends null {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    
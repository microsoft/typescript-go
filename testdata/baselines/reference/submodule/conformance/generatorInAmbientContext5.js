//// [tests/cases/conformance/es6/yieldExpressions/generatorInAmbientContext5.ts] ////

//// [generatorInAmbientContext5.ts]
class C {
    *generator(): any { }
}

//// [generatorInAmbientContext5.js]
"use strict";
class C {
    *generator() { }
}


//// [generatorInAmbientContext5.d.ts]
class C {
    generator(): any;
}


//// [DtsFileErrors]


generatorInAmbientContext5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== generatorInAmbientContext5.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        generator(): any;
    }
    
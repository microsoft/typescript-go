//// [tests/cases/compiler/enumDecl1.ts] ////

//// [enumDecl1.ts]
declare namespace mAmbient {
    enum e {
        x,
        y,
        z
    }
}


//// [enumDecl1.js]
"use strict";


//// [enumDecl1.d.ts]
namespace mAmbient {
    enum e {
        x,
        y,
        z
    }
}


//// [DtsFileErrors]


enumDecl1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== enumDecl1.d.ts (1 errors) ====
    namespace mAmbient {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        enum e {
            x,
            y,
            z
        }
    }
    
//// [tests/cases/compiler/declFileModuleAssignmentInObjectLiteralProperty.ts] ////

//// [declFileModuleAssignmentInObjectLiteralProperty.ts]
namespace m1 {
    export class c {
    }
}
var d = {
    m1: { m: m1 },
    m2: { c: m1.c },
};

//// [declFileModuleAssignmentInObjectLiteralProperty.js]
"use strict";
var m1;
(function (m1) {
    class c {
    }
    m1.c = c;
})(m1 || (m1 = {}));
var d = {
    m1: { m: m1 },
    m2: { c: m1.c },
};


//// [declFileModuleAssignmentInObjectLiteralProperty.d.ts]
namespace m1 {
    class c {
    }
}
var d: {
    m1: {
        m: typeof m1;
    };
    m2: {
        c: typeof m1.c;
    };
};


//// [DtsFileErrors]


declFileModuleAssignmentInObjectLiteralProperty.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileModuleAssignmentInObjectLiteralProperty.d.ts (1 errors) ====
    namespace m1 {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        class c {
        }
    }
    var d: {
        m1: {
            m: typeof m1;
        };
        m2: {
            c: typeof m1.c;
        };
    };
    
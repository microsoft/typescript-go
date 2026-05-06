//// [tests/cases/compiler/declFileImportChainInExportAssignment.ts] ////

//// [declFileImportChainInExportAssignment.ts]
namespace m {
    export namespace c {
        export class c {
        }
    }
}
import a = m.c;
import b = a;
export = b;

//// [declFileImportChainInExportAssignment.js]
"use strict";
var m;
(function (m) {
    let c;
    (function (c_1) {
        class c {
        }
        c_1.c = c;
    })(c = m.c || (m.c = {}));
})(m || (m = {}));
var a = m.c;
var b = a;
module.exports = b;


//// [declFileImportChainInExportAssignment.d.ts]
namespace m {
    namespace c {
        class c {
        }
    }
}
import a = m.c;
import b = a;
export = b;


//// [DtsFileErrors]


declFileImportChainInExportAssignment.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileImportChainInExportAssignment.d.ts (1 errors) ====
    namespace m {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace c {
            class c {
            }
        }
    }
    import a = m.c;
    import b = a;
    export = b;
    
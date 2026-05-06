//// [tests/cases/compiler/declarationEmitImportInExportAssignmentModule.ts] ////

//// [declarationEmitImportInExportAssignmentModule.ts]
namespace m {
    export namespace c {
        export class c {
        }
    }
    import x = c;
    export var a: typeof x;
}
export = m;

//// [declarationEmitImportInExportAssignmentModule.js]
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
module.exports = m;


//// [declarationEmitImportInExportAssignmentModule.d.ts]
namespace m {
    namespace c {
        class c {
        }
    }
    import x = c;
    var a: typeof x;
}
export = m;


//// [DtsFileErrors]


declarationEmitImportInExportAssignmentModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitImportInExportAssignmentModule.d.ts (1 errors) ====
    namespace m {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace c {
            class c {
            }
        }
        import x = c;
        var a: typeof x;
    }
    export = m;
    
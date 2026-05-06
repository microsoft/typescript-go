//// [tests/cases/compiler/declFileWithInternalModuleNameConflictsInExtendsClause1.ts] ////

//// [declFileWithInternalModuleNameConflictsInExtendsClause1.ts]
namespace X.A.C {
    export interface Z {
    }
}
namespace X.A.B.C {
    namespace A {
    }
    export class W implements X.A.C.Z { // This needs to be referred as X.A.C.Z as A has conflict
    }
}

//// [declFileWithInternalModuleNameConflictsInExtendsClause1.js]
"use strict";
var X;
(function (X) {
    var A;
    (function (A) {
        var B;
        (function (B) {
            var C;
            (function (C) {
                class W {
                }
                C.W = W;
            })(C = B.C || (B.C = {}));
        })(B = A.B || (A.B = {}));
    })(A = X.A || (X.A = {}));
})(X || (X = {}));


//// [declFileWithInternalModuleNameConflictsInExtendsClause1.d.ts]
namespace X.A.C {
    interface Z {
    }
}
namespace X.A.B.C {
    class W implements X.A.C.Z {
    }
}


//// [DtsFileErrors]


declFileWithInternalModuleNameConflictsInExtendsClause1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileWithInternalModuleNameConflictsInExtendsClause1.d.ts (1 errors) ====
    namespace X.A.C {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Z {
        }
    }
    namespace X.A.B.C {
        class W implements X.A.C.Z {
        }
    }
    
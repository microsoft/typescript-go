//// [tests/cases/compiler/declFileWithInternalModuleNameConflictsInExtendsClause2.ts] ////

//// [declFileWithInternalModuleNameConflictsInExtendsClause2.ts]
namespace X.A.C {
    export interface Z {
    }
}
namespace X.A.B.C {
    export class W implements A.C.Z { // This can refer to it as A.C.Z
    }
}

namespace X.A.B.C {
    namespace A {
    }
}

//// [declFileWithInternalModuleNameConflictsInExtendsClause2.js]
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


//// [declFileWithInternalModuleNameConflictsInExtendsClause2.d.ts]
namespace X.A.C {
    interface Z {
    }
}
namespace X.A.B.C {
    class W implements A.C.Z {
    }
}
namespace X.A.B.C {
}


//// [DtsFileErrors]


declFileWithInternalModuleNameConflictsInExtendsClause2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileWithInternalModuleNameConflictsInExtendsClause2.d.ts (1 errors) ====
    namespace X.A.C {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Z {
        }
    }
    namespace X.A.B.C {
        class W implements A.C.Z {
        }
    }
    namespace X.A.B.C {
    }
    
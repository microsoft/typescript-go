//// [tests/cases/compiler/declFileWithInternalModuleNameConflictsInExtendsClause3.ts] ////

//// [declFileWithInternalModuleNameConflictsInExtendsClause3.ts]
namespace X.A.C {
    export interface Z {
    }
}
namespace X.A.B.C {
    export class W implements X.A.C.Z { // This needs to be referred as X.A.C.Z as A has conflict
    }
}

namespace X.A.B.C {
    export namespace A {
    }
}

//// [declFileWithInternalModuleNameConflictsInExtendsClause3.js]
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


//// [declFileWithInternalModuleNameConflictsInExtendsClause3.d.ts]
namespace X.A.C {
    interface Z {
    }
}
namespace X.A.B.C {
    class W implements X.A.C.Z {
    }
}
namespace X.A.B.C {
    namespace A {
    }
}


//// [DtsFileErrors]


declFileWithInternalModuleNameConflictsInExtendsClause3.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileWithInternalModuleNameConflictsInExtendsClause3.d.ts (1 errors) ====
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
    namespace X.A.B.C {
        namespace A {
        }
    }
    
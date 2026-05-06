//// [tests/cases/compiler/declFileModuleContinuation.ts] ////

//// [declFileModuleContinuation.ts]
namespace A.C {
    export interface Z {
    }
}

namespace A.B.C {
    export class W implements A.C.Z {
    }
}

//// [declFileModuleContinuation.js]
"use strict";
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
})(A || (A = {}));


//// [declFileModuleContinuation.d.ts]
namespace A.C {
    interface Z {
    }
}
namespace A.B.C {
    class W implements A.C.Z {
    }
}


//// [DtsFileErrors]


declFileModuleContinuation.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileModuleContinuation.d.ts (1 errors) ====
    namespace A.C {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Z {
        }
    }
    namespace A.B.C {
        class W implements A.C.Z {
        }
    }
    